// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"context"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// class that represents data about a the weather on a city.
/* Sample JSON
{
"coord": {
	"lon": -79.9,
	"lat": -2.1667
},
"weather": [
	{
	"id": 804,
	"main": "Clouds",
	"description": "overcast clouds",
	"icon": "04n"
	}
],
"base": "stations",
"main": {
	"temp": 296.31,
	"feels_like": 296.9,
	"temp_min": 295.9,
	"temp_max": 296.31,
	"pressure": 1012,
	"humidity": 85
},
"visibility": 10000,
"wind": {
	"speed": 2.68,
	"deg": 285,
	"gust": 8.05
},
"clouds": {
	"all": 90
},
"dt": 1629516875,
"sys": {
	"type": 2,
	"id": 2008064,
	"country": "EC",
	"sunrise": 1629458495,
	"sunset": 1629501876
},
"timezone": -18000,
"id": 3657509,
"name": "Guayaquil",
"cod": 200
}
*/
type CityWeather struct {
	coord      Coordinate `json:"coord"`
	weather    []Weather  `json:"weather"`
	base       string     `json:"base"`
	main       Main       `json:"main"`
	visibility int64      `json:"visibility"`
	wind       Wind       `json:"wind"`
	clouds     Clouds     `json:"clouds"`
	dt         int64      `json:"dt"`
	sys        Sys        `json:"sys"`
	timezone   int64      `json:"timezone"`
	id         int64      `json:"id"`
	name       string     `json:"name"`
	cod        int64      `json:"cod"`
}

type Coordinate struct {
	lon float64 `json:"lon"`
	lat float64 `json:"lat"`
}

type Weather struct {
	id          int64  `json:"id"`
	main        string `json:"main"`
	description string `json:"description"`
	icon        string `json:"icon"`
}

type Main struct {
	temp       float64 `json:"temp"`
	feels_like float64 `json:"feels_like"`
	temp_min   float64 `json:"temp_min"`
	temp_max   float64 `json:"temp_max"`
	pressure   int64   `json:"pressure"`
	humidity   int64   `json:"humidity"`
}

type Wind struct {
	speed float64 `json:"speed"`
	deg   int64   `json:"deg"`
	gust  float64 `json:"gust"`
}

type Clouds struct {
	all int64 `json:"all"`
}

type Sys struct {
	id       int64  `json:"id"`
	sys_type int64  `json:"type"`
	country  string `json:"country"`
	sunrise  int64  `json:"sunrise"`
	sunset   int64  `json:"sunset"`
}

const redirectURI = "http://localhost:9001/callback"

var (
	auth   = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate), spotifyauth.WithClientID(viperEnvVariable("SPOTIFY_ID")), spotifyauth.WithClientSecret(viperEnvVariable("SPOTIFY_SECRET")))
	client *spotify.Client
	state  = "abc123"
)

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string) string {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {

	os.Setenv("SPOTIFY_ID", viperEnvVariable("SPOTIFY_ID"))
	os.Setenv("SPOTIFY_SECRET", viperEnvVariable("SPOTIFY_SECRET"))

	//redirectURI := viperEnvVariable("SPOTIFY_REDIRECT_URI")

	//auth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))

	/*

		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client := <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)

		// search for playlists and albums containing "holiday"
		results, err := client.Search("holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
		if err != nil {
			log.Fatal(err)
		}

		// handle album results
		if results.Albums != nil {
			fmt.Println("Albums:")
			for _, item := range results.Albums.Albums {
				fmt.Println("   ", item.Name)
			}
		}
		// handle playlist results
		if results.Playlists != nil {
			fmt.Println("Playlists:")
			for _, item := range results.Playlists.Playlists {
				fmt.Println("   ", item.Name)
			}
		}
	*/

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	r.GET("/callback", func(c *gin.Context) {
		completeAuth(c.Writer, c.Request)
		c.Redirect(http.StatusMovedPermanently, "/user")
	})
	r.GET("/login", func(c *gin.Context) {
		fmt.Println("Got login request")
		//construir la URL
		url := auth.AuthURL(state)
		fmt.Println(url)
		c.Redirect(http.StatusMovedPermanently, url)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/user", func(c *gin.Context) {
		// use the client to make calls that require authorization
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
		fmt.Println("User ID:", user.ID)
		fmt.Println("Display name:", user.DisplayName)
		fmt.Println("Spotify URI:", string(user.URI))
		fmt.Println("Endpoint:", user.Endpoint)
		fmt.Println("Followers:", user.Followers.Count)
		c.HTML(http.StatusOK, "user.html", gin.H{
			"user": user,
		})
	})
	r.GET("/climate", func(c *gin.Context) {
		resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=Guayaquil&appid=1cdbcd14a6e201f2b5d091e4b1c53b8a")
		if err != nil {
			log.Fatalln(err)
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//Convert the body to type string
		sb := string(body)
		c.IndentedJSON(http.StatusOK, sb)
	})
	r.Run(":9001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback handler for state:", state)
	fmt.Println(r.FormValue("state"))

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client = spotify.New(auth.Client(r.Context(), tok))
	fmt.Println("Client a!")
}
