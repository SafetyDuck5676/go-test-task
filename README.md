# Queue Server on Windows

Это краткое руководство по запуску проекта на Windows.

---

## 1. Установка Go

1. Скачайте и установите Go с официального сайта: [https://golang.org/dl/](https://golang.org/dl/).
2. Проверьте установку:
   ```cmd
   go version
   ```
   Вы должны увидеть что-то вроде:
   ```
   go version go1.20.5 windows/amd64
   ```

---

## 2. Клонирование проекта

1. Создайте папку для проекта, например `C:\queue-server`.
2. Скопируйте все файлы проекта (`main.go`, `queue.go`, `server.go`) в эту папку.

---

## 3. Компиляция программы

Откройте командную строку (`cmd`) или PowerShell и перейдите в папку проекта:
```cmd
cd C:\queue-server
```

Скомпилируйте программу:
```cmd
go build -o queue-server.exe
```

---

## 4. Запуск сервера

Запустите скомпилированный файл с параметрами:
```cmd
queue-server.exe -port=8080 -max-queues=10 -default-size=100 -timeout=5
```

- `-port`: порт, на котором будет работать сервер (например, `8080`).
- `-max-queues`: максимальное количество очередей (например, `10`).
- `-default-size`: размер очереди по умолчанию (например, `100`).
- `-timeout`: таймаут ожидания сообщения в секундах (например, `5`).

Если все настроено правильно, вы увидите сообщение:
```
Starting server on port 8080
```

---

## 5. Тестирование

Для тестирования можно использовать `curl` или Postman.

### Добавить сообщение в очередь:
```cmd
curl -XPUT http://localhost:8080/queue/pet -d "{\"message\":\"hello\"}"
```

### Получить сообщение из очереди:
```cmd
curl http://localhost:8080/queue/pet?timeout=10
```

---

## 6. Альтернатива: Использование PowerShell

Если `curl` не работает, используйте PowerShell:

### Добавить сообщение в очередь:
```powershell
curl -Method PUT -Uri "http://localhost:8080/queue/pet" -Body '{"message":"hello"}' -ContentType "application/json"
```

### Получить сообщение из очереди:
```powershell
curl -Uri "http://localhost:8080/queue/pet?timeout=10"
```

---

## 7. Возможные проблемы

1. **Go не найден**:
   - Убедитесь, что Go добавлен в переменную среды `PATH`. Проверьте это командой:
     ```cmd
     echo %PATH%
     ```
   - Если Go отсутствует в `PATH`, добавьте путь к папке `bin` Go (например, `C:\Go\bin`) в системные переменные среды.

2. **Проблемы с curl**:
   - Если `curl` не работает, скачайте его с [https://curl.se/download.html](https://curl.se/download.html) или используйте PowerShell.

---