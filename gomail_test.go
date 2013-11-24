package main

import (
	_ "fmt"
	"net/smtp"
	"testing"
)

func TestWelcomeMessage(t *testing.T) {
	startServer(3005)

	client, err := smtp.Dial("localhost:3005")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Hello("cheeri.os")
	if err != nil {
		t.Fatal(err)
	}

	stopServer()
}
