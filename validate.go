package ingester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
)

// for ingest, we only need to find pairs of
// ingest receipt + matching pdf input
// we assume someone has replaced non-pdf with pdf
// so we look for version of file with pdf suffix
func handleIngestLearnReceipt(path string) error {

	// read the details
	sub, err := parselearn.ParseLearnReceipt(path)
	if err != nil {
		return err
	}

	if gradexpath.IsPdf(sub.Filename) {
		if _, err := os.Stat(sub.Filename); os.IsNotExist(err) {
			return err
		}

	}
	return nil
}

func IdentifyReplacementPdf(path string) error {

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			//expect file at same level as receipt, save time
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	return err
}
