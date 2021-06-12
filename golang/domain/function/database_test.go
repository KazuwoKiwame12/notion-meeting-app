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

/*
*test用の関数
 */
func createData(sh *infrastructure.SqlHandler) {
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

func deleteAllData(sh *infrastructure.SqlHandler) {
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
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}
	return dbOp
}

func hasSameUserRecord(result *model.User, want *model.User) bool {
	isSameID := result.ID == want.ID
	isSameName := result.Name == want.Name
	isSameSlackUser := result.SlackUserID == want.SlackUserID
	isSameWorkspace := result.WorkspaceID == want.WorkspaceID
	isSameAd := result.IsAdministrator == want.IsAdministrator
	return isSameID && isSameName && isSameSlackUser && isSameWorkspace && isSameAd
}

func hasSameNotionRecord(result *model.Notion, want *model.Notion) bool {
	isSameUser := result.UserID == want.UserID
	isSameDate := result.Date == want.Date
	isSameToken := reflect.DeepEqual(result.NotionToken, want.NotionToken)
	isSameDatabaseID := reflect.DeepEqual(result.NotionDatabaseID, want.NotionDatabaseID)
	isSameContent := result.NotionPageContent == want.NotionPageContent
	return isSameUser && isSameDate && isSameToken && isSameDatabaseID && isSameContent
}

/*
*全てのテストを実行する
 */
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

	fmt.Println("test start")
	code := m.Run()

	fmt.Println("after all...")
	deleteAllData(sh)
	os.Exit(code)
}

/*
*t_workspaceテーブルに対する操作のテスト
 */
func TestRegisterWorkspace(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

	tests := []struct {
		name       string
		insertData *model.Workspace
	}{
		{
			name:       "workspace tableにデータを登録する",
			insertData: w,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterWorkspace(test.insertData); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
		})
	}
}

/*
*t_userテーブルに対する操作のテスト
 */

func TestGetUser(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dbOp.GetUser(test.workspaceID, test.slackUserID)
			if err != nil || !hasSameUserRecord(result, test.want) {
				t.Errorf("get user error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
	}
}

func TestGetAdminUser(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

	tests := []struct {
		name        string
		workspaceID string
		slackUserID string
		want        *model.User
	}{
		{
			name:        "admin_userのデータを取得する",
			workspaceID: w.ID,
			slackUserID: u.SlackUserID,
			want:        u,
		},
	}

	for _, test := range tests {
		t.Run(test.name+"administoratorバージョン", func(t *testing.T) {
			result, err := dbOp.GetAdministrator()
			if err != nil || !hasSameUserRecord(result, test.want) {
				t.Errorf("get user error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
	}
}

func TestGetUserNameList(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

	tests := []struct {
		name string
		ids  []int
		want []string
	}{
		{
			name: "userの名前を一括取得する",
			ids:  []int{u.ID},
			want: []string{u.Name},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nameList, err := dbOp.GetUserNameList(test.ids)
			if err != nil {
				t.Errorf("get user name list error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, nameList, test.want)
			}

			for index, name := range nameList {
				if test.want[index] != name {
					t.Errorf("Failed to get the correct name: %+v\nresult is %+v\nwant is %+v\n", err, name, test.want[index])
				}
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

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

/*
*t_notionテーブルに対する操作のテスト
 */
func TestGetNotionInfo(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := dbOp.GetNotionInfo(test.userID)
			if err != nil || !hasSameNotionRecord(result, test.want) {
				t.Errorf("get notion error ocurred: %+v\nresult is %+v\nwant is %+v\n", err, result, test.want)
			}
		})
	}
}

func TestRegisterNotion(t *testing.T) {
	dbOp := getDbOperatorInstance()
	deleteAllData(dbOp.SqlHandler.(*infrastructure.SqlHandler))
	createData(dbOp.SqlHandler.(*infrastructure.SqlHandler))

	tests := []struct {
		name       string
		insertData *model.Notion
	}{
		{
			name:       "notion tableにデータを登録する",
			insertData: n,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterNotionInfo(test.insertData); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
			result, _ := dbOp.GetNotionInfo(test.insertData.UserID)
			if !hasSameNotionRecord(result, test.insertData) {
				t.Errorf("insert wrong notion's properties...result is %+v\nwant is %+v\n", result, test.insertData)
			}
		})
	}
}
