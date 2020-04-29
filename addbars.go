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
		TaskName:     "add-mark-bar",
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
		Name:     "mark-bar",
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

func AddModerateActiveBar(exam string, moderator string, mch chan chmsg.MessageInfo) error {

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-moderate-active-bar",
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
		Name:     "moderate-active-bar",
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
		Name:     "moderating",
		UnixTime: time.Now().UnixNano(),
	}

	oc := OverlayCommand{
		FromPath:          gradexpath.ModerateActive(exam),
		ToPath:            gradexpath.ModeratorReady(exam, moderator),
		ExamName:          exam,
		TemplatePath:      gradexpath.OverlayLayoutSVG(),
		SpreadName:        "moderate-active",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    gradexpath.ModeratorABCDecoration(moderator),
	}

	err = OverlayPapers(oc)

	cm.Send(fmt.Sprintf("Finished Processing add-moderate-active-bar UUID=%s\n", uuidStr))

	return err
}

func AddModerateInActiveBar(exam string, mch chan chmsg.MessageInfo) error {

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-moderate-inactive-bar",
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
		Name:     "moderate-inactive-bar",
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
		Name:     "moderating",
		UnixTime: time.Now().UnixNano(),
	}

	oc := OverlayCommand{
		FromPath:          gradexpath.ModerateInActive(exam),
		ToPath:            gradexpath.ModeratedInActiveBack(exam),
		ExamName:          exam,
		TemplatePath:      gradexpath.OverlayLayoutSVG(),
		SpreadName:        "moderate-inactive",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    gradexpath.ModeratorABCDecoration("X"),
	}

	err = OverlayPapers(oc)

	cm.Send(fmt.Sprintf("Finished Processing add-moderate-inactive-bar UUID=%s\n", uuidStr))

	return err
}

func AddCheckBar(exam string, checker string, mch chan chmsg.MessageInfo) error {

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-check-bar",
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
		Name:     "check-bar",
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
		Name:     "checking",
		UnixTime: time.Now().UnixNano(),
	}

	oc := OverlayCommand{
		FromPath:          gradexpath.ModeratedReady(exam),
		ToPath:            gradexpath.CheckerReady(exam, checker),
		ExamName:          exam,
		TemplatePath:      gradexpath.OverlayLayoutSVG(),
		SpreadName:        "check",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    gradexpath.CheckerABCDecoration(checker),
	}

	err = OverlayPapers(oc)

	cm.Send(fmt.Sprintf("Finished Processing add-check-bar UUID=%s\n", uuidStr))

	return err
}
