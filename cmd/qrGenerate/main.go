package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
	"github.com/stretchr/testify/assert/yaml"
)

func main() {
	// read config
	var alleycat models.Alleycat
	ymlData, err := os.ReadFile("./config/config.yml")
	util.CheckPanic(err)
	err = yaml.Unmarshal(ymlData, &alleycat.Config)
	util.CheckPanic(err)

	// read manifest
	var manifest []models.Checkpoint
	ymlData, err = os.ReadFile("./config/manifest.yml")
	util.CheckPanic(err)
	err = yaml.Unmarshal(ymlData, &manifest)
	util.CheckPanic(err)

	tmp, err := os.MkdirTemp("/tmp/", "tmp_")
	util.CheckPanic(err)

	// generate dummy

	for _, ch := range manifest {
		fmt.Printf("Generating %d\n", ch.ID)
		generateImage(ch.Location, alleycat.Config.Url, ch.Link, alleycat.Config.RemovalDate, tmp)
	}

	fmt.Println("Generating Grid")

	// pagecount = floored + 1 if overflow
	pages := (len(manifest) / 8)
	overflow := len(manifest) % 8
	if overflow != 0 {
		pages++
	}

	for p := range pages {
		fileList := ""
		for i := range 8 {
			if (p*8)+i+1 <= len(manifest) {
				fileList += tmp + "/" + manifest[(p*8)+i].Link + ".png "
			} else {
				fileList += "null:"
			}
		}

		//cmdTemplate := "montage -mode concatenate -tile 4x2 %s %s/grid%d.png"
		cmdTemplate := `montage -density 300 %s -tile 4x2 -geometry 874x1240+0+0 -background white -extent 874x1240 %s/grid%d.png`
		cmd := fmt.Sprintf(cmdTemplate, fileList, tmp, p)
		c := exec.Command("bash", "-c", cmd)
		err = c.Run()
		util.CheckPanic(err)
	}

	cmdTemplate := "convert %s/grid*.png grid.pdf"
	cmd := fmt.Sprintf(cmdTemplate, tmp)
	c := exec.Command("bash", "-c", cmd)
	err = c.Run()
	util.CheckPanic(err)

	fmt.Println("Generated grid.pdf")

	os.RemoveAll(tmp)
}

func generateImage(title string, url string, link string, removalDate string, outDir string) {
	qrPath := fmt.Sprintf("%s/%s_qr.png", outDir, link)
	cardPath := fmt.Sprintf("%s/%s.png", outDir, link)

	cmdTemplate := `qrencode -s 30 -d 30 -o %s %s`
	cmd := fmt.Sprintf(cmdTemplate, qrPath, url+link)
	c := exec.Command("bash", "-c", cmd)
	err := c.Run()
	util.CheckPanic(err)

	cmdTemplate = `
		magick \
 			\( -background none -size 874x100 canvas:none \) \
  			\( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 60 -size 874x caption:"%s" \) \
  			\( -background none -size 1x40 canvas:none \) \
  			\( "%s" -resize 600x600 \) \
  			\( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 36 -size 874x caption:"%s" \) \
  			\( -background none -size 1x200 canvas:none \) \
  			\( -background none -fill black -font "Montserrat-Bold" -gravity center -pointsize 28 -size 800x caption:"Dieser Sticker ist Teil einer Schnitzeljagd und wird am %s wieder entfernt." \) \
  			-background white -append \
  			-gravity north -extent 874x1240 \
  			"%s"
	`
	cmd = fmt.Sprintf(cmdTemplate, title, qrPath, url, removalDate, cardPath)
	c = exec.Command("bash", "-c", cmd)
	err = c.Run()
	util.CheckPanic(err)

	cmdTemplate = `rm %s`
	cmd = fmt.Sprintf(cmdTemplate, qrPath)
	c = exec.Command("bash", "-c", cmd)
	err = c.Run()
	util.CheckPanic(err)
}
