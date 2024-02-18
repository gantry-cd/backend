package main

import (
	"fmt"

	"github.com/gantrycd/backend/internal/utils/branch"
)

func main() {
	sample := "feature/#1"

	fmt.Print("Transpile1123: ", branch.Transpile1123(sample))
}
