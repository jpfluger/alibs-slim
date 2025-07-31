package alegal

import (
	"github.com/jpfluger/alibs-slim/acontact"
	"github.com/jpfluger/alibs-slim/aemail"
)

type LegalOperatorView struct {
	CompanyName   string
	LegalName     string
	BusinessURL   string
	LegalURL      string
	BusinessLabel string

	LegalEmail aemail.EmailAddress
	LegalPhone string
	LegalMail  []string // ⬅ text-formatted multiline mailing address
}

func NewLegalOperatorView(lo *LegalOperator) *LegalOperatorView {
	if lo == nil {
		return nil
	}

	view := &LegalOperatorView{
		CompanyName:   lo.Name.Company,
		LegalName:     lo.Name.Legal,
		BusinessLabel: lo.Name.MustGetShort(),
	}

	// URLs
	if u := lo.Urls.FindByType(acontact.URLTYPE_BUSINESS); u != nil && u.Link != nil {
		view.BusinessURL = u.GetLinkWithOptions("raw")
	}
	if u := lo.Urls.FindByType(acontact.URLTYPE_LEGAL); u != nil && u.Link != nil {
		view.LegalURL = u.GetLinkWithOptions("raw")
	}

	// Email
	for _, e := range lo.Emails {
		if e.Type == acontact.EMAILTYPE_LEGAL {
			view.LegalEmail = e.Address
			break
		}
	}

	// Phone
	for _, p := range lo.Phones {
		if p.Type == acontact.PHONETYPE_LEGAL {
			view.LegalPhone = p.Number
			break
		}
	}

	// Mailing address
	for _, m := range lo.Mails {
		if m.Type == acontact.MAILTYPE_LEGAL {
			view.LegalMail = m.Address.ToLines()
			break
		}
	}

	return view
}
