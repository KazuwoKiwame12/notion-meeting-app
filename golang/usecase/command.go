package usecase

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type CommandUsecase struct {
	ProcessManager map[string]chan<- struct{}
	DBOperator     *function.DBOperator
}

func (c *CommandUsecase) Start(userProcessID string) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}
	// nowJST := time.Now().UTC().In(jst)
	// date := nowJST.Format("2006-01-02")

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase(c.dbOperator)
	s.Every(30).Seconds().Do(templateUC.CreateForTeamMeeting, userProcessID)
	s.Every(30).Seconds().Do(templateUC.CreateForGeneralMeeting, userProcessID)
	// s.Every(1).Week().Tuesday().At("09:00").Tag("default").Do(templateUC.CreateForTeamMeeting, date)
	// s.Every(1).Week().Wednesday().At("09:00").Tag("default").Do(templateUC.CreateForGeneralMeeting, date)
	s.StartAsync()
	s.RunByTag("default")

	// cancel commandでスケジューラを停止させる
	cancelCh := make(chan struct{})
	c.ProcessManager[userProcessID] = cancelCh
	select {
	case <-cancelCh:
		fmt.Println("----done-----")
		s.Stop()
		return
	}
}

func (c *CommandUsecase) Stop(userProcessID string) {
	cancelCh := c.ProcessManager[userProcessID]
	close(cancelCh)
	delete(c.ProcessManager, userProcessID)
}
