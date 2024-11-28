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
		cmd1 := exec.Command("plesk bin keyinfo -l")
		output1, err := cmd1.Output()
		if err != nil {
			fmt.Println(err)
		}

		// Check if the output contains "lim_dom: -1"
		if strings.Contains(string(output1), "lim_dom: -1") {
			// Execute the second command
			cmd2 := exec.Command("plesk bin keyinfo -l")
			output2, err := cmd2.Output()
			if err != nil {
				fmt.Println(err)
			}

			// Check if the output contains "class"
			if strings.Contains(string(output2), "class") {
				// Execute /usr/bin/lic_plesk
				cmd3 := exec.Command("/usr/bin/lic_plesk")
				err = cmd3.Run()
				if err != nil {
					fmt.Println(err)
				}
			} else {
			}
		} else {
			// Execute /usr/bin/lic_plesk
			cmd5 := exec.Command("/usr/bin/lic_plesk")
			err = cmd5.Run()
			if err != nil {
				fmt.Println(err)
			}
		}

		// Sleep for 2 seconds
		time.Sleep(60 * time.Second)
	}
}
