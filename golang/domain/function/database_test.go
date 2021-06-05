package function_test

import (
	"app/domain/function"
	"app/domain/model"
	"app/infrastructure"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

var w *model.Workspace = &model.Workspace{ID: "W100", Name: "test1"}
var u *model.User = &model.User{SlackUserID: "U100", WorkspaceID: w.ID, IsAdministrator: true, Name: "sample100"}
var n *model.Notion = &model.Notion{UserID: u.ID, Date: 2, NotionToken: []byte("tokenNo1"), NotionDatabaseID: []byte("DatabaseIDNo1"), NotionPageContent: "sample"}

// TODO テストコードの記入
func CreateData(sh *infrastructure.SqlHandler) {
	queryForWorkspace := "INSERT INTO t_workspace(id, name) VALUES($1, $2)"
	queryForUser := "INSERT INTO t_user(id, slack_user_id, t_workspace_id, is_administrator, name) VALUES($1, $2, $3, $4, $5)"
	queryForNotion := "INSERT INTO t_notion(t_user_id, date, notion_token, notion_database_id, notion_page_content) VALUES($1, $2, $3, $4, $5)"
	if _, err := sh.Execute(queryForWorkspace, w.ID, w.Name); err != nil {
		fmt.Printf("insert workspace query err: %+v\n", err)
	}
	if _, err := sh.Execute(queryForUser, u.ID, u.SlackUserID, u.WorkspaceID, u.IsAdministrator, u.Name); err != nil {
		fmt.Printf("insert user query err: %+v\n", err)
	}
	if _, err := sh.Execute(queryForNotion, n.UserID, n.Date, n.NotionToken, n.NotionDatabaseID, n.NotionPageContent); err != nil {
		fmt.Printf("insert notion query err: %+v\n", err)
	}
}

func DeleteAllData(sh *infrastructure.SqlHandler) {
	if _, err := sh.Execute("DELETE FROM t_workspace"); err != nil {
		fmt.Printf("delete workspace err: %+v\n", err)
	}
	if _, err := sh.Execute("DELETE FROM t_user"); err != nil {
		fmt.Printf("delete user err: %+v\n", err)
	}
	if _, err := sh.Execute("DELETE FROM t_notion"); err != nil {
		fmt.Printf("delete notion err: %+v\n", err)
	}
}

func getDbOperatorInstance() *function.DatabaseOperater {
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		fmt.Printf("newSqlHandler err: %+v", err)
	}
	// defer func() {
	// 	if err := sh.DB.Close(); err != nil {
	// 		fmt.Printf("closed err: %+v", err)
	// 	}
	// }()
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}
	return dbOp
}

func hasSameRecordOfNotion(result *model.Notion, want *model.Notion) bool {
	isSameUser := result.UserID == want.UserID
	isSameDate := result.Date == want.Date
	isSameToken := reflect.DeepEqual(result.NotionToken, want.NotionToken)
	isSameDatabaseID := reflect.DeepEqual(result.NotionDatabaseID, want.NotionDatabaseID)
	isSameContent := result.NotionPageContent == want.NotionPageContent
	return isSameUser && isSameDate && isSameToken && isSameDatabaseID && isSameContent
}

func TestMain(m *testing.M) {
	fmt.Println("before all...")
	envfilePath := fmt.Sprintf("%s/src/app/.env", os.Getenv("GOPATH"))
	if err := godotenv.Load(envfilePath); err != nil {
		fmt.Printf("%v\n", err)
	}
	sh, err := infrastructure.NewSqlHandler()
	defer func() {
		if err := sh.DB.Close(); err != nil {
			fmt.Printf("closed err: %+v", err)
		}
	}()

	if err != nil {
		fmt.Printf("newSqlHandler err: %+v", err)
	}

	fmt.Println("delete all table data and create data")
	DeleteAllData(sh)
	CreateData(sh)

	fmt.Println("test start")
	code := m.Run()

	fmt.Println("after all...")
	DeleteAllData(sh)
	os.Exit(code)
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID string
		slackUserID string
		want        *model.User
	}{
		{
			name:        "userのデータを取得する",
			workspaceID: w.ID,
			slackUserID: u.SlackUserID,
			want:        u,
		},
	}

	hasSameRecord := func(result *model.User, want *model.User) bool {
		isSameID := result.Name == want.Name
		isSameName := result.Name == want.Name
		isSameSlackUser := result.SlackUserID == want.SlackUserID
		isSameWorkspace := result.WorkspaceID == want.WorkspaceID
		isSameAd := result.IsAdministrator == want.IsAdministrator
		return isSameID && isSameName && isSameSlackUser && isSameWorkspace && isSameAd
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dbOp.GetUser(test.workspaceID, test.slackUserID)
			if err != nil || !hasSameRecord(result, test.want) {
				t.Errorf("get user error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
		t.Run(test.name+"administoratorバージョン", func(t *testing.T) {
			result, err := dbOp.GetAdministrator()
			if err != nil || !hasSameRecord(result, test.want) {
				t.Errorf("get user error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
	}
}

func TestGetNotionInfo(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   *model.Notion
	}{
		{
			name:   "notionのデータを取得する",
			userID: u.ID,
			want:   n,
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dbOp.GetNotionInfo(test.userID)
			if err != nil || !hasSameRecordOfNotion(result, test.want) {
				t.Errorf("get notion error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
	}
}

func TestRegisterWorkspace(t *testing.T) {
	tests := []struct {
		name       string
		insertData *model.Workspace
	}{
		{
			name:       "workspace tableにデータを登録する",
			insertData: &model.Workspace{ID: "W101", Name: "test2"},
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterWorkspace(test.insertData); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	users := []model.User{
		{
			SlackUserID:     "U101",
			WorkspaceID:     w.ID,
			IsAdministrator: false,
			Name:            "sample101",
		},
		{
			SlackUserID:     "U102",
			WorkspaceID:     w.ID,
			IsAdministrator: false,
			Name:            "sample102",
		},
		{
			SlackUserID:     "U103",
			WorkspaceID:     w.ID,
			IsAdministrator: false,
			Name:            "sample103",
		},
	}

	tests := []struct {
		name       string
		insertData []model.User
	}{
		{
			name:       "user tableにデータを登録する",
			insertData: users,
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterUsers(test.insertData); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
			for index, data := range test.insertData {
				result, _ := dbOp.GetUser(data.WorkspaceID, data.SlackUserID)
				if result.SlackUserID != data.SlackUserID {
					t.Errorf("insert wrong users[%d]' properties...result is %s\nwant is %s\n", index, result.SlackUserID, data.SlackUserID)
				}
			}
		})
	}
}

func TestRegisterNotion(t *testing.T) {
	tests := []struct {
		name       string
		insertData *model.Notion
	}{
		{
			name:       "notion tableにデータを登録する",
			insertData: &model.Notion{UserID: u.ID, Date: 3, NotionToken: []byte("tokenNo2"), NotionDatabaseID: []byte("DatabaseIDNo2"), NotionPageContent: "sample"},
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterNotionInfo(test.insertData); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
			result, _ := dbOp.GetNotionInfo(test.insertData.UserID)
			if !hasSameRecordOfNotion(result, test.insertData) {
				t.Errorf("insert wrong notion's properties...result is %+v\nwant is %+v\n", result, test.insertData)
			}
		})
	}
}

func TestUpdateNotionInfo(t *testing.T) {
	dbOp := getDbOperatorInstance()

	type kindOfUpdate int
	const (
		token kindOfUpdate = iota
		databaseID
		pageContent
	)
	tests := []struct {
		name       string
		updateData string
		userID     int
		kind       kindOfUpdate
		function   func(string, int) error
	}{
		{
			name:       "notionのtokenを更新する",
			updateData: "updated_token",
			userID:     n.UserID,
			kind:       token,
			function:   dbOp.UpdateNotionToken,
		},
		{
			name:       "notionのdatabase_idを更新する",
			updateData: "updated_dataabseID",
			userID:     n.UserID,
			kind:       databaseID,
			function:   dbOp.UpdateNotionDatabaseID,
		},
		{
			name:       "notionのpage contentを更新する",
			updateData: "updated_page_content",
			userID:     n.UserID,
			kind:       pageContent,
			function:   dbOp.UpdateNotionPageContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.function(test.updateData, test.userID); err != nil {
				t.Errorf("update notion table error ocurred: %+v\n", err)
			}
			result, _ := dbOp.GetNotionInfo(test.userID)
			var sameRecord bool
			switch test.kind {
			case token:
				sameRecord = reflect.DeepEqual(result.NotionToken, test.updateData)
			case databaseID:
				sameRecord = reflect.DeepEqual(result.NotionDatabaseID, test.updateData)
			case pageContent:
				sameRecord = result.NotionPageContent == test.updateData
			}
			if !sameRecord {
				t.Errorf("update wrong notion's property...result is %+v\nwant is %s\n", result, test.updateData)
			}
		})
	}
}

func TestUpdateNotionDate(t *testing.T) {
	tests := []struct {
		name       string
		updateData int
		userID     int
	}{
		{
			name:       "notionのdateを更新する",
			updateData: 1,
			userID:     n.UserID,
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.UpdateDateToCreate(test.updateData, test.userID); err != nil {
				t.Errorf("update notion table error ocurred: %+v\n", err)
			}

			result, _ := dbOp.GetNotionInfo(test.userID)
			if result.Date != test.updateData {
				t.Errorf("update wrong notion's property...result is %d\nwant is %d\n", result.Date, test.updateData)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name   string
		userID int
	}{
		{
			name:   "notionのdateを更新する",
			userID: n.UserID,
		},
	}

	dbOp := getDbOperatorInstance()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RemoveaNotionToken(test.userID); err != nil {
				t.Errorf("remove notion token error ocurred: %+v\n", err)
			}

			result, _ := dbOp.GetNotionInfo(test.userID)
			if len(result.NotionToken) == 0 {
				t.Errorf("can't remove notion's token...result is %s\nwant is '' ", result.NotionToken)
			}
		})
	}
}
