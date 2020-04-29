package ingester

import (
	"fmt"
	"os/exec"
)

// simplified https://github.com/catherinelu/evangelist/blob/master/server.go

func convertPDFToJPEGs(pdfPath string, jpegPath string, outputFile string) error {

	//fmt.Printf("%s\n%s\n%s\n", pdfPath, jpegPath, outputFile)
	//info, err := os.Stat(pdfPath)
	//fmt.Println(err, info)
	//info, err = os.Stat(filepath.Dir(jpegPath))
	//fmt.Println(err, info)
	//info, err = os.Stat(filepath.Dir(outputFile))
	//fmt.Println(err, info)

	outputFileOption := fmt.Sprintf("-sOutputFile=%s", outputFile)

	cmd := exec.Command("gs", "-dNOPAUSE", "-sDEVICE=jpeg", outputFileOption, "-dJPEGQ=90", "-r175", "-q", pdfPath,
		"-c", "quit")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("gs command failed: %s\n", err.Error())
		return err
	}

	return nil
}

// This worked
// gs -dNOPAUSE -sDEVICE=jpeg -sOutputFile=edited-%d.jpg -dJPEGQ=95 -r300 -q edited5-covered.pdf -c quit
