package ingester

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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

//check we can move files without adjusting the modification time
func TestStructFileMod(t *testing.T) {

	mch := make(chan chmsg.MessageInfo)
	g, err := NewIngester("./tmp-delete-me", mch)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	d1 := []byte("Gradex Testing\n")
	basepath := filepath.Join(g.Root(), "tmp")
	err = g.EnsureDir(basepath)
	assert.NoError(t, err)
	testPath := filepath.Join(basepath, "test.txt")
	err = ioutil.WriteFile(testPath, d1, 0755)
	assert.NoError(t, err)
	err = os.Chmod(testPath, 0755)
	assert.NoError(t, err)

	info, err := os.Stat(testPath)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.NotEqual(t, info.ModTime(), time.Now())

	newPath := filepath.Join(g.Root(), "tmp", "new.txt")
	err = os.Rename(testPath, newPath)
	infoNew, err := os.Stat(newPath)

	assert.NoError(t, err)
	assert.Equal(t, info.ModTime(), infoNew.ModTime())
}

func TestNewFileMove(t *testing.T) {
	mch := make(chan chmsg.MessageInfo)
	g, err := NewIngester("./tmp-delete-me", mch)

	d0 := []byte("Gradex Testing\n")
	basepath := filepath.Join(g.Root(), "tmp")
	err = g.EnsureDir(basepath)
	assert.NoError(t, err)

	test0 := filepath.Join(basepath, "test0.txt")
	err = ioutil.WriteFile(test0, d0, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test0, 0755)
	assert.NoError(t, err)
	info0, err := os.Stat(test0)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test1 := filepath.Join(basepath, "test1.txt")
	d1 := []byte("XXXX\n")
	err = ioutil.WriteFile(test1, d1, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test1, 0755)
	assert.NoError(t, err)
	info1, err := os.Stat(test1)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test2 := filepath.Join(basepath, "test2.txt")
	d2 := []byte("YYYY\n")
	err = ioutil.WriteFile(test2, d2, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test2, 0755)
	assert.NoError(t, err)
	info2, err := os.Stat(test2)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test3 := filepath.Join(basepath, "test3.txt")
	d3 := []byte("ZZZZ\n")
	err = ioutil.WriteFile(test3, d3, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test3, 0755)
	assert.NoError(t, err)
	info3, err := os.Stat(test3)
	assert.NoError(t, err)

	// check file modtimes

	assert.True(t, info3.ModTime().After(info2.ModTime()))
	assert.True(t, info2.ModTime().After(info1.ModTime()))
	assert.True(t, info1.ModTime().After(info0.ModTime()))

	//should move
	err = g.MoveIfNewerThanDestination(test1, test0)
	assert.NoError(t, err)

	//should NOT move - but throw no error
	err = g.MoveIfNewerThanDestination(test2, test3)
	assert.NoError(t, err)

	info0new, err := os.Stat(test0)
	assert.NoError(t, err)
	_, err = os.Stat(test1)
	assert.Error(t, err) // ERROR should have moved!
	_, err = os.Stat(test2)
	assert.NoError(t, err) // no error - should NOT have moved
	info3new, err := os.Stat(test3)
	assert.NoError(t, err)

	if !info0new.ModTime().After(info0.ModTime()) {
		t.Error("first file mod time should have changed")
	}

	if !info3new.ModTime().Equal(info3.ModTime()) {
		t.Error("last file mod time should NOT have changed")
	}

	c0, err := ioutil.ReadFile(test0)
	assert.NoError(t, err)
	c2, err := ioutil.ReadFile(test2)
	assert.NoError(t, err)
	c3, err := ioutil.ReadFile(test3)
	assert.NoError(t, err)

	assert.Equal(t, c0, d1) //content changed
	assert.Equal(t, c2, d2)
	assert.Equal(t, c3, d3) //content not changed

}
