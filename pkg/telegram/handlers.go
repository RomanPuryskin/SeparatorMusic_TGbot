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
		setNewState(msg.Chat.ID, choosingFormatAudioForChangeState)
		handlerChoosingFormatAudioForChangeState(bot, msg)
	case replyButtonSeparateAudio:
		setNewState(msg.Chat.ID, choosingFormatAudioForSeparateState)
		handlerChoosingFormatAudioForSeparateState(bot, msg)
	default:
		handlerCommandStart(bot, msg)
	}
}

func handlerChoosingFormatAudioForSeparateState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	msgText1 := `Выберите формат, в котором хотите получить аудио`
	requestMsg1 := tgbotapi.NewMessage(msg.Chat.ID, msgText1)
	requestMsg1.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(requestMsg1)

	msgText2 := "Поддерживаемые форматы:"
	requestMsg2 := tgbotapi.NewMessage(msg.Chat.ID, msgText2)
	requestMsg2.ReplyMarkup = inlineButtonsForChooseFormateToSeparateKeyBoard()
	bot.Send(requestMsg2)
}

func handlerWaitingAudioForSeparateState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	switch {

	case msg.Text == replyButtonMainMenu:
		handlerCommandStart(bot, msg)
		return

	case msg.Text == replyButtonReturn:
		setNewState(msg.Chat.ID, choosingFormatAudioForSeparateState)
		handlerChoosingFormatAudioForSeparateState(bot, msg)
		return

	case msg.Audio != nil:
		err := handlerSeparateAudio(bot, msg)
		if err != nil {
			//TODO: заблокировать действия пользователя пока обрабатывается аудио
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Произошла ошибка на сервере, попробуйте еще раз позднее"))
		}

	default:
		msgText := `Загрузите аудиофайл для обработки.
		Для возврата к выбору формата нажмите Return
		Для возврата в главное меню нажмите Main menu`

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

func handlerChoosingFormatAudioForChangeState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	msgText1 := `Выберите формат, в котором хотите получить аудио`
	requestMsg1 := tgbotapi.NewMessage(msg.Chat.ID, msgText1)
	requestMsg1.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(requestMsg1)

	msgText2 := "Поддерживаемые форматы:"
	requestMsg2 := tgbotapi.NewMessage(msg.Chat.ID, msgText2)
	requestMsg2.ReplyMarkup = inlineButtonsForChooseFormateToChangeKeyBoard()
	bot.Send(requestMsg2)
}

func handlerWaitingAudioForChangeFormatState(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	switch {
	case msg.Text == replyButtonMainMenu:
		handlerCommandStart(bot, msg)
		return

	case msg.Text == replyButtonReturn:
		setNewState(msg.Chat.ID, choosingFormatAudioForChangeState)
		handlerChoosingFormatAudioForChangeState(bot, msg)
		return

	case msg.Audio != nil:
		// обработка
		err := handlerChangeFormatAudio(bot, msg)
		if err != nil {
			//TODO: заблокировать действия пользователя пока обрабатывается аудио
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Произошла ошибка на сервере, попробуйте еще раз позднее"))
		}

	// любое другое сообщение
	default:
		msgText := `Загрузите аудиофайл для обработки.
		Для возврата к выбору формата нажмите Return
		Для возврата в главное меню нажмите Main menu`

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

func handlerSeparateAudio(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {

	// первое сообщение при начале обработки
	msgText := "Идет обработка..."
	requestMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	requestMsg.ReplyToMessageID = msg.MessageID

	// отправим сообщение
	firstMessage, _ := bot.Send(requestMsg)

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
	switch getCurrentFormatForSeparate(msg.Chat.ID) {
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
		// если не удалось запустить контейнер, удалим уже скачанный файл
		os.Remove(inputDir + `\` + audioUniqueName)
		return err
	}

	//второе сообщение, о том, что файл обработан и теперь отправляется
	editedText := "Файл обработан. Идет отправка..."
	editedMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, firstMessage.MessageID, editedText)
	bot.Send(editedMsg)

	err = sendAudioFiles(outputDir+`\`+audioUniqueName, bot, msg)
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

func handlerChangeFormatAudio(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	// первое сообщение при начале обработки
	msgText := "Идет обработка..."
	requestMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	requestMsg.ReplyToMessageID = msg.MessageID

	// отправим сообщение
	firstMessage, _ := bot.Send(requestMsg)

	//получим URL аудиофайла
	audioURL, err := bot.GetFileDirectURL(msg.Audio.FileID)
	if err != nil {
		return errors.New("ошибка при получении файла с сервера")
	}

	inputDir := `.\pkg\ffmpeg\input`   // Папка для входных файлов
	outputDir := `.\pkg\ffmpeg\output` // папка для выходных файлов

	audioUniqueName := generateAudioFileName()

	// скачиваем файл
	err = downloadAudio(audioURL, inputDir+`\`+audioUniqueName)
	if err != nil {
		return err
	}

	// получим явно формат аудио, который желает пользователь
	format := ""
	switch getCurrentFormatForChange(msg.Chat.ID) {
	case mp3ChangeFormat:
		format = "mp3"
	case flacChangeFormat:
		format = "flac"
	case aacChangeFormat:
		format = "aac"
	case oggChangeFormat:
		format = "ogg"
	case m4aChangeFormat:
		format = "m4a"
	default:
		format = "wav"
	}

	//сгенерируем имя аудиофайла, в который сохраним результат
	audioOutputName := msg.Audio.FileName + "." + format
	// Запускаем ffmpeg
	err = runFFmpeg(outputDir, inputDir, audioUniqueName, audioOutputName)
	if err != nil {
		// если не удалось запустить контейнер, удалим уже скачанный файл
		os.Remove(inputDir + `\` + audioUniqueName)
		return err
	}

	//второе сообщение, о том, что файл обработан и теперь отправляется
	editedText := "Файл обработан. Идет отправка..."
	editedMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, firstMessage.MessageID, editedText)
	bot.Send(editedMsg)

	// TODO:
	err = sendAudioFiles(outputDir, bot, msg)
	if err != nil {
		// если произойдет ошибка отправки по какой то причине, так же почистим файлы
		os.Remove(inputDir + `\` + audioUniqueName)
		os.Remove(outputDir + `\` + audioOutputName)
		return err
	}

	// почистим после отправки аудиофайлы созданные локально
	os.Remove(inputDir + `\` + audioUniqueName)
	os.Remove(outputDir + `\` + audioOutputName)

	return nil
}
