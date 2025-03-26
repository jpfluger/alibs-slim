package ahttp

import (
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/jpfluger/alibs-slim/azb"
)

type NewPageData func(activeURL string, title string, us asessions.ILoginSessionPerm, data interface{}) map[string]interface{}
type NewPageDataPaginate func(activeURL string, title string, us asessions.ILoginSessionPerm, paginate azb.IPaginate, data interface{}) map[string]interface{}

func NewRHPageData(activeURL string, title string, us asessions.ILoginSessionPerm, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"URL_ACTIVE": activeURL,
		"TITLE":      title,
		"US":         us,
		"DATA":       data,
	}
}

func NewRHPageDataPaginate(activeURL string, title string, us asessions.ILoginSessionPerm, paginate azb.IPaginate, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"URL_ACTIVE": activeURL,
		"TITLE":      title,
		"US":         us,
		"PAGINATE":   paginate,
		"DATA":       data,
	}
}
