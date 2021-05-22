package handler

import (
	"app/usecase"

	"github.com/labstack/echo/v4"
)

type CommandHandler struct {
	commandUC *usecase.CommandUsecase
}

func NewCommandHandler(commandUC *usecase.CommandUsecase) *CommandHandler {
	return &CommandHandler{
		commandUC: commandUC,
	}
}

func (ch *CommandHandler) StartScheduler(c echo.Context) error {
	go ch.commandUC.Start()
	return c.JSON(204, nil)
}

func (ch *CommandHandler) CancelScheduler(c echo.Context) error {
	go ch.commandUC.Cancel()
	return c.JSON(204, nil)
}
