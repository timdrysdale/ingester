package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
)

func CollectFilesFrom(path string) {

}

func TestAddBars(t *testing.T) {
	verbose := false

	//collectOutputs := true

	gradexpath.SetTesting()

	root := gradexpath.Root()

	if root != "./tmp-delete-me" {
		t.Errorf("test root set up wrong %s", root)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
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

	//>>>>>>>>>>>>>>>>>>>>>>>>> INGEST >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

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

	//>>>>>>>>>>>>>>>>>>>>>>>>> VALIDATE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

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

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP FOR FLATTEN/RENAME  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// Now we test Flatten
	//copy in the identity database
	src := "./test-fs/etc/identity/identity.csv"
	dest := "./tmp-delete-me/etc/identity/identity.csv"
	err = gradexpath.Copy(src, dest)
	assert.NoError(t, err)
	_, err = os.Stat(dest)

	//>>>>>>>>>>>>>>>>>>>>>>>>> FLATTEN/RENAME  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
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

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP FOR OVERLAY (via ADDBARS) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	templateFiles, err = gradexpath.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(gradexpath.OverlayTemplate(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	mch := make(chan chmsg.MessageInfo)

	closed := make(chan struct{})
	defer close(closed)
	go func() {
		for {
			select {
			case <-closed:
				break
			case msg := <-mch:
				if verbose {
					fmt.Printf("MC:%s\n", msg.Message)
				}
			}

		}
	}()

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD MARKBAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	marker := "tddrysdale"
	err = AddMarkBar(exam, marker, mch)
	assert.NoError(t, err)

	expectedMarker1Pdf := []string{
		"Practice Exam Drop Box-B999995-maTDD.pdf",
		"Practice Exam Drop Box-B999997-maTDD.pdf",
		"Practice Exam Drop Box-B999998-maTDD.pdf",
		"Practice Exam Drop Box-B999999-maTDD.pdf",
	}

	readyPdf, err := gradexpath.GetFileList(gradexpath.MarkerReady(exam, marker))

	assert.NoError(t, err)

	assert.Equal(t, len(expectedMarker1Pdf), len(readyPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedMarker1Pdf, readyPdf))

	pds, err = pdfpagedata.GetPageDataFromFile(readyPdf[0])
	assert.NoError(t, err)
	pd = pds[0]
	assert.Equal(t, pd[0].Questions[0].Name, "marking")

	for _, file := range readyPdf[0:2] {
		destination := filepath.Join(gradexpath.ModerateActive(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}
	for _, file := range readyPdf[2:4] {
		destination := filepath.Join(gradexpath.ModerateInActive(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD ACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	moderator := "ABC"
	err = AddModerateActiveBar(exam, moderator, mch)
	assert.NoError(t, err)

	expectedActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
	}

	activePdf, err := gradexpath.GetFileList(gradexpath.ModeratorReady(exam, moderator))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActive), len(activePdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedActive, activePdf))

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD INACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = AddModerateInActiveBar(exam, mch)
	assert.NoError(t, err)

	expectedInActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	inActivePdf, err := gradexpath.GetFileList(gradexpath.ModeratedInActiveBack(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInActive), len(inActivePdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedInActive, inActivePdf))

	// copy files to common area (as if have processed them - not checked here)

	for _, file := range activePdf {
		destination := filepath.Join(gradexpath.ModeratedReady(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	for _, file := range inActivePdf {
		destination := filepath.Join(gradexpath.ModeratedReady(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	expectedModeratedReadyPdf := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	moderatedReadyPdf, err := gradexpath.GetFileList(gradexpath.ModeratedReady(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedModeratedReadyPdf), len(moderatedReadyPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedModeratedReadyPdf, moderatedReadyPdf))

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD CHECK BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	checker := "LD"

	err = AddCheckBar(exam, checker, mch)
	assert.NoError(t, err)
	expectedChecked := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX-cLD.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX-cLD.pdf",
	}

	checkedPdf, err := gradexpath.GetFileList(gradexpath.CheckerReady(exam, checker))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedChecked), len(checkedPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedChecked, checkedPdf))

}
