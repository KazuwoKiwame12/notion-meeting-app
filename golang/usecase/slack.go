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

func (su *SlackUsecase) GetModalView(userID int, triggerID string) error {
	jsonfilePath := fmt.Sprintf("%s/src/app/asset/slack/modalview.json", os.Getenv("GOPATH")) // 絶対パスで参照(相対パスの場合、当然だが実行ファイルからの相対パスとして認識される)
	viewBytes, err := os.ReadFile(jsonfilePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %+v", err)
	}
	var viewObj slack.ModalViewRequest
	if err := json.Unmarshal(viewBytes, &viewObj); err != nil {
		return fmt.Errorf("jsonファイルのデータを構造体にマウントウトする際のエラー: %+v", err)
	}

	if notion, err := su.DBOperator.GetNotionInfo(userID); err == nil {
		plainToken, plainDatabaseID, err := notion.GetDecyptInfo()
		if err != nil {
			return fmt.Errorf("get decrypt token and databaseID error: %+v", err)
		}
		su.embedInCurrentNotionInfos(&viewObj, plainToken, plainDatabaseID, notion.NotionPageContent, notion.Date)
	}

	slackClient := slack.New(config.SLACK_TOKEN())
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
	if err := notion.SetEncryptInfo(); err != nil {
		return err
	}

	if err := su.DBOperator.RegisterNotionInfo(&notion); err != nil {
		return fmt.Errorf("register notion error: %+v, notion info: %+v", err, notion)
	}
	return nil
}

func (su *SlackUsecase) embedInCurrentNotionInfos(viewObj *slack.ModalViewRequest, token, databaseID, pageContent string, date int) {
	dateStringList := []string{"月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日", "日曜日"}
	embeddedData := []string{token, databaseID, pageContent, dateStringList[date]}
	var index int = 3
	for _, data := range embeddedData {
		sec := viewObj.Blocks.BlockSet[index].(*slack.SectionBlock)
		viewObj.Blocks.BlockSet[index] = &slack.SectionBlock{
			Type: sec.Type,
			Text: &slack.TextBlockObject{
				Type: sec.Text.Type,
				Text: sec.Text.Text + fmt.Sprintf("```%s```", data),
			},
		}
		index += 2
	}
}
