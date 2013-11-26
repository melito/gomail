package main

import (
	"net/smtp"
	"testing"
)

func TestHelloQuit(t *testing.T) {

	go startServer(3005)

	client, err := smtp.Dial("localhost:3005")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Hello("cheeri.os")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Quit()
	if err != nil {
		t.Fatal(err)
	}

}
