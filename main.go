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
const dirName = "./testdata"
const filesDir = "files"

func ConvertMP4ToHLS(fileName string, destPath string) error {
	_, file := filepath.Split(fileName)
	createVideoDirName := strings.ToLower(strings.TrimSuffix(file, filepath.Ext(file)))
	err := os.MkdirAll(strings.Join([]string{destPath, createVideoDirName}, "/"), 0777)
	if err != nil {
		return fmt.Errorf("failed creating the dir %s with the following error: %s", strings.Join([]string{destPath, createVideoDirName}, "/"), err)
	}

	cmd := exec.Command(commandName, "-i", fileName, "-c:a", "libmp3lame", "-b:a", "128k", "-map", "0:0", "-f", "segment", "-segment_time", "10", "-segment_list", strings.Join([]string{destPath, createVideoDirName, "outputlist.m3u8"}, "/"), "-segment_format", "mpegts", strings.Join([]string{destPath, createVideoDirName, "output%03d.ts"}, "/"))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg failed to splice the file %s with the following error: %#v \n", fileName, err)
	}

	return nil
}

func main() {
	_, err := exec.LookPath(commandName)
	if err != nil {
		fmt.Println("install ffmpeg to the system and run again")
		os.Exit(1)
	}

	var files []string

	err = filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "Zone.Identifier") && !strings.Contains(path, "fixtures") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("data file opened incorrectly with the following error: %#v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, fileName := range files[1:] {
		wg.Add(1)
		go func(fileName string, wg *sync.WaitGroup) {
			defer wg.Done()
			err = ConvertMP4ToHLS(fileName, filesDir)
			if err != nil {
				fmt.Printf("Conversion of MP4 to HLS failed with the following error: %s", err)
			}
		}(fileName, &wg)
	}

	wg.Wait()
}
