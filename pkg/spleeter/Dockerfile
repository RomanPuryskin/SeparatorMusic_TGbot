# Используем официальный образ Spleeter от Deezer
FROM researchdeezer/spleeter:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Создаем папки для входных и выходных данных
RUN mkdir -p /app/input 
RUN mkdir -p /app/output

# Создайте директорию для моделей
RUN mkdir -p /model/2stems

# Скопируйте скачанную модель в образ
COPY 2stems.tar.gz /model/2stems/

# Распакуйте модель
RUN tar -xzf /model/2stems/2stems.tar.gz -C /model/2stems/ && \
    rm /model/2stems/2stems.tar.gz

# Укажите путь к модели в переменной окружения
ENV MODEL_PATH /model