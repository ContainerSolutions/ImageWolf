package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	url := "http://localhost:8000/healt"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Failed to GET healt from ImageWolf service ", err)
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}
