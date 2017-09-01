# Twitch Demo
Simple demo for Twitch API access

Provides the following information about a Twitch user, given their username

* user's channel's **# of views**
* user's channel's **# of followers**
* user's channel's **game**
* user's channel's **language**
* if the user is **currently streaming**
* user's **display name**
* user's **bio**
* user's **account creation date**

## Setup
* Acquire a Client ID from Twitch for API access
* Put your Client ID in a file called 'twitchClientId' in the base of the repo (it should not contain any whitespace)
* Make sure you have Go installed
* Make sure port 8080 is not already in use
* Start the server with 'go run main.go'

## Frontend Demo View
Once the server is running, it will serve a simple webpage on [http://localhost:8080](http://localhost:8080) which provides a human-friendly interface for testing the APIs. Enter a Twitch username and click the buttons to fetch information and display it in human-friendly formatting.

## Backend APIs
The server exposes APIs under http://localhost:8080/api/. All endpoints provide JSON data for status 200 responses, and will use different status codes for errors.

* GET /api/user?name=<username> -> returns information about a Twitch user
* GET /api/channel?name=<username> -> returns information about a Twitch user's channel
* GET /api/stream?name=<username> -> returns information about a Twitch user's current live stream, if they are broadcasting (otherwise 'null')

## Example Request/Response
Request

    GET http://localhost:8080/api/user?name=sp4zie

Response

    //
