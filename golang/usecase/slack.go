package usecase

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
	"encoding/json"
	"fmt"

	"io/ioutil"

	"github.com/slack-go/slack"
)

type SlackUsecase struct {
	DBOperator *function.DatabaseOperater
}

func (su *SlackUsecase) GetModalView(triggerID string) error {
	slackClient := slack.New(config.SLACK_TOKEN())
	viewJson, err := ioutil.ReadFile("./modal.json")
	if err != nil {
		return fmt.Errorf("jsonファイル読み込みのエラー: %+v", err)
	}
	var viewObj slack.ModalViewRequest
	if err := json.Unmarshal(viewJson, &viewObj); err != nil {
		return fmt.Errorf("jsonファイルのデータを構造体にマウントウトする際のエラー: %+v", err)
	}

	if _, err := slackClient.OpenView(triggerID, viewObj); err != nil {
		return fmt.Errorf("views.open api error: %+v", err)
	}

	return nil
}

func (su *SlackUsecase) RegisterNotionInfo(workspaceID, slackUserID string, notion model.Notion) error {
	user, err := su.DBOperator.GetUser(workspaceID, slackUserID)
	if err != nil {
		return fmt.Errorf("get user error: %+v", err)
	}
	notion.UserID = user.ID

	if err := su.DBOperator.RegisterNotionInfo(&notion); err != nil {
		return fmt.Errorf("register notion error: %+v", err)
	}
	return nil
}
