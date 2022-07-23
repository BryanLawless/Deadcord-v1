/*
 * Copyright (C) 2022 GalaxzyDev - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the AGPL-3.0 License.
 *
 * You should have received a copy of the AGPL-3.0 License with
 * this file. If not, please refer to the license linked below.
 * https://github.com/GalaxzyDev/Deadcord/blob/main/LICENSE.md
 */

package core

import (
	"math/rand"
	"sync"
)

var (
	wg               sync.WaitGroup
	ActionFlag       = 0
	DeadcordVersion  = 1.4
	RawTokensLoaded  []string
	Profiles         map[int]map[string]string
	BuiltTokenStruct map[int]map[string]string

	RawProxiesLoaded []string
	BuiltProxyStruct map[int]map[string]string

	TokensInServers map[string][]string
)

func SetTokens(raw_token_list []string, built_token_list map[int]map[string]string) int {
	RawTokensLoaded = raw_token_list
	BuiltTokenStruct = built_token_list

	return len(RawTokensLoaded)
}

func SetProxies(raw_proxy_list []string, built_proxy_list map[int]map[string]string) int {
	RawProxiesLoaded = raw_proxy_list
	BuiltProxyStruct = built_proxy_list

	return len(RawProxiesLoaded)
}

func RandomToken() string {
	random_token := RawTokensLoaded[rand.Intn(len(RawTokensLoaded))]
	return random_token
}

func RandomTokenInServer(server_id string) string {
	random_token_in_server := TokensInServers[server_id][rand.Intn(len(TokensInServers))]

	return random_token_in_server
}
