package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const testFileDir = "./testdata/fixtures"

func TestConvertToHLSSingleFile(t *testing.T) {
	fileName := strings.Join([]string{testFileDir, "single", "Drone.mp4"}, "/")
	destPath := strings.Join([]string{testFileDir, "single"}, "/")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Error("file directory was not created correctly")
	}

	err := ConvertMP4ToHLS(fileName, destPath)
	if err != nil {
		t.Errorf("convert files to hls failed with error: %v\n", err)
	}

	files, err := ioutil.ReadDir(strings.Join([]string{destPath, "drone"}, "/"))
	if err != nil {
		t.Errorf("The directory of files that were created from the ConvertMP4ToHLS method failed with the following error: %s\n", err)
	}
	// check that the correct amount of files has been created to be used for HLS and the drone video
	// This should work because the drone video and ffmpeg encoding method we call all use set time intervals, so the
	// amount of outputted files should be deterministic
	expectedOutputNumOfFiles := 5
	actualOutputNum := len(files)
	if expectedOutputNumOfFiles != actualOutputNum {
		t.Errorf("Drone.mp4 was not converted into the correct amount of output files.  Expected: %d, Actual: %d\n", expectedOutputNumOfFiles, actualOutputNum)
	}

	// cleanup test data
	os.RemoveAll(strings.Join([]string{destPath, "drone"}, "/"))
}
