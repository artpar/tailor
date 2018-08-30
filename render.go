package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"os"
	"sync"
)

func RenderLogFile(location LogfileLocation, lineNumber int64, keepWatching bool, follow bool, usePolling bool) {
	filePath := fmt.Sprintf("%v%v", location.DirectoryPath, location.Filename)
	t, err := tail.TailFile(filePath, tail.Config{
		Location: &tail.SeekInfo{
			Offset: lineNumber,
			Whence: 2,
		},
		Follow: follow,
		Poll:   usePolling,
		ReOpen: keepWatching,
		Logger: log.New(os.Stdout, location.ApplicationName+": ", log.LstdFlags),
	})
	if err != nil {
		fmt.Printf("[%v] failed to tail file: %v", location.Filename, err)
		return
	}


	for line := range t.Lines {
		if line.Err != nil {
			fmt.Printf("[%v] failed to tail file: %v", location.Filename, line.Err)
			return
		}
		fmt.Printf("%v: %v\n", location.Filename, line.Text)
	}
}

func RenderLogfileLocation(location LogfileLocation, lineNumber int64, keepWatching bool, follow bool, usePolling bool, wg sync.WaitGroup) {
	if keepWatching {
		for {
			RenderLogFile(location, lineNumber, keepWatching, follow, usePolling)
			lineNumber = 0
		}
	} else {
		RenderLogFile(location, lineNumber, keepWatching, follow, usePolling)
		wg.Done()
	}

}
