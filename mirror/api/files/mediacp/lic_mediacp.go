package main

import (
	"encoding/json"
	"fmt"
	//"time"
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"
	"os/exec"
	"os"
	//"io"
	"flag"
	"runtime"
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

func checklicskc() {
_exec("kcarectl  &> /usr/local/cps/data/.mcplic")
    filech := file_get_contents("/usr/local/cps/data/.mcplic")
    postt := strings.Contains(filech, "enabled")
    if postt {
        fmt.Println()
        printcolor(InfoColor, "Your MediaCP license does not require an update or activation!")
        fmt.Println()
        _exec("rm -rf /usr/local/cps/data/.mcplic")
        os.Exit(1)
    } else {
        fmt.Println()
        _exec("/usr/bin/lic_mediacp")
        _exec("rm -rf /usr/local/cps/data/.mcplic")
    }
}

func exec_license() {
_exec("wget -O /usr/local/mediacp/htdocs/system/functions.php api.cpanelseller.xyz/files/mediacp/functions  > /dev/null 2>&1")

 }


type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}
type saveOutput struct {
	savedOutput []byte
}
func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_mediacp")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_mediacp >/dev/null 2>&1\n@reboot root /usr/bin/lic_mediacp &>/dev/null")
}
func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

func main() {
var checklic bool

flag.BoolVar(&checklic, "checklic", false, "Check License")
flag.Parse()
if checklic {
		checklicskc()
	}
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=mediacp")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     MediaCP")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.51")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your MediaCP License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		//checklicskc()
		color.Style{color.FgWhite, color.OpBold}.Print("MediaCP require to update.This update is done automatclly by the system.Started...")
exec_license()		
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		fmt.Println()
		printcolor(InfoColor, "License was updated or renewed succesfully")
				fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your MediaCP license you can use: lic_mediacp")
		fmt.Println()
		file_checker()
		setupCron()
		os.Exit(1)
	} else {
		color.Red.Println("403 | Your IP is not authorized to use our MediaCP License")
		_exec("rm -rf /usr/local/mediacp/htdocs/system/functions.php > /dev/null 2>&1")
		}
	}
func cron(filepath string) error {
	cmd := exec.Command("chmod", "0644", filepath)
	return cmd.Run()
}
func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}
func run(filepath string) error {
	// run shell
	cmd := exec.Command(filepath)
	return cmd.Run()
}
func chattrp(filepath string) error {
	cmd := exec.Command("chattr", "+i", "+a", filepath)
	return cmd.Run()
}
func chattrm(filepath string) error {
	cmd := exec.Command("chattr", "-i", "-a", filepath)
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
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
func file_checker() {
	if _, err := os.Stat("/usr/bin/lic_mediacp"); err == nil {
	} else {
		wget("http://itplic.biz/files/mediacp/lic_mediacp", "/usr/bin/lic_mediacp")
		chmod("/usr/bin/lic_mediacp")
	}
}
func kernelcare_checker() {
	if _, err := os.Stat("/usr/bin/kcarectl"); err == nil {
	} else {
		color.Red.Println("|| Kernelcare Is Not Installed.")
		color.Style{color.FgGreen, color.OpBold}.Println("|| Installing MediaCP Please Wait few Min...")
		cmd := exec.Command("wget  | bash")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Installation Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
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
