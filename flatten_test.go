package ingester

/*
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/pdfpagedata"
)

func TestFlatten(t *testing.T) {

	gp, err := New("./tmp-delete-me")
	assert.NoError(t, err)

	assert.Equal(gp.Root(), "./tmp-delete-me")

	// don't use GetRoot() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	gp.EnsureDirectoryStructure()

	testfiles, err := gp.GetFileList("./test")

	assert.NoError(t, err)

	for _, file := range testfiles {
		destination := filepath.Join(gp.Ingest(), filepath.Base(file))
		err := gp.Copy(file, destination)
		assert.NoError(t, err)
	}

	templateFiles, err := gp.GetFileList("./test-fs/etc/ingest/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(gp.IngestTemplate(), filepath.Base(file))
		err := gp.Copy(file, destination)
		assert.NoError(t, err)
	}

	ingestfiles, err := gp.GetFileList(gp.Ingest())
	assert.NoError(t, err)

	assert.True(t, gp.CopyIsComplete(testfiles, ingestfiles))

	gp.StageFromIngest()

	expectedRejects, err := gp.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := gp.GetFileList(gp.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, gp.CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := gp.GetFileList("./expected/temp-txt")
	assert.NoError(t, err)

	actualTxt, err := gp.GetFileList(gp.TempTxt())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gp.CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := gp.GetFileList("./expected/temp-pdf")
	assert.NoError(t, err)

	actualPdf, err := gp.GetFileList(gp.TempPdf())
	assert.NoError(t, err)

	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gp.CopyIsComplete(expectedPdf, actualPdf))

	assert.NoError(t, gp.ValidateNewPapers())

	exam := "Practice Exam Drop Box"

	actualPdf, err = gp.GetFileList(gp.AcceptedPapers(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gp.CopyIsComplete(expectedPdf, actualPdf))

	actualTxt, err = gp.GetFileList(gp.AcceptedReceipts(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gp.CopyIsComplete(expectedTxt, actualTxt))

	tempPdf, err := gp.GetFileList(gp.TempPdf())
	assert.NoError(t, err)
	assert.Equal(t, len(tempPdf), 0)

	tempTxt, err := gp.GetFileList(gp.TempTxt())
	assert.NoError(t, err)
	assert.Equal(t, len(tempTxt), 0)

	// Now we test Flatten

	//copy in the identity database
	src := "./test-fs/etc/identity/identity.csv"
	dest := "./tmp-delete-me/etc/identity/identity.csv"
	err = gp.Copy(src, dest)
	assert.NoError(t, err)
	_, err = os.Stat(dest)

	// do flatten
	err = gp.FlattenNewPapers("Practice Exam Drop Box")
	assert.NoError(t, err)

	// check files exist

	expectedAnonymousPdf := []string{
		"Practice Exam Drop Box-B999995.pdf",
		"Practice Exam Drop Box-B999997.pdf",
		"Practice Exam Drop Box-B999998.pdf",
		"Practice Exam Drop Box-B999999.pdf",
	}

	anonymousPdf, err := gp.GetFileList(gp.AnonymousPapers(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(anonymousPdf), len(expectedAnonymousPdf))

	assert.True(t, gp.CopyIsComplete(expectedAnonymousPdf, anonymousPdf))

	// check data extraction

	pds, err := pdfpagedata.GetPageDataFromFile(anonymousPdf[0])
	assert.NoError(t, err)
	pd := pds[0]
	assert.Equal(t, pd[0].Exam.CourseCode, "Practice Exam Drop Box")

}
*/
