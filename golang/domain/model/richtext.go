package model

import "time"

type RichText struct {
	Type        RichTextType `json:"type,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`

	PlainText string    `json:"plain_text,omitempty"`
	HRef      *string   `json:"href,omitempty"`
	Text      *Text     `json:"text,omitempty"`
	Equation  *Equation `json:"equation,omitempty"`
}

type Equation struct {
	Expression string `json:"expression"`
}

type Annotations struct {
	Bold          bool  `json:"bold,omitempty"`
	Italic        bool  `json:"italic,omitempty"`
	Strikethrough bool  `json:"strikethrough,omitempty"`
	Underline     bool  `json:"underline,omitempty"`
	Code          bool  `json:"code,omitempty"`
	Color         Color `json:"color,omitempty"`
}

type Date struct {
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end,omitempty"`
}

type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

type Link struct {
	URL string `json:"url"`
}

type (
	RichTextType string
	Color        string
)

const (
	RichTextTypeText     RichTextType = "text"
	RichTextTypeEquation RichTextType = "equation"
)

const (
	ColorDefault  Color = "default"
	ColorGray     Color = "gray"
	ColorBrown    Color = "brown"
	ColorOrange   Color = "orange"
	ColorYellow   Color = "yellow"
	ColorGreen    Color = "green"
	ColorBlue     Color = "blue"
	ColorPurple   Color = "purple"
	ColorPink     Color = "pink"
	ColorRed      Color = "red"
	ColorGrayBg   Color = "gray_background"
	ColorBrownBg  Color = "brown_background"
	ColorOrangeBg Color = "orange_background"
	ColorYellowBg Color = "yellow_background"
	ColorGreenBg  Color = "green_background"
	ColorBlueBg   Color = "blue_background"
	ColorPurpleBg Color = "purple_background"
	ColorPinkBg   Color = "pink_background"
	ColorRedBg    Color = "red_background"
)
