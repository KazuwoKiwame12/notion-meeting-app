package handler

import (
	"app/config"
	"app/domain/model"
	"app/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

type ModalHandler struct {
	slackUC *usecase.SlackUsecase
}

func NewModalHandler(slackUC *usecase.SlackUsecase) *ModalHandler {
	return &ModalHandler{
		slackUC: slackUC,
	}
}

func (mh *ModalHandler) CallModalOperation(c echo.Context) error {
	payload := c.FormValue("payload")
	var interactionObj slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &interactionObj); err != nil {
		log.Printf("error occurred for invalid payload: %+v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var callbackID string = ""
	switch interactionObj.TriggerID {
	case string(slack.InteractionTypeShortcut):
		callbackID = interactionObj.CallbackID
	case string(slack.InteractionTypeViewSubmission):
		callbackID = interactionObj.View.CallbackID
	}

	switch callbackID {
	case "notion-info__form":
		userID, _ := strconv.Atoi(c.FormValue("user_id"))
		if err := mh.slackUC.GetModalView(userID, interactionObj.TriggerID); err != nil {
			log.Printf("Failed to call views.open api: %+v\n", err)
		}
	case "notion-info__record":
		date, _ := strconv.Atoi(interactionObj.View.State.Values["scheduler-date"]["static_select-action"].Value)
		notion := model.Notion{
			Date:              date,
			NotionToken:       []byte(interactionObj.View.State.Values["notion-token"]["plain_text_input-action"].Value),
			NotionDatabaseID:  []byte(interactionObj.View.State.Values["notion-database_id"]["plain_text_input-action"].Value),
			NotionPageContent: interactionObj.View.State.Values["notion-page_content"]["plain_text_input-action"].Value,
		}
		text := fmt.Sprintf("@%s\nnotionの情報を登録しました。", interactionObj.User.Name)
		if err := mh.slackUC.RegisterNotionInfo(interactionObj.Team.ID, interactionObj.User.ID, notion); err != nil {
			log.Printf("Failed to store notion info: %+v\n", err)
			text = fmt.Sprintf("@%s\nnotionの情報の登録に失敗しました。", interactionObj.User.Name)
		}
		if err := slack.PostWebhook(config.WEBHOOK_URL(), mh.createResponseMessage(text)); err != nil {
			log.Printf("Failed to call incominng web hook: %+v\n", err)
		}
	default:
		log.Println("error occurred for an unexpected request.")
		text := fmt.Sprintf("@%s\n想定外のリクエストのためエラー。", interactionObj.User.Name)
		if err := slack.PostWebhook(config.WEBHOOK_URL(), mh.createResponseMessage(text)); err != nil {
			log.Printf("Failed to call incominng web hook: %+v\n", err)
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func (mh *ModalHandler) createResponseMessage(text string) *slack.WebhookMessage {
	msg := &slack.WebhookMessage{
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					&slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: text,
					},
					nil,
					nil,
				),
			},
		},
	}
	return msg
}
