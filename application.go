package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"runtime"
	"os"
	"sync"
)

type Application struct {
	Locations []LogfileLocation
}

func (a *Application) List(c *cli.Context) error {

	availableApps := make([]LogfileLocation, 0)

	for _, appConfig := range a.Locations {

		if appConfig.OperatingSystem == runtime.GOOS {
			availableApps = append(availableApps, appConfig)
			break
		}
	}

	if len(availableApps) == 0 {
		fmt.Printf("no applications available for %v", runtime.GOOS)
		os.Exit(0)
	}

	for _, app := range availableApps {
		fmt.Printf("%v: %v\n", app.ApplicationName, strings.Join(app.Tags, ", "))
	}

	return nil
}



func (a *Application) ViewAction(c *cli.Context) error {

	keepWatching := false
	follow := false
	usePolling := false
	lineNumber := int64(0)
	if c.Bool("F") {
		keepWatching = true
		follow = true
	}
	if c.Bool("f") {
		follow = true
	}

	lineNumber = c.Int64("n")

	if runtime.GOOS == "windows" {
		usePolling = true
	}

	if !c.Args().Present() {
		fmt.Errorf("failed")
	}
	applicationNames := c.Args()
	fmt.Printf("Show logs from: %v\n", strings.Join(applicationNames, ", "))

	locations := make([]LogfileLocation, 0)
	foundArray := make([]string, 0)

	for _, config := range a.Locations {
		if IsConfigMatch(config, applicationNames) {
			foundArray = append(foundArray, config.ApplicationName)
			locations = append(locations, config)
		}
	}

	notFound := make([]string, 0)

	for _, app := range applicationNames {
		appName := strings.Split(app, ":")[0]
		if !InArray(appName, foundArray) {
			notFound = append(notFound, appName)
		}
	}

	if len(notFound) > 0 {
		fmt.Printf("no log file location configured for: %v\n", strings.Join(notFound, ", "))
	}

	if len(foundArray) == 0 {
		fmt.Errorf("no applications to tail logs, exiting")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for _, location := range locations {
		fmt.Printf("%v: %v%v\n", location.ApplicationName, location.DirectoryPath, location.Filename)
		wg.Add(1)
		go RenderLogfileLocation(location, int64(-1*(10*lineNumber)), keepWatching, follow, usePolling, wg)

	}

	//wg.Wait()

	for {
		runtime.Gosched()
	}

	return nil
}
