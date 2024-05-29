package config

import "github.com/gorilla/sessions"

// install dulu session go get github.com/gorilla/sessions

const SESSION_ID = "go_auth_sess"

var Store = sessions.NewCookieStore([]byte("bebasapasaja"))
