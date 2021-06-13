package usecase_test

import (
	"app/domain/model"
	"app/usecase"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/slack-go/slack"
)

const (
	monday int = 0
)

func createProperties(text string) model.Properties {
	props := model.Properties{}
	props.Name.Title = append(props.Name.Title, struct {
		model.Text `json:"text"`
	}{
		Text: model.Text{
			Content: text,
		},
	})
	return props
}

func TestCreatePageFormat(t *testing.T) {
	type insertData struct {
		databaseID  string
		pageContent string
	}

	testsTitle := []struct {
		name       string
		insertData insertData
		want       *model.Template
	}{
		{
			name: "titleの作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: "title sample",
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
			},
		},
		{
			name: "titleの作成 (添字なし)",
			insertData: insertData{
				databaseID:  "",
				pageContent: "sample",
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
			},
		},
		{
			name: string(model.BlockTypeParagraph) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeParagraph)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeParagraph)),
				},
			},
		},
		{
			name: string(model.BlockTypeParagraph) + "の作成 (添字なし)",
			insertData: insertData{
				databaseID:  "",
				pageContent: "sample\nsample",
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeParagraph)),
				},
			},
		},
		{
			name: string(model.BlockTypeHeading1) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeHeading1)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeHeading1)),
				},
			},
		},
		{
			name: string(model.BlockTypeHeading2) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeHeading2)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeHeading2)),
				},
			},
		},
		{
			name: string(model.BlockTypeHeading3) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeHeading3)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeHeading3)),
				},
			},
		},
		{
			name: string(model.BlockTypeBulletedListItem) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeBulletedListItem)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeBulletedListItem)),
				},
			},
		},
		{
			name: string(model.BlockTypeNumberedListItem) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeNumberedListItem)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeNumberedListItem)),
				},
			},
		},
		{
			name: string(model.BlockTypeToDo) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeToDo)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeToDo)),
				},
			},
		},
		{
			name: string(model.BlockTypeToggle) + "の作成 (添字あり)",
			insertData: insertData{
				databaseID:  "",
				pageContent: fmt.Sprintf("sample\n%s sample", string(model.BlockTypeToggle)),
			},
			want: &model.Template{
				Parent: model.Parent{
					DatabaseID: "",
				},
				Properties: createProperties("sample"),
				Children: []model.Block{
					*(usecase.ExportCreateBlock("sample", model.BlockTypeToggle)),
				},
			},
		},
	}

	for _, test := range testsTitle {
		t.Run(test.name, func(t *testing.T) {
			result := usecase.ExportCreateTemplateFormat(test.insertData.databaseID, test.insertData.pageContent)
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("Failed to create corretc template format:\n\tresult = %+v\n\twant = %+v\n", result, test.want)
			}
		})
	}
}

func TestEmbedInCurrentNotionInfos(t *testing.T) {
	requestForTest, err := getModalFormat(false)
	if err != nil {
		t.Error(err)
	}
	want, err := getModalFormat(true)
	if err != nil {
		t.Error(err)
	}

	usecase.ExportEmbedInCurrentNotionInfos(requestForTest, "sample", "sample", "sample", monday)
	if !reflect.DeepEqual(requestForTest, want) {
		t.Errorf("Failed to embed infos in modal:\n\tresult = %+v\n\twant = %+v\n", requestForTest, want)
	}
}

func getModalFormat(requireWantData bool) (*slack.ModalViewRequest, error) {
	var jsonfilePath string
	if requireWantData {
		jsonfilePath = fmt.Sprintf("%s/src/app/asset/testdata/modalview_monday.json", os.Getenv("GOPATH"))
	} else {
		jsonfilePath = fmt.Sprintf("%s/src/app/asset/slack/modalview.json", os.Getenv("GOPATH"))
	}

	viewBytes, err := os.ReadFile(jsonfilePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile error: %+v", err)
	}
	var viewObj slack.ModalViewRequest
	if err := json.Unmarshal(viewBytes, &viewObj); err != nil {
		return nil, fmt.Errorf("jsonファイルのデータを構造体にマウントウトする際のエラー: %+v", err)
	}

	return &viewObj, nil
}
