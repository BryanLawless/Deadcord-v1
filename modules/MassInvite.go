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
	"Deadcord/constants"
	"Deadcord/core"
	"Deadcord/requests"
	"Deadcord/util"
	"encoding/json"
	"fmt"
)

func StartMassInviteThreads(channel_id string) bool {
	wg.Add(len(core.RawTokensLoaded))

	for _, token := range core.RawTokensLoaded {
		go func(server_id string, token string) {
			defer wg.Done()
			inviteCreateWorker(channel_id, token)
		}(channel_id, token)
	}

	wg.Wait()

	return false
}

func inviteCreateWorker(channel_id string, token string) {

	week_seconds := 604800

	var status bool
	var status_code int
	var invite_json []byte

	for i := 1; i <= 1000; i++ {

		if core.ActionFlag == 1 {
			return
		}

		if i < 101 {
			status, status_code, invite_json = requests.SendDiscordRequest("channels/"+channel_id+"/invites", "POST", token, "general", map[string]interface{}{}, map[string]interface{}{
				"max_age":   0,
				"max_uses":  i,
				"temporary": false,
			})
		} else {
			status, status_code, invite_json = requests.SendDiscordRequest("channels/"+channel_id+"/invites", "POST", token, "general", map[string]interface{}{}, map[string]interface{}{
				"max_age":  week_seconds,
				"max_uses": 0,
			})

			week_seconds--
		}

		if status {
			switch status_code {
			case 200:
				util.WriteToConsole("Bot succesfully created invite.", 2)
			case 403:
				util.WriteToConsole("Bot could not create invite, no access.", 1)
			case 429:
				var invite_create constants.RateLimit
				err := json.Unmarshal(invite_json, &invite_create)
				util.CheckError(err)

				retry_when := util.RoundFloat(invite_create.RetryAfter, 0)

				util.WriteToConsole("Thread Paused: "+fmt.Sprintf("%.2f", retry_when)+" seconds.", 1)

				util.SleepExact(retry_when)
			default:
				util.WriteToConsole("Bot could not create invite, request failed. Code: "+fmt.Sprintf("%d", status_code), 1)
			}
		}
	}
}
