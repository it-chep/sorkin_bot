package middleware

type MessageLogMiddleware struct {
}

func NewMessageLogMiddleware() MessageLogMiddleware {
	return MessageLogMiddleware{}
}

type SentryMiddleware struct {
}

func NewSentryMiddleware() SentryMiddleware {
	return SentryMiddleware{}
}

// TgAdminWarningMiddleware todo подумать надо ли мне это
type TgAdminWarningMiddleware struct {
}

func NewTgAdminWarningMiddleware() TgAdminWarningMiddleware {
	return TgAdminWarningMiddleware{}
}
