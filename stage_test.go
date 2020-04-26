package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradexpath"
)

var rejectedFiles = []string{
	"bar.jpg",
	"foo.doc",
	"Practice Exam Drop Box_s00000002_attempt_2020-04-22-10-43-23_my exam.doc",
	"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam - Copy (copy).jpg",
	"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam - Copy.jpg",
	"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam.jpg",
}

func TestStage(t *testing.T) {

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

	//fmt.Println(testfiles)

	for _, file := range testfiles {
		destination := filepath.Join(gradexpath.Ingest(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)

	}

	//fmt.Println(gradexpath.Ingest())
	ingestfiles, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	fmt.Println(len(ingestfiles))

	assert.True(t, gradexpath.CopyIsComplete(testfiles, ingestfiles))

	StageFromIngest()

}
