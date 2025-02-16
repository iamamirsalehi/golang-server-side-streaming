package notificationservice

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/server-side-streaming/notificationservice/notificationproto"
	"google.golang.org/grpc"
	"time"
)

var _ notificationproto.NotificationServiceServer = (*Handler)(nil)

type Handler struct {
	notificationproto.UnimplementedNotificationServiceServer
	redisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{
		redisClient: redisClient,
	}
}

func (h *Handler) GetNotifications(request *notificationproto.NotificationRequest, server grpc.ServerStreamingServer[notificationproto.Notification]) error {
	pubSub := h.redisClient.Subscribe(server.Context(), fmt.Sprintf("notification:%s", request.GetUserId()))
	for {
		select {
		case <-server.Context().Done():
			return server.Context().Err()
		case msg := <-pubSub.Channel():
			if err := server.Send(&notificationproto.Notification{
				UserId:    request.GetUserId(),
				Content:   fmt.Sprintf("new notification at %s: %s", time.Now().String(), msg.Payload),
				CreatedAt: time.Now().UnixMilli(),
			}); err != nil {
				return fmt.Errorf("could not send notification: %w", err)
			}

		}
	}
}
