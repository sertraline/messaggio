# Запуск проекта
1. Склонировать репозиторий:  
`git clone https://github.com/sertraline/messaggio`
2. Для persisted storage нужно создать директории которые будут спроксированы в докер, иначе данные будут утеряны после остановки контейнера.
  * `mkdir -p volumes/psql-write`
  * `mkdir -p volumes/pgadmin`
3. Контейнеру нужны права на запись/чтение данных в директорию, postgres использует UID 1001, pgadmin использует 5050.  
  * `chown -R 1001:1001 volumes/`
  * `chown -R 5050:5050 volumes/pgadmin`
4. `cp .env.example .env` -> заполнить переменные в .env или оставить как есть
5. Можно запустить контейнеры через `docker compose up -d`.
6. Теперь нужно создать топик в kafka.  
  * Зайдем в шелл контейнера: `docker exec -it messaggio.kafka /bin/bash`  
  * Создадим топик:  
  `kafka-topics --create --topic messages --zookeeper messaggio.zookeeper:2181 --partitions 2 --replication-factor 1`
6. Запустить проект: `go run .`


# Структура
Проект придерживается структуры routes -> controllers -> services.  

* main - определяет роуты и эндпоинты  
* controllers - отвечает за логику эндпоинтов  
* services - содержит бизнес-логику эндпоинтов  
* validators - валидация входящих запросов и данных в POST формах  
* errors - ошибки используемые в проекте  
* database - содержит конфигурацию sqlx и kafka-go  
  * database/models - сериализация/десериализация данных из постгреса


# API
Для запросов к API я использую утилиту curl.

Модель Messages:
```
Id        	int       `json:"id"`
Content		string    `json:"content"`
CreatedAt 	time.Time `json:"created_at" db:"created_at"`
Processed 	bool `json:"processed"`
```

1. Создать сообщение  
```bash
curl -X "POST" http://localhost:3333/messages -d '{"content": "123"}'
```  

Сообщение будет сохранено в базе данных постгреса и отправлено в topic 'messages' кафки с флагом processed=true. Если отправка в кафку не удалась, сообщение будет помечено processed=false.  

Для ключей в kafka я использую генератор UUID.

2. Получить сообщение из БД по ID:  
```bash
curl http://localhost:3333/messages/1
```  

Если сообщение не найдено, будет возвращен статус `{"status":"Resource not found."}`.  

3. Получить статистику сообщений:
```bash
curl http://localhost:3333/stats  
```  

Будет возвращен JSON объект вида `{"Всего сообщений":8,"Обработано сообщений":2}`.   

4. Прочитать сообщение из kafka:  
```bash
curl http://localhost:3333/kafka/read
```  

При наличии сообщения в очереди, будет возвращен объект вида:  
`{"key":"d09995bd-f77f-420e-b3ef-6f5b7eef9fac","offset":2,"partition":0 "topic":"messages","value":"123"}`.  

Если сообщения нет, метод будет блокировать выдачу ответа до таймаута или появления сообщения в очереди.
