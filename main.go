package main

import (
	"fmt"
	"log"

	"limbo.services/trace"

	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	span, _ := trace.New(ctx, "New Send Wercker Message Job")
	defer span.Close()

	env, err := wercker.New(ctx)

	fmt.Println("Notify")

	if err != nil {
		fmt.Println(err)
		span.Error(err)
		return
	}

	message, err := env.NewMessage(ctx)

	conn, err := grpc.Dial(env.Host(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := wercker.NewNotificationServiceClient(conn)
	_, err = c.Notify(ctx, message)
	if err != nil {
		fmt.Println(err)
		span.Error(err)
	}

}
