/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package core

import (
	"Deadcord/util"
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func HarvestFromURL(url string, timeout int) ([]string, error) {
	proxy_match := regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):([0-9]){1,4}`)
	proxy_client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	response, err := proxy_client.Get(url)
	util.CheckError(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	proxies := proxy_match.FindAllString(string(body), -1)
	return proxies, nil
}

func LoadProxies() (bool, []string, map[int]map[string]string) {
	//var found_proxies_raw []string

	if _, err := os.Stat("./proxies.txt"); os.IsNotExist(err) {
		return false, nil, nil
	} else {
		proxies_loaded, err := parseProxyFile("./proxies.txt")
		util.CheckError(err)

		if len(proxies_loaded) > 0 {
			return true, proxies_loaded, nil
		} else {
			util.WriteToConsole("No proxies could be loaded from ./proxies.txt.", 3)
			return false, nil, nil
		}
	}
}

func parseProxyFile(path string) ([]string, error) {
	file, err := os.Open(path)
	util.CheckError(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
