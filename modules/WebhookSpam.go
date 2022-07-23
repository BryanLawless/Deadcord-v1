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
	"net/http"
)

func StartWebhookSpam(webhook string, username string, message string) {
	go func(webhook string) {
		for {

			if core.ActionFlag == 1 {
				return
			}

			webhook_headers := http.Header{"Content-type": []string{"application/json"}}

			status, status_code, webhook_json := requests.RequestTemplate(http.Client{}, "POST", webhook, webhook_headers, map[string]interface{}{
				"username": username,
				"content":  message,
			})

			if status {
				switch status_code {
				case 429:
					var rate_limit constants.RateLimit
					err := json.Unmarshal(webhook_json, &rate_limit)
					util.CheckError(err)

					pause_time := int(rate_limit.RetryAfter) / 1000

					if pause_time > 0 {
						util.WriteToConsole("Webhook rate limited, pausing for: "+fmt.Sprintf("%d", pause_time)+" seconds.", 1)
					}

					util.Sleep(int(rate_limit.RetryAfter) / 1000)
				case 404:
					util.WriteToConsole("Webhook no longer exists.", 3)
					return
				case 401:
					util.WriteToConsole("An error occured while trying to spam webhook.", 3)
					return
				}
			}
		}
	}(webhook)
}
