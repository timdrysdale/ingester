package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/unzip"
)

// wait for user to press an "do ingest button", then filewalk to get the paths
func StageFromIngest() error {

	ingestPath := gradexpath.Ingest()

	// consider listing paths then moving....
	//pdfPaths := []string{}
	//txtPaths := []string{}

LOOP:
	for {
		passAgain := false

		err := filepath.Walk(ingestPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			switch {
			//case gradexpath.IsZip(path):
			//	passAgain = true
			//	handleIngestZip(path)
			// TODO try another zip library
			case gradexpath.IsTxt(path):
				err := gradexpath.MoveIfNewerThanDestinationInDir(path, gradexpath.TempTxt())
				if err != nil {
					return err
				}
			case gradexpath.IsPdf(path):

				gradexpath.MoveIfNewerThanDestinationInDir(path, gradexpath.TempPdf())
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		if !passAgain {
			break LOOP
		}
	}

	//TODO some reporting on what is left over? or another tool can do that?
	// and overall file system status tool?
	return nil
}

func handleIngestZip(zipPath string) error {
	suffix := filepath.Ext(zipPath)
	zipBase := filepath.Base(zipPath)
	zipName := strings.TrimSuffix(zipBase, suffix)
	temploc := fmt.Sprintf("tmp-unzip-%s", strings.Replace(zipName, " ", "", -1))
	extractPath := filepath.Join(gradexpath.Ingest(), temploc)
	err := handleZip(zipPath, extractPath)
	return err
}

func handleZip(zipPath, extractPath string) error {
	uz := unzip.New(zipPath, extractPath)
	err := uz.Extract()
	if err != nil {
		return err
	}
	err = os.Remove(zipPath)
	return err
}
