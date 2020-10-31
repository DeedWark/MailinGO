package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	greenTxt = "\033[92m%s\033[00m"
	cyanTxt  = "\033[96m%s\033[00m"
	errorTxt = "\033[91m%s\033[00m"
)

var contentmore string
var hdate string
var hfrom string
var hto string
var hrid string
var ctype string
var b64 string
var encoding string

func main() {
	//INIT THE KEYBOARD SCANNER
	scanner := bufio.NewScanner(os.Stdin)

	//MAIL FROM
	fmt.Print("MAIL FROM: ")
	scanner.Scan()
	mailFrom := scanner.Text()

	//RCPT TO
	fmt.Print("RCPT TO: ")
	scanner.Scan()
	rcptTo := scanner.Text()

	//Get rcptTo domain
	cutDomain := strings.Split(rcptTo, "@") //Get only domain
	domainOnly := cutDomain[len(cutDomain)-1]

	//Auto detect mx for a domain
	mxServ := []string{}
	mxs, _ := net.LookupMX(domainOnly) //Resolve MX
	var indx string
	if len(mxs) != 0 {
		for _, mx := range mxs {
			mxRaw := strings.TrimRight(mx.Host, ".") //cut . at the end
			mxServ = append(mxServ, mxRaw)           //get MX
		}
	}
	cutMx := strings.Join(mxServ, "\n")  //join mx
	mxList := strings.Split(cutMx, "\n") //slice
	indx = mxList[0]                     //Get only 1st mx
	var defMx string
	if indx != "" {
		defMx = " (default: " + indx + "): "
	} else {
		defMx = ": "
	}

	//MX
	fmt.Print("SMTP" + defMx)
	scanner.Scan()
	smtpServ := scanner.Text()
	if smtpServ == "" {
		smtpServ = indx
	}

	//Header From
	fmt.Print("From: ")
	scanner.Scan()
	hfrom := scanner.Text()

	//Header To
	fmt.Print("To: ")
	scanner.Scan()
	hto := scanner.Text()

	//Subject
	fmt.Print("Subject: ")
	scanner.Scan()
	hsub := scanner.Text()

	//Current date
	dt := time.Now()
	hdate := dt.Format("Mon, 02 Jan 2006 15:04:05 -0700")

	//BODY
	fmt.Println("CONTENT [. to quit]")
	block := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "." {
			ctype = "text/plain"
			if line == "<html>" || line == "</html>" {
				ctype = "text/html"
				block = append(block, line)
				continue
			} else {
				ctype = "text/plain"
				block = append(block, line)
				continue
			}
		} else {
			break
		}
	}

	//Join multiline [Body]
	content := strings.Join(block, "\n")

	//Random ID -> Message-ID
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		hrid = ""
	} else {
		randid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		hid := "<" + randid + "@golangmail.this>"
		hrid = hid
	}

	if ctype == "text/html" {
		b64 = content
		encoding = "7bit"
	} else {
		fmt.Print("Encode body in base64 [Y/n]: ")
		scanner.Scan()
		choice := scanner.Text()
		switch strings.ToLower(choice) {
		case "y", "yes":
			b64 = base64.URLEncoding.EncodeToString([]byte(content))
			encoding = "base64"
			ctype = "text/plain"
		case "n", "no":
			b64 = content
			encoding = "7bit"
		default:
			b64 = content
			encoding = "7bit"
		}
	}

	//Add Attachment ?
	fmt.Print("Attachment [Y/n]: ")
	scanner.Scan()
	att_ch := scanner.Text()

	var contentmore string
	var baseContentType string

	switch strings.ToLower(att_ch) {

	case "y", "yes":

		fmt.Print("File: ")
		scanner.Scan()
		filePath := scanner.Text()
		fileraw, err := os.Open(filePath) //Open the attachment file
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(fileraw)       //Init file reader
		content, _ := ioutil.ReadAll(reader)     //Read and get the file content
		fileOnly := strings.Split(filePath, "/") //Split at / -> see just below
		filename := fileOnly[len(fileOnly)-1]    //Get the filename in case of path is like "../dir/dir/image.png"

		//Encode file/attachment in base64
		encodedFile := base64.StdEncoding.EncodeToString(content)
		//76 char per line for b64 attach
		for i, g := 0, len(encodedFile); i < g; i++ {
			b64buf := bytes.NewBuffer(nil)
			b64buf.WriteByte(encodedFile[i])
			if (i+1)%76 == 0 {
				b64buf.WriteString("\r\n")
			}
		}
		//All the data to send
		baseContentType = "Content-Type: multipart/mixed; boundary=\"----=_MIME_BOUNDARY_GOO_LANG\"" + "\r\n\r\n" +
			"------=_MIME_BOUNDARY_GOO_LANG" + "\r\n" +
			"Content-Type: " + ctype + "; charset=\"UTF-8\"" + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + b64 + "\r\n" +
			"------=_MIME_BOUNDARY_GOO_LANG" + "\r\n" +
			"Content-Type: application/octet-stream; name=\"" + filename + "\"" + "\r\n" +
			"Content-Description: " + filename + "\r\n" +
			"Content-Disposition: attachment; filename=\"" + filename + "\"" + "\r\n" +
			"Content-Transfer-Encoding: base64" + "\r\n\r\n" +
			encodedFile + "\r\n\r\n" + "------=_MIME_BOUNDARY_GOO_LANG--"

	case "n", "no":

		//All the data without attachment
		baseContentType = "Content-Type: " + ctype + "; charset=\"UTF-8\"" + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + b64
	}

	contentmore = "Date: " + hdate + "\r\n" +
		"From: " + hfrom + "\r\n" +
		"To: " + hto + "\r\n" +
		"Subject: " + hsub + "\r\n" +
		"Message-ID: " + hrid + "\r\n" +
		"X-Mailer: SendMail-Golang v1.0" + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		baseContentType

	//Print Overview
	fmt.Println("\r\n" + "---------------Overview---------------" + "\n" + contentmore + "\n" + "--------------------------------------")

	fmt.Printf("\n"+cyanTxt, "Sending in progress... please wait!"+"\n")

	// Connect to SMTP server
	mx, err := smtp.Dial(smtpServ + ":25")
	if err != nil {
		fmt.Printf("\n"+errorTxt, "ERROR: Cannot connect to "+smtpServ+":25"+"\n")
		log.Fatal(err)
	}
	defer mx.Close()

	// Set MailFrom and RcptTo
	mx.Mail(mailFrom)
	mx.Rcpt(rcptTo)

	// Send email body
	mxc, err := mx.Data()
	if err != nil {
		fmt.Printf("\n"+errorTxt, "Body Error!"+"\n")
		log.Fatal(err)
	}
	defer mxc.Close()
	buf := bytes.NewBufferString(contentmore)
	if _, err = buf.WriteTo(mxc); err != nil {
		fmt.Println(errorTxt, "500: Mail not sent!")
		log.Fatal(err)
	} else {
		fmt.Printf("\n"+greenTxt, "250: Mail sent!  -->  Message-ID: "+hrid+"\r\n")
	}
}
