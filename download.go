package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/bogem/id3v2"
)

func downloadMP3(u string) (mp3Path string, err error) {
	data, err := getAndParseHTML(u)
	if err != nil {
		return
	}

	os.MkdirAll("downloads/img", 0755)

	mp3Path = path.Join("downloads", fmt.Sprintf("%s - %s - %s.mp3", data.artist(), data.songTitle(), data.ShareID))
	cmd := exec.Command("ffmpeg", "-i", data.songURL(), mp3Path)
	cmd.Run()

	imgPath := path.Join("downloads", "img", fmt.Sprintf("%s - %s - %s.jpeg", data.artist(), data.songTitle(), data.ShareID))
	resp, err := http.Get(data.albumURL())
	if err != nil {
		err = fmt.Errorf("Failed to fetch album image: \n%s", data.albumURL())
		return
	}
	defer resp.Body.Close()
	imgFile, err := os.Create(imgPath)
	if err != nil {
		err = fmt.Errorf("Failed to create album image file: \n%s", imgPath)
	}
	defer imgFile.Close()
	_, err = io.Copy(imgFile, resp.Body)
	if err != nil {
		err = fmt.Errorf("Failed to download album image file: \n%s \n==> %s", data.albumURL(), imgPath)
	}

	tag, err := id3v2.Open(mp3Path, id3v2.Options{Parse: true})
	if err != nil {
		err = fmt.Errorf("Failed to open/parse mp3 file %s: %s", mp3Path, err)
	}
	defer tag.Close()

	artwork, err := ioutil.ReadFile(imgPath)
	if err != nil {
		err = fmt.Errorf("Failed to open read cover image file %s: %s", imgPath, err)
	}

	t := time.Now()
	tag.SetArtist(data.artist())
	tag.SetTitle(data.songTitle())
	tag.SetAlbum(t.Format(fmt.Sprintf("%s 2006", data.artist())))
	tag.SetYear(t.Format("2006"))
	tag.AddAttachedPicture(id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: data.songTitle(),
		Picture:     artwork,
	})

	if err = tag.Save(); err != nil {
		err = fmt.Errorf("Error while saving a tag: %v", err)
	}

	return
}
