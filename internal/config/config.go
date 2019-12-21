package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// env 実行環境
var env *Env

// init 実行環境の判定、設定ファイルのロードを行う
func init() {
	// タイムゾーンの設定
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	time.Local = loc

	// 実行環境の設定
	goEnv := Env(os.Getenv("GO_ENV"))
	if !goEnv.IsInitialized() {
		msg := fmt.Sprintf("Env is not initialized by GO_ENV=%v, env=%v", os.Getenv("GO_ENV"), string(goEnv))
		panic(msg)
	}
	env = &goEnv

	// 一般設定ファイルのロード
	viper.SetConfigName(env.ToString())
	viper.AddConfigPath("./config/")
	// vscode での test 用
	configDir := fmt.Sprintf("%s/config/", os.Getenv("PROJECT_ROOT"))
	viper.AddConfigPath(configDir)
	if err := viper.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("Error reading config file, %s", err)
		panic(msg)
	}
}

// GetEnv 実行環境を返す
func GetEnv() *Env {
	return env
}

// GetString は指定されたpathの設定値を返します
func GetString(path string) string {
	return viper.GetString(path)
}
