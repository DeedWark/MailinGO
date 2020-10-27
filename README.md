# SendMail (Golang)

Allow you to send email with this Go program

## How to use
You can build this program
```bash
go build -o gomail SendMail.go
#and execute the built program
./gomail

#mv gomail /usr/bin/gomail
gomail
```

You can now launch this program
```bash
go run SendMail.go
```

## Usage (Multiline content is possible)

```bash
FROM: 
TO:
SMTP:
From: 
To: 
Subject: 
CONTENT [. to quit]

.
Encode body in base64 [Y/n]:

250: Message sent
```
Example:
```bash
FROM: sender@domain.com
TO: receiver@domain.com
SMTP (default: smtp.domain.com):
From: Me <sender@domain.com>
To: You <receiver@domain.com>
Subject: Hello
CONTENT [. to quit]
Hello, 

This mail is sent with Golang.

Bye,
.
Encode body in base64 [Y/n]: Y

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

Example with attachment:
```
FROM: sender@domain.com
TO: receiver@domain.com
SMTP (default: smtp.domain.com):
From: Me <sender@domain.com>
To: You <receiver@domain.com
Subject: Hello
CONTENT [. to quit]
Hello,

This mail with attachment is sent with Golang.

Bye,
.
Encode body in base64 [Y/n]: Y
Attachment [Y/n]: y
File: ./test.png

---------------Overview---------------
Date: Tue, 27 Oct 2020 14:03:49 +0100
From: Me <sender@domain.com>
To: You <receiver@domain.com>
Subject: Hello
Message-ID: <687a38e5-6d7e-499c-1607-2d696574c354@golangmail.this>
X-Mailer: SendMail-Golang v1.0
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
