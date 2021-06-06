package function

import (
	"app/domain/adapter"
	"app/domain/model"
	"fmt"
	"strings"
)

type DatabaseOperater struct {
	SqlHandler adapter.SqlHandlerAdapter
}

func (do *DatabaseOperater) RegisterWorkspace(workspace *model.Workspace) error {
	query := "INSERT INTO t_workspace(id, name) values ($1, $2) ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name"
	if _, err := do.SqlHandler.Execute(query, workspace.ID, workspace.Name); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) RegisterUsers(users []model.User) error {
	convertDataForExec := func(us []model.User, numOfColumns int) []interface{} {
		bind := make([]interface{}, len(us)*numOfColumns)
		var index int = 0
		for _, u := range us {
			bind[index] = u.SlackUserID
			bind[index+1] = u.WorkspaceID
			bind[index+2] = u.IsAdministrator
			bind[index+3] = u.Name
			index += numOfColumns
		}
		return bind
	}

	queryForUpdatePart := "UPDATE SET is_administrator = EXCLUDED.is_administrator, name = EXCLUDED.name"
	query := fmt.Sprintf("INSERT INTO t_user(slack_user_id, t_workspace_id, is_administrator, name) values %s ON CONFLICT(slack_user_id) DO %s", do.makePlaceHolders(len(users), 4), queryForUpdatePart)
	if _, err := do.SqlHandler.Execute(query, convertDataForExec(users, 4)...); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) GetUser(workspaceID, slackUserID string) (*model.User, error) {
	row := do.SqlHandler.QueryRow("SELECT * FROM t_user WHERE t_workspace_id = $1 AND slack_user_id = $2", workspaceID, slackUserID)

	user := &model.User{}
	if err := row.Scan(
		&user.ID,
		&user.SlackUserID,
		&user.WorkspaceID,
		&user.IsAdministrator,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (do *DatabaseOperater) GetAdministrator() (*model.User, error) {
	row := do.SqlHandler.QueryRow("SELECT * FROM t_user WHERE is_administrator = $1", true)

	user := &model.User{}
	if err := row.Scan(
		&user.ID,
		&user.SlackUserID,
		&user.WorkspaceID,
		&user.IsAdministrator,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (do *DatabaseOperater) RegisterNotionInfo(n *model.Notion) error {
	queryForUpdatePart := "UPDATE SET date = EXCLUDED.date, notion_token = EXCLUDED.notion_token, notion_database_id=EXCLUDED.notion_database_id, notion_page_content=EXCLUDED.notion_page_content"
	query := fmt.Sprintf("INSERT INTO t_notion(t_user_id, date, notion_token, notion_database_id, notion_page_content) values %s ON CONFLICT(t_user_id) DO %s", do.makePlaceHolders(1, 5), queryForUpdatePart)
	if _, err := do.SqlHandler.Execute(query, n.UserID, n.Date, n.NotionToken, n.NotionDatabaseID, n.NotionPageContent); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) GetNotionInfo(userID int) (*model.Notion, error) {
	row := do.SqlHandler.QueryRow("SELECT * FROM t_notion WHERE t_user_id = $1", userID)

	notion := &model.Notion{}
	if err := row.Scan(
		&notion.ID,
		&notion.UserID,
		&notion.Date,
		&notion.NotionToken,
		&notion.NotionDatabaseID,
		&notion.NotionPageContent,
	); err != nil {
		return nil, err
	}
	return notion, nil
}

func (do *DatabaseOperater) makePlaceHolders(numOfInsertValues, numOfColumns int) string {
	placeholdersList := make([]string, numOfInsertValues)
	var count int = 1
	for i := 0; i < numOfInsertValues; i++ {
		placeholders := make([]string, numOfColumns)
		for j := 0; j < numOfColumns; j++ {
			placeholders[j] = fmt.Sprintf("$%d", count)
			count += 1
		}
		joined := strings.Join(placeholders, ",")
		placeholdersList[i] = "(" + joined + ")"
	}
	return strings.Join(placeholdersList, ",")
}
