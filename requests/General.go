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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetCfData() (string, string) {
	req, err := http.Get("https://discord.com")

	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	if req.StatusCode == http.StatusOK {
		body_bytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		body_data := string(body_bytes)

		reg_a := regexp.MustCompile("r:'[^']*'")
		r := strings.Replace(strings.Replace(reg_a.FindStringSubmatch(body_data)[0], "r:'", "", 1), "'", "", 1)

		reg_b := regexp.MustCompile("m:'[^']*'")
		m := strings.Replace(strings.Replace(reg_b.FindStringSubmatch(body_data)[0], "m:'", "", 1), "'", "", 1)

		return r, m

	} else {
		return "", ""
	}
}

func CfbmCookieValue(previous_cookie_string string, r string, m string) string {
	var hexes []string

	for i := 0; i < 2; i++ {
		token := make([]byte, 16)
		if _, err := rand.Read(token); err != nil {
			log.Fatal(err)
		}

		if len(hexes) != 2 {
			hex := string(hex.EncodeToString(token))
			hexes = append(hexes, hex)
		}
	}

	cfbm_payload := constants.CfbmPayload{
		M: m,
		Results: []interface{}{
			hexes[0],
			hexes[1],
		},
		Timing: rand.Intn(120-40) + 40,
		Fp: constants.Fp{
			ID: 3,
			E: constants.E{
				R:  []int{1920, 1080},
				Ar: []int{1040, 1920},
				Pr: 1,
				Cd: 24,
				Wb: false,
				Wp: false,
				Wn: false,
				Ch: true,
				Ws: false,
				Wd: true,
			},
		},
	}

	cookie_client := http.Client{}
	cookie_json, err := json.Marshal(cfbm_payload)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://discord.com/cdn-cgi/bm/cv/result?req_id="+r, bytes.NewBuffer(cookie_json))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := cookie_client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return resp.Cookies()[0].Value
}

func GetDiscordCookies() string {
	r, m := GetCfData()

	req, err := http.Get("https://discord.com")

	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	cookie_string_one := "__dcfduid=" + req.Cookies()[0].Value + "; __sdcfduid=" + req.Cookies()[1].Value + "; locale=en-GB;"
	cfbm_cookie_string := CfbmCookieValue(cookie_string_one, r, m)
	return cookie_string_one + "__cf_bm" + cfbm_cookie_string

}

func GetDiscordFingerprint() string {
	resp, err := http.Get("https://discord.com/api/v9/experiments")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var science constants.Science
	if err := json.Unmarshal(body, &science); err != nil {
		log.Fatal(err)
	}

	return science.Fingerprint
}

func GetNonce() int64 {
	nonce_raw := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(nonce_raw, 2, 64)
	return nonce
}

func SimpleGet(url string) (int, []byte) {
	req, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
	}

	return req.StatusCode, body
}
