package ingester

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
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
				handleIngestPdf(path)
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

	// TODO check raw pdf?

	//TODO some reporting on what is left over? or another tool can do that?
	// and overall file system status tool?
	return nil
}

func handleIngestArchive(archivePath string) error {
	err := archiver.Unarchive(archivePath, gradexpath.Ingest())
	if err != nil {
		return err
	}
	err = os.Remove(archivePath)
	return err
}

func handleIngestPdf(path string) error {
	//return gradexpath.MoveIfNewerThanDestinationInDir(path, gradexpath.TempPdf())

	type PdfSummary struct {
		CourseCode  string
		PreparedFor string
		ToDo        string
	}

	t, err := pdfpagedata.TriagePdf(path)

	if err != nil {
		// no page data so likely a raw script
		return gradexpath.MoveIfNewerThanDestinationInDir(path, gradexpath.TempPdf())
	}

	switch t.ToDo {

	case "marking":

		origin := gradexpath.MarkerSent(t.CourseCode, t.PreparedFor)

		preOrigin := gradexpath.MarkerReady(t.CourseCode, t.PreparedFor)

		if gradexpath.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return err
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return err
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := gradexpath.MarkerBack(t.CourseCode, t.PreparedFor)
			return gradexpath.MoveIfNewerThanDestinationInDir(path, destination)
		}
	case "moderating":

		origin := gradexpath.ModeratorSent(t.CourseCode, t.PreparedFor)

		preOrigin := gradexpath.ModeratorReady(t.CourseCode, t.PreparedFor)

		if gradexpath.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return err
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return err
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := gradexpath.ModeratorBack(t.CourseCode, t.PreparedFor)
			return gradexpath.MoveIfNewerThanDestinationInDir(path, destination)
		}
	case "checking":

		origin := gradexpath.CheckerSent(t.CourseCode, t.PreparedFor)

		preOrigin := gradexpath.CheckerReady(t.CourseCode, t.PreparedFor)

		if gradexpath.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return err
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return err
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := gradexpath.CheckerBack(t.CourseCode, t.PreparedFor)
			return gradexpath.MoveIfNewerThanDestinationInDir(path, destination)
		}
	default:
		// check later to see if it has a learn receipt, etc
		return gradexpath.MoveIfNewerThanDestinationInDir(path, gradexpath.TempPdf())

	}

	return errors.New("Didn't know how to handle pdf ingest")
}
