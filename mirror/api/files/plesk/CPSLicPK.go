package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// runCommand executes the given command and returns its output and any error.
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command %s %v: %v\nCombined Output:\n%s", command, args, err, combinedOutput)
	}
	return string(combinedOutput), nil
}

func main() {
	// Handle interrupt signals for graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal. Exiting...")
		os.Exit(0)
	}()

	// Main loop
	for {
		log.Println("Starting...")

		// Execute the first command
		output1, err := runCommand("plesk", "bin", "keyinfo", "-l")
		if err != nil {
			log.Println("Error executing command:", err)
			log.Println("Skipping further processing...")
			time.Sleep(25 * time.Second) // Sleep for 25 seconds before next iteration
			continue
		}

		// Check if the output does not contain "lim_dom: -1"
		if !strings.Contains(output1, "lim_dom: -1") {
			// Execute /usr/bin/lic_plesk
			log.Println("Plesk is not running. Starting Plesk...")
			if _, err := runCommand("/usr/bin/lic_plesk"); err != nil {
				log.Println("Error starting Plesk:", err)
			} else {
				log.Println("Plesk started successfully.")
			}
		}

		log.Println("End...")
		// Sleep for 25 seconds
		time.Sleep(25 * time.Second)
	}
}
