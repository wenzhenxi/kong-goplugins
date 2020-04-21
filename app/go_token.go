/*
A "hello world" plugin in Go,
which reads a request header and sets a response header.
*/
package main

import (
	"github.com/sunmi-OS/go-pdk"
	"kong-goplugins/pkg/db"
	"strings"
)

var i = 0

type Config struct {
	DBConnect string
	DBType    string
	DBMaxIdle int
	DBMaxOpen int
	TokenName string
	SetIDName string
	SecretSQL string
	IsDebug   bool
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {

	// 获取参数进行判断是否具备所需条件
	token, err := kong.Request.GetHeader(strings.ToLower(conf.TokenName))
	if err != nil {
		if conf.IsDebug {
			kong.Response.Exit(200, `{"code":20000,"data":"","msg":"missing Parameters Token"}`, nil)
		} else {
			kong.Response.Exit(200, `{"code":20000,"data":"","msg":"missing Parameters"}`, nil)
		}
		return
	}

	// 获取数据库连接
	orm, err := db.GetOrm(conf.DBConnect, conf.DBType, conf.DBMaxIdle, conf.DBMaxOpen)
	if err != nil {
		if conf.IsDebug {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error : `+err.Error()+`"}`, nil)
		} else {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error"}`, nil)
		}
		return
	}
	// 字段存在去数据库获取 包名和secret
	id, err := db.GetSecretByAppId(orm, conf.SecretSQL, token)
	if err != nil {
		if err == db.ErrRecordNotFound {
			kong.Response.Exit(200, `{"code":30001,"data":"","msg":"id Nil Error"}`, nil)
			return
		}
		if conf.IsDebug {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error : `+err.Error()+`"}`, nil)
		} else {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error"}`, nil)
		}
		return
	}
	// 未获取到直接报错返回
	if id == "" {
		kong.Response.Exit(200, `{"code":30001,"data":"","msg":"id Nil Error"}`, nil)
		return
	}

	err = kong.ServiceRequest.SetHeader(conf.SetIDName, id)
	if err != nil {
		if conf.IsDebug {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error : `+err.Error()+`"}`, nil)
		} else {
			kong.Response.Exit(200, `{"code":50001,"data":"","msg":"service Error"}`, nil)
		}
		return
	}
}
