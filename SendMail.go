//@Kenji DURIEZ - 2020
//Send an email like Telnet in Golang with KeyboardScanner

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	yellowTxt = "\033[93m%s\033[00m"
	greenTxt  = "\033[92m%s\033[00m"
	cyanTxt   = "\033[96m%s\033[00m"
	errorTxt  = "\033[91m%s\033[00m"
)

func main() {
	//INITIALIZE THE KEYBOARD SCANNER
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("SMTP: ")
	scanner.Scan()
	smtpServ := scanner.Text()

	fmt.Print("FROM: ")
	scanner.Scan()
	mailFrom := scanner.Text()

	fmt.Print("TO: ")
	scanner.Scan()
	rcptTo := scanner.Text()

	fmt.Print("From: ")
	scanner.Scan()
	hfrom := scanner.Text()

	fmt.Print("To: ")
	scanner.Scan()
	hto := scanner.Text()

	fmt.Print("Subject: ")
	scanner.Scan()
	hsub := scanner.Text()

	dt := time.Now()
	hdate := dt.Format("Mon, 02 Jan 2006 15:04:05 -0700")

	fmt.Println("CONTENT [. to quit]")
	block := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "." {
			block = append(block, line)
			continue
		} else {
			break
		}
	}

	//Random ID
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	randid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	hid := "<" + randid + "@golangmail.this>"

	content := strings.Join(block, "\n")

	//Base64 encoding
	fmt.Print("Encode body in base64 [Y/n]: ")
	scanner.Scan()
	choice := scanner.Text()
	b64 := content
	switch choice {
	case "Y", "y", "yes", "Yes", "YES":
		b64 = base64.URLEncoding.EncodeToString([]byte(content))
	case "N", "n", "no", "No", "NO":
		b64 = content
	default:
		b64 = content
	}

	//Body
	contentmore := "From: " + hfrom + "\r\n" +
		"To: " + hto + "\r\n" +
		"Subject: " + hsub + "\r\n" +
		"Date: " + hdate + "\r\n" +
		"Message-ID: " + hid + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + b64

	fmt.Println("\r\n" + "---------------Overview---------------" + "\n" + contentmore + "\n" + "--------------------------------------")

	fmt.Printf("\n"+yellowTxt, "Sending in progress... please wait!"+"\n")

	// Connect to SMTP server
	mx, err := smtp.Dial(smtpServ + ":25")
	if err != nil {
		log.Fatal(err)
	}
	defer mx.Close()
	// Set the sender and recipient.
	mx.Mail(mailFrom)
	mx.Rcpt(rcptTo)
	// Send the email body.
	mxc, err := mx.Data()
	if err != nil {
		log.Fatal(err)
	}

	defer mxc.Close()
	buf := bytes.NewBufferString(contentmore)
	if _, err = buf.WriteTo(mxc); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n"+greenTxt, "250: Mail sent")
}
