package config

import "os"

func CORSAllowOrigin() string {
	return os.Getenv("CORS_ALLOW_ORIGIN")
}

func DatabaseID() string {
	return os.Getenv("DATABASE_ID")
}

func NotionToken() string {
	return os.Getenv("NOTION_TOKEN")
}

func Port() string {
	return os.Getenv("PORT")
}
