package usecase

import (
	"app/domain/model"

	"github.com/slack-go/slack"
)

var ExportEmbedInCurrentNotionInfos func(*slack.ModalViewRequest, string, string, string, int) = embedInCurrentNotionInfos
var ExportCreateTemplateFormat func(string, string) *model.Template = createTemplateFormat
var ExportCreateBlock func(string, model.BlockType) *model.Block = createBlock
