package function_test

import (
	"app/domain/function"
	"app/domain/model"
	"app/infrastructure"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var w *model.Workspace = &model.Workspace{ID: "W100", Name: "test1"}
var u *model.User = &model.User{SlackUserID: "U100", WorkspaceID: w.ID, IsAdministrator: true, Name: "sample100"}
var n *model.Notion = &model.Notion{Date: 2, NotionToken: "tokenNo1", NotionDatabaseID: "DatabaseIDNo1", NotionPageContent: "sample"}

// TODO テストコードの記入
func CreateData(sh *infrastructure.SqlHandler) {
	queryForWorkspace := "INSERT INTO t_workspace(id, name) VALUES($1, $2)"
	queryForUser := "INSERT INTO t_user(slack_user_id, t_workspace_id, is_administrator, name) VALUES($1, $2, $3, $4)"
	queryForNotion := "INSERT INTO t_notion(t_user_id, date, notion_token, notion_database_id, notion_page_content) VALUES($1, $2, $3, $4, $5)"
	if _, err := sh.Execute(queryForWorkspace, w.ID, w.Name); err != nil {
		fmt.Printf("insert workspace query err: %+v\n", err)
	}
	if _, err := sh.Execute(queryForUser, u.SlackUserID, u.WorkspaceID, u.IsAdministrator, u.Name); err != nil {
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

func TestMain(m *testing.M) {
	fmt.Println("before all...")
	if err := godotenv.Load(); err != nil {
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
	CreateData(sh)
	code := m.Run()

	fmt.Println("after all...")
	DeleteAllData(sh)
	os.Exit(code)
}
func TestRegisterWorkspace(t *testing.T) {
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		fmt.Printf("newSqlHandler err: %+v", err)
	}
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}

	tests := []struct {
		name        string
		insert_data *model.Workspace
	}{
		{
			name:        "workspace tableにデータを登録する",
			insert_data: &model.Workspace{ID: "W101", Name: "test2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterWorkspace(test.insert_data); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		fmt.Printf("newSqlHandler err: %+v", err)
	}
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}

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
		name        string
		insert_data []model.User
	}{
		{
			name:        "user tableにデータを登録する",
			insert_data: users,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterUsers(test.insert_data); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
		})
	}
}

func TestRegisterNotion(t *testing.T) {
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		fmt.Printf("newSqlHandler err: %+v", err)
	}
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}

	tests := []struct {
		name        string
		insert_data *model.Notion
	}{
		{
			name:        "workspaceにデータを登録する",
			insert_data: &model.Notion{UserID: u.ID, Date: 3, NotionToken: "tokenNo2", NotionDatabaseID: "DatabaseIDNo2", NotionPageContent: "sample"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := dbOp.RegisterNotionInfo(test.insert_data); err != nil {
				t.Errorf("insert error ocurred: %+v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestRemove(t *testing.T) {

}
