package config_test

import (
	"app/config"
	"os"
	"testing"
)

func TestCanGetEnvInfo(t *testing.T) {

	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "環境変数が設定されている場合にしっかりと取得できるかのテスト",
			key:  "TEST_DATA",
			want: "this is a test data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.key, test.want)
			if config.ExportGetterEnvInfo(test.key) != test.want {
				t.Error("configのgetterEnvInfo関数は誤った動作をしています")
			}
		})
	}
}
