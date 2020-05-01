package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
)

func TestStageArchive(t *testing.T) {
	verbose := true

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
	g, err := NewIngester("./tmp-delete-me", mch)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	// don't use GetRoot() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	file := "./test-zip/test.zip"
	destination := filepath.Join(g.Ingest(), filepath.Base(file))
	err = g.Copy(file, destination)
	assert.NoError(t, err)

	ingestfiles, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, g.CopyIsComplete([]string{file}, ingestfiles))

	g.StageFromIngest()

	expectedRejects, err := g.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, g.CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := g.GetFileList("./expected/temp-txt")
	assert.NoError(t, err)

	actualTxt, err := g.GetFileList(g.TempTXT())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, g.CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := g.GetFileList("./expected/temp-pdf")
	assert.NoError(t, err)

	actualPdf, err := g.GetFileList(g.TempPDF())
	assert.NoError(t, err)

	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, g.CopyIsComplete(expectedPdf, actualPdf))

}

func TestStageUnModified(t *testing.T) {
	verbose := true
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

	g, err := NewIngester("./tmp-delete-me", mch)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	// don't use GetRoot() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	testfiles, err := g.GetFileList("./test")

	assert.NoError(t, err)

	//fmt.Println(testfiles)

	for _, file := range testfiles {
		destination := filepath.Join(g.Ingest(), filepath.Base(file))
		err := g.Copy(file, destination)
		assert.NoError(t, err)

	}

	ingestfiles, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, g.CopyIsComplete(testfiles, ingestfiles))

	g.StageFromIngest()

	expectedRejects, err := g.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, g.CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := g.GetFileList("./expected/temp-txt")
	assert.NoError(t, err)

	actualTxt, err := g.GetFileList(g.TempTXT())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, g.CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := g.GetFileList("./expected/temp-pdf")
	assert.NoError(t, err)

	actualPdf, err := g.GetFileList(g.TempPDF())
	assert.NoError(t, err)

	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, g.CopyIsComplete(expectedPdf, actualPdf))

}
