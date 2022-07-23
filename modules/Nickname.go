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
	"Deadcord/util"
	"strconv"
)

func StartNickThreads(server_id string, nickname string) {
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, server_id string, nickname string) {
			defer wg.Done()
			nickWorker(token, server_id, nickname)
		}(token, server_id, nickname)
	}

	wg.Wait()
}

func nickWorker(token string, server_id string, nickname string) {

	nickname_string := ""

	switch nickname {
	case "reset":
		nickname_string = " "
	default:
		nickname_string = nickname
	}

	status, status_code, _ := requests.SendDiscordRequest("guilds/"+server_id+"/members/@me", "PATCH", token, "general", map[string]interface{}{}, map[string]interface{}{
		"nick": nickname_string,
	})

	if status {
		switch status_code {
		case 200:
			util.WriteToConsole("Bot chnaged nickname to: "+nickname_string+".", 2)
		case 429:
			util.WriteToConsole("Change nickname request rate limited.", 1)
		default:
			util.WriteToConsole("Bot could not change nickname, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
