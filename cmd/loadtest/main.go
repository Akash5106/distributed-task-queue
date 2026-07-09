package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	for i := 1; i <= 1000; i++ {
		body := []byte(`{"payload":"benchmark task"}`)

		_, err := http.Post(
			"http://localhost:8080/tasks",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Submitted 1000 tasks")
}
