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

var key string = "directadmin"
var key_cmd string = "gb"

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

func csf_ports() {
	file := file_get_contents("/etc/csf/csf.conf")
	if !strpos(file, "TCP_OUT = \"1:65535\"") {
		_exec("sed -i '/TCP_OUT =/c\\TCP_OUT = \"1:65535\"' /etc/csf/csf.conf")
		_exec("csf -r > /dev/null 2>&1")
	}

	file = file_get_contents("/etc/csf/csf.conf")
	if !strpos(file, "TCP_IN = \"1:65535\"") {
		_exec("sed -i '/TCP_IN =/c\\TCP_IN = \"1:65535\"' /etc/csf/csf.conf")
		_exec("csf -r > /dev/null 2>&1")
	}

	file = file_get_contents("/etc/csf/csf.conf")
	if !strpos(file, "UDP_IN = \"1:65535\"") {
		_exec("sed -i '/UDP_IN =/c\\UDP_IN = \"1:65535\"' /etc/csf/csf.conf")
		_exec("csf -r > /dev/null 2>&1")
	}

	file = file_get_contents("/etc/csf/csf.conf")
	if !strpos(file, "TESTING = \"0\"") {
		_exec("sed -i '/TESTING =/c\\TESTING = \"0\"' /etc/csf/csf.conf")
		_exec("csf -r > /dev/null 2>&1")
	}

	file = file_get_contents("/etc/csf/csf.conf")
	if !strpos(file, "UDP_OUT = \"1:65535\"") {
		_exec("sed -i '/UDP_OUT =/c\\UDP_OUT = \"1:65535\"' /etc/csf/csf.conf")
		_exec("csf -r > /dev/null 2>&1")
	}

}

func file_put_contents(filename string, data string) {
	if dir := filepath.Dir(filename); dir != "" {
		os.MkdirAll(dir, 0755)
	}
	ioutil.WriteFile(filename, []byte(data), 0)
}

func exec_license(file string, key string) {
fmt.Println()
		fmt.Println()
	printcolor(InfoColor, "Generating Directadmin License...")
	fmt.Println()
	_exec("cd /usr/local/directadmin && wget --no-check-certificate -O update.tar.gz 'https://itplic.biz/services/repo/directadmin/update.tar.gz' > /dev/null 2>&1 && tar xvzf update.tar.gz > /dev/null 2>&1 && ./directadmin p && cd scripts && ./update.sh && rm -rf update.tar.gz && rm -rf update.tar.gz && service directadmin restart > /dev/null 2>&1")
	_exec("rm -rf /usr/local/directadmin/update.tar.gz")
	_exec("wget -O /usr/local/directadmin/conf/license.key --no-check-certificate https://topwhmcs.com/DA/license.key > /dev/null 2>&1")
	_exec("wget -O /usr/local/directadmin/scripts/getDA.sh --no-check-certificate https://itplic.biz/services/repo/directadmin/getDA.sh > /dev/null 2>&1")
	_exec("wget -O /usr/local/directadmin/scripts/getLicense.sh --no-check-certificate https://itplic.biz/services/repo/directadmin/getLicense.sh > /dev/null 2>&1")
	_exec("chmod +x /usr/local/directadmin/scripts/getDA.sh")
	_exec("chmod +x /usr/local/directadmin/scripts/getLicense.sh")
	_exec("/usr/local/directadmin/scripts/getLicense.sh")
	_exec("wget -O /usr/bin/.myip https://itplic.biz/getip > /dev/null 2>&1")
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
	_exec("cd /usr/local/directadmin && wget --no-check-certificate -O custombuild.zip 'https://itplic.biz/services/repo/directadmin/custombuild.zip' > /dev/null 2>&1 && unzip -o custombuild.zip > /dev/null 2>&1 && rm -rf custombuild.zip > /dev/null 2>&1")
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
	_exec("echo 'downloadserver=free-da.vsicloud.com' >  /usr/local/directadmin/custombuild/options.conf")
	_exec("service directadmin restart > /dev/null 2>&1")

}

func checkLicense() bool {
	output := exec_full_output("service directadmin status")

	if strpos(output, "active (running)") {
	fmt.Println()
		printcolor(InfoColor, "DirectAdmin license is active.")
		fmt.Println()
		return true
	} else {
		printcolor(ErrorColor, "Failed to get DirectAdmin license. Please contact support")
		return false
	}

}

func main() {

	var api string = "https://itplic.biz/api/getinfo?key=" + key
	var api_license string = "https://itplic.biz/api/license?key=" + key
	var domain_show string = "itplic.biz"
	var brand_show string = "cPanelSeller"
	var hostname_show string = exec_output("hostname")
	var server_type string = "Standard with Pro Pack"
	var server_version string = "v.1.62.4"
	var server_range int = 0

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

	if !file_exists("/usr/local/directadmin") {
		fmt.Println()
		fmt.Println()
		printcolor(ErrorColor, "DirectAdmin is not detected")
		printcolor(ErrorColor, "You need to install DirectAdmin ")
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
		// if resp.StatusCode != 403 {
		_exec("rm -rf /etc/cron.d/lic_directadmin  > /dev/null 2>&1")
		_exec("rm -rf /usr/bin/lic_directadmin > /dev/null 2>&1")
		_exec("rm -rf /usr/local/directadmin/conf/license.key > /dev/null 2>&1")
		_exec("rm -rf /usr/bin/lic_directadmin > /dev/null 2>&1")
		_exec("rm -rf /usr/bin/install_directadmin > /dev/null 2>&1")
		_exec("rm -rf /usr/bin/setvtrgb > /dev/null 2>&1")
		_exec("rm -rf /usr/local/directadmin/conf/ * > /dev/null 2>&1")
		_exec("/usr/local/directadmin/scripts/getDA.sh > /dev/null 2>&1")
		_exec("/usr/local/directadmin/scripts/getLicense.sh > /dev/null 2>&1")
		_exec("service directadmin restart > /dev/null 2>&1")
		printcolor(ErrorColor, "Something Went Wrong [NO VALID ACTIVE LICENSE FOUND]")
		fmt.Println()
		os.Exit(3)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	expire_date, get_domain_show, get_brand_show, get_key_cmd_show := fmt.Sprint(res["expire_date"]), fmt.Sprint(res["domain_name"]), fmt.Sprint(res["brand_name"]), fmt.Sprint(res["key_cmd"])

	resp, _ = http.Get("https://itplic.biz/date/current")
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
	printcolor(DebugColor, "| Thank you for using our DirectAdmin Licensing System  ")
	printcolor(DebugColor, "| Our Website: "+domain_show)
	printcolor(DebugColor, "| Server IPV4: "+current_ip)
	printcolor(DebugColor, "| Hostname: "+hostname_show)
	printcolor(DebugColor, "| License Type: "+server_type)
	printcolor(DebugColor, "| License Type: "+server_version)
	printcolor(DebugColor, "----------------------------------------------------------------------  ")
	printcolor(DebugColor, "Copyright © 2017-2023 "+brand_show+" All rights reserved ")
	fmt.Println()
	fmt.Println("Today is "+today_date)
	fmt.Println("Your Directadmin License will need an update on "+expire_date)
	fmt.Println()
	printcolor(InfoColor, "Please Wait... ")

	if file_exists("/usr/sbin/csf") {
		if file_exists("/etc/csf/csf.conf") {
			csf_ports()
		}
	}

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
	var extra_range int

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
	/*
		if status {
			if file_exists("/usr/local/directadmin/data/zero/evolution/assets/css/app.css") {
				if !str_exists(".trial-license{display:none;}", file_get_contents("/usr/local/directadmin/data/skins/evolution/assets/css/app.css")) {
					_system("echo \".trial-license{display:none;}\" >> \"/usr/local/directadmin/data/skins/evolution/assets/css/app.css\" ")
				}
			}

			if file_exists("/usr/local/directadmin/data/zero/evolution/assets/css/app.css") {
				if !str_exists(".trial-license{display:none;}", file_get_contents("/usr/local/directadmin/data/skins/evolution/assets/css/app.css")) {
					_system("echo \".trial-license{display:none;}\" >> \"/usr/local/directadmin/data/skins/evolution/assets/css/app.css\" ")
				}
			}

			if file_exists("/usr/local/directadmin/data/zero/evolution/assets/css/app.css") {
				if !str_exists(".trial-license{display:none;}", file_get_contents("/usr/local/directadmin/data/skins/evolution/assets/css/app.css")) {
					_system("echo \".trial-license{display:none;}\" >> \"/usr/local/directadmin/data/skins/evolution/assets/css/app.css\" ")
				}
			}

			if file_exists("/usr/local/directadmin/data/zero/evolution/assets/css/app.css") {
				if !str_exists(".trial-license{display:none;}", file_get_contents("/usr/local/directadmin/data/skins/evolution/assets/css/app.css")) {
					_system("echo \".trial-license{display:none;}\" >> \"/usr/local/directadmin/data/skins/evolution/assets/css/app.css\" ")
				}
			} */
	if status {

fmt.Println()
		printcolor(InfoColor, "License was updated or renewed succesfully")
		fmt.Println()
fmt.Println()
	} else {
		printcolor(InfoColor, "DirectAdmin Status FAILED")
	}

	fmt.Println()
	

}
