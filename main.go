/*
 * Copyright @random_robbie (c) 2018.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/logrusorgru/aurora"
	"github.com/remeh/sizedwaitgroup"
	"io/ioutil"
	"net/url"
)

var (
	fileofurls  = "list.txt"
	outputfile  = "./cfg/"
	filepathurl = "/assets../settings/90-local.conf"
	au          aurora.Aurora
	colors      = flag.Bool("colors", true, "enable or disable colors")
)

func init() {
	flag.Parse()
	au = aurora.NewAurora(*colors)
}

func grabURL(URL string, output string, filepathurl string, swg *sizedwaitgroup.SizedWaitGroup) {

	defer swg.Done()


	newurl := URL + filepathurl



	req, err := http.NewRequest("GET", newurl, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	req.Header.Set("Connection", "close")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(au.Red("[*] There was an issue connecting to the server. [*]"))
		return 
	}
	defer resp.Body.Close()


	fmt.Println("[*] Testing ", newurl, "[*]")





	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		fmt.Println(au.Red("[*] Redirected ... Not Vulnerable [*]"))
	}


		if resp.Body != nil {
			bodyBytes, nil := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			htmlcontent := strings.Contains(bodyString, "STATIC_URL")
			if htmlcontent == true {

				u, err := url.Parse(URL)
				if err != nil {
					log.Fatal(err)
				}


				capturedfile := outputfile + u.Hostname() + ".txt"

				htmlfile, err := os.Create(capturedfile)

				if err != nil {
					log.Fatalf("could not create file: %v", err)
					os.Exit(1)
				}

				defer htmlfile.Close()
				fmt.Println(au.Green("[*] Saving to file [*]"))
				io.Copy(htmlfile, resp.Body)

			} else {
				fmt.Println(au.Red("[*] Not Vulnerable  [*]"))
			}

		}

	}


func main() {
	swg := sizedwaitgroup.New(10)
	fmt.Println(au.Blue("[*] Nginx Alias Checker - By @random_robbie [*]"))
	file, err := os.Open(fileofurls)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		swg.Add()
		URL := strings.TrimSpace(scanner.Text())
		go grabURL(URL, outputfile, filepathurl, &swg)
		swg.Wait()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}