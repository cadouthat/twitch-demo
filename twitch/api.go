/*
Utilities for accessing Twitch APIs
by: Connor Douthat
9/1/2017
*/
package twitch

import (
	"log"
	"errors"
	"time"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

const apiBase = "https://api.twitch.tv/kraken/"
const apiAccept = "application/vnd.twitchtv.v5+json"
const clientIdFile = "twitchClientId"
var clientId = ""

type CacheEntry struct {
	key string
	expires time.Time
	value interface{}
}

// The time (in seconds) it will take for each type of cache entry to expire
const userCacheTTL = 60
const channelCacheTTL = 60
const streamCacheTTL = 5

// TODO: These basic maps demonstrate caching, but in production we would
//	need something with eviction so we can cap the memory usage (e.g. LRU cache)
var userCache = make(map[string]CacheEntry)
var channelCache = make(map[string]CacheEntry)
var streamCache = make(map[string]CacheEntry)

func init() {
	clientIdBytes, err := ioutil.ReadFile(clientIdFile)
	if err != nil {
		log.Fatal("Failed to read Client ID")
	}
	clientId = string(clientIdBytes)
}

func cacheFetch(cache map[string]CacheEntry, key string) interface{} {
	entry, ok := cache[key]

	if !ok {
		return nil
	}

	if entry.expires.Before(time.Now()) {
		return nil
	}

	return entry.value
}

func cachePut(cache map[string]CacheEntry, key string, value interface{}, ttl int) {
	entry := CacheEntry{}
	entry.key = key
	entry.expires = time.Now().Add(time.Duration(ttl) * time.Second)
	entry.value = value
	cache[key] = entry
}

func doGetAPI(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", apiAccept)
	req.Header.Add("Client-ID", clientId)

	//
	log.Println(uri)
	//

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type TwitchUser struct {
	Id string `json:"_id"`
	DisplayName string `json:"display_name"`
	Bio string `json:"bio"`
	CreatedAt string `json:"created_at"`
}

type TwitchUsersResponse struct {
	Users []TwitchUser `json:"users"`
}

func GetUserByName(name string) (*TwitchUser, error) {
	cached := cacheFetch(userCache, name)
	if cached != nil {
		return cached.(*TwitchUser), nil
	}

	uriName := url.QueryEscape(name)
	uri := apiBase + "users?login=" + uriName

	body, err := doGetAPI(uri)
	if err != nil {
		return nil, err
	}

	response := &TwitchUsersResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if len(response.Users) != 1 {
		return nil, errors.New("No unique user found")
	}

	user := &response.Users[0]
	cachePut(userCache, name, user, userCacheTTL)
	return user, nil
}

type TwitchChannel struct {
   Game string `json:"game"`
   Language string `json:"language"`
   Views int `json:"views"`
   Followers int `json:"followers"`
}

func GetChannelByUser(userId string) (*TwitchChannel, error) {
	cached := cacheFetch(channelCache, userId)
	if cached != nil {
		return cached.(*TwitchChannel), nil
	}

	uriId := url.QueryEscape(userId)
	uri := apiBase + "channels/" + uriId

	body, err := doGetAPI(uri)
	if err != nil {
		return nil, err
	}

	response := &TwitchChannel{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	cachePut(channelCache, userId, response, channelCacheTTL)
	return response, nil
}

type TwitchStream struct {
	Viewers int `json:"viewers"`
}

type TwitchStreamResponse struct {
	Stream *TwitchStream `json:"stream"`
}

func GetStreamByUser(userId string) (*TwitchStream, error) {
	cached := cacheFetch(streamCache, userId)
	if cached != nil {
		return cached.(*TwitchStream), nil
	}

	uriId := url.QueryEscape(userId)
	uri := apiBase + "streams/" + uriId

	body, err := doGetAPI(uri)
	if err != nil {
		return nil, err
	}

	response := &TwitchStreamResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	cachePut(streamCache, userId, response.Stream, streamCacheTTL)
	return response.Stream, nil
}
