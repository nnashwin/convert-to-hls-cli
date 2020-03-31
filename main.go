package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const commandName = "ffmpeg"
const dirName = "./data"
const filesDir = "files"

func main() {
	_, err := exec.LookPath(commandName)
	if err != nil {
		fmt.Println("install ffmpeg to the system and run again")
		os.Exit(1)
	}

	var files []string

	err = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "Zone.Identifier") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("data file opened incorrectly with the following error: %#v", err)
	}

	var wg sync.WaitGroup

	for _, fileName := range files[1:] {

		wg.Add(1)
		go func(fileName string, wg *sync.WaitGroup) {
			defer wg.Done()
			createVideoDirName := strings.ToLower(fileName[len(dirName)-1 : len(fileName)-4])
			err = os.MkdirAll(strings.Join([]string{filesDir, createVideoDirName}, "/"), 0777)
			if err != nil {
				fmt.Printf("not able to create the following dir: %s", strings.Join([]string{filesDir, createVideoDirName}, "/"))
			}

			cmd := exec.Command(commandName, "-i", fileName, "-c:a", "libmp3lame", "-b:a", "128k", "-map", "0:0", "-f", "segment", "-segment_time", "10", "-segment_list", strings.Join([]string{filesDir, createVideoDirName, "outputlist.m3u8"}, "/"), "-segment_format", "mpegts", strings.Join([]string{filesDir, createVideoDirName, "output%03d.ts"}, "/"))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("ffmpeg failed to splice the file %s with the following error: %#v \n", fileName, err)
			}
		}(fileName, &wg)
	}

	wg.Wait()
}
