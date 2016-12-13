package backends

import (
	"fmt"
	"net/mail"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hatchitup/go-email/email"
	guerrilla "github.com/hatchitup/go-guerrilla"
	enmime "github.com/jhillyerd/go.enmime"
	"github.com/laurent22/ical-go/ical"
)

func init() {
	backends["dummy"] = &DummyBackend{}
}

type DummyBackend struct {
	config dummyConfig
}

type dummyConfig struct {
	LogReceivedMails bool `json:"log_received_mails"`
}

func (b *DummyBackend) loadConfig(backendConfig guerrilla.BackendConfig) error {
	var converted bool
	b.config.LogReceivedMails, converted = backendConfig["log_received_mails"].(bool)
	if !converted {
		return fmt.Errorf("failed to load backend config (%v)", backendConfig)
	}
	return nil
}

func (b *DummyBackend) Initialize(backendConfig guerrilla.BackendConfig) error {
	return b.loadConfig(backendConfig)
}

func (b *DummyBackend) Finalize() error {
	return nil
}

func (b *DummyBackend) Process(client *guerrilla.Client, from *guerrilla.EmailParts, to []*guerrilla.EmailParts) string {
	if b.config.LogReceivedMails {

		reader := strings.NewReader(client.Data)
		msg, _ := email.ParseMessage(reader)

		Email, _ := mail.ReadMessage(strings.NewReader(client.Data))
		MIME, _ := enmime.ParseMIMEBody(Email)

		// Headers are in the net/mail Message
		log.Infof("From: %v\n", Email.Header.Get("From"))

		log.Infof("From: %v\n", MIME.GetHeader("From"))

		// enmime can decode quoted-printable headers
		log.Infof("Subject: %v\n", MIME.GetHeader("Subject"))

		// The plain text body is available as MIME.Text
		log.Infof("Text Body: %v chars\n", len(MIME.Text))

		log.Infof("Text Body: %s\n", MIME.Text)

		// The HTML body is stored in MIME.HTML
		log.Infof("HTML Body: %v chars\n", len(MIME.HTML))

		// MIME.Inlines is a slice of inlined attacments
		log.Infof("Inlines: %v\n", len(MIME.Inlines))

		// MIME.Attachments contains the non-inline attachments
		log.Infof("Attachments: %v\n", len(MIME.Attachments))

		//log.Infof("Body: %v\n", msg.Body)

		for _, part := range msg.PartsContentTypePrefix("text/html") {
			log.Infof("Part: %v\n", string(part.Body))
		}
		for _, part := range msg.PartsContentTypePrefix("text/calendar") {
			cal, _ := ical.ParseCalendar(string(part.Body))
			log.Infof("Calendar: %s\n", cal)
		}

	}
	return fmt.Sprintf("250 OK : queued as %s", client.Hash)
}
