package chat
import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)
func initBotAPI(tgToken string) (*api.BotAPI, error) {
    bot, err  := api.NewBotAPI(tgToken)
    if err != nil {
        return nil, err
    }
    log.Printf("Authorized on account %s \n", bot.Self.UserName)
    return bot, nil
}