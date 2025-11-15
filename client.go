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

