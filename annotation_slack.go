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
		outputFilename = NewOutputFilename()
		fmt.Fprintln(os.Stderr, "'output-filename' is not specified. "+outputFilename+" is used as output-filename instead.")
	}

	if channelID == "" {
		_ = cli.ShowCommandHelp(c, "slack")
		return cli.NewExitError("`channel` is a required field.", 1)
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	cacheFilename := CacheFilename
	cache, err := LoadCache(cacheFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	examples, err := ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	stat := GetStat(examples)
	msg := rtm.NewOutgoingMessage(fmt.Sprintf("Positive:%d, Negative:%d, Unlabeled:%d", stat["positive"], stat["negative"], stat["unlabeled"]), channelID)
	rtm.SendMessage(msg)

	AttachMetaData(cache, examples)
	if filterStatusCodeOk {
		examples = FilterStatusCodeOkExamples(examples)
	}
	model := NewPerceptronClassifier(examples)
	example := NextExampleToBeAnnotated(model, examples)
	if example == nil {
		return errors.New("No example to annotate")
	}

	rtm.SendMessage(rtm.NewOutgoingMessage("Ready to annotate!", channelID))
	showExample(rtm, model, example, channelID)
	prevTimestamp := ""

annotationLoop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.AckMessage:
				prevTimestamp = ev.Timestamp
			case *slack.MessageEvent:
				if ev.Channel != channelID {
					break
				}
				text := ev.Text
				if len(text) > 1 || len(text) == 0 {
					break
				}
				r := []rune(text)[0]
				act := rune2ActionType(r)

				switch act {
				case LABEL_AS_POSITIVE:
					example.Annotate(POSITIVE)
					model = NewPerceptronClassifier(examples)
					rtm.AddReaction("heavy_plus_sign", slack.NewRefToMessage(channelID, prevTimestamp))
				case LABEL_AS_NEGATIVE:
					example.Annotate(NEGATIVE)
					model = NewPerceptronClassifier(examples)
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
				showExample(rtm, model, example, channelID)
			case *slack.InvalidAuthEvent:
				return errors.New("Invalid credentials")
			default:
			}
		}
	}
	WriteExamples(examples, outputFilename)
	cache.Save(cacheFilename)
	return nil
}

func showExample(rtm *slack.RTM, model BinaryClassifier, example *Example, channelID string) {
	activeFeaturesStr := "Active Features: "
	for _, pair := range SortedActiveFeatures(model, *example, 5) {
		activeFeaturesStr += fmt.Sprintf("%s(%+0.1f) ", pair.Feature, pair.Weight)
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s\nScore: %+0.2f\n%s", example.Url, example.Score, activeFeaturesStr), channelID))
}
