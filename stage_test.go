package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradexpath"
)

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

	fmt.Println(testfiles)

	for _, file := range testfiles {
		destination := filepath.Join(gradexpath.Ingest(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)

	}

	fmt.Println(gradexpath.Ingest())
	ingestfiles, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	fmt.Println(ingestfiles)

	StageFromIngest()

}
