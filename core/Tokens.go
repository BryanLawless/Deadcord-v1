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
	b64 "encoding/base64"
	"fmt"
	"os"
	"strings"

	ua "github.com/mileusna/useragent"
)

var (
	duplicates   int    = 0
	token_file   string = "./tokens.txt"
	token_struct        = map[int]map[string]string{}
)

func LoadTokens() (bool, []string, map[int]map[string]string) {

	if _, err := os.Stat(token_file); os.IsNotExist(err) {
		util.WriteToConsole("No token file found, Deadcord will create one for you. Please restart Deadcord.", 3)
		token_file, err := os.Create(token_file)
		util.CheckError(err)

		token_file.Close()
		util.SoftQuit()
	}

	tokens_loaded, err := parseTokenFile(token_file)

	if duplicates > 0 {
		util.WriteToConsole(fmt.Sprintf("Removed %d duplicate tokens.", duplicates), 2)
		ResetTokenServiceWithManualTokens(tokens_loaded)
	}

	util.CheckError(err)

	if len(tokens_loaded) > 0 {
		buildProfiles(tokens_loaded)

		return true, tokens_loaded, token_struct
	} else {
		util.WriteToConsole("No tokens could be loaded from "+token_file, 3)
		util.SoftQuit()
		return false, nil, nil
	}
}

func buildProfiles(token_list []string) {

	buildSingleProfile := func(label_num int, token string) {
		random_agent := util.RandomUserAgent()
		parse_agent := ua.Parse(random_agent)

		token_parts := strings.Split(token, ".")

		if len(token_parts) == 3 {
			user_id_decode, err := b64.StdEncoding.DecodeString(token_parts[0])
			util.CheckError(err)

			token_struct[label_num]["user_id"] = string(user_id_decode[:])
			token_struct[label_num]["browser"] = parse_agent.Name
			token_struct[label_num]["token"] = token
			token_struct[label_num]["agent"] = random_agent
			token_struct[label_num]["os"] = parse_agent.OS
			token_struct[label_num]["browser_version"] = parse_agent.Version
			token_struct[label_num]["os_version"] = parse_agent.OSVersion
		} else {
			util.WriteToConsole("Invalid token found, please make sure your tokens.txt contains only valid tokens.", 3)
			return
		}
	}

	label_num := 1

	for _, token := range token_list {
		token_struct[label_num] = map[string]string{}
		buildSingleProfile(label_num, token)
		label_num++
	}

	Profiles = token_struct
}

func parseTokenFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	final_lines := util.RemoveDuplicates(lines)
	duplicates = len(lines) - len(final_lines)

	return final_lines, scanner.Err()
}

func ResetTokenServiceWithManualTokens(token_list []string) int {
	err := os.Remove(token_file)
	util.CheckError(err)

	WriteLines(token_list, token_file)

	status, raw_tokens, built_tokens := LoadTokens()

	if status {
		return SetTokens(raw_tokens, built_tokens)
	}

	return 0
}

func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
