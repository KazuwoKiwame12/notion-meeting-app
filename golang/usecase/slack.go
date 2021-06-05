package usecase

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
	"encoding/json"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type SlackUsecase struct {
	DBOperator *function.DatabaseOperater
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

func (su *SlackUsecase) GetModalView(triggerID string) error {
	jsonfilePath := fmt.Sprintf("%s/src/app/asset/slack/modalview.json", os.Getenv("GOPATH")) // 絶対パスで参照(相対パスの場合、当然だが実行ファイルからの相対パスとして認識される)
	viewBytes, err := os.ReadFile(jsonfilePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %+v", err)
	}
	var viewObj slack.ModalViewRequest
	if err := json.Unmarshal(viewBytes, &viewObj); err != nil {
		return fmt.Errorf("jsonファイルのデータを構造体にマウントウトする際のエラー: %+v", err)
	}

	slackClient := slack.New(config.SLACK_TOKEN())
	if _, err := slackClient.OpenView(triggerID, viewObj); err != nil {
		return fmt.Errorf("views.open api error: %+v", err)
	}

	return nil
}
