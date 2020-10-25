package sendEmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type SendMail struct {
	clientID     string
	clientSecret string
	accessToken  string
	refreshToken string
	gmailService *gmail.Service
}

func NewSendMail(clientID, clientSecret, accessToken, refreshToken string) *SendMail {
	return &SendMail{clientID: clientID,
		clientSecret: clientSecret,
		accessToken:  accessToken,
		refreshToken: refreshToken}
}

func (sm *SendMail) Send(to, subject, emailBody string) error {

	if sm.gmailService == nil {
		err := sm.setupOAuthGmailService()
		if err != nil {
			return err
		}
	}

	err := sm.sendEmailOAuth2(to, subject, emailBody)
	if err != nil {
		return err
	}
	return nil
}

func (sm *SendMail) setupOAuthGmailService() error {
	config := oauth2.Config{
		ClientID:     sm.clientID,
		ClientSecret: sm.clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  sm.accessToken,
		RefreshToken: sm.refreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	sm.gmailService = srv
	if sm.gmailService != nil {
		log.Println("Email service is initialized")
	}

	return nil
}

func (sm *SendMail) sendEmailOAuth2(to string, subject string, emailBody string) error {

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subjectLine := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subjectLine + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	if sm.gmailService == nil {
		//log.Println("Unable to send email, cannot get GmailService")
		return fmt.Errorf("Unable to send email, cannot get GmailService")
	}

	_, err := sm.gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return err
	}
	return nil
}
