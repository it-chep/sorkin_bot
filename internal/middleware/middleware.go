package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"sorkin_bot/internal/controller/dto/tg"
)

type MessageLogMiddleware struct {
	messageService MessageService
}

func NewMessageLogMiddleware(messageService MessageService) MessageLogMiddleware {
	return MessageLogMiddleware{
		messageService: messageService,
	}
}

func (middleware MessageLogMiddleware) ProcessRequest(c *gin.Context) {
	var update tgbotapi.Update
	// todo выглядит как костыль, потому что решил логать сообщение в мидлварине, исправить
	// todo костыль потому что если прочитать контекст, то он потом будет пустым
	// Чтение тела запроса
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	// Восстановление тела запроса обратно в тело запроса (поскольку Read() его потребовал)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Декодирование JSON в объект update
	if err = json.Unmarshal(body, &update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sourceMessage *tgbotapi.Message
	if update.CallbackQuery != nil {
		sourceMessage = update.CallbackQuery.Message
	} else if update.Message != nil {
		sourceMessage = update.Message
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No message or callback query data found"})
		return
	}
	chatId := tg.Chat{ID: sourceMessage.Chat.ID}

	tgMessage := tg.MessageDTO{
		MessageID: int64(sourceMessage.MessageID),
		Text:      sourceMessage.Text,
		Chat:      &chatId,
	}

	go func() {
		if err := middleware.messageService.SaveMessageLog(context.Background(), tgMessage); err != nil {
			log.Println("Failed to log incoming message:", err)
		}
	}()
	//c.Set("update", byte(update))
	c.Next()
}

// todo

type SentryMiddleware struct {
}

func NewSentryMiddleware() SentryMiddleware {
	return SentryMiddleware{}
}
func (middleware SentryMiddleware) ProcessRequest(c *gin.Context) {
	c.Next()
}

// TgAdminWarningMiddleware todo подумать надо ли мне это
type TgAdminWarningMiddleware struct {
}

func NewTgAdminWarningMiddleware() TgAdminWarningMiddleware {
	return TgAdminWarningMiddleware{}
}
func (middleware TgAdminWarningMiddleware) ProcessRequest(c *gin.Context) {
	c.Next()
}
