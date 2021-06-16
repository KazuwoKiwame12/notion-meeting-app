package function_test

import (
	"app/domain/function"
	"app/domain/model"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var properties = model.Properties{
	Name: struct {
		Title []struct {
			model.Text "json:\"text\""
		} "json:\"title\""
	}{
		Title: []struct {
			model.Text "json:\"text\""
		}{
			{
				Text: model.Text{
					Content: "sample",
				},
			},
		},
	},
}
var children = []model.Block{
	{
		Object: "block",
		Type:   model.BlockTypeHeading1,
		Heading1: &model.Heading{
			Text: []model.RichText{
				{
					Type: "text",
					Text: &model.Text{
						Content: "sample",
					},
				},
			},
		},
	},
}

func loadEnv() {
	envfilePath := fmt.Sprintf("%s/src/app/.env", os.Getenv("GOPATH"))
	if err := godotenv.Load(envfilePath); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestCreatePageSuccess(t *testing.T) {
	loadEnv()

	var requestbodyCorrect = model.Template{
		Parent: model.Parent{
			DatabaseID: os.Getenv("TEST_NOTION_DATABASE_ID"),
		},
		Properties: properties,
		Children:   children,
	}

	tests := []struct {
		name       string
		insertData struct {
			token  string
			params model.Template
		}
	}{
		{
			name: "通常のリクエスト",
			insertData: struct {
				token  string
				params model.Template
			}{
				token:  os.Getenv("TEST_NOTION_TOKEN"),
				params: requestbodyCorrect,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nc := function.NewNotionClient()
			if err := nc.CreatePage(test.insertData.token, test.insertData.params); err != nil {
				t.Errorf("error内容: %+v", err)
			}
		})
	}
}

func TestCreatePageFail(t *testing.T) {
	loadEnv()

	var requestbodyInvalied = model.Template{
		Parent: model.Parent{
			DatabaseID: os.Getenv("TEST_NOTION_DATABASE_ID"),
		},
		Properties: properties,
		Children: []model.Block{
			{
				Object: "invailed object", //ここは常に"block"である必要がある
			},
		},
	}
	var requestbodyInvaliedD = model.Template{
		Parent: model.Parent{
			DatabaseID: "invalied databaseID",
		},
		Properties: properties,
		Children:   children,
	}

	tests := []struct {
		name       string
		insertData struct {
			token  string
			params model.Template
		}
	}{
		{
			name: "無効なリクエストボディ",
			insertData: struct {
				token  string
				params model.Template
			}{
				token:  os.Getenv("TEST_NOTION_TOKEN"),
				params: requestbodyInvalied,
			},
		},
		{
			name: "無効なtoken",
			insertData: struct {
				token  string
				params model.Template
			}{
				token:  "",
				params: requestbodyInvalied,
			},
		},
		{
			name: "無効なdatabaseID",
			insertData: struct {
				token  string
				params model.Template
			}{
				token:  os.Getenv("TEST_NOTION_TOKEN"),
				params: requestbodyInvaliedD,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nc := function.NewNotionClient()
			if err := nc.CreatePage(test.insertData.token, test.insertData.params); err == nil {
				t.Error("errorが存在しないため、testは失敗")
			} else {
				t.Logf("error内容(errorを起こすテストであるため、成功): %+v", err)
			}
		})
	}
}
