package main

import (
	"bytes"
	"code.google.com/p/go.crypto/scrypt"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	maildirPathEnv = "GOMAIL_MAILDIR_PATH"
	salt           = "bGrQjAI810a81janJJJAHSBCXXXXTZcx"
)

func maildirExists() bool {
	if _, err := os.Stat(mailDirPath()); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func mailDirPath() string {
	path := os.Getenv(maildirPathEnv)
	if path == "" {
		return "./Maildir"
	}
	return path
}

func createMailDirForUser(emailStr string) {
	path := pathToMailDirForEmail(emailStr)
	os.MkdirAll(filepath.Join(path, "new"), 0700)
	os.MkdirAll(filepath.Join(path, "cur"), 0700)
}

func pathToMailDirForEmail(emailStr string) string {
	user := usernameForAddress(emailStr)
	host := hostnameForAddress(emailStr)
	return filepath.Join(mailDirPath(), host, user)
}

// createUnqiqueFileName doesn't adhere to the traditional Maildir message format
// as closely as it should, but for the time being this should be good enough
func createUniqueFileName() string {

	// Get a timestamp
	tv := syscall.Timeval{}
	syscall.Gettimeofday(&tv)
	left := fmt.Sprintf("%d.%d", tv.Sec, tv.Usec)

	// Just generate a random number for now
	b := make([]byte, 16)
	rand.Read(b)
	middle := fmt.Sprintf("%x", b)

	// The right dot should be the hostname
	right, _ := os.Hostname()

	// Put the pieces together
	combined := []string{left, middle, right}
	return strings.Join(combined, ".")
}

func addUser(address string, password string) error {
	createMailDirForUser(address)
	passwd, _ := mkpasswd(password)

	// base64 encode the scrypt'd
	passwdStr := base64.StdEncoding.EncodeToString(passwd)

	passwdFile, err := os.Create(filepath.Join(pathToMailDirForEmail(address), "passwd"))
	if err != nil {
		return err
	}

	passwdFile.WriteString(passwdStr)
	return nil
}

func authUser(address string, password string) bool {

	passwdFile, err := os.Open(filepath.Join(pathToMailDirForEmail(address), "passwd"))

	userPasswd, err := ioutil.ReadAll(passwdFile)
	if err != nil {
		return false
	}

	inPasswd, _ := mkpasswd(password)
	inPasswdStr := base64.StdEncoding.EncodeToString(inPasswd)

	if inPasswdStr == string(userPasswd) {
		return true
	}

	return false
}

func usernameForAddress(emailStr string) string {
	pieces := strings.Split(emailStr, "@")
	user := pieces[0]
	user = strings.Split(user, "+")[0]
	return user
}

func hostnameForAddress(emailStr string) string {
	pieces := strings.Split(emailStr, "@")
	host := strings.Join(pieces[1:], "")
	return host
}

func mkpasswd(password string) ([]byte, error) {

	saltBuf := bytes.NewBufferString(salt)

	dk, err := scrypt.Key([]byte(password), saltBuf.Bytes(), 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return dk, nil

}
