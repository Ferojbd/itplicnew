package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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
func str_exists(str string, subject string) bool {
	return strings.Contains(subject, str)
}

func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

var file string

func checkliccpanel() {
	_exec("php /usr/local/cpanel/whostmgr/cgi/softaculous/cli.php -l &> /usr/local/cps/data/.softlic")
	filech := file_get_contents("/usr/local/cps/data/.softlic")
	postt := strings.Contains(filech, "Premium")
	if postt {
		fmt.Println()
		printcolor(InfoColor, "Your Softaculous license does not require an update or activation!")
		_exec("rm -rf /usr/local/cps/license")
		_exec("rm -rf /usr/local/cps/sanity")
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
		os.Exit(1)
	} else {
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
	}
}

func checklicplesk() {
	_exec("php /usr/local/softaculous/cli.php -l &> /usr/local/cps/data/.softlic")
	filech := file_get_contents("/usr/local/cps/data/.softlic")
	postt := strings.Contains(filech, "Premium")
	if postt {
		fmt.Println()
		printcolor(InfoColor, "Your Softaculous license does not require an update or activation!")
		_exec("rm -rf /usr/local/cps/license")
		_exec("rm -rf /usr/local/cps/sanity")
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
		os.Exit(1)
	} else {
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
	}
}

func checklicda() {
	_exec("php /usr/local/directadmin/plugins/softaculous/cli.php -l &> /usr/local/cps/data/.softlic")
	filech := file_get_contents("/usr/local/cps/data/.softlic")
	postt := strings.Contains(filech, "Premium")
	if postt {
		fmt.Println()
		printcolor(InfoColor, "Your Softaculous license does not require an update or activation!")
		_exec("rm -rf /usr/local/cps/license")
		_exec("rm -rf /usr/local/cps/sanity")
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
		os.Exit(1)
	} else {
		fmt.Println()
		_exec("rm -rf /usr/local/cps/data/.softlic")
	}
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
	cronfile, err := os.Create("/etc/cron.d/lic_softaculous")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */4 * * * root /usr/bin/lic_softaculous -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_softaculous -checklic &>/dev/null")
}

func _system(of string) string {
	var out bytes.Buffer
	cmd := exec.Command("bash", "-c", of)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	} else {
		lines := strings.Split(out.String(), "\n")
		if len(lines) > 1 {
			fmt.Println(lines[len(lines)-2])
			return lines[len(lines)-2]
		} else {
			return ""
		}
	}

}

func _exec(of string) string {
	//fmt.Println(of)
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

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

func file_exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func stripos(haystack string, needle string) bool {

	if strings.Index(strings.ToUpper(haystack), strings.ToUpper(needle)) > -1 {
		return true
	}

	return false
}

func main() {
	var checklic bool

	flag.BoolVar(&checklic, "checklic", false, "Check License")
	flag.Parse()
	if checklic {
		checkliccpanel()
	}
	var fleetssl bool
	var ssl_services bool
	var uninstall bool
	var upcp bool

	flag.BoolVar(&fleetssl, "fleetssl", false, "Install FleetSSL Premium")
	flag.BoolVar(&ssl_services, "ssl_services", false, "Install SSL on Hostname")
	flag.BoolVar(&upcp, "upcp", false, "Upgrade/Downgrade to the Supported cPanel Version")
	flag.BoolVar(&uninstall, "uninstall", false, "Remove Our License tem")
	flag.Parse()

	if upcp {
		update()
	}
	if fleetssl {
		fleet()
	}
	if ssl_services {
		installssl()
	}
	if uninstall {
		remove()
	}

	resp, err := http.Get("https://itplic.biz/api/getinfo?key=softaculous")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgWhite, color.OpBold}.Printf("Today is ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Softaculous License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Checking Softaculous License Files...")
		cpanel_checker()
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		checkliccpanel()
		checklicplesk()
		checklicda()
		color.Style{color.FgWhite, color.OpBold}.Print("Softaculous License require to update.This update is done automatclly by the system.Started...")

		cpanel := "/usr/local/cpanel/whostmgr/cgi/softaculous/cli.php"
		directadmin := "/usr/local/directadmin/plugins/softaculous/cli.php"
		plesk := "/usr/local/softaculous/cli.php"

		if fileExists(cpanel) {
			fmt.Println()
			fmt.Println("\x1B[32m\ncPanel detected. Installing license...\x1B[0m")
			fmt.Println()
			_, _ = exec.Command("bash", "-c", "chattr -i /usr/local/cpanel/whostmgr/cgi/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "/usr/local/cpanel/3rdparty/bin/php /usr/local/cpanel/whostmgr/cgi/softaculous/cli.php --refresh-license").Output()
			_, _ = exec.Command("bash", "-c", "wget -O /usr/local/cpanel/whostmgr/cgi/softaculous/enduser/license.php https://itplic.biz/api/softaculous?key=softaculous").Output()
			_, _ = exec.Command("bash", "-c", "chattr +i /usr/local/cpanel/whostmgr/cgi/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "//usr/local/cpanel/3rdparty/bin/php /usr/local/cpanel/whostmgr/docroot/cgi/softaculous/cli.php --enable_script --all").Output()
		}

		if fileExists(directadmin) {
			fmt.Println()
			fmt.Println("\x1B[32m\nDirectAdmin detected. Installing license...\x1B[0m")
			fmt.Println()
			_, _ = exec.Command("bash", "-c", "chattr -i /usr/local/directadmin/plugins/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "/usr/local/bin/php", "-d", "open_basedir=\"\"", "-d", "safe_mode=0", "-d", "disable_functions=\"\"", "/usr/local/directadmin/plugins/softaculous/cli.php", "--refresh-license").Output()
			_, _ = exec.Command("bash", "-c", "wget -O /usr/local/directadmin/plugins/softaculous/enduser/license.php https://itplic.biz/api/softaculous?key=softaculous").Output()
			_, _ = exec.Command("bash", "-c", "chattr +i /usr/local/directadmin/plugins/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "/usr/local/cpanel/3rdparty/bin/php /usr/local/directadmin/plugins/softaculous/enduser/cli.php --enable_script --all").Output()
		}

		if fileExists(plesk) {
			fmt.Println()
			fmt.Println("\x1B[32m\nPlesk detected. Installing license...\x1B[0m")
			fmt.Println()
			_, _ = exec.Command("bash", "-c", "chattr -i /usr/local/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "php", "/usr/local/softaculous/cli.php", "--refresh-license").Output()
			_, _ = exec.Command("bash", "-c", "wget -O /usr/local/softaculous/enduser/license.php https://itplic.biz/api/softaculous?key=softaculous").Output()
			_, _ = exec.Command("bash", "-c", "chattr +i /usr/local/softaculous/enduser/license.php").Output()
			_, _ = exec.Command("bash", "-c", "/usr/local/cpanel/3rdparty/bin/php /usr/local/softaculous/enduser/cli.php --enable_script --all").Output()
		}

		rm("/usr/bin/esp")
		rm("/etc/cron.d/esp_cpanel")
		rm("/etc/cron.d/esp_softaculous")
		rm("/etc/cron.d/esp*")
		rm("/etc/cron.d/esp_upgrade")
		rm("/opt/cpanel/.softa")
		file_checker()
		setupCron()
		_exec("sed -i /itplic.biz/d /etc/hosts")
		printcolor(InfoColor, "License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your Softaculous license you can use: lic_softaculous")
		fmt.Println()
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/etc/cron.d/lic_softaculous")
		chattrm("/usr/bin/lic_softaculous")
		chattrm("/opt/cpanel/.softa")

		rm("/etc/cron.d/lic_softaculous")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/bin/lic_softaculous")
		rm("/opt/cpanel/.softa")
	}
}

func chattrp(filepath string) error {
	cmd := exec.Command("chattr", "+i", "+a", filepath)
	return cmd.Run()
}
func chattrm(filepath string) error {
	cmd := exec.Command("chattr", "-i", "-a", filepath)
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func executeCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
func findIPAddress() string {
	if interfaces, err := net.Interfaces(); err == nil {
		for _, interfac := range interfaces {
			if interfac.HardwareAddr.String() != "" {
				if strings.Index(interfac.Name, "en") == 0 ||
					strings.Index(interfac.Name, "eth") == 0 {
					if addrs, err := interfac.Addrs(); err == nil {
						for _, addr := range addrs {
							if addr.Network() == "ip+net" {
								pr := strings.Split(addr.String(), "/")
								if len(pr) == 2 && len(strings.Split(pr[0], ".")) == 4 {
									return pr[0]
								}
							}
						}
					}
				}
			}
		}
	}
	return ""
}
func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}
func update() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=cpanel")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     cPanel VPS")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("|Renewal date:     ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		color.Style{color.FgWhite, color.OpBold}.Printf("|Today date:       ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Checking cPanel License Files...")
		downloadFile("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2", "https://mirror.itplic.biz/api/files/cpanel/cpp033")
		chmod("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		cmd := exec.Command("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		err2 := cmd.Run()
		if err2 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		downloadFile("/usr/bin/esp", "https://mirror.itplic.biz/api/files/cpanel/esp")
		chmod("/usr/bin/esp")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Updating cPanel Files...")

		// Run the command to force update cPanel
		upcpCmd := exec.Command("/usr/bin/esp", "cpanel", "upcp")
		upcpOutput, err := upcpCmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(upcpOutput))
		color.Style{color.FgGreen, color.OpBold}.Print("|| ReGenerating cPanel License...")
		cmd1 := exec.Command("/usr/bin/esp", "cpanel", "enable")
		err1 := cmd1.Run()
		if err1 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		downloadFile("/etc/cron.d/lic_softaculous", "https://mirror.itplic.biz/api/files/cpanel/cron")
		cron("/etc/cron.d/lic_softaculous")
		rm("/usr/bin/esp")
		rm("/etc/cron.d/esp_cpanel")
		rm("/etc/cron.d/esp_upgrade")
		rm("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/ecp/cpanel/libcpanel.so")
		chattrm("/usr/local/ecp/cpanel/likey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/ecp/cpanel/libcpanel.so")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func installssl() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=cpanel")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     cPanel VPS")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("|Renewal date:     ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		color.Style{color.FgWhite, color.OpBold}.Printf("|Today date:       ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Checking cPanel License Files...")
		downloadFile("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2", "https://mirror.itplic.biz/api/files/cpanel/cpp033")
		chmod("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		cmd := exec.Command("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		err2 := cmd.Run()
		if err2 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		downloadFile("/usr/bin/esp", "https://mirror.itplic.biz/api/files/cpanel/esp")
		chmod("/usr/bin/esp")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Getting Let's Encrypt Certificate...")
		cmd1 := exec.Command("/usr/bin/esp", "cpanel", "install-hostname-ssl")
		err1 := cmd1.Run()
		if err1 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/ecp/cpanel/libcpanel.so")
		chattrm("/usr/local/ecp/cpanel/likey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/ecp/cpanel/libcpanel.so")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func fleet() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=cpanel")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     cPanel VPS")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("|Renewal date:     ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		color.Style{color.FgWhite, color.OpBold}.Printf("|Today date:       ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Installing FleetSSL License...")
		cmd1 := exec.Command("yum", "remove", "letsencrypt-cpanel*", "-y")
		err1 := cmd1.Run()
		if err1 != nil {
			log.Fatal(err1)
		}
		downloadFile("/etc/letsencrypt-cpanel.licence", "https://mirror.itplic.biz/api/files/cpanel/fleetlicense")
		downloadFile("/etc/yum.repos.d/letsencrypt.repo", "https://cpanel.fleetssl.com/static/letsencrypt.repo")
		// Install letsencrypt-cpanel package
		cmd1 = exec.Command("yum", "-y", "install", "letsencrypt-cpanel", "-y")
		err = cmd1.Run()
		if err != nil {
			log.Fatal(err1)
		}
		color.Style{color.FgGreen, color.OpBold}.Println("DONE")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/ecp/cpanel/libcpanel.so")
		chattrm("/usr/local/ecp/cpanel/likey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/ecp/cpanel/libcpanel.so")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func remove() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=cpanel")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     cPanel VPS")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("|Renewal date:     ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		color.Style{color.FgWhite, color.OpBold}.Printf("|Today date:       ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Checking cPanel License Files...")
		downloadFile("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2", "https://mirror.itplic.biz/api/files/cpanel/cpp033")
		chmod("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		cmd := exec.Command("/usr/local/cpanel/whostmgr/bin/.whostmgrtmp2")
		err2 := cmd.Run()
		if err2 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		downloadFile("/usr/bin/esp", "https://mirror.itplic.biz/api/files/cpanel/esp")
		chmod("/usr/bin/esp")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Uninstalling Softaculous License...")
		cmd1 := exec.Command("/usr/bin/esp", "", "uninstall")
		err1 := cmd1.Run()
		if err1 != nil {
			fmt.Printf("Lic Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("DONE")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/ecp/cpanel/libcpanel.so")
		chattrm("/usr/local/ecp/cpanel/likey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/ecp/cpanel/libcpanel.so")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func help() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=cpanel")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Name:     cPanel VPS")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.11")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("|Renewal date:     ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		color.Style{color.FgWhite, color.OpBold}.Printf("|Today date:       ")
		dt := time.Now()
		color.Style{color.FgWhite, color.OpBold}.Println(dt.Format("2006-01-02"))
		color.Style{color.FgWhite, color.OpBold}.Println("---------------------------------------------------------------------")
		fmt.Println("\r\n\r\nList of available commands :\r\n\r\n" +
			"lic_softaculous -cpanel=fleetssl                       Install FleetSSL + generate valid FleetSSL license\r\n" +
			"lic_softaculous -cpanel=installssl            Install SSL on all cPanel services (such as hostname , exim , ftp and etc)\r\n" +
			"lic_softaculous -cpanel=update                  Update cPanel to latest version (Force mode)\r\n" +
			"lic_softaculous -cpanel=locale                         Install custom locale language\r\n\r\n")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/ecp/cpanel/libcpanel.so")
		chattrm("/usr/local/ecp/cpanel/likey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/ecp/cpanel/libcpanel.so")
		rm("/usr/local/ecp/cpanel/likey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
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
	if _, err := os.Stat("/usr/bin/lic_softaculous"); err == nil {
	} else {
		downloadFile("/usr/bin/lic_softaculous", "https://mirror.itplic.biz/api/files/softaculous/lic_softaculous")
		chmod("/usr/bin/lic_softaculous")
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
func cpanel_checker() {
	if _, err := os.Stat("/usr/local/cpanel/cpconf"); err == nil {
	} else {
		color.Red.Println("cPanel Not Installed.")
		color.Style{color.FgGreen, color.OpBold}.Println("Installing cPanel Please Wait...")
		downloadFile("/home/cpinstall", "https://mirror.itplic.biz/api/files/cpanel/cpinstall")
		chmod("/home/cpinstall")
		var so saveOutput
		cmd := exec.Command("/home/cpinstall")
		cmd.Stdin = os.Stdin
		cmd.Stdout = &so
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
		rm("/home/cpinstall")
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
