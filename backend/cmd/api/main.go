// Command api is the Lambda entrypoint placeholder for the health-log service.
// HTTP routing is introduced during the protected API implementation phase.
package main

import (
	"log"

	"github.com/divanovSXM/health-log/internal/config"
)

func main() {
	if _, err := config.LoadFromOS(); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	log.Print("health-log configuration validated")
}
