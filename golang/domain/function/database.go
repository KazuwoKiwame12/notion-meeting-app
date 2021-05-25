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
	query := fmt.Sprintf("INSERT INTO t_workspace(id, name) values %s", do.makePlaceHolders(1, 2))
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

	query := fmt.Sprintf("INSERT INTO t_user(slack_user_id, t_workspace_id, is_administrator, name) values %s", do.makePlaceHolders(len(users), 4))
	if _, err := do.SqlHandler.Execute(query, convertDataForExec(users, 4)...); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) makePlaceHolders(numOfInsertValues, numOfColumns int) string {
	placeholdersList := make([]string, numOfInsertValues)
	for i := 0; i < numOfInsertValues; i++ {
		placeholders := make([]string, numOfColumns)
		for j := 0; j < numOfColumns; j++ {
			placeholders[j] = "?"
		}
		joined := strings.Join(placeholders, ",")
		placeholdersList[i] = "(" + joined + ")"
	}
	return strings.Join(placeholdersList, ",")
}

func (do *DatabaseOperater) GetUser(workspaceID, slackUserID string) (*model.User, error) {
	row := do.SqlHandler.QueryRow("SELECT * FROM t_user WHERE t_workspace_id = ? AND slack_user_id = ?", workspaceID, slackUserID)

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
	row := do.SqlHandler.QueryRow("SELECT * FROM t_user WHERE is_administrator = ?", true)

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
	query := fmt.Sprintf("INSERT INTO t_notion(t_user_id, date, notion_token, notion_database_id, notion_page_content) values %s", do.makePlaceHolders(1, 5))
	if _, err := do.SqlHandler.Execute(query, n.UserID, n.Date, n.NotionToken, n.NotionDatabaseID, n.NotionPageContent); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) UpdateNotionToken(token string) error {
	if _, err := do.SqlHandler.Execute("UPDATE t_notion SET notion_token = ?", token); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) UpdateNotionDatabaseID(databaseID string) error {
	if _, err := do.SqlHandler.Execute("UPDATE t_notion SET notion_database_id = ?", databaseID); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) UpdateNotionPageContent(pageContent string) error {
	if _, err := do.SqlHandler.Execute("UPDATE t_notion SET notion_page_content = ?", pageContent); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) UpdateDateToCreate(date string) error {
	if _, err := do.SqlHandler.Execute("UPDATE t_notion SET date = ?", date); err != nil {
		return err
	}
	return nil
}

func (do *DatabaseOperater) RemoveaNotionToken() error {
	if _, err := do.SqlHandler.Execute("UPDATE t_notion SET notion_token = ?", ""); err != nil {
		return err
	}
	return nil
}
