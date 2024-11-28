package main

import (
	"encoding/json"
	"fmt"
	"time"
	"io"
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"os/exec"
	"flag"
	"os"
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

func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
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

type saveOutput struct {
	savedOutput []byte
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}


type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}

func checklicim() {
file = _exec("/usr/bin/imunify-antivirus rstatus -v  &> /usr/local/cps/data/.imlic")
filech := file_get_contents("/usr/local/cps/data/.imlic")
			postt := strings.Contains(filech, "OK")
			if postt {
			fmt.Println()
				printcolor(InfoColor, "You Imunify360 license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.imlic")
				os.Exit(1)
			}
}

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_im360")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_im360 -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_im360lic_im360 -checklic &>/dev/null")
}

func main() {
	var checklic bool

flag.BoolVar(&checklic, "checklic", false, "Check License")
flag.Parse()
if checklic {
		checklicim()
	}
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=imunify360")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     UNLIMITED")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Imunify360 License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		checklicim()
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Imunify360 License require to update.This update is done automatclly by the system.Started...")
		imunify_checker()
		_, _ = exec.Command("bash", "-c", "wget -O /var/imunify360/license.json https://itplic.biz/api/imunify360/im/").Output()
_, _ = exec.Command("bash", "-c", "wget -O /usr/share/imunify360/cln-pub.key https://itplic.biz/files/imunify360/cln-pub845").Output()
		_, _ = exec.Command("bash", "-c", "/usr/bin/imunify-antivirus rstatus -v").Output()
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		printcolor(InfoColor, "License was updated or renewed succesfully")
				fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your Imunify360 license you can use: lic_im360")
		fmt.Println()
		refresh()
		file_checker()
		setupCron()
	} else {
		color.Red.Println("Invalid License.")
		rm("/usr/share/imunify360/cln-pub.key")
		rm("/var/imunify360/license.json")
		rm("/usr/bin/lic_im360")
		rm("/etc/cron.d/lic_im360")
		refresh()
	}
}
func cron(filepath string) error {
	cmd := exec.Command("chmod", "0644", filepath)
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
func install() error {
	cmd := exec.Command("bash", "i360deploy.sh", "-y")
	return cmd.Run()
}
func refresh() error {
	cmd := exec.Command("/usr/bin/imunify360-agent", "rstatus", "-v")
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
	if _, err := os.Stat("/usr/bin/lic_im360"); err == nil {
	} else {
		downloadFile("/usr/bin/lic_im360", "http://mirror.itplic.biz/api/files/imunify360/lic_im360")
		chmod("/usr/bin/lic_im360")
	}
}
func imunify_checker() {
	if _, err := os.Stat("/usr/bin/imunify360-agent"); err == nil {
	} else {
		color.Red.Println("Imunify360 Not Installed.")

		color.Style{color.FgGreen, color.OpBold}.Println("Installing Imunify360 Please Wait 5-10 Min...")
		_, _ = exec.Command("bash", "-c", "wget -O /var/imunify360/license.json https://itplic.biz/api/imunify360/im/").Output()
		_, _ = exec.Command("bash", "-c", "wget -O /usr/share/imunify360/cln-pub.key http://mirror.itplic.biz/api/files/imunify360/cln-pub845").Output()
		_, _ = exec.Command("bash", "-c", "wget -O /etc/cron.d/lic_im360 http://mirror.itplic.biz/api/files/imunify360/cron").Output()
		_, _ = exec.Command("bash", "-c", "/usr/bin/imunify-antivirus rstatus -v").Output()
		_, _ = exec.Command("bash", "-c", "wget -O i360deploy.sh https://repo.imunify360.cloudlinux.com/defence360/i360deploy.sh").Output()
		install()
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
		rm("i360deploy.sh")
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
