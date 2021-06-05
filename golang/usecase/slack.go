package usecase

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/slack-go/slack"
)

type SlackUsecase struct {
	DBOperator *function.DatabaseOperater
}

func (su *SlackUsecase) GetModalView(triggerID string) error {
	jsonfilePath := fmt.Sprintf("%s/src/app/asset/slack/modalview.json", os.Getenv("GOPATH")) // 絶対パスで参照(相対パスの場合、当然だが実行ファイルからの相対パスとして認識される)
	viewBytes, err := os.ReadFile(jsonfilePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %+v", err)
	}
	var viewObj slack.ModalViewRequest
	if err := json.Unmarshal(viewBytes, &viewObj); err != nil {
		return fmt.Errorf("jsonファイルのデータを構造体にマウントウトする際のエラー: %+v", err)
	}

	slackClient := slack.New(config.SLACK_TOKEN())
	if _, err := slackClient.OpenView(triggerID, viewObj); err != nil {
		return fmt.Errorf("views.open api error: %+v", err)
	}

	return nil
}

func (su *SlackUsecase) RegisterNotionInfo(workspaceID, slackUserID string, notion model.Notion) error {
	user, err := su.DBOperator.GetUser(workspaceID, slackUserID)
	if err != nil {
		return fmt.Errorf("get user error: %+v", err)
	}
	notion.UserID = user.ID
	notion.NotionToken, notion.NotionDatabaseID, err = encryptInfos(notion.NotionToken, notion.NotionDatabaseID)
	if err != nil {
		return err
	}

	if err := su.DBOperator.RegisterNotionInfo(&notion); err != nil {
		return fmt.Errorf("register notion error: %+v, notion info: %+v", err, notion)
	}
	return nil
}

func encryptInfos(token, databaseID []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher([]byte(config.ENCRYPTION_KEY()))
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("make cipher.Block error: %+v", err)
	}

	cipheTextForToken, cipheTextForDatabaseID := make([]byte, aes.BlockSize+len(token)), make([]byte, aes.BlockSize+len(databaseID))
	ivForToken, ivForDatabaseID := cipheTextForToken[:aes.BlockSize], cipheTextForDatabaseID[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, ivForToken); err != nil {
		return []byte{}, []byte{}, fmt.Errorf("make iv for token error: %+v", err)
	}
	if _, err := io.ReadFull(rand.Reader, ivForDatabaseID); err != nil {
		return []byte{}, []byte{}, fmt.Errorf("make iv for databaseID error: %+v", err)
	}
	encryptStreamForToken, encryptStreamForDatabaseID := cipher.NewCTR(block, ivForToken), cipher.NewCTR(block, ivForDatabaseID)
	encryptStreamForToken.XORKeyStream(cipheTextForToken[aes.BlockSize:], token) // バイナリ形式でxorされるので、存在しない文字コードが発生するために、string変換できない
	encryptStreamForDatabaseID.XORKeyStream(cipheTextForDatabaseID[aes.BlockSize:], databaseID)
	return cipheTextForToken, cipheTextForDatabaseID, nil
}
