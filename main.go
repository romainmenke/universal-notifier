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

	fmt.Println("Notify")

	ctx := context.Background()

	span, _ := trace.New(ctx, "New Send Wercker Message Job")
	defer span.Close()

	env, err := wercker.New(ctx)
	if err != nil {
		fmt.Println(err)
		span.Error(err)
		return
	}

	fmt.Printf("Getting ready to dial %s", env.Host())

	message, err := env.NewMessage(ctx)
	if err != nil {
		fmt.Println(err)
		span.Error(err)
		return
	}

	conn, err := grpc.Dial(env.Host(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := wercker.NewNotificationServiceClient(conn)
	_, err = c.Notify(ctx, message)
	if err != nil {
		span.Error(err)
		log.Fatal(err)
		return
	}

	fmt.Println("Successfully notified your service")

}
