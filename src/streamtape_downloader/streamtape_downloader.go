package streamtape_downloader

import (
	"encoding/json"
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

type VideoData struct {
	Token  string `json:"token"`
	Fileid string `json:"fileid"`
}

func Streamtape_download(url string, folder string, errs chan error, w *astilectron.Window, del int) error {
	var data []byte
	var _ int
	var vid VideoData
	var _ *http.Response
	var retour reqLien

	fmt.Print("Now downloading from streamtape : ")
	fmt.Println(url)

	_ = os.Mkdir(folder, os.ModePerm)

	resp, err := initiateDownload(url)
	if err != nil {
		log.Println("main.go:33 " + err.Error())
		errs <- err
		return err
	}

	if data, err = parseVideoInfo(resp); err != nil {
		log.Println("streamtape_downloader.go:39 " + err.Error())
		errs <- err
		return err
	}

	if vid, err = extractVideoInfo(data); err != nil {
		log.Println("streamtape_downloader.go:45 " + err.Error())
		errs <- err
		return err
	}

	if retour, err = getDownloadLink(vid); err != nil {
		log.Println("streamtape_downloader.go:51 " + err.Error())
		errs <- err
		return err
	}

	fmt.Println(retour)
	if retour.Status != 200 {
		errs <- fmt.Errorf("streamtape est actuellement indisponible, merci de réessayer plus tard")
		log.Println(retour)
		return fmt.Errorf("streamtape est actuellement indisponible, merci de réessayer plus tard")
	}

	target := folder + "/" + retour.Result.Name

	err = DownloadFile(target, retour.Result.Url, w, del)
	if err != nil {
		log.Println("streamtape_downloader.go:55 " + err.Error())
		errs <- err
		return err
	}

	errs <- fmt.Errorf("fin")
	return nil
}

// DownloadFile will initiateDownload a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string, w *astilectron.Window, del int) error {

	if strings.Contains(filepath, "/") {
		tmp := strings.Split(filepath, "/")
		r := regexp.MustCompile(".*VOSTFR.*.mp4")
		for i := 0; i < len(tmp); i++ {
			if r.MatchString(tmp[i]) {
				filepath = tmp[i]
			}
		}

		if strings.Contains(filepath, "/") || filepath == "" {
			tmp := strings.Split(filepath, "/")
			r := regexp.MustCompile(".*.mp4")
			for i := 0; i < len(tmp); i++ {
				if r.MatchString(tmp[i]) {
					filepath = tmp[i]
				}
			}
		}
	}

	strings.Replace(filepath, ":", "_", -1)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	go Track_dl(out, float32(resp.ContentLength), w, del)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extractVideoInfo(data []byte) (VideoData, error) {
	var vid VideoData
	var err error

	data = []byte(strings.Replace(string(data), ";", "", -1))
	if err = json.Unmarshal(data, &vid); err != nil {
		log.Println("streamtape_downloader.go:36 " + err.Error())
		return vid, err
	}

	return vid, nil
}

func parseVideoInfo(resp []byte) ([]byte, error) {
	r := regexp.MustCompile("{\"token\":\".*\",\"adblock\":.*,\"blockadblock\":.*,\"noads\":.*,\"fileid\":\".*\",\"ampallow\":.*};")
	resp = r.Find(resp)
	if string(resp) == "" {
		return nil, fmt.Errorf("l'épisode n'est pas téléchargeable pour le moment")
	}

	return resp, nil
}

func initiateDownload(url string) ([]byte, error) {
	var reponse = make([]byte, 10000)
	var req *http.Request
	var resp *http.Response
	var err error
	var client = &http.Client{}

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println(err)
		return nil, err
	}

	if resp, err = client.Do(req); err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if _, err = resp.Body.Read(reponse); err != nil {
		log.Println(err)
		return nil, err
	}

	return reponse, nil
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
		fmt.Printf("%2f", percentageDL)
		fmt.Println("%")
		s := fmt.Sprintf("%f", percentageDL)
		w.SendMessage("progress>"+strconv.Itoa(del)+">"+s, func(m *astilectron.EventMessage) {
		})
		time.Sleep(5 * time.Second)
	}
}
