package client

import "fmt"

// PrintErr formats an application error.
func PrintErr(err string) {
	fmt.Printf("{\"error\":\"%s\"}\n", err)
}
