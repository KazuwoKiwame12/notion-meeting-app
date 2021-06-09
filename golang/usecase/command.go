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
	text := fmt.Sprintf("@%s\n", name) +
		"## notion-meeting-app\n"+
		"### できること\n"+
		"1. notionに議事録のテンプレートページを定期的に自動作成できる\n"+
		"2. 作成するテンプレートページの内容・場所・間隔を簡単にカスタマイズできる\n"+
		"3. いつでも自動作成するのを停止できます。自動作成の軌道も同様に簡単にできる\n"+
		"### 使い方\n"+
		"1. ショートカット``Register the notion info`を選択する`\n"+
		"2. 表示されたmodalの項目を埋める\n"+
		"```"+
		"1. notionのtoken\n"+
		"2. notionのページを作成する場所\n"+
		"3. notionのテンプレートページの内容\n"+
		"4. テンプレートページを作成する曜日\n"+
		"```\n\n"+
		"3. ``/start``のslash commandでスケジューラを起動させる\n"+
		"※スケジューラを停止させたい場合は``/stop``のslash commandで停止させる\n"+
		"**上記のテンプレートページの内容の作り方**\n"+
		"以下に例を示す。\n"+
		"```"+
		"title ミーティング\n"+
		"heading_2 概要\n"+
		"タイトル:\n\n"+
		"heading_2 内容\n"+
		"numbered_list_item sample\n\n"+
		"numbered_list_item sample\n\n"+
		"heading_2 次回\n"+
		"次回の日程:\n"+
		"書記: sample, リード: sample, 計測係: sample\n"+
		"```\n\n"+
		"1行目は、ページのタイトルに該当し、2行目以降は、ページの中身に該当する。\n"+
		"番号付きリスト・todoリスト・リスト・heading1・heading2・heading3・テキスト・トグルなどのテキストタイプが使えます。\n"+
		"上記を使うためには、以下の形式で記入する必要があります。\n"+
		"```"+
		"`テキストタイプ` `テキスト`..."+
		"```\n\n"+
		"以下にテキストタイプの種類を示す。\n"+
		"- paragraph: 文章\n"+
		"- heading_1: 1番大きいheader\n"+
		"- heading_2: 2番大きいheader\n"+
		"- heading_3: 3番大きいheader\n"+
		"- bulleted_list_item: リスト\n"+
		"- numbered_list_item: 番号つきリスト\n"+
		"- to_do: todoリスト\n"+
		"- toggle: toggleリスト\n"+
		"ページ内容の例と実際のnotionのページの比較図を以下に掲載します。\n"+
		"### コマンド一覧\n"+
		"1. ``/explain``...アプリの使い方の説明を得る\n"+
		"2. ショートカット``Register the notion info``...使用するnotion関連の情報を登録・更新する\n"+
		"3. ``/start``...議事録のテンプレートページを定期的に自動作成するスケジューラを起動する\n"+
		"4. ``/stop``...スケジューラを停止する\n"+
		"5. ``/all``...スケジューラを起動しているユーザを確認する(admin用)\n"+
		"6. ``/all-stop``...全てのスケジューラを停止する(admin用)\n"+

	msg := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: text,
				},
				nil,
				nil,
			),
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
