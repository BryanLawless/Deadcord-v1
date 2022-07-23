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
	b64 "encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

var image_urls = [2]string{
	"https://picsum.photos/512/512",
	"https://cataas.com/cat?width=512&height=512",
}

func StartDisguiseThreads() {

	wg.Add(len(core.RawTokensLoaded))

	for _, token := range core.RawTokensLoaded {
		go func(token string) {
			defer wg.Done()
			disguiseWorker(token)
		}(token)
	}

	wg.Wait()
}

func disguiseWorker(token string) {
	random_image_api := image_urls[rand.Intn(len(image_urls))]

	resp, err := http.Get(random_image_api)
	util.CheckError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckError(err)

	image := "data:image/png;base64," + b64.StdEncoding.EncodeToString(body)

	status, status_code, _ := requests.SendDiscordRequest("users/@me", "PATCH", token, "general", map[string]interface{}{}, map[string]interface{}{
		"avatar": image,
	})

	if status {
		switch status_code {
		case 200:
			util.WriteToConsole("Successfully changed bot avatar.", 2)
		case 429:
			util.WriteToConsole("IP ratelimited or Cloudflare banned.", 1)
		case 400:
			util.WriteToConsole("Bot could not change avatar, most likely due to rate-limits.", 3)
		default:
			util.WriteToConsole("Bot could not change avatar, request failed. Code: "+strconv.Itoa(status_code), 3)
		}
	}
}
