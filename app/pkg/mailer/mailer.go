package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)




func SendEmailVerification(to []string, recieverName string, publicId string) error {
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := "smtp.gmail.com"
	AppURL := os.Getenv("APP_URL")
	t, err := template.ParseFiles("app/pkg/mailer/templates/account-verify.html")
	if err != nil {
		fmt.Println(err)
		return err
	}
	link := fmt.Sprintf("%s/clients/verification/%s", AppURL, publicId)
	var body bytes.Buffer;
	subject := "Subject: Email Verification of LIFT Fitness Gym Account\n"
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("%s%s\n\n", subject, headers)))
	err = t.Execute(&body, map[string]string{"name": recieverName, "link": link})
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", SMTPEmail, SMTPPassword,  SMTPHost), SMTPEmail, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
        return err
    }
	return nil
}


func SendEmailPasswordReset(to []string, recieverName string, publicId string) error {
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := "smtp.gmail.com"
	AppURL := os.Getenv("APP_URL")
	t, err := template.ParseFiles("app/pkg/mailer/templates/password-reset.html")
	if err != nil {
		fmt.Println(err)
		return err
	}
	link := fmt.Sprintf("%s/change-password?key=%s", AppURL, publicId)
	var body bytes.Buffer;
	subject := "Subject: Password Reset of LIFT Fitness Gym Account\n"
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("%s%s\n\n", subject, headers)))
	err = t.Execute(&body, map[string]string{"name": recieverName, "link": link})
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", SMTPEmail, SMTPPassword,  SMTPHost), SMTPEmail, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
        return err
    }
	return nil
}

