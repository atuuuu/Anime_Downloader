package clicthune_extractor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Clicthune(urlS string) (string, error) {
	err := fmt.Errorf("init")
	var url string

	url, err = getVidURL(urlS)
	if err != nil {
		err = fmt.Errorf("extractor.go:19 " + err.Error())
		return "", err
	}

	time.Sleep(11 * time.Second)

	return url, nil
}

func getVidURL(clictuneURL string) (string, error) {
	retour, err := getRequestContent(clictuneURL)
	if err != nil {
		log.Println("extractor.go:27 " + err.Error())
		return "", err
	}
	return extractURL(retour)
}

func extractURL(site []byte) (string, error) {
	var retour string

	r := regexp.MustCompile("<a href=\"https:\\/\\/www.mylink..*\">")
	match := r.Find(site)
	candidate := strings.Split(string(match), "\"")

	r = regexp.MustCompile("^(https://)")

	for i := 0; i < len(candidate); i++ {
		if r.MatchString(candidate[i]) {
			retour = candidate[i]
		}
	}
	if retour == "" {
		err := fmt.Errorf("extractor.go:50 no url found")
		return "", err
	}

	return retour, nil
}

func getRequestContent(clictuneURL string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", clictuneURL, nil)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	var tmp = make([]byte, 10000)
	var retour = make([]byte, 10000)

	for {
		n, err := resp.Body.Read(retour)
		if err == io.EOF {
			break
		}

		fmt.Print("Lu : " + strconv.Itoa(n) + "\n\n")
		fmt.Println(string(retour))
		for i := 0; i < len(tmp); i++ {
			retour = append(retour, tmp[i])
		}
	}

	return retour, nil
}
