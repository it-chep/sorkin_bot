package notification

import "context"

type notifyGateway interface {
	SendNotification(ctx context.Context, to []string, message string) error
}
