/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package requests

import (
	"Deadcord/constants"
	"Deadcord/core"
	"Deadcord/util"
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var (
	CookieString      = GetDiscordCookies()
	GlobalFingerprint = GetDiscordFingerprint()
	BaseURLs          = []string{"https://discord.com/api/v9/", "https://ptb.discord.com/api/v9/", "https://canary.discord.com/api/v9/"}
)

func SendDiscordRequest(endpoint string, method string, token string, header_key string, custom_context map[string]interface{}, data map[string]interface{}) (bool, int, []byte) {

	build_x_prop := func(x_prop_data map[string]interface{}) string {
		x_prop_json, err := json.Marshal(custom_context)
		util.CheckError(err)

		x_prop_refined := b64.StdEncoding.EncodeToString(x_prop_json)

		return x_prop_refined
	}

	used_base_url := BaseURLs[rand.Intn(len(BaseURLs))]

	discord_client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			},
			DisableKeepAlives: true,
			ForceAttemptHTTP2: true,
		},
	}

	release_channel := "stable"
	if strings.Contains(used_base_url, "ptb") {
		release_channel = "ptb"
	}

	if strings.Contains(used_base_url, "canary") {
		release_channel = "canary"
	}

	current_profile := util.ReadTokenProfile(core.Profiles, token)

	x_super_props := map[string]interface{}{
		"os":                       current_profile["os"],
		"browser":                  current_profile["browser"],
		"device":                   "",
		"system_locale":            "en-US",
		"browser_user_agent":       current_profile["agent"],
		"browser_version":          current_profile["browser_version"],
		"os_version":               current_profile["os_version"],
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          release_channel,
		"client_build_number":      130153,
		"client_event_source":      "None",
	}

	x_track_props := map[string]interface{}{
		"os":                       current_profile["os"],
		"browser":                  current_profile["browser"],
		"device":                   "",
		"system_locale":            "en-US",
		"browser_user_agent":       current_profile["agent"],
		"browser_version":          current_profile["browser_version"],
		"os_version":               current_profile["os_version"],
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          release_channel,
		"client_build_number":      "9999",
		"client_event_source":      "None",
	}

	x_super := build_x_prop(x_super_props)
	x_track := build_x_prop(x_track_props)
	x_context := build_x_prop(custom_context)

	discord_headers := http.Header{}

	switch header_key {
	case "all":
		discord_headers = constants.AllHeaders(CookieString, x_super, x_context, x_track, GlobalFingerprint, current_profile)
	case "general":
		discord_headers = constants.GeneralHeaders(CookieString, x_super, current_profile)
	case "join":
		discord_headers = constants.JoinGuildHeaders(CookieString, x_super, x_context, current_profile)
	case "message":
		discord_headers = constants.SendMessageHeaders(CookieString, x_super, current_profile)
	default:
		discord_headers = constants.GeneralHeaders(CookieString, x_super, current_profile)
	}

	switch method {
	case "GET":
		status, status_code, body := GetRequestTemplate(discord_client, used_base_url+endpoint, discord_headers)
		return status, status_code, body
	case "POST":
		status, status_code, body := RequestTemplate(discord_client, "POST", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "PUT":
		status, status_code, body := RequestTemplate(discord_client, "PUT", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "PATCH":
		status, status_code, body := RequestTemplate(discord_client, "PATCH", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	case "DELETE":
		status, status_code, body := RequestTemplate(discord_client, "DELETE", used_base_url+endpoint, discord_headers, data)
		return status, status_code, body
	}

	return false, 0, nil
}

func RequestTemplate(client http.Client, request_type string, url string, headers http.Header, json_payload map[string]interface{}) (bool, int, []byte) {

	patch_json, err := json.Marshal(json_payload)
	util.CheckError(err)

	req, err := http.NewRequest(request_type, url, bytes.NewBuffer(patch_json))
	util.CheckError(err)

	req.Close = true

	req.Header = headers

	res, err := client.Do(req)
	util.CheckError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	util.CheckError(err)

	return true, res.StatusCode, []byte(body)
}

func GetRequestTemplate(client http.Client, url string, headers http.Header) (bool, int, []byte) {

	req, err := http.NewRequest("GET", url, nil)
	util.CheckError(err)

	req.Close = true

	req.Header = headers

	res, err := client.Do(req)
	util.CheckError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	util.CheckError(err)

	return true, res.StatusCode, []byte(body)
}
