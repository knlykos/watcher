package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	url := "http://cove.nkodexsoft.com:3000/cove/csv"

	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	// w.SetMaxEvents(1)

	// Only notify rename and move events.
	w.FilterOps(watcher.Move, watcher.Create)

	// Only files that match the regular expression during file listings
	// will be watched.
	// r := regexp.MustCompile("^abc$")
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))
	r := regexp.MustCompile(".csv$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event.Op)
				if event.FileInfo.IsDir() == false {
					dat, err := ioutil.ReadFile(event.Path)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(string(dat))
					payload := strings.NewReader(string(dat))
					req, _ := http.NewRequest("POST", url, payload)
					fmt.Println(req)

					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("cache-control", "no-cache")
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Fatalln(err)
					}
					// fmt.Println(res)
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(body)

					fmt.Println(res)
					fmt.Println(string(body))
				}

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add("./input"); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	// if err := w.AddRecursive("test_folder"); err != nil {
	// 	log.Fatalln(err)
	// }

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	// Trigger 2 events after watcher started.
	// go func() {
	// 	w.Wait()
	// 	w.TriggerEvent(watcher.Create, nil)
	// 	w.TriggerEvent(watcher.Remove, nil)
	// }()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
