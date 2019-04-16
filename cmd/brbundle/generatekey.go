package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func generateKey() {
	key := make([]byte, 32+12)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fail to generate key")
		os.Exit(1)
	}
	fmt.Printf("%s", base64.StdEncoding.EncodeToString(key))
}
