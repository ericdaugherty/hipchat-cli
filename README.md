# Hipchat-CLI

This is a command line interface (CLI) tool to send notifications to a HipChat room using HipChat API V2.

This was developed to serve as a Nagios Notifaction mechnaims, but is general purpose enough to be used to
send any messages to HipChat.

## Usage

Prerequisites:
- Go Lang (https://golang.org/doc/install)

Install:
- Download or Clone this repository (go get github.com/ericdaugherty/hipchat-cli)
- Build (go build .)
- Run the executable (./hipchat-cli)

Usage:
hipchat-cli uses command line flags to specify all the arguments.  Here are the available flags.

       -t, --token=                                      HiptChat V2 API Token
       -r, --room=                                       The HipChat room name
       -m, --message=                                    The Message body to send. Available in custom templates as {{.message}}
       -p, --param=                                      A key-value pair to use to fill the template. -p name:value Available
                                                         in custom templates as {{.name}}
           --template-file=                              Template file to use in place of the default template.
           --template-body=                              Template definition to use in place of the default template.
       -c, --color=[yellow|green|red|purple|gray|random] The color to use to display the notifcation
       -n, --notify                                      If present, the message will trigger a HipChat user notification.

The --token and --room parameters are required for all calls.

The token flag must be passed a HipChat API V2 Token.
The room flag must be passed a room name that the Token has permission to publish to.

The --message token can be used to send simple messages without formatting.

If you wish to define more robust messages with formatting, please specify a template file or an inline template body.
(see the Template section below)  If both are specified, the template file is used.  If neither is specified, the
default template is used, which is just the message body.

## Templates

Template files are defined as Go HTML templates (see https://golang.org/pkg/html/template/ for detailed instructions).

Generally, the template file can support basic HTML tags: (a, b, i, strong, em, br, img, pre, code, lists, tables)
and template values defined as {{.name}}.  So a simple template file that just displays a bold message would be:

    <b>{{.message}<b>

Any values passed in using the -p flag are also available. A slightly more complex template could be:

    <strong>{{.host}}<strong>
    {{.message}}

This could be called with the following flags:

    ./hipchat-cli -t <token> -r <room-name> -m "The host is inaccessible" -p host:example.com

## Nagios

The Hipchat-CLI can easily be used to send nagios notifications.

Define a command to send notifications using Hipchat-CLI

    define command{
            command_name    notify-hipchat
            command_line    <path to hipchat-cli>/hipchat-cli -t <token> -r <room> -m <Your Message>
    }

Use can use the various Nagios Macros to pass information about the current state/notification to HipChat.  The full
list of available macros is available here: https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/macrolist.html

Hipchat-CLI contains a default nagios template that can be activated using the --nagios flag.  Here is a recommended
command definition for use with the default template:

    define command{
            command_name    notify-hipchat-service
            command_line    $USER1$/hipchat-cli -t <token> -r <room> -p host:$HOSTNAME$ -p status:$SERVICEOUTPUT$ --message $LONGSERVICEOUTPUT$
    }

The default nagios template is defined as:

    <strong>Nagios Notification</strong><br/>
    Host: {{.host}}<br/>
    Status: <b>{{.status}}</b><br/>
    {{.message}}

 