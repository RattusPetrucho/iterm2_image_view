package viewer

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
)

// View struct
type Viewer struct {
	images    []string
	sizes     []string
	current   int
	autoclose bool
}

// View Constructor
func NewViewer(autoclose bool, dirname, filename string) (v *Viewer, err error) {
	v = new(Viewer)
	v.autoclose = autoclose

	v.images, v.current, err = getFilesList(dirname, filename)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Main loop for key events
func (v *Viewer) MainLoop() error {
	if err := v.display(); err != nil {
		return err
	}

	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyEsc || char == 'q' {
			fmt.Print("\033[2J")
			fmt.Print("\033[0;0H")
			if v.autoclose {
				exec.Command("osascript", "-e", `tell application "iTerm" to tell current tab of current window to close`).Start()
			}
			break
		} else if char == 'n' {
			v.nextDisplay()
		} else if char == 'p' {
			v.prevDisplay()
		} else {
			fmt.Printf("You pressed: %q\r\n", char)
			fmt.Println(key)
		}
	}

	return nil
}

var img = map[string]struct{}{
	".jpg":  struct{}{},
	".jpeg": struct{}{},
	".png":  struct{}{},
	".gif":  struct{}{},
}

// Get all image files in directory
func getFilesList(dirname, filename string) ([]string, int, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, 0, err
	}

	i, cur := 0, 0
	var images []string
	var ext string
	for _, val := range files {
		ext = strings.ToLower(path.Ext(val.Name()))
		if _, ok := img[ext]; ok {
			images = append(images, filepath.Join(dirname, val.Name()))
			if filename == val.Name() {
				cur = i
			}
			i++
		}
	}

	return images, cur, nil
}
