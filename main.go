package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type LogfileLocation struct {
	Filename        string
	DirectoryPath   string
	ApplicationName string
	FileType        string
	Tags            []string
	OperatingSystem string
}


func main() {

	conf := make([]LogfileLocation, 0)

	contents, err := ioutil.ReadFile("locations.json")
	if err != nil {
		fmt.Printf("Failed to find config")
		os.Exit(1)
	}
	err = json.Unmarshal(contents, &conf)
	if err != nil {
		fmt.Printf("Failed to read config")
		os.Exit(1)
	}

	application := &Application{
		Locations: conf,
	}

	cliManager := cli.NewApp()
	cliManager.Name = "tailor"
	cliManager.Usage = "watch logs from applications by name"

	cliManager.Commands = []cli.Command{
		{
			Name:    "view",
			Aliases: []string{"v"},
			Usage:   "view logs of an application",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f",
					Usage: "The -f option causes tail to not stop when end of file is reached, but rather to wait for additional data to be appended to the input.  The -f option is ignored if the standard input is a pipe, but not if it is a FIFO.",
				}, cli.BoolFlag{
					Name:  "F",
					Usage: "The -F option implies the -f option, but tail will also check to see if the file being followed has been renamed or rotated.  The file is closed and reopened when tail detects that the filename being read from has a new inode number.  The -F option is ignored if reading from standard input rather than a file.",
				}, cli.Int64Flag{
					Name:  "n",
					Usage: "The location is number lines",
				}, cli.BoolFlag{
					Name:  "q",
					Usage: "Suppresses printing of headers when multiple files are being examined",
				},
			},
			Action: application.ViewAction,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list names of supported applications",
			Action:  application.List,
		},
	}
	err = cliManager.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}
