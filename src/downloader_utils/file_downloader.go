package downloader_utils

import (
	"clicthune_extractor"
	"doods_downloader"
	"fmt"
	"github.com/asticode/go-astilectron"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"streamtape_downloader"
	"strings"
)

func detectServiceProvider(url string) (func(string, string, chan error, *astilectron.Window, int) error, error) {
	req, err := http.Get(url)
	if err != nil {
		log.Println("main.go:31 " + err.Error())
		return nil, err
	}

	var tmp = make([]byte, 50000)

	_, err = req.Body.Read(tmp)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	doodR := regexp.MustCompile(".*dood.*")
	streamTapeR := regexp.MustCompile(".*streamtape.*")
	googleR := regexp.MustCompile(".*google.*")

	if streamTapeR.MatchString(string(tmp)) {
		return streamtape_downloader.Streamtape_download, nil
	}
	if doodR.Match(tmp) {
		fmt.Println("dood")
		return doods_downloader.Dood_download, nil
	}
	if googleR.Match(tmp) {
		fmt.Println("google")
		return nil, fmt.Errorf("les liens google drive ne sont pas supportés")
	}

	return nil, fmt.Errorf("support non pris en charge")
}

func InitDownload(url string, file string, w *astilectron.Window, del int) {
	var errs = make(chan error, 1)
	var err = fmt.Errorf("init")
	var lien, errMessage string
	var download func(string, string, chan error, *astilectron.Window, int) error
	var i, j, k, l = 0, 0, 0, 0

	url = adaptUrl(url)

	for l < 5 {
		for err != nil && i < 15 {
			fmt.Println("Clictune ?")

			err = nil
			lien, err = clicthune_extractor.Clicthune(url)
			if err != nil {
				log.Println("main.go:33 " + err.Error())
			}
			i++
		}

		fmt.Println(url)

		if lien == "" {
			lien = url
		}

		err = fmt.Errorf("init")
		for err != nil && lien != "" && j < 15 {
			fmt.Println("detection...")

			err = nil
			download, err = detectServiceProvider(lien)
			if err != nil {
				log.Println("main.go:49 " + err.Error())
			}
			j++
		}

		fmt.Println("téléchargement")

		if download != nil && j < 15 && i < 15 {
			for k < 15 {
				k++
				go download(lien, file, errs, w, del)
				err = <-errs
				if err != nil {
					if strings.Contains(err.Error(), "fin") {
						fmt.Println("téléchargement fini")
						w.SendMessage("success>"+strconv.Itoa(del), func(m *astilectron.EventMessage) {})
						return
					}
					errMessage = err.Error()
					fmt.Println(errMessage)
					i = 0
					j = 0
				}
			}
		}
		l++
	}
	w.SendMessage("error>"+strconv.Itoa(del)+">"+errMessage, func(m *astilectron.EventMessage) {})
}

func adaptUrl(url string) string {

	regexWww := regexp.MustCompile("^(www.)")

	url = strings.Trim(url, " ")
	url = strings.Trim(url, "\t")
	url = strings.Trim(url, "\n")

	if regexWww.MatchString(url) {
		url = "https://" + url
	}

	return url
}
