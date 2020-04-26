package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
)

func ValidateNewPapers() error {

	// wait for user to press an "do import new scripts button", then check the temp-txt and temp-pdf dirs
	possibleReceipts, err := gradexpath.GetFileList(gradexpath.TempTxt())
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
		_, err = os.Stat(gradexpath.GetExamPath(sub.Assignment))
		if os.IsNotExist(err) {
			err = gradexpath.SetupExamPaths(sub.Assignment)
			if err != nil {
				return err // If we can't set up a new exam, we may as well bail out
			}
		}

		// Get the corresponding PDF from the temp-pdf dir
		pdfFilename := sub.Filename

		// if the original receipt says the submission was not pdf
		// we need to find a handmade PDF with possibly non-lower case suffix
		// so search for matching basename
		if !gradexpath.IsPdf(pdfFilename) {

			possibleFiles, err := gradexpath.GetFileList(gradexpath.TempPdf())
			if err != nil {
				continue //TODO flag we've an issue?
			}

			for _, file := range possibleFiles {
				want := gradexpath.BareFile(sub.Filename)
				got := gradexpath.BareFile(file)
				equal := strings.Compare(want, got) == 0
				if equal {
					pdfFilename = file
				}
			}
		}

		// file we want to get from the temp-pdf dir
		currentPath := filepath.Join(gradexpath.TempPdf(), filepath.Base(pdfFilename))

		_, err = os.Stat(currentPath)

		if !os.IsNotExist(err) { //double negative, file exists

			err = gradexpath.MoveIfNewerThanDestinationInDir(currentPath, gradexpath.AcceptedPapers(sub.Assignment))
			if err != nil {
				// reject receipt visibly, if problem copying in the pdf
				fmt.Printf("wanted to copy [%s] but %v\n", currentPath, err)
				err = gradexpath.MoveIfNewerThanDestinationInDir(receipt, gradexpath.Ingest())
				if err != nil {
					fmt.Println(err)
					continue //carry on with the rest ... TODO flag this in case not actually a lost cause
				}
			}
			err = gradexpath.MoveIfNewerThanDestinationInDir(receipt, gradexpath.AcceptedReceipts(sub.Assignment))
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
