package main

import "fmt"
import "regexp"
import "os"
import "log"
import "net/url"
import "github.com/joho/godotenv"
import anaconda "github.com/ChimeraCoder/anaconda"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("ConsumerKey"))
	anaconda.SetConsumerSecret(os.Getenv("ConsumerSecret"))
	api := anaconda.NewTwitterApi(os.Getenv("AccessToken"), os.Getenv("AccessTokenSecret"))

	v := url.Values{}
	v.Set("count", "30")

	since_id := "0"
	barf := url.Values{}
	timeline, err := api.GetHomeTimeline(barf)

	for _, tweet := range timeline {
		if tweet.IdStr > since_id {
			since_id = tweet.IdStr
		}
	}

	v.Set("since_id", since_id)
	name := os.Getenv("Name")

	coolWords := []string{"awesome", "great", "rad", "cool", "really cool", "the best"}
	searchTerms := ""

	for index, word := range coolWords {
		if index != 0 {
			searchTerms = searchTerms + " OR "
		}
		searchTerms = searchTerms + " \"" + name + " is " + word + "\" "
	}

	searchTerms = "-RT " + searchTerms

	fmt.Println(searchTerms)

	searchResult, err := api.GetSearch(searchTerms, v)
	if err != nil {
		fmt.Println(err)
	}

	matches := make(map[string]anaconda.Tweet)
	var maxId int64

	for _, tweet := range searchResult.Statuses {
		isRT, _ := regexp.MatchString("RT |Marano", tweet.Text)
		maxId = tweet.Id

		if !isRT {
			matches[tweet.Text] = tweet
		}
	}

	for text, tweet := range matches {
		//_, err := api.Retweet(tweet.Id, false)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(text, tweet.Id, tweet.FilterLevel, tweet.PossiblySensitive, tweet.RetweetedStatus)
	}
	fmt.Println(maxId)
}
