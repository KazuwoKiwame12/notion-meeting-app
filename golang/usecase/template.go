package usecase

import (
	"app/domain/function"
	"app/domain/model"
	"strings"
)

type TemplateUsecase struct {
	client     *function.NotionClient
	dbOperator *function.DatabaseOperater
}

const (
	blockTypeIndex int = iota
	textIndex
	numOfContents
)

func NewTemplateUsecase(dbOperator *function.DatabaseOperater) *TemplateUsecase {
	return &TemplateUsecase{
		client:     function.NewNotionClient(),
		dbOperator: dbOperator,
	}
}

func (t *TemplateUsecase) CreateForMeeting(databaseID, token, pageContent string) error {
	params := t.createTemplateFormat(databaseID, pageContent)
	return t.client.CreatePage(token, *params)
}

func (t *TemplateUsecase) createTemplateFormat(databaseID, pageContent string) *model.Template {
	params := &model.Template{}
	// database idの指定
	params.Parent.DatabaseID = databaseID

	blocksByString := strings.Split(pageContent, "\n")
	title := strings.Split(blocksByString[0], " ")

	existBlockType := func(typeName string) bool {
		checker := map[model.BlockType]struct{}{
			model.BlockTypeParagraph:        struct{}{},
			model.BlockTypeHeading1:         struct{}{},
			model.BlockTypeHeading2:         struct{}{},
			model.BlockTypeHeading3:         struct{}{},
			model.BlockTypeBulletedListItem: struct{}{},
			model.BlockTypeNumberedListItem: struct{}{},
			model.BlockTypeToDo:             struct{}{},
			model.BlockTypeToggle:           struct{}{},
			model.BlockType("title"):        struct{}{},
		}
		_, ok := checker[model.BlockType(typeName)]
		return ok
	}

	createTextFromContents := func(contents []string) string {
		var text string
		if len(contents) > numOfContents {
			if existBlockType(contents[blockTypeIndex]) {
				text = strings.Join(contents[textIndex:], " ")
			} else {
				text = strings.Join(contents, " ")
			}
		} else if len(contents) == numOfContents {
			if existBlockType(contents[blockTypeIndex]) {
				text = contents[textIndex]
			} else {
				text = strings.Join(contents, " ")
			}
		} else if len(contents) == 1 {
			text = contents[0]
		} else {
			text = ""
		}

		return text
	}

	params.Properties.Name.Title = append(params.Properties.Name.Title, struct {
		model.Text `json:"text"`
	}{
		Text: model.Text{
			Content: createTextFromContents(title),
		},
	})

	params.Children = make([]model.Block, 0, len(blocksByString[1:]))
	for _, blockByString := range blocksByString[1:] {
		contents := strings.Split(blockByString, " ")
		text := createTextFromContents(contents)
		block := t.createBlock(text, model.BlockType(contents[blockTypeIndex]))
		params.Children = append(params.Children, *block)
	}

	return params
}

func (t *TemplateUsecase) createBlock(text string, blocktype model.BlockType) *model.Block {
	var block *model.Block
	switch blocktype {
	case model.BlockTypeParagraph:
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
		block = &model.Block{
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
	default:
		block = &model.Block{
			Object: "block",
			Type:   model.BlockTypeParagraph,
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
	}

	return block
}
