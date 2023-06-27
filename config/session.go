package config

import "github.com/gorilla/sessions"

const SESSION_ID = "go_auth_session"

var Store = sessions.NewCookieStore([]byte("asldkjalsihtasdnsdbasd"))
