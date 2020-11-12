//@Kenji DURIEZ - [DeedWark] - 2020
//Send email with Go

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const ( //COLOR
	greenTXT  = "\033[92m"   //OK
	cyanTXT   = "\033[96m"   //INFO
	yellowTXT = "\033[1;32m" //Others
	redTXT    = "\033[91m"   //ERROR
	endTXT    = "\033[00m"   //Ending color
)

var optSmtpServ string //MX/SMTP flag
var smtpServ string    //MX/SMTP server
var port string        //PORT
var mailFrom string    //MAIL FROM
var rcptTo string      //RCPT TO
var hFrom string       //Header From
var hTo string         //Header To
var hSub string        //Subject
var body string        //Body
var content string     //Content
var date string        //Date
var attach string      //Attachment
var auth bool          //Allow auth (Gmail...)
var ctype string       //Content-Type
var encoding string    //Encoding

//OS STDIN SCANNER
var sc = bufio.NewScanner(os.Stdin)

//CURRENT DATE
var cDate = time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

//MORE OPTIONS
var mid string             //Message-ID
var xmailer string         //X-Mailer
var charset string         //Encoding
var promptContent bool     //Write Content with prompt (Allow HTML)
var htmlFile string        //Read HTML file as Body
var htmlFileContent []byte //HTML file content
var txtFile string         //Read txt file content
var txtFileContent []byte  //Txt file content
var bs64 bool              //Set base64 encoding
var xprio string           //X-Priority
var boundary string        //Custom Boundary

func usage() {
	fmt.Println(`
  -s  	         Set SMTP/MX server (default "Autodetect with domain")
  -p  	         Set TCP Port (default "25/SMTP")
  -f             Set MAIL FROM (protocolar)
  -t  	         Set RCPT TO (protocolar)
--hfrom          Set Header From (ex "Me <go@lang.org>")
--hto            Set Header To (ex "You <go@pher.org>")
--subject        Set a subject
--date           Set a custom date (default "current date")
--body           Write content to Body
--attach         Add an attachment/file
--auth           Enable authentication (Gmail, Outlook...)
--x-mailer       Set a custom X-Mailer (default "SendMail-Golang v2.0")
--x-priority     Set a custom X-Priority
--charset        Set a custom charset (default "UTF-8")
--html-file      Import a HTML file as body
--text-file      Import a TXT file as body
--boundary       Set a custom boundary (default "------=_MIME_BOUNDARY_GOO_LANG--")
--content-type   Set a custom Content-Type (default "text/plain")
--base64 Encode  body in base64
--prompt         Get a prompt to write on your terminal 
`)
}

func flags() {
	//Define FLAGS
	//    TYPE       VAR      ARGS,DEFAULT    HELP
	flag.StringVar(&optSmtpServ, "s", "", "Set SMTP/MX server")
	flag.StringVar(&port, "p", "25", "Set TCP port")
	flag.StringVar(&mailFrom, "f", "", "Mail From address (MAIL FROM - Protocolar)")
	flag.StringVar(&rcptTo, "t", "", "Recipient To address (RCPT TO - Protocolar)")
	flag.StringVar(&hFrom, "hfrom", "", "Set Header From (From:)")
	flag.StringVar(&hTo, "hto", "", "Set Header To (To:)")
	flag.StringVar(&hSub, "subject", "", "Set a subject")
	flag.StringVar(&date, "date", cDate, "Set a custom date")
	flag.StringVar(&body, "body", "", "Content in body")
	flag.StringVar(&attach, "attach", "", "Add an attachment")
	flag.BoolVar(&auth, "auth", false, "Enable authentication (for Gmail/Outlook...)")
	//MORE OPTIONS
	//flag.StringVar(&mid, "mid", "<c1882e5b-18b0-3ab5-89a0-ce6a534da8d4@golangmail.this>", "Set a custom Message-ID")
	flag.StringVar(&xmailer, "x-mailer", "SendMail-Golang v2.0", "Set a custom X-Mailer")
	flag.StringVar(&charset, "charset", "UTF-8", "Set a charset format")
	flag.StringVar(&htmlFile, "html-file", "", "Import HTML file as Body")
	flag.StringVar(&txtFile, "text-file", "", "Import Text file as body")
	flag.StringVar(&xprio, "x-priority", "1", "Set a custom X-Priority")
	flag.StringVar(&boundary, "boundary", "------=_MIME_BOUNDARY_GOO_LANG--", "Set a custom Boudnary")
	flag.StringVar(&ctype, "content-type", "text/plain", "Set a custom Content-Type")
	flag.BoolVar(&bs64, "base64", false, "Encode body in base64")
	flag.BoolVar(&promptContent, "body-prompt", false, "Write content with a Prompt")

	flag.Parse()
}

func sendMail() {
	flags()

	if rcptTo == "" {
		fmt.Print("RCPT TO: ")
		sc.Scan()          //Get
		rcptTo = sc.Text() //Store os stdin
	}

	/////////////////////////////////////
	//      RESOLVE MX WITH DOMAIN     //
	/////////////////////////////////////
	cutAddress := strings.Split(rcptTo, "@") //remove @
	domainOnly := cutAddress[len(cutAddress)-1]

	mxServ := []string{}
	mxs, _ := net.LookupMX(domainOnly) //Resolve MX

	if len(mxs) != 0 {
		for _, mx := range mxs {
			mxRaw := strings.TrimRight(mx.Host, ".") //Cut ending "."
			mxServ = append(mxServ, mxRaw)           //Store MX in mxServ list
		}
	}
	cutMx := strings.Join(mxServ, "\n")  //join MX with \n
	mxList := strings.Split(cutMx, "\n") //Slice at \n

	var rMx = mxList[0]
	if rMx != "" {
		smtpServ = rMx
	} else {
		fmt.Println("SMTP server not found!" + "\n")
		fmt.Print("SMTP: ")
		sc.Scan()            //Get
		smtpServ = sc.Text() //Store os stdin
	}

	if optSmtpServ != "" {
		smtpServ = optSmtpServ
	}

	////////////////////////////////////////////////////////////////////////
	//Message-ID -> <c1882e5b-18b0-3ab5-89a0-ce6a534da8d4@golangmail.this>//
	////////////////////////////////////////////////////////////////////////
	b := make([]byte, 16)
	rand.Read(b)
	randomId := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	var messageId = "<" + randomId + "@golangmail.this>"

	//////////////////////
	// CONTENT - PROMPT //
	//////////////////////
	if promptContent == true {
		fmt.Println("CONTENT [. to quit]")
		block := []string{}
		for sc.Scan() {
			line := sc.Text()
			if line != "." {
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
		body = strings.Join(block, "\n") //Join multiline content
	}

	/////////////
	// Charset //
	/////////////
	if charset != "" {
		switch strings.ToLower(charset) {
		case "utf-8", "utf8":
			charset = "UTF-8"
		case "usascii", "us", "us-ascii":
			charset = "US-ASCII"
		case "quoted-printable", "qp", "quoted", "printable":
			charset = "quoted-printable"
		default:
			charset = "UTF-8"
		}
	}

	//////////////////////
	// HTML File Import //
	//////////////////////
	if htmlFile != "" {
		htmlFileRaw, err := os.Open(htmlFile) //Open the HTML file
		if err != nil {
			fmt.Println(redTXT + "Cannot open HTML file" + endTXT)
			log.Fatalln(err)
		}

		reader := bufio.NewReader(htmlFileRaw)      //Init the file reader
		htmlFileContent, _ = ioutil.ReadAll(reader) //Read and get HTML file content
		body = string(htmlFileContent)
		ctype = "text/html"
	}

	//////////////////////
	// TEXT File Import //
	//////////////////////
	if txtFile != "" {
		txtFileRaw, err := os.Open(txtFile) //Open txt file
		if err != nil {
			fmt.Println(redTXT + "Cannot open TEXT file" + endTXT)
			log.Fatalln(err)
		}

		reader := bufio.NewReader(txtFileRaw)      //Init the file reader
		txtFileContent, _ = ioutil.ReadAll(reader) //Read and get HTML file content
		body = string(txtFileContent)
		ctype = "text/plain"
	}

	///////////////////////////////
	// Content-Transfer Encoding //
	///////////////////////////////
	if bs64 == true {
		if ctype != "text/html" {
			encoding = "base64"
			body = base64.URLEncoding.EncodeToString([]byte(body))
		} else {
			encoding = "7bit"
		}
	} else {
		encoding = "7bit"
	}

	////////////////
	// Attachment //
	////////////////
	if attach != "" {
		fileRaw, err := os.Open(attach) //Open file
		if err != nil {
			log.Fatalln(redTXT + "Cannot open the File as attachment" + endTXT)
		}

		fileReader := bufio.NewReader(fileRaw)       //Init the file reader
		contentFile, _ := ioutil.ReadAll(fileReader) //Read and get content file
		fileOnly := strings.Split(attach, "/")       //Split at "/" in case of Unix Path
		filename := fileOnly[len(fileOnly)-1]        //Get only filename

		//
		//ENCODE FILE/ATTACHMENT IN BASE64
		//
		encodedFile := base64.StdEncoding.EncodeToString(contentFile)

		content = "Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n" +
			boundary + "\r\n" +
			"Content-Type: " + ctype + "; charset=" + charset + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + body + "\r\n" +
			boundary + "\r\n" +
			"Content-Type: application/octet-stream; name=\"" + filename + "\"" + "\r\n" +
			"Content-Description: " + filename + "\r\n" +
			"Content-Disposition: attachment; filename=\"" + filename + "\"" + "\r\n" +
			"Content-Transfer-Encoding: base64" + "\r\n\r\n" +
			encodedFile + "\r\n\r\n" + boundary
	} else {
		content = "Content-Type: " + ctype + "; charset=" + charset + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + body
	}

	var baseContent string = "Date: " + date + "\r\n" +
		"From: " + hFrom + "\r\n" +
		"To: " + hTo + "\r\n" +
		"Subject: " + hSub + "\r\n" +
		"Message-ID: " + messageId + "\r\n" +
		"X-Mailer: " + xmailer + "\r\n" +
		"X-Priority: " + xprio + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		content

	fmt.Println("\r\n" + yellowTXT + "---------------Overview---------------" + endTXT + "\n" + baseContent + "\n" + yellowTXT + "--------------------------------------" + endTXT)
	fmt.Println(cyanTXT + "Sending in progress... please wait!" + "\n" + endTXT)

	if auth != false {
		if mailFrom != "" {
			//ASK password
			fmt.Print("Password: ")
			password, _ := terminal.ReadPassword(0)

			from := mailFrom
			err := smtp.SendMail("smtp.gmail.com:587",
				smtp.PlainAuth("", from, string(password), "smtp.gmail.com"),
				from, []string{rcptTo}, []byte(body))
			fmt.Println(string(password))
			if err != nil {
				fmt.Println(redTXT + "Error with Auth" + endTXT)
				log.Fatalln(err)
			}
		}
	} else {
		//
		//Connect to SMTP serv
		mx, err := smtp.Dial(smtpServ + ":" + port)
		if err != nil {
			fmt.Println(redTXT + "Error: Cannot connect to " + smtpServ + ":" + port + "\n" + endTXT)
			log.Fatalln(err)
		}
		defer mx.Close()

		//
		//Set MailFrom and RcptTo
		mx.Mail(mailFrom)
		mx.Rcpt(rcptTo)

		//
		//Send email body
		mxc, err := mx.Data()
		if err != nil {
			fmt.Println(redTXT + "Error: " + endTXT)
			log.Fatalln(err)
		}
		defer mxc.Close()
		buf := bytes.NewBufferString(body)
		if _, err = buf.WriteTo(mxc); err != nil {
			fmt.Println(redTXT + "500: Mail not sent!" + endTXT)
		} else {
			fmt.Println(greenTXT + "250: Mail sent!  -->  Message-ID: " + messageId + "\r\n")
		}

	}
}

func main() {
	sendMail()
}
