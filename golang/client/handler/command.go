package handler

import (
	"app/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
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
	userProcessID := c.FormValue("id")
	if _, alreadyExsit := ch.commandUC.ProcessManager[userProcessID]; alreadyExsit {
		return c.JSON(http.StatusLocked, responseJson{StatusCode: http.StatusLocked, Message: "your scheduler is already runned"})
	}
	go ch.commandUC.Start(userProcessID)
	return c.JSON(http.StatusAccepted, responseJson{StatusCode: http.StatusAccepted, Message: "your request successed!! scheduler is runnning ..."})
}

func (ch *CommandHandler) StopScheduler(c echo.Context) error {
	userProcessID := c.FormValue("id")
	if _, alreadyExsit := ch.commandUC.ProcessManager[userProcessID]; !alreadyExsit {
		return c.JSON(http.StatusBadRequest, responseJson{StatusCode: http.StatusBadRequest, Message: "your scheduler is already canceled or you didn't start your scheduler"})
	}
	go ch.commandUC.Stop(userProcessID)
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
	return nil
}