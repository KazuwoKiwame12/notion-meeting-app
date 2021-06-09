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

						msg := &slack.WebhookMessage{
							Blocks: &slack.Blocks{
								BlockSet: createResponseBlocks(fmt.Sprintf("@%s\nnotionの情報を登録しました。", payload.User.Name)),
							},
						}
						if err := slackUC.RegisterNotionInfo(payload.Team.ID, payload.User.ID, notion); err != nil {
							client.Debugf("\nFailed to register notion info: %v\n", err)
							// msg.Text = fmt.Sprintf("@%s\nnotionの情報の登録に失敗しました。", payload.User.Name)
						}
						if err := slack.PostWebhook(config.WEBHOOK_URL(), msg); err != nil {
							client.Debugf("Failed to call incominng web hook: %+v\n", err)
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

				switch cmd.Command {
				case "/explain":
					client.Ack(*envelope.Request, commandUC.GetExplainMessage(cmd.UserName))
				case "/start":
					user, err := dbOp.GetUser(cmd.TeamID, cmd.UserID)
					if err != nil {
						client.Debugf("\nFailed to get user: %v\n", err)
					}
					go commandUC.Start(user.ID)

					client.Ack(*envelope.Request, map[string]interface{}{"blocks": createResponseBlocks("your request successed!! scheduler is runnning ...")})
				case "/stop":
					user, err := dbOp.GetUser(cmd.TeamID, cmd.UserID)
					if err != nil {
						client.Debugf("\nFailed to get user: %v\n", err)
					}
					go commandUC.Stop(user.ID)
					client.Ack(*envelope.Request, map[string]interface{}{"blocks": createResponseBlocks("your request successed!! scheduler is canceled")})
				case "/all":
					text, err := commandUC.All()
					if err != nil {
						text = "名前の取得時にエラーが発生しました。"
						client.Ack(*envelope.Request, map[string]interface{}{"blocks": createResponseBlocks(text)})
					} else {
						client.Ack(*envelope.Request, map[string]interface{}{"blocks": createResponseBlocks("```" + text + "```")})
					}
				case "/all-stop":
					commandUC.AllStop()
					msg := &slack.WebhookMessage{
						Blocks: &slack.Blocks{
							BlockSet: createResponseBlocks("@channel\nメンテナンスのために、スケジューラを全て停止しました。再度スケジューラをスタート可能になった際に通知いたします。"),
						},
					}

					if err := slack.PostWebhook(config.WEBHOOK_URL(), msg); err != nil {
						log.Printf("Failed to call incominng web hook: %+v\n", err)
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

func createResponseBlocks(text string) []slack.Block {
	responseBlocks := []slack.Block{
		slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: text,
			},
			nil,
			nil,
		),
	}
	return responseBlocks
}
