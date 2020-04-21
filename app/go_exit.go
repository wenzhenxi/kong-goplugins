/*
A "hello world" plugin in Go,
which reads a request header and sets a response header.
*/
package main

import (
	"github.com/sunmi-OS/go-pdk"
)

var i = 0

type Config struct {
	Body   string
	Status int
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {

	host, err := kong.Request.GetHeader("host")
	if err != nil {
		kong.Log.Err(err.Error())
	}
	message := conf.Body
	if message == "" {
		message = "hello"
	}
	status := conf.Status
	if status == 0 {
		status = 200
	}

	kong.Response.Exit(status, host+":"+message, nil)
}
