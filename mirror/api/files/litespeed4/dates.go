package main

import (
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Data struct {
	Status string `json:"status"`
	Brand  string `json:"brand_name"`
	Domain string `json:"domain_name"`
	Expiry string `json:"expire_date"`
}

func main() {
	resp, err := http.Get("http://itplic.biz/api/getinfo?key=litespeed4")
	if err != nil {
		os.Exit(1)
	}
	byteResult, err := ioutil.ReadAll(resp.Body)

	var f Data
	err = json.Unmarshal(byteResult, &f)
	var res map[string]interface{}
	json.Unmarshal([]byte(byteResult), &res)
	if f.Status == "success" {
		dtstr1 := fmt.Sprint(res["expire_date"])
		dt, _ := time.Parse("2006-01-02", dtstr1)
		dtstr2, _ := strconv.Atoi(dt.Format("2006"))
		dtstr3, _ := strconv.Atoi(dt.Format("01"))
		dtstr4, _ := strconv.Atoi(dt.Format("02"))
		dt1 := time.Now()
		dtstr5, _ := strconv.Atoi(dt1.Format("2006"))
		dtstr6, _ := strconv.Atoi(dt1.Format("01"))
		dtstr7, _ := strconv.Atoi(dt1.Format("02"))
		t1 := Date(dtstr5, dtstr6, dtstr7)
    		t2 := Date(dtstr2, dtstr3, dtstr4)
    		days := t2.Sub(t1).Hours() / 24
    		fmt.Println(days)
		os.Exit(1)
	}else {
		fmt.Println("0")
		chattrm("/usr/local/lsws/conf/license.key")
		chattrm("/usr/local/lsws/conf/trial.key")
		chattrm("/usr/local/lsws/conf/serial.no")
		chattrm("/usr/local/lsws/conf/serial2.no")
		rm("/usr/local/lsws/conf/trial.key")
		rm("/usr/local/lsws/conf/license.key")
		rm("/usr/local/lsws/conf/serial2.no")
		rm("/usr/local/lsws/conf/serial.no")
		cmd := exec.Command("/usr/local/lsws/bin/lswsctrl", "restart")
		err := cmd.Run()
		if err != nil {}
	}
}
func chattrm(filepath string) error {
      cmd := exec.Command("chattr", "-i", "-a", filepath)
      return cmd.Run()
}
func rm(filepath string) error {
      cmd := exec.Command("rm", "-rf", filepath)
      return cmd.Run()
}
func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
