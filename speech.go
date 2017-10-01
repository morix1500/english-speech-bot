package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/pkg/errors"
	"io"
	"os"
)

func GetSpeech(text string, output_path string) error {
	text = fmt.Sprintf("<speak xml:lang=\"en-US\"><prosody rate=\"default\">%s</prosody><break time=\"1s\" /><prosody rate=\"x-slow\">%s</prosody></speak>", text, text)

	region := "us-east-1"

	conf := &aws.Config{
		Region: &region,
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		return errors.Wrap(err, "Faild create new session.")
	}

	svc := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		SampleRate:   aws.String("8000"),
		Text:         aws.String(text),
		TextType:     aws.String("ssml"),
		VoiceId:      aws.String("Salli"),
	}
	result, err := svc.SynthesizeSpeech(input)
	if err != nil {
		return errors.Wrap(err, "Failed get aws polly.")
	}

	wf, err := os.OpenFile(output_path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "Failed file open.")
	}
	defer wf.Close()

	io.Copy(wf, result.AudioStream)

	return nil
}
