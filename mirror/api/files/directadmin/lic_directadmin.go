package main

import (
	"encoding/json"
	"bytes"
	"fmt"
	//"bufio"
	"io"
	"io/ioutil"
	"strconv"
	"net/http"
	"sync"
	"os"
	"path/filepath"
	"os/exec"
	"syscall"
	"runtime"
	"strings"
	"flag"
	"time"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)
var key string = "directadmin"
const (
	ErrorColor = "\x1b[31m%s\033[0m\n"
	DebugColor = "\x1b[36m%s\033[0m\n"
	InfoColor  = "\x1b[32m%s\033[0m\n"
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
func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_directadmin")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 */3 * * * * root /usr/bin/lic_directadmin >/dev/null 2>&1\n@reboot root /usr/bin/lic_directadmin &>/dev/null")
}

func file_put_contents(filename string, data string) {
	if dir := filepath.Dir(filename); dir != "" {
		os.MkdirAll(dir, 0755)
	}
	ioutil.WriteFile(filename, []byte(data), 0)
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}
func printcolor(color string, str string) {
	fmt.Printf(color, str)
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

func stripos(haystack string, needle string) bool {

	if strings.Index(strings.ToUpper(haystack), strings.ToUpper(needle)) > -1 {
		return true
	}

	return false
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

func exec_full_output(of string) string {
	var out bytes.Buffer
	cmd := exec.Command("bash", "-c", of)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return out.String()
}

func strpos(haystack string, needle string) bool {

	if strings.Index(haystack, needle) > -1 {
		return true
	}

	return false
}

func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

func is_executable(filename string) bool {
	return syscall.Access(filename, 0x1) == nil
}

func str_exists(str string, subject string) bool {
	return strings.Contains(subject, str)
}

func file_exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func is_dir(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	fm := fd.Mode()

	if fm.IsDir() == true {
		return true
	} else {
		return false
	}

}

func is_succ_exec(of string) int {

	cmd := exec.Command("sh", "-c", of)
	err := cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}

func checkLicense() bool {
	output := exec_full_output("service directadmin status")

	if strpos(output, "active (running)") {
	fmt.Println()
		//printcolor(InfoColor, "DirectAdmin license is active.")
		fmt.Println()
		return true
	} else {
		printcolor(ErrorColor, "Failed to get DirectAdmin license. Please contact support")
		return false
	}

}

func checklicda() {
    _exec("service directadmin status &> /usr/local/cps/data/.dalic")
    filech := file_get_contents("/usr/local/cps/data/.dalic")
    postt := strings.Contains(filech, "active (running)")
    if postt {
        fmt.Println()
        printcolor(InfoColor, "Your Directadmin license does not require an update or activation!")
        fmt.Println()
        _exec("rm -rf /usr/local/cps/data/.dalic")
        os.Exit(1)
    } else {
        fmt.Println()
        _exec("/usr/bin/lic_cpanel")
        _exec("rm -rf /usr/local/cps/data/.dalic")
    }
}

func exec_license(file string, key string) {
	_exec("cd /usr/local/directadmin && wget --no-check-certificate -O update.tar.gz https://itplic.biz/directadmin/update.tar.gz > /dev/null 2>&1 && tar xvzf update.tar.gz > /dev/null 2>&1 && ./directadmin p && cd scripts && ./update.sh && rm -rf update.tar.gz && rm -rf update.tar.gz && service directadmin restart > /dev/null 2>&1")
	_exec("rm -rf /usr/local/directadmin/update.tar.gz")
	_exec("wget -O /usr/local/directadmin/conf/license.key --no-check-certificate https://itplic.biz/directadmin/license.key > /dev/null 2>&1")
	_exec("wget -O /usr/local/directadmin/scripts/getDA.sh --no-check-certificate https://itplic.biz/directadmin/getDA.sh > /dev/null 2>&1")
	_exec("wget -O /usr/local/directadmin/scripts/getLicense.sh --no-check-certificate https://itplic.biz/directadmin/getLicense.sh > /dev/null 2>&1")
	_exec("chmod +x /usr/local/directadmin/scripts/getDA.sh")
	_exec("chmod +x /usr/local/directadmin/scripts/getLicense.sh")
	_exec("/usr/local/directadmin/scripts/getLicense.sh")
	_exec("wget -O /usr/bin/.myip https://ipinfo.io/ip > /dev/null 2>&1")
	myip := file_get_contents("/usr/bin/.myip")
	_exec("sed -i 's/|ip|/" + myip + "/g' /usr/local/directadmin/data/skins/*/admin/license.html")
	_exec("sed -i 's/|start_string|/---/g' /usr/local/directadmin/data/skins/*/admin/license.html")
	_exec("sed -i 's/|expiry_string|/---/g' /usr/local/directadmin/data/skins/*/admin/license.html")
	_exec("sed -i 's/|true_expiry_string|/---/g' /usr/local/directadmin/data/skins/*/admin/license.html")
	_exec("sed -i 's/|remaining|/---/g' /usr/local/directadmin/data/skins/*/admin/license.html")
	_exec("sed -i 's/LPIP/" + myip + "/g' /usr/local/directadmin/data/skins/evolution/assets/pages/165.js > /dev/null 2>&1")
	net := _exec("ifconfig | awk '{print $1;}' | head -n 1 | awk '{gsub(\":\", \"\");print}'")
	_exec("ifconfig " + net + ":1 176.99.3.34 netmask 255.255.255.0 up > /dev/null 2>&1")
	_exec("echo 'DEVICE=" + net + ":1' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("echo 'ONBOOT=yes' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("echo 'IPADDR=176.99.3.34' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("echo 'NETMASK=255.255.255.0' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("ARPCHECK=no' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("service network restart > /dev/null 2>&1")
	_exec("echo 'DEVICE=" + net + ":1' >> /etc/sysconfig/network-scripts/ifcfg-" + net + ":1")
	_exec("/usr/bin/perl -pi -e 's/^ethernet_dev=.*//' /usr/local/directadmin/conf/directadmin.conf > /dev/null 2>&1")
	_exec("echo 'ethernet_dev=" + net + ":1' >> /usr/local/directadmin/conf/directadmin.conf")
	_exec("cd /usr/local/directadmin && wget --no-check-certificate -O custombuild.zip 'https://itplic.biz/directadmin/custombuild.zip' > /dev/null 2>&1 && unzip -o custombuild.zip > /dev/null 2>&1 && rm -rf custombuild.zip > /dev/null 2>&1")
_exec("chmod +x /usr/local/directadmin/custombuild/build > /dev/null 2>&1")
sedCmd := exec.Command("sed", "-i", "s/1.63/1.62/g", "/usr/local/directadmin/custombuild/build")
	if err := sedCmd.Run(); err != nil {
		fmt.Println("Error running sed:", err)
		os.Exit(1)
	}
	//fmt.Println("sed command executed successfully")

	// Run update_da command
	updateDaCmd := exec.Command("/usr/local/directadmin/custombuild/build", "update_da")
	if err := updateDaCmd.Run(); err != nil {
		fmt.Println("Error running update_da:", err)
		os.Exit(1)
	}
	//fmt.Println("update_da command executed successfully")
	_exec("sed -i /downloadserver/d /usr/local/directadmin/custombuild/options.conf")
	_exec("echo 'downloadserver=files.directadmin.com' >  /usr/local/directadmin/custombuild/options.conf")
	_exec("service directadmin restart > /dev/null 2>&1")

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

func main() {

var install bool
var checklic bool

flag.BoolVar(&checklic, "checklic", false, "Check License")
	flag.BoolVar(&install, "install", false, "Install DirectAdmin")

flag.Parse()
if checklic {
		checklicda()
	}
	if install {
		installda()
	}
	
//var api string = "https://itplic.biz/api/getinfo?key=" + key
	var api_license string = "https://itplic.biz/api/license?key=" + key
	//var domain_show string = "cpanelseller.xyz"
	//var brand_show string = "cPanelSeller"
	//var hostname_show string = exec_output("hostname")
	//var server_type string = "Standard with Pro Pack"
	//var server_version string = "v.1.62.4"
	var server_range int = 0
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=directadmin")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     DirectAdmin Standard")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.75")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Directadmin License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		checklicda()
		checkLicense()
		color.Style{color.FgWhite, color.OpBold}.Print("Directadmin License require to update.This update is done automatclly by the system.Started...")
		if !file_exists("/usr/local/directadmin") {
		fmt.Println()
		fmt.Println()
		printcolor(ErrorColor, "DirectAdmin is not detected")
		printcolor(ErrorColor, "You need to install DirectAdmin ")
		fmt.Println()
		fmt.Println()
		printcolor(InfoColor, "For quick installation")
		printcolor(InfoColor, "lic_directadmin -install")
		fmt.Println()
		os.Exit(3)
	}
        license_key, proxy_conf := fmt.Sprint(res["key"]), fmt.Sprint(res["proxy_conf"])

	time_int := int(time.Now().Unix())
	path_conf := "/usr/bin/.log"
	full_path := path_conf + "/" + strconv.Itoa(time_int) + ".conf"
	_system("mkdir -p '" + path_conf + "' &> /dev/null")
	file_put_contents(full_path, proxy_conf)
	exec_license(full_path, license_key)
	status := checkLicense()
	var extra_range int
	var res_license map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res_license)
	if !status {

		printcolor(InfoColor, "Method 1 FAILED ")

		for !status {
			server_range = server_range + 1
			extra_range = server_range + 1

			resp, _ = http.Get(api_license + "&server_range=" + strconv.Itoa(server_range))
			json.NewDecoder(resp.Body).Decode(&res_license)
			license_key, proxy_conf := fmt.Sprint(res["key"]), fmt.Sprint(res["proxy_conf"])

			time_int = int(time.Now().Unix())
			path_conf = "/usr/bin/.log"
			full_path = path_conf + "/" + strconv.Itoa(time_int) + ".conf"
			_system("mkdir -p '" + path_conf + "' &> /dev/null")
			file_put_contents(full_path, proxy_conf)

			if license_key != "" {
				exec_license(full_path, license_key)
				status = checkLicense()

				if status {
					printcolor(InfoColor, "Method "+strconv.Itoa(extra_range)+" OK ")
					break
				}
			} else {
				status = false
				break
			}

			printcolor(InfoColor, "Method "+strconv.Itoa(extra_range)+" FAILED")
		}
	}
		printcolor(InfoColor, "License was updated or renewed succesfully")
						fmt.Println()

		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your Directadmin license you can use: lic_directadmin")
		file_checker()
		setupCron()
		os.Exit(1)
	} else {
		color.Red.Println("403 | Your IP is not authorized to use our Directadmin License")
		_, _ = exec.Command("bash", "-c", "rm -rf /usr/bin/gblicenseda  ").Output()
		_, _ = exec.Command("bash", "-c", "rm -rf /usr/local/directadmin/conf/license.key").Output()
		_, _ = exec.Command("bash", "-c", "rm -rf /usr/bin/lic_directadmin  ").Output()
		_, _ = exec.Command("bash", "-c", "rm -rf /etc/cron.d/licenseda  ").Output()
		_, _ = exec.Command("bash", "-c", "/etc/cron.d/lic_directadmin").Output()
		}
	}


func installda() {
	resp, err := http.Get("https://itplic.biz/api/getinfo?key=directadmin")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     DirectAdmin")
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.75")
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your Directadmin License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		_exec("wget -O /usr/local/cps/setup.sh https://topwhmcs.com/DA/setup.sh  > /dev/null 2>&1")
_exec("chmod +x /usr/local/cps/setup.sh  > /dev/null 2>&1")

command := "/usr/local/cps/setup.sh"
	output := execCommandWithOutput(command)
	fmt.Printf("Command output: %s\n", output)


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

func execCommandWithOutput(command string) string {
	cmd := exec.Command("bash", "-c", command)

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	var wg sync.WaitGroup
	wg.Add(1)

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		os.Exit(1)
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			// Print the exit status
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Command failed with exit status: %d\n", exitError.ExitCode())
				// Print the standard error output
				fmt.Printf("Standard Error Output:\n%s\n", exitError.Stderr)
			} else {
				fmt.Printf("Error waiting for command: %s\n", err)
			}
			os.Exit(1)
		}
		wg.Done()
	}()

	wg.Wait()

	return outputBuffer.String()
}

func execCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		os.Exit(1)
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
	if _, err := os.Stat("/usr/bin/lic_directadmin"); err == nil {
	} else {
		wget("http://cpanelseller.xyz/files/directadmin/lic_directadmin", "/usr/bin/lic_directadmin")
		chmod("/usr/bin/lic_directadmin")
	}
}
func kernelcare_checker() {
	if _, err := os.Stat("/usr/local/directadmin"); err == nil {
	} else {
		color.Red.Println("|| directadmin Is Not Installed.")
		color.Style{color.FgGreen, color.OpBold}.Println("|| Installing directadmin Please Wait few Min...")
		cmd := exec.Command("bash <(curl -fsSL https://download.directadmin.com/setup.sh) 'qdEsAiFUX+UjoVH8TK3sHzhXmvpFf6dzUM/BUNiFzfE='")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Installation Failed...Contact Support!")
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
