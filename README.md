# SendMail (Golang)

Send email with this Go program

## Setup
```bash
go build -o gomail sendMailV2.go

# And run the built program
./gomail

# mv gomail /usr/bin/gomail
gomail
```

## How to use

```bash
  -s  	         Set SMTP/MX server (default "Autodetect with domain")
  -p  	         Set TCP Port (default "25/SMTP")
  -f             Set MAIL FROM (protocolar)
  -t  	         Set RCPT TO (protocolar)
--hfrom          Set Header From (ex "Me <go@lang.org>")
--hto            Set Header To (ex "You <go@pher.org>")
--subject        Set a subject
--date           Set a custom date (default "current date")
--body           Add content to Body
--attach         Add an attachment/file
--auth           Enable authentication (Gmail, Outlook...)
--x-mailer       Set a custom X-Mailer (default "SendMail-Golang v2.0")
--x-priority     Set a custom X-Priority (default "1")
--charset        Set a custom charset (default "UTF-8")
--html-file      Import a HTML file as body
--text-file      Import a TXT file as body
--boundary       Set a custom boundary (default "------=_MIME_BOUNDARY_GOO_LANG--")
--content-type   Set a custom Content-Type (default "text/plain")
--encoding       Set an encoding (default "7bit")
--base64         Encode body in base64 (default no)
--prompt         Write body with a Prompt (HTML allowed) 
--save           Save email to an EML file
```
## Overview
```
---------------Overview---------------

Date: Mon, 12 Oct 2020 20:00:00 +0200From: Me <sender@domain.com>
To: You <receiver@domain.com>
Subject: Hello
Message-ID: <c5896269-c2c4-77e3-4bd7-a3b5feBc71a@golangmail.this>
Content-Type: text/plain; charset="UTF-8"
Content-Transfer-Encoding: base64

SGVsbG8sCgpUaGlzIG1haWwgaXMgc2VudCB3aXRoIEdvbGFuZy4KCkJ5ZSw=
--------------------------------------

Sending in progress... please wait!
250: Mail sent!  -->  Message-ID: <c5896269-c2c4-77e3-4bd7-a3b5feBc71a@golangmail.this>
```

- With attachment
```
---------------Overview---------------
Date: Tue, 27 Oct 2020 14:03:49 +0100
From: Me <sender@domain.com>
To: You <receiver@domain.com>
Subject: Hello
Message-ID: <687a38e5-6d7e-499c-1607-2d696574c354@golangmail.this>
X-Mailer: SendMail-Golang v1.0
X-Priority: 1
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="----=_MIME_BOUNDARY_GOO_LANG"

------=_MIME_BOUNDARY_GOO_LANG
Content-Type: text/plain
Content-Transfer-Encoding: base64

SGksIApUaGVyZSBpcyBhIHRoaW5nIGZvciB5b3UKKEl0J3MgYSB0ZXN0IGZvciBpbmZyYSBhbmQgZ
GV2KQoKRE8gTk9UIEJMT0NLIE1ZIElQIFBMRUFTRQpCZXN0IHJlZ2FyZHMs
------=_MIME_BOUNDARY_GOO_LANG
Content-Type: application/octet-stream; name="test.png"
Content-Description: test.png
Content-Disposition: attachment; filename="test.png"
Content-Transfer-Encoding: base64

iVBORw0KGgoAAAANSUhEUgAAAJMAAACUCAYAAACX4ButAAAAAXNSR0IArs4c6QAAAARnQU1BAACxj
wv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAJrSURBVHhe7dpNSltRGIBh5xEHIjhx0mmX0H24AH
dRFdyF0KHSodAVdKDgEjrQXWQHx1tb09hKK/SVkPAQHgj3y70QeLnn5Gdre3s2oCAmMmIiIyYyYiI

------=_MIME_BOUNDARY_GOO_LANG--
--------------------------------------

Sending in progress... please wait!

250: Mail sent!  -->  Message-ID: <687a38e5-6d7e-499c-1607-2d696574c354@golangmail.this>

```
## Me
[LinkedIn](https://fr.linkedin.com/in/kenji-duriez-9b93bb141)
