package app

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileInfo struct {
	FilePath  string
	Content   []byte
	Sum       string
	IsRenamed bool
}

var numOfWorker int

func renamefilegeneratorConcurrent() float64 {
	log.Println("start")
	start := time.Now()

	flag.IntVar(&numOfWorker, "w", runtime.NumCPU(), "number of workers")

	chanFileContent := readFiles()

	chanFileSums := make([]<-chan FileInfo, numOfWorker)

	for i := 0; i < numOfWorker; i++ {
		chanFileSums[i] = getSum(chanFileContent)
	}
	chanFileSum := mergeChanFileInfo(chanFileSums...)

	chanRenames := make([]<-chan FileInfo, numOfWorker)

	for i := 0; i < numOfWorker; i++ {
		chanRenames[i] = rename(chanFileSum)
	}
	chanRename := mergeChanFileInfo(chanRenames...)

	counterRenamed := 0
	counterTotal := 0

	for fileInfo := range chanRename {
		if fileInfo.IsRenamed {
			counterRenamed++
		}

		counterTotal++
	}

	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)

	duration := time.Since(start)
	log.Println("done in", duration.Seconds(), "seconds")

	return duration.Seconds()
}

func readFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo, numOfWorker)

	go func() {

		defer close(chanOut)

		err := filepath.Walk(tempPath, func(path string, info fs.FileInfo, err error) error {
			time.Sleep(time.Nanosecond)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			buf, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			chanOut <- FileInfo{
				FilePath: path,
				Content:  buf,
			}

			log.Println("success read", path)

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
	}()

	return chanOut
}

func getSum(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo, numOfWorker)

	go func() {
		defer close(chanOut)
		for fileInfo := range chanIn {
			time.Sleep(time.Nanosecond)
			fileInfo.Sum = fmt.Sprintf("%x", md5.Sum(fileInfo.Content))

			chanOut <- fileInfo

			log.Println("success sum", fileInfo.FilePath)
		}
	}()

	return chanOut
}

func mergeChanFileInfo(chanInMany ...<-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo, numOfWorker)
	wg := new(sync.WaitGroup)

	wg.Add(len(chanInMany))

	for _, eachChan := range chanInMany {
		go func(eachChan <-chan FileInfo) {
			for eachChanData := range eachChan {
				chanOut <- eachChanData
			}
			wg.Done()
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func rename(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo, numOfWorker)

	go func() {
		defer close(chanOut)
		for fileInfo := range chanIn {
			time.Sleep(time.Nanosecond)
			newPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", fileInfo.Sum))

			err := os.Rename(fileInfo.FilePath, newPath)
			fileInfo.IsRenamed = err == nil
			chanOut <- fileInfo

			log.Println("success rename", fileInfo.FilePath)
		}
	}()

	return chanOut
}
