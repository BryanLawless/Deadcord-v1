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
	"errors"
	"strconv"
	"strings"
)

func CheckInServerStruct(server_id string) {
	if len(core.TokensInServers[server_id]) == 0 {
		tokens_in_server := map[string][]string{}
		for _, token := range core.RawTokensLoaded {
			if InServer(server_id, token) {
				tokens_in_server[server_id] = append(tokens_in_server[server_id], token)
			}
		}

		core.TokensInServers = tokens_in_server
	}
}

func InServer(server_id string, token string) bool {
	status, status_code, _ := requests.SendDiscordRequest("guilds/"+server_id, "GET", token, "general", map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		return true
	}

	return false
}

func GetMessages(channel_id string, amount int, token string) ([]byte, error) {
	status, status_code, messages_json := requests.SendDiscordRequest("channels/"+channel_id+"/messages?limit="+strconv.Itoa(amount), "GET", token, "general", map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		return messages_json, nil
	}

	return nil, errors.New("get messages request failed, code not ok")
}

func GetGuildIdAndChannelIdFromInvite(invite string) (string, string, error) {
	invite_parts := strings.Split(invite, "/")
	invite_code := ""

	if len(invite_parts) > 1 {
		if invite_parts[3] == "invite" {
			invite_code = invite_parts[4]
		} else {
			invite_code = invite_parts[3]
		}
	} else {
		invite_code = invite_parts[0]
	}

	status, status_code, invite_json := requests.SendDiscordRequest("invites/"+invite_code, "GET", core.RandomToken(), "general", map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		var invite constants.Invite
		err := json.Unmarshal(invite_json, &invite)
		util.CheckError(err)

		return invite.Guild.ID, invite.Channel.ID, nil
	}

	return "", "", errors.New("get guild from invite request failed, code not ok")
}

func GetChannels(server_id string) (int, []string) {
	var channels []string

	CheckInServerStruct(server_id)

	channel_status, channel_status_code, channel_json := requests.SendDiscordRequest("guilds/"+server_id+"/channels", "GET", core.RandomTokenInServer(server_id), "general", map[string]interface{}{}, map[string]interface{}{})

	if channel_status && channel_status_code == 200 {
		var result constants.GuildChannels
		err := json.Unmarshal(channel_json, &result)
		util.CheckError(err)

		for _, channel := range result {
			if channel.Type == 0 {
				channels = append(channels, channel.Name+":"+channel.ID)
			}
		}
	}

	return channel_status_code, channels
}

func GuildHasCommunityVerification(server_id string) bool {
	CheckInServerStruct(server_id)

	status, status_code, guild_json := requests.SendDiscordRequest("guilds/"+server_id, "GET", core.RandomTokenInServer(server_id), "general", map[string]interface{}{}, map[string]interface{}{})

	if status && status_code == 200 {
		var guild constants.Guild
		err := json.Unmarshal(guild_json, &guild)
		util.CheckError(err)

		for _, feature := range guild.Features {
			if feature == "MEMBER_VERIFICATION_GATE_ENABLED" {
				return true
			}
		}
	}

	return false
}
