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

func SALCK_TOKEN() string {
	return getterEnvInfo("SALCK_TOKEN")
}

func WEBHOOK_URL() string {
	return getterEnvInfo("WEBHOOK_URL")
}

func getterEnvInfo(key string) string {
	return os.Getenv(key)
}
