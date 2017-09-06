package viewer

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// display current image
func (v *Viewer) display() error {
	file, err := os.Open(v.images[v.current])
	if err != nil {
		return err
	}
	defer file.Close()

	return display(file, v.images[v.current])
}

// display next image
func (v *Viewer) nextDisplay() error {
	v.current++
	if v.current >= len(v.images) {
		v.current = 0
	}

	return v.display()
}

// display previous image
func (v *Viewer) prevDisplay() error {
	v.current--
	if v.current < 0 {
		v.current = len(v.images) - 1
	}

	return v.display()
}

// Read from reader, convert image into base64 string and print to iterm2
func display(r io.Reader, filename string) error {
	fmt.Print("\033[2J")
	fmt.Print("\033[0;0H")

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	s := string(out)
	s = strings.TrimSpace(s)

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	height, err := strconv.Atoi(strings.Split(s, " ")[0])
	if err != nil {
		return err
	}

	fmt.Print("\033]1337;")
	fmt.Printf("File=inline=1")
	fmt.Printf(";height=%d", height-1)
	fmt.Printf(";size=1781800")
	fmt.Print(":")
	fmt.Printf("%s", base64.StdEncoding.EncodeToString(data))
	fmt.Print("\a\n")

	fmt.Print(filename)

	return nil
}
