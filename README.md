# Wercker Universal Notifier

[![wercker status](https://app.wercker.com/status/437bb822b3448571623ea7384bac8d57/m "wercker status")](https://app.wercker.com/project/bykey/437bb822b3448571623ea7384bac8d57)

Forwards Wercker results to a remote service.
Does no internal handling of anything. Just takes the results and sends them to a specified host.

## Options

### required

- `host` - Your Service URL.
- `port` - Your Service Port (gRPC has 50051 as default).

## Example with step variables

```yml
build:
	after-steps:
		- romainmenke/universal-notifier:
			host: "52.51.22.243"
			port: "50051"
```

To create a Service implement a gRPC server based on Universal Notifier

[universal notifier source](https://github.com/romainmenke/universal-notifier)

```go
package main

import (
  "errors"
  "log"
  "net"

  "golang.org/x/net/context"
  "google.golang.org/grpc"

  "github.com/romainmenke/universal-notifier/pkg/wercker"
)

const (
	port = ":50051"
)

func main() {

  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
    return
  }

  s := grpc.NewServer()
  wercker.RegisterNotificationServiceServer(s, &server{})
  s.Serve(lis)

}

type server struct{}

func (s *server) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {

  // do stuff with results from wercker
  return &wercker.WerckerResponse{Success: false}, errors.New("unimplemented")

}
```

Proto

```protoc
message Git {

	string domain = 100;
	string owner = 101;
	string repository = 102;
	string branch = 103;
	string commit = 104;

}

message Result {

	bool result = 100;
	string failed_step_name = 101;
	string failed_step_message = 102;

}

message Build {

	int64 started = 100;
	string url = 101;
	string user = 102;

	enum Action {

		BUILD = 0;
		DEPLOY = 1;

	}
	Action action = 103;

}

message WerckerMessage {

	string url = 100;

	Build build = 200;
	Result result = 201;
	Git git = 202;

}
```
