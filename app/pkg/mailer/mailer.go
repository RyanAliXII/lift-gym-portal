package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

// func SendEmailVerification(to []string, recieverName string, publicId string) error {
// 	SMTPEmail := os.Getenv("SMTP_EMAIL")
// 	SMTPPassword := os.Getenv("SMTP_PASSWORD")
// 	SMTPHost := "smtp.gmail.com"
// 	AppURL := os.Getenv("APP_URL")
// 	t, err := template.ParseFiles("app/pkg/mailer/templates/account-verify.html")
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	link := fmt.Sprintf("%s/clients/verification/%s", AppURL, publicId)
// 	var body bytes.Buffer;
// 	subject := "Subject: Email Verification of LIFT Fitness Gym Account\n"
// 	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
// 	body.Write([]byte(fmt.Sprintf("%s%s\n\n", subject, headers)))
// 	err = t.Execute(&body, map[string]string{"name": recieverName, "link": link})
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println("Sending email verification.")
// 	start := time.Now()
// 	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", SMTPEmail, SMTPPassword,  SMTPHost), SMTPEmail, to, body.Bytes())
// 	if err != nil {
// 		fmt.Println(err)
//         return err
//     }
// 	elapsed := time.Since(start)
// 	fmt.Printf("Email verification has been sent. %v", elapsed)
// 	return nil
// }

// func SendEmailPasswordReset(to []string, recieverName string, publicId string) error {
// 	SMTPEmail := os.Getenv("SMTP_EMAIL")
// 	SMTPPassword := os.Getenv("SMTP_PASSWORD")
// 	SMTPHost := "smtp.gmail.com"
// 	AppURL := os.Getenv("APP_URL")
// 	t, err := template.ParseFiles("app/pkg/mailer/templates/password-reset.html")
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	link := fmt.Sprintf("%s/change-password?key=%s", AppURL, publicId)
// 	var body bytes.Buffer;
// 	subject := "Subject: Password Reset of LIFT Fitness Gym Account\n"
// 	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
// 	body.Write([]byte(fmt.Sprintf("%s%s\n\n", subject, headers)))
// 	err = t.Execute(&body, map[string]string{"name": recieverName, "link": link})
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println("Sending password reset email.")
// 	start := time.Now()
// 	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", SMTPEmail, SMTPPassword,  SMTPHost), SMTPEmail, to, body.Bytes())
// 	if err != nil {
// 		fmt.Println(err)
//         return err
//     }
// 	elapsed := time.Since(start)
// 	fmt.Printf("Password reset email has been sent. %v", elapsed)
// 	return nil
// }



func SendEmailVerification(to []string, recieverName string, publicId string) error {
	if(len(to) == 0) {
		return nil
	}
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := "smtp.gmail.com"
	AppURL := os.Getenv("APP_URL")
	t, err := template.ParseFiles("app/pkg/mailer/templates/account-verify.html")
	if err != nil {
		fmt.Println(err)
		return err
	}
	mail := gomail.NewMessage()

	mail.SetHeader("From", SMTPEmail)
	mail.SetHeader("To", to[0])
	mail.SetHeader("Subject", "Email Verification of LIFT Fitness Gym Account")
	link := fmt.Sprintf("%s/clients/verification/%s", AppURL, publicId)
	var body bytes.Buffer;
	err = t.Execute(&body, map[string]string{"name": recieverName, "link": link})
	if err != nil {
		fmt.Println(err)
		return err
	}
	mail.SetBody("text/html", body.String())
	d := gomail.NewDialer(SMTPHost, 587, SMTPEmail, SMTPPassword)
	if err := d.DialAndSend(mail); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Email sent successfully.")
	return nil
}

func SendEmailPasswordReset(to []string, recieverName string, publicId string) error {
	if(len(to) == 0) {
		return nil
	}
	SMTPEmail := os.Getenv("SMTP_EMAIL")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := "smtp.gmail.com"
	AppURL := os.Getenv("APP_URL")
	tmpl, err := template.ParseFiles("app/pkg/mailer/templates/password-reset.html")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	link := fmt.Sprintf("%s/change-password?key=%s", AppURL, publicId)
	var body bytes.Buffer;
	err = tmpl.Execute(&body, map[string]string{"name": recieverName, "link": link})
	if err != nil {
		fmt.Println(err)
		return err
	}
    mail := gomail.NewMessage()
	mail.SetHeader("From", SMTPEmail)
	mail.SetHeader("To", to[0])
	mail.SetHeader("Subject", "Password Reset of LIFT Fitness Gym Account")
	mail.SetBody("text/html", body.String())
	d := gomail.NewDialer(SMTPHost, 587, SMTPEmail, SMTPPassword)
	if err := d.DialAndSend(mail); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Email sent successfully.")

	return nil
}
