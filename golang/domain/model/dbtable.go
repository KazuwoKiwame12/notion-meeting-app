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
	IsAdministrator bool
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Notion struct {
	ID                int
	UserID            int
	Date              int
	NotionToken       string
	NotionDatabaseID  string
	NotionPageContent string
}
