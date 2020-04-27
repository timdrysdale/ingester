package ingester

import (
	"errors"
	"fmt"
	"strings"
)

func MakeAnonymous(exam string) error {

	// assume someone hits a button - for this exam
	//a, err := anon.New(gradexpath.IdentityCSV())
	//
	//receipts, err := gradexpath.GetFileList(gradexpath.AcceptedReceipts(exam))
	//
	//if err != nil {
	//	return err
	//}

	return nil
}

// rename files
func checkMatriculation(m string) (bool, error) {
	expectedLength := 8
	actualLength := len(m)
	if actualLength != expectedLength {
		return false, errors.New(fmt.Sprintf("Wrong length got %d not %d", actualLength, expectedLength))
	}
	if strings.HasPrefix(strings.ToLower(m), "s") {
		return false, errors.New(fmt.Sprintf("Does not start with s"))
	}
	return true, nil
}

func checkExamNumber(m string) (bool, error) {
	expectedLength := 7
	actualLength := len(m)
	if actualLength != expectedLength {
		return false, errors.New(fmt.Sprintf("Wrong length got %d not %d", actualLength, expectedLength))
	}

	if strings.HasPrefix(strings.ToLower(m), "b") {
		return false, errors.New(fmt.Sprintf("Does not start with b"))

	}
	return true, nil
}
