package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"sorkin_bot/internal/controller/dto/tg"
)

type MessageLogMiddleware struct {
	messageService MessageService
}

func NewMessageLogMiddleware() MessageLogMiddleware {
	return MessageLogMiddleware{}
}

func (middleware MessageLogMiddleware) ProcessRequest(c *gin.Context) {
	var update tgbotapi.Update

	if err := c.ShouldBindJSON(&update); err != nil {
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

	tgMessage := tg.MessageDTO{
		MessageID: int64(sourceMessage.MessageID),
		Text:      sourceMessage.Text,
	}

	go func() {
		if err := middleware.messageService.SaveMessageLog(context.Background(), tgMessage); err != nil {
			log.Println("Failed to log incoming message:", err)
		}
	}()

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
