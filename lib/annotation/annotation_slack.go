package annotation

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/file"
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
		outputFilename = util.NewOutputFilename()
		fmt.Fprintln(os.Stderr, "'output-filename' is not specified. "+outputFilename+" is used as output-filename instead.")
	}

	if channelID == "" {
		_ = cli.ShowCommandHelp(c, "slack")
		return cli.NewExitError("`channel` is a required field.", 1)
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	cache, err := cache.NewCache()
	if err != nil {
		return err
	}

	examples, err := file.ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	stat := example.GetStat(examples)
	msg := rtm.NewOutgoingMessage(fmt.Sprintf("Positive:%d, Negative:%d, Unlabeled:%d", stat["positive"], stat["negative"], stat["unlabeled"]), channelID)
	rtm.SendMessage(msg)

	cache.AttachMetaData(examples)
	if filterStatusCodeOk {
		examples = util.FilterStatusCodeOkExamples(examples)
	}
	model := classifier.NewBinaryClassifier(examples)
	e := NextExampleToBeAnnotated(model, examples)
	if e == nil {
		return errors.New("No e to annotate")
	}

	rtm.SendMessage(rtm.NewOutgoingMessage("Ready to annotate!", channelID))
	showExample(rtm, model, e, channelID)
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
					e.Annotate(example.POSITIVE)
					model = classifier.NewBinaryClassifier(examples)
					rtm.AddReaction("heavy_plus_sign", slack.NewRefToMessage(channelID, prevTimestamp))
				case LABEL_AS_NEGATIVE:
					e.Annotate(example.NEGATIVE)
					model = classifier.NewBinaryClassifier(examples)
					rtm.AddReaction("heavy_minus_sign", slack.NewRefToMessage(channelID, prevTimestamp))
				case SKIP:
					rtm.SendMessage(rtm.NewOutgoingMessage("Skiped this e", channelID))
					break
				case SAVE:
					rtm.SendMessage(rtm.NewOutgoingMessage("Saved labeld examples", channelID))
					file.WriteExamples(examples, outputFilename)
				case HELP:
					rtm.SendMessage(rtm.NewOutgoingMessage(ActionHelpDoc, channelID))
				case EXIT:
					rtm.SendMessage(rtm.NewOutgoingMessage("EXIT", channelID))
					break annotationLoop
				default:
					break annotationLoop
				}
				e = NextExampleToBeAnnotated(model, examples)
				if e == nil {
					return errors.New("No e to annotate")
				}
				showExample(rtm, model, e, channelID)
			case *slack.InvalidAuthEvent:
				return errors.New("Invalid credentials")
			default:
			}
		}
	}
	file.WriteExamples(examples, outputFilename)
	return nil
}

func showExample(rtm *slack.RTM, model classifier.BinaryClassifier, example *example.Example, channelID string) {
	activeFeaturesStr := "Active Features: "
	for _, pair := range SortedActiveFeatures(model, *example, 5) {
		activeFeaturesStr += fmt.Sprintf("%s(%+0.1f) ", pair.Feature, pair.Weight)
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s\nScore: %+0.2f\n%s", example.Url, example.Score, activeFeaturesStr), channelID))
}
