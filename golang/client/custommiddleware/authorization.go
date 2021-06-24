package custommiddleware

import (
	"app/config"
	"app/usecase"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			if err := verifySignature(c.Request().Header, reqBody, config.SLACK_SECRET()); err != nil {
				log.Printf("error occured! in verify signature:\n\t %+v, %+v, %+v\n", c.Request().Header, reqBody, err)
				return c.JSON(http.StatusBadRequest, err)
			}

			workspaceID, slackUserID, err := getWorkspaceIDAndSlackUserID(c.FormValue("type"), c.FormValue("payload"))
			if err != nil {
				log.Printf("error occured! in getWorkspaceIDAndSlackUserID:\n\t %+v, %+v, %+v\n", workspaceID, slackUserID, err)
				return c.JSON(http.StatusBadRequest, err)
			}
			userID, userName, err := auc.IsUser(workspaceID, slackUserID)
			if err != nil {
				log.Printf("error occured! in auc.IsUser:\n\t %+v, %+v, %+v\n", userID, userName, err)
				return c.JSON(http.StatusBadRequest, err)
			}
			c.Request().Form.Set("user_id", strconv.Itoa(userID))
			c.Request().Form.Set("user_name", userName)
			log.Printf("form_values: %+v\n", c.Request().Form)
			return next(c)
		}
	}
}

func AuthAdminMiddleware(auc *usecase.AuthorizationUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := verifySignature(c.Request().Header, []byte(c.FormValue("payload")), config.SLACK_SECRET()); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			workspaceID, slackUserID, err := getWorkspaceIDAndSlackUserID(c.FormValue("type"), c.FormValue("payload"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			userID, userName, err := auc.IsAdmin(workspaceID, slackUserID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			c.Request().Form.Set("user_id", strconv.Itoa(userID))
			c.Request().Form.Set("user_name", userName)
			return next(c)
		}
	}
}

/*リクエストがslackから送られてきたものであることの保証する関数*/
func verifySignature(header http.Header, body []byte, secret string) error {
	sv, err := slack.NewSecretsVerifier(header, secret)
	if err != nil {
		return err
	}
	sv.Write(body)                      // 生成するsignatureのベースとなる値に、payloadを付け加えている
	if err := sv.Ensure(); err != nil { // Ensureにて，slackのsignareと生成したsignatureの比較
		return err
	}

	return nil
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
