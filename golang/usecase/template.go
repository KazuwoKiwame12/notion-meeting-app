package usecase

import (
	"app/config"
	"app/domain/function"
	"app/domain/model"
)

type category int

const (
	isTeam category = iota
	isGeneral
)

type TemplateUsecase struct {
	client     *function.NotionClient
	dbOperator *function.DBOperator
}

func NewTemplateUsecase(dbOperator *function.DBOperator) *TemplateUsecase {
	 // TDODO 構成を考え直す＝アーキテクチャ考える...usecaseでuseidに紐着くtokenやpage content取得してインスタンス作るのあり
	return &TemplateUsecase{
		client:     function.NewNotionClient(),
		dbOperator: dbOperator,
	}
}

func (t *TemplateUsecase) CreateForTeamMeeting(date string) error {
	params := t.createTemplateFormat(date, isTeam)
	err := t.client.CreatePage(*params)
	if err != nil {
		// fmt.Printf("error内容: %w", err)
	}

	return err
}

func (t *TemplateUsecase) CreateForGeneralMeeting(date string) error {
	params := t.createTemplateFormat(date, isGeneral)
	return t.client.CreatePage(*params)
}

func (t *TemplateUsecase) createTemplateFormat(date string, categ category) *model.Template {
	params := &model.Template{}
	// database idの指定
	params.Parent.DatabaseID = config.DatabaseID()

	// タイトルの作成
	var title string
	if categ == isTeam {
		title = date + "チームミーティング"
	} else {
		title = date + "全体ミーティング"
	}
	params.Properties.Name.Title = append(params.Properties.Name.Title, struct {
		model.Text `json:"text"`
	}{
		Text: model.Text{
			Content: title,
		},
	})
	/*
		blockの作成
		共通で使うブロックの作成
	*/
	//Header2のブロック作成
	defaultTexts := []string{"進捗報告", "勉強会", "英語論文"}
	defaultHeaders := t.createBlocks(defaultTexts, model.BlockTypeHeading2)
	// 空白ブロックの作成
	blanckBlock := model.Block{
		Object: "block",
		Type:   model.BlockTypeParagraph,
		Paragraph: &model.RichTextBlock{
			Text: []model.RichText{
				{
					Type: "text",
					Text: &model.Text{
						Content: "",
					},
				},
			},
		},
	}
	// ブロックのリストの作成
	listBlock := make([]model.Block, 6)
	for i := 0; i < len(listBlock); i++ {
		if i%2 == 0 {
			listBlock[i] = defaultHeaders[i/2]
		} else {
			listBlock[i] = blanckBlock
		}
	}

	var insertBlocks []model.Block
	/*Childrenの作成*/
	if categ == isTeam {
		insertBlocks = listBlock
	} else {
		teams := []string{"チームN", "チームI", "チームS"}
		teamsHeaders := t.createBlocks(teams, model.BlockTypeHeading1)
		numOfBlocks := len(teams)*len(listBlock) + len(teams)
		insertBlocks = make([]model.Block, numOfBlocks)
		for i := 0; i < len(teams); i++ {
			startIndex := i * 7
			insertBlocks[startIndex] = teamsHeaders[i]
			for j, block := range listBlock {
				startIndexForListBlock := startIndex + 1
				insertBlocks[startIndexForListBlock+j] = block
			}
		}
	}
	params.Children = insertBlocks

	return params
}

func (t *TemplateUsecase) createBlocks(texts []string, blocktype model.BlockType) []model.Block {
	blocks := make([]model.Block, len(texts))
	for index, text := range texts {
		switch blocktype {
		case model.BlockTypeParagraph:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				Paragraph: &model.RichTextBlock{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeHeading1:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				Heading1: &model.Heading{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeHeading2:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				Heading2: &model.Heading{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeHeading3:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				Heading3: &model.Heading{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeBulletedListItem:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				BulletedListItem: &model.RichTextBlock{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeNumberedListItem:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				NumberedListItem: &model.RichTextBlock{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		case model.BlockTypeToDo:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				ToDo: &model.ToDo{
					RichTextBlock: model.RichTextBlock{
						Text: []model.RichText{
							{
								Type: "text",
								Text: &model.Text{
									Content: text,
								},
							},
						},
					},
				},
			}
		case model.BlockTypeToggle:
			blocks[index] = model.Block{
				Object: "block",
				Type:   blocktype,
				Toggle: &model.RichTextBlock{
					Text: []model.RichText{
						{
							Type: "text",
							Text: &model.Text{
								Content: text,
							},
						},
					},
				},
			}
		}
	}
	return blocks
}
