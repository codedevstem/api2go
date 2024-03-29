package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const chunkSize = 64000

func DeepCompareDirectories(dir1, dir2 string) (bool, error) {
	filesDir1, err := ioutil.ReadDir(dir1)
	if err != nil {
		return false, err
	}
	filesDir2, err := ioutil.ReadDir(dir2)
	if err != nil {
		return false, err
	}
	if len(filesDir1) != len(filesDir2) {
		return false, nil
	}
	for i := 0; i < len(filesDir1); i++ {
		fileEq, err := DeepCompareFiles(fmt.Sprintf("%s/%s", dir1, filesDir1[i].Name()), fmt.Sprintf("%s/%s", dir2, filesDir2[i].Name()))
		if err != nil {
			return false, err
		}
		if !fileEq {
			fmt.Printf("Failed: %s\n", filesDir1[i].Name())
			return false, nil
		}
		fmt.Printf("Matched: %s\n", filesDir1[i].Name())
	}
	return true, nil
}

func DeepCompareFiles(file1, file2 string) (bool, error) {

	// Open file1
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f1 *os.File) {
		_ = f1.Close()
	}(f1)

	// Open file2
	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f2 *os.File) {
		_ = f2.Close()
	}(f2)

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			output, err := exec.Command("diff", file1, file2).Output()
			if err != nil {
				switch err.(type) {
				case *exec.ExitError:
					println(string(output))
				default:
					log.Fatal("failed to execute command: ", err)
				}
			}
			return false, nil
		}
	}
}
