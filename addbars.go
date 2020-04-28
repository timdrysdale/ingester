package ingester

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
)

func AddMarkBar(exam string, marker string, mch chan chmsg.MessageInfo) error {

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "markbar",
	}

	cm := chmsg.New(mc, mch, 100*time.Millisecond)

	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuidStr,
		Previous: "", //dynamic
		UnixTime: time.Now().UnixNano(),
		Name:     "markbar",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0, //dynamic
	}

	UUIDBytes, err = uuid.NewRandom()
	uuidStr = UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	markDetails := pdfpagedata.QuestionDetails{
		UUID:     uuidStr,
		Name:     "marking",
		UnixTime: time.Now().UnixNano(),
	}

	oc := OverlayCommand{
		FromPath:          gradexpath.AnonymousPapers(exam),
		ToPath:            gradexpath.MarkerReady(exam, marker),
		ExamName:          exam,
		TemplatePath:      gradexpath.OverlayLayoutSVG(),
		SpreadName:        "mark",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    gradexpath.MarkerABCDecoration(marker),
	}

	err = OverlayPapers(oc)

	cm.Send(fmt.Sprintf("Finished Processing markbar UUID=%s\n", uuidStr))

	return err
}

/*
func AddMarkBar() {
}

func AddModActiveBar() {
}

func AddModInActiveBar() {
}

func AddCheckBar() {
}

//TODO ensure/add these paths to gradexpath
func FilesFrom(exam, stage string) {
	switch stage {
	case "mark":
		return gradexpath.AnonymousPapers(exam)
	case "moderate-active":
		return gradexpath.ModerateActivePapers(exam)
	case "moderate-inactive":
		return gradexpath.ModerateInactivePapers(exam)
	case "check":
		return gradexpath.CheckPapers(exam)
	default:
		return ""
	}
}

//TODO ensure/add these paths to gradexpath
func FilesTo(exam, stage string) {
	switch stage {
	case "mark":
		return gradexpath.ToMarkPapers(exam) // consider issue of number of markers separately ... ?
	case "moderate-active":
		return gradexpath.ToModerateActivePapers(exam)
	case "moderate-inactive":
		return gradexpath.ToModerateInactivePapers(exam)
	case "check":
		return gradexpath.ToCheckPapers(exam)
	default:
		return ""
	}
}
*/
