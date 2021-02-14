package main

import (
	"fmt"
	_ "github.com/asticode/go-astikit"
	_ "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
	"os"
	v "view"
)

// Vars injected via ldflags by bundler
var (
	BuiltAt string
)

func main() {
	if err := buildAppPage(); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	// Run bootstrap
	v.L.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(v.App); err != nil {
		v.L.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}

func buildAppPage() error {
	appdata := os.Getenv("appdata")
	fmt.Println(appdata)

	if _, err := os.Stat(appdata + "\\resources"); os.IsNotExist(err) {
		os.Mkdir(appdata+"\\resources", 666)
	}
	if _, err := os.Stat(appdata + "\\resources\\app"); os.IsNotExist(err) {
		os.Mkdir(appdata+"\\resources\\app", 666)
	}

	//var fd uintptr
	var appPage, jsPage *os.File
	var err error

	if appPage, err = os.Create(appdata + "\\resources\\app\\view.html"); err != nil {
		return err
	}

	if _, err := appPage.WriteString(v.SourceHTML); err != nil {
		return err
	}

	if jsPage, err = os.Create(appdata + "\\resources\\app\\index.js"); err != nil {
		return err
	}

	if _, err := jsPage.WriteString(v.SourceJS); err != nil {
		return err
	}

	return nil
}
