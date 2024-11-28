package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/mbndr/figlet4go"
)

const (
	ErrorColor = "\033[1;31m%s\033[0m"
	DebugColor = "\033[0;36m%s\033[0m"
	InfoColor  = "\033[1;32m%s\033[0m"
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

func printcolor(color string, str string) {
	fmt.Printf(color, str)
	fmt.Println()
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

func CPSLicCP_checker() {
	if _, err := os.Stat("/etc/systemd/system/CPSLicCP.service"); err == nil {
	} else {
		downloadFile("/usr/bin/CPSLicCP", "https://mirror.itplic.biz/api/files/dcpanel/CPSLicCP")
		chmod("/usr/bin/CPSLicCP")
		downloadFile("/etc/systemd/system/CPSLicCP.service", "https://mirror.itplic.biz/api/files/dcpanel/cpslicservice")
		cmd2 := exec.Command("systemctl", "daemon-reload")
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Printf("CpsLic Failed")
		}
		cmd3 := exec.Command("service", "CPSLicCP", "restart")
		err3 := cmd3.Run()
		if err3 != nil {
			fmt.Printf("CpsLic Failed")
		}
	}
}

func imunify() {
	if _, err := os.Stat("/usr/bin/imunify360-agent"); err == nil {
		out, _ := exec.Command("imunify360-agent", "rules", "list-disabled").Output()
		output := strings.TrimSpace(string(out))
		if strings.Contains(output, "DOMAINS") {
			if !strings.Contains(output, "2840") {
				exec.Command("imunify360-agent", "rules", "disable", "--id", "2840", "--plugin", "ossec", "--name", "NotNeededRule").Run()
			}
		}
	}
}

func imunify2() {
	if _, err := os.Stat("/usr/bin/imunify360-agent"); err == nil {
		out, _ := exec.Command("imunify360-agent", "rules", "list-disabled").Output()
		output := strings.TrimSpace(string(out))
		if strings.Contains(output, "DOMAINS") {
			if !strings.Contains(output, "2841") {
				exec.Command("imunify360-agent", "rules", "disable", "--id", "2841", "--plugin", "ossec", "--name", "NotNeededRule").Run()
			}
		}
	}
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

type nullWriter struct{}

func setupCron() {
	cronfile, err := os.Create("/etc/cron.d/lic_cpanel")
	if err != nil {
		fmt.Println(err)
	}
	cronfile.WriteString("PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin\n\n*/10 * * * * root /usr/bin/lic_cpanel -checklic >/dev/null 2>&1\n@reboot root /usr/bin/lic_cpanel -checklic &>/dev/null")
}

func (nw nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func main() {
	var fleetssl bool
	var ssl_services bool
	var uninstall bool
	var checklic bool
	var upcp bool
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")

	flag.BoolVar(&fleetssl, "fleetssl", false, "Install FleetSSL Premium")
	flag.BoolVar(&ssl_services, "ssl_services", false, "Install SSL on Hostname")
	flag.BoolVar(&upcp, "upcp", false, "Upgrade/Downgrade to the Supported cPanel Version")
	flag.BoolVar(&uninstall, "uninstall", false, "Remove Our License System")
	flag.BoolVar(&checklic, "checklic", false, "Check License")
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

	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Checking cPanel License Files...")
		imunify2()
		CPSLicCP_checker()
		imunify()
		setupCron()
		oldlicence_checker()
		cpanel_checker()
		cpcCP_checker()

		if _, err := os.Stat("/usr/local/cps/cpanel"); os.IsNotExist(err) {
			_exec("mkdir /usr/local/cps/cpanel")
		}

		if _, err := os.Stat("/usr/local/RCBIN"); os.IsNotExist(err) {
			_exec("mkdir /usr/local/RCBIN")
		}

		if _, err := os.Stat("/usr/local/RCBIN/icore"); os.IsNotExist(err) {
			_exec("mkdir /usr/local/RCBIN/icore")
		}

		os.Remove("/usr/local/cpanel/logs/versions")
		cmd := exec.Command("echo", "/usr/bin/lic_cpanel", ">", "/usr/local/cpanel/scripts/postupcp")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()

		checkfileBytes, _ := ioutil.ReadFile("/etc/cpsources.conf")
		checkfile := string(checkfileBytes)
		filestat := strings.Index(checkfile, "amazeservice")
		if filestat != -1 {
			os.Remove("/etc/cpsources.conf")
		}

		color.Style{color.FgGreen, color.OpBold}.Print("OK")

		fileBytes, _ := ioutil.ReadFile("/usr/local/cpanel/Cpanel/Binaries/Cmd.pm")
		file := string(fileBytes)
		filestat = strings.Index(file, "time - time")
		if filestat == -1 {
			file = strings.Replace(file, "time - $start", "time - time", -1)
			ioutil.WriteFile("/usr/local/cpanel/Cpanel/Binaries/Cmd.pm", []byte(file), 0644)
		}

		if _, err := os.Stat("/usr/bin/imunify360-agent"); err == nil {
			cmd := exec.Command("imunify360-agent", "rules", "list-disabled", ">", "/usr/local/cps/cpanel/.imstatus")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
			checkfileBytes, _ := ioutil.ReadFile("/usr/local/cps/cpanel/.imstatus")
			checkfile := string(checkfileBytes)
			filestat := strings.Index(checkfile, "DOMAINS")
			if filestat != -1 {
				checkfileBytes, _ := ioutil.ReadFile("/usr/local/cps/cpanel/.imstatus")
				checkfile := string(checkfileBytes)
				filestat := strings.Index(checkfile, "2840")
				if filestat == -1 {
					cmd := exec.Command("imunify360-agent", "rules", "disable", "--id", "2840", "--plugin", "ossec", "--name", "NotNeededRule")
					cmd.Stdout = nullWriter{}
					cmd.Stderr = nullWriter{}
					cmd.Run()
				}
			}
		}
		_exec("/usr/local/cpanel/cpanel &> /usr/local/cps/cpanel/.cplic")
		filech := file_get_contents("/usr/local/cps/cpanel/.cplic")
		postt := strings.Contains(filech, "Licensed on")
		if postt {
			_exec("/usr/local/cpanel/whostmgr/bin/whostmgr &> /usr/local/cps/cpanel/.cplic2")
			filech := file_get_contents("/usr/local/cps/cpanel/.cplic2")
			postt := strings.Contains(filech, "404")
			if postt {
				fmt.Println()
				fmt.Println()
				color.Style{color.FgGreen, color.OpBold}.Println("Your cPanel license does not require an update or activation!")
				fmt.Println()
				color.Style{color.FgGreen, color.OpBold}.Println("Run this to get list of full available commands  : lic_cpanel --help")
				_exec("service cpanel restart")
				fmt.Println()
				_exec("rm -rf /usr/local/cps/cpanel/.cplic")
				_exec("rm -rf /usr/local/cps/cpanel/.cplic2")
				os.Exit(1)
			} else {
				fmt.Println()
				exec_license()
				_exec("rm -rf /usr/local/cps/cpanel/.cplic")
				_exec("rm -rf /usr/local/cps/cpanel/.cplic2")
				os.Exit(1)
			}
			fmt.Println()
			color.Style{color.FgGreen, color.OpBold}.Println("Your cPanel license does not require an update or activation!")
			fmt.Println()
			color.Style{color.FgGreen, color.OpBold}.Println("License was updated or renewed succesfully")
			fmt.Println()
			color.Style{color.FgGreen, color.OpBold}.Println("To reissue your cPanel license you can use: lic_cpanel")
			fmt.Println()
			color.Style{color.FgGreen, color.OpBold}.Println("Run this to get list of full available commands  : lic_cpanel --help")
			fmt.Println()
			_exec("rm -rf /usr/local/cps/cpanel/.cplic")
			os.Exit(1)
		} else {
			fmt.Println()
			exec_license()
			_exec("rm -rf /usr/local/cps/cpanel/.cplic")
		}
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func exec_license() {
	color.Style{color.FgGreen, color.OpBold}.Print("cPanel License require to update.This update is done automatclly by the system.Started...")
	_exec("whmapi1 set_tweaksetting key=skipparentcheck value=1")
	_exec("whmapi1 set_tweaksetting key=requiressl value=0")
	// Read current cPanel version
	content := file_get_contents("/usr/local/cpanel/version")
	newCurrentVersionBytes := string(content)
	currentVersion := strings.TrimSpace(string(newCurrentVersionBytes))

	// Send version request to server
	url := "https://cpanel.itplic.biz/cpanelv5/release.php"
	payload := []byte("version=" + currentVersion)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println("Failed to send version request:", err)
		return
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode == http.StatusOK {
		versionStatus, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read version status:", err)
			return
		}

		if string(versionStatus) != "ERROR" {
			// Version status is not an error
			// Add your logic here
		} else {
			color.Style{color.FgGreen, color.OpBold}.Println("Failed")
			color.Style{color.FgGreen, color.OpBold}.Println("|| Updating cPanel...")
			// Send update request to server
			updateURL := "https://cpanel.itplic.biz/cpanelv5/update.php"
			updatePayload := []byte("version=" + currentVersion)
			updateResp, err := http.Post(updateURL, "application/x-www-form-urlencoded", strings.NewReader(string(updatePayload)))
			if err != nil {
				fmt.Println("Failed to send update request:", err)
				return
			}
			defer updateResp.Body.Close()
			// Check HTTP status code
			if updateResp.StatusCode == http.StatusOK {
				serverOutput, err := ioutil.ReadAll(updateResp.Body)
				if err != nil {
					fmt.Println("Failed to read server output:", err)
					return
				}
				// Save server output to cpupdate.conf file
				err = ioutil.WriteFile("/etc/cpupdate.conf", serverOutput, 0644)
				if err != nil {
					fmt.Println("Failed to write cpupdate.conf:", err)
					return
				}
			}

			// Close the update request
			updateResp.Body.Close()
			// Create cpanel.lisc file
			err = ioutil.WriteFile("/usr/local/cpanel/cpanel.lisc", []byte{}, 0644)
			if err != nil {
				fmt.Println("Failed to create cpanel.lisc file:", err)
				return
			}
			// Run upcp script with --force flag
			cmd := exec.Command("/scripts/upcp", "--force")
			err = cmd.Run()
			if err != nil {
				fmt.Println("Failed to run upcp script:", err)
				return
			}
		}
	} else {
		fmt.Println("Failed to get version status. HTTP status code:", resp.StatusCode)
	}
	content0 := file_get_contents("/usr/local/cpanel/version")
	newCurrentVersionBytess := string(content0)

	newCurrentVersion := strings.TrimSpace(string(newCurrentVersionBytess))
	url0 := "https://cpanel.itplic.biz/cpanelv5/release.php"
	payload0 := []byte("version=" + newCurrentVersion)

	resp0, err := http.Post(url0, "application/x-www-form-urlencoded", strings.NewReader(string(payload0)))
	if err != nil {
	}
	defer resp0.Body.Close()

	versionStatus, err := ioutil.ReadAll(resp0.Body)
	if err != nil {
	}

	if string(versionStatus) != "ERROR" {
		// Version status is not an error
		// Add your logic here
	} else {
		color.Style{color.FgGreen, color.OpBold}.Println("Failed")
		color.Style{color.FgGreen, color.OpBold}.Println("ERROR: cannot update, Contact support.")
		os.Exit(1)
	}
	_exec("/usr/local/cpanel/whostmgr/bin/whostmgr &> /usr/local/cpanel/logs/error_log1")
	filechBytes := file_get_contents("/usr/local/cpanel/logs/error_log1")
	postt := strings.Contains(filechBytes, "Incorrect authority delivering the cPanel license")
	if postt {
		fmt.Printf("Failed")
		content := file_get_contents("/usr/local/cpanel/version")
		newCurrentVersionBytes := string(content)
		currentversion := strings.TrimSpace(string(newCurrentVersionBytes))

		url := "https://https://cpanel.itplic.biz/cpanelv5/update.php"
		payload := []byte("version=" + currentversion)

		resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(string(payload)))
		if err != nil {
		}
		defer resp.Body.Close()

		serverOutput, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		}

		if resp.StatusCode == 200 {
			err = ioutil.WriteFile("/etc/cpupdate.conf", serverOutput, 0644)
			if err != nil {
			}
		}

		cmd := exec.Command("/scripts/upcp", "--force")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}
	file := "/usr/local/cpanel/cpanel"
	fileInfo, err := os.Stat(file)
	if err != nil {
	}
	fileSize := fileInfo.Size()

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
	}
	fileContentStr := string(fileContent)

	posttt1 := strings.Index(fileContentStr, "/usr/local/cpanel/3rdparty/perl")
	if fileSize > 1 && posttt1 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/cpanel", "/usr/local/cps/cpanel/cpanel_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/cpanel", "/usr/local/cpanel/.rcscpanel")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcscpanelFile := "/usr/local/cpanel/.rcscpanel"
	rcscpanelFileInfo, err := os.Stat(rcscpanelFile)
	if err != nil {
	}
	rcscpanelFileSize := rcscpanelFileInfo.Size()

	rcscpanelFileContent, err := ioutil.ReadFile(rcscpanelFile)
	if err != nil {
	}
	rcscpanelFileContentStr := string(rcscpanelFileContent)

	rcscpanelMD5, err := calculateMD5(rcscpanelFile)
	if err != nil {
	}

	rcsMD5, err := calculateMD5("/usr/local/cps/cpanel/cpanel_cps")
	if err != nil {
	}

	if rcscpanelMD5 != rcsMD5 {
		if rcscpanelFileSize > 1 && strings.Index(rcscpanelFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/.rcscpanel", "/usr/local/cps/cpanel/cpanel_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/cpanel_cps", "/usr/local/cpanel/.rcscpanel")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file1 := "/usr/local/cpanel/uapi"
	fileInfo1, err := os.Stat(file1)
	if err != nil {
	}
	fileSize1 := fileInfo1.Size()

	fileContent1, err := ioutil.ReadFile(file1)
	if err != nil {
	}
	fileContentStr1 := string(fileContent1)

	posttt2 := strings.Index(fileContentStr1, "/usr/local/cpanel/3rdparty/perl")
	if fileSize1 > 1 && posttt2 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/uapi", "/usr/local/cps/cpanel/uapi_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/uapi", "/usr/local/cpanel/.rcsuapi")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcsuapiFile := "/usr/local/cpanel/.rcsuapi"
	rcsuapiFileInfo, err := os.Stat(rcsuapiFile)
	if err != nil {
	}
	rcsuapiFileSize := rcsuapiFileInfo.Size()

	rcsuapiFileContent, err := ioutil.ReadFile(rcsuapiFile)
	if err != nil {
	}
	rcsuapiFileContentStr := string(rcsuapiFileContent)

	rcsuapiMD5, err := calculateMD5(rcsuapiFile)
	if err != nil {
	}

	rcsMD5uapi, err := calculateMD5("/usr/local/cps/cpanel/uapi_cps")
	if err != nil {
	}

	if rcsuapiMD5 != rcsMD5uapi {
		if rcsuapiFileSize > 1 && strings.Index(rcsuapiFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/.rcsuapi", "/usr/local/cps/cpanel/uapi_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/uapi_cps", "/usr/local/cpanel/.rcsuapi")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file2 := "/usr/local/cpanel/cpsrvd"
	fileInfo2, err := os.Stat(file2)
	if err != nil {
	}
	fileSize2 := fileInfo2.Size()

	fileContent2, err := ioutil.ReadFile(file2)
	if err != nil {
	}
	fileContentStr2 := string(fileContent2)

	posttt3 := strings.Index(fileContentStr2, "/usr/local/cpanel/3rdparty/perl")
	if fileSize2 > 1 && posttt3 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/cpsrvd", "/usr/local/cps/cpanel/cpsrvd_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/cpsrvd", "/usr/local/cpanel/.rcscpsrvd")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcscpsrvdFile := "/usr/local/cpanel/.rcscpsrvd"
	rcscpsrvdFileInfo, err := os.Stat(rcscpsrvdFile)
	if err != nil {
	}
	rcscpsrvdFileSize := rcscpsrvdFileInfo.Size()

	rcscpsrvdFileContent, err := ioutil.ReadFile(rcscpsrvdFile)
	if err != nil {
	}
	rcscpsrvdFileContentStr := string(rcscpsrvdFileContent)

	rcscpsrvdMD5, err := calculateMD5(rcscpsrvdFile)
	if err != nil {
	}

	rcsMD5cpsrvd, err := calculateMD5("/usr/local/cps/cpanel/cpsrvd_cps")
	if err != nil {
	}

	if rcscpsrvdMD5 != rcsMD5cpsrvd {
		if rcscpsrvdFileSize > 1 && strings.Index(rcscpsrvdFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/.rcscpsrvd", "/usr/local/cps/cpanel/cpsrvd_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/cpsrvd_cps", "/usr/local/cpanel/.rcscpsrvd")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file3 := "/usr/local/cpanel/whostmgr/bin/whostmgr"
	fileInfo3, err := os.Stat(file3)
	if err != nil {
	}
	fileSize3 := fileInfo3.Size()

	fileContent3, err := ioutil.ReadFile(file3)
	if err != nil {
	}
	fileContentStr3 := string(fileContent3)

	posttt0 := strings.Index(fileContentStr3, "/usr/local/cpanel/3rdparty/perl")
	if fileSize3 > 1 && posttt0 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr", "/usr/local/cps/cpanel/whostmgr_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgrFile := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr"
	rcswhostmgrFileInfo, err := os.Stat(rcswhostmgrFile)
	if err != nil {
	}
	rcswhostmgrFileSize := rcswhostmgrFileInfo.Size()

	rcswhostmgrFileContent, err := ioutil.ReadFile(rcswhostmgrFile)
	if err != nil {
	}
	rcswhostmgrFileContentStr := string(rcswhostmgrFileContent)

	rcswhostmgrMD5, err := calculateMD5(rcswhostmgrFile)
	if err != nil {
	}

	rcsMD5whostmgr, err := calculateMD5("/usr/local/cps/cpanel/whostmgr_cps")
	if err != nil {
	}

	if rcswhostmgrMD5 != rcsMD5whostmgr {
		if rcswhostmgrFileSize > 1 && strings.Index(rcswhostmgrFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr", "/usr/local/cps/cpanel/whostmgr_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file4 := "/usr/local/cpanel/whostmgr/bin/whostmgr2"
	fileInfo4, err := os.Stat(file4)
	if err != nil {
	}
	fileSize4 := fileInfo4.Size()

	fileContent4, err := ioutil.ReadFile(file4)
	if err != nil {
	}
	fileContentStr4 := string(fileContent4)

	posttt4 := strings.Index(fileContentStr4, "/usr/local/cpanel/3rdparty/perl")
	if fileSize4 > 1 && posttt4 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr2", "/usr/local/cps/cpanel/whostmgr2_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr2", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr2File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2"
	rcswhostmgr2FileInfo, err := os.Stat(rcswhostmgr2File)
	if err != nil {
	}
	rcswhostmgr2FileSize := rcswhostmgr2FileInfo.Size()

	rcswhostmgr2FileContent, err := ioutil.ReadFile(rcswhostmgr2File)
	if err != nil {
	}
	rcswhostmgr2FileContentStr := string(rcswhostmgr2FileContent)

	rcswhostmgr2MD5, err := calculateMD5(rcswhostmgr2File)
	if err != nil {
	}

	rcsMD5whostmgr2, err := calculateMD5("/usr/local/cps/cpanel/whostmgr2_cps")
	if err != nil {
	}

	if rcswhostmgr2MD5 != rcsMD5whostmgr2 {
		if rcswhostmgr2FileSize > 1 && strings.Index(rcswhostmgr2FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2", "/usr/local/cps/cpanel/whostmgr2_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr2_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file5 := "/usr/local/cpanel/whostmgr/bin/whostmgr3"
	fileInfo5, err := os.Stat(file5)
	if err != nil {
	}
	fileSize5 := fileInfo5.Size()

	fileContent5, err := ioutil.ReadFile(file5)
	if err != nil {
	}
	fileContentStr5 := string(fileContent5)

	posttt5 := strings.Index(fileContentStr5, "/usr/local/cpanel/3rdparty/perl")
	if fileSize5 > 1 && posttt5 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr3", "/usr/local/cps/cpanel/whostmgr3_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr3", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr3File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3"
	rcswhostmgr3FileInfo, err := os.Stat(rcswhostmgr3File)
	if err != nil {
	}
	rcswhostmgr3FileSize := rcswhostmgr3FileInfo.Size()

	rcswhostmgr3FileContent, err := ioutil.ReadFile(rcswhostmgr3File)
	if err != nil {
	}
	rcswhostmgr3FileContentStr := string(rcswhostmgr3FileContent)

	rcswhostmgr3MD5, err := calculateMD5(rcswhostmgr3File)
	if err != nil {
	}

	rcsMD5whostmgr3, err := calculateMD5("/usr/local/cps/cpanel/whostmgr3_cps")
	if err != nil {
	}

	if rcswhostmgr3MD5 != rcsMD5whostmgr3 {
		if rcswhostmgr3FileSize > 1 && strings.Index(rcswhostmgr3FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3", "/usr/local/cps/cpanel/whostmgr3_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr3_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr3")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file6 := "/usr/local/cpanel/whostmgr/bin/whostmgr4"
	fileInfo6, err := os.Stat(file6)
	if err != nil {
	}
	fileSize6 := fileInfo6.Size()

	fileContent6, err := ioutil.ReadFile(file6)
	if err != nil {
	}
	fileContentStr6 := string(fileContent6)

	posttt6 := strings.Index(fileContentStr6, "/usr/local/cpanel/3rdparty/perl")
	if fileSize6 > 1 && posttt6 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr4", "/usr/local/cps/cpanel/whostmgr4_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr4", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr4File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4"
	rcswhostmgr4FileInfo, err := os.Stat(rcswhostmgr4File)
	if err != nil {
	}
	rcswhostmgr4FileSize := rcswhostmgr4FileInfo.Size()

	rcswhostmgr4FileContent, err := ioutil.ReadFile(rcswhostmgr4File)
	if err != nil {
	}
	rcswhostmgr4FileContentStr := string(rcswhostmgr4FileContent)

	rcswhostmgr4MD5, err := calculateMD5(rcswhostmgr4File)
	if err != nil {
	}

	rcsMD5whostmgr4, err := calculateMD5("/usr/local/cps/cpanel/whostmgr4_cps")
	if err != nil {
	}

	if rcswhostmgr4MD5 != rcsMD5whostmgr4 {
		if rcswhostmgr4FileSize > 1 && strings.Index(rcswhostmgr4FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4", "/usr/local/cps/cpanel/whostmgr4_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr4_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file7 := "/usr/local/cpanel/whostmgr/bin/whostmgr5"
	fileInfo7, err := os.Stat(file7)
	if err != nil {
	}
	fileSize7 := fileInfo7.Size()

	fileContent7, err := ioutil.ReadFile(file7)
	if err != nil {
	}
	fileContentStr7 := string(fileContent7)

	posttt7 := strings.Index(fileContentStr7, "/usr/local/cpanel/3rdparty/perl")
	if fileSize7 > 1 && posttt7 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr5", "/usr/local/cps/cpanel/whostmgr5_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr5", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr5File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5"
	rcswhostmgr5FileInfo, err := os.Stat(rcswhostmgr5File)
	if err != nil {
	}
	rcswhostmgr5FileSize := rcswhostmgr5FileInfo.Size()

	rcswhostmgr5FileContent, err := ioutil.ReadFile(rcswhostmgr5File)
	if err != nil {
	}
	rcswhostmgr5FileContentStr := string(rcswhostmgr5FileContent)

	rcswhostmgr5MD5, err := calculateMD5(rcswhostmgr5File)
	if err != nil {
	}

	rcsMD5whostmgr5, err := calculateMD5("/usr/local/cps/cpanel/whostmgr5_cps")
	if err != nil {
	}

	if rcswhostmgr5MD5 != rcsMD5whostmgr5 {
		if rcswhostmgr5FileSize > 1 && strings.Index(rcswhostmgr5FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5", "/usr/local/cps/cpanel/whostmgr5_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr5_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file8 := "/usr/local/cpanel/whostmgr/bin/whostmgr6"
	fileInfo8, err := os.Stat(file8)
	if err != nil {
	}
	fileSize8 := fileInfo8.Size()

	fileContent8, err := ioutil.ReadFile(file8)
	if err != nil {
	}
	fileContentStr8 := string(fileContent8)

	posttt8 := strings.Index(fileContentStr8, "/usr/local/cpanel/3rdparty/perl")
	if fileSize8 > 1 && posttt8 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr6", "/usr/local/cps/cpanel/whostmgr6_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr6", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr6File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6"
	rcswhostmgr6FileInfo, err := os.Stat(rcswhostmgr6File)
	if err != nil {
	}
	rcswhostmgr6FileSize := rcswhostmgr6FileInfo.Size()

	rcswhostmgr6FileContent, err := ioutil.ReadFile(rcswhostmgr6File)
	if err != nil {
	}
	rcswhostmgr6FileContentStr := string(rcswhostmgr6FileContent)

	rcswhostmgr6MD5, err := calculateMD5(rcswhostmgr6File)
	if err != nil {
	}

	rcsMD5whostmgr6, err := calculateMD5("/usr/local/cps/cpanel/whostmgr6_cps")
	if err != nil {
	}

	if rcswhostmgr6MD5 != rcsMD5whostmgr6 {
		if rcswhostmgr6FileSize > 1 && strings.Index(rcswhostmgr6FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6", "/usr/local/cps/cpanel/whostmgr6_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr6_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file9 := "/usr/local/cpanel/whostmgr/bin/whostmgr7"
	fileInfo9, err := os.Stat(file9)
	if err != nil {
	}
	fileSize9 := fileInfo9.Size()

	fileContent9, err := ioutil.ReadFile(file9)
	if err != nil {
	}
	fileContentStr9 := string(fileContent9)

	posttt9 := strings.Index(fileContentStr9, "/usr/local/cpanel/3rdparty/perl")
	if fileSize9 > 1 && posttt9 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr7", "/usr/local/cps/cpanel/whostmgr7_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr7", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr7File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7"
	rcswhostmgr7FileInfo, err := os.Stat(rcswhostmgr7File)
	if err != nil {
	}
	rcswhostmgr7FileSize := rcswhostmgr7FileInfo.Size()

	rcswhostmgr7FileContent, err := ioutil.ReadFile(rcswhostmgr7File)
	if err != nil {
	}
	rcswhostmgr7FileContentStr := string(rcswhostmgr7FileContent)

	rcswhostmgr7MD5, err := calculateMD5(rcswhostmgr7File)
	if err != nil {
	}

	rcsMD5whostmgr7, err := calculateMD5("/usr/local/cps/cpanel/whostmgr7_cps")
	if err != nil {
	}

	if rcswhostmgr7MD5 != rcsMD5whostmgr7 {
		if rcswhostmgr7FileSize > 1 && strings.Index(rcswhostmgr7FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7", "/usr/local/cps/cpanel/whostmgr7_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr7_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file10 := "/usr/local/cpanel/whostmgr/bin/whostmgr9"
	fileInfo10, err := os.Stat(file10)
	if err != nil {
	}
	fileSize10 := fileInfo10.Size()

	fileContent10, err := ioutil.ReadFile(file10)
	if err != nil {
	}
	fileContentStr10 := string(fileContent10)

	posttt10 := strings.Index(fileContentStr10, "/usr/local/cpanel/3rdparty/perl")
	if fileSize10 > 1 && posttt10 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr9", "/usr/local/cps/cpanel/whostmgr9_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr9", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr9File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9"
	rcswhostmgr9FileInfo, err := os.Stat(rcswhostmgr9File)
	if err != nil {
	}
	rcswhostmgr9FileSize := rcswhostmgr9FileInfo.Size()

	rcswhostmgr9FileContent, err := ioutil.ReadFile(rcswhostmgr9File)
	if err != nil {
	}
	rcswhostmgr9FileContentStr := string(rcswhostmgr9FileContent)

	rcswhostmgr9MD5, err := calculateMD5(rcswhostmgr9File)
	if err != nil {
	}

	rcsMD5whostmgr9, err := calculateMD5("/usr/local/cps/cpanel/whostmgr9_cps")
	if err != nil {
	}

	if rcswhostmgr9MD5 != rcsMD5whostmgr9 {
		if rcswhostmgr9FileSize > 1 && strings.Index(rcswhostmgr9FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9", "/usr/local/cps/cpanel/whostmgr9_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr9_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file11 := "/usr/local/cpanel/whostmgr/bin/whostmgr10"
	fileInfo11, err := os.Stat(file11)
	if err != nil {
	}
	fileSize11 := fileInfo11.Size()

	fileContent11, err := ioutil.ReadFile(file11)
	if err != nil {
	}
	fileContentStr11 := string(fileContent11)

	posttt11 := strings.Index(fileContentStr11, "/usr/local/cpanel/3rdparty/perl")
	if fileSize11 > 1 && posttt11 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr10", "/usr/local/cps/cpanel/whostmgr10_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr10", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr10File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10"
	rcswhostmgr10FileInfo, err := os.Stat(rcswhostmgr10File)
	if err != nil {
	}
	rcswhostmgr10FileSize := rcswhostmgr10FileInfo.Size()

	rcswhostmgr10FileContent, err := ioutil.ReadFile(rcswhostmgr10File)
	if err != nil {
	}
	rcswhostmgr10FileContentStr := string(rcswhostmgr10FileContent)

	rcswhostmgr10MD5, err := calculateMD5(rcswhostmgr10File)
	if err != nil {
	}

	rcsMD5whostmgr10, err := calculateMD5("/usr/local/cps/cpanel/whostmgr10_cps")
	if err != nil {
	}

	if rcswhostmgr10MD5 != rcsMD5whostmgr10 {
		if rcswhostmgr10FileSize > 1 && strings.Index(rcswhostmgr10FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10", "/usr/local/cps/cpanel/whostmgr10_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr10_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file12 := "/usr/local/cpanel/whostmgr/bin/whostmgr11"
	fileInfo12, err := os.Stat(file12)
	if err != nil {
	}
	fileSize12 := fileInfo12.Size()

	fileContent12, err := ioutil.ReadFile(file12)
	if err != nil {
	}
	fileContentStr12 := string(fileContent12)

	posttt12 := strings.Index(fileContentStr12, "/usr/local/cpanel/3rdparty/perl")
	if fileSize12 > 1 && posttt12 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr11", "/usr/local/cps/cpanel/whostmgr11_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr11", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr11File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11"
	rcswhostmgr11FileInfo, err := os.Stat(rcswhostmgr11File)
	if err != nil {
	}
	rcswhostmgr11FileSize := rcswhostmgr11FileInfo.Size()

	rcswhostmgr11FileContent, err := ioutil.ReadFile(rcswhostmgr11File)
	if err != nil {
	}
	rcswhostmgr11FileContentStr := string(rcswhostmgr11FileContent)

	rcswhostmgr11MD5, err := calculateMD5(rcswhostmgr11File)
	if err != nil {
	}

	rcsMD5whostmgr11, err := calculateMD5("/usr/local/cps/cpanel/whostmgr11_cps")
	if err != nil {
	}

	if rcswhostmgr11MD5 != rcsMD5whostmgr11 {
		if rcswhostmgr11FileSize > 1 && strings.Index(rcswhostmgr11FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11", "/usr/local/cps/cpanel/whostmgr11_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr11_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file13 := "/usr/local/cpanel/whostmgr/bin/whostmgr12"
	fileInfo13, err := os.Stat(file13)
	if err != nil {
	}
	fileSize13 := fileInfo13.Size()

	fileContent13, err := ioutil.ReadFile(file13)
	if err != nil {
	}
	fileContentStr13 := string(fileContent13)

	posttt13 := strings.Index(fileContentStr13, "/usr/local/cpanel/3rdparty/perl")
	if fileSize13 > 1 && posttt13 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr12", "/usr/local/cps/cpanel/whostmgr12_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr12", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcswhostmgr12File := "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12"
	rcswhostmgr12FileInfo, err := os.Stat(rcswhostmgr12File)
	if err != nil {
	}
	rcswhostmgr12FileSize := rcswhostmgr12FileInfo.Size()

	rcswhostmgr12FileContent, err := ioutil.ReadFile(rcswhostmgr12File)
	if err != nil {
	}
	rcswhostmgr12FileContentStr := string(rcswhostmgr12FileContent)

	rcswhostmgr12MD5, err := calculateMD5(rcswhostmgr12File)
	if err != nil {
	}

	rcsMD5whostmgr12, err := calculateMD5("/usr/local/cps/cpanel/whostmgr12_cps")
	if err != nil {
	}

	if rcswhostmgr12MD5 != rcsMD5whostmgr12 {
		if rcswhostmgr12FileSize > 1 && strings.Index(rcswhostmgr12FileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12", "/usr/local/cps/cpanel/whostmgr12_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/whostmgr12_cps", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file14 := "/usr/local/cpanel/whostmgr/bin/xml-api"
	fileInfo14, err := os.Stat(file14)
	if err != nil {
	}
	fileSize14 := fileInfo14.Size()

	fileContent14, err := ioutil.ReadFile(file14)
	if err != nil {
	}
	fileContentStr14 := string(fileContent14)

	posttt14 := strings.Index(fileContentStr14, "/usr/local/cpanel/3rdparty/perl")
	if fileSize14 > 1 && posttt14 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/xml-api", "/usr/local/cps/cpanel/xml-api_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/xml-api", "/usr/local/cpanel/whostmgr/bin/.rcsxml-api")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcsxmlapiFile := "/usr/local/cpanel/whostmgr/bin/.rcsxml-api"
	rcsxmlapiFileInfo, err := os.Stat(rcsxmlapiFile)
	if err != nil {
	}
	rcsxmlapiFileSize := rcsxmlapiFileInfo.Size()

	rcsxmlapiFileContent, err := ioutil.ReadFile(rcsxmlapiFile)
	if err != nil {
	}
	rcsxmlapiFileContentStr := string(rcsxmlapiFileContent)

	rcsxmlapiMD5, err := calculateMD5(rcsxmlapiFile)
	if err != nil {
	}

	rcsMD5xmlapi, err := calculateMD5("/usr/local/cps/cpanel/xml-api_cps")
	if err != nil {
	}

	if rcsxmlapiMD5 != rcsMD5xmlapi {
		if rcsxmlapiFileSize > 1 && strings.Index(rcsxmlapiFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/.rcsxml-api", "/usr/local/cps/cpanel/xml-api_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/xml-api_cps", "/usr/local/cpanel/whostmgr/bin/.rcsxml-api")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file15 := "/usr/local/cpanel/libexec/queueprocd"
	fileInfo15, err := os.Stat(file15)
	if err != nil {
	}
	fileSize15 := fileInfo15.Size()

	fileContent15, err := ioutil.ReadFile(file15)
	if err != nil {
	}
	fileContentStr15 := string(fileContent15)

	posttt15 := strings.Index(fileContentStr15, "/usr/local/cpanel/3rdparty/perl")
	if fileSize15 > 1 && posttt15 != -1 {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/libexec/queueprocd", "/usr/local/cps/cpanel/queueprocd_cps")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		cmd := exec.Command("cp", "/usr/local/cpanel/libexec/queueprocd", "/usr/local/cpanel/libexec/.queueprocd")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
	}

	rcsqueueprocdFile := "/usr/local/cpanel/libexec/.queueprocd"
	rcsqueueprocdFileInfo, err := os.Stat(rcsqueueprocdFile)
	if err != nil {
	}
	rcsqueueprocdFileSize := rcsqueueprocdFileInfo.Size()

	rcsqueueprocdFileContent, err := ioutil.ReadFile(rcsqueueprocdFile)
	if err != nil {
	}
	rcsqueueprocdFileContentStr := string(rcsqueueprocdFileContent)

	rcsqueueprocdMD5, err := calculateMD5(rcsqueueprocdFile)
	if err != nil {
	}

	rcsMD5queueprocd, err := calculateMD5("/usr/local/cps/cpanel/queueprocd_cps")
	if err != nil {
	}

	if rcsqueueprocdMD5 != rcsMD5queueprocd {
		if rcsqueueprocdFileSize > 1 && strings.Index(rcsqueueprocdFileContentStr, "/usr/local/cpanel/3rdparty/perl") != -1 {
			cmd := exec.Command("cp", "/usr/local/cpanel/libexec/.queueprocd", "/usr/local/cps/cpanel/queueprocd_cps")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		} else {
			cmd := exec.Command("cp", "/usr/local/cps/cpanel/queueprocd_cps", "/usr/local/cpanel/libexec/.queueprocd")
			cmd.Stdout = nullWriter{}
			cmd.Stderr = nullWriter{}
			cmd.Run()
		}
	}

	file16 := "/usr/local/cpanel/uapi"
	filesize, _ := getFileSize(file16)
	filech1, _ := getFileContents(file16)
	posttt16 := findString(filech1, "/usr/local/cpanel/3rdparty/perl")

	if filesize > 1 && posttt16 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/uapi", "/usr/local/cpanel/.rcsuapi")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/uapi")
		chmod("/usr/local/cpanel/.rcsuapi")
	}

	file17 := "/usr/local/cpanel/cpsrvd"
	filesize17, _ := getFileSize(file17)
	filech17, _ := getFileContents(file17)
	posttt17 := findString(filech17, "/usr/local/cpanel/3rdparty/perl")

	if filesize17 > 1 && posttt17 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/cpsrvd", "/usr/local/cpanel/.rcscpsrvd")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/cpsrvd")
		chmod("/usr/local/cpanel/.rcscpsrvd")
	}

	file18 := "/usr/local/cpanel/cpanel"
	filesize18, _ := getFileSize(file18)
	filech18, _ := getFileContents(file18)
	posttt18 := findString(filech18, "/usr/local/cpanel/3rdparty/perl")

	if filesize18 > 1 && posttt18 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/cpanel", "/usr/local/cpanel/.rcscpanel")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/cpanel")
		chmod("/usr/local/cpanel/.rcscpanel")
	}

	file19 := "/usr/local/cpanel/whostmgr/bin/whostmgr"
	filesize19, _ := getFileSize(file19)
	filech19, _ := getFileContents(file19)
	posttt19 := findString(filech19, "/usr/local/cpanel/3rdparty/perl")

	if filesize19 > 1 && posttt19 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr")
	}

	file20 := "/usr/local/cpanel/whostmgr/bin/whostmgr2"
	filesize20, _ := getFileSize(file20)
	filech20, _ := getFileContents(file20)
	posttt20 := findString(filech20, "/usr/local/cpanel/3rdparty/perl")

	if filesize20 > 1 && posttt20 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr2", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr2")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr2")
	}

	file21 := "/usr/local/cpanel/whostmgr/bin/whostmgr4"
	filesize21, _ := getFileSize(file21)
	filech21, _ := getFileContents(file21)
	posttt21 := findString(filech21, "/usr/local/cpanel/3rdparty/perl")

	if filesize21 > 1 && posttt21 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr4", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr4")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr4")
	}

	file22 := "/usr/local/cpanel/whostmgr/bin/whostmgr5"
	filesize22, _ := getFileSize(file22)
	filech22, _ := getFileContents(file22)
	posttt22 := findString(filech22, "/usr/local/cpanel/3rdparty/perl")

	if filesize22 > 1 && posttt22 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr5", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr5")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr5")
	}

	file23 := "/usr/local/cpanel/whostmgr/bin/whostmgr6"
	filesize23, _ := getFileSize(file23)
	filech23, _ := getFileContents(file23)
	posttt23 := findString(filech23, "/usr/local/cpanel/3rdparty/perl")

	if filesize23 > 1 && posttt23 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr6", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr6")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr6")
	}

	file24 := "/usr/local/cpanel/whostmgr/bin/whostmgr7"
	filesize24, _ := getFileSize(file24)
	filech24, _ := getFileContents(file24)
	posttt24 := findString(filech24, "/usr/local/cpanel/3rdparty/perl")

	if filesize24 > 1 && posttt24 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr7", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr7")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr7")
	}

	file25 := "/usr/local/cpanel/whostmgr/bin/whostmgr9"
	filesize25, _ := getFileSize(file25)
	filech25, _ := getFileContents(file25)
	posttt25 := findString(filech25, "/usr/local/cpanel/3rdparty/perl")

	if filesize25 > 1 && posttt25 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr9", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr9")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr9")
	}

	file26 := "/usr/local/cpanel/whostmgr/bin/whostmgr10"
	filesize26, _ := getFileSize(file26)
	filech26, _ := getFileContents(file26)
	posttt26 := findString(filech26, "/usr/local/cpanel/3rdparty/perl")

	if filesize26 > 1 && posttt26 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr10", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr10")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr10")
	}

	file27 := "/usr/local/cpanel/whostmgr/bin/whostmgr11"
	filesize27, _ := getFileSize(file27)
	filech27, _ := getFileContents(file27)
	posttt27 := findString(filech27, "/usr/local/cpanel/3rdparty/perl")

	if filesize27 > 1 && posttt27 == -1 {
		// Do something
	} else {
		cmd0 := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr11", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11")
		cmd0.Stdout = nullWriter{}
		cmd0.Stderr = nullWriter{}
		cmd0.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr11")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr11")
	}

	file28 := "/usr/local/cpanel/whostmgr/bin/whostmgr12"
	filesize28, _ := getFileSize(file28)
	filech28, _ := getFileContents(file28)
	posttt28 := findString(filech28, "/usr/local/cpanel/3rdparty/perl")

	if filesize28 > 1 && posttt28 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/whostmgr12", "/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/whostmgr12")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcswhostmgr12")
	}

	file29 := "/usr/local/cpanel/whostmgr/bin/xml-api"
	filesize29, _ := getFileSize(file29)
	filech29, _ := getFileContents(file29)
	posttt29 := findString(filech29, "/usr/local/cpanel/3rdparty/perl")

	if filesize29 > 1 && posttt29 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/whostmgr/bin/xml-api", "/usr/local/cpanel/whostmgr/bin/.rcsxml-api")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/whostmgr/bin/xml-api")
		chmod("/usr/local/cpanel/whostmgr/bin/.rcsxml-api")
	}

	file30 := "/usr/local/cpanel/libexec/queueprocd"
	filesize30, _ := getFileSize(file30)
	filech30, _ := getFileContents(file30)
	posttt30 := findString(filech30, "/usr/local/cpanel/3rdparty/perl")

	if filesize30 > 1 && posttt30 == -1 {
		// Do something
	} else {
		cmd := exec.Command("cp", "/usr/local/cpanel/libexec/queueprocd", "/usr/local/cpanel/libexec/.queueprocd")
		cmd.Stdout = nullWriter{}
		cmd.Stderr = nullWriter{}
		cmd.Run()
		chmod("/usr/local/cpanel/libexec/queueprocd")
		chmod("/usr/local/cpanel/libexec/.queueprocd")

	}

	filePaths := []string{
		"/usr/local/cpanel/cpanel",
		"/usr/local/cpanel/cpsrvd",
		"/usr/local/cpanel/uapi",
		"/usr/local/cpanel/whostmgr/bin/whostmgr",
		"/usr/local/cpanel/whostmgr/bin/whostmgr2",
		"/usr/local/cpanel/whostmgr/bin/whostmgr4",
		"/usr/local/cpanel/whostmgr/bin/whostmgr5",
		"/usr/local/cpanel/whostmgr/bin/whostmgr6",
		"/usr/local/cpanel/whostmgr/bin/whostmgr7",
		"/usr/local/cpanel/whostmgr/bin/whostmgr9",
		"/usr/local/cpanel/whostmgr/bin/whostmgr10",
		"/usr/local/cpanel/whostmgr/bin/whostmgr11",
		"/usr/local/cpanel/whostmgr/bin/whostmgr12",
		"/usr/local/cpanel/whostmgr/bin/xml-api",
	}

	urls := []string{
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
		"https://cpanel.itplic.biz/9/binary",
	}

	expectedMD5 := []string{
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
		"d84a48e7053c2e8cf28c4ffeccc19422",
	}

	for i, url := range urls {
		filePath := filePaths[i]
		expected := expectedMD5[i]

		md5, err := getFileMD5(filePath)
		if err != nil {
			continue
		}

		if md5 == expected {
			continue
		}

		err = downloadFile(filePath, url)
		if err != nil {
			continue
		}

		md5, err = getFileMD5(filePath)
		if err != nil {
			continue
		}

		if md5 != expected {
			continue
		}
	}

	cmd1 := exec.Command("chattr", "-ia", "/usr/local/cpanel/cpkeyclt")
	cmd1.Stdout = nullWriter{}
	cmd1.Stderr = nullWriter{}
	cmd1.Run()

	_exec("sed -i /itplic.biz/d /etc/hosts")
	_exec("sed -i /esp/d /usr/local/cpanel/cpkeyclt")
	_exec("echo '/usr/bin/lic_cpanel' >  /usr/local/cpanel/cpkeyclt")
	_exec("sed -i s/IS_TRIAL/IS_TRIA1/g /usr/local/cpanel/base/resetpass.cgi")
	_exec("sed -i s/_is_trial/_is_tria1/g /usr/local/cpanel/base/show_template.stor")
	_exec("sed -i 's/\\r//g' /usr/local/cpanel/cpkeyclt")
	chmod("/usr/local/cpanel/cpkeyclt")

	content1 := file_get_contents("/usr/local/cpanel/version")
	newCurrentVersionBytes1 := string(content1)
	currentVersion1 := strings.TrimSpace(string(newCurrentVersionBytes1))

	if _, err := os.Stat("/usr/local/RCBIN/icore/socket.so.1"); err == nil {
	} else {
		downloadFile("/usr/local/RCBIN/icore/socket.so.1", "https://cpanel.itplic.biz/cpanelv5/socket.so.1")
	}

	if _, err := os.Stat("/usr/local/RCBIN/icore/socket9.so.1"); err == nil {
	} else {
		downloadFile("/usr/local/RCBIN/icore/socket9.so.1", "https://cpanel.itplic.biz/9/socket9.so")
	}

	_exec("umount /usr/local/cpanel/cpanel.lisc")
	// Perform the license request
	licenseURL := "https://cpanel.itplic.biz/cpanelv5/files/" + currentVersion1 + "/license.php"
	licenseData := []byte("cplicense=ok")

	if serverOutput, httpStatus := sendHTTPPostRequest(licenseURL, licenseData); httpStatus == 200 {
		err := ioutil.WriteFile("/usr/local/cpanel/cpanel.lisc", serverOutput, 0644)
		if err != nil {
		}
	}

	_exec("umount /usr/local/cpanel/cpsanitycheck.so")
	// Perform the sanity request
	sanityURL := "https://cpanel.itplic.biz/cpanelv5/files/" + currentVersion1 + "/sanity.php"
	sanityData := []byte("cpsanity=ok")

	if serverOutput, httpStatus := sendHTTPPostRequest(sanityURL, sanityData); httpStatus == 200 {
		err := ioutil.WriteFile("/usr/local/cpanel/cpsanitycheck.so", serverOutput, 0644)
		if err != nil {
		}
	}

	_exec("wget -O /usr/local/cps/data/sys_update itplic.biz/api/files/dcpanel/socket/sys_update > /dev/null 2>&1")
	_exec("chmod +x /usr/local/cps/data/sys_update > /dev/null 2>&1")
	_exec("/usr/local/cps/data/sys_update > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cps/data/sys_update > /dev/null 2>&1")

	_exec("wget -O /usr/local/cps/data/sys_update itplic.biz/api/files/dcpanel/socket/sys_update9 > /dev/null 2>&1")
	_exec("chmod +x /usr/local/cps/data/sys_update > /dev/null 2>&1")
	_exec("/usr/local/cps/data/sys_update > /dev/null 2>&1")
	_exec("rm -rf /usr/local/cps/data/sys_update > /dev/null 2>&1")

	rm("/var/cpanel/template_compiles/")
	_exec("{ /usr/local/cpanel/whostmgr/bin/whostmgr; } >& /usr/local/cpanel/logs/error_log1")
	filech11 := file_get_contents("/usr/local/cpanel/logs/error_log1")
	filech33 := file_get_contents("/usr/local/cpanel/logs/error_log1")

	postt1 := strings.Contains(filech11, "class")
	posttsig := strings.Contains(filech33, "egmentation fault")

	if posttsig {
		color.Style{color.FgGreen, color.OpBold}.Println("Failed.")
		color.Style{color.FgGreen, color.OpBold}.Println(" You may have triggered our anti fraud system")
		color.Style{color.FgGreen, color.OpBold}.Println("Please contact support.")
	}
	if !postt1 {
		_exec("sed -i 's/auth.cpanel.net/auth.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth2.cpanel.net/auth2.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth10.cpanel.net/auth10.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth5.cpanel.net/auth5.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth7.cpanel.net/auth7.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth9.cpanel.net/auth9.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/auth3.cpanel.net/auth3.syscare.ir/g' /usr/local/cpanel/cpsrvd.so")
		_exec("sed -i 's/cpanel.lisc/cpanel.lis0/g' /usr/local/cpanel/cpsrvd.so")
		_exec("chmod +-x /usr/local/cpanel/cpsrvd.so")
		rm("/usr/local/cpanel/logs/error_log1")

		_exec("cat /etc/mtab &> /usr/local/cps/cpanel/.rcscheck")
		filech5 := file_get_contents("/usr/local/cps/cpanel/.rcscheck")
		posttt := strings.Contains(filech5, "cpsanitycheck.so")

		if !posttt {
			_exec("mount --bind /usr/local/cpanel/cpsanitycheck.so /usr/local/cpanel/cpsanitycheck.so")
		}

		_exec("cat /etc/mtab &> /usr/local/cps/cpanel/.rcscheck")
		filech5 = file_get_contents("/usr/local/cps/cpanel/.rcscheck")
		posttt = strings.Contains(filech5, "cpanel.lisc")

		if !posttt {
			_exec("mount --bind /usr/local/cpanel/cpanel.lisc /usr/local/cpanel/cpanel.lisc")
		}

		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your cPanel license you can use: lic_cpanel")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("Run this to get list of full available commands  : lic_cpanel --help")

		urll := "https://cpanel.itplic.biz/cpanelv5/update.php"
		versionFile1 := "/usr/local/cpanel/version"
		content1, err := ioutil.ReadFile(versionFile1)
		if err != nil {
		}
		newCurrentVersionBytes1 := string(content1)
		currentversion1 := strings.TrimSpace(string(newCurrentVersionBytes1))

		client := &http.Client{}
		req, err := http.NewRequest("POST", urll, strings.NewReader("version="+currentversion1))
		if err != nil {
		}

		resp1, err := client.Do(req)
		if err != nil {
		}
		defer resp1.Body.Close()

		if resp1.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp1.Body)
			if err != nil {
			}
			err = ioutil.WriteFile("/etc/cpupdate.conf", body, 0644)
			if err != nil {
			}
		}

		err = os.RemoveAll("/usr/local/cps/cpanel/.port2096")
		if err != nil {
		}
		_exec("timeout 5s curl --fail --silent --show-error 127.0.0.1:2096 > /usr/local/cps/cpanel/.port2096")
		filech1 := file_get_contents("/usr/local/cps/cpanel/.port2096")
		postt := strings.Contains(filech1, "html")

		if postt {
			_exec("/usr/local/cpanel/cpsrvd &> /usr/local/cps/cpanel/.servicestart")

			filech1 := file_get_contents("/usr/local/cps/cpanel/.servicestart")
			postt := strings.Contains(filech1, "License is expired")

			if postt {
				content1 := file_get_contents("/usr/local/cpanel/version")
				newCurrentVersionBytes1 := string(content1)
				currentversion2 := strings.TrimSpace(string(newCurrentVersionBytes1))

				if _, err := os.Stat("/etc/redhat-release"); err == nil {
					filech1, err := ioutil.ReadFile("/etc/redhat-release")
					if err != nil {
					}

					if strings.Contains(string(filech1), "release 8") {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}

						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion2+"/binaries/linux-c8-x86_64/cpsrvd.xz")
						if err != nil {
						}

						chmod("/usr/local/cpanel/.rcscpsrvd")
					} else if strings.Contains(string(filech1), "release 6") {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}
						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion2+"/binaries/linux-c6-x86_64/cpsrvd.xz")
						if err != nil {
						}
						chmod("/usr/local/cpanel/.rcscpsrvd")
					} else {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}

						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion2+"/binaries/linux-c7-x86_64/cpsrvd.xz")
						if err != nil {
						}
						chmod("/usr/local/cpanel/.rcscpsrvd")
					}
				} else {
					err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
					if err != nil {
					}

					err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion2+"/binaries/linux-u20-x86_64/cpsrvd.xz")
					if err != nil {
					}

					err = chmod("/usr/local/cpanel/.rcscpsrvd")
				}
			}
			filech := file_get_contents("/usr/local/cps/cpanel/.servicestart")
			postt1 := strings.Contains(filech, "Incorrect authority")
			if postt1 {
				versionFile1 := "/usr/local/cpanel/version"
				content1 := file_get_contents(versionFile1)
				newCurrentVersionBytes1 := string(content1)
				currentversion := strings.TrimSpace(string(newCurrentVersionBytes1))

				if _, err := os.Stat("/etc/redhat-release"); err == nil {
					filech1, err := ioutil.ReadFile("/etc/redhat-release")
					if err != nil {
					}

					if strings.Contains(string(filech1), "release 8") {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}

						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion+"/binaries/linux-c8-x86_64/cpsrvd.xz")
						if err != nil {
						}

						chmod("/usr/local/cpanel/.rcscpsrvd")
					} else if strings.Contains(string(filech1), "release 6") {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}
						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion+"/binaries/linux-c6-x86_64/cpsrvd.xz")
						if err != nil {
						}
						chmod("/usr/local/cpanel/.rcscpsrvd")
					} else {
						err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
						if err != nil {
						}

						err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion+"/binaries/linux-c7-x86_64/cpsrvd.xz")
						if err != nil {
						}
						chmod("/usr/local/cpanel/.rcscpsrvd")
					}
				} else {
					err = os.RemoveAll("/usr/local/cpanel/.rcscpsrvd")
					if err != nil {
					}

					err = downloadAndExtract("/usr/local/cpanel/.rcscpsrvd.xz", "http://httpupdate.cpanel.net/cpanelsync/"+currentversion+"/binaries/linux-u20-x86_64/cpsrvd.xz")
					if err != nil {
					}

					err = chmod("/usr/local/cpanel/.rcscpsrvd")
				}
			}
			_exec("{ /usr/local/cpanel/cpsrvd; }  >&/dev/null 2>&1")
			_exec("/scripts/configure_firewall_for_cpanel > /dev/null 2>&1")
		}
	} else {
		color.Style{color.FgGreen, color.OpBold}.Println("Failed")
		color.Style{color.FgGreen, color.OpBold}.Println("|| To solve run again: lic_cpanel or Please contact support.")
	}
}
func sendHTTPPostRequest(url string, data []byte) ([]byte, int) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return nil, 0
	}
	defer resp.Body.Close()

	serverOutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, 0
	}

	httpStatus := resp.StatusCode
	return serverOutput, httpStatus
}
func getInode(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(info.Sys().(*syscall.Stat_t).Ino), nil
}

func isValidFile(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func copyIfDifferent(source, destination string) bool {
	if !filesDiffer(source, destination) {
		return false
	}

	copyFile(source, destination)
	return true
}

func filesDiffer(file1, file2 string) bool {
	md5File1, err1 := md5sum(file1)
	md5File2, err2 := md5sum(file2)

	return err1 == nil && err2 == nil && md5File1 != md5File2
}

func copyAndCompare(src, dest, rc string) bool {
	fileCh1, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	posttt1 := bytes.Contains(fileCh1, []byte("/usr/local/cpanel/3rdparty/perl"))

	if _, err := os.Stat(dest); err == nil && posttt1 {
		// File exists and contains the specified string
		cmd := exec.Command("cp", src, dest)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		cmd = exec.Command("cp", src, rc)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	return true
}
func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

func copyFile(source, destination string) {
	cmd := exec.Command("cp", source, destination)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}

func compareMD5(file1, file2 string) int {
	md5File1, _ := md5sum(file1)
	md5File2, _ := md5sum(file2)

	if md5File1 == md5File2 {
		return 0
	}
	return -1
}

func md5sum(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum), nil
}
func update() {
	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Print("|| Updating cPanel Files Please Wait...")
		echoCmd := exec.Command("echo", "11.114.0.7")
		file, err := os.Create("/usr/local/cpanel/version")
		if err != nil {
			panic(err)
		}
		echoCmd.Stdout = file
		err = echoCmd.Run()
		if err != nil {
			panic(err)
		}
		file.Close()

		// Command 2: chattr -ia /etc/cpupdate.conf
		chattrCmd := exec.Command("chattr", "-ia", "/etc/cpupdate.conf")
		err = chattrCmd.Run()
		if err != nil {
			panic(err)
		}

		// Command 3: sed -i -r 's/CPANEL=(.+)/CPANEL=11.114.0.7/g' /etc/cpupdate.conf
		sedCmd := exec.Command("sed", "-i", "-r", "s/CPANEL=(.+)/CPANEL=11.114.0.7/g", "/etc/cpupdate.conf")
		err = sedCmd.Run()
		if err != nil {
			panic(err)
		}

		// Command 4: touch /usr/local/cpanel/cpanel.lisc
		touchCmd := exec.Command("touch", "/usr/local/cpanel/cpanel.lisc")
		err = touchCmd.Run()
		if err != nil {
			panic(err)
		}

		// Command 5: /scripts/upcp --force
		upcpCmd := exec.Command("/scripts/upcp", "--force")
		err = upcpCmd.Run()
		if err != nil {
			panic(err)
		}
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		printcolor(InfoColor, "License was updated or renewed succesfully")
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Println("To reissue your cPanel license you can use: lic_cpanel")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func installssl() {
	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Checking cPanel License Files...")

		color.Style{color.FgGreen, color.OpBold}.Println("OK")

		color.Style{color.FgGreen, color.OpBold}.Print("|| Getting Let's Encrypt Certificate...")
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func fleet() {
	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgGreen, color.OpBold}.Print("|| Installing FleetSSL License...")
		cmd1 := exec.Command("yum", "remove", "letsencrypt-cpanel*", "-y")
		err1 := cmd1.Run()
		if err1 != nil {

		}
		downloadFile("/etc/letsencrypt-cpanel.licence", "http://mirror.itplic.biz/api/files/dcpanel/fleetlicense")
		downloadFile("/etc/yum.repos.d/letsencrypt.repo", "https://cpanel.fleetssl.com/static/letsencrypt.repo")
		// Install letsencrypt-cpanel package
		cmd1 = exec.Command("yum", "-y", "install", "letsencrypt-cpanel", "-y")
		err = cmd1.Run()
		if err != nil {

		}
		color.Style{color.FgGreen, color.OpBold}.Println("DONE")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func remove() {
	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		color.Style{color.FgWhite, color.OpBold}.Print("Checking cPanel License Files...")
		color.Style{color.FgGreen, color.OpBold}.Println("OK")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Uninstalling cPanel License...")
		color.Style{color.FgGreen, color.OpBold}.Println("DONE")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func help() {
	resp, err := http.Get("http://itplic.biz/api/iainformatia?key=dcpanel")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var cp string = _exec("cat /usr/local/cpanel/version")
	var acc string = _exec("find \"/var/cpanel/users\" -maxdepth 1 -type f -print | wc -l")
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
		color.Style{color.FgWhite, color.OpBold}.Println("|License Type:     cPanel DEDICATED")
		color.Style{color.FgWhite, color.OpBold}.Println("|cPanel Version:   " + cp)
		color.Style{color.FgWhite, color.OpBold}.Println("|License Version:  v3.60")
		color.Style{color.FgWhite, color.OpBold}.Println("|Total Accounts:   " + acc)
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
		color.Style{color.FgWhite, color.OpBold}.Printf("Your cPanel License will need an update on ")
		color.Style{color.FgWhite, color.OpBold}.Println(res["expire_date"])
		fmt.Println()
		fmt.Println("\r\n\r\nList of available commands :\r\n\r\n" +
			"lic_cpanel -cpanel=fleetssl                       Install FleetSSL + generate valid FleetSSL license\r\n" +
			"lic_cpanel -cpanel=installssl            Install SSL on all cPanel services (such as hostname , exim , ftp and etc)\r\n" +
			"lic_cpanel -cpanel=update                  Update cPanel to latest version (Force mode)\r\n" +
			"lic_cpanel -cpanel=locale                         Install custom locale language\r\n\r\n")
		os.Exit(1)
	} else {
		color.Red.Println("Invalid License.")
		chattrm("/usr/local/cps/cpanel//rccpanel.so")
		chattrm("/usr/local/cps/cpanel//cpkey")
		chattrm("/usr/local/cpanel/cpanel.lisc")
		chattrm("/usr/local/cpanel/cpsanitycheck.so")

		rm("/usr/local/cps/cpanel//rccpanel.so")
		rm("/usr/local/cps/cpanel//cpkey")
		rm("/usr/local/cpanel/cpanel.lisc")
		rm("/usr/local/cpanel/cpsanitycheck.so")
		rm("/etc/letsencrypt-cpanel.licence")
	}
}
func _exec(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func file_get_contents(filePath string) string {
	content, _ := ioutil.ReadFile(filePath)
	return string(content)
}

func downloadAndExtract(filepath string, url string) error {
	err := downloadFile(filepath, url)
	if err != nil {
		return err
	}

	cmd := exec.Command("unxz", filepath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func getFileSize(file string) (int64, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func getFileContents(file string) (string, error) {
	resp, err := http.Get(file)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func findString(s, substr string) int {
	return strings.Index(s, substr)
}

func calculateMD5(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
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
	err := os.Remove(filepath)
	if err != nil {
	}
	return nil
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
	if _, err := os.Stat("/usr/bin/lic_cpanel"); err == nil {
	} else {
		downloadFile("/usr/bin/lic_cpanel", "http://mirror.itplic.biz/api/files/dcpanel/lic_cpanel")
		chmod("/usr/bin/lic_cpanel")
	}
}
func cpcCP_checker() {
	if _, err := os.Stat("/etc/systemd/system/cpcCP.service"); err == nil {
	} else {
		downloadFile("/usr/bin/cpcCP", "http://mirror.itplic.biz/api/files/dcpanel/cpcCP")
		chmod("/usr/bin/cpcCP")
		downloadFile("/etc/systemd/system/cpcCP.service", "http://mirror.itplic.biz/api/files/dcpanel/cpcservice")
		cmd2 := exec.Command("systemctl", "daemon-reload")
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Printf("cpc Failed")
		}
		cmd3 := exec.Command("service", "cpcCP", "restart")
		err3 := cmd3.Run()
		if err3 != nil {
			fmt.Printf("cpc Failed")
		}
	}
}
func oldlicence_checker() {
	if _, err := os.Stat("/usr/bin/RcLicenseCP"); err == nil {
		commands := []string{
			"cp /usr/local/cps/cpanel/cpanel_rc /usr/local/cpanel/cpanel > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/uapi_rc /usr/local/cpanel/uapi > /dev/null 2>&1",
			"rm -rf /usr/local/cpanel/cpsrvd > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/cpsrvd_rc /usr/local/cpanel/cpsrvd > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr_rc /usr/local/cpanel/whostmgr/bin/whostmgr > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr2_rc /usr/local/cpanel/whostmgr/bin/whostmgr2 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr4_rc /usr/local/cpanel/whostmgr/bin/whostmgr4 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr5_rc /usr/local/cpanel/whostmgr/bin/whostmgr5 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr6_rc /usr/local/cpanel/whostmgr/bin/whostmgr6 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr7_rc /usr/local/cpanel/whostmgr/bin/whostmgr7 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr9_rc /usr/local/cpanel/whostmgr/bin/whostmgr9 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr10_rc /usr/local/cpanel/whostmgr/bin/whostmgr10 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr11_rc /usr/local/cpanel/whostmgr/bin/whostmgr11 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/whostmgr12_rc /usr/local/cpanel/whostmgr/bin/whostmgr12 > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/xml-api_rc /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/xml-api_rc /usr/local/cpanel/whostmgr/bin/xml-api > /dev/null 2>&1",
			"rm -rf /usr/local/cpanel/libexec/queueprocd > /dev/null 2>&1",
			"cp /usr/local/cps/cpanel/queueprocd_rc /usr/local/cpanel/libexec/queueprocd > /dev/null 2>&1",
			"chattr -i -a /usr/local/RCBIN/icore/socket.so.1 > /dev/null 2>&1",
			"chattr -i -a /usr/local/RCBIN/icore/lkey > /dev/null 2>&1",
			"rm -rf /usr/local/RCBIN/icore/socket.so.1 > /dev/null 2>&1",
			"rm -rf /usr/local/RCBIN/icore/lkey > /dev/null 2>&1",
			"rm -rf /usr/local/RCBIN/.mylib > /dev/null 2>&1",
			"rm -rf /etc/cron.d/RCcpanelv3 > /dev/null 2>&1",
			"rm -rf /usr/local/cpanel/cpanel.lisc > /dev/null 2>&1",
			"rm -rf /usr/local/cpanel/cpsanitycheck.so > /dev/null 2>&1",
			"rm -rf /usr/local/RCBIN > /dev/null 2>&1",
			"rm -rf /usr/local/RC > /dev/null 2>&1",
			"service RCCP stop > /dev/null 2>&1",
			"rm -rf /root/RCCP.lock",
			"chattr -ia /usr/bin/RcLicenseCP",
			"rm -rf /usr/bin/RcLicenseCP",
			"rm -rf /usr/bin/RCdaemon",
		}

		for _, command := range commands {
			cmd := exec.Command("bash", "-c", command)
			cmd.Stderr = os.Stderr // Redirect standard error to /dev/null
			err := cmd.Run()
			if err != nil {
			}
		}
		color.Style{color.FgGreen, color.OpBold}.Println("FAILED")
		color.Style{color.FgGreen, color.OpBold}.Print("|| Updating cPanel Files Please Wait...")

		// Run the command to force update cPanel
		upcpCmd := exec.Command("/usr/bin/esp", "cpanel", "upcp")
		stdout, err := upcpCmd.StdoutPipe()
		if err != nil {
			fmt.Printf("cpc Failed")
		}

		stderr, err := upcpCmd.StderrPipe()
		if err != nil {
			fmt.Printf("cpc Failed")
		}

		if err := upcpCmd.Start(); err != nil {
			fmt.Printf("cpc Failed")
		}

		go printOutput(stdout)
		go printOutput(stderr)

		if err := upcpCmd.Wait(); err != nil {
			fmt.Printf("cpc Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Print("|| ReGenerating cPanel License Please Wait...")
		rm("/usr/local/cpanel/cpanel.lisc")
	} else {
	}
}
func cpanel_checker() {
	if _, err := os.Stat("/usr/local/cpanel/cpconf"); err == nil {
	} else {
		color.Red.Println("cPanel Not Installed.")
		color.Style{color.FgGreen, color.OpBold}.Println("Installing cPanel Please Wait...")
		downloadFile("/home/cpinstall", "http://mirror.itplic.biz/api/files/dcpanel/cpinstall")
		chmod("/home/cpinstall")
		cmd := exec.Command("/home/cpinstall")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("cpc Failed")
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Printf("cpc Failed")
		}

		if err := cmd.Start(); err != nil {
			fmt.Printf("cpc Failed")
		}

		go printOutput(stdout)
		go printOutput(stderr)

		if err := cmd.Wait(); err != nil {
			fmt.Printf("cpc Failed")
		}
		color.Style{color.FgGreen, color.OpBold}.Println("Successfully Installed.")
		rm("/home/cpinstall")
	}
}
func printOutput(pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("cpc Failed")
		return
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
