package doods_downloader

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Dood_download(url string, folder string, errs chan error, w *astilectron.Window, del int) error {
	var download string

	fmt.Println(url)

	link, err := getVideoLink(url)
	if err != nil {
		errs <- err
		return err
	}

	time.Sleep(1 * time.Second)

	download, err = getDownloadPageLink(link)
	if err != nil {
		download = url
	}

	fmt.Println(download)
	time.Sleep(1 * time.Second)

	finalLink, err := getDirectDownloadLink(download)
	if err != nil {
		errs <- err
		return err
	}

	fmt.Println(finalLink)

	time.Sleep(1 * time.Second)

	err = DownloadFile(folder, finalLink, w, del)
	if err != nil {
		errs <- err
		return err
	}

	errs <- fmt.Errorf("fin")
	return nil
}

func getVideoLink(url string) (string, error) {

	fmt.Println("getVideoLink : " + url)

	fmt.Println("Get video link from : " + url)

	r := regexp.MustCompile("(%2Fe%2F.*%0D)|(%2Fd%2F.*%0D)|(%2Fe%2F.*&id)|(%2Fd%2F.*&id)")
	link := r.FindString(url)

	link = strings.Replace(link, "%2Fe%2F", "", -1)
	link = strings.Replace(link, "%2Fd%2F", "", -1)
	link = strings.Replace(link, "%0D", "", -1)
	link = strings.Replace(link, "&id", "", -1)

	fmt.Println("lien récupéré GetVideoLink : " + link)

	if link == "" {
		return "", fmt.Errorf("pas de lien")
	}

	link = "https://dood.so/d/" + link

	return link, nil
}

func getDownloadPageLink(url string) (string, error) {
	var client = &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")
	req.Header.Add("Referer", "https://dood.so")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	n, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		n = 500000
	}

	var tmp = make([]byte, n)

	n, err = resp.Body.Read(tmp)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile("/download/.*/n/.*\" ")
	url = r.FindString(string(tmp[:n]))
	url = strings.Trim(url, "\" ")

	if url == "" {
		return "", fmt.Errorf("pas de lien")
	}

	var link = "https://dood.so" + url

	return link, nil
}

func getDirectDownloadLink(url string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Referer", "https://dood.so/")
	req.Header.Add("Host", "dood.so")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Accept-Encoding", "UTF-8")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("TE", "Trailers")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	n := int(resp.ContentLength)

	if int(resp.ContentLength) <= 0 {
		n = 50000
	}

	var tmp = make([]byte, n)

	_, err = resp.Body.Read(tmp)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile("https://.*.dood.video/.*/.*.?token=.*',")
	link := r.FindString(string(tmp))

	link = strings.Trim(link, "',")

	if link == "" {
		return "", fmt.Errorf("directDownloadLink : pas de lien")
	}

	return link, nil
}

// DownloadFile will initiateDownload a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string, w *astilectron.Window, del int) error {
	fmt.Println(url)

	r := regexp.MustCompile("([A-Za-z0-9-_\\[\\]])*.([A-Za-z0-9-_\\[\\]])*.mp4")
	filepath = r.FindString(url)

	if filepath == "" {
		filepath = "unnamed.mp4"
	}
	if strings.Contains(filepath, "/") {
		tmp := strings.Split(filepath, "/")
		r := regexp.MustCompile(".*.mp4")
		for i := 0; i < len(tmp); i++ {
			if r.MatchString(tmp[i]) {
				filepath = tmp[i]
			}
		}
	}

	filepath = strings.Replace(filepath, ":", "_", -1)
	fmt.Println(filepath)

	fmt.Println("fichier : " + filepath)

	var client = &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Referer", "https://dood.so")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	go Track_dl(out, float32(resp.ContentLength), w, del)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Track_dl(file *os.File, byteLength float32, w *astilectron.Window, del int) {
	var percentageDL float64

	for percentageDL = 0; percentageDL*100 < 99; {
		stat, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}

		percentageDL = float64(stat.Size()) / float64(byteLength)

		fmt.Print("Téléchargement : ")
		fmt.Printf("%.2f", percentageDL*100)
		fmt.Println("%")
		s := fmt.Sprintf("%f", percentageDL)
		w.SendMessage("progress>"+strconv.Itoa(del)+">"+s, func(m *astilectron.EventMessage) {})
		time.Sleep(5 * time.Second)
	}
}
