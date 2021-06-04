package usecase

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
	"encoding/json"
	"fmt"

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
	viewBytes := []byte(getViewWithJson())
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

func getViewWithJson() string {
	return `
	{
		"callback_id": "notion-info__record",
		"type": "modal",
		"title": {
			"type": "plain_text",
			"text": "My App",
			"emoji": true
		},
		"submit": {
			"type": "plain_text",
			"text": "Submit",
			"emoji": true
		},
		"close": {
			"type": "plain_text",
			"text": "Cancel",
			"emoji": true
		},
		"blocks": [
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "初めまして、このアプリは議事録フォーマットを定期的に自動作成するアプリです.\n実際に使用するためには以下の情報が必要です。\n1. notionのtoken\n2. notionのdatabase_id(議事録フォーマットを生成するページ)\n3. notionのpage_content(デフォルトで存在します)\n4. 議事録フォーマットを自動生成する曜日(デフォルトは水曜日です)\n\n *情報を入力してください*"
				}
			},
			{
				"type": "divider"
			},
			{
				"type": "input",
				"block_id": "notion-token",
				"element": {
					"type": "plain_text_input",
					"action_id": "plain_text_input-action"
				},
				"label": {
					"type": "plain_text",
					"text": "1. notionのtoken",
					"emoji": true
				}
			},
			{
				"type": "input",
				"block_id": "notion-database_id",
				"element": {
					"type": "plain_text_input",
					"action_id": "plain_text_input-action"
				},
				"label": {
					"type": "plain_text",
					"text": "2. notionのdatabase_id",
					"emoji": true
				}
			},
			{
				"type": "input",
				"block_id": "notion-page_content",
				"element": {
					"type": "plain_text_input",
					"multiline": true,
					"action_id": "plain_text_input-action"
				},
				"label": {
					"type": "plain_text",
					"text": "3. notionのpage_content",
					"emoji": true
				}
			},
			{
				"type": "input",
				"block_id": "scheduler-date",
				"element": {
					"type": "static_select",
					"placeholder": {
						"type": "plain_text",
						"text": "曜日を選択してください",
						"emoji": true
					},
					"options": [
						{
							"text": {
								"type": "plain_text",
								"text": "月曜日",
								"emoji": true
							},
							"value": "0"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "火曜日",
								"emoji": true
							},
							"value": "1"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "水曜日",
								"emoji": true
							},
							"value": "2"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "木曜日",
								"emoji": true
							},
							"value": "3"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "金曜日",
								"emoji": true
							},
							"value": "4"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "土曜日",
								"emoji": true
							},
							"value": "5"
						},
						{
							"text": {
								"type": "plain_text",
								"text": "日曜日",
								"emoji": true
							},
							"value": "6"
						}
					],
					"action_id": "static_select-action"
				},
				"label": {
					"type": "plain_text",
					"text": "4. 議事録フォーマットを自動生成する曜日",
					"emoji": true
				}
			},
			{
				"type": "divider"
			}
		]
	}`
}
