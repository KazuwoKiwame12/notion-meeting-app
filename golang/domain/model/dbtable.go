package model

import "time"

type User struct {
	ID              int
	SlackUserID     string
	WorkspaceID     string
	IsAdministrator string
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Workspace struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
