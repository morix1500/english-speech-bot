# english-speech-bot
english-speech-bot is a Twitter bot that converts English sentences into voices.

## Description
This bot converts the English tweets of the target account into voice and tweets it.  

## Usage
You need to set the follwing environment variables.  

| name                        | description                                     |
| --------------------------- | ----------------------------------------------  |
| USER_LIST_PATH              | yaml file with target twitter id list.          | 
| TWITTER_CONSUMER_KEY        | consumer key of twitter application.            |
| TWITTER_CONSUMER_SECRET     | consumer secret key of twitter application.     |
| TWITTER_ACCESS_TOKEN        | access token of twitter application.            |
| TWITTER_ACCESS_TOKEN_SECRET | access token secret key of twitter application. |

Define target Twitter ID in yaml file.  

```
userist:
  - XXXXXX
  - YYYYYY
```

Place the image to be used in the movie.  
Image name is "logo.png".  


Execute this bot.  

```
./english-speech-bot
```

The following is a DEMO.  
<https://twitter.com/km_eng_speech>

## Requirement
In order to use english-speech-bot you need the following.  

* Twitter Application
* AWS Account
* ffmpeg

## Installation
First installation of ffmpeg is required.  
<https://github.com/FFmpeg/FFmpeg/blob/master/INSTALL.md>

Next you will install english-speech-bot.  

```
$ go get github.com/morix1500/english-speech-bot
```

### Setup policy for Amazon Polly
In order to use Amazon Polly, the following policy setting is required.  

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "getSpeech",
            "Effect": "Allow",
            "Action": [
                "polly:SynthesizeSpeech"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## License
Please see the [LICENSE](./LICENSE) file for details.  

## Author
Shota Omori(Morix)  
https://github.com/morix1500
