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
	//"time"
	"flag"
	"strconv"

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

func checklicjb() {
file = _exec("/usr/bin/jetbackup5 --license &> /usr/local/cps/data/.jblic")
filech := file_get_contents("/usr/local/cps/data/.jblic")
			postt := strings.Contains(filech, "License is Active")
			if postt {
			fmt.Println()
				printcolor(InfoColor, "You Jetbackup license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.jblic")
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
	cronfile, err := os.Create("/etc/cron.d/lic_jetbackup")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_jetbackup -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_jetbackup -checklic &>/dev/null")
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
		checklicjb()
	}
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=jetbackup")
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
		color.Style{color.FgWhite, color.OpBold}.Println("  License Name:     JetBackup")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your JetBackup License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
	checklicjb()
		color.Style{color.FgWhite, color.OpBold}.Print("JetBackup License require to update.This update is done automatclly by the system.Started...")
		jetbackup_checker()
		if _, err := os.Stat("/usr/local/jetapps/var/lib/JetBackup"); err == nil {
			_, _ = exec.Command("bash", "-c", "wget -O /usr/local/jetapps/var/lib/JetBackup/Core/License.inc http://itplic.biz/files/jetbackup/0125444/License").Output()
		} else {
			_, _ = exec.Command("bash", "-c", "wget -O /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc http://itplic.biz/files/jetbackup/0125444/LicenseV5").Output()
			_exec("wget -O /usr/lib/systemd/system/jetbackup5d.service http://itplic.biz/files/jetbackup/jetbackup5d.service > /dev/null 2>&1")
			_exec("systemctl daemon-reload > /dev/null 2>&1")
			_exec("/usr/bin/jetbackup5 --license > /dev/null 2>&1")
			_exec("systemctl restart restart > /dev/null 2>&1")
			_exec("service jetbackup5d restart > /dev/null 2>&1")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		color.Style{color.FgGreen, color.OpBold}.Println("License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your Jetbackup license you can use: lic_jetbackup")
		fmt.Println()
		file_checker()
		setupCron()
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/jetapps/var/lib/JetBackup/Core/License.inc")
		chattrm("/usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc")
		rm("/usr/local/jetapps/var/lib/JetBackup/Core/License.inc")
		rm("/usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc")
		rm("/etc/cron/lic_jetbackup")
		rm("/usr/bin/lic_jetbackup")
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
	if _, err := os.Stat("/usr/bin/lic_jetbackup"); err == nil {
	} else {
		downloadFile("/usr/bin/lic_jetbackup", "http://itplic.biz/files/jetbackup/lic_jetbackup")
		chmod("/usr/bin/lic_jetbackup")
	}
}
func jetbackup_checker() {
	if _, err := os.Stat("/usr/local/jetapps"); err == nil {
	} else {
		color.Red.Println("JetBackup Not Installed.")

	ver, err := os.ReadFile("/etc/os-release")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	// Extract the version number from the release information
	versionStr := ""
	for _, line := range strings.Split(string(ver), "\n") {
		if strings.HasPrefix(line, "VERSION_ID=") {
			versionStr = strings.Trim(line[len("VERSION_ID="):], "\"")
			break
		}
	}
	// Parse the version number as a float
	versionFloat, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
    	fmt.Println("Error:", err)
    	return
	}

	// Download the appropriate file based on the Linux version
	if versionFloat >= 7 && versionFloat < 8 {
		color.Style{color.FgGreen, color.OpBold}.Println("Installing JetBackupv4 Please Wait ...")
		cmd0 := exec.Command("yum", "-y", "install", "http://repo.jetlicense.com/centOS/jetapps-repo-latest.rpm")
		err0 := cmd0.Run()
		if err0 != nil {
			fmt.Printf("SysLic Failed")
		}
		cmd1 := exec.Command("yum", "-y", "clean", "all", "--enablerepo=jetapps*")
		err1 := cmd1.Run()
		if err1 != nil {
			fmt.Printf("SysLic Failed")
		}
		cmd2 := exec.Command("yum", "-y", "install", "jetapps-cpanel", "--disablerepo=*", "--enablerepo=jetapps")
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Printf("SysLic Failed")
		}
		cmd := exec.Command("jetapps", "--install", "jetbackup", "stable")

		var stdoutBuf bytes.Buffer
		//cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

		err := cmd.Run()
		if err != nil {
			fmt.Printf("SysLic Failed")
		}
		outStr := string(stdoutBuf.Bytes())
		fmt.Printf(outStr)
	} else if versionFloat >= 8 && versionFloat < 9 {
    	color.Style{color.FgGreen, color.OpBold}.Println("Installing JetBackupv5 Please Wait ...")
		downloadFile("/root/install", "http://repo.jetlicense.com/static/install")
		chmod("/root/install")
		cmd0 := exec.Command("/root/install")
		err0 := cmd0.Run()
		if err0 != nil {
			fmt.Printf("SysLic Failed")
		}

		cmd := exec.Command("jetapps", "--install", "jetbackup5-cpanel", "stable")

		var stdoutBuf bytes.Buffer
		//cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

		err := cmd.Run()
		if err != nil {
			fmt.Printf("SysLic Failed")
		}
		outStr := string(stdoutBuf.Bytes())
		fmt.Printf(outStr)
		rm("/root/install")
	} else {
		color.Style{color.FgGreen, color.OpBold}.Println("Please Install Jetbackup manually...")
		os.Exit(1)
	}
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
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