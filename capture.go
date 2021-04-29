package main

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"time"

	"httpcap/common"
	"httpcap/config"
	"httpcap/reader"
	"httpcap/writer"
)

func startCapture() {
	input := reader.NewRAWInput(config.Setting.InterfaceName, config.Setting.Port)

	go CopyMulty(input)

	for {
		select {
		// case <-stop:
		// 	return
		case <-time.After(1 * time.Second):
		}
	}
}

func CopyMulty(src reader.InputReader) (err error) {
	// Don't exit on panic
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); !ok {
				fmt.Printf("PANIC: pkg: %v %s \n", r, debug.Stack())
			} else {
				fmt.Printf("PANIC: pkg: %s \n", debug.Stack())
			}
			log.Fatal(r.(error))
		}
	}()

	http := writer.NewHttpOutput("")
	memcache := writer.NewMemcacheOutput("")

	buf := make([]byte, 5*1024*1024)

	for {
		nr, raw, er := src.Read(buf)
		if nr > 0 && len(buf) > nr {
			common.Debug("[DEBUG]", raw.SrcAddr, ":", raw.SrcPort, "->", raw.DestAddr, ":", raw.DestPort, "size:", nr)
			if config.Setting.Service == "memcache" {
				memcache.Write(buf[0:nr], int(raw.SrcPort), int(raw.DestPort), raw.SrcAddr, raw.DestAddr, raw.Seq)
			} else {
				if !(config.Setting.Service != "" && config.Setting.Service != "http") {
					http.Write(buf[0:nr], int(raw.SrcPort), int(raw.DestPort), raw.SrcAddr, raw.DestAddr, raw.Seq)
				}
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return err
}
