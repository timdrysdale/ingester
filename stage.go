package ingester

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
	"github.com/timdrysdale/pdfpagedata"
)

func (g *Ingester) CleanFromIngest() error {

	files, err := g.GetFileList(g.Ingest())
	if err != nil {
		return err
	}
	errorCache := error(nil)

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			//count errors?
			errorCache = err
		}
	}
	return errorCache
}

// wait for user to press an "do ingest button", then filewalk to get the paths
func (g *Ingester) StageFromIngest() error {

	ingestPath := g.Ingest()

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
			case g.IsArchive(path):
				passAgain = true
				g.HandleIngestArchive(path)
			case IsTXT(path):
				err := g.MoveIfNewerThanDestinationInDir(path, g.TempTXT())
				if err != nil {
					return err
				}
			case IsPDF(path):
				g.HandleIngestPDF(path)

			case IsCSV(path):
				g.HandleIngestCSV(path)
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

func (g *Ingester) HandleIngestArchive(archivePath string) error {
	err := archiver.Unarchive(archivePath, g.Ingest())
	if err != nil {
		return err
	}
	err = os.Remove(archivePath)
	return err
}

func (g *Ingester) HandleIngestCSV(path string) error {
	if strings.ToLower(filepath.Base(path)) == "identity.csv" {
		return g.MoveIfNewerThanDestinationInDir(path, g.Identity())
	}
	return nil
	// leave file in ingest if not newer - to overwrite current file with an older version
	// e.g. to roll back a change, you have to roll forward by modifying the old file,
	// saving it to get a new modtime (can change back the mod before ingesting if needed)
	// just need the new mod time on the file
}

func (g *Ingester) HandleIngestPDF(path string) error {
	//return g.MoveIfNewerThanDestinationInDir(path, g.TempPDF())

	type PDFSummary struct {
		CourseCode  string
		PreparedFor string
		ToDo        string
	}

	t, err := pdfpagedata.TriagePdf(path)

	if err != nil {
		// no page data so likely a raw script
		return g.MoveIfNewerThanDestinationInDir(path, g.TempPDF())
	}

	switch t.ToDo {

	case "flattening":

		// these aren't usually exported, but we may be repopulating a new ingester or
		// manually correcting something, so we consider our options
		origin := g.AnonymousPapers(t.CourseCode)
		return g.MoveIfNewerThanDestinationInDir(path, origin)
		// leave the file in ingest if we don't want it

	case "marking":
		// these could be marked, or just being returned by DSA if prematurely exported
		origin := g.MarkerSent(t.CourseCode, t.PreparedFor)

		preOrigin := g.MarkerReady(t.CourseCode, t.PreparedFor)

		if g.IsSameAsSelfInDir(path, origin) {
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
			destination := g.MarkerBack(t.CourseCode, t.PreparedFor)
			return g.MoveIfNewerThanDestinationInDir(path, destination)
		}
	case "moderating":

		origin := g.ModeratorSent(t.CourseCode, t.PreparedFor)

		preOrigin := g.ModeratorReady(t.CourseCode, t.PreparedFor)

		if g.IsSameAsSelfInDir(path, origin) {
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
			destination := g.ModeratorBack(t.CourseCode, t.PreparedFor)
			return g.MoveIfNewerThanDestinationInDir(path, destination)
		}
	case "checking":

		origin := g.CheckerSent(t.CourseCode, t.PreparedFor)

		preOrigin := g.CheckerReady(t.CourseCode, t.PreparedFor)

		if g.IsSameAsSelfInDir(path, origin) {
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
			destination := g.CheckerBack(t.CourseCode, t.PreparedFor)
			return g.MoveIfNewerThanDestinationInDir(path, destination)
		}
	default:
		// check later to see if it has a learn receipt, etc
		return g.MoveIfNewerThanDestinationInDir(path, g.TempPDF())

	}

	return errors.New("Didn't know how to handle pdf ingest")
}
