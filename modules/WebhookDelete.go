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
	"Deadcord/requests"
	"net/http"
)

func StartWebhookDelete(webhook string) bool {
	webhook_headers := http.Header{"Content-type": []string{"application/json"}}

	status, status_code, _ := requests.RequestTemplate(http.Client{}, "DELETE", webhook, webhook_headers, map[string]interface{}{})

	if status && status_code == 204 {
		return true
	}

	return true
}
