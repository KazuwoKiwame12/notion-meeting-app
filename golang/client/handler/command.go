package handler

import (
	"app/config"
	"app/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

type CommandHandler struct {
	commandUC *usecase.CommandUsecase
}

type responseJson struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewCommandHandler(commandUC *usecase.CommandUsecase) *CommandHandler {
	return &CommandHandler{
		commandUC: commandUC,
	}
}

func (ch *CommandHandler) StartScheduler(c echo.Context) error {
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	if _, alreadyExsit := ch.commandUC.ProcessManager[userID]; alreadyExsit {
		return c.JSON(http.StatusLocked, responseJson{StatusCode: http.StatusLocked, Message: "your scheduler is already runned"})
	}
	go ch.commandUC.Start(userID)
	return c.JSON(http.StatusAccepted, responseJson{StatusCode: http.StatusAccepted, Message: "your request successed!! scheduler is runnning ..."})
}

func (ch *CommandHandler) StopScheduler(c echo.Context) error {
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	if _, alreadyExsit := ch.commandUC.ProcessManager[userID]; !alreadyExsit {
		return c.JSON(http.StatusBadRequest, responseJson{StatusCode: http.StatusBadRequest, Message: "your scheduler is already canceled or you didn't start your scheduler"})
	}
	go ch.commandUC.Stop(userID)
	return c.JSON(http.StatusOK, responseJson{StatusCode: http.StatusOK, Message: "your request successed!! scheduler is canceled"})
}

// TODO implementation
func (ch *CommandHandler) RegisterNotionInfo(c echo.Context) error {
	return nil
}

func (ch *CommandHandler) UpdateNotionToken(c echo.Context) error {
	return nil
}

func (ch *CommandHandler) UpdateNotionDatabaseID(c echo.Context) error {
	return nil
}

func (ch *CommandHandler) UpdateNotionPageContent(c echo.Context) error {
	return nil
}

func (ch *CommandHandler) UpdateSchedulerDate(c echo.Context) error {
	return nil
}

func (ch *CommandHandler) ExplainHowToUse(c echo.Context) error {
	text := "このアプリでは、notionに議事録のテンプレートページを定期的に自動生成するスケジューラを起動・停止することができます。\n" +
		"スケジューラを動かすためには、以下の手順を行います。\n" +
		"1. ショートカット'Register the notion info'を選択し、表示されるモーダルにnotion情報を登録します。\n" +
		"2. /startというslash commandを呼び出すことで、スケジューラが起動します。\n\n" +
		"スケジューラを停止させるためには、/stopを実行すればスケジューラは停止します。\n" +
		"※1ユーザにつき1スケジューラであるために、現時点では複数台のスケジューラを起動させることができません。" +
		"そのような機能が必要であれば、管理人に連絡してください。"
	data := &slack.WebhookMessage{
		Text: text,
	}
	if err := slack.PostWebhook(config.WEBHOOK_URL(), data); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
