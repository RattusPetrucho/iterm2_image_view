package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/RattusPetrucho/iterm2_image_view/viewer"
)

func init() {
	ec := make(chan os.Signal, 2)
	signal.Notify(ec, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ec
		fmt.Print("\033[2J")
		fmt.Print("\033[0;0H")
		os.Exit(1)
	}()
}

var dir_path = flag.String("d", "", "images directory path")
var file_path = flag.String("f", "", "image file path")
var autoclose = flag.Bool("c", false, "close terminal tab/window when quit from application")

func main() {
	flag.Parse()

	if *dir_path != "" {
	} else if *file_path != "" {
		*dir_path, *file_path = filepath.Split(*file_path)
		if *file_path == "" {
			log.Fatal("did not set file or dir path")
		}
	} else {
		log.Fatal("did not set file or dir path")
	}

	v, err := viewer.NewViewer(*autoclose, *dir_path, *file_path)
	if err != nil {
		log.Fatal(err)
	}

	if err = v.MainLoop(); err != nil {
		log.Fatal(err)
	}
}
