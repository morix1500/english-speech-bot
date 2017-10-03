package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var logger = logrus.New()

type UserIdList struct {
	Ids []string `yaml:"userlist"`
}

const (
	image_path         string = "logo.png"
)

func loadConfig() (string, error) {
	fp, err := ioutil.ReadFile(os.Getenv("USER_LIST_PATH"))
	if err != nil {
		return "", errors.Wrap(err, "Faild user list config.")
	}

	u := UserIdList{}
	err = yaml.Unmarshal(fp, &u)
	if err != nil {
		return "", errors.Wrap(err, "Faild parse from yaml.")
	}

	return strings.Join(u.Ids, ","), nil
}

func uploadTweet(a *anaconda.TwitterApi, video string, speech string, text string) error {
	defer func() {
		os.Remove(speech)
		os.Remove(video)
	}()

	fp, err := ioutil.ReadFile(video)
	if err != nil {
		return errors.Wrap(err, "Faild video open.")
	}

	total_bytes := len(fp)
	chanked_media, err := a.UploadVideoInit(total_bytes, "video/mp4")
	if err != nil {
		return errors.Wrap(err, "Faild video init upload.")
	}

	media_max_len := 5 * 1024 * 1024
	segment := 0

	for i := 0; i < total_bytes; i += media_max_len {
		var media_data string
		if i+media_max_len < total_bytes {
			media_data = base64.StdEncoding.EncodeToString(fp[i : i+media_max_len])
		} else {
			media_data = base64.StdEncoding.EncodeToString(fp[i:])
		}
		if err = a.UploadVideoAppend(chanked_media.MediaIDString, segment, media_data); err != nil {
			break
		}
		segment += 1
	}

	video_media, err := a.UploadVideoFinalize(chanked_media.MediaIDString)
	if err != nil {
		return errors.Wrap(err, "Faild video upload.")
	}

	params := url.Values{}
	params.Add("media_ids", video_media.MediaIDString)
	_, err = a.PostTweet(text, params)
	if err != nil {
		return errors.Wrap(err, "Faild post of tweet.")
	}
	logger.Info("Post success.")
	return nil
}

func createVideo(tweet, speech_output_path, video_output_path string) error {
	err := GetSpeech(tweet, speech_output_path)
	if err != nil {
		return err
	}
	err = EncodeVideo(image_path, speech_output_path, video_output_path)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	logger.Out = os.Stdout

	user_ids, err := loadConfig()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	v := url.Values{}
	v.Set("follow", user_ids)

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	stream := api.PublicStreamFilter(v)
	logger.Info("Start stream...")
	for {
		select {
		case stream := <-stream.C:
			switch status := stream.(type) {
			case anaconda.Tweet:
				if status.RetweetedStatus == nil && status.InReplyToStatusID == 0 {
					tweet := status.Text
					tweet = strings.Split(tweet, "\n")[0]
					logger.Info("Tweet: " + tweet)

					base := fmt.Sprintf("%x", md5.Sum([]byte(tweet)))
					speech_output_path := base + ".mp3"
					video_output_path := base + ".mp4"

					err := createVideo(tweet, speech_output_path, video_output_path)
					if err != nil {
						logger.Fatal(err)
					}

					_, err = api.Retweet(status.Id, true)
					if err != nil {
						logger.Fatal(err)
					}

					err = uploadTweet(api, video_output_path, speech_output_path, tweet)
					if err != nil {
						logger.Fatal(err)
					}
				}
			}
		}
	}
}
