package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMailDirDetection(t *testing.T) {

	if mailDirPath() == "" {
		t.Fatal("mailDirPath was nil")
	}

	if mailDirPath() != "./Maildir" {
		t.Fatal("mailDirPath was incorrect default")
	}

	os.Setenv(maildirPathEnv, "/tmp/MaildirBlah")
	if mailDirPath() != "/tmp/MaildirBlah" {
		t.Fatal("mailDirPath didn't change the directory path")
	}

}

func TestDirectoryCreation(t *testing.T) {

	os.RemoveAll(mailDirPath())

	if maildirExists() {
		t.Fatal("Maildir shouldn't exist yet")
	}

	users := []string{"melgray@gmail.com", "mel@clevercollie.com"}
	for _, user := range users {
		createMailDirForUser(user)
	}

	if maildirExists() != true {
		t.Fatal("Maildir should exist now")
	}

	os.RemoveAll(mailDirPath())

}

func TestUniqueFileName(t *testing.T) {

	a := createUniqueFileName()
	b := createUniqueFileName()

	if a == b {
		t.Fatal("Unique file names collided")
	}
}

func TestCreateUser(t *testing.T) {

	os.RemoveAll(mailDirPath())

	path := filepath.Join(pathToMailDirForEmail("melgray@gmail.com"))

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("User directory already existed")
	}

	err := addUser("melgray@gmail.com", "123456")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("User directory didn't exist")
	}

	passwdPath := filepath.Join(path, "passwd")
	if _, err := os.Stat(passwdPath); os.IsNotExist(err) {
		t.Fatal("User passwd file didn't exist")
	}

	os.RemoveAll(mailDirPath())

}

func TestUserAuth(t *testing.T) {

	err := addUser("melgray@gmail.com", "123456")
	if err != nil {
		t.Fatal(err)
	}

	if !authUser("melgray@gmail.com", "123456") {
		t.Fatal("User authentication not working")
	}
}
