package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

// функция для создания экземпляра бота по токену
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{bot: bot}, nil
}

// главная функция для бота, запускает возможность получать обновления
func (b *Bot) Start() {

	// устанавливаем способ получения обновлений (используем long polling)
	b.longPolling()
}

// реализация способа получения обновлений long - polling
func (b *Bot) longPolling() {

	// инициализируем канал обновлений UpdatesChannel с которого будем получать все обновления
	updates := b.initUpdatesChannel()

	for update := range updates {

		if update.Message != nil {
			if update.Message.IsCommand() {
				commandsRouter(b.bot, update.Message)
			} else {
				statesRouter(b.bot, update.Message)
			}
			log.Println(update.Message.Chat.ID, getCurrentState(update.Message.Chat.ID))
		}

		if update.CallbackQuery != nil {
			inlineButtonsRouter(b.bot, update.CallbackQuery)
		}

	}

}

// функция для создания канала обновлений UpdatesChannel (используется только для метода long-polling)
func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	return b.bot.GetUpdatesChan(updateConfig)
}
