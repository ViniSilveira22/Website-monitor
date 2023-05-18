package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 3
const delay = 10

func main() {
	showIntroduction()
	for {
		showOptions()
		options()
	}
}

func showIntroduction() {
	var version = 1.1
	fmt.Println("Hello")
	fmt.Println("This program is currently in version", version)
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func showOptions() {
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("3 - Exit")
}

func options() {
	switch readCommand() {
	case 1:
		startMonitoring()
	case 2:
		showLogs()
	case 3:
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Option does not exist, provide another")
	}
}

func startMonitoring() {
	fmt.Println("Monitoring...")
	sites := getSitesFromFile()
	for i := 0; i < monitoring; i++ {
		for _, site := range sites {
			getStatusSite(site)
		}
		time.Sleep(delay * time.Second)
	}

}

func getStatusSite(site string) {
	response, error := http.Get(site)
	if error != nil {
		fmt.Println("Error", error)
	}
	if response.StatusCode == 200 {
		fmt.Println("Website:", site, "is online")
		registerLogs(site, true)
	} else {
		fmt.Println("Website:", site, "have problems. Status Code:", response.StatusCode)
		registerLogs(site, false)
	}
}

func getSitesFromFile() []string {

	var sites []string
	file, error := os.Open("websites.txt")
	if error != nil {
		fmt.Println("Error", error)
	}
	reader := bufio.NewReader(file)

	for {
		line, error := reader.ReadString('\n')
		sites = append(sites, strings.TrimSpace(line))

		if error == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func registerLogs(site string, status bool) {
	file, error := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if error != nil {
		fmt.Println("Error", error)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - Online: " + strconv.FormatBool(status) + "\n")
	file.Close()
}

func showLogs() {
	file, error := ioutil.ReadFile("log.txt")
	if error != nil {
		fmt.Println("Error", error)
	}
	fmt.Println(string(file))
}
