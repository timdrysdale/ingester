package ingester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/timdrysdale/parselearn"
)

func (g *Ingester) ValidateNewPapers() error {

	// wait for user to press an "do import new scripts button", then check the temp-txt and temp-pdf dirs
	possibleReceipts, err := g.GetFileList(g.TempTXT())
	if err != nil {
		return err
	}

	for _, receipt := range possibleReceipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)

		if err != nil {
			fmt.Println(err)
			continue // assume there may be others uses for txt, and that clean up will happen at end of the ingest
		}

		// assume we want to process this exam at some point - so set up the structure now
		// if it does not exist already
		_, err = os.Stat(g.GetExamPath(sub.Assignment))
		if os.IsNotExist(err) {
			err = g.SetupExamPaths(sub.Assignment)
			if err != nil {
				return err // If we can't set up a new exam, we may as well bail out
			}
		}

		pdfFilename, err := GetPDFPath(sub.Filename, g.TempPDF())
		if err != nil {
			continue
		}

		// file we want to get from the temp-pdf dir
		currentPath := filepath.Join(g.TempPDF(), filepath.Base(pdfFilename))

		_, err = os.Stat(currentPath)

		if !os.IsNotExist(err) { //double negative, file exists

			err = g.MoveIfNewerThanDestinationInDir(currentPath, g.AcceptedPapers(sub.Assignment))
			if err != nil {
				// reject receipt visibly, if problem copying in the pdf
				fmt.Printf("wanted to copy [%s] but %v\n", currentPath, err)
				err = g.MoveIfNewerThanDestinationInDir(receipt, g.Ingest())
				if err != nil {
					fmt.Println(err)
					continue //carry on with the rest ... TODO flag this in case not actually a lost cause
				}
			}
			err = g.MoveIfNewerThanDestinationInDir(receipt, g.AcceptedReceipts(sub.Assignment))
			if err != nil {
				continue
			}

		} else {
			// TODO Need to flag this to the user
			fmt.Printf("wanted [%s] but does not exist?\n", currentPath)
		}

	}
	return nil
}
