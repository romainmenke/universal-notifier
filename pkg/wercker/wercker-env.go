package wercker

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/net/context"
	"limbo.services/trace"
)

type WerckerEnv struct {
	port string
	host string

	started        string
	ci             string
	buildURL       string
	domain         string
	owner          string
	repo           string
	branch         string
	commit         string
	root           string
	source         string
	cache          string
	output         string
	user           string
	applicationURL string
	result         string
	stepName       string
	stepMessage    string
	action         Build_Action
}

func New(ctx context.Context) (*WerckerEnv, error) {

	span, _ := trace.New(ctx, "client.grpc.newEnv")
	defer span.Close()

	env := WerckerEnv{}

	env.port = os.Getenv("WERCKER_UNIVERSAL_NOTIFIER_PORT")
	env.host = os.Getenv("WERCKER_UNIVERSAL_NOTIFIER_HOST")

	if env.port == "" || env.host == "" {
		return nil, span.Error("Set REMOTE_HOST and REMOTE_PORT")
	}

	env.started = os.Getenv("WERCKER_MAIN_PIPELINE_STARTED")

	env.ci = os.Getenv("CI")
	env.buildURL = os.Getenv("WERCKER_RUN_URL")

	env.domain = os.Getenv("WERCKER_GIT_DOMAIN")
	env.owner = os.Getenv("WERCKER_GIT_OWNER")
	env.repo = os.Getenv("WERCKER_GIT_REPOSITORY")
	env.branch = os.Getenv("WERCKER_GIT_BRANCH")
	env.commit = os.Getenv("WERCKER_GIT_COMMIT")

	env.root = os.Getenv("WERCKER_ROOT")
	env.source = os.Getenv("WERCKER_SOURCE_DIR")
	env.cache = os.Getenv("WERCKER_CACHE_DIR")
	env.output = os.Getenv("WERCKER_OUTPUT_DIR")

	env.user = os.Getenv("WERCKER_STARTED_BY")
	env.applicationURL = os.Getenv("WERCKER_APPLICATION_URL")

	env.result = os.Getenv("WERCKER_RESULT")
	env.stepName = os.Getenv("WERCKER_FAILED_STEP_DISPLAY_NAME")
	env.stepMessage = os.Getenv("WERCKER_FAILED_STEP_MESSAGE")

	action := os.Getenv("DEPLOY")

	fmt.Println(action)
	if action == "" {
		env.action = Build_DEPLOY
	} else {
		env.action = Build_BUILD
	}

	return &env, nil
}

func (e *WerckerEnv) NewGit() *Git {
	return &Git{
		Domain:     e.domain,
		Owner:      e.owner,
		Repository: e.repo,
		Branch:     e.branch,
		Commit:     e.commit,
	}
}

func (e *WerckerEnv) NewResult() *Result {
	var result bool
	if e.result == "failed" {
		result = false
	} else {
		result = true
	}

	return &Result{
		Result:            result,
		FailedStepName:    e.stepName,
		FailedStepMessage: e.stepMessage,
	}
}

func (e *WerckerEnv) NewBuild(ctx context.Context) (*Build, error) {

	span, _ := trace.New(ctx, "client.grpc.newBuild")
	defer span.Close()

	started, err := strconv.ParseInt(e.started, 10, 64)
	if err != nil {
		return nil, span.Error(err)
	}

	return &Build{
		Started: started,
		Url:     e.buildURL,
		User:    e.user,
		Action:  e.action,
	}, nil
}

func (e *WerckerEnv) NewMessage(ctx context.Context) (*WerckerMessage, error) {

	span, _ := trace.New(ctx, "client.grpc.newMessage")
	defer span.Close()

	build, err := e.NewBuild(ctx)
	if err != nil {
		return nil, span.Error(err)
	}

	result := e.NewResult()
	git := e.NewGit()

	return &WerckerMessage{
		Url:    e.applicationURL,
		Build:  build,
		Result: result,
		Git:    git,
	}, nil
}

func (e *WerckerEnv) Host() string {
	return fmt.Sprintf("%s:%s", e.host, e.port)
}

func TestEnv() *WerckerEnv {
	return &WerckerEnv{
		port:           "50051",
		host:           "52.51.22.243",
		started:        "12345",
		ci:             "",
		buildURL:       "localhost",
		domain:         "github.com",
		owner:          "romainmenke",
		repo:           "test",
		branch:         "master",
		commit:         "update",
		root:           "",
		source:         "",
		cache:          "",
		output:         "",
		user:           "romainmenke",
		applicationURL: "",
		result:         "false",
		stepName:       "blah",
		stepMessage:    "stuf",
		action:         Build_BUILD,
	}
}
