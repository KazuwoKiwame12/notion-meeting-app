package main

import (
	"app/config"
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

	processManager := make(map[int]chan<- struct{})
	defer func() {
		for _, process := range processManager {
			close(process)
		}
	}()
	commandUC := &usecase.CommandUsecase{
		ProcessManager: processManager,
		DBOperator:     dbOp,
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
						user, err := dbOp.GetUser(payload.User.TeamID, payload.User.ID)
						if err != nil {
							client.Debugf("\nFailed to get user from db: %v\n", err)
						}
						err = slackUC.GetModalView(user.ID, payload.TriggerID)
						if err != nil {
							client.Debugf("\nFailed to opemn a modal: %v\n", err)
						}
					}
					client.Ack(*envelope.Request, createResponseMessage("shortcut called!"))
				case slack.InteractionTypeViewSubmission:
					if payload.View.CallbackID == "notion-info__record" {
						date, err := strconv.Atoi(payload.View.State.Values["scheduler-date"]["static_select-action"].SelectedOption.Value)
						if err != nil {
							client.Debugf("\nFailed to cast to int: %v, value: %d, original value: %s\n", err, date, payload.View.State.Values["scheduler-date"]["static_select-action"].SelectedOption.Value)
						}

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
						client.Ack(*envelope.Request, createResponseMessage("registered!"))
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

				switch cmd.Command {
				case "/explain":
					text := fmt.Sprintf("@%s\n", cmd.UserName) +
						"このアプリでは、notionに議事録のテンプレートページを定期的に自動生成するスケジューラを起動・停止することができます。\n" +
						"スケジューラを動かすためには、以下の手順を行います。\n" +
						"```1. ショートカット'Register the notion info'を選択し、表示されるモーダルにnotion情報を登録します。\n" +
						"2. /startというslash commandを呼び出すことで、スケジューラが起動します。```\n\n" +
						"スケジューラを停止させるためには、/stopを実行すればスケジューラは停止します。\n" +
						"また、notion情報を更新する際には、再度'1'の手順を実行してください。\n" +
						"※1ユーザにつき1スケジューラであるために、現時点では複数台のスケジューラを起動させることができません。" +
						"そのような機能が必要であれば、管理人に連絡してください。"
					client.Ack(*envelope.Request, createResponseMessage(text))
				case "/start":
					user, err := dbOp.GetUser(cmd.TeamID, cmd.UserID)
					if err != nil {
						client.Debugf("\nFailed to get user: %v\n", err)
					}
					go commandUC.Start(user.ID)
					client.Ack(*envelope.Request, createResponseMessage("executed!"))
				case "/stop":
					user, err := dbOp.GetUser(cmd.TeamID, cmd.UserID)
					if err != nil {
						client.Debugf("\nFailed to get user: %v\n", err)
					}
					go commandUC.Stop(user.ID)
					client.Ack(*envelope.Request, createResponseMessage("called!"))
				case "/all":
					text, err := commandUC.All()
					if err != nil {
						text = "名前の取得時にエラーが発生しました。"
						client.Ack(*envelope.Request, createResponseMessage(text))
					} else {
						client.Ack(*envelope.Request, createResponseMessage("```"+text+"```"))
					}
				case "/all-stop":
					commandUC.AllStop()
					msg := &slack.WebhookMessage{
						Text: "@channel\nメンテナンスのために、スケジューラを全て停止しました。再度スケジューラをスタート可能になった際に通知いたします。",
					}
					if err := slack.PostWebhook(config.WEBHOOK_URL(), msg); err != nil {
						client.Ack(*envelope.Request, "")
					} else {
						client.Ack(*envelope.Request, "全てのスケジューラを停止させました。")
					}
				}
			default:
				client.Debugf("\nSkipped: %v\n", envelope.Type)
			}
		}
	}()
	client.Run()
}

func createResponseMessage(text string) map[string]interface{} {
	responseMsg := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: text,
				},
				nil,
				nil,
			),
		},
	}

	return responseMsg
}
