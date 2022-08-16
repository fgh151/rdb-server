package main

import (
	err2 "db-server/err"
	"db-server/plugins"
	"github.com/RattusPetrucho/smsc"
)

type smscPlugin string

type SmsCPluginParams struct {
	Sender  string
	Message string
	Phones  string

	Login    string
	Password string
}

func (p smscPlugin) Run(params plugins.PluginParams) plugins.PluginResult {
	pp := params.Data.(SmsCPluginParams)

	sc, err := smsc.New(pp.Login, pp.Password)
	err2.PanicErr(err)

	err = sc.SetSenderName(pp.Sender)
	err2.PanicErr(err)

	resp, err := sc.SendSms("", pp.Message, pp.Phones)

	return plugins.PluginResult{
		Payload: resp,
		Err:     err,
	}
}

var Run smscPlugin
