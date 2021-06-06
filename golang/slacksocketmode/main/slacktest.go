package main

import (
	"app/domain/function"
	"app/domain/model"
	"app/infrastructure"
	"app/usecase"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func init() {
	envfilePath := fmt.Sprintf("%s/src/app/.env", os.Getenv("GOPATH"))
	if err := godotenv.Load(envfilePath); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func main() {
	// infrastructure初期化
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		log.Fatalf("newSqlHandler err: %+v", err)
	}
	defer func() {
		if err := sh.DB.Close(); err != nil {
			log.Fatalf("closed err: %+v", err)
		}
	}()

	// domain初期化
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}

	slackUC := &usecase.SlackUsecase{
		DBOperator: dbOp,
	}

	webApi := slack.New(
		os.Getenv("SLACK_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_SOCKET_TOKEN")),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	client := socketmode.New(
		webApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for envelope := range client.Events {
			switch envelope.Type {
			case socketmode.EventTypeInteractive:
				payload, _ := envelope.Data.(slack.InteractionCallback)
				switch payload.Type {
				case slack.InteractionTypeShortcut:
					if payload.CallbackID == "notion-info__form" {
						client.Ack(*envelope.Request)
						err = slackUC.GetModalView(payload.TriggerID)
						if err != nil {
							client.Debugf("\nFailed to opemn a modal: %v\n", err)
						}
					}
				case slack.InteractionTypeViewSubmission:
					if payload.View.CallbackID == "notion-info__record" {
						date, _ := strconv.Atoi(payload.View.State.Values["scheduler-date"]["static_select-action"].Value)
						notion := model.Notion{
							Date:              date,
							NotionToken:       []byte(payload.View.State.Values["notion-token"]["plain_text_input-action"].Value),
							NotionDatabaseID:  []byte(payload.View.State.Values["notion-database_id"]["plain_text_input-action"].Value),
							NotionPageContent: payload.View.State.Values["notion-page_content"]["plain_text_input-action"].Value,
						}
						err = slackUC.RegisterNotionInfo(payload.Team.ID, payload.User.ID, notion)
						if err != nil {
							client.Debugf("\nFailed to register notion info: %v\n", err)
						}
					}
				default:
					client.Debugf("\nSkipped: %v\n", payload)
				}
			case socketmode.EventTypeSlashCommand:
				cmd, ok := envelope.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", envelope)

					continue
				}

				client.Debugf("Slash command received: %+v", cmd)

				text := fmt.Sprintf("@%s\n", cmd.UserName) +
					"このアプリでは、notionに議事録のテンプレートページを定期的に自動生成するスケジューラを起動・停止することができます。\n" +
					"スケジューラを動かすためには、以下の手順を行います。\n" +
					"1. ショートカット'Register the notion info'を選択し、表示されるモーダルにnotion情報を登録します。\n" +
					"2. /startというslash commandを呼び出すことで、スケジューラが起動します。\n\n" +
					"スケジューラを停止させるためには、/stopを実行すればスケジューラは停止します。\n" +
					"また、notion情報を更新する際には、再度'1'の手順を実行してください。\n" +
					"※1ユーザにつき1スケジューラであるために、現時点では複数台のスケジューラを起動させることができません。" +
					"そのような機能が必要であれば、管理人に連絡してください。"
				payload := map[string]interface{}{
					"blocks": []slack.Block{
						slack.NewSectionBlock(
							&slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: text,
							},
							nil,
							nil,
						),
					}}
				client.Ack(*envelope.Request, payload)
			default:
				client.Debugf("\nSkipped: %v\n", envelope.Type)
			}
		}
	}()
	client.Run()
}
