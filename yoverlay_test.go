package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
)

func TestOverlay(t *testing.T) {

	gradexpath.SetTesting()

	root := gradexpath.Root()

	if root != "./tmp-delete-me" {
		t.Errorf("test root set up wrong %s", root)
	}

	// don't use GetRoot() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	EnsureDirectoryStructure()

	testfiles, err := gradexpath.GetFileList("./test")

	assert.NoError(t, err)

	for _, file := range testfiles {
		destination := filepath.Join(gradexpath.Ingest(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)

	}

	templateFiles, err := gradexpath.GetFileList("./test-fs/etc/ingest/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(gradexpath.IngestTemplate(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	ingestfiles, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	assert.True(t, gradexpath.CopyIsComplete(testfiles, ingestfiles))

	StageFromIngest()

	expectedRejects, err := gradexpath.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, gradexpath.CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := gradexpath.GetFileList("./expected/temp-txt")
	assert.NoError(t, err)

	actualTxt, err := gradexpath.GetFileList(gradexpath.TempTxt())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gradexpath.CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := gradexpath.GetFileList("./expected/temp-pdf")
	assert.NoError(t, err)

	actualPdf, err := gradexpath.GetFileList(gradexpath.TempPdf())
	assert.NoError(t, err)

	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gradexpath.CopyIsComplete(expectedPdf, actualPdf))

	assert.NoError(t, ValidateNewPapers())

	exam := "Practice Exam Drop Box"

	actualPdf, err = gradexpath.GetFileList(gradexpath.AcceptedPapers(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gradexpath.CopyIsComplete(expectedPdf, actualPdf))

	actualTxt, err = gradexpath.GetFileList(gradexpath.AcceptedReceipts(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gradexpath.CopyIsComplete(expectedTxt, actualTxt))

	tempPdf, err := gradexpath.GetFileList(gradexpath.TempPdf())
	assert.NoError(t, err)
	assert.Equal(t, len(tempPdf), 0)

	tempTxt, err := gradexpath.GetFileList(gradexpath.TempTxt())
	assert.NoError(t, err)
	assert.Equal(t, len(tempTxt), 0)

	// Now we test Flatten

	//copy in the identity database
	src := "./test-fs/etc/identity/identity.csv"
	dest := "./tmp-delete-me/etc/identity/identity.csv"
	err = gradexpath.Copy(src, dest)
	assert.NoError(t, err)
	_, err = os.Stat(dest)

	// do flatten
	err = FlattenNewPapers("Practice Exam Drop Box")
	assert.NoError(t, err)

	// check files exist

	expectedAnonymousPdf := []string{
		"Practice Exam Drop Box-B999995.pdf",
		"Practice Exam Drop Box-B999997.pdf",
		"Practice Exam Drop Box-B999998.pdf",
		"Practice Exam Drop Box-B999999.pdf",
	}

	anonymousPdf, err := gradexpath.GetFileList(gradexpath.AnonymousPapers(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(anonymousPdf), len(expectedAnonymousPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedAnonymousPdf, anonymousPdf))

	// check data extraction

	pds, err := pdfpagedata.GetPageDataFromFile(anonymousPdf[0])
	assert.NoError(t, err)
	pd := pds[0]
	assert.Equal(t, pd[0].Exam.CourseCode, "Practice Exam Drop Box")

	// Now we test Overlay

	templateFiles, err = gradexpath.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(gradexpath.OverlayTemplate(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	mch := make(chan chmsg.MessageInfo)

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "Overlay",
		TaskName:     "Mark",
	}

	closed := make(chan struct{})

	go func() {
		for {
			select {
			case <-closed:
				break
			case msg := <-mch:
				fmt.Printf("MC:%s\n", msg.Message)
			}

		}
	}()

	cm := chmsg.New(mc, mch, 100*time.Millisecond)

	var UUIDBytes uuid.UUID

	UUIDBytes, err = uuid.NewRandom()
	uuid := UUIDBytes.String()
	if err != nil {
		uuid = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuid,
		Previous: "", //dynamic / Not Implemented
		UnixTime: time.Now().UnixNano(),
		Name:     "MarkBars",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 1, //dynamic
	}

	oc := OverlayCommand{
		FromPath:          gradexpath.AnonymousPapers(exam),
		ToPath:            gradexpath.Export(),
		ExamName:          exam,
		TemplatePath:      gradexpath.OverlayLayoutSVG(),
		SpreadName:        "mark",
		ProcessingDetails: procDetails,
		Msg:               cm,
		PathDecoration:    gradexpath.Marker1,
	}

	err = OverlayPapers(oc)
	// check files exist

	expectedMarker1Pdf := []string{
		"Practice Exam Drop Box-B999995-ma1.pdf",
		"Practice Exam Drop Box-B999997-ma1.pdf",
		"Practice Exam Drop Box-B999998-ma1.pdf",
		"Practice Exam Drop Box-B999999-ma1.pdf",
	}

	exportPdf, err := gradexpath.GetFileList(gradexpath.Export())
	assert.NoError(t, err)

	assert.Equal(t, len(expectedMarker1Pdf), len(exportPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedMarker1Pdf, exportPdf))

}
