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
		"このアプリでは、notionに議事録のテンプレートページを定期的に自動生成するスケジューラを起動・停止することができます。\n" +
		"スケジューラを動かすためには、以下の手順を行います。\n" +
		"```1. ショートカット'Register the notion info'を選択し、表示されるモーダルにnotion情報を登録します。\n" +
		"2. /startというslash commandを呼び出すことで、スケジューラが起動します。```\n\n" +
		"スケジューラを停止させるためには、/stopを実行すればスケジューラは停止します。\n" +
		"また、notion情報を更新する際には、再度'1'の手順を実行してください。\n" +
		"※1ユーザにつき1スケジューラであるために、現時点では複数台のスケジューラを起動させることができません。" +
		"そのような機能が必要であれば、管理人に連絡してください。"

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
