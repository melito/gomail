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

	// Say helo
	err = client.Hello("cheeri.os")
	if err != nil {
		t.Fatal(err)
	}

	// Run vrfy
	err = client.Verify("melgray@gmail.com")
	if err == nil {
		t.Fatal("We didn't receive an error when issuing a VRFY")
	}

	// Run RCPT to
	err = client.Rcpt("mel@clevercollie.com")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Mail("melgray@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Say goodbye
	err = client.Quit()
	if err != nil {
		t.Fatal(err)
	}

}
