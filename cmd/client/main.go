package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/server-side-streaming/notificationservice"
	"github.com/server-side-streaming/notificationservice/notificationproto"
	"io"
	"log"
)

func main() {
	client, err := notificationservice.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	stream, err := client.GetNotifications(ctx, &notificationproto.NotificationRequest{
		UserId: "123",
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		notification, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed to read notification: %v", err)
		}

		b, err := json.MarshalIndent(notification, "", "\t")
		if err != nil {
			log.Fatalf("failed to marshal notification: %v", err)
		}

		fmt.Println(string(b))
	}
}
