# Normify Go SDK

Go клиент для [Normify API](https://normify.ru) — сервиса нормализации и обработки данных.

На данный момент реализован **синхронный метод** `/api/v1/auth/process` для мгновенной нормализации небольших объёмов данных (до 10 записей). Асинхронные методы будут добавлены позже.

## Установка

```bash
go get github.com/normify-ru/normify-go
```

## Получение API ключа

Для использования SDK необходим API ключ. Получить его можно:

1. Зарегистрируйтесь на [normify.ru](https://normify.ru)
2. Перейдите в личный кабинет
3. В разделе "API ключи" создайте новый ключ
4. Скопируйте полученный ключ и используйте его в коде

## Использование

### Основной пример

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/normify-ru/normify-go"
)

func main() {
	// Инициализация клиента с вашим API ключом
	client := normify.NewClient("ВАШ_API_КЛЮЧ_ЗДЕСЬ")

	// Подготовка запроса на нормализацию
	req := &normify.ProcessRequest{
		Entity: "company_name", // Тип сущности для нормализации
		Data: []normify.DataItem{
			{ID: "1", Value: "ООО Рога и Копыта"},
			{ID: "2", Value: "ИП Иванов Иван Иванович"},
		},
	}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Вызов API
	resp, err := client.Process(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	// Обработка результата
	fmt.Printf("Успешно: %v\n", resp.Success)
	fmt.Printf("Сущность: %s\n", resp.Data.Entity)
	
	for _, out := range resp.Data.Result.Output {
		fmt.Printf("ID: %s, Значение: %v\n", out.ID, out.Value)
		if len(out.Metadata) > 0 {
			fmt.Printf("  Метаданные: %v\n", out.Metadata)
		}
	}
	
	if len(resp.Data.Result.Errors) > 0 {
		fmt.Printf("Ошибки: %v\n", resp.Data.Result.Errors)
	}
}
```

### Расширенная конфигурация клиента

```go
package main

import (
	"net/http"
	"time"

	"github.com/normify-ru/normify-go"
)

func main() {
	// Клиент с кастомным HTTP клиентом и базовым URL
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	
	client := normify.NewClient(
		"ВАШ_API_КЛЮЧ_ЗДЕСЬ",
		normify.WithHTTPClient(httpClient),
		normify.WithBaseURL("https://api.normify.ru"), // Опционально
	)
	
	// ... использование клиента
}
```

## Структуры данных

### ProcessRequest
```go
type ProcessRequest struct {
	Entity string     `json:"entity"` // Ти
п сущности (например: "company_name", "person_name")
	Data   []DataItem `json:"data"`   // Данные для нормализации
}

type DataItem struct {
	ID    string      `json:"id"`    // Уникальный идентификатор записи
	Value interface{} `json:"value"` // Значение для нормализации
}
```

### ProcessResponse
```go
type ProcessResponse struct {
	Success bool        `json:"success"` // Успешность операции
	Data    ProcessData `json:"data"`    // Результаты нормализации
}

type ProcessData struct {
	Entity string        `json:"entity"` // Тип обработанной сущности
	Result ProcessResult `json:"result"` // Результаты и ошибки
}

type ProcessResult struct {
	Output []ProcessOutputItem `json:"output"` // Нормализованные данные
	Errors []interface{}       `json:"errors"` // Ошибки обработки
}

type ProcessOutputItem struct {
	ID       string                 `json:"id"`       // Идентификатор исходной записи
	Value    interface{}            `json:"value"`    // Нормализованное значение
	Metadata map[string]interface{} `json:"metadata"` // Дополнительная информация
}
```

## Обработка ошибок

```go
resp, err := client.Process(ctx, req)
if err != nil {
	// Проверка типа ошибки
	if apiErr, ok := err.(*normify.APIError); ok {
		fmt.Printf("API ошибка: статус %d, тело: %s\n", 
			apiErr.StatusCode, string(apiErr.Body))
	} else {
		fmt.Printf("Другая ошибка: %v\n", err)
	}
	return
}
```

## Ограничения

- **Синхронный метод**: поддерживает до 10 записей за один запрос
- **Таймаут по умолчанию**: 30 секунд (можно изменить через `WithHTTPClient`)
- **Аутентификация**: требуется валидный API ключ


## Дополнительные ресурсы

- [Официальная документация Normify](https://normify.ru/app/docs/api)
- [Примеры использования](https://github.com/normify-ru/normify-go/tree/main/examples)
