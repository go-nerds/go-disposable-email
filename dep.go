package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gookit/color"
)

func getAvailableDomains() ([]string, error) {
	var availableDomains []string
	response, err := http.Get("https://www.1secmail.com/api/v1/?action=getDomainList")

	if err != nil {
		return availableDomains, errors.New("error in getting domain list from server")
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return availableDomains, errors.New("error in parsing the response")
	}

	_ = json.Unmarshal(responseData, &availableDomains)

	return availableDomains, nil

}

func createEmail(name string, domain string) (string, error) {
	info := fmt.Sprintf("https://www.1secmail.com/?login=%v&domain=%v", name, domain)
	resp, err := http.Get(info)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("email was not created")
	}

	email := name + "@" + domain
	return email, nil

}

func checkMail(name string, domainOnly string) int {
	s := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=getMessages&login=%v&domain=%v", name, domainOnly)
	resp, _ := http.Get(s)
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(body), &mailsResponse)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	totalEmails := len(mailsResponse)

	switch totalEmails {
	case 0:
		color.Yellow.Println("Mailbox is empty")
	case 1:
		color.Cyan.Println("You received 1 mail in your mailbox")
	default:
		color.Green.Println("You received,", totalEmails, "mail in your mailbox")
	}

	return totalEmails
}

func deleteMail(name string, domainOnly string) error {
	delUrl := "https://www.1secmail.com/mailbox"
	data := url.Values{}
	data.Set("action", "deleteMailbox")
	data.Set("login", name)
	data.Set("domain", domainOnly)
	_, err := http.PostForm(delUrl, data)
	return err
}
