package error

import "fmt"

var (
	ErrDeploymentsNotFound = fmt.Errorf("deployment not found")
)
