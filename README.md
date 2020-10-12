# SendMail (Golang)

Allow you to send email with this Go program

## How to use
You can build this program
```bash
go build -o gomail SendMail.go
#and execute the built program
./gomail
```

You can now launch this program
```bash
go run SendMail.go
```

## Usage (Multiline content is possible)

```bash
SMTP: 
FROM: 
TO: 
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
SMTP: mail.domain.com
FROM: sender@domain.com
TO: receiver@domain.com
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
From: Me <sender@domain.com>
To: You <receiver@domain.com>
Subject: Hello
Date: Mon, 12 Oct 2020 20:00:00 +0200
Message-ID: <c5896269-c2c4-77e3-4bd7-a3b5feBc71a@golangmail.this>
Content-Type: text/plain; charset="UTF-8"
Content-Transfer-Encoding: base64

SGVsbG8sCgpUaGlzIG1haWwgaXMgc2VudCB3aXRoIEdvbGFuZy4KCkJ5ZSw=
--------------------------------------

Sending in progress... please wait!
250: Mail sent!
```

## Me
[LinkedIn](https://fr.linkedin.com/in/kenji-duriez-9b93bb141)
