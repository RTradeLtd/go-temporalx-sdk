package cmd

import (
	"errors"
	"fmt"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

var (
	pubSubData                      string
	pubSubTopic                     string
	pubSubDiscover, pubSubAdvertise bool
)

// LoadPubSubCommands returns a cli commands object for the grpc pubsub client
func LoadPubSubCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "pubsub",
			Usage:       "pubsub commands",
			Description: "Enables access to PubSubAPI",
			Subcommands: cli.Commands{
				pubSubPublish(),
				pubSubSubscribe(),
				pubSubListPeers(),
				pubSubListTopics(),
			},
		},
	}
}

func pubSubPublish() *cli.Command {
	return &cli.Command{
		Name:  "publish",
		Usage: "publish a pubsub message",
		Action: func(c *cli.Context) error {
			if pubSubTopic == "" {
				return errors.New("topic flag is empty")
			}
			if pubSubData == "" {
				return errors.New("data flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			return cl.PSPublish(
				ctx,
				pubSubTopic,
				[]byte(pubSubData),
			)
		},
		Flags: pubSubFlags(),
	}
}

func pubSubSubscribe() *cli.Command {
	return &cli.Command{
		Name:  "subscribe",
		Usage: "subscribe to a pubsub topic",
		Action: func(c *cli.Context) error {
			if pubSubTopic == "" {
				return errors.New("topic flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			errChan := make(chan error, 1)
			go func() {
				if err := cl.PSSubscribe(
					ctx,
					pubSubTopic,
					pubSubDiscover,
					nil,
				); err != nil {
					errChan <- err
					return
				}
			}()
			select {
			case <-errChan:
				return err
			case <-ctx.Done():
				return nil
			}
		},
		Flags: pubSubFlags(),
	}
}

func pubSubListPeers() *cli.Command {
	return &cli.Command{
		Name:  "list-peers",
		Usage: "list pubsub peers for a topic",
		Action: func(c *cli.Context) error {
			if pubSubTopic == "" {
				return errors.New("topic flag is empty")
			}
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			peers, err := cl.PSListPeers(ctx, pubSubTopic)
			if err != nil {
				return err
			}
			for _, peer := range peers {
				fmt.Printf(
					"%s %s\n%s %s\n",
					au.Bold(au.Green("peerid:")),
					au.Bold(au.White(fmt.Sprint(peer.GetPeerID()))),
					au.Bold(au.Green("topic:")),
					au.Bold(au.White(fmt.Sprint(peer.GetTopic()))),
				)
			}
			return nil
		},
		Flags: pubSubFlags(),
	}
}

func pubSubListTopics() *cli.Command {
	return &cli.Command{
		Name:  "list-topics",
		Usage: "list known pubsub topics",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			resp, err := cl.PSGetTopics(ctx)
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s\n",
				au.Bold(au.Green(fmt.Sprint(resp))),
			)
			return nil
		},
		Flags: pubSubFlags(),
	}
}

func pubSubFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "topic",
			Usage:       "pubsub topic",
			Destination: &pubSubTopic,
		},
		&cli.StringFlag{
			Name:        "data",
			Usage:       "pubsub data",
			Destination: &pubSubData,
		},
		&cli.BoolFlag{
			Name:        "discover",
			Usage:       "enable pubsub discovery",
			Destination: &pubSubDiscover,
		},
		&cli.BoolFlag{
			Name:        "advertise",
			Usage:       "enable pubsub advertise",
			Destination: &pubSubAdvertise,
		},
	}
}
