package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/pterm/pterm"
)

type CheckMail struct {
	Id          int    `json:"id"`
	From        string `json:"from"`
	Subject     string `json:"subject"`
	Date        string `json:"date"`
	Attachments []struct {
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		Size        int    `json:"size"`
	}
	Body     string `json:"body"`
	Textbody string `json:"textBody"`
}

var mailsResponse []CheckMail

func main() {

	red := color.FgRed.Render
	info := color.FgLightCyan.Render

	pterm.DefaultBox.
		WithRightPadding(10).
		WithLeftPadding(10).
		WithTopPadding(2).
		WithBottomPadding(2).
		Println(info("Go Disposable email created by AAVision"))

	availableDomains, err := getAvailableDomains()

	if err != nil {
		fmt.Println(err)
		return
	}

	selectedDomain, _ := pterm.DefaultInteractiveSelect.WithOptions(availableDomains).Show()

	generatedName := generateRandomString(10)
	email, err := createEmail(generatedName, string(selectedDomain))

	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := os.Stat(selectedDomain); errors.Is(err, os.ErrNotExist) {
		creationError := os.Mkdir(selectedDomain, os.ModePerm)
		if creationError != nil {
			fmt.Println(err)
			return
		}
	}

	clearConsole()

	color.Cyan.Println("Your Temporary Email is:", red(email))
	color.Cyan.Println("Mailbox content is refreshed automatically every 5 seconds.")
	color.Cyan.Println("All Emails are saved in", string(selectedDomain), "folder and", generatedName, "file")

	fmt.Println(red("--------------------------------------------------------"))

	handleInterrupt(generatedName, string(selectedDomain))

	for {
		checkMail(generatedName, string(selectedDomain))
		toggleMap(mailsResponse)
		saveMailsToFile(generatedName, string(selectedDomain))
		time.Sleep(5 * time.Second)
	}
}
