package main

import "github.com/cagnosolutions/web"

var WEBMASTER = web.Auth{
	Roles:    []string{"webmaster"},
	Redirect: "/",
	Msg:      "You are not authorized",
}
