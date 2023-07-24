package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gookit/color"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

var ids = make(map[int]int)

func toggleMap(response []CheckMail) {
	for _, v := range response {
		if _, found := ids[v.Id]; !found {
			ids[v.Id] = 0
		}
	}
}

func handleInterrupt(name string, domainOnly string) {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		err := deleteMail(name, domainOnly)
		if err != nil {
			log.Fatal(err)
		}
		clearConsole()
		fmt.Println("Emails Deleted. Exiting.")
		os.Exit(0)
	}()
}

func clearConsole() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform. Cannot clear console.")
	}
}

func saveMailsToFile(name string, domain string) {
	for k := range ids {
		if ids[k] == 0 {
			s := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=readMessage&login=%v&domain=%v&id=%v", name, domain, k)
			resp, _ := http.Get(s)
			body, _ := io.ReadAll(resp.Body)
			time.Sleep(1 * time.Second)
			resp.Body.Close()
			var receivedEmail CheckMail
			err := json.Unmarshal([]byte(body), &receivedEmail)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {

				receivedMail := fmt.Sprintf("ID: %v\nFrom: %v\nSubject: %v\nDate: %v\nText: %v\n", receivedEmail.Id, receivedEmail.From, receivedEmail.Subject, receivedEmail.Date, receivedEmail.Textbody)

				path := domain + "/" + name + ".txt"

				f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

				defer f.Close()

				if err != nil {
					log.Fatal(err)
				}

				if _, err := f.WriteString(string(receivedMail)); err != nil {
					log.Println(err)
				}

				for _, value := range receivedEmail.Attachments {
					s := fmt.Sprintf("https://www.1secmail.com/api/v1/?action=download&login=%v&domain=%v&id=%v&file=%v", name, domain, k, value.Filename)
					fileResponse, _ := http.Get(s)
					fmt.Println(value.Filename)
					fmt.Println(s)

					defer fileResponse.Body.Close()

					out, err := os.Create(domain + "/" + value.Filename)

					if err != nil {
						return
					}

					defer out.Close()

					io.Copy(out, fileResponse.Body)

					color.Success.Println(value.Filename, "downloaded successfully!")
					time.Sleep(1 * time.Second)
				}

				ids[k] = 1
			}
		}
	}
}
