package ingester

import (
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/timdrysdale/gradexpath"
)

// wait for user to press an "do ingest button", then filewalk to get the paths
func StageFromIngest() error {

	ingestPath := gradexpath.Ingest()

	// TODO consider listing paths then moving....
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
			case gradexpath.IsArchive(path):
				passAgain = true
				handleIngestArchive(path)
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

func handleIngestArchive(archivePath string) error {
	//suffix := filepath.Ext(archivePath)
	//archiveBase := filepath.Base(archivePath)
	//archiveName := strings.TrimSuffix(archiveBase, suffix)
	//temploc := fmt.Sprintf("tmp-unarchive-%s", strings.Replace(archiveName, " ", "", -1))
	//extractPath := filepath.Join(gradexpath.Ingest(), temploc)
	err := archiver.Unarchive(archivePath, gradexpath.Ingest())
	if err != nil {
		return err
	}
	err = os.Remove(archivePath)
	return err
}
