package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/yagop/jumble-bot/config"
	"gopkg.in/telegram-bot-api.v4"
)

func Fetch(user, repo, authUser, authPassword string) (<-chan Commit, error) {
	channel := make(chan Commit, 10)
	timeInterval := time.Second * 3

	client := &http.Client{}
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commits?pagelen=1", user, repo)

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(authUser, authPassword)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("StatusCode %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var commits Commits
	err = json.Unmarshal(body, &commits)

	if err != nil {
		fmt.Println(err)
		fmt.Println(body)
		return nil, err
	}

	// Take last commit
	last_commit := commits.Values[0]
	eTag := resp.Header.Get("ETag")

	fmt.Println("ETag:", eTag)
	fmt.Println("Last commit:", last_commit.Hash)

	go func() {
		// for range []byte{1} {
		for {
			req, err := http.NewRequest("GET", url, nil)
			// Conditional request
			req.Header.Add("If-None-Match", eTag)
			req.SetBasicAuth(authUser, authPassword)
			resp, err := client.Do(req)

			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				if err == nil {
					if resp.StatusCode == 200 {
						eTag = resp.Header.Get("ETag")
						fmt.Println("ETag:", eTag)
						err = json.Unmarshal(body, &commits)
						if err != nil {
							fmt.Println(err)
							fmt.Println(body)
						} else {
							// New commit
							fmt.Println("Last commit:", last_commit.Hash)
							if commits.Values[0].Hash != last_commit.Hash {
								last_commit = commits.Values[0]
								channel <- last_commit
							}
						}
					} else if resp.StatusCode != 304 { // Not Modified
						fmt.Println("StatusCode:", resp.StatusCode)
					}
				}
			} else {
				fmt.Println(err)
			}
			time.Sleep(timeInterval)
		}
		close(channel)
	}()
	return channel, nil
}

func BitBucket(config *config.TomlConfig) (<-chan tgbotapi.MessageConfig, error) {
	message_channel := make(chan tgbotapi.MessageConfig, 10)
	user := config.BitBucketRepoUser
	repo := config.BitBucketRepo
	authUser := config.BitBucketUser
	authPassword := config.BitBucketPassword
	telegramChat := config.BitBucketTelegramChat

	commit_channel, err := Fetch(user, repo, authUser, authPassword)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	go func() {
		for commit := range commit_channel {
			fmt.Println("Commit:", commit)
			commit_url := fmt.Sprintf("https://bitbucket.org/%s/commits/%s", commit.Repository.FullName, commit.Hash)
			text := fmt.Sprintf("New commit on *%s* by %s\n[%s](%s): `%s`", commit.Repository.FullName, commit.Author.User.DisplayName, commit.Hash[0:7], commit_url, commit.Message)
			message := tgbotapi.NewMessage(telegramChat, text)
			message.ParseMode = "markdown"
			message_channel <- message
		}
	}()
	return message_channel, nil
}
