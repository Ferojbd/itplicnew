package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"flag"
	//"strconv"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)

type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}

type saveOutput struct {
	savedOutput []byte
}
const (
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
	InfoColor  = "\033[1;32m%s\033[0m"
)
func printcolor(color string, str string) {
	fmt.Printf(color, str)
	fmt.Println()
}
func str_exists(str string, subject string) bool {
	return strings.Contains(subject, str)
}

func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

var file string

func checklicaa() {
file = _exec("service bt status &> /usr/local/cps/data/.aaplic")
filech := file_get_contents("/usr/local/cps/data/.aaplic")
			postt := strings.Contains(filech, "running")
			if postt {
			fmt.Println()
				printcolor(InfoColor, "You Aapanel license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.aaplic")
				os.Exit(1)
			} 
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
func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_aapanel")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_aapanel -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_aapanel -checklic &>/dev/null")
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
		checklicaa()
	}
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=aapanel")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("  Our Website:      ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["domain_name"])
		color.Style{color.FgWhite, color.OpBold}.Println("  License Name:     AaPanel")
		color.Style{color.FgWhite, color.OpBold}.Println("  License Version:  v3.41")
		host, _ := os.Hostname()
		color.Style{color.FgWhite, color.OpBold}.Printf("  Hostname:         ")
		color.Style{color.FgWhite, color.OpBold}.Println(host)
		color.Style{color.FgWhite, color.OpBold}.Printf("  Server IP:        ")
		curl := exec.Command("curl", "-s", "https://ipinfo.io/ip")
		out, err := curl.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		ip := string(out)
		color.Style{color.FgWhite, color.OpBold}.Println(ip)
		color.Style{color.FgWhite, color.OpBold}.Println("----------------------------------------------------------------------")
		color.Style{color.FgWhite, color.OpBold}.Printf("Today is ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Printf("Your AaPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
	checklicaa()
		color.Style{color.FgWhite, color.OpBold}.Print("AaPanel License require to update.This update is done automatclly by the system.Started...")
		aapanel_checker()
		_exec("chattr -i /www/server/panel/data/plugin.json  > /dev/null 2>&1")
		_exec("cp /www/server/panel/data/plugin.json /www/server/panel/data/pluggin.json  > /dev/null 2>&1")
		cmd1 := exec.Command("sed", "-i", "s|\"endtime\": -1|\"endtime\": 999999999999|g", "/www/server/panel/data/plugin.json")
    if err := cmd1.Run(); err != nil {
        panic(err)
    }

    // Command to replace "-1" with "0" for "pro" key
    cmd2 := exec.Command("sed", "-i", "s|\"pro\": -1|\"pro\": 0|g", "/www/server/panel/data/plugin.json")
    if err := cmd2.Run(); err != nil {
        panic(err)
    }

    // Command to set immutable attribute for the file
    cmd3 := exec.Command("chattr", "+i", "/www/server/panel/data/plugin.json")
    if err := cmd3.Run(); err != nil {
        panic(err)
    }
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		color.Style{color.FgGreen, color.OpBold}.Println("License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your AaPanel license you can use: lic_aapanel")
		fmt.Println()
		file_checker()
		setupCron()
	} else {
		color.Red.Println("Invalid License.")
		_exec("chattr -i /www/server/panel/data/plugin.json  > /dev/null 2>&1")
		cmd4 := exec.Command("sed", "-i", "s|\"endtime\": 999999999999|\"endtime\": -1|g", "/www/server/panel/data/plugin.json")
    if err := cmd4.Run(); err != nil {
        panic(err)
    }

    // Command to replace "-1" with "0" for "pro" key
    cmd5 := exec.Command("sed", "-i", "s|\"pro\": 0|\"pro\": -1|g", "/www/server/panel/data/plugin.json")
    if err := cmd5.Run(); err != nil {
        panic(err)
    }

    // Command to set immutable attribute for the file
    cmd6 := exec.Command("chattr", "+i", "/www/server/panel/data/plugin.json")
    if err := cmd6.Run(); err != nil {
        panic(err)
    }
		rm("/etc/cron/lic_aapanel")
		rm("/usr/bin/lic_aapanel")
	}
}

func cron(filepath string) error {
	cmd := exec.Command("chmod", "0644", filepath)
	return cmd.Run()
}
func chattrm(filepath string) error {
	cmd := exec.Command("chattr", "-i", "-a", filepath)
	return cmd.Run()
}
func run(filepath string) error {
	// run shell
	cmd := exec.Command(filepath)
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
	if _, err := os.Stat("/usr/bin/lic_aapanel"); err == nil {
	} else {
		downloadFile("/usr/bin/lic_aapanel", "http://mirror.itplic.biz/api/files/aapanel/lic_aapanel")
		chmod("/usr/bin/lic_aapanel")
	}
}
func aapanel_checker() {
	if _, err := os.Stat("/www/server/panel/"); err == nil {
	} else {
		color.Red.Println("AaPanel Not Installed.You need to install Aapanel First!")
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
	//_, err = io.Copy(out, resp.Body)
	if err != nil { return err }

	return nil
}