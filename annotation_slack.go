package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func doAnnotateWithSlack(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	outputFilename := c.String("output-filename")
	channelID := c.String("channel")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "slack")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	if outputFilename == "" {
		_ = cli.ShowCommandHelp(c, "slack")
		return cli.NewExitError("`output-filename` is a required field.", 1)
	}

	if channelID == "" {
		_ = cli.ShowCommandHelp(c, "slack")
		return cli.NewExitError("`channel` is a required field.", 1)
	}

	cacheFilename := CacheFilename
	cache, err := LoadCache(cacheFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	examples, err := ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	AttachMetaData(cache, examples)
	if filterStatusCodeOk {
		examples = FilterStatusCodeOkExamples(examples)
	}
	model := TrainedModel(examples)

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	example := NextExampleToBeAnnotated(model, examples)
	if example == nil {
		return errors.New("No example to annotate")
	}

	rtm.SendMessage(rtm.NewOutgoingMessage("Ready to annotate!", channelID))
	rtm.SendMessage(rtm.NewOutgoingMessage(example.Url, channelID))
	prevTimestamp := ""

annotationLoop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.AckMessage:
				prevTimestamp = ev.Timestamp
			case *slack.MessageEvent:
				text := ev.Text
				if len(text) > 1 || len(text) == 0 {
					break
				}
				r := []rune(text)[0]
				act := rune2ActionType(r)

				switch act {
				case LABEL_AS_POSITIVE:
					example.Annotate(POSITIVE)
					model = TrainedModel(examples)
					rtm.AddReaction("heavy_plus_sign", slack.NewRefToMessage(channelID, ev.Timestamp))
				case LABEL_AS_NEGATIVE:
					example.Annotate(NEGATIVE)
					model = TrainedModel(examples)
					rtm.AddReaction("heavy_minus_sign", slack.NewRefToMessage(channelID, prevTimestamp))
				case SKIP:
					rtm.SendMessage(rtm.NewOutgoingMessage("Skiped this example", channelID))
					break
				case SAVE:
					rtm.SendMessage(rtm.NewOutgoingMessage("Saved labeld examples", channelID))
					WriteExamples(examples, outputFilename)
				case HELP:
					rtm.SendMessage(rtm.NewOutgoingMessage(ActionHelpDoc, channelID))
				case EXIT:
					rtm.SendMessage(rtm.NewOutgoingMessage("EXIT", channelID))
					break annotationLoop
				default:
					break annotationLoop
				}
				example = NextExampleToBeAnnotated(model, examples)
				if example == nil {
					return errors.New("No example to annotate")
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(example.Url, channelID))
			case *slack.InvalidAuthEvent:
				return errors.New("Invalid credentials")
			default:
			}
		}
	}
	return nil
}
