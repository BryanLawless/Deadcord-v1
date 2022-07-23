/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package util

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	ColorReset = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	Black      = "\u001b[30;1m"
)

func CheckError(err error) {
	if err != nil {
		log_file, _ := os.OpenFile("./deadcord.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

		defer log_file.Close()

		log_file.WriteString(err.Error() + "\n")
	}
}

func SoftQuit() {
	WriteToConsole("Press the enter key to exit...", 0)
	fmt.Scanln()
	os.Exit(1)
}

func GetQuote() string {
	quotes := []string{
		"{ Grass }",
		"This ain't hacking",
		"Lagging Discord since last month",
		"What are you here for",
		"Great to see you again",
		"I am sentient",
		"Tokens not included",
		"R.I.P Groovy & Rhythm o7",
		"We built an entire GUI before switching to a terminal window",
		"Deadcord never dies. Its already dead",
		"Some assembly required",
		"Deadcord is a rat ( This is a joke lmao )",
		"GravityNet was here :0",
		"Tokens not included",
		"Discord rick-rolled us",
		"Sqeaky clean on VirusTotal |:)",
		"Game over man, game over",
		"Discord Hammer -> Hammercord -> Deadcord",
	}

	random_quote := quotes[rand.Intn(len(quotes))]

	return random_quote
}

func GetTimestamp() string {
	current_time := time.Now()
	return current_time.Format("15:04:05")
}

func WriteToConsole(status string, mode int) {
	switch mode {
	case 0:
		fmt.Fprintln(color.Output, White+"[ INFO ] "+"[ "+GetTimestamp()+" ] "+status+ColorReset)
	case 1:
		fmt.Fprintln(color.Output, Yellow+"[ WARNING ] "+"[ "+GetTimestamp()+" ] "+status+ColorReset)
	case 2:
		fmt.Fprintln(color.Output, Purple+"[ SUCCESS ] "+"[ "+GetTimestamp()+" ] "+status+ColorReset)
	case 3:
		fmt.Fprintln(color.Output, Red+"[ ERROR ] "+"[ "+GetTimestamp()+" ] "+status+ColorReset)
	}
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func SleepExact(seconds_nearest float64) {
	time.Sleep(time.Duration(seconds_nearest) * time.Second)
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func NumberSliceCounts(arr []int) map[int]int {
	dict := make(map[int]int)
	for _, num := range arr {
		dict[num] = dict[num] + 1
	}

	return dict
}

func Contains(elements []string, value string) bool {
	for _, search := range elements {
		if value == search {
			return true
		}
	}
	return false
}

func RemoveFromSlice(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	slice[len(slice)-1] = ""
	slice = slice[:len(slice)-1]

	return slice
}

func AllParameters(parameters []string) bool {
	needed_paramters := len(parameters)
	parameters_filled := 0

	for _, parameter := range parameters {
		if len(parameter) > 0 {
			parameters_filled++
		}
	}

	if needed_paramters == parameters_filled {
		return true
	}

	return false
}

func ReadTokenProfile(profile_build map[int]map[string]string, token string) map[string]string {
	for _, check_token_struct := range profile_build {
		for _, data_value := range check_token_struct {
			if data_value == token {
				return check_token_struct
			}
		}
	}

	return nil
}

func RemoveDuplicates(str_slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range str_slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}
