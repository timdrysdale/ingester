package ingester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
)

//func ValidateNewPapers() error {
//
//	// wait for user to press an "do import new scripts button", then check the temp-txt and temp-pdf dirs
//	possibleReceipts, err := gradexpath.GetFileList(gradexpath.TempTxt())
//	if err != nil {
//		return err
//	}
//	fmt.Println(len(possibleReceipts))
//	for i, receipt := range possibleReceipts {
//		fmt.Printf("%d %s\n", i, receipt)
//		sub, err := parselearn.ParseLearnReceipt(receipt)
//
//		if err != nil {
//			fmt.Println(err)
//			continue // assume there may be others uses for txt, and that clean up will happen at end of the ingest
//		}
//		//parsesvg.PrettyPrintStruct(sub)
//
//		if !gradexpath.IsPdf(sub.Filename) {
//
//			possibleFiles, err := gradexpath.GetFileList(gradexpath.TempPdf())
//
//			fullbase := filepath.Base(sub.File)
//			ext := 	filepath.Ext(base)
//
//			base := strings.TrimSuffix(filepath.Base(sub.File), filepath.Ext(sub.File))
//
//			possibleFiles, err := gradexpath.GetFileList(gradexpath.TempPdf())
//
//			for _, file := range possibleFiles {
//
//			}
//
//		}
//
//			currentPath := filepath.Join(gradexpath.TempPdf(), filepath.Base(sub.Filename))
//
//			fmt.Println(currentPath)
//			info, err := os.Stat(currentPath)
//			fmt.Println(info)
//			// source newer by definition if destination does not exist
//			if !os.IsNotExist(err) {
//				// TODO flag error somehow ... want to carry on processing though
//
//				_, err := os.Stat(gradexpath.GetExamPath(sub.Assignment))
//				if os.IsNotExist(err) {
//					err = gradexpath.SetupExamPaths(sub.Assignment)
//					if err != nil {
//						fmt.Println("Returning")
//						return err // If we can't set up a new exam, we may as well bail out
//					}
//				}
//
//				err = gradexpath.MoveIfNewerThanDestination(currentPath, gradexpath.AcceptedPapers(sub.Assignment))
//				if err != nil {
//					// reject receipt visibly, if problem copying in the pdf
//					err = gradexpath.MoveIfNewerThanDestination(receipt, gradexpath.Ingest())
//					if err != nil {
//						continue //carry on with the rest ... TODO flag this in case not actually a lost cause
//					}
//				}
//				err = gradexpath.MoveIfNewerThanDestination(receipt, gradexpath.AcceptedReceipts(sub.Assignment))
//				if err != nil {
//					continue
//				}
//
//			}
//		}
//	}
//	return nil
//}

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
