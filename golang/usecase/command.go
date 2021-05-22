package usecase

import (
	"time"

	"github.com/go-co-op/gocron"
)

type CommandUsecase struct {
}

func (c *CommandUsecase) Start() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}
	now := time.Now()
	nowUTC := now.UTC()
	nowJST := nowUTC.In(jst)
	s := gocron.NewScheduler(jst)
	// s.Sunday().At("09:00").Do(scheduler.CreateDailyDiaryTemplates, token, api_url, database_id)
	templateUC := NewTemplateUsecase()
	s.Every(1).Minutes().Do(templateUC.CreateForTeamMeeting, nowJST.Format("2006-01-02"))
	s.Every(1).Minutes().Do(templateUC.CreateForGeneralMeeting, nowJST.Format("2006-01-02"))
	s.StartAsync()
}

func (c *CommandUsecase) Cancel() {
}
