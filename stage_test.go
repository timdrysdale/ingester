package ingester

import (
	"os"
	"testing"

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

	gradexpath.Copy("./test/*", gradexpath.Ingest())

}
