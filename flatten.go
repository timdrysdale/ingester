package ingester

import (
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
	"github.com/timdrysdale/pdfpagedata"
)

type FlattenTask struct {
	Path      string
	PageCount int
	Data      pdfpagedata.PageData
}

func FlattenNewPapers(exam string) error {

	//assume someone hits a button to ask us to do this ...

	tasks := []FlattenTask{}

	receipts, err := gradexpath.GetFileList(gradexpath.AcceptedReceipts(exam))
	if err != nil {
		return err
	}

	for _, receipt := range receipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)
		if err != nil {
			continue
			// TODO need to flag to user as we shouldn't fail to read a learn receipt here
		}

		pdfPath, err := GetPdfPath(sub.Filename, gradexpath.AcceptedPapers(exam))

		if err != nil {
			continue
			// TODO need to flag to user as we shouldn't fail to find a PDF here
		}

		count, err := countPages(pdfPath)

		if err != nil {
			continue
			// TODO need to flag to user as we shouldn't fail to count pages here
		}
		shortDate, err := GetShortLearnDate(sub)
		if err != nil {
			continue
			// TODO need to flag to user as we shouldn't fail to read sub here
		}

		pagedata := pdfpagedata.PageData{
			Exam: pdfpagedata.ExamDetails{
				CourseCode: sub.Assignment,
				Date:       shortDate,
			},
			Author: pdfpagedata.AuthorDetails{
				Identity: sub.Matriculation,
			},
		}
		//TODO fill this out a bit more...
		tasks = append(tasks, FlattenTask{Path: pdfPath, PageCount: count, Data: pagedata})
	}

	//parsesvg.PrettyPrintStruct(tasks)

	return nil

}
