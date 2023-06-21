package main

import (
	"fmt"
	"github.com/electricbubble/gadb"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var CoreStatus = map[string]string{"Off": "0", "On": "1"}
var CoreNumber = map[string]string{"Core_0": "0", "Core_1": "1", "Core_2": "2", "Core_3": "3", "Core_4": "4", "Core_5": "5", "Core_6": "6", "Core_7": "7"}
var FreqNumber = map[string]string{"Freq_0": "0", "Freq_1": "1", "Freq_2": "2", "Freq_3": "3", "Freq_4": "4", "Freq_5": "5", "Freq_6": "6", "Freq_7": "7"}

type AndroidName struct {
	Model string
}

func Simple(w http.ResponseWriter, r *http.Request) {
	adbClient, _ := gadb.NewClient()
	devices, _ := adbClient.DeviceList()
	dev := devices[0]

	name := AndroidName{Model: dev.DeviceInfo()["model"]}

	if r.Method == "GET" {

		t, _ := template.ParseFiles("simple.html")
		t.Execute(w, name)
	}
	if r.Method == "POST" {

		r.ParseForm()

		for i, v := range r.Form {

			if strings.Contains(i, "Core") {

				dev.RunShellCommand("su -c 'echo " + CoreStatus[v[0]] + " > /sys/devices/system/cpu/cpu" + CoreNumber[i] + "/online'")
			}

			if strings.Contains(i, "Freq") {
				if i == "Freq_0" {
					for j := 0; j < 4; j++ {
						dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + strconv.Itoa(j) + "/cpufreq/scaling_setspeed'")
						dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + strconv.Itoa(j) + "/cpufreq/scaling_max_freq'")

					}
				}
				if i == "Freq_4" {
					for j := 4; j < 8; j++ {

						dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + strconv.Itoa(j) + "/cpufreq/scaling_setspeed'")
						dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + strconv.Itoa(j) + "/cpufreq/scaling_max_freq'")
					}
				}
			}
		}

		t, _ := template.ParseFiles("simple.html")
		t.Execute(w, name)

	}

	//http.Redirect(w, r, "/simple", 200)

}

func MainPage(w http.ResponseWriter, r *http.Request) {
	adbClient, _ := gadb.NewClient()
	devices, _ := adbClient.DeviceList()
	dev := devices[0]

	name := AndroidName{Model: dev.DeviceInfo()["model"]}
	if r.Method == "GET" {

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, name)
	} else {

		r.ParseForm()
		//fmt.Println(r.Form)
		for i, v := range r.Form {
			fmt.Println("err111111")
			if strings.Contains(i, "Core") {
				dev.RunShellCommand("su -c 'echo " + CoreStatus[v[0]] + " > /sys/devices/system/cpu/cpu" + CoreNumber[i] + "/online'")
			}

			if strings.Contains(i, "Freq") {
				dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + FreqNumber[i] + "/cpufreq/scaling_setspeed'")
				dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + FreqNumber[i] + "/cpufreq/scaling_max_freq'")

			}
		}

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, name)

	}
}

func main() {

	http.HandleFunc("/main", MainPage)
	http.HandleFunc("/simple", Simple)

	err := http.ListenAndServe(":80", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
