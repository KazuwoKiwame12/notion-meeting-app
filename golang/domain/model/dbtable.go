package model

import "time"

type Workspace struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type User struct {
	ID              int
	SlackUserID     string
	WorkspaceID     string
	IsAdministrator string
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Notion struct {
	ID                int
	UserID            string
	Date              int
	NotionToken       string
	NotionDatabaseID  string
	NotionPageContent string
}
