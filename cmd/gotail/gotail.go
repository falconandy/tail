// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"os"
	"os/signal"
)

func args2config() (tail.Config, int64) {
	config := tail.Config{Follow: true}
	n := int64(0)
	maxlinesize := int(0)
	flag.Int64Var(&n, "n", 0, "tail from the last Nth location")
	flag.IntVar(&maxlinesize, "max", 0, "max line size")
	flag.BoolVar(&config.Follow, "f", false, "wait for additional data to be appended to the file")
	flag.BoolVar(&config.ReOpen, "F", false, "follow, and track file rename/rotation")
	flag.BoolVar(&config.Poll, "p", false, "use polling, instead of inotify")
	flag.StringVar(&config.TempLinkDirectory, "t", "", "directory for temp hard links (Windows only)")
	flag.Parse()
	if config.ReOpen {
		config.Follow = true
	}
	config.MaxLineSize = maxlinesize
	return config, n
}

func main() {
	config, n := args2config()
	if flag.NFlag() < 1 {
		fmt.Println("need one or more files as arguments")
		os.Exit(1)
	}

	if n != 0 {
		config.Location = &tail.SeekInfo{-n, os.SEEK_END}
	}

	done := make(chan bool)
	tails := make(chan *tail.Tail, len(flag.Args()))
	for _, filename := range flag.Args() {
		go tailFile(filename, config, done, tails)
	}

	processInterrupts(tails)

	for _, _ = range flag.Args() {
		<-done
	}
}

func tailFile(filename string, config tail.Config, done chan bool, tails chan *tail.Tail) {
	defer func() { done <- true }()
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	tails <- t
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

func processInterrupts(tails chan *tail.Tail) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func(){
		for range signalChan {
			for {
				select {
				case t := <-tails:
					t.Cleanup()
				default:
					os.Exit(0)
				}
			}
		}
	}()
}