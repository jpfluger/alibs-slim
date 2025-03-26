package aclient_smtp

import (
	"github.com/jhillyerd/enmime"
	"github.com/jpfluger/alibs-slim/aconns"
)

// ISMTPAuth defines an interface for SMTP authentication
type ISMTPAuth interface {
	GetSender() (enmime.Sender, error)
	GetName() aconns.AdapterName
}
