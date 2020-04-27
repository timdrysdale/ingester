package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/timdrysdale/anon"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
	"github.com/timdrysdale/parsesvg"
	"github.com/timdrysdale/pdfcomment"
	"github.com/timdrysdale/pdfpagedata"
	"github.com/timdrysdale/pool"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

type FlattenTask struct {
	InputPath  string
	PageCount  int
	Data       pdfpagedata.PageData
	OutputPath string
}

func FlattenNewPapers(exam string) error {

	//assume someone hits a button to ask us to do this ...

	// we'll use this same set of procDetails for flattens that we do in this batch
	// that means we can use the uuid to map the processing in graphviz later, for example
	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuid := UUIDBytes.String()
	if err != nil {
		uuid = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuid,
		Previous: "none",
		UnixTime: time.Now().UnixNano(),
		Name:     "flatten",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0,
	}

	// load our identity database
	identity, err := anon.New(gradexpath.IdentityCSV())
	if err != nil {
		return err
	}
	fmt.Println("Got Identity DB")
	flattenTasks := []FlattenTask{}

	receipts, err := gradexpath.GetFileList(gradexpath.AcceptedReceipts(exam))
	if err != nil {
		return err
	}

	for _, receipt := range receipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)
		if err != nil {
			fmt.Printf("couldn't parse receipt %s because %v", receipt, err)
			continue
			// TODO need to flag to user as we shouldn't fail to read a learn receipt here
		}

		pdfPath, err := GetPdfPath(sub.Filename, gradexpath.AcceptedPapers(exam))

		if err != nil {
			fmt.Printf("couldn't get PDF filename for %s because %v\n", sub.Filename, err)
			continue
			// TODO need to flag to user as we shouldn't fail to find a PDF here
		}

		count, err := countPages(pdfPath)

		if err != nil {
			fmt.Printf("couldn't countPages for %s because %v\n", pdfPath, err)
			continue
			// TODO need to flag to user as we shouldn't fail to count pages here
		}
		shortDate, err := GetShortLearnDate(sub)
		if err != nil {
			fmt.Printf("couldn't get shortlearndate for %s because %v\n", receipt, err)
			continue
			// TODO need to flag to user as we shouldn't fail to read sub here
		}

		//TODO If identity not known, need to flag to user, and not process paper just now

		anonymousIdentity, err := identity.GetAnonymous(sub.Matriculation)
		if err != nil {
			fmt.Printf("couldn't get identity for for %s because %v\n", sub.Matriculation, err)
			continue
			// TODO need to flag to user as we should have all IDs in our dictionary
		}

		pagedata := pdfpagedata.PageData{
			Exam: pdfpagedata.ExamDetails{
				CourseCode: sub.Assignment,
				Date:       shortDate,
			},
			Author: pdfpagedata.AuthorDetails{
				Anonymous: anonymousIdentity,
			},
			Processing: []pdfpagedata.ProcessingDetails{procDetails},
		}

		renamedBase := gradexpath.GetAnonymousFileName(sub.Assignment, anonymousIdentity)
		outputPath := filepath.Join(gradexpath.AnonymousPapers(sub.Assignment), renamedBase)

		flattenTasks = append(flattenTasks, FlattenTask{
			InputPath:  pdfPath,
			OutputPath: outputPath,
			PageCount:  count,
			Data:       pagedata})
	}

	// now process the files
	N := len(flattenTasks)

	pcChan := make(chan int, N)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		inputPath := flattenTasks[i].InputPath
		outputPath := flattenTasks[i].OutputPath
		pd := flattenTasks[i].Data

		newtask := pool.NewTask(func() error {
			pc, err := FlattenOnePdf(inputPath, outputPath, pd)
			pcChan <- pc
			return err
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))

	closed := make(chan struct{})

	//	h := thist.NewHist(nil, "Page count", "fixed", 10, false)
	//
	//	go func() {
	//	LOOP:
	//		for {
	//			select {
	//			case pc := <-pcChan:
	//				h.Update(float64(pc))
	//				fmt.Println(h.Draw())
	//			case <-closed:
	//				break LOOP
	//			}
	//		}
	//	}()
	//
	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			fmt.Println(task.Err)
			numErrors++
		}
	}
	close(closed)

	return nil

}

func FlattenOnePdf(inputPath, outputPath string, pageData pdfpagedata.PageData) (int, error) {

	if strings.ToLower(filepath.Ext(inputPath)) != ".pdf" {
		return 0, errors.New(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
	}

	// need page count to find the jpeg files again later
	numPages, err := countPages(inputPath)

	// render to images
	jpegPath := gradexpath.AcceptedPaperImages(pageData.Exam.CourseCode)

	suffix := filepath.Ext(inputPath)
	basename := strings.TrimSuffix(filepath.Base(inputPath), suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("FLATTEN Can't open pdf")
		return 0, err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		fmt.Println("FLATTEN Can't read pdf")
		return 0, err
	}

	comments, err := pdfcomment.GetComments(pdfReader)

	f.Close()

	err = convertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
	if err != nil {
		return 0, err
	}

	// convert images to individual pdfs, with form overlay

	pagePath := gradexpath.AcceptedPaperPages(pageData.Exam.CourseCode)
	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)
	fmt.Printf("Page file: %s\n", pageFileOption)
	mergePaths := []string{}

	pageData.Page.Of = numPages

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		//TODO select Layout to suit landscape or portrait
		svgLayoutPath := gradexpath.FlattenLayoutSVG()

		pageNumber := imgIdx - 1

		pageData.Page.Number = pageNumber + 1
		pageData.Page.Filename = filepath.Base(pageFilename)

		var pageUUIDBytes uuid.UUID

		pageUUIDBytes, err = uuid.NewRandom()

		pageUUID := pageUUIDBytes.String()

		if err != nil {
			pageUUID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		pageData.Page.UUID = pageUUID

		contents := parsesvg.SpreadContents{
			SvgLayoutPath:         svgLayoutPath,
			SpreadName:            "flatten",
			PreviousImagePath:     previousImagePath,
			PageNumber:            pageNumber,
			PdfOutputPath:         pageFilename,
			Comments:              comments,
			PageData:              pageData,
			TemplatePathsRelative: true,
		}

		err := parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			fmt.Println(err)
			return 0, err

		}

		mergePaths = append(mergePaths, pageFilename)
	}
	fmt.Printf("MergePaths: %v\v", mergePaths)
	err = mergePdf(mergePaths, outputPath)
	if err != nil {
		fmt.Printf("MERGE: %v", err)
		return 0, err
	}

	return numPages, nil

}
