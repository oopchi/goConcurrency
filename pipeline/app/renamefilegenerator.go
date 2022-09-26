package app

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func renamefilegenerator() float64 {
	log.Println("start")
	start := time.Now()

	proceed()

	duration := time.Since(start)
	log.Println("done in", duration.Seconds(), "seconds")

	return duration.Seconds()
}

func proceed() {
	counterTotal := 0
	counterRenamed := 0

	err := filepath.Walk(tempPath, func(path string, info fs.FileInfo, err error) error {
		time.Sleep(time.Nanosecond)

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		counterTotal++

		buf, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		log.Println("success read without concurrent", path)
		time.Sleep(time.Nanosecond)

		sum := fmt.Sprintf("%x", md5.Sum(buf))

		log.Println("success sum without concurrent", path)
		time.Sleep(time.Nanosecond)

		destinationPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", sum))

		err = os.Rename(path, destinationPath)

		if err != nil {
			return err
		}

		counterRenamed++

		log.Println("success rename without concurrent", path)

		return nil
	})

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)
}
