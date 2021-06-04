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

func (cu *CommandUsecase) Start(userID int) {
	notion, err := cu.DBOperator.GetNotionInfo(userID)
	if err != nil {
		log.Printf("database notion get error: %+v", err)
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}
	// nowJST := time.Now().UTC().In(jst)
	// date := nowJST.Format("2006-01-02")

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase(cu.DBOperator)
	s.Every(30).Seconds().Do(templateUC.CreateForTeamMeeting, notion.NotionDatabaseID, notion.NotionToken, userID) // TODO tokenの暗号化を複合する処理を記述する
	s.Every(30).Seconds().Do(templateUC.CreateForGeneralMeeting, notion.NotionDatabaseID, notion.NotionToken, userID)
	// s.Every(1).Week().Tuesday().At("09:00").Tag("default").Do(templateUC.CreateForTeamMeeting, date)
	// s.Every(1).Week().Wednesday().At("09:00").Tag("default").Do(templateUC.CreateForGeneralMeeting, date)
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
