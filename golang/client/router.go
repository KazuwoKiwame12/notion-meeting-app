package client

import (
	"app/client/custommiddleware"
	"app/client/handler"
	"app/config"
	"app/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(commandUC *usecase.CommandUsecase, authorizationUC *usecase.AuthorizationUsecase, slackUC *usecase.SlackUsecase) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.CORSAllowOrigin()},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	commandH := handler.NewCommandHandler(commandUC)
	modalH := handler.NewModalHandler(slackUC)

	userRoute := e.Group("/user", custommiddleware.AuthUserMiddleware(authorizationUC))
	userRoute.POST("/start", commandH.StartScheduler)
	userRoute.POST("/start", commandH.StartScheduler)
	userRoute.POST("/stop", commandH.StopScheduler)
	userRoute.POST("/modal/operation", modalH.CallModalOperation)
	userRoute.POST("/explain", commandH.ExplainHowToUse)

	adminRoute := e.Group("/admin", custommiddleware.AuthAdminMiddleware(authorizationUC))
	adminRoute.POST("/all", commandH.CheckAllProcess)
	adminRoute.POST("/all/stop", commandH.StopAllProcess)
	return e
}
