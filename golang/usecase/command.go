package usecase

import (
	"app/domain/function"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

type CommandUsecase struct {
	ProcessManager map[int]chan<- struct{}
	DBOperator     *function.DatabaseOperater
}

const (
	monday int = iota
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

func (cu *CommandUsecase) Start(userID int) {
	notion, err := cu.DBOperator.GetNotionInfo(userID)
	if err != nil {
		log.Printf("database notion get error: %+v", err)
		return
	}
	plainTextForToken, plainTextForDatabaseID, err := notion.GetDecyptInfo()
	if err != nil {
		log.Printf("failed to decrypt: %+v", err)
		return
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase(cu.DBOperator)
	switch notion.Date {
	case monday:
		s = s.Every(1).Week().Monday()
	case tuesday:
		s = s.Every(1).Week().Tuesday()
		// s.Every(20).Seconds().Tag("default").Do(templateUC.CreateForMeeting, plainTextForDatabaseID, plainTextForToken, notion.NotionPageContent)
	case wednesday:
		s = s.Every(1).Week().Wednesday()
	case thursday:
		s = s.Every(1).Week().Thursday()
	case friday:
		s = s.Every(1).Week().Friday()
	case saturday:
		s = s.Every(1).Week().Saturday()
	case sunday:
		s = s.Every(1).Week().Sunday()
	default:
		log.Println("invailed value for date")
		return
	}

	s.At("09:00").Tag("default").Do(templateUC.CreateForMeeting, plainTextForDatabaseID, plainTextForToken, notion.NotionPageContent)
	s.StartAsync()
	s.RunByTag("default")

	// cancel commandでスケジューラを停止させる
	cancelCh := make(chan struct{})
	cu.ProcessManager[userID] = cancelCh
	select {
	case <-cancelCh:
		fmt.Println("----done-----")
		s.Stop()
		return
	}
}

func (cu *CommandUsecase) Stop(userID int) {
	cancelCh := cu.ProcessManager[userID]
	close(cancelCh)
	delete(cu.ProcessManager, userID)
}
