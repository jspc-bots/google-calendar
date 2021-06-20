package main

import "os"

const (
	Chan = "#dashboard"
)

var (
	Username        = os.Getenv("SASL_USER")
	Password        = os.Getenv("SASL_PASSWORD")
	Server          = os.Getenv("SERVER")
	VerifyTLS       = os.Getenv("VERIFY_TLS") == "true"
	CredentialsFile = os.Getenv("CREDENTIALS_FILE")
	TokenFile       = os.Getenv("TOKEN_FILE")
	Timezone        = os.Getenv("TZ")
)

func main() {
	g, err := NewGoogle(CredentialsFile, TokenFile, Timezone)
	if err != nil {
		panic(err)
	}

	c, err := New(Username, Password, Server, VerifyTLS, g)
	if err != nil {
		panic(err)
	}

	c.bottom.Client.Connect()
}
