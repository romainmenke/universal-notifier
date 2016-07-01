package main

import (
	"fmt"

	"limbo.services/trace"
	"limbo.services/trace/dev"

	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	trace.DefaultHandler = dev.NewHandler(nil)

	fmt.Println("Notify")

	ctx := context.Background()

	span, ctx := trace.New(ctx, "client.grpc.notify")
	defer span.Close()

	env, err := wercker.New(ctx)
	if err != nil {
		span.Error(err)
		return
	}

	fmt.Printf("Getting ready to dial %s", env.Host())

	message, err := env.NewMessage(ctx)
	if err != nil {
		span.Error(err)
		return
	}

	if message.Git.Branch != "master" {
		span.Error("Not on the master branch")
		return
	}

	conn, err := grpc.Dial(env.Host(), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := wercker.NewNotificationServiceClient(conn)
	_, err = c.Notify(ctx, message)
	if err != nil {
		span.Error(err)
		return
	}

	fmt.Println("Successfully notified your service")

}
