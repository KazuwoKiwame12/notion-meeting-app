package model_test

import (
	"app/domain/model"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func loadEnv() {
	envfilePath := fmt.Sprintf("%s/src/app/.env", os.Getenv("GOPATH"))
	if err := godotenv.Load(envfilePath); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestEncryptionAndDecryption(t *testing.T) {
	loadEnv()

	sameToken, sameDatabaseID := "same token", "same databaseID"
	result := model.Notion{NotionToken: []byte(sameToken), NotionDatabaseID: []byte(sameDatabaseID)}
	want := result

	sameValue := func(result model.Notion, want model.Notion) bool {
		isSameToken := reflect.DeepEqual(result.NotionToken, want.NotionToken)
		isSameDatabseID := reflect.DeepEqual(result.NotionDatabaseID, want.NotionDatabaseID)
		return isSameToken && isSameDatabseID
	}

	if err := (&result).SetEncryptInfo(); err != nil {
		t.Errorf("encrypting error: %+v", err)
	}

	if sameValue(result, want) {
		t.Errorf("failed to encrypt: result = %+v, want = %+v", result, want)
	}

	token, databaseID, err := (&result).GetDecyptInfo()
	if err != nil {
		t.Errorf("decrpting error: %+v", err)
	}

	result.NotionToken, result.NotionDatabaseID = []byte(token), []byte(databaseID)
	if !sameValue(result, want) {
		t.Errorf("failed to encrypt: result = %+v, want = %+v", result, want)
	}
}
