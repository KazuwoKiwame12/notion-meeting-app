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
	nowJST := time.Now().UTC().In(jst)
	date := nowJST.Format("2006-01-02")

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase()
	// s.Every(1).Minutes().Do(templateUC.CreateForTeamMeeting, date)
	// s.Every(1).Minutes().Do(templateUC.CreateForGeneralMeeting, date)
	s.Tuesday().At("09:00").Do(templateUC.CreateForTeamMeeting, date)
	s.Wednesday().At("09:00").Do(templateUC.CreateForGeneralMeeting, date)

	s.StartAsync()
}

func (c *CommandUsecase) Cancel() {
}
