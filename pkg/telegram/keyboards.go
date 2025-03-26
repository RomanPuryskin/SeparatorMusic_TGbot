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

/*
func replyButtonReturnKeyBoard() tgbotapi.ReplyKeyboardMarkup {
	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(replyButtonReturn),
		),
	)

	return replyMarkup
}*/

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
			tgbotapi.NewInlineKeyboardButtonData("wav", wavSeparate),
			tgbotapi.NewInlineKeyboardButtonData("mp3", mp3Separate),
			tgbotapi.NewInlineKeyboardButtonData("flac", flacSeparate),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("aac", aacSeparate),
			tgbotapi.NewInlineKeyboardButtonData("ogg", oggSeparate),
			tgbotapi.NewInlineKeyboardButtonData("m4a", m4aSeparate),
		),
	)
	return replyMarkup
}
