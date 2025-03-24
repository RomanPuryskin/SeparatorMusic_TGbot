package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// в этом файле расположены конфигурации кнопок

func commandStartKeyBoard() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(replyButtonSeparateAudio),
			tgbotapi.NewKeyboardButton(replyButtonChangeFormat),
		),
	)

	return replyMarkup
}

func replyButtonReturnKeyBoard() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(replyButtonReturn),
		),
	)

	return replyMarkup
}

func replyButtonReturnAndMainMenuKeyBoard() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(replyButtonReturn),
			tgbotapi.NewKeyboardButton(replyButtonMainMenu),
		),
	)

	return replyMarkup
}

func inlineButtonsForChooseFormateToSeparateKeyBoard() tgbotapi.InlineKeyboardMarkup {
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("wav", "wav_separate"),
			tgbotapi.NewInlineKeyboardButtonData("mp3", "mp3_separate"),
			tgbotapi.NewInlineKeyboardButtonData("flac", "flac_separate"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("wma", "wma_separate"),
			tgbotapi.NewInlineKeyboardButtonData("ogg", "ogg_separate"),
			tgbotapi.NewInlineKeyboardButtonData("m4a", "m4a_separate"),
		),
	)
	return replyMarkup
}

func empty() tgbotapi.ReplyKeyboardMarkup {
	replyMarup := tgbotapi.NewReplyKeyboard()
	return replyMarup
}
