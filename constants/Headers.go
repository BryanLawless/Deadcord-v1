/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package constants

import (
	"net/http"
)

var (
	BasicUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9004 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36"
)

func SimpleHeaders(x_super string) http.Header {
	simple_headers := http.Header{
		"Accept-Language":    []string{"en-US"},
		"Content-Type":       []string{"application/json"},
		"User-Agent":         []string{BasicUserAgent},
		"X-Super-Properties": []string{x_super},
		"X-Debug-Options":    []string{"bugReporterEnabled"},
		"X-Discord-Locale":   []string{"en-US"},
		"Sec-Fetch-Site":     []string{"same-origin"},
		"Sec-Fetch-Mode":     []string{"cors"},
		"Sec-Fetch-Dest":     []string{"empty"},
		"TE":                 []string{"trailers"},
	}

	return simple_headers
}

func GeneralHeaders(cookie_string string, x_super_properties string, current_profile map[string]string) http.Header {
	general_headers := http.Header{
		"Host":               []string{"discord.com"},
		"Accept":             []string{"*/*"},
		"Accept-Language":    []string{"en-US,en;q=0.5"},
		"Content-Type":       []string{"application/json"},
		"Authorization":      []string{current_profile["token"]},
		"User-Agent":         []string{current_profile["agent"]},
		"X-Super-Properties": []string{x_super_properties},
		"X-Discord-Locale":   []string{"en-US"},
		"X-Debug-Options":    []string{"bugReporterEnabled"},
		"Cookie":             []string{cookie_string},
		"TE":                 []string{"trailers"},
	}

	return general_headers
}

func GeneralContextHeaders(cookie_string string, x_super_properties string, x_context_properties string, current_profile map[string]string) http.Header {
	general_context_headers := http.Header{
		"Host":                 []string{"discord.com"},
		"Accept":               []string{"*/*"},
		"Accept-Language":      []string{"en-US,en;q=0.5"},
		"Authorization":        []string{current_profile["token"]},
		"X-Context-Properties": []string{x_context_properties},
		"X-Super-Properties":   []string{x_super_properties},
		"X-Discord-Locale":     []string{"en-US"},
		"X-Debug-Options":      []string{"bugReporterEnabled"},
		"Cookie":               []string{cookie_string},
		"User-Agent":           []string{current_profile["agent"]},
	}

	return general_context_headers
}

func JoinGuildHeaders(cookie_string string, x_super_properties string, x_context_properties string, current_profile map[string]string) http.Header {
	join_guild_headers := http.Header{
		"Host":                 []string{"discord.com"},
		"Accept":               []string{"*/*"},
		"Accept-Language":      []string{"en-US,en;q=0.5"},
		"Content-Type":         []string{"application/json"},
		"Authorization":        []string{current_profile["token"]},
		"User-Agent":           []string{current_profile["agent"]},
		"X-Context-Properties": []string{x_context_properties},
		"X-Super-Properties":   []string{x_super_properties},
		"X-Discord-Locale":     []string{"en-US"},
		"X-Debug-Options":      []string{"bugReporterEnabled"},
		"Cookie":               []string{cookie_string},
	}

	return join_guild_headers
}

func SendMessageHeaders(cookie_string string, x_super_properties string, current_profile map[string]string) http.Header {
	send_message_headers := http.Header{
		"Host":               []string{"discord.com"},
		"Accept":             []string{"*/*"},
		"Accept-Language":    []string{"en-US,en;q=0.5"},
		"Content-Type":       []string{"application/json"},
		"User-Agent":         []string{current_profile["agent"]},
		"Origin":             []string{"https://discord.com"},
		"DNT":                []string{"1"},
		"Referer":            []string{"https://discord.com/channels/@me"},
		"Sec-Fetch-Dest":     []string{"empty"},
		"Sec-Fetch-Mode":     []string{"cors"},
		"Sec-Fetch-Site":     []string{"same-origin"},
		"Authorization":      []string{current_profile["token"]},
		"X-Super-Properties": []string{x_super_properties},
		"X-Discord-Locale":   []string{"en-US"},
		"X-Debug-Options":    []string{"bugReporterEnabled"},
		"Cookie":             []string{cookie_string},
	}

	return send_message_headers
}

func AllHeaders(cookie_string string, x_super_properties string, x_context_properties string, x_track_properties string, fingerprint string, current_profile map[string]string) http.Header {
	all_headers := http.Header{
		"Host":               []string{"discord.com"},
		"Accept":             []string{"*/*"},
		"Accept-language":    []string{"en-GB"},
		"Authorization":      []string{current_profile["token"]},
		"Alt-Used":           []string{"discord.com"},
		"Content-type":       []string{"application/json"},
		"Cookie":             []string{cookie_string},
		"DNT":                []string{"1"},
		"Sec-ch-ua":          []string{"Not A;Brand';v='99', 'Chromium';v='96', 'Google Chrome';v='96'"},
		"Sec-ch-ua-mobile":   []string{"0"},
		"Sec-ch-ua-platform": []string{current_profile["os"]},
		"TE":                 []string{"Trailers"},
		"User-Agent":         []string{current_profile["agent"]},
		"X-Debug-options":    []string{"bugReporterEnabled"},
		"X-Discord-Locale":   []string{"en-US"},
		"X-Track":            []string{x_track_properties},
		"X-Fingerprint":      []string{fingerprint},
		"X-Super-Properties": []string{x_super_properties},
	}

	return all_headers
}
