package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"html/template"
	"io/ioutil"
	"log"
)

const defaultTemplate = `
<strong>{{.Host}}</strong><br/>
<strong>{{.Message}}</strong><br/>
`

var token string
var room string
var message string
var templatePath string

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	flag.StringVar(&token, "token", "", "HiptChat V2 API Token")
	flag.StringVar(&room, "room", "", "The HipChat room name")
	flag.StringVar(&message, "message", "", "The notification to send to HipChat")
	flag.StringVar(&templatePath, "template", "", "Template File to use in place of default")
	flag.Parse()
}

func main() {
	// Validate the input.
	invalid := false
	if token == "" {
		invalid = true
		fmt.Println("You must specify a token.")
	}
	if room == "" {
		invalid = true
		fmt.Println("You must specify a room name.")
	}
	if message == "" {
		invalid = true
		fmt.Println("You must specify a message.")
	}
	if invalid {
		flag.Usage()
		return
	}

	client := hipchat.NewClient(token)

	notifReq := &hipchat.NotificationRequest{Message: formatMessage(), MessageFormat: "html"}

	resp, err := client.Room.Notification(room, notifReq)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error sending notification.", err)
		fmt.Println(string(body))
	} else {
		fmt.Println("OK")
	}
}

func formatMessage() string {

	var t *template.Template
	var err error

	if templatePath != "" {
		t, err = template.ParseFiles(templatePath)
		if err != nil {
			log.Fatalln("Unable to parse the specified template file.", err)
		}
	} else {
		t, err = template.New("m").Parse(defaultTemplate)
		if err != nil {
			log.Fatalln("Unable to parse the message template.", err)
		}
	}

	data := make(map[string]string)
	data["Host"] = "www.nba.com"
	data["Message"] = "Yo dawg, your server be down"

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		log.Fatalln("Unable to execute the message template.", err)
	}

	return out.String()
}
