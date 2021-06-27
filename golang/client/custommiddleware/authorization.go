package custommiddleware

import (
	"app/config"
	"app/usecase"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

func AuthUserMiddleware(auc *usecase.AuthorizationUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return authFunc(c, next, auc.IsUser)
		}
	}
}

func AuthAdminMiddleware(auc *usecase.AuthorizationUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return authFunc(c, next, auc.IsAdmin)
		}
	}
}

func authFunc(c echo.Context, next echo.HandlerFunc, userAuthFunc func(string, string) (int, string, error)) error {
	reqBody := []byte{}
	if c.Request().Body != nil { // Read
		reqBody, _ = ioutil.ReadAll(c.Request().Body)
	}

	if err := verifySignature(c.Request().Header, reqBody, config.SLACK_SECRET()); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	workspaceID, slackUserID, payload, err := getIDsAndPayload(reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	userID, userName, err := userAuthFunc(workspaceID, slackUserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	c.Request().Form = url.Values{
		"user_id":   {fmt.Sprintf("%d", userID)},
		"user_name": {userName},
		"payload":   {payload},
	}
	return next(c)
}

func getIDsAndPayload(reqBody []byte) (string, string, string, error) {
	reqbodyString, err := url.QueryUnescape(string(reqBody))
	if err != nil {
		return "", "", "", err
	}

	if reqbodyString[0:8] == "payload=" {
		// interactive message
		var interactionObj slack.InteractionCallback

		if err := json.Unmarshal([]byte(reqbodyString[8:]), &interactionObj); err != nil {
			return "", "", "", err
		}
		return interactionObj.User.TeamID, interactionObj.User.ID, reqbodyString[8:], nil
	}

	//slash commands
	bodyContents := strings.Split(reqbodyString, "&")
	const (
		indexValue  int = 1
		indexTeamID int = 1
		indexUserID int = 5
	)
	teamIDContent, userIDContent := strings.Split(bodyContents[indexTeamID], "="), strings.Split(bodyContents[indexUserID], "=")
	workspaceID, err := url.QueryUnescape(teamIDContent[indexValue])
	if err != nil {
		return "", "", "", err
	}
	slackUserID, err := url.QueryUnescape(userIDContent[indexValue])
	if err != nil {
		return "", "", "", err
	}
	return workspaceID, slackUserID, "", nil
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
