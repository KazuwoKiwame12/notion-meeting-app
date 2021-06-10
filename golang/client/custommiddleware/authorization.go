package custommiddleware

import (
	"app/usecase"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

const (
	failed string = ""
)

func AuthUserMiddleware(auc *usecase.AuthorizationUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			workspaceID, slackUserID, err := getWorkspaceIDAndSlackUserID(c.FormValue("type"), c.FormValue("payload"))
			if err != nil {
				c.Error(err)
			}
			userID, userName, err := auc.IsUser(workspaceID, slackUserID)
			if err != nil {
				c.Error(err)
			}
			c.Request().Form.Set("user_id", strconv.Itoa(userID))
			c.Request().Form.Set("user_name", userName)
			return next(c)
		}
	}
}

func AuthAdminMiddleware(auc *usecase.AuthorizationUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			workspaceID, slackUserID, err := getWorkspaceIDAndSlackUserID(c.FormValue("type"), c.FormValue("payload"))
			if err != nil {
				c.Error(err)
			}
			userID, userName, err := auc.IsAdmin(workspaceID, slackUserID)
			if err != nil {
				c.Error(err)
			}
			c.Request().Form.Set("user_id", strconv.Itoa(userID))
			c.Request().Form.Set("user_name", userName)
			return next(c)
		}
	}
}

func getWorkspaceIDAndSlackUserID(requestType, payload string) (string, string, error) {
	switch requestType {
	case "slash_commands":
		var obj *slack.SlashCommand
		if err := json.Unmarshal([]byte(payload), obj); err != nil {
			return failed, failed, err
		}
		return obj.TeamID, obj.UserID, nil
	case "shortcut":
		var obj *slack.InteractionCallback
		if err := json.Unmarshal([]byte(payload), obj); err != nil {
			return failed, failed, err
		}
		return obj.User.TeamID, obj.User.ID, nil
	default:
		log.Println("invalid request type")
		return failed, failed, errors.New("invalid request type")
	}
}
