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
	socketMode := socketmode.New(
		webApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeInteractive:
				payload, _ := envelope.Data.(slack.InteractionCallback)
				switch payload.Type {
				case slack.InteractionTypeShortcut:
					if payload.CallbackID == "notion-info__form" {
						socketMode.Ack(*envelope.Request)
						err = slackUC.GetModalView(payload.TriggerID)
						if err != nil {
							socketMode.Debugf("\nFailed to opemn a modal: %v\n", err)
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
							socketMode.Debugf("\nFailed to register notion info: %v\n", err)
						}
					}
				default:
					socketMode.Debugf("\nSkipped: %v\n", payload)
				}
			default:
				socketMode.Debugf("\nSkipped: %v\n", envelope.Type)
			}
		}
	}()
	socketMode.Run()
}
