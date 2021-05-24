package config

import "os"

func CORSAllowOrigin() string {
	return getterEnvInfo("CORS_ALLOW_ORIGIN")
}

func DatabaseID() string {
	return getterEnvInfo("DATABASE_ID")
}

func NotionToken() string {
	return getterEnvInfo("NOTION_TOKEN")
}

func Port() string {
	return getterEnvInfo("PORT")
}

func DSN() string {
	return getterEnvInfo("DSN")
}

func getterEnvInfo(key string) string {
	return os.Getenv(key)
}
