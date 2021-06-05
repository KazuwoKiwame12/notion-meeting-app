package usecase

import (
	"app/config"
	"app/domain/function"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

type CommandUsecase struct {
	ProcessManager map[int]chan<- struct{}
	DBOperator     *function.DatabaseOperater
}

func (cu *CommandUsecase) Start(userID int) {
	notion, err := cu.DBOperator.GetNotionInfo(userID)
	if err != nil {
		log.Printf("database notion get error: %+v", err)
	}
	// 機密情報を解読
	plainTextForToken, plainTextForDatabaseID, err := decryptInfos(notion.NotionToken, notion.NotionDatabaseID)
	if err != nil {
		log.Printf("failed to decrypt: %+v", err)
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}
	// nowJST := time.Now().UTC().In(jst)
	// date := nowJST.Format("2006-01-02")

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase(cu.DBOperator)
	s.Every(30).Seconds().Do(templateUC.CreateForTeamMeeting, plainTextForDatabaseID, plainTextForToken, userID) // TODO tokenの暗号化を複合する処理を記述する
	s.Every(30).Seconds().Do(templateUC.CreateForGeneralMeeting, plainTextForDatabaseID, plainTextForToken, userID)
	// s.Every(1).Week().Tuesday().At("09:00").Tag("default").Do(templateUC.CreateForTeamMeeting, date)
	// s.Every(1).Week().Wednesday().At("09:00").Tag("default").Do(templateUC.CreateForGeneralMeeting, date)
	s.StartAsync()
	s.RunByTag("default")

	// cancel commandでスケジューラを停止させる
	cancelCh := make(chan struct{})
	cu.ProcessManager[userID] = cancelCh
	select {
	case <-cancelCh:
		fmt.Println("----done-----")
		s.Stop()
		return
	}
}

func (cu *CommandUsecase) Stop(userID int) {
	cancelCh := cu.ProcessManager[userID]
	close(cancelCh)
	delete(cu.ProcessManager, userID)
}

func decryptInfos(token, databaseID []byte) (string, string, error) {
	block, err := aes.NewCipher([]byte(config.ENCRYPTION_KEY()))
	if err != nil {
		return "", "", fmt.Errorf("make cipher.Block error: %+v", err)
	}

	decryptedTextForToken, decryptedTextForDatabaseID := make([]byte, len(token[aes.BlockSize:])), make([]byte, len(databaseID[aes.BlockSize:]))
	decryptStreamForToken, decryptStreamForDatabaseID := cipher.NewCTR(block, token[:aes.BlockSize]), cipher.NewCTR(block, databaseID[:aes.BlockSize])
	decryptStreamForToken.XORKeyStream(decryptedTextForToken, token[aes.BlockSize:])
	decryptStreamForDatabaseID.XORKeyStream(decryptedTextForDatabaseID, token[aes.BlockSize:])
	return string(decryptedTextForToken), string(decryptedTextForDatabaseID), nil
}
