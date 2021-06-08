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
		return c.JSON(http.StatusOK, ch.createResponseMessage("your scheduler is already runned"))
	}
	go ch.commandUC.Start(userID)
	return c.JSON(http.StatusOK, ch.createResponseMessage("your request successed!! scheduler is runnning ..."))
}

func (ch *CommandHandler) StopScheduler(c echo.Context) error {
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	if _, alreadyExsit := ch.commandUC.ProcessManager[userID]; !alreadyExsit {
		return c.JSON(http.StatusOK, ch.createResponseMessage("your scheduler is already canceled or you didn't start your scheduler"))
	}
	go ch.commandUC.Stop(userID)
	return c.JSON(http.StatusOK, ch.createResponseMessage("your request successed!! scheduler is canceled"))
}

func (ch *CommandHandler) ExplainHowToUse(c echo.Context) error {
	name := c.FormValue("user_name")
	return c.JSON(http.StatusOK, ch.commandUC.GetExplainMessage(name))
}

func (ch *CommandHandler) CheckAllProcess(c echo.Context) error {
	text, err := ch.commandUC.All()
	if err != nil {
		text = "名前の取得時にエラーが発生しました。"
		return c.JSON(http.StatusOK, ch.createResponseMessage(text))
	}
	return c.JSON(http.StatusOK, ch.createResponseMessage("```"+text+"```"))
}

func (ch *CommandHandler) StopAllProcess(c echo.Context) error {
	ch.commandUC.AllStop()
	msg := &slack.WebhookMessage{
		Text: "@channel\nメンテナンスのために、スケジューラを全て停止しました。再度スケジューラをスタート可能になった際に通知いたします。",
	}
	if err := slack.PostWebhook(config.WEBHOOK_URL(), msg); err != nil {
		return c.JSON(http.StatusOK, "")
	}
	return c.JSON(http.StatusOK, ch.createResponseMessage("全てのスケジューラを停止させました。"))
}

func (ch *CommandHandler) createResponseMessage(text string) map[string]interface{} {
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
