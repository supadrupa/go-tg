package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mr-linch/go-tg"
)

var (
	token = os.Getenv("BOT_TOKEN")
)

func init() {
	flag.StringVar(&token, "token", token, "Telegram Bot API token [$BOT_TOKEN]")
	flag.Parse()

	if token == "" {
		fmt.Println("token is required, provide it via -token or $BOT_TOKEN")
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	client := tg.NewClient(token)

	me, err := client.GetMe(ctx)
	if err != nil {
		fmt.Printf("call error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bot: %+v\n", me)
}
