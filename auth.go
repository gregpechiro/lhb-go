package main

import "github.com/cagnosolutions/web"

var ADMIN = web.Auth{
	Roles:    []string{"admin"},
	Redirect: "/",
	Msg:      "You are not authorized",
}
