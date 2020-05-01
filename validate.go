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

	// Map receipts, keeping only the latest revision for any given filename, ignoring dir and ext
	// so as to capture files in different dirs e.g. patch dirs, and with renamed filetypes
	receiptMap := make(map[string]parselearn.Submission)

	for _, receipt := range possibleReceipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)

		if err != nil {
			fmt.Println(err) //TODO spit this out on msgchan
			continue         // assume there may be others uses for txt, and that clean up will happen at end of the ingest
		}

		if existingSub, ok := receiptMap[fileKey(sub.Filename)]; ok {
			if sub.Revision > existingSub.Revision {
				receiptMap[fileKey(sub.Filename)] = sub
			}
		} else {
			receiptMap[fileKey(sub.Filename)] = sub
		}

	}

	// >>>>>>>>>>>>> drop IGNORE receipts >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	parselearn.HandleIgnoreReceipts(&receiptMap)

	// >>>>>>>>>>>>>>> drop multiple file submissions>>>>>>>>>>>>>>>>>>>>
	// look for, and reject, any multiple file submissions
	// these need flattening before merging so automatic merging
	// is a TODO - automatic flatten and merge multple  pdf submission
	// these need explicit patching to distinguish between us just taking
	// the first file before it is merged, with taking a merged file named
	// the same as the first file - with a patch receipt and a manual merge
	// we can name it what we like
	for k, v := range receiptMap {
		if v.NumberOfFiles > 1 {
			list, err := parselearn.GetFilePaths(v.OwnPath)
			if err != nil {
				continue
			}
			if len(list) > 1 {
				//TODO pipe these out over message channel
				fmt.Printf("REJECTING %s need to merge: %v\n", k, list)
				delete(receiptMap, k)
			}
		}
	}

	for _, sub := range receiptMap {

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
				err = g.MoveIfNewerThanDestinationInDir(sub.OwnPath, g.Ingest())
				if err != nil {
					fmt.Println(err)
					continue //carry on with the rest ... TODO flag this in case not actually a lost cause
				}
			}
			err = g.MoveIfNewerThanDestinationInDir(sub.OwnPath, g.AcceptedReceipts(sub.Assignment))
			if err != nil {
				continue
			}

		} else {
			// TODO Need to flag this to the user
			fmt.Printf("wanted [%s] but does not exist?\n", currentPath)
		}

	}

	// reject back to ingest anything we didn't take further
	rejectPDF, err := g.GetFileList(g.TempPDF())

	for _, reject := range rejectPDF {
		g.MoveToDir(reject, g.Ingest())
	}

	rejectTXT, err := g.GetFileList(g.TempTXT())

	for _, reject := range rejectTXT {
		g.MoveToDir(reject, g.Ingest())
	}

	return nil
}
