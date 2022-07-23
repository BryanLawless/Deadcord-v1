/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package main

import (
	"Deadcord/core"
	"Deadcord/modules"
	"Deadcord/util"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var (
	DeadcordVersionNumberString = fmt.Sprintf("%g", core.DeadcordVersion)
)

func readyRequestCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func jsonResponse(code int, message string, data map[string]interface{}) []byte {
	response := make(map[string]interface{})
	response["code"] = code
	response["message"] = message
	response["data"] = data

	switch code {
	case 200:
		util.WriteToConsole(response["message"].(string), 2)
	case 400:
		util.WriteToConsole(response["message"].(string), 1)
	case 500:
		util.WriteToConsole(response["message"].(string), 3)
	default:
		util.WriteToConsole(response["message"].(string), 0)
	}

	raw_json_response, err := json.Marshal(response)
	util.CheckError(err)

	return raw_json_response
}

func errorResponse(message string) []byte {
	return []byte(jsonResponse(500, message, map[string]interface{}{}))
}

func allParametersError() []byte {
	return []byte(jsonResponse(400, "All parameters must be provided.", map[string]interface{}{}))
}

func pingTokens(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Pinging tokens...", 0)

	var alive []string
	var locked []string
	var limited []string
	var invalid []string
	var cloudflare []string

	token_results := modules.StartPingTokens()

	for _, token_ping_result := range token_results {
		token_ping_split := strings.Split(token_ping_result, ":")
		token := token_ping_split[1]

		switch token_ping_split[0] {
		case "0":
			alive = append(alive, token)
		case "1":
			invalid = append(invalid, token)
		case "2":
			locked = append(locked, token)
		case "3":
			cloudflare = append(cloudflare, token)
		case "4":
			limited = append(limited, token)
		}
	}

	if len(locked) > 0 || len(invalid) > 0 {
		alive_token_list := append(alive, limited...)
		dead_token_list := append(locked, invalid...)
		core.WriteLines(dead_token_list, "dead-tokens.txt")
		amount := core.ResetTokenServiceWithManualTokens(alive_token_list)

		if amount == 0 {
			util.WriteToConsole("No tokens could be reclaimed after ping. All tokens dead or invalid.", 3)
			util.Sleep(5)
			os.Exit(0)
		}
	}

	result_string := fmt.Sprintf("\n"+util.Green+"Alive: "+util.ColorReset+"%d tokens.\n"+util.Red+"Locked: "+util.ColorReset+"%d tokens.\n"+util.Yellow+"Invalid: "+util.ColorReset+"%d tokens.\n"+util.Blue+"Rate Limited: "+util.ColorReset+"%d tokens.\n"+util.Cyan+"Cloudflare Banned: "+util.ColorReset+"%d tokens.\n", len(alive), len(invalid), len(locked), len(limited), len(cloudflare))
	fmt.Fprintln(color.Output, result_string)

	w.Write(jsonResponse(200, "All tokens pinged: "+strconv.Itoa(len(alive))+" alive tokens.", map[string]interface{}{}))
}

func startSpam(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting Spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	channels := r.Form.Get("channels")
	messages := r.Form.Get("messages")
	spam_delay := r.Form.Get("delay")
	spam_mode := r.Form.Get("mode")
	spam_tts := r.Form.Get("tts")

	if util.AllParameters([]string{server_id, channels, messages, spam_delay, spam_mode, spam_tts}) {
		spam_mode_num, err := strconv.Atoi(spam_mode)

		if err != nil {
			w.Write(errorResponse("Invalid spam mode parameter type."))
			return
		}

		delay_num, err := strconv.Atoi(spam_delay)

		if err != nil {
			w.Write(errorResponse("Invalid delay parameter type."))
			return
		}

		spam_tts_bool, err := strconv.ParseBool(spam_tts)

		if err != nil {
			w.Write(errorResponse("Invalid TTS parameter type."))
			return
		}

		messages := strings.Split(messages, "\n")

		start_spam_routines := modules.StartSpamThreads(server_id, channels, messages, delay_num, spam_mode_num, spam_tts_bool)

		switch start_spam_routines {
		case 1:
			w.Write(errorResponse("Could not start spam, message content hit the character limit, or something went wrong."))
			return
		case 2:
			w.Write(errorResponse("Could not start spam, no open channels found."))
			return
		}

	} else {
		w.Write(allParametersError())
	}

}

func stopAll(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Stopping all running actions...", 0)

	core.ActionFlag = 1
	w.Write(jsonResponse(200, "Attempted to stop running actions.", map[string]interface{}{}))
}

func startTypingSpam(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting typing spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")

	if util.AllParameters([]string{channel_id}) {
		modules.StartTypingSpamThreads(channel_id)

		w.Write(jsonResponse(200, "Attempted to start typing spam.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func react(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting mass reacting...", 0)

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")
	message_id := r.Form.Get("message_id")
	emoji := r.Form.Get("emoji")

	if util.AllParameters([]string{channel_id, message_id, emoji}) {
		modules.StartReactThreads(channel_id, message_id, emoji, true)

		w.Write(jsonResponse(200, "Bots attempted to react.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func changeNick(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting mass nickname...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	nickname := r.Form.Get("nickname")

	if util.AllParameters([]string{server_id, nickname}) {
		modules.StartNickThreads(server_id, nickname)

		w.Write(jsonResponse(200, "Bots attempted to change their nickname.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func joinGuild(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to join guild...", 0)

	r.ParseForm()
	guild_invite := r.Form.Get("invite")
	join_delay := r.Form.Get("delay")

	if util.AllParameters([]string{guild_invite, join_delay}) {
		join_result_number := 0

		delay, err := strconv.Atoi(join_delay)
		if err != nil {
			w.Write(errorResponse("Invalid delay parameter type."))
			return
		}

		server_id, channel_id, err := modules.GetGuildIdAndChannelIdFromInvite(guild_invite)

		if err != nil {
			w.Write(errorResponse("Unable to get guild ID from invite. Invite invalid."))
			return
		}

		join_result_number = modules.StartJoinGuildThreads(guild_invite, server_id, channel_id, delay)

		if join_result_number > 0 {

			if len(server_id) > 0 {
				util.WriteToConsole("Attempting to auto-verify bots.", 2)
				status := modules.StartAutoVerifyThreads(server_id)

				switch status {
				case 1:
					w.Write(errorResponse("No verification messages found."))
				case 2:
					w.Write(errorResponse("Automatic verification request failed. Code not ok."))
				}
			}

		} else {
			w.Write(jsonResponse(500, "Tokens could not join guild.", map[string]interface{}{}))
		}

		w.Write(jsonResponse(200, strconv.Itoa(join_result_number)+"/"+strconv.Itoa(len(core.RawTokensLoaded))+" tokens joined guild.", map[string]interface{}{}))

	} else {
		w.Write(allParametersError())
	}
}

func leaveGuild(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to leave guild...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")

	if util.AllParameters([]string{server_id}) {
		modules.StartLeaveGuildThreads(server_id)

		w.Write(jsonResponse(200, "Bots attempted to leave the target guild.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func sendFriendRequests(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to send friend requests...", 0)

	r.ParseForm()
	user_id := r.Form.Get("user_id")

	if util.AllParameters([]string{user_id}) {
		modules.StartFriendThreads(user_id)

		w.Write(jsonResponse(200, "Bots attempted to send friend requests.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func removeFriendRequests(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to remove friend requests...", 0)

	r.ParseForm()
	user_id := r.Form.Get("user_id")

	if util.AllParameters([]string{user_id}) {
		modules.StartRemoveFriendThreads(user_id)

		w.Write(jsonResponse(200, "Bots attempted to remove friend requests from target.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func speak(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting mass speak...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")
	message := r.Form.Get("message")

	if util.AllParameters([]string{server_id, message}) {
		modules.StartSpeakThreads(server_id, message)

		w.Write(jsonResponse(200, "Bots attempted to send messages in all open channels.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func startWebhookSpam(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Starting webhook spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	webhook := r.Form.Get("webhook")
	username := r.Form.Get("username")
	message := r.Form.Get("message")

	if util.AllParameters([]string{webhook, username, message}) {

		modules.StartWebhookSpam(webhook, username, message)

		w.Write(jsonResponse(200, "Attempting to start webhook spam.", map[string]interface{}{}))
	} else {
		w.Write(allParametersError())
	}
}

func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to delete webhook...", 0)

	r.ParseForm()
	webhook := r.Form.Get("webhook")

	if util.AllParameters([]string{webhook}) {

		if modules.StartWebhookDelete(webhook) {
			w.Write(jsonResponse(200, "Webhook has been deleted.", map[string]interface{}{}))
		} else {
			w.Write(errorResponse("Could not delete webhook."))
		}

	} else {
		w.Write(allParametersError())
	}
}

func disguiseTokens(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to disguise tokens...", 0)

	modules.StartDisguiseThreads()

	w.Write(jsonResponse(200, "Bots attempted to disguise.", map[string]interface{}{}))
}

func startThreadSpam(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to start thread spam...", 0)

	core.ActionFlag = 0

	r.ParseForm()
	channel_id := r.Form.Get("channel_id")
	thread_name := r.Form.Get("thread_name")

	if util.AllParameters([]string{channel_id, thread_name}) {
		modules.StartMassThreadCreateThreads(channel_id, thread_name)
	} else {
		w.Write(allParametersError())
	}
}

func fetchChannels(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	util.WriteToConsole("Attempting to fetch channels...", 0)

	r.ParseForm()
	server_id := r.Form.Get("server_id")

	channel_status_code, found_channels := modules.GetChannels(server_id)

	if channel_status_code != 200 {
		w.Write(errorResponse("An error occured when attempting to fetch guild channels. Code: " + strconv.Itoa(channel_status_code)))
		return
	}

	if len(found_channels) > 0 {
		w.Write(jsonResponse(200, "Successfully fetched guild channels.", map[string]interface{}{
			"channels": found_channels,
		}))
	} else {
		w.Write(errorResponse("No open channels available."))
	}
}

func createMassInvites(w http.ResponseWriter, r *http.Request) {
	readyRequestCors(w)

	core.ActionFlag = 0

	util.WriteToConsole("Starting create invite spam...", 0)
	r.ParseForm()

	channel_id := r.Form.Get("channel_id")

	if util.AllParameters([]string{channel_id}) {
		modules.StartMassInviteThreads(channel_id)
	} else {
		w.Write(allParametersError())
	}
}

var deadcord_banner string = `
    ██████╗ ███████╗ █████╗ ██████╗  ██████╗ ████████╗ ██████╗ ██████╗   ┏━━━━━━━━━━━━━━━━━━ Info ━━━━━━━━━━━━━━━━┓
    ██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝██████████╗██╔══██╗██╔══██╗   ` + util.Purple + `@ Package:` + util.ColorReset + ` Deadcord-Engine
    ██║  ██║█████╗  ███████║██║  ██║██║     ██║ ██  ██║██████╔╝██║  ██║   ` + util.Purple + `@ Tokens:` + util.ColorReset + ` %d tokens loaded.
    ██║  ██║██╔══╝  ██╔══██║██║  ██║██║     ████  ████║██╔══██╗██║  ██║   ` + util.Purple + `@ Warning:` + util.Red + ` Use at your own risk.` + util.ColorReset + ` 
    ██████╔╝███████╗██║  ██║██████╔╝╚██████╗╚████████╔╝██║  ██║██████╔╝   ` + util.Purple + `@ Author:` + util.ColorReset + ` https://github.com/GalaxzyDev` + util.Purple + `
    ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═════╝  ╚═════╝ █═█═█═█═╝ ╚═╝  ╚═╝╚═════╝   ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

   ` + util.White + ` You need to download our Better Discord (https://betterdiscord.app/) plugin to interact with Deadcord. Having
       trouble with Deadcord? Read our Github README and create a ticket in our community Discord or Telegram. ` + util.Purple + `

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

`

func main() {

	rand.Seed(time.Now().UnixNano())

	util.WriteToConsole("Starting error logger.", 0)

	token_status, raw_tokens, built_tokens := core.LoadTokens()
	proxy_status, raw_proxies, built_proxies := core.LoadProxies()

	if token_status {

		returned_token_amount := core.SetTokens(raw_tokens, built_tokens)

		if proxy_status {
			returned_proxy_amount := core.SetProxies(raw_proxies, built_proxies)
			util.WriteToConsole(strconv.Itoa(returned_proxy_amount)+" active proxies loaded.", 2)
		}

		banner_template := fmt.Sprintf(strings.ReplaceAll(deadcord_banner, "█", util.White+"█"+util.Purple), returned_token_amount)
		fmt.Fprintln(color.Output, banner_template)

		if len(raw_tokens) > 350 {
			util.WriteToConsole("Your token file exceeds the safely tested 350 token amount. Using this amount of tokens may lead to unexpected side-effects. Deadcord and the developers are not responsible for any loss of tokens. Continue at your own risk.", 1)
		}

		util.WriteToConsole("Starting Deadcord version: "+DeadcordVersionNumberString, 0)

		util.WriteToConsole(util.GetQuote(), 0)

		main_router := mux.NewRouter()

		api_router := main_router.PathPrefix("/deadcord/").Subrouter()
		api_router.HandleFunc("/ping-tokens", pingTokens).Methods("GET")
		api_router.HandleFunc("/start-spam", startSpam).Methods("POST")
		api_router.HandleFunc("/stop-all", stopAll).Methods("GET")
		api_router.HandleFunc("/start-typing-spam", startTypingSpam).Methods("POST")
		api_router.HandleFunc("/join-guild", joinGuild).Methods("POST")
		api_router.HandleFunc("/leave-guild", leaveGuild).Methods("POST")
		api_router.HandleFunc("/react", react).Methods("POST")
		api_router.HandleFunc("/nick", changeNick).Methods("POST")
		api_router.HandleFunc("/disguise", disguiseTokens).Methods("GET")
		api_router.HandleFunc("/friend", sendFriendRequests).Methods("POST")
		api_router.HandleFunc("/remove-friend", removeFriendRequests).Methods("POST")
		api_router.HandleFunc("/speak", speak).Methods("POST")
		api_router.HandleFunc("/start-webhook-spam", startWebhookSpam).Methods("POST")
		api_router.HandleFunc("/start-thread-spam", startThreadSpam).Methods("POST")
		api_router.HandleFunc("/delete-webhook", deleteWebhook).Methods("POST")
		api_router.HandleFunc("/fetch-channels", fetchChannels).Methods("POST")
		api_router.HandleFunc("/mass-invites", createMassInvites).Methods("POST")

		// Switch back to a static port for now, due to some issues with the client.
		open_port, err := net.Listen("tcp", ":30603")
		if err != nil {
			log.Fatal(err)
		}

		use_open_port := strings.Split(open_port.Addr().String(), ":")[3]

		util.WriteToConsole("Deadcord is ready and running on port: "+use_open_port+".", 2)

		go func() {
			util.WriteToConsole("Gateway connection started.", 2)
			core.SendOnline()
		}()

		util.CheckError(http.Serve(open_port, main_router))

	}
}
