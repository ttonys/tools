package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v2"
	"gopkg.in/yaml.v2"
)

var configPath = filepath.Join(os.Getenv("HOME"), ".config/glark/config.yaml")

type Config struct {
	Lark []LarkConfig `yaml:"lark"`
}

type LarkConfig struct {
	ID     string `yaml:"id"`
	BotKey string `yaml:"botKey"`
}

type TextPusher interface {
	PushText(string) error
	PushMarkdown(string, string) error
}

type Lark struct {
	bot  *lark.CustomerBot
	sign string
}

func NewLark(botKey, sign string) TextPusher {
	if !strings.HasPrefix(botKey, "http") {
		botKey = "https://open.feishu.cn/open-apis/bot/v2/hook/" + botKey
	}
	bot := lark.NewCustomerBot(botKey, sign)
	return &Lark{
		bot:  bot,
		sign: sign,
	}
}

func (d *Lark) PushText(s string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	msg := lark.MessageText{Text: s}
	resp, err := d.bot.SendMessage(ctx, "text", msg)
	if err != nil {
		return fmt.Errorf("failed to send lark text, %s", err)
	}
	if resp.CodeError.Code != 0 {
		return fmt.Errorf("failed to send lark text, %v", resp.CodeError)
	}
	return nil
}

func (d *Lark) PushMarkdown(title, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	title = strings.ReplaceAll(title, "&nbsp;", "")
	content = strings.ReplaceAll(content, "&nbsp;", "")
	msg := &lark.MessageCardDiv{
		Text: &lark.MessageCardLarkMd{Content: content},
	}
	card := lark.MessageCard{
		Header: &lark.MessageCardHeader{
			Title: &lark.MessageCardPlainText{Content: title},
		},
		Elements: []lark.MessageCardElement{msg},
	}
	resp, err := d.bot.SendMessage(ctx, "interactive", card)
	if err != nil {
		return fmt.Errorf("failed to send lark markdown, %s", err)
	}
	if resp.CodeError.Code != 0 {
		return fmt.Errorf("failed to send lark markdown, %v", resp.CodeError)
	}
	return nil
}

func loadConfig() (Config, error) {
	var config Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return config, err
		}
		defaultConfig := Config{
			Lark: []LarkConfig{
				{ID: "subs", BotKey: "xxxxxx"},
				{ID: "scan", BotKey: "xxxxxx"},
			},
		}
		data, _ := yaml.Marshal(&defaultConfig)
		os.WriteFile(configPath, data, 0644)
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

func main() {
	id := flag.String("id", "", "The ID of the Lark bot to use")
	title := flag.String("t", "title", "The title of the message")
	msg := flag.String("msg", "", "The message to send")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	var botKey string
	for _, larkConfig := range config.Lark {
		if larkConfig.ID == *id {
			botKey = larkConfig.BotKey
			break
		}
	}

	if botKey == "" {
		fmt.Printf("No bot found with ID: %s\n", *id)
		return
	}

	pusher := NewLark(botKey, "")

	if *msg != "" {
		err = pusher.PushMarkdown(*title, *msg)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	} else {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading stdin: %v\n", err)
			return
		}
		err = pusher.PushMarkdown(*title, string(input))
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}
	}
}
