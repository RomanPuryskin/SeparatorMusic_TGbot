package telegram

import (
	"errors"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ------------------------Хэндлеры команд------------------------//

// обработчик для главной команды старт
func handlerCommandStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// установим состояние пользователя
	setNewState(msg.Chat.ID, mainMenuState)

	textMessage :=
		`Привет!
	Я бот для обработки музыки
	Выбери действие, которое тебе нужно:
	Separate audio - отправь мне аудио и я отделю вокал от инструментала
	Change format - отправь мне аудио и укажи в какой формат его преобразовать`

	requestMsg := tgbotapi.NewMessage(msg.Chat.ID, textMessage)

	// добавление кнопок для ответного сообщения на команду старт
	requestMsg.ReplyMarkup = commandStartKeyBoard()

	// отправим сообщение
	_, err := bot.Send(requestMsg)
	if err != nil {
		log.Println("Ошибка отправки сообщения в handlerCommadStart", err)
	}

}

// ------------------------Хэндлеры команд------------------------//

//------------------------Хэндлеры состояний------------------------//

func handlerMainMenuState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Text {
	case replyButtonChangeFormat:
		//ToDo сделать хендлер для этого
	case replyButtonSeparateAudio: // по сути нажали на соответствующую кнопку
		// установим соответствующее состояние пользователя
		//setNewState(msg.Chat.ID, waitingAudioForSeparateState)
		//handlerWaitingAudioForSeparateState(bot, msg)
		setNewState(msg.Chat.ID, choosingFormatAudioForSeparateState)
		handlerChoosingFormatAudioForSeparateState(bot, msg)
		//TODO: если нажата кнопка separate audio перейти в состояние choosingFormatAudioForSeparateState
	default:
		handlerCommandStart(bot, msg)
	}
}

func handlerChoosingFormatAudioForSeparateState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	msgText1 := `Выберите формат, в котором хотите получить аудио`
	requestMsg1 := tgbotapi.NewMessage(msg.Chat.ID, msgText1)
	requestMsg1.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(requestMsg1)

	msgText2 := "Форматы:"
	requestMsg2 := tgbotapi.NewMessage(msg.Chat.ID, msgText2)
	requestMsg2.ReplyMarkup = inlineButtonsForChooseFormateToSeparateKeyBoard()
	bot.Send(requestMsg2)
}

func handlerWaitingAudioForSeparateState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	// если нажата кнопка выхода в главное меню
	if msg.Text == replyButtonMainMenu {
		// изменим состояние и вернем клавиатуру главного меню
		handlerCommandStart(bot, msg)
		return
	}

	// если нажата кнопка Return идем в состояние выбора формата аудио
	if msg.Text == replyButtonReturn {
		// TODO: сделать переход в состояние
		handlerChoosingFormatAudioForSeparateState(bot, msg)
		return
	}

	if msg.Audio != nil {
		msgText := "Идет обработка..."
		requestMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
		requestMsg.ReplyToMessageID = msg.MessageID

		// отправим сообщение
		bot.Send(requestMsg)

		err := handlerAudio(bot, msg)
		if err != nil {
			log.Println(err)
			// ToDo как то обработать ошибку
		}

	} else {
		msgText := `Загрузите аудиофайл для обработки.
		Для выхода в главное меню нажмите Return`

		requestMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)

		// установим соответствующую клавиатуру
		requestMsg.ReplyMarkup = replyButtonReturnAndMainMenuKeyBoard()

		// отправим сообщение
		_, err := bot.Send(requestMsg)
		if err != nil {
			log.Println("Ошибка отправки сообщения в handlerWaitingAudioForSeparateState", err)
		}
	}

}

//------------------------Хэндлеры состояний------------------------//

func handlerAudio(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {

	//получим URL аудиофайла
	audioURL, err := bot.GetFileDirectURL(msg.Audio.FileID)
	if err != nil {
		return errors.New("ошибка при получении файла с сервера")
	}

	inputDir := `.\pkg\spleeter\input`   // Папка для входных файлов
	outputDir := `.\pkg\spleeter\output` // папка для выходных файлов

	audioUniqueName := generateAudioFileName()

	// скачиваем файл
	err = downloadAudio(audioURL, inputDir+`\`+audioUniqueName)
	if err != nil {
		return err
	}

	// получим явно формат аудио, который желает пользователь
	format := ""
	switch getCurrentFormat(msg.Chat.ID) {
	case mp3Separate:
		format = "mp3"
	case flacSeparate:
		format = "flac"
	case aacSeparate:
		format = "aac"
	case oggSeparate:
		format = "ogg"
	case m4aSeparate:
		format = "m4a"
	default:
		format = "wav"
	}
	// Запускаем Spleeter
	err = runSpleeter(outputDir, inputDir, audioUniqueName, format)
	if err != nil {
		return err
	}

	err = sendAudioFiles(audioUniqueName, outputDir, bot, msg)
	if err != nil {
		// если произойдет ошибка отправки по какой то причине, так же почистим файлы
		os.Remove(inputDir + `\` + audioUniqueName)
		os.RemoveAll(outputDir + `\` + audioUniqueName)
		return err
	}

	// почистим после отправки аудиофайлы созданные локально
	os.Remove(inputDir + `\` + audioUniqueName)
	os.RemoveAll(outputDir + `\` + audioUniqueName)

	return nil
}
