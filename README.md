# Image Archiver Service

HTTP-сервис для создания задач архивации файлов `.pdf`, `.jpeg`, `.jpg` (в том числе с параметрами в URL) и скачивания архива.

---

## API

### 1. Создать задачу

**POST** `/task/create`
curl -X POST http://localhost:8080/task/create

ответ:
{
  "id": "<task_id>",
  "status": "pending",
  "files": [],
  "failed_files": [],
  "url": ""
}

### 2. Добавить файлы в задачу
POST /task/add


curl -X POST http://localhost:8080/task/add \
-H "Content-Type: application/json" \
-d '{
  "id": "<task_id>",
  "files": [
    "https://example.com/image.jpeg?size=large",
    "https://example.com/doc.pdf"
  ]
}'

### 3. Получить статус задачи
GET /task/status?id=<task_id>
curl "http://localhost:8080/task/status?id=<task_id>"

ответ:
{
  "id": "<task_id>",
  "status": "complete",
  "files": [...],
  "failed_files": [],
  "url": "http://localhost:8080/download/<task_id>.zip"
}
4. Скачать архив
GET /download/<task_id>.zip
curl -O http://localhost:8080/download/<task_id>.zip

поддерживает - .pdf, .jpeg, .jpg
проверка расширения игнорирует параметры URL

Архивы доступны по /download/

Запуск go run cmd/main.go
```



