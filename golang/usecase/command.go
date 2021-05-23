package usecase

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/net/context"
)

// TODO contextを構造体に入れるのはNGなので対策を考える
type CommandUsecase struct {
	ProcessManager map[string]context.CancelFunc
	ParentContext  context.Context
}

func (c *CommandUsecase) Start(userProcessID string) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}
	nowJST := time.Now().UTC().In(jst)
	date := nowJST.Format("2006-01-02")

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase()
	// s.Every(30).Seconds().Do(templateUC.CreateForTeamMeeting, userProcessID)
	// s.Every(30).Seconds().Do(templateUC.CreateForGeneralMeeting, userProcessID)
	s.Every(1).Week().Tuesday().At("09:00").Tag("default").Do(templateUC.CreateForTeamMeeting, date)
	s.Every(1).Week().Wednesday().At("09:00").Tag("default").Do(templateUC.CreateForGeneralMeeting, date)
	s.StartAsync()
	s.RunByTag("default")

	// cancel commandでスケジューラを停止させる
	ctx, cancel := context.WithCancel(c.ParentContext)
	c.ProcessManager[userProcessID] = cancel
	select {
	case <-ctx.Done():
		fmt.Println("----done-----")
		s.Stop()
		return
	}
}

func (c *CommandUsecase) Stop(userProcessID string) {
	stopScheduler := c.ProcessManager[userProcessID]
	stopScheduler()
	delete(c.ProcessManager, userProcessID)
}
