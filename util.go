package main

import (
	"fmt"
	"net/http"
	"strings"
)

func SplitName(name string) (firstName, lastName string) {
	s := strings.Split(name, " ")
	if len(s) <= 1 {
		return name, ""
	}
	if len(s) == 2 {
		return s[0], s[1]
	}

	return strings.Join(s[:len(s)-1], " "), s[len(s)-1]
}

func WriteError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%d %s", code, http.StatusText(code))))
}
