package model

type Block struct {
	Object string    `json:"object"`
	Type   BlockType `json:"type"`

	Paragraph        *RichTextBlock `json:"paragraph,omitempty"`
	Heading1         *Heading       `json:"heading_1,omitempty"`
	Heading2         *Heading       `json:"heading_2,omitempty"`
	Heading3         *Heading       `json:"heading_3,omitempty"`
	BulletedListItem *RichTextBlock `json:"bulleted_list_item,omitempty"`
	NumberedListItem *RichTextBlock `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo          `json:"to_do,omitempty"`
	Toggle           *RichTextBlock `json:"toggle,omitempty"`
}

type RichTextBlock struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

type Heading struct {
	Text []RichText `json:"text"`
}

type ToDo struct {
	RichTextBlock
	Checked *bool `json:"checked,omitempty"`
}

type ChildPage struct {
	Title string `json:"title"`
}

type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeUnsupported      BlockType = "unsupported"
)
