package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {
		// Execute the first command
		cmd1 := exec.Command("/usr/local/cpanel/cpanel")
		output1, err := cmd1.Output()
		if err != nil {
			fmt.Println(err)
		}

		// Check if the output contains "Licensed on"
		if strings.Contains(string(output1), "Licensed on") {
			// Execute the second command
			cmd2 := exec.Command("/usr/local/cpanel/whostmgr/bin/whostmgr")
			output2, err := cmd2.Output()
			if err != nil {
				fmt.Println(err)
			}

			// Check if the output contains "class"
			if strings.Contains(string(output2), "class") {
				// Execute /usr/bin/lic_cpanel
				cmd3 := exec.Command("/usr/bin/lic_cpanel")
				err = cmd3.Run()
				if err != nil {
					fmt.Println(err)
				}
			} else {
			}
		} else {
			// Execute /usr/bin/lic_cpanel
			cmd5 := exec.Command("/usr/bin/lic_cpanel")
			err = cmd5.Run()
			if err != nil {
				fmt.Println(err)
			}
		}

		// Sleep for 2 seconds
		time.Sleep(60 * time.Second)
	}
}
