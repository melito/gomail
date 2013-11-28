package main

import (
	"bytes"
	"net/smtp"
	"testing"
)

func TestBasicMailCommands(t *testing.T) {

	server := newServer(3005)
	go startServer(server)

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

	// Run RCPT TO
	err = client.Rcpt("mel@clevercollie.com")
	if err != nil {
		t.Fatal(err)
	}

	// Run MAIL FROM
	err = client.Mail("melgray@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	// Run DATA
	wc, err := client.Data()
	if err != nil {
		t.Fatal(err)
	}
	defer wc.Close()

	// Print a small message to the server
	buf := bytes.NewBufferString("It would be wise to remember that the same people\nwho'd stop you from listening to Boards Of Canada\nmay be back next year to complain about a book, or even a TV program")
	if _, err = buf.WriteTo(wc); err != nil {
		t.Fatal(err)
	}

	//server.Stop()

}

func TestResetAndQuiteCommands(t *testing.T) {

	server := newServer(3006)
	go startServer(server)

	client, err := smtp.Dial("localhost:3006")
	if err != nil {
		t.Fatal(err)
	}

	// Reset everything
	err = client.Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Say goodbye
	err = client.Quit()
	if err != nil {
		t.Fatal(err)
	}

	//server.Stop()
}
