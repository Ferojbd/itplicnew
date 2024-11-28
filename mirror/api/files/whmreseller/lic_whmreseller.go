package main

import (
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"strings"
	"os/exec"
	"os"
	//"net"
	"bytes"
	"io"
	"runtime"
	"log"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)
const (
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
	InfoColor  = "\033[1;32m%s\033[0m"
)
func printcolor(color string, str string) {
	fmt.Printf(color, str)
	fmt.Println()
}

type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}
func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_whmreseller")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_whmreseller >/dev/null 2>&1\n@reboot root /usr/bin/lic_whmreseller &>/dev/null")
}
func main() {

	resp, err := http.Get("http://trlisans.org/api/getinfo?key=whmreseller")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)
	
	var f Data
	err = json.Unmarshal(byteResult, &f)
	var res map[string]interface{}
	json.Unmarshal([]byte(byteResult), &res)
	if f.Status == "success" {
	CallClear()
	ascii := figlet4go.NewAsciiRender()

	options := figlet4go.NewRenderOptions()
	options.FontName = "slant"
	ascii.LoadFont("/usr/local/go/bin")
	str := fmt.Sprint(res["brand_name"])
	renderStr, _ := ascii.RenderOpts(str, options)
	color.Style{color.FgWhite, color.OpBold}.Printf(renderStr)
	color.Style{color.FgWhite, color.OpBold}.Println("---------------------- Licensing System Started ----------------------")
		color.Style{color.FgWhite, color.OpBold}.Printf("|Our Website:      ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["domain_name"])
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     WhmReseller")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.50")
		host, _ := os.Hostname()
		color.Style{color.FgWhite, color.OpBold}.Printf("|Hostname:         ")
		color.Style{color.FgWhite, color.OpBold}.Println(host)
		color.Style{color.FgWhite, color.OpBold}.Printf("|Server IP:        ")
		curl := exec.Command("curl", "-s", "https://ipinfo.io/ip")
		out, err := curl.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		ip := string(out)
		color.Style{color.FgWhite, color.OpBold}.Println(ip)
		color.Style{color.FgWhite, color.OpBold}.Printf("")
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgWhite, color.OpBold}.Printf("Today is ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Printf("Your WhmReseller License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
	color.Style{color.FgWhite, color.OpBold}.Print("WhmReseller require to update.This update is done automatclly by the system.Started...")
	whmreseller_checker()
    downloadFile("/root/subreseller.cpp", "http://trlisans.org/api/files/whmreseller/subreseller27")
	cmd2 := exec.Command("g++", "/root/subreseller.cpp", "-o", "/usr/local/cpanel/whostmgr/cgi/whmreseller/subreseller.cgi")
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Printf("lic Failed")
	}
	color.Style{color.FgGreen, color.OpBold}.Println("OK")
	printcolor(InfoColor, "License was updated or renewed succesfully")
				fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your cPanel license you can use: lic_whmreseller")
	rm("/root/subreseller.cpp")
	setupCron()
	file_checker()
	} else {
		color.Red.Println("Invalid License.")
		rm("/root/subreseller.cpp")
		rm("/usr/local/cpanel/whostmgr/cgi/whmreseller/subreseller.cgi")
	}
}
func cron(filepath string) error {
      cmd := exec.Command("chmod", "0644", filepath)
      return cmd.Run()
}
func run(filepath string) error {
      // run shell
      cmd := exec.Command(filepath)
      return cmd.Run()
}
func wget(url, filepath string) error {
      // run shell `wget URL -O filepath`
      cmd := exec.Command("wget", url, "-O", filepath)
      return cmd.Run()
}
func chmod(filepath string) error {
      cmd := exec.Command("chmod", "+x", filepath)
      return cmd.Run()
}
func rm(filepath string) error {
      cmd := exec.Command("rm", "-rf", filepath)
      return cmd.Run()
}
var clear map[string]func() //create a map for storing clear funcs

func init() {
    clear = make(map[string]func()) //Initialize it
    clear["linux"] = func() { 
        cmd := exec.Command("clear") //Linux example, its tested
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested 
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func CallClear() {
    value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
    if ok { //if we defined a clear func for that platform:
        value()  //we execute it
    } else { //unsupported platform
        panic("Your platform is unsupported! I can't clear terminal screen :(")
    }
}

func file_checker() {
   if _, err := os.Stat("/usr/bin/lic_whmreseller"); err == nil {
   } else {
      downloadFile("/usr/bin/lic_whmreseller", "http://trlisans.org/api/files/whmreseller/lic_whmreseller")
	  chmod("/usr/bin/lic_whmreseller")
   }
}
func whmreseller_checker() {
   if _, err := os.Stat("/usr/local/cpanel/whostmgr/docroot/cgi/whmreseller/index.cgi"); err == nil {
   } else {
    color.Red.Println("WhmReseller Not Installed.")
	color.Style{color.FgGreen, color.OpBold}.Println("Installing WhmReseller Please Wait...")
		// Change directory to /usr/local/cpanel/whostmgr/docroot/cgi
	if err := os.Chdir("/usr/local/cpanel/whostmgr/docroot/cgi"); err != nil {
		log.Fatal(err)
	}

	// Download install.cpp file from the given URL
	downloadFile("install.cpp", "https://deasoft.com/whmreseller/install.cpp")

	// Compile install.cpp and create an executable named install
	compileCmd := exec.Command("g++", "install.cpp", "-o", "install")
	if err := compileCmd.Run(); err != nil {
		log.Fatal(err)
	}

	// Set permissions for the install executable
	chmodCmd := exec.Command("chmod", "700", "install")
	if err := chmodCmd.Run(); err != nil {
		log.Fatal(err)
	}

	// Run the install executable
	cmd := exec.Command("./install")

	var stdoutBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("lic Failed")
	}
	outStr := string(stdoutBuf.Bytes())
	fmt.Printf(outStr)
	color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
	rm("install")
	rm("install.cpp")
   }
}
func getData(fileurl string) string {
	resp, err := http.Get(fileurl)
	if err != nil {
		fmt.Println("Unable to get Data")
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)

	}
	data := string(html[:])
	data = strings.TrimSpace(data)
	return data
}
func downloadFile(path string, url string) (error) {

	// Create the file
	out, err := os.Create(path)
	if err != nil { return err }
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil { return err }
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil { return err }

	return nil
}

