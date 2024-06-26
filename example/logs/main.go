package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/innoobijr/faas-provider/types"
	"github.com/innoobijr/go-sdk"
)

func main() {
	// NOTE: You can have any name for environment variables. below defined variables names are not standard names
	username := os.Getenv("OPENFAAS_USERNAME")
	password := os.Getenv("OPENFAAS_PASSWORD")

	gatewayURL, _ := url.Parse(os.Getenv("OPENFAAS_GATEWAY_URL"))
	auth := &sdk.BasicAuth{
		Username: username,
		Password: password,
	}

	client := sdk.NewClient(gatewayURL, auth, http.DefaultClient)

	// Deploy function
	status, err := client.Deploy(context.Background(), types.FunctionDeployment{
		Service:    "env-store-test",
		Image:      "ghcr.io/openfaas/alpine:latest",
		Namespace:  "openfaas-fn",
		EnvProcess: "env",
	})
	// non 200 status value will have some error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Status: %d Deploy Failed: %s", status, err)
		os.Exit(1)
	}

	// Follow is allows the user to request a stream of logs until the timeout
	follow := false
	// Tail sets the maximum number of log messages to return, <=0 means unlimited
	tail := 5
	// Since is the optional datetime value to start the logs from
	since := time.Now().Add(-30 * time.Second)

	logsChan, err := client.GetLogs(context.Background(), "env-store-test", "openfaas-fn", follow, tail, &since)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get Logs Failed: %s", err)
		os.Exit(1)
	}

	fmt.Println("Logs Received....")
	for line := range logsChan {
		fmt.Println(line)
	}
}
