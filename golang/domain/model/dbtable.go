package model

import (
	"app/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

type Workspace struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type User struct {
	ID              int
	SlackUserID     string
	WorkspaceID     string
	IsAdministrator bool
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Notion struct {
	ID                int
	UserID            int
	Date              int
	NotionToken       []byte
	NotionDatabaseID  []byte
	NotionPageContent string
}

func (n *Notion) SetEncryptInfo() error {
	block, err := aes.NewCipher([]byte(config.ENCRYPTION_KEY()))
	if err != nil {
		return fmt.Errorf("make cipher.Block error: %+v", err)
	}

	cipheTextForToken, cipheTextForDatabaseID := make([]byte, aes.BlockSize+len(n.NotionToken)), make([]byte, aes.BlockSize+len(n.NotionDatabaseID))
	ivForToken, ivForDatabaseID := cipheTextForToken[:aes.BlockSize], cipheTextForDatabaseID[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, ivForToken); err != nil {
		return fmt.Errorf("make iv for token error: %+v", err)
	}
	if _, err := io.ReadFull(rand.Reader, ivForDatabaseID); err != nil {
		return fmt.Errorf("make iv for databaseID error: %+v", err)
	}
	encryptStreamForToken, encryptStreamForDatabaseID := cipher.NewCTR(block, ivForToken), cipher.NewCTR(block, ivForDatabaseID)
	encryptStreamForToken.XORKeyStream(cipheTextForToken[aes.BlockSize:], n.NotionToken) // バイナリ形式でxorされるので、存在しない文字コードが発生するために、string変換できない
	encryptStreamForDatabaseID.XORKeyStream(cipheTextForDatabaseID[aes.BlockSize:], n.NotionDatabaseID)
	n.NotionToken, n.NotionDatabaseID = cipheTextForToken, cipheTextForDatabaseID
	return nil
}

func (n *Notion) GetDecyptInfo() (string, string, error) {
	block, err := aes.NewCipher([]byte(config.ENCRYPTION_KEY()))
	if err != nil {
		return "", "", fmt.Errorf("make cipher.Block error: %+v", err)
	}

	decryptedTextForToken, decryptedTextForDatabaseID := make([]byte, len(n.NotionToken[aes.BlockSize:])), make([]byte, len(n.NotionDatabaseID[aes.BlockSize:]))
	decryptStreamForToken, decryptStreamForDatabaseID := cipher.NewCTR(block, n.NotionToken[:aes.BlockSize]), cipher.NewCTR(block, n.NotionDatabaseID[:aes.BlockSize])
	decryptStreamForToken.XORKeyStream(decryptedTextForToken, n.NotionToken[aes.BlockSize:])
	decryptStreamForDatabaseID.XORKeyStream(decryptedTextForDatabaseID, n.NotionDatabaseID[aes.BlockSize:])
	return string(decryptedTextForToken), string(decryptedTextForDatabaseID), nil
}
