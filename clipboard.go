package main

import (
	"time"

	"golang.design/x/clipboard"
)

func copyToClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))

	go func() {
		time.Sleep(15 * time.Second)
		clipboard.Write(clipboard.FmtText, []byte(""))
	}()
}