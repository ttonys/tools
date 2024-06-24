package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/larksuite/oapi-sdk-go/v2"
	"gopkg.in/yaml.v2"
)

var (
	configFile = flag.String("config", "~/.config/glark/config.yaml", "Path to configuration file")
	id         = flag.String("id", "", "ID of the Lark bot to use")
	msg        = flag.String("msg", "", "Message to send")
)

type Config struct {
	Lark []struct {
		ID     string `yaml:"id"`
		BotKey string `yaml:"botKey"`
	} `yaml:"glark"`
}

type TextPusher interface {
	PushText(s string) error
	PushMarkdown(title, content string) error
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

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %s", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config file: %s", err)
	}
	return cfg, nil
}

func findBotByID(cfg Config, id string) (string, error) {
	for _, bot := range cfg.Lark {
		if bot.ID == id {
			return bot.BotKey, nil
		}
	}
	return "", fmt.Errorf("bot with ID '%s' not found in config", id)
}

func main() {
	flag.Parse()

	if *id == "" {
		log.Fatal("ID must be provided")
	}

	if *msg == "" {
		// Check if there is input piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Failed to read from stdin: %s", err)
			}
			*msg = string(bytes)
		} else {
			log.Fatal("Message must be provided via -msg or piped input")
		}
	}

	configPath := os.ExpandEnv(*configFile)
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	botKey, err := findBotByID(cfg, *id)
	if err != nil {
		log.Fatalf("Failed to find bot ID: %s", err)
	}

	larkBot := NewLark(botKey, "")

	err = larkBot.PushText(*msg)
	if err != nil {
		log.Fatalf("Failed to push message: %s", err)
	}

	fmt.Println("Message sent successfully")
}
