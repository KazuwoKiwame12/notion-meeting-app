package handler

import (
	"app/usecase"
	"fmt"
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
		return c.JSON(http.StatusLocked, ch.createResponseMessage("your scheduler is already runned"))
	}
	go ch.commandUC.Start(userID)
	return c.JSON(http.StatusAccepted, ch.createResponseMessage("your request successed!! scheduler is runnning ..."))
}

func (ch *CommandHandler) StopScheduler(c echo.Context) error {
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	if _, alreadyExsit := ch.commandUC.ProcessManager[userID]; !alreadyExsit {
		return c.JSON(http.StatusBadRequest, responseJson{StatusCode: http.StatusBadRequest, Message: "your scheduler is already canceled or you didn't start your scheduler"})
	}
	go ch.commandUC.Stop(userID)
	return c.JSON(http.StatusOK, responseJson{StatusCode: http.StatusOK, Message: "your request successed!! scheduler is canceled"})
}

func (ch *CommandHandler) ExplainHowToUse(c echo.Context) error {
	name := c.FormValue("user_name")
	text := fmt.Sprintf("@%s\n", name) +
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
	return c.JSON(http.StatusOK, payload)
}

func (ch *CommandHandler) CheckAllProcess(c echo.Context) error {
	// TODO 実装
}

func (ch *CommandHandler) StopAllProcess(c echo.Context) error {
	// TODO 実装
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
