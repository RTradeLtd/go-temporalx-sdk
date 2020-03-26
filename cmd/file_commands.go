package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	pb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	au "github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

// LoadFileCommands returns a file commands object
func LoadFileCommands() cli.Commands {
	return cli.Commands{
		&cli.Command{
			Name:        "file",
			Usage:       "file upload/download commands",
			Description: "Enables access to the FileAPI",
			Subcommands: cli.Commands{fileUpload(), fileDownload()},
		},
	}
}

func fileUpload() *cli.Command {
	return &cli.Command{
		Name: "upload",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			fileName := c.String("file.name")
			if fileName == "" {
				return errors.New("file.name flag empty")
			}
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			stats, err := file.Stat()
			if err != nil {
				return err
			}
			resp, err := cl.UploadFile(ctx, file, stats.Size(), &pb.UploadOptions{
				MultiHash: c.String("multi.hash"),
				Layout:    c.String("layout"),
				Chunker:   c.String("chunker"),
			}, c.Bool("print.progress"))
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s %s\n",
				au.Bold(au.Green("hash of file:")),
				au.Bold(au.White(resp.GetHash())),
			)
			return nil
		},
		Flags: []cli.Flag{
			printProgressFlag(),
			&cli.StringFlag{
				Name:    "file.name",
				Aliases: []string{"fn"},
				Usage:   "file to upload",
			},
			&cli.StringFlag{
				Name:    "multi.hash",
				Aliases: []string{"mh"},
				Usage:   "multihash function to use",
				Value:   "sha2-256",
			},
			&cli.StringFlag{
				Name:  "layout",
				Usage: "dag layout to use",
				Value: "balanced",
			},
			&cli.StringFlag{
				Name:  "chunker",
				Usage: "chunking algorithm to use",
				Value: "default",
			},
		},
	}
}

func fileDownload() *cli.Command {
	return &cli.Command{
		Name: "download",
		Action: func(c *cli.Context) error {
			cl, err := client.NewClient(optsFromFlags(c))
			if err != nil {
				return err
			}
			if c.String("cid") == "" {
				return errors.New("cid flag empty")
			}
			if c.String("save.path") == "" {
				return errors.New("save.path flag empty")
			}
			resp, err := cl.DownloadFile(ctx, &pb.DownloadRequest{Hash: c.String("cid")}, c.Bool("print.progress"))
			if err != nil {
				return err
			}
			return ioutil.WriteFile(c.String("save.path"), resp.Bytes(), os.FileMode(0640))
		},
		Flags: []cli.Flag{
			cidFlag("cid of file to download"),
			printProgressFlag(),
			&cli.StringFlag{
				Name:    "save.path",
				Aliases: []string{"sp"},
				Usage:   "path to save file too",
			},
		},
	}
}
