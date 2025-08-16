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

func newClient(server, username, password string, insecure bool) (*caldav.Client, error) {
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
	inner, err := caldav.NewClient(webdavHttp, server)
	return inner, err
}

func getClientFromEnv(ctx context.Context, call *nu.ExecCommand) (client *caldav.Client, err error) {
	urlVar, err := call.GetEnvVar(ctx, "NU_PLUGIN_CALDAV_URL")
	if err != nil {
		return
	}
	usernameVar, err := call.GetEnvVar(ctx, "NU_PLUGIN_CALDAV_USERNAME")
	if err != nil {
		return
	}
	passwordVar, err := call.GetEnvVar(ctx, "NU_PLUGIN_CALDAV_PASSWORD")
	if err != nil {
		return
	}
	insecureVar, err := call.GetEnvVar(ctx, "NU_PLUGIN_CALDAV_INSECURE")
	if err != nil {
		return
	}

	if urlVar == nil {
		err = fmt.Errorf("NU_PLUGIN_CALDAV_URL is not set")
		return
	}
	url, err := tryCast[string](*urlVar)
	if err != nil {
		return
	}
	if url == "" {
		err = fmt.Errorf("NU_PLUGIN_CALDAV_URL is an empty string")
		return
	}

	var username string
	if usernameVar != nil {
		username, err = tryCast[string](*usernameVar)
		if err != nil {
			return
		}
	}
	var password string
	if passwordVar != nil {
		password, err = tryCast[string](*passwordVar)
		if err != nil {
			return
		}
	}
	var insecure string
	if insecureVar != nil {
		insecure, err = tryCast[string](*insecureVar)
		if err != nil {
			return
		}
	}

	client, err = newClient(url, username, password, insecure == "1")
	return
}

func tryCast[T any](val nu.Value) (T, error) {
	cast, ok := val.Value.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("expected type %T, got %T", zero, val.Value)
	}
	return cast, nil
}
