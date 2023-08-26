package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func SendVerifyAccountEmail() {
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := "smtp.gmail.com"
	t, err := template.ParseFiles("app/pkg/mailer/templates/account-verify.html")
	if err != nil {
		fmt.Println(err)
	}
	
	var body bytes.Buffer;
	subject := "Subject: Account Verification\n"
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("%s%s\n\n", subject, headers)))
	t.Execute(&body, nil)
	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", SMTPEmail, SMTPPassword,  SMTPHost), SMTPEmail, []string{"ryanali456@gmail.com"}, body.Bytes())
	if err != nil {
        fmt.Printf("smtp error: %s", err)
        return
    }
}

