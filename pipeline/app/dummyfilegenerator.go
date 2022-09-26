package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const totalFile = 3000
const contentLength = 5000

func init() {
	rand.Seed(time.Now().Unix())
}

func dummyfilegenerator() {
	log.Println("start")
	start := time.Now()

	generateFiles()

	duration := time.Since(start)

	log.Println("done in", duration.Seconds(), "seconds")
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func generateFiles() {
	err := os.RemoveAll(tempPath)

	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(tempPath, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < totalFile; i++ {
		filename := filepath.Join(tempPath, fmt.Sprintf("file-%d.txt", i))

		content := randomString(contentLength)

		err := ioutil.WriteFile(filename, []byte(content), os.ModePerm)

		if err != nil {
			log.Fatal("Error writing file", filename)
		}

		if (i+1)%100 == 0 {
			log.Println(i+1, "files created")

		}

	}

	log.Printf("%d files out of %d total files created", totalFile, totalFile)
}
