package config

import "os"

func CORSAllowOrigin() string {
	return getterEnvInfo("CORS_ALLOW_ORIGIN")
}

func Port() string {
	return getterEnvInfo("PORT")
}

func DSN() string {
	return getterEnvInfo("DSN")
}

func SLACK_TOKEN() string {
	return getterEnvInfo("SLACK_TOKEN")
}

func WEBHOOK_URL() string {
	return getterEnvInfo("WEBHOOK_URL")
}

func NOTION_API_VERSION() string {
	return getterEnvInfo("NOTION_API_VERSION")
}

func NOTION_API_URL() string {
	return getterEnvInfo("NOTION_API_URL")
}

func getterEnvInfo(key string) string {
	return os.Getenv(key)
}
