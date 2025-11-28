package aclient_smtp

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/aemail"
	"github.com/stretchr/testify/assert"
)

func TestAClientSMTP_Validate(t *testing.T) {
	tests := []struct {
		name    string
		client  AClientSMTP
		wantErr bool
	}{
		{
			name: "Valid client",
			client: AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "name",
					Host: "smtp.example.com",
					Port: 587,
				},
				Username: "user",
				Password: "pass",
				AuthType: AUTHTYPE_SENDER_PLAIN,
			},
			wantErr: false,
		},
		{
			name: "Invalid auth type",
			client: AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "name",
					Host: "smtp.example.com",
					Port: 587,
				},
				Username: "user",
				Password: "pass",
				AuthType: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Missing username",
			client: AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "name",
					Host: "smtp.example.com",
					Port: 587,
				},
				Password: "pass",
				AuthType: AUTHTYPE_SENDER_PLAIN,
			},
			wantErr: true,
		},
		{
			name: "Missing password",
			client: AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "name",
					Host: "smtp.example.com",
					Port: 587,
				},
				Username: "user",
				AuthType: AUTHTYPE_SENDER_PLAIN,
			},
			wantErr: true,
		},
		{
			name: "Missing identity for identity auth type",
			client: AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "name",
					Host: "smtp.example.com",
					Port: 587,
				},
				Username: "user",
				Password: "pass",
				AuthType: AUTHTYPE_SENDER_IDENTITY,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AClientSMTP.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAClientSMTP_CertificateConfigurations(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		port           int
		dialMode       DialMode
		insecure       bool
		expectedToPass bool
	}{
		{
			name:           "Plain SMTP (No TLS)",
			host:           "localhost",
			port:           1025,
			dialMode:       DIALMODE_NOTLS,
			insecure:       false,
			expectedToPass: true,
		},
		//{
		// name:           "SMTP with Valid TLS",
		// host:           "testmail.example.com",
		// port:           1465,
		// dialMode:       DIALMODE_TLS,
		// insecure:       false,
		// expectedToPass: true,
		//},
		//{
		// name:           "SMTP with Self-Signed TLS (Insecure)",
		// host:           "localhost",
		// port:           1565,
		// dialMode:       DIALMODE_TLS,
		// insecure:       true,
		// expectedToPass: true,
		//},
		{
			name:           "SMTP with Self-Signed TLS (Strict)",
			host:           "localhost",
			port:           1565,
			dialMode:       DIALMODE_TLS,
			insecure:       false,
			expectedToPass: false,
		},
		//{
		// name:           "SMTP with Valid STARTTLS",
		// host:           "testmail.example.com",
		// port:           1587,
		// dialMode:       DIALMODE_STARTTLS,
		// insecure:       false,
		// expectedToPass: true,
		//},
		//{
		// name:           "SMTP with Self-Signed STARTTLS (Insecure)",
		// host:           "localhost",
		// port:           1687,
		// dialMode:       DIALMODE_STARTTLS,
		// insecure:       true,
		// expectedToPass: true,
		//},
		{
			name:           "SMTP with Self-Signed STARTTLS (Strict)",
			host:           "localhost",
			port:           1687,
			dialMode:       DIALMODE_STARTTLS,
			insecure:       false,
			expectedToPass: false,
		},
		{
			name:           "Auto-Detect on Plain SMTP",
			host:           "localhost",
			port:           1025,
			dialMode:       DIALMODE_UNKNOWN,
			insecure:       false,
			expectedToPass: true,
		},
		{
			name:           "Force Plain SMTP (NoTLS)",
			host:           "localhost",
			port:           1025,
			dialMode:       DIALMODE_NOTLS,
			insecure:       false,
			expectedToPass: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &AClientSMTP{
				Adapter: aconns.Adapter{
					Type: ADAPTERTYPE_SMTP,
					Name: "mailhog-test",
					Host: tc.host,
					Port: tc.port,
				},
				Username:           "user",
				Password:           "pass",
				AuthType:           AUTHTYPE_SENDER_NONE,
				ConnectionTimeout:  10, // 10-second timeout
				InsecureSkipVerify: tc.insecure,
				DialMode:           tc.dialMode,
			}

			// REMEMBER!
			// Prior to running tests, start `mailhog`, which is imperfect due to TLS errors.
			// docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog
			ok, _, err := client.Test()
			if tc.expectedToPass {
				assert.True(t, ok, "Test should succeed for '%s'", tc.name)
				assert.NoError(t, err, "Test should succeed for '%s'", tc.name)
				assert.True(t, client.HasInitialized(), "Client should be initialized for '%s'", tc.name)
			} else {
				assert.False(t, ok, "Test should fail for '%s'", tc.name)
				assert.Error(t, err, "Test should fail for '%s'", tc.name)
			}

			mailPiece := &MailPiece{
				From:    aemail.Address{Address: "sender@example.com"},
				To:      aemail.Addresses{aemail.Address{Address: "recipient@example.com"}},
				Subject: "Test Email",
				Text:    "This is a test email sent to MailHog.",
			}

			if tc.expectedToPass && err == nil {
				err = client.SendMail(mailPiece)
				assert.NoError(t, err, "SendMail should succeed for '%s'", tc.name)
			}
		})
	}
}
