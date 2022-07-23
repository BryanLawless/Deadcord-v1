/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 *
 * Thanks to Zeveryone for this boilerplate websocket code to help me futher understand the gorilla websocket module.
 *
 */

package core

import (
	"Deadcord/constants"
	"Deadcord/util"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

var (
	GatewayConnection bool = false
)

func DisconnectFromGateway() {
	GatewayConnection = false
}

func StartWebsocketConnection(token string) *websocket.Conn {
	GatewayConnection = true

	ws := constants.WebsocketConnection{
		Connection: &websocket.Conn{},
	}

	var dialErr error
	ws.Connection, _, dialErr = websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=9&encoding=json", nil)
	if dialErr != nil {
		return ws.Connection
	}
	return ws.Connection
}

func showIdentification(ws *websocket.Conn, token string) error {
	current_profile := util.ReadTokenProfile(Profiles, token)

	identification := constants.DiscordGatewayPayload{
		Opcode: 2,
		EventData: constants.DiscordGatewayEventDataIdentify{
			Token:   token,
			Intents: 32767,
			Properties: map[string]interface{}{
				"$os":      current_profile["os"],
				"$browser": current_profile["browser"],
				"$device":  "pc",
			},
		},
		EventName: "IDENTIFY",
	}

	b, marshalErr := json.Marshal(identification)

	if marshalErr != nil {
		return marshalErr
	}

	err := ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

func RecieveIncomingPayloads(ws *websocket.Conn, token string) error {

	for GatewayConnection {

		_, p, readErr := ws.ReadMessage()
		if readErr != nil {
			return readErr
		}

		var decoded constants.DiscordGatewayPayload
		err := json.Unmarshal(p, &decoded)
		if err != nil {
			return err
		}

		if decoded.Opcode == 10 {
			data := decoded.EventData.(map[string]interface{})
			heartbeatInterval := data["heartbeat_interval"].(float64)

			go setupHeartbeat(heartbeatInterval, ws)
			showIdentification(ws, token)
		}

		if decoded.Opcode == 9 {
			util.WriteToConsole("Token could not connect to gateway. Most likely due to duplicate connection.", 3)
			return nil
		}
	}

	return nil
}

func setupHeartbeat(interval float64, ws *websocket.Conn) error {
	c := time.Tick(time.Duration(interval) * time.Millisecond)
	for range c {
		b, marshalErr := json.Marshal(constants.DiscordGatewayPayload{
			Opcode:    1,
			EventData: nil,
			EventName: "",
		})

		if marshalErr != nil {
			return marshalErr
		}

		err := ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendOnline() {
	wg.Add(len(RawTokensLoaded))

	for _, token := range RawTokensLoaded {
		go func(token string) {
			defer wg.Done()
			ws := StartWebsocketConnection(token)
			err := RecieveIncomingPayloads(ws, token)
			util.CheckError(err)

		}(token)
	}
	wg.Wait()
}
