package client

import (
	"app/client/handler"
	"app/config"
	"app/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(commandUC *usecase.CommandUsecase) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.CORSAllowOrigin()},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	commandH := handler.NewCommandHandler(commandUC)
	e.POST("/start", commandH.StartScheduler)
	e.POST("/stop", commandH.StopScheduler)
	e.POST("/register/notion", commandH.RegisterNotionInfo)
	e.POST("/update/notion/token", commandH.UpdateNotionToken)
	e.POST("/update/notion/databaseID", commandH.UpdateNotionDatabaseID)
	e.POST("/update/notion/pageContent", commandH.UpdateNotionPageContent)
	e.POST("/update/scheduler/date", commandH.UpdateSchedulerDate)
	e.POST("/explain", commandH.ExplainHowToUse)
	return e
}
