# Nagios-Hipchat

Send Nagios notifications to HipChat via the HipChat API


## Templates

From the HipChat API Docs:

Must be valid HTML and entities must be escaped. May contain basic tags: a, b, i, strong, em, br, img, pre, code, lists, tables.

To specify paramterized data to add, use the format {{.key}}

ex. if called with -p host:www.example.com -p "message:Hello World" the template:

    {{.host}} is down.<br/>
    {{.message}}

would produce:

    www.example.com is down.<br/>
    Hello World

## Colors

yellow, green, red, purple, gray, random.
