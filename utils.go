package main

import (
	"strings"
	"runtime"
)

func InArray(needle string, hay []string) bool {
	for _, h := range hay {
		if needle == h {
			return true
		}
	}
	return false
}

func IsConfigMatch(location LogfileLocation, inArray []string) bool {

	if location.OperatingSystem != runtime.GOOS {
		return false
	}
	for _, ina := range inArray {
		parts := strings.Split(ina, ":")
		appName := parts[0]
		if appName == location.ApplicationName {
			if len(parts) > 1 {
				tags := strings.Split(parts[1], ",")
				for _, tag := range tags {
					if InArray(tag, location.Tags) {
						return true
					}
				}
			} else {
				return true
			}
		}
	}
	return false
}
