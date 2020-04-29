package ingester

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func countPages(inputPath string) (int, error) {

	numPages := 0

	f, err := os.Open(inputPath)
	if err != nil {
		return numPages, err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return numPages, err
	}

	defer f.Close()

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return numPages, err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return numPages, err
		}
	}

	numPages, err = pdfReader.GetNumPages()
	if err != nil {
		return numPages, err
	}

	return numPages, nil

}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func itemExists(sliceType interface{}, item interface{}) bool {
	slice := reflect.ValueOf(sliceType)

	if slice.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// when we read the Learn receipt, we might get a suffix for a word doc etc
// so find the pdf file in the target directory with the same base prefix name
// but possibly variable capitalisation of the suffix (handmade file!)
func GetPdfPath(filename, directory string) (string, error) {

	// if the original receipt says the submission was not pdf
	// we need to find a handmade PDF with possibly non-lower case suffix
	// so search for matching basename
	if !gradexpath.IsPdf(filename) {

		possibleFiles, err := gradexpath.GetFileList(directory)
		if err != nil {
			return "", err
		}

	LOOP:
		for _, file := range possibleFiles {
			want := gradexpath.BareFile(filename)
			got := gradexpath.BareFile(file)
			equal := strings.Compare(want, got) == 0
			if equal {
				filename = file
				break LOOP
			}
		}

	} else { //assume the file is there
		filename = filepath.Join(directory, filename)
	}
	return filename, nil
}

func GetShortLearnDate(sub parselearn.Submission) (string, error) {

	if sub == (parselearn.Submission{}) {
		return "", errors.New("Empty submission")
	}
	newDate := sub.DateSubmitted
	//Example: "Tuesday, 23 April 2020 10-43-23 o'clock BST"

	tokens := strings.Split(newDate, " ")

	if len(tokens) == 7 {
		day := tokens[1]
		month := tokens[2]
		year := tokens[3]
		if len(month) >= 3 {
			month = month[0:3]
		}
		s := []string{day, month, year}
		newDate = strings.Join(s, "-")
	}

	// no change if we don't understand the format?
	return newDate, nil
}
