package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

type Env struct {
	URL                string
	Username, Password string
	Insecure           bool
}

func getEnvString(ctx context.Context, call *nu.ExecCommand, env string) (val string, err error) {
	variable, err := call.GetEnvVar(ctx, env)
	if err != nil {
		return
	}
	if variable == nil {
		err = fmt.Errorf("%s not set", env)
		return
	}
	str, err := tryCast[string](*variable)
	if err != nil {
		return
	}
	val = str
	return
}

func getEnvBool(ctx context.Context, call *nu.ExecCommand, env string) (val bool, err error) {
	variable, err := call.GetEnvVar(ctx, env)
	if err != nil {
		return
	}
	if variable == nil {
		err = fmt.Errorf("%s not set", env)
		return
	}
	str, err := tryCast[string](*variable)
	if err != nil {
		return
	}
	val = str == "1" || str == "true"
	return
}

func getClient(ctx context.Context, call *nu.ExecCommand) (client *caldav.Client, err error) {
	url, err := getEnvString(ctx, call, "NU_PLUGIN_CALDAV_URL")
	if err != nil {
		return
	}
	username, err := getEnvString(ctx, call, "NU_PLUGIN_CALDAV_USERNAME")
	if err != nil {
		return
	}
	password, err := getEnvString(ctx, call, "NU_PLUGIN_CALDAV_PASSWORD")
	if err != nil {
		return
	}
	insecure, err := getEnvBool(ctx, call, "NU_PLUGIN_CALDAV_INSECURE")
	if err != nil {
		return
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	webdavHttp := webdav.HTTPClient(httpClient)
	if username != "" && password != "" {
		webdavHttp = webdav.HTTPClientWithBasicAuth(httpClient, username, password)
	}
	client, err = caldav.NewClient(webdavHttp, url)

	return
}

type job interface {
	Do(ctx context.Context) error
}

func worker(ctx context.Context, jobs chan job, status chan error) {
	for j := range jobs {
		if j == nil {
			continue
		}
		status <- j.Do(ctx)
	}
}

func parallelizeJobs(ctx context.Context, jobs []job, parallel int) (err error) {
	if len(jobs) == 0 {
		return
	}

	jobsChan := make(chan job)
	status := make(chan error)
	defer close(jobsChan)
	defer close(status)

	for range parallel {
		go worker(ctx, jobsChan, status)
	}

	current := jobs[0]
	currentIdx := 0
	finished := 0
	for {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("context canceled")
			return
		case err = <-status:
			if err != nil {
				return
			}
			finished++
			if finished >= len(jobs) {
				return
			}
		case jobsChan <- current:
			if currentIdx >= len(jobs) {
				continue
			}
			currentIdx++
			if currentIdx >= len(jobs) {
				current = nil
				continue
			}
			current = jobs[currentIdx]
		}
	}
}

func caldavKeywordsQuery(additional ...string) []string {
	return append([]string{"caldav", "query", "search", "find", "pull", "filter"}, additional...)
}
