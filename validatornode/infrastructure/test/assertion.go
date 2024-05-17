package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Assert(t testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

func AssertThatMessageIsLogged(t testing.TB, loggedMessageStructs []struct{ Msg string }, expectedMessages ...string) {
	isLoggedByExpectedMessage := make(map[string]bool)
	var loggedMessages []string
	for _, call := range loggedMessageStructs {
		loggedMessages = append(loggedMessages, call.Msg)
		for _, expectedMessage := range expectedMessages {
			if strings.Contains(call.Msg, expectedMessage) {
				isLoggedByExpectedMessage[expectedMessage] = true
				break
			}
		}
	}
	for _, expectedMessage := range expectedMessages {
		_, ok := isLoggedByExpectedMessage[expectedMessage]
		Assert(t, ok, fmt.Sprintf("expected message was not logged: %s\nlogged messages: %v", expectedMessage, loggedMessages))
	}
}
