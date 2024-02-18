package main

import (
	"fmt"
	"os"

	"github.com/peterhellberg/sseclient"
)

func main() {

	events, err := sseclient.OpenURL("http://localhost:8080/usage?organization=gantry-sample&span=10")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for event := range events {
		fmt.Println(event.Name, event.Data)
	}
}
