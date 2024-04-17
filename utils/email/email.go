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

	t, err := template.ParseFiles("utils/email/sighting.html")
	if err != nil {
		return err
	}

	var o bytes.Buffer
	err = t.ExecuteTemplate(&o, "sighting", s)
	if err != nil {
		return err
	}

	html := o.String()

	fmt.Println(html)

	message := mail.NewSingleEmail(from, subject, to, plain, html)

	res, err := c.sg.Send(message)
	if err != nil {
		return err
	}

	log.Println(res.StatusCode)
	log.Println(res.Body)
	log.Println(res.Headers)
	return nil
}

func (c *EmailClient) TestEmail() {
	from := mail.NewEmail("Example User", "tigerhall-kittens@mwyndham.dev")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "kazeam.plus@gmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	res, err := c.sg.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(res.StatusCode)
		fmt.Println(res.Body)
		fmt.Println(res.Headers)
	}
}
