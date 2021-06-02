package observation

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
)

type SlackObserver struct {
	DBOperator *function.DatabaseOperater
}

func (so *SlackObserver) KeepSatate() {
	// slack client作成
	slackClient := slack.New(config.SLACK_TOKEN())

	// scheduler作成
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}

	s := gocron.NewScheduler(jst)
	s.Every(1).Month(1).At("09:00").Tag("default").Do(so.keepStateOfSlackTeam, slackClient)
	s.Every(1).Hour().Tag("default").Do(so.keepStateOfSlackUsers, slackClient)
	s.StartAsync()
	s.RunByTag("default")
}

func (so *SlackObserver) keepStateOfSlackUsers(sc *slack.Client) {
	usersOfSlack, err := sc.GetUsers()
	if err != nil {
		log.Printf("users.list api error: %+v", err)
	} else {
		usersOfModel := make([]model.User, len(usersOfSlack))
		for i := 0; i < len(usersOfSlack); i++ {
			usersOfModel[i] = model.User{
				SlackUserID:     usersOfSlack[i].ID,
				WorkspaceID:     usersOfSlack[i].TeamID,
				IsAdministrator: usersOfSlack[i].IsOwner,
				Name:            usersOfSlack[i].Name,
			}
		}

		errOfRegister := so.DBOperator.RegisterUsers(usersOfModel)
		if errOfRegister != nil {
			log.Printf("register users error: %+v", errOfRegister)
		}
	}
}

func (so *SlackObserver) keepStateOfSlackTeam(sc *slack.Client) {
	team, err := sc.GetTeamInfo()
	if err != nil {
		log.Printf("team.info api error: %+v", err)
	} else {
		workspace := &model.Workspace{
			ID:   team.ID,
			Name: team.Name,
		}
		errOfRegister := so.DBOperator.RegisterWorkspace(workspace)
		if errOfRegister != nil {
			log.Printf("register users error: %+v", errOfRegister)
		}
	}
}
