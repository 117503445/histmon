package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/117503445/goutils"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

var cli struct {
	Command    string `env:"COMMAND"`
	Output     string `env:"OUTPUT"`
	StartAt    int    `env:"START_AT"` // 毫秒时间戳
	EndAt      int    `env:"END_AT"`   // 毫秒时间戳
	ExitStatus int    `env:"EXIT_STATUS"`

	Token    string `env:"TOKEN"`
	Endpoint string `env:"ENDPOINT"`
}

func main() {
	goutils.InitZeroLog()

	kong.Parse(&cli)
	log.Info().Interface("cli", cli).Msg("Starting histmon")

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get hostname")
	}

	content := fmt.Sprintf(`# 命令执行完成
- hostname: %s`, hostname)

	{
		payload := map[string]any{
			"title":       "命令执行完成",
			"description": "命令执行完成",
			"content":     content,
			"channel":     "ding",
			"token":       cli.Token,
		}

		// 将数据转换为 JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to marshal JSON")
			return
		}

		// 创建 HTTP 请求
		resp, err := http.Post(
			cli.Endpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to send message")
			return
		}
		defer resp.Body.Close()
	}

}
