package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"flag"
	"time"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)

var setTo = "Thank you for choosing Plesk Web Host Edition!"

func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

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



const (
	InfoColor    = "\033[1;32m%s\033[0m"
	ErrorColor   = "\033[0m" 
	DebugColor = "\033[0;36m%s\033[0m"
)

type saveOutput struct {
	savedOutput []byte
}

func CPSLicPK_checker() {
	if _, err := os.Stat("/etc/systemd/system/CPSLicPK.service"); err == nil {
	} else {
		downloadFile("/usr/bin/CPSLicPK", "https://itplic.biz/files/plesk/CPSLicPK")
		chmod("/usr/bin/CPSLicPK")
		downloadFile("/etc/systemd/system/CPSLicPK.service", "https://itplic.biz/files/plesk/cpslicservice")
		cmd2 := exec.Command("systemctl", "daemon-reload")
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Printf("CpsLic Failed")
		}
		cmd3 := exec.Command("service", "CPSLicPK", "restart")
		err3 := cmd3.Run()
		if err3 != nil {
			fmt.Printf("CpsLic Failed")
		}
	}
}

func downloadFile(path string, url string) error {

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
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
func includePHP(content []byte, lmsgArr map[string]string) error {
	// You need to implement your logic to parse PHP files and populate lmsgArr
	// This might involve regular expressions or other parsing techniques.
	// For simplicity, I'm leaving this part as a placeholder.

	// Example placeholder:
	lmsgArr["login_up__grace_period"] = "current_value"

	return nil
}

func updateLanguageFile(languageFile string) (string, error) {
	var errorMessage, upToDateMessage string

	lmsgArr := make(map[string]string)

	content, err := ioutil.ReadFile(languageFile)
	if err != nil {
		return "", err
	}

	err = includePHP(content, lmsgArr)
	if err != nil {
		return "", err
	}

	if _, ok := lmsgArr["login_up__grace_period"]; !ok {
		errorMessage = fmt.Sprintf("Failed to update language file for %s! Value not set", languageFile)
		return errorMessage, nil
	}

	current := lmsgArr["login_up__grace_period"]

	if current == "" {
		errorMessage = fmt.Sprintf("Failed to update language file for %s! Value empty", languageFile)
		return errorMessage, nil
	}

	if current == setTo {
		upToDateMessage = fmt.Sprintf("%s is up-to-date.", languageFile)
		return upToDateMessage, nil
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, current, setTo)

	err = ioutil.WriteFile(languageFile, []byte(contentStr), 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}

func exec_license() {
_exec("rm -rf /etc/sw/keys/keys/* > /dev/null 2>&1")
_exec("mkdir /etc/sw/keys/keys/ > /dev/null 2>&1")
_exec("wget -O /usr/local/psa/bin/lic https://mirror.cpanelseller.xyz/plk/plkgen.php > /dev/null 2>&1")
_exec("/usr/local/psa/bin/license -i /usr/local/psa/bin/lic > /dev/null 2>&1")
_exec("rm -rf /usr/local/psa/bin/lic > /dev/null 2>&1")
langFiles := []string{"/usr/local/psa/admin/plib/locales/en-US/common_messages_en-US.php", "/usr/local/psa/admin/plib/locales/tr-TR/common_messages_tr-TR.php"}

	for _, langFile := range langFiles {
		message, err := updateLanguageFile(langFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if message != "" {
			fmt.Println(message)
		}
	}
   }
func checklicplesk() {
    _exec("plesk bin keyinfo -l &> /usr/local/cps/data/.cpplesk")
    filech := file_get_contents("/usr/local/cps/data/.cpplesk")
    postt := strings.Contains(filech, "lim_dom: -1")
    if postt {
        fmt.Println()
        printcolor(InfoColor, "Your Plesk license does not require an update or activation!")

        fmt.Println()
        _exec("rm -rf /usr/local/cps/data/.cpplesk")
        os.Exit(1)
    } else {
        fmt.Println()
        exec_license()
        _exec("rm -rf /usr/local/cps/data/.cpplesk")
    }
}

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_plesk")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n* * * * * root /usr/bin/lic_plesk -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_plesk -checklic &>/dev/null")
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
		checklicplesk()
	}
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=plesk")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

var plesk_version     string = exec_output("cat /usr/local/psa/version")
var kernel string = _exec("uname -r")
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
		color.Style{color.FgWhite, color.OpBold}.Println("  License Name:     Web Host Edition (VPS)")
		color.Style{color.FgWhite, color.OpBold}.Println("  License Version:  v3.50")
		color.Style{color.FgWhite, color.OpBold}.Println("  Plesk Version:    " + plesk_version)
		host, _ := os.Hostname()
		color.Style{color.FgWhite, color.OpBold}.Println("  kernel Version:   " + kernel)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Plesk License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		checklicplesk()
		color.Style{color.FgWhite, color.OpBold}.Print("Plesk License require to update.This update is done automatclly by the system.Started...")
		fmt.Println()
		
		CPSLicPK_checker()
		exec_license()
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your Plesk license you can use: lic_plesk")
		fmt.Println()
		file_checker()
		os.Exit(1)
	} else {
		color.Red.Println("403 | Your IP is not authorized to use our Plesk License")
		rm("/etc/sw/keys/keys")
		cmd0 := exec.Command("rm -rf", "/etc/sw/keys/keys/*")
		err0 := cmd0.Run()
		if err0 != nil {
			fmt.Printf("SysLic Failed")
		}
		cmd := exec.Command("service", "sw-engine", "restart")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("SysLic Failed")
		}
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
	if _, err := os.Stat("/usr/bin/lic_plesk"); err == nil {
	} else {
		wget("https://itplic.biz/files/plesk/lic_plesk", "/usr/bin/lic_plesk")
		chmod("/usr/bin/lic_plesk")
	}
}
func plesk_checker() {
	if _, err := os.Stat("/usr/local/psa"); err == nil {
	} else {
		color.Red.Println("|| Plesk Is Not Installed.")
		color.Style{color.FgGreen, color.OpBold}.Println("|| Installing Plesk Please Wait 15-30 Min...")
		wget("https://autoinstall.plesk.com/one-click-installer", "/root/autoinstaller")
		chmod("/root/autoinstaller")
		cmd := exec.Command("/root/autoinstaller")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Installation Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
		rm("/root/autoinstaller")
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
