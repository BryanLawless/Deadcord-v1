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
	"strings"

	"github.com/enescakir/emoji"
)

func StartReactThreads(channel_id string, message_id string, emoji string, suffix bool) {
	wg.Add(len(core.RawTokensLoaded))
	for _, token := range core.RawTokensLoaded {
		go func(token string, channel_id string, message_id string, emoji string, suffix bool) {
			defer wg.Done()
			reactWorker(token, channel_id, message_id, emoji, suffix)
		}(token, channel_id, message_id, emoji, suffix)
	}

	wg.Wait()
}

func reactWorker(token string, channel_id string, message_id string, emoji_string string, suffix bool) {
	react_emoji := ""

	if suffix {
		react_emoji = strings.TrimSuffix(emoji.Parse(":"+emoji_string+":"), " ")
	} else {
		react_emoji = emoji_string
	}

	if !strings.Contains(react_emoji, ":") {
		status, status_code, _ := requests.SendDiscordRequest("channels/"+channel_id+"/messages/"+message_id+"/reactions/"+react_emoji+"/@me", "PUT", token, "general", map[string]interface{}{}, map[string]interface{}{})

		if status {
			switch status_code {
			case 204:
				util.WriteToConsole("Bot reacted with: [ "+react_emoji+" ].", 2)
			case 429:
				util.WriteToConsole("Reaction request was rate limited.", 1)
			default:
				util.WriteToConsole("Bot could not react, request failed. Code:  "+strconv.Itoa(status_code), 3)
			}
		}

	} else {
		util.WriteToConsole("Do not use ':' in the emoji name.", 1)
	}
}
