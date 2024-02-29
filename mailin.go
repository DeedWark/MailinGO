// @Kenji DURIEZ - [DeedWark] - 2020
// Build an email and send it in Go

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

const ( // COLOR
	greenTXT  = "\033[92m"   // OK
	cyanTXT   = "\033[96m"   // INFO
	yellowTXT = "\033[1;32m" // Others
	redTXT    = "\033[91m"   // ERROR
	endTXT    = "\033[00m"   // Ending color
)

var (
	optSmtpServ string // MX/SMTP flag
	smtpServ    string // MX/SMTP server
	port        string // PORT
	mailFrom    string // MAIL FROM
	rcptTo      string // RCPT TO
	hFrom       string // Header From
	hTo         string // Header To
	hSub        string // Subject
	body        string // Body
	content     string // Content
	date        string // Date
	attach      string // Attachment
	ctype       string // Content-Type

	// OS STDIN SCANNER
	sc = bufio.NewScanner(os.Stdin)

	// CURRENT DATE
	cDate = time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

	// MORE OPTIONS
	messageId       string // Message-ID
	xmailer         string // X-Mailer
	charset         string // Encoding
	promptContent   bool   // Write Content with prompt (Allow HTML)
	htmlFile        string // Read HTML file as Body
	htmlFileContent []byte // HTML file content
	txtFile         string // Read txt file content
	txtFileContent  []byte // Txt file content
	bs64            bool   // Set base64 encoding
	xprio           string // X-Priority
	boundary        string // Custom Boundary
	encoding        string // Change encode (7bit / 8bit / binary)
	gmail           bool   // Allow auth (Gmail...)
	saveEml         bool   // Save email to an EML file
	silent          bool   // Silent mode - Do not disaply overview or info
)

// ALL OPTIONS
func usage() {
	fmt.Printf(`
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
--gmail          Enable authentication (Gmail)
--x-mailer       Set a custom X-Mailer (default "MailinGO v1.0")
--x-priority     Set a custom X-Priority (default "1")
--charset        Set a custom charset (default "UTF-8")
--html-file      Import a HTML file as body
--text-file      Import a TXT file as body
--boundary       Set a custom boundary (default "------=_MIME_BOUNDARY_MAILIN_GO--")
--content-type   Set a custom Content-Type (default "text/plain")
--encoding       Set an encoding (default "7bit")
--base64         Encode body in base64 (default no)
--prompt         Write body with a Prompt (HTML allowed) 
--save           Save email to an EML file 
--silent         Silent mode - Do not display overview or info ` + "\r\n\n")
}

func flags() {
	// Define FLAGS
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
	flag.BoolVar(&gmail, "gmail", false, "Enable authentication (for Gmail)")
	// MORE OPTIONS
	// flag.StringVar(&mid, "mid", "<c1882e5b-18b0-3ab5-89a0-ce6a534da8d4@golangmail.this>", "Set a custom Message-ID")
	flag.StringVar(&xmailer, "x-mailer", "MailinGO v1.0", "Set a custom X-Mailer")
	flag.StringVar(&xprio, "x-priority", "1", "Set a custom X-Priority")
	flag.StringVar(&charset, "charset", "UTF-8", "Set a charset format")
	flag.StringVar(&htmlFile, "html-file", "", "Import HTML file as Body")
	flag.StringVar(&txtFile, "text-file", "", "Import Text file as Body")
	flag.StringVar(&boundary, "boundary", "----=_MIME_BOUNDARY_GOO_LANG--", "Set a custom Boudnary")
	flag.StringVar(&ctype, "content-type", "text/plain", "Set a custom Content-Type")
	flag.StringVar(&encoding, "encoding", "7bit", "Set an encoding")
	flag.BoolVar(&bs64, "base64", false, "Encode body in base64")
	flag.BoolVar(&promptContent, "body-prompt", false, "Write content with a Prompt (HTML allowed)")
	flag.BoolVar(&saveEml, "save", false, "Save email to an EML file")
	flag.BoolVar(&silent, "silent", false, "Silent mode - Do not display overview or info")

	flag.Parse()

	if flag.Arg(0) == "help" {
		usage()
		os.Exit(0)
	}
}

func setCharset(charset string) string {
	/////////////
	// Charset //
	/////////////
	if charset != "" {
		switch strings.ToLower(charset) {
		case "utf-8", "utf8":
			charset = "\"UTF-8\""
		case "usascii", "us", "us-ascii":
			charset = "\"US-ASCII\""
		default:
			charset = "\"UTF-8\""
		}
	}

	return charset
}

func setMessageID() string {
	//////////////////////////////////////////////////////////////////////////
	// Message-ID -> <c1882e5b-18b0-3ab5-89a0-ce6a534da8d4@golangmail.this> //
	//////////////////////////////////////////////////////////////////////////
	b := make([]byte, 16)
	rand.Read(b)
	randomId := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	messageId := "<" + randomId + "@golangmail.this>"

	return messageId
}

func setEncoding(encoding string) string {
	//////////////////////
	// Choosen Encoding //
	//////////////////////
	if encoding != "" {
		switch strings.ToLower(encoding) {
		case "7bit", "7-bit":
			encoding = "7bit"
		case "8bit", "8-bit":
			encoding = "8bit"
		case "qp", "quoted", "quoted-printable", "printable":
			encoding = "quoted-printable"
		default:
			encoding = "7bit"
		}
	}

	return encoding
}

func resolveMX(rcptTo string) string {
	/////////////////////////////////////
	//      RESOLVE MX WITH DOMAIN     //
	/////////////////////////////////////
	cutAddress := strings.Split(rcptTo, "@") // [1] // remove @
	domainOnly := cutAddress[len(cutAddress)-1]

	mxServ := []string{}
	mxs, _ := net.LookupMX(domainOnly) // Resolve MX

	if len(mxs) != 0 {
		for _, mx := range mxs {
			mxRaw := strings.TrimRight(mx.Host, ".") // Cut ending "."
			mxServ = append(mxServ, mxRaw)           // Store MX in mxServ list
		}
	}
	cutMx := strings.Join(mxServ, "\n")  // join MX with \n
	mxList := strings.Split(cutMx, "\n") // Slice at \n

	rMx := mxList[0]
	if rMx != "" {
		smtpServ = rMx
	} else {
		fmt.Printf("SMTP server not found!" + "\n\n")
		fmt.Print("SMTP: ")
		sc.Scan()            // Get
		smtpServ = sc.Text() // Store os stdin
	}

	if optSmtpServ != "" {
		smtpServ = optSmtpServ
	}

	return smtpServ
}

func sendMail() {
	//////////////////////
	// CONTENT - PROMPT //
	//////////////////////
	if promptContent {
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
		body = strings.Join(block, "\n") // Join multiline content
	}

	//////////////////////
	// HTML File Import //
	//////////////////////
	if htmlFile != "" {
		htmlFileRaw, err := os.Open(htmlFile) // Open the HTML file
		if err != nil {
			fmt.Println(redTXT + "Cannot open HTML file" + endTXT)
			log.Fatalln(err)
		}

		reader := bufio.NewReader(htmlFileRaw)  // Init the file reader
		htmlFileContent, _ = io.ReadAll(reader) // Read and get HTML file content
		body = string(htmlFileContent)
		ctype = "text/html"
	}

	//////////////////////
	// TEXT File Import //
	//////////////////////
	if txtFile != "" {
		txtFileRaw, err := os.Open(txtFile) // Open txt file
		if err != nil {
			fmt.Println(redTXT + "Cannot open TEXT file" + endTXT)
			log.Fatalln(err)
		}

		reader := bufio.NewReader(txtFileRaw)  // Init the file reader
		txtFileContent, _ = io.ReadAll(reader) // Read and get HTML file content
		body = string(txtFileContent)
		ctype = "text/plain"
	}

	///////////////////////////////
	// Content-Transfer-Encoding //
	///////////////////////////////
	encoding := setEncoding(encoding)

	if bs64 && ctype == "text/html" {
		encoding = "7bit"
	} else if bs64 && ctype != "text/html" {
		encoding = "base64"
		body = base64.URLEncoding.EncodeToString([]byte(body))
		if len(body) > 77 {
			body = rfcSplit(body, 76, "\n")
		}
	} else {
		if len(body) > 77 {
			body = rfcSplit(body, 76, "\n")
		}
	}

	////////////////
	// Attachment //
	////////////////
	charset := setCharset(charset)

	if attach != "" {
		fileRaw := attach

		contentFile, err := os.ReadFile(fileRaw) // Read and get content file
		if err != nil {
			log.Fatalln(redTXT+"File error:"+endTXT, err)
		}

		mimeFile := http.DetectContentType(contentFile)

		fileOnly := strings.Split(attach, "/") // Split at "/" in case of Unix Path
		filename := fileOnly[len(fileOnly)-1]  // Get only filename

		//
		// ENCODE FILE/ATTACHMENT IN BASE64
		//
		encodedFile := base64.StdEncoding.EncodeToString(contentFile)

		if len(encodedFile) > 77 {
			encodedFile = rfcSplit(encodedFile, 76, "\n")
		}

		content = "Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n" +
			"--" + boundary + "\r\n" +
			"Content-Type: " + ctype + "; charset=" + charset + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + body + "\r\n" +
			"--" + boundary + "\r\n" +
			"Content-Type: " + mimeFile + "; name=\"" + filename + "\"" + "\r\n" +
			"Content-Description: " + filename + "\r\n" +
			"Content-Disposition: attachment; filename=\"" + filename + "\"" + "\r\n" +
			"Content-Transfer-Encoding: base64" + "\r\n\r\n" +
			encodedFile + "\r\n\r\n" + "--" + boundary
	} else {
		content = "Content-Type: " + ctype + "; charset=" + charset + "\r\n" +
			"Content-Transfer-Encoding: " + encoding + "\r\n" +
			"\r\n" + body
	}

	messageId := setMessageID()

	baseContent := "Date: " + date + "\r\n" +
		"From: " + hFrom + "\r\n" +
		"To: " + hTo + "\r\n" +
		"Subject: " + hSub + "\r\n" +
		"Message-ID: " + messageId + "\r\n" +
		"X-Mailer: " + xmailer + "\r\n" +
		"X-Priority: " + xprio + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		content

	if !silent {
		fmt.Println("\r\n" + yellowTXT + "---------------Overview---------------" + endTXT + "\n" + baseContent + "\n" + yellowTXT + "--------------------------------------" + endTXT)
	}

	// SAVE EML
	if saveEml {
		write := []byte(baseContent + "\r\n")
		err := os.WriteFile("./savedEmail.eml", write, 0644)
		if err != nil {
			fmt.Println(redTXT + "Cannot save this email to an EML file!")
		}
	}

	if !silent {
		fmt.Println(cyanTXT + "I am trying to send that... please wait!" + "\n" + endTXT)
	}

	resolveMX(rcptTo)

	if gmail {
		if mailFrom != "" {
			// ASK password
			fmt.Print("Password: ")
			password, _ := term.ReadPassword(int(syscall.Stdin))

			from := mailFrom
			err := smtp.SendMail("smtp.gmail.com:587",
				smtp.PlainAuth("", from, string(password), "smtp.gmail.com"),
				from, []string{rcptTo}, []byte(baseContent))
			if err != nil {
				fmt.Println(redTXT + "Error with Auth" + endTXT)
				log.Fatalln(err)
			}
		}
	} else {
		//
		// Connect to SMTP serv
		mx, err := smtp.Dial(smtpServ + ":" + port)
		if err != nil {
			fmt.Println(redTXT + "Error: Cannot connect to " + smtpServ + ":" + port + "\n" + endTXT)
			log.Fatalln(err)
		}
		defer mx.Close()

		//
		// Set MailFrom and RcptTo
		mx.Mail(mailFrom)
		mx.Rcpt(rcptTo)

		//
		// Send email body
		mxc, err := mx.Data()
		if err != nil {
			fmt.Println(redTXT + "Error: " + endTXT)
			log.Fatalln(err)
		}
		defer mxc.Close()
		buf := bytes.NewBufferString(baseContent)
		if _, err = buf.WriteTo(mxc); err != nil {
			if !silent {
				fmt.Println(redTXT + "500: Mail not sent!" + endTXT)
			}
		} else {
			if !silent {
				fmt.Println(greenTXT + "250: Mail sent!  -->  Message-ID: " + messageId + "\r\n")
			}
		}

	}
}

func rfcSplit(body string, limit int, end string) string {
	///////////////////////////////////////////////////////////////////////////////
	// Split attachment base64 encoding according to RFC (max 76 chars by line) //
	/////////////////////////////////////////////////////////////////////////////
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result string

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}

	return result
}

func main() {
	flags() // CALL FLAGS

	// Check if rcptTo is empty
	if rcptTo == "" {
		fmt.Print("RCPT TO: ")
		sc.Scan()          // Get
		rcptTo = sc.Text() // Store os stdin
	}

	sendMail()
}
