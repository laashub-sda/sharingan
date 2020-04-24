package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	_ "github.com/didi/sharingan"
)

const waitFlagParseTime = 10

var systemTest *bool
var endRunning chan bool

func stop() {
	endRunning <- true
}

func signalHandler() {
	var callback sync.Once
	// 定义并监听 kill信号, On ^C or SIGTERM
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		callback.Do(stop)
	}()
}

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

// TestMain Test started when the test binary is started. Only calls main.
func TestMain(m *testing.M) {
	if *systemTest {
		go main()
		signalHandler()
		endRunning = make(chan bool, 1)
		// Maximum waiting time(10s) for flag.Parse.
		// If the flag still missed to execute after 10 seconds, check your logic with main function.
		checkTime := time.After(waitFlagParseTime * time.Second)
		for {
			if flag.Parsed() {
				break
			}
			select {
			case <-checkTime:
				if !flag.Parsed() {
					flag.Parse()
				}
				break
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
		<-endRunning
		os.Exit(m.Run())
	}
}
