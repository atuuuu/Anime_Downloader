package streamtape_downloader

import (
	"authentication"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type reqLien struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Result result `json:"result"`
}

type result struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Url    string `json:"url"`
	Ticket string `json:"ticket"`
}

//Récupère le lien de téléchargement
func getDownloadLink(vid VideoData) (reqLien, error) {
	var ticket string
	var url reqLien
	var err error

	if ticket, err = getTicket(vid); err != nil {
		log.Println("requestBuilder.go:32 " + err.Error())
		return url, err
	}

	if url, err = getURL(ticket, vid); err != nil {
		log.Println("requestBuilder.go:37 " + err.Error())
		return url, err
	}

	return url, nil
}

/////////////////PREPARE THE DOWNLOAD///////////////////

//Execute la requète visant à récupérer le ticket, puis récupère le ticket
func getTicket(vid VideoData) (string, error) {
	var req *http.Request
	var client = &http.Client{}
	var rep *http.Response
	var ticket string
	var err error

	if req, err = ReqBuilderLogin(vid); err != nil {
		log.Println("requestBuilder.go:55 " + err.Error())
		return "", err
	}

	if rep, err = client.Do(req); err != nil {
		log.Println("requestBuilder.go:60 " + err.Error())
		return "", err
	}
	defer rep.Body.Close()

	if ticket, err = extractTicket(rep); err != nil {
		log.Println("requestBuilder.go:66 " + err.Error())
		return "", err
	}

	return ticket, nil
}

//Construit la requète pour obtenir le ticket de téléchargement
func ReqBuilderLogin(info VideoData) (*http.Request, error) {
	var url = strings.Builder{}
	var req *http.Request
	var err error

	url.WriteString("https://api.streamtape.com/file/dlticket?file=")
	url.WriteString(info.Fileid)
	url.WriteString("&login=")
	url.WriteString(authentication.Streamtape_login)
	url.WriteString("&key=")
	url.WriteString(authentication.Streamtape_pass)

	//fmt.Println(url.String())

	if req, err = http.NewRequest("GET", url.String(), nil); err != nil {
		log.Println("requestBuilder.go:86 " + err.Error())
		return nil, err
	}

	return req, nil
}

//Récupère le ticket de téléchargement de la requète
func extractTicket(req *http.Response) (string, error) {
	var resp = make([]byte, 10000)
	var err error
	var ticket reqLien
	var n int

	if n, err = req.Body.Read(resp); err != nil {
		log.Println("requestBuilder.go:101 " + err.Error())
		return "", err
	}

	err = json.Unmarshal(resp[:n], &ticket)
	if err != nil {
		log.Println("requestBuilder.go:107 " + err.Error())
		return "", err
	}

	if ticket.Status != 200 {
		log.Println("requestBuilder.go:116 " + string(resp))
		return "", fmt.Errorf("les serveurs de streamtape bloquent actuellement les téléchargements, merci de réessayer plus tard")
	}

	return ticket.Result.Ticket, nil
}

/////////////////////ACTUAL DOWNLOAD////////////////////////

//Construit le lien de téléchargement à partir du ticket
func ReqBuilderLink(ticket string, info VideoData) (*http.Request, error) {
	var url strings.Builder
	var req *http.Request
	var err error

	url.WriteString("https://api.streamtape.com/file/dl?")
	url.WriteString("file=")
	url.WriteString(info.Fileid)
	url.WriteString("&ticket=")
	url.WriteString(ticket)
	url.WriteString("&captcha_response=true")

	if req, err = http.NewRequest("GET", url.String(), nil); err != nil {
		log.Println("requestBuilder.go:135 " + err.Error())
		return nil, err
	}

	return req, nil
}

//Execute la requète récupérant l'URL et récupère l'URL de téléchargement
func getURL(ticket string, vid VideoData) (reqLien, error) {
	var req *http.Request
	var client = &http.Client{}
	var rep *http.Response
	var err error
	var reqlien reqLien
	var n int

	var resp = make([]byte, 10000)
	if req, err = ReqBuilderLink(ticket, vid); err != nil {
		log.Println("requestBuilder.go:153 " + err.Error())
		return reqlien, err
	}

	time.Sleep(5 * time.Second)

	if rep, err = client.Do(req); err != nil {
		log.Println("requestBuilder.go:160 " + err.Error())
		return reqlien, err
	}
	defer rep.Body.Close()

	if n, err = rep.Body.Read(resp); err != nil {
		log.Println("requestBuilder.go:166 " + err.Error())
		return reqlien, err
	}

	if err = json.Unmarshal(resp[:n], &reqlien); err != nil {
		log.Println("requestBuilder.go:171 " + err.Error())
		return reqlien, err
	}

	return reqlien, nil
}
