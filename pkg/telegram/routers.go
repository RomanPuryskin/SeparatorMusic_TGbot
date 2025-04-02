package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// команды бота:
const (
	// команды
	commandStart = "start" // главная команда бота

	// reply кнопки  (по сути будем обрабатывать как команды, но...)
	// (нажатие reply кнопки сопровождается Update.Message, как обычный текст, поэтому обрабатываем
	// эти события как обычный текст)
	replyButtonSeparateAudio = "Separate audio" // кнопка для перехода в состояние waitingAudioForSeparateState
	replyButtonChangeFormat  = "Change format"  // кнопка для перехода в состояни choosingFormatAudioForChangeState
	replyButtonReturn        = "Return"         // кнопка для возврата состояния в предыдущее значение
	replyButtonMainMenu      = "Main menu"      //кнопка для возврата в главное меню
)

// возможные состояния:
const (
	mainMenuState = "main_menu" // пользователь в главном меню и может выбрать дальнейшие действия

	choosingFormatAudioForSeparateState = "choosing_format_audio_for_separate" // пользовательно выбрал действие separate audio и теперь бот ждет от него выбора формата

	waitingAudioForSeparateState = "waiting_audio_for_separate" // пользователь выбрал действие separate audio и бот
	// находится в ожидании аудио

	choosingFormatAudioForChangeState = "choosing_format_audio_for_change" // пользователь выбрал действие change format и теперь бот ждет от него выбора формата

	waitingAudioForChangeFormatState = "waiting_audio_for_change_format" // пользователь выбрал формат аудио, в который хочет переформатировать аудио и теперь бот ждет от него аудио
)

// inline кнопки
const (
	// для опции separate audio
	wavSeparate  = "wav_separate"
	mp3Separate  = "mp3_separate"
	flacSeparate = "flac_separate"
	aacSeparate  = "aac_separate"
	oggSeparate  = "ogg_separate"
	m4aSeparate  = "m4a_separate"

	// для опции change format
	wavChangeFormat  = "wav_change_format"
	mp3ChangeFormat  = "mp3_change_format"
	flacChangeFormat = "flac_change_format"
	aacChangeFormat  = "aac_change_format"
	oggChangeFormat  = "ogg_change_format"
	m4aChangeFormat  = "m4a_change_format"
)

func commandsRouter(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	switch msg.Command() {
	case commandStart:
		handlerCommandStart(bot, msg)
	default:
		requestMsg := tgbotapi.NewMessage(msg.Chat.ID, "Я не знаю такой команды")
		requestMsg.ReplyToMessageID = msg.MessageID

		_, err := bot.Send(requestMsg)
		if err != nil {
			log.Println("Ошибка отправки сообщения в commandsRouter , default:", err)
		}
	}
}

func statesRouter(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	switch getCurrentState(msg.Chat.ID) {
	case choosingFormatAudioForSeparateState:
		handlerChoosingFormatAudioForSeparateState(bot, msg)
	case waitingAudioForSeparateState:
		handlerWaitingAudioForSeparateState(bot, msg)
	case choosingFormatAudioForChangeState:
		handlerChoosingFormatAudioForChangeState(bot, msg)
	case waitingAudioForChangeFormatState:
		handlerWaitingAudioForChangeFormatState(bot, msg)
	default:
		handlerMainMenuState(bot, msg)
	}
}

func inlineButtonsRouter(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {

	if callback.Data == wavSeparate || callback.Data == mp3Separate ||
		callback.Data == flacSeparate || callback.Data == oggSeparate ||
		callback.Data == aacSeparate || callback.Data == m4aSeparate {
		setNewState(callback.Message.Chat.ID, waitingAudioForSeparateState)
		handlerWaitingAudioForSeparateState(bot, callback.Message)
		setNewCurrentFormatForSeparate(callback.Message.Chat.ID, callback.Data)
	}

	if callback.Data == wavChangeFormat || callback.Data == mp3ChangeFormat ||
		callback.Data == flacChangeFormat || callback.Data == oggChangeFormat ||
		callback.Data == aacChangeFormat || callback.Data == m4aChangeFormat {
		setNewState(callback.Message.Chat.ID, waitingAudioForChangeFormatState)
		handlerWaitingAudioForChangeFormatState(bot, callback.Message)
		setNewCurrentFormatForChange(callback.Message.Chat.ID, callback.Data)
	}

}
