package telegram

// временное решение для хранения текущего состояния пользователя в диалоге
// userState хранит пару ключ - значение:
// ключ - chatID - индентификатор чата
// значение - строка текущения состояния пользователя в диалоге

var userState = make(map[int64]string)

// функция для установки нового состояния
func setNewState(chatID int64, state string) {
	userState[chatID] = state
}

// функция для получения текущего состояния

func getCurrentState(chatID int64) string {
	return userState[chatID]
}

// временное решение для хранения текущего желаемого формата аудио для Separate
// userState хранит пару ключ - значение:
// ключ - chatID - индентификатор чата
// значение - строка текущего желаемого формата
var userCurrentFormatForSeparate = make(map[int64]string)

// функция для установки текущего желаемого формата для Separate
func setNewCurrentFormat(chatID int64, state string) {
	userCurrentFormatForSeparate[chatID] = state
}

// функция для получения текущего формата желаемого аудио для Separate
func getCurrentFormat(chatID int64) string {
	return userCurrentFormatForSeparate[chatID]
}
