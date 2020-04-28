package ingester

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
)

func TestFlatten(t *testing.T) {
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

}
