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
	if err := ch.commandUC.Start(); err != nil {
		return err
	}
	return c.JSON(204, nil)
}

func (ch *CommandHandler) CancelScheduler(c echo.Context) error {
	if err := ch.commandUC.Cancel(); err != nil {
		return err
	}
	return c.JSON(204, nil)
}
