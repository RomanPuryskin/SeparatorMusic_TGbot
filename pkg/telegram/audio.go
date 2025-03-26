package telegram

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// файл с методами для обработки аудио

func generateAudioFileName() string {
	currentTime := time.Now().UnixNano()
	return "input" + strconv.Itoa(int(currentTime))
}

func downloadAudio(audioURL string, audioInputPath string) error {

	// получаем файл с сервера
	response, err := http.Get(audioURL)
	if err != nil {
		return errors.New("ошибка получения файла с сервера")
	}
	defer response.Body.Close()

	// создадим файл
	dir, err := os.Create(audioInputPath)
	if err != nil {
		return errors.New("ошибка сохранения файла локально")
	}
	defer dir.Close()

	io.Copy(dir, response.Body)

	return nil
}

func runSpleeter(output, input, filename, formatAudio string) error {
	cmd := exec.Command("docker", "run", "--rm", "-v", input+":/app/input", "-v", output+":/app/output", "spleeter-image:latest", "separate", "-i", "/app/input/"+filename, "-c", formatAudio, "-o", "/app/output")
	err := cmd.Run()
	if err != nil {
		return errors.New("ошибка запуска контейнера Spleeter")
	}

	return nil
}

func sendAudioFiles(inputAudioName, outputDir string, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	// прочитаем файлы которые сделал Spleeter
	files, err := os.ReadDir(outputDir + `\` + inputAudioName)
	if err != nil {
		return errors.New("ошибка работы Spleeter или прочтения файлов")
	}

	// отправим файлы пользователю
	for _, file := range files {
		filePath := filepath.Join(outputDir+`\`+inputAudioName, file.Name())
		audio := tgbotapi.NewAudio(msg.Chat.ID, tgbotapi.FilePath(filePath))
		_, err = bot.Send(audio)
		if err != nil {
			return errors.New("ошибка отправки файла")
		}
	}

	return nil
}
