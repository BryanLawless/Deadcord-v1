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
	"math/rand"
	"strconv"
)

func StartMassThreadCreateThreads(channel_id string, thread_name string) {
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, channel_id string, thread_name string) {
			defer wg.Done()
			massThreadWorker(token, channel_id, thread_name)
		}(token, channel_id, thread_name)
	}

	wg.Wait()
}

func massThreadWorker(token string, channel_id string, thread_name string) {
	status, status_code, _ := requests.SendDiscordRequest("channels/"+channel_id+"/threads", "POST", token, "general", map[string]interface{}{}, map[string]interface{}{
		"name":                  thread_name + strconv.Itoa(rand.Intn(1000-1)+1),
		"type":                  11,
		"auto_archive_duration": 1440,
	})

	if status {
		switch status_code {
		case 201:
			util.WriteToConsole("Bot created thread: "+thread_name+".", 2)
		case 429:
			util.WriteToConsole("Thread request was rate limited.", 1)
		default:
			util.WriteToConsole("Bot could not create thread, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
