package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailClient struct {
	sg          *sendgrid.Client
	senderEmail string
}

func NewEmailClient() *EmailClient {
	sg := sendgrid.NewSendClient(config.Get(config.SENDGRID_API_KEY))
	return &EmailClient{sg, config.Get(config.SENDGRID_SENDER_EMAIL)}
}

type SightingEmail struct {
	DestinationEmail  string
	TigerName         string
	SightingDate      string
	SightingLatitude  string
	SightingLongitude string
	ImageURL          string
}

func (c *EmailClient) SendSightingEmail(s *SightingEmail) error {
	from := mail.NewEmail("Tigerhall Kittens", c.senderEmail)
	subject := fmt.Sprintf("New Sightings for %s the Tiger!", s.TigerName)
	to := mail.NewEmail("Example User", s.DestinationEmail)

	html, err := c.RenderHTMLStr(s)
	if err != nil {
		return err
	}

	plain := c.RenderPlainStr(s)

	message := mail.NewSingleEmail(from, subject, to, plain, html)

	_, err = c.sg.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func (c *EmailClient) RenderPlainStr(s *SightingEmail) string {
	plain := fmt.Sprintf(`
	  New Sighting for %s the Tiger Confirmed!
	  Tiger Name: %s
	  Sighting Date: %s
	  Location
	  Latitude: %s
	  Longitude: %s
	  `,
		s.TigerName,
		s.TigerName,
		s.SightingDate,
		s.SightingLatitude,
		s.SightingLongitude,
	)
	return plain
}

func (c *EmailClient) RenderHTMLStr(s *SightingEmail) (string, error) {
	t, err := template.ParseFiles("utils/email/sighting.html")
	if err != nil {
		return "", err
	}

	var o bytes.Buffer
	err = t.ExecuteTemplate(&o, "sighting", s)
	if err != nil {
		return "", err
	}

	html := o.String()

	return html, nil
}

func (c *EmailClient) QueueConsumer(msg <-chan SightingEmail) {
	for m := range msg {
		err := c.SendSightingEmail(&m)
		if err != nil {
			log.Println(err)
		}
	}
}
