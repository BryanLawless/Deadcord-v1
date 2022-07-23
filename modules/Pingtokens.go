/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package modules

import (
	"Deadcord/core"
	"Deadcord/requests"
	"strconv"
	"strings"
)

var (
	TokenPingResults []string
)

func StartPingTokens() []string {
	if len(TokenPingResults) > 0 {
		TokenPingResults = nil
	}

	wg.Add(len(core.RawTokensLoaded))

	ping_channel := make(chan []string)

	for _, token := range core.RawTokensLoaded {
		go func(token string, ping_channel chan []string) {
			defer wg.Done()
			tokenPingWorker(token, ping_channel)
		}(token, ping_channel)
	}

	ping_channel_results := <-ping_channel

	close(ping_channel)

	wg.Wait()

	return ping_channel_results

}

func tokenPingWorker(token string, ping_results chan []string) {
	status, status_code, token_ping_json := requests.SendDiscordRequest("users/@me/library", "GET", token, "general", map[string]interface{}{}, map[string]interface{}{})
	token_ping_json_string := string(token_ping_json[:])

	if status {
		switch status_code {
		case 200:
			TokenPingResults = append(TokenPingResults, strconv.Itoa(0)+":"+token)
		case 403:
			if strings.Contains(token_ping_json_string, "Cloudflare") {
				TokenPingResults = append(TokenPingResults, strconv.Itoa(3)+":"+token)
			} else {
				TokenPingResults = append(TokenPingResults, strconv.Itoa(1)+":"+token)
			}
		case 401:
			TokenPingResults = append(TokenPingResults, strconv.Itoa(2)+":"+token)
		case 429:
			TokenPingResults = append(TokenPingResults, strconv.Itoa(4)+":"+token)
		}
	}

	if len(TokenPingResults) == len(core.RawTokensLoaded) {
		ping_results <- TokenPingResults
		return
	}
}
