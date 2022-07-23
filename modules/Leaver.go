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

func StartLeaveGuildThreads(server_id string) bool {
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(server_id string, token string) {
			defer wg.Done()
			leaveWorker(server_id, token)
		}(server_id, token)
	}

	wg.Wait()

	return false
}

func leaveWorker(server_id string, token string) {
	status, status_code, _ := requests.SendDiscordRequest("users/@me/guilds/"+server_id, "DELETE", token, "general", map[string]interface{}{}, map[string]interface{}{
		"lurking": false,
	})

	if status {
		if status_code == 204 {
			util.WriteToConsole("Bot successfully left guild.", 2)
		} else {
			util.WriteToConsole("Bot could not leave guild, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
