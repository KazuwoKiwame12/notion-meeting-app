package usecase

import (
	"app/domain/function"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
)

type CommandUsecase struct {
	ProcessManager map[int]chan<- struct{}
	DBOperator     *function.DatabaseOperater
}

const (
	monday int = iota
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

func (cu *CommandUsecase) Start(userID int) {
	notion, err := cu.DBOperator.GetNotionInfo(userID)
	if err != nil {
		log.Printf("database notion get error: %+v", err)
		return
	}
	plainTextForToken, plainTextForDatabaseID, err := notion.GetDecyptInfo()
	if err != nil {
		log.Printf("failed to decrypt: %+v", err)
		return
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("JST", +9*60*60)
	}

	s := gocron.NewScheduler(jst)
	templateUC := NewTemplateUsecase(cu.DBOperator)
	switch notion.Date {
	case monday:
		s = s.Every(1).Week().Monday()
	case tuesday:
		s = s.Every(1).Week().Tuesday()
		// s.Every(20).Seconds().Tag("default").Do(templateUC.CreateForMeeting, plainTextForDatabaseID, plainTextForToken, notion.NotionPageContent)
	case wednesday:
		s = s.Every(1).Week().Wednesday()
	case thursday:
		s = s.Every(1).Week().Thursday()
	case friday:
		s = s.Every(1).Week().Friday()
	case saturday:
		s = s.Every(1).Week().Saturday()
	case sunday:
		s = s.Every(1).Week().Sunday()
	default:
		log.Println("invailed value for date")
		return
	}

	s.At("09:00").Tag("default").Do(templateUC.CreateForMeeting, plainTextForDatabaseID, plainTextForToken, notion.NotionPageContent)
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

func (cu *CommandUsecase) GetExplainMessage(name string) map[string]interface{} {
	target := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("@%s\n", name),
		},
		nil,
		nil,
	)
	headerCanDo := slack.NewHeaderBlock(
		&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "このアプリでできること",
		},
	)
	descriptionCanDo := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "1. notionに議事録のテンプレートページを定期的に自動作成できる\n" +
				"2. 作成するテンプレートページの内容・場所・間隔を簡単にカスタマイズできる\n" +
				"3. いつでも自動作成するのを停止できます。自動作成の軌道も同様に簡単にできる\n",
		},
		nil,
		nil,
	)
	headerHowToUse := slack.NewHeaderBlock(
		&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "このアプリの使い方",
		},
	)
	descriptionHowToUse := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "1. ショートカット `Register the notion info`を選択する\n" +
				"2. 表示されたmodalにある以下の項目を埋める\n" +
				"```" +
				"1. notionのtoken\n" +
				"2. notionのページを作成する場所\n" +
				"3. notionのテンプレートページの内容\n" +
				"4. テンプレートページを作成する曜日\n" +
				"```\n\n" +
				"3. `/start`のslash commandでスケジューラを起動させる\n" +
				"※スケジューラを停止させたい場合は `/stop`のslash commandで停止させる\n",
		},
		nil,
		nil,
	)
	headerHowToCreatePage := slack.NewHeaderBlock(
		&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "notion情報登録の際の、ページ内容の作り方",
		},
	)
	descriptionHowToCreatePage := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "1行目はページのタイトルに該当し、2行目以降はページの中身に該当する。\n" +
				"文章・heading1・heading2・heading3・リスト・番号付きリスト・todoリスト・toggleリストなどのテキストタイプが使える。\n" +
				"テキストタイプを変更する場合には `Enterキー`で改行する必要がある。\n" +
				"ページ内容を記入する際には以下の形式である必要がある。\n" +
				"```" +
				"`テキストタイプ` `テキスト`..." +
				"```\n\n" +
				"以下にテキストタイプの種類を示す。\n",
		},
		[]*slack.TextBlockObject{
			{
				Type: "mrkdwn",
				Text: "*text type*",
			},
			{
				Type: "mrkdwn",
				Text: "*about*",
			},
			{
				Type: "mrkdwn",
				Text: "`paragraph` or `空白`",
			},
			{
				Type: "mrkdwn",
				Text: "文章",
			},
			{
				Type: "mrkdwn",
				Text: "`heading_1`・ `heading_2`・ `heading_3`",
			},
			{
				Type: "mrkdwn",
				Text: "1・2・3番大きいheader",
			},
			{
				Type: "mrkdwn",
				Text: "`toggle`・ `bulleted_list_item`",
			},
			{
				Type: "mrkdwn",
				Text: "toggleリスト・リスト",
			},
			{
				Type: "mrkdwn",
				Text: "`to_do`・ `numbered_list_item`",
			},
			{
				Type: "mrkdwn",
				Text: "todoリスト・番号付きリスト",
			},
		},
		nil,
	)
	descriptionExample := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "以下にページ内容の登録例と実際のnotionの画面を示す。\n" +
				"```" +
				"title ミーティング\n" +
				"heading_2 概要\n" +
				"タイトル:\n\n" +
				"heading_2 内容\n" +
				"numbered_list_item sample\n\n" +
				"numbered_list_item sample\n\n" +
				"heading_2 次回\n" +
				"次回の日程:\n" +
				"書記: sample, リード: sample, 計測係: sample\n" +
				"```\n",
		},
		nil,
		nil,
	)
	img := slack.NewImageBlock(
		"http://placekitten.com/500/500", // TODO 写真をnotionの画像に差し替える
		"image",
		"image",
		nil,
	)

	msg := map[string]interface{}{
		"blocks": []slack.Block{
			target,
			headerCanDo,
			descriptionCanDo,
			headerHowToUse,
			descriptionHowToUse,
			headerHowToCreatePage,
			descriptionHowToCreatePage,
			descriptionExample,
			img,
		},
	}
	return msg
}

func (cu *CommandUsecase) All() (string, error) {
	userIDs := make([]int, 0, len(cu.ProcessManager))
	for key, _ := range cu.ProcessManager {
		userIDs = append(userIDs, key)
	}

	nameList, err := cu.DBOperator.GetUserNameList(userIDs)
	if err != nil {
		return "", fmt.Errorf("get user name list error: %+v", err)
	}
	return strings.Join(nameList, "\n"), nil
}

func (cu *CommandUsecase) AllStop() {
	for key, _ := range cu.ProcessManager {
		cu.Stop(key)
	}
}
