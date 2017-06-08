package main

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

const defaultTemplate = `
<strong>{{.host}}</strong><br/>
<strong>{{.message}}</strong><br/>
`

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

type Options struct {
	Token        string            `short:"t" long:"token" description:"HiptChat V2 API Token" required:"true"`
	Room         string            `short:"r" long:"room" description:"The HipChat room name" required:"true"`
	Data         map[string]string `short:"p" long:"param" description:"A key-value pair to use to fill the template. -p name:value"`
	TemplatePath string            `long:"template" description:"Template File to use in place of default"`
	Color        string            `short:"c" long:"color" description:"The color to use to display the notifcation" choice:"yellow" choice:"green" choice:"red" choice:"purple" choice:"gray" choice:"random"`
	Notify       bool              `short:"n" long:"notify" description:"If present, the message will trigger a HipChat user notification."`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func init() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})
}

func main() {

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	fmt.Println("Option Notify", options.Notify)

	client := hipchat.NewClient(options.Token)

	notifReq := &hipchat.NotificationRequest{
		Message:       formatMessage(),
		MessageFormat: "html",
		Color:         parseColor(),
		Notify:        options.Notify,
	}

	resp, err := client.Room.Notification(options.Room, notifReq)
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

	if options.TemplatePath != "" {
		t, err = template.ParseFiles(options.TemplatePath)
		if err != nil {
			log.Fatalln("Unable to parse the specified template file.", err)
		}
	} else {
		t, err = template.New("m").Parse(defaultTemplate)
		if err != nil {
			log.Fatalln("Unable to parse the message template.", err)
		}
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, options.Data)
	if err != nil {
		log.Fatalln("Unable to execute the message template.", err)
	}

	return out.String()
}

func parseColor() hipchat.Color {
	switch options.Color {
	case "yellow":
		return hipchat.ColorYellow
	case "green":
		return hipchat.ColorGreen
	case "red":
		return hipchat.ColorRed
	case "purple":
		return hipchat.ColorPurple
	case "gray":
		return hipchat.ColorGray
	case "random":
		return hipchat.ColorRandom
	default:
		return hipchat.ColorYellow
	}
}
