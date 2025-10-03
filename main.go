package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

func init() {
	// 设置全局时区为 UTC+8
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("UTC+8", 8*60*60)
	}
	time.Local = loc
}

func main() {
	goutils.InitZeroLog()

	kong.Parse(&cli)
	log.Info().Interface("cli", cli).Msg("Starting histmon")
	// 如果执行时间小于 5 秒，则不发送消息
	if cli.EndAt-cli.StartAt < 5000 {
		log.Info().Msg("Command execution time is less than 5 second, no message will be sent")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get hostname")
	}

	output := cli.Output
	if len(output) > 1000 {
		output = output[:1000] + "..."
	}

	content := fmt.Sprintf(`# 命令执行完成
- **主机名**: %s
- **命令**: %s
- **退出状态码**: %d
- **开始时间**: %s
- **结束时间**: %s
- **执行时长**: %s
- **输出**: %s`,
		hostname,
		cli.Command,
		cli.ExitStatus,
		time.UnixMilli(int64(cli.StartAt)).Format("2006-01-02 15:04:05"),
		time.UnixMilli(int64(cli.EndAt)).Format("2006-01-02 15:04:05"),
		goutils.DurationToStr(time.Duration(cli.EndAt-cli.StartAt)*time.Millisecond),
		output,
	)

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
