package main

import (
	"fmt"
	"github.com/electricbubble/gadb"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type output struct {
	Name     string
	Target   int
	Filename string
}

var CoreStatus = map[string]string{"Off": "0", "On": "1"}
var CoreNumber = map[string]string{"Core_0": "0", "Core_1": "1", "Core_2": "2", "Core_3": "3", "Core_4": "4", "Core_5": "5", "Core_6": "6", "Core_7": "7"}
var FreqNumber = map[string]string{"Freq_0": "0", "Freq_1": "1", "Freq_2": "2", "Freq_3": "3", "Freq_4": "4", "Freq_5": "5", "Freq_6": "6", "Freq_7": "7"}

func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	adbClient, _ := gadb.NewClient()
	devices, _ := adbClient.DeviceList()
	dev := devices[0]
	if r.Method == "GET" {

		template.ParseFiles("index.html")
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, output{"", 0, ""})
		//memory := r.Form["memory"][0]
		//main_title:=output{out,target,"static/"+name+"_"+strconv.Itoa(target)+".png"}
		//t, _ := template.ParseFiles("index.html")
		//t.Execute(w, main_title)
	} else {

		r.ParseForm()
		//fmt.Println(r.Form)
		for i, v := range r.Form {
			if strings.Contains(i, "Core") {
				dev.RunShellCommand("su -c 'echo " + CoreStatus[v[0]] + " > /sys/devices/system/cpu/cpu" + CoreNumber[i] + "/online'")
			}

			if strings.Contains(i, "Freq") {
				dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + FreqNumber[i] + "/cpufreq/scaling_setspeed'")
				dev.RunShellCommand("su -c 'echo " + v[0] + " >  /sys/devices/system/cpu/cpu" + FreqNumber[i] + "/cpufreq/scaling_max_freq'")

			}
			//fmt.Println(i, v[0])
		}

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, output{"", 0, ""})

	}
}

func main() {

	http.HandleFunc("/", MainPage)
	err := http.ListenAndServe(":80", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
