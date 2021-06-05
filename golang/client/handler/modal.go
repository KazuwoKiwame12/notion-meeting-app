package handler

import (
	"app/domain/model"
	"app/usecase"
	"encoding/json"
	"errors"
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
	log.Println("passed")

	payload := c.FormValue("payload")
	var interactionObj slack.InteractionCallback
	if err := json.Unmarshal([]byte(payload), &interactionObj); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var callbackID string = ""
	switch interactionObj.TriggerID {
	case string(slack.InteractionTypeShortcut):
		callbackID = interactionObj.CallbackID
	case string(slack.InteractionTypeViewSubmission):
		callbackID = interactionObj.View.CallbackID
	}

	var err error
	switch callbackID {
	case "notion-info__form":
		err = mh.slackUC.GetModalView(interactionObj.TriggerID)
	case "notion-info__record":
		date, _ := strconv.Atoi(interactionObj.View.State.Values["scheduler-date"]["static_select-action"].Value)
		notion := model.Notion{
			Date:              date,
			NotionToken:       []byte(interactionObj.View.State.Values["notion-token"]["plain_text_input-action"].Value),
			NotionDatabaseID:  []byte(interactionObj.View.State.Values["notion-database_id"]["plain_text_input-action"].Value),
			NotionPageContent: interactionObj.View.State.Values["notion-page_content"]["plain_text_input-action"].Value,
		}
		err = mh.slackUC.RegisterNotionInfo(interactionObj.Team.ID, interactionObj.User.ID, notion)
	default:
		err = errors.New("想定外の呼び出し方のためエラー")
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}
