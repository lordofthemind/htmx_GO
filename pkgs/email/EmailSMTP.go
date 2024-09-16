package email

// import (
// 	"net/smtp"
// )

// type EmailService struct {
// 	SMTPServer   string
// 	SMTPPort     string
// 	SMTPUsername string
// 	SMTPPassword string
// }

// func NewEmailService() *EmailService {
// 	return &EmailService{
// 		SMTPServer:   "smtp.example.com",
// 		SMTPPort:     "587",
// 		SMTPUsername: "your-email@example.com",
// 		SMTPPassword: "your-email-password",
// 	}
// }

// func (e *EmailService) SendSMTPMail(from string, to []string, subject, body string) error {
// 	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPServer)
// 	msg := []byte("To: " + to[0] + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body + "\r\n")
// 	return smtp.SendMail(e.SMTPServer+":"+e.SMTPPort, auth, from, to, msg)
// }
