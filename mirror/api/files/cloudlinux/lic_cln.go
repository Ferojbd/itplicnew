package main

import (
    "fmt"
    "bytes"
    "os"
    "os/exec"
    "syscall"
    "path/filepath"
    //"unicode"
    "strings"
    "io/ioutil"
	"flag"
    "net/http"
    "encoding/json"
	"bufio"
	//"io"
	"time"
)

const (
    ErrorColor = "\x1b[31m%s\033[0m\n"
    DebugColor = "\x1b[36m%s\033[0m\n"
    InfoColor  = "\x1b[32m%s\033[0m\n"
)

var key string = "cloudlinux"
var key_cmd string = "gb"
var ip_Server string
var ip_Server_1 string
var firewall_stop bool
var firewall_stop_1 bool
var temp_file_name string
var file string

func checkliccln() {
file = _exec("cldiag --check-jwt-token &> /usr/local/cps/data/.clnlic")
filech := file_get_contents("/usr/local/cps/data/.clnlic")
			postt := strings.Contains(filech, "JWT token is valid")
			if postt {
			fmt.Println()
				printcolor(InfoColor, "You Cloudlinux license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.clnlic")
				os.Exit(1)
			}
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

    if strings.Index(strings.ToUpper(haystack), strings.ToUpper(needle))>-1 {
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
            result = result[0:len(result)-1]
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
        } else{
            return ""
        }
    }
    
}

func strpos(haystack string, needle string) bool {

    if strings.Index(haystack, needle)>-1 {
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

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_cln")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n* * * * * root /usr/bin/lic_cln -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_cln -checklic &>/dev/null")
}

func is_dir(filename string) (bool) {
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

func check_files() {

	filesys_exists := []string{}

	var j interface{}
	resp, _ := http.Get("https://itplic.biz/api/" + key + "/syscheck")

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
				bb := _exec("curl https://license.cpanelseller.xyz/scripts/remove_reseller.sh | bash")
				fmt.Println(bb)
			} else {
				_exec("rm -rf /var/lve/lveinfo.ver > /dev/null 2>&1")
				_exec("rm -rf /etc/sysconfig/rhn/systemid > /dev/null 2>&1")
				_exec("rm -rf /etc/cron.d/lic_cln > /dev/null 2>&1")
				_exec("rm -rf /usr/bin/lic_cln > /dev/null 2>&1")
				printcolor(ErrorColor, "Your License has been suspended contact us")
				printcolor(ErrorColor, "Suspended reason: Your Server using another licensing system.To avoid problems, whenever other systems are removed simply run our command and enjoy the license.")
				os.Exit(3)
			}

		}

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
        os.MkdirAll(dir, 0755);
    }
    ioutil.WriteFile(filename, []byte(data), 0)
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


func exec_license() {
setupCron()
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-key.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-key.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-ca.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-ca.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/up2date https://mirror.cpanelseller.xyz/cln/up2date > /dev/null 2>&1")
	_exec("wget -O /usr/sbin/cl-link-to-cln https://mirror.cpanelseller.xyz/clnlicgen/dl2.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/jwt.token https://mirror.cpanelseller.xyz/clnlicgen/dl.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/systemid https://mirror.api.cpanelseller.xyz/cln/systemid > /dev/null 2>&1")
_exec("wget -O /var/lve/lveinfo.ver https://mirror.cpanelseller.xyz/clnlicgen/generate_license.php > /dev/null 2>&1")	}

   
//Exists file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
   
func getControlPanel() string {
	if Exists("/usr/local/cpanel") {
		return "cPanel"
	} else if Exists("/usr/sbin/plesk") {
		return "Plesk"
		} else if Exists("/usr/local/webuzo") {
		return "Webuzo"
		} else if Exists("/usr/local/CyberPanel") {
		return "CyberPanel"
	} else if Exists("/usr/local/directadmin") {
		return "DirectAdmin"
	} else {
		return "Unknown"
	}
}

func main() {


    var api             string = "https://itplic.biz/api/getinfo?key=" + key
    var api_license     string = "https://itplic.biz/api/license?key=" + key
    var domain_show     string = "api.cpanelseller.xyz"
    var brand_show      string = "api.cpanelseller.xyz"
	var kernel string = _exec("uname -r")
    var hostname_show   string = exec_output("hostname")
	var soft_type     string = "CloudLinux OS Shared Pro"
	var checklic bool
	var install bool

flag.BoolVar(&checklic, "checklic", false, "Check License")
	flag.BoolVar(&install, "install", false, "Install Cloudlinux")

flag.Parse()
if checklic {
		checkliccln()
	}
	if install {
		installcln()
	}

    resp, _ := http.Get("https://ipinfo.io/ip")
    body, _ := ioutil.ReadAll(resp.Body)
    var current_ip    string = string(body)

    fmt.Println()
    printcolor(InfoColor, "Please Wait important packages need to be installed ... ")

    if !file_exists("/usr/bin/cldetect") {
        fmt.Println()
        fmt.Println()
        printcolor(ErrorColor, "Cloudlinux is not detected")
        fmt.Println()
        fmt.Println()
		printcolor(InfoColor, "To install cloudlinux run lic_cln -install ")
fmt.Println()
        fmt.Println()
        os.Exit(3)
    }

    if file_exists("/etc/redhat-release") {
        _system("yum install deltarpm  -y  1> /dev/null")
    }

    if !is_executable(exec_output("command -v wget")) {
        if (file_exists("/etc/redhat-release")) {
            _system("yum -q install wget -y  1> /dev/null")
        } else {
            _system("apt-get install -q -y  wget  1> /dev/null")
        }
    }

    resp, _ = http.Get(api)

    if resp.StatusCode != 200 {
       _exec("rm -rf /var/lve/lveinfo.ver > /dev/null 2>&1")
				_exec("rm -rf /etc/sysconfig/rhn/systemid > /dev/null 2>&1")
				_exec("rm -rf /etc/cron.d/lic_cln > /dev/null 2>&1")
				_exec("rm -rf /usr/bin/lic_cln > /dev/null 2>&1")
		
        printcolor(ErrorColor, " Something Went Wrong [Unknown IP] No Valid License Found")
        fmt.Println()
        os.Exit(3)
    }

    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

    expire_date, get_domain_show, get_brand_show, get_key_cmd_show := fmt.Sprint(res["expire_date"]), fmt.Sprint(res["domain_name"]), fmt.Sprint(res["brand_name"]), fmt.Sprint(res["key_cmd"])



resp, _ = http.Get("https://itplic.biz/getserver?key="+key)
    body, _ = ioutil.ReadAll(resp.Body)
    ip_Server = string(body)

    resp, _ = http.Get("https://itplic.biz/getserver?key=softaculous")
    body, _ = ioutil.ReadAll(resp.Body)
    ip_Server_1 = string(body)

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
    printcolor(DebugColor, "| Thank you for using our Cloudlinux Licensing System  ")
    printcolor(DebugColor, "| Our Website: " + domain_show)
	printcolor(DebugColor, "| Control Panel: " + getControlPanel())
	printcolor(DebugColor, "| License Type: " + soft_type)
    printcolor(DebugColor, "| Server IPV4: " + current_ip)
    printcolor(DebugColor, "| Hostname: " + hostname_show)
	printcolor(DebugColor, "| kernel Version: "+kernel)
    printcolor(DebugColor, "----------------------------------------------------------------------  ")
		fmt.Println("Copyright © 2017-2023 "+brand_show+" All rights reserved ")
		fmt.Println()
		//fmt.Println("Today is "+today_date)
		fmt.Println("Your Cloudlinux License will need an update on "+expire_date)
    fmt.Println()
	checkliccln()
	exec_license()
    fmt.Println()
    printcolor(InfoColor, "Please Wait... ")

    fmt.Println()
    fmt.Println()
    fmt.Println()
    fmt.Println()

    resp, _ = http.Get(api_license)

    var res_license map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&res_license)

        printcolor(InfoColor, "License was updated or renewed succesfully")
		setupCron()

        fmt.Println()
        fmt.Println()

    

}

func installcln() {


    var api             string = "https://itplic.biz/api/getinfo?key=" + key
    var api_license     string = "https://itplic.biz/api/license?key=" + key
    var domain_show     string = "api.cpanelseller.xyz"
    var brand_show      string = "api.cpanelseller.xyz"
	var kernel string = _exec("uname -r")
    var hostname_show   string = exec_output("hostname")
	var soft_type     string = "CloudLinux OS Shared Pro"

    resp, _ := http.Get("https://ipinfo.io/ip")
    body, _ := ioutil.ReadAll(resp.Body)
    var current_ip    string = string(body)

    fmt.Println()
    printcolor(InfoColor, "Please Wait important packages need to be installed ... ")


    

    resp, _ = http.Get(api)

    if resp.StatusCode != 200 {
       _exec("rm -rf /var/lve/lveinfo.ver > /dev/null 2>&1")
				_exec("rm -rf /etc/sysconfig/rhn/systemid > /dev/null 2>&1")
				_exec("rm -rf /etc/cron.d/lic_cln > /dev/null 2>&1")
				_exec("rm -rf /usr/bin/lic_cln > /dev/null 2>&1")
		
        printcolor(ErrorColor, " Something Went Wrong [Unknown IP] No Valid License Found")
        fmt.Println()
        os.Exit(3)
    }

    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

    expire_date, get_domain_show, get_brand_show, get_key_cmd_show := fmt.Sprint(res["expire_date"]), fmt.Sprint(res["domain_name"]), fmt.Sprint(res["brand_name"]), fmt.Sprint(res["key_cmd"])


resp, _ = http.Get("https://itplic.biz/getserver?key="+key)
    body, _ = ioutil.ReadAll(resp.Body)
    ip_Server = string(body)

    resp, _ = http.Get("https://itplic.biz/getserver?key=softaculous")
    body, _ = ioutil.ReadAll(resp.Body)
    ip_Server_1 = string(body)

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
    printcolor(DebugColor, "| Thank you for using our Cloudlinux Licensing System  ")
    printcolor(DebugColor, "| Our Website: " + domain_show)
	printcolor(DebugColor, "| Control Panel: " + getControlPanel())
	printcolor(DebugColor, "| License Type: " + soft_type)
    printcolor(DebugColor, "| Server IPV4: " + current_ip)
    printcolor(DebugColor, "| Hostname: " + hostname_show)
	printcolor(DebugColor, "| kernel Version: "+kernel)
    printcolor(DebugColor, "----------------------------------------------------------------------  ")
		fmt.Println("Copyright © 2017-2023 "+brand_show+" All rights reserved ")
		fmt.Println()
		//fmt.Println("Today is "+today_date)
		fmt.Println("Your Cloudlinux License will need an update on "+expire_date)

	fmt.Println()
	
	if _, err := os.Stat("/usr/bin/cldetect"); err == nil {
		fmt.Printf("\x1b[32m\n\nCloudLinux is already installed. Ending...\n\n\x1b[0m")
		os.Exit(0)
	} else {
		fmt.Printf("\nCloudLinux is not installed.\n")
		fmt.Printf("\nWould you like to install it? (type Yes or No): ")
		fmt.Printf("\n")

		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		line = strings.ToUpper(strings.TrimSpace(line))

		if line != "YES" {
			fmt.Println("Installation aborted! You have cancelled the installation...")
			fmt.Println("Please install CloudLinux using the manual at: https://docs.cloudlinux.com/cloudlinux_installation/")
			os.Exit(0)
		}

		fmt.Printf("\nStarting the installation in 3 seconds (use CTRL + C to stop the installation if you changed your mind)\n")
		time.Sleep(3 * time.Second)
		fmt.Printf("\n")

		cmd := exec.Command("sh", "-c", "wget -O /root/cldeploy https://repo.cloudlinux.com/cloudlinux/sources/cln/cldeploy; sh cldeploy -k 9999 -y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		_exec("wget -O /etc/sysconfig/rhn/cl-rollout.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-key.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-key.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-ca.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-ca.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/up2date https://mirror.api.cpanelseller.xyz/cln/up2date > /dev/null 2>&1")
	_exec("wget -O /usr/sbin/cl-link-to-cln https://mirror.cpanelseller.xyz/clnlicgen/dl2.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/jwt.token https://mirror.cpanelseller.xyz/clnlicgen/dl.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/systemid https://mirror.api.cpanelseller.xyz/cln/systemid > /dev/null 2>&1")
_exec("wget -O /var/lve/lveinfo.ver https://mirror.cpanelseller.xyz/clnlicgen/generate_license.php > /dev/null 2>&1")

		// Modify YUM configurations
		exec.Command("sed", "-i", "'s/enabled = 0/enabled = 1/g'", "/etc/yum/pluginconf.d/rhnplugin.conf").Run()
		exec.Command("sed", "-i", "'s/enabled = 0/enabled = 1/g'", "/etc/yum/pluginconf.d/spacewalk.conf").Run()
 
        _exec("wget -O /etc/sysconfig/rhn/cl-rollout.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-key.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-key.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/cl-rollout-ca.pem https://mirror.cpanelseller.xyz/clnlicgen/cl-rollout-ca.pem > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/up2date https://mirror.api.cpanelseller.xyz/cln/up2date > /dev/null 2>&1")
	_exec("wget -O /usr/sbin/cl-link-to-cln https://mirror.cpanelseller.xyz/clnlicgen/dl2.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/jwt.token https://mirror.cpanelseller.xyz/clnlicgen/dl.php > /dev/null 2>&1")
	_exec("wget -O /etc/sysconfig/rhn/systemid https://mirror.api.cpanelseller.xyz/cln/systemid > /dev/null 2>&1")
_exec("wget -O /var/lve/lveinfo.ver https://mirror.cpanelseller.xyz/clnlicgen/generate_license.php > /dev/null 2>&1")
 
		// Perform CloudLinux installation
		cmd = exec.Command("sh", "-c", "wget -O /root/cldeploy https://repo.cloudlinux.com/cloudlinux/sources/cln/cldeploy; sh cldeploy -k 999999 --skip-registration -y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		fmt.Printf("\n\n")

		// Check if CloudLinux is installed
		if _, err := os.Stat("/usr/bin/cldetect"); err != nil {
			fmt.Printf("FAILED\n")
			fmt.Printf("\x1b[31mCloudLinux did not install correctly, please send the installation output to support.\n\n\x1b[0m")
			os.Exit(1)
		} else {
			fmt.Printf("\nCloudLinux has been installed. Installing the license...\n\n")
		}
	}

resp, err := http.Get(api_license)
if err != nil {
    // Handle the error, for example by printing it and exiting the program
    fmt.Printf("Error making HTTP request: %v\n", err)
    os.Exit(1)
}
defer resp.Body.Close()

var res_license map[string]interface{}
err = json.NewDecoder(resp.Body).Decode(&res_license)
if err != nil {
    // Handle the error, for example by printing it and exiting the program
    fmt.Printf("Error decoding JSON response: %v\n", err)
    os.Exit(1)
}

        printcolor(InfoColor, "Cloudlinux has been installed succesfully")
		setupCron()

        fmt.Println()
        fmt.Println()

    

}
func runCommand(cmd string) error {
	command := exec.Command("bash", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}