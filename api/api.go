// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"bytes"
	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"html/template"
	"net/http"

	_ "github.com/cloudfoundry/jibber_jabber"
	_ "github.com/nicksnyder/go-i18n/i18n"
)

var ServerTemplates *template.Template

type ServerTemplatePage Page

func NewServerTemplatePage(templateName string) *ServerTemplatePage {
	return &ServerTemplatePage{
		TemplateName: templateName,
		Props:        make(map[string]string),
		ClientCfg:    utils.ClientCfg,
		Locale:       model.DEFAULT_LOCALE,
	}
}

func (me *ServerTemplatePage) Render() string {
	var text bytes.Buffer
	if err := ServerTemplates.ExecuteTemplate(&text, me.TemplateName, me); err != nil {
		l4g.Error(utils.T("api.api.render.error"), me.TemplateName, err)
	}

	return text.String()
}

func InitApi() {
	r := Srv.Router.PathPrefix("/api/v1").Subrouter()
	InitUser(r)
	InitTeam(r)
	InitChannel(r)
	InitPost(r)
	InitWebSocket(r)
	InitFile(r)
	InitCommand(r)
	InitAdmin(r)
	InitOAuth(r)
	InitWebhook(r)
	InitPreference(r)
	InitLicense(r)

	templatesDir := utils.FindDir("api/templates")
	l4g.Debug(utils.T("api.api.init.parsing_templates.debug"), templatesDir)
	var err error
	if ServerTemplates, err = template.ParseGlob(templatesDir + "*.html"); err != nil {
		l4g.Error(utils.T("api.api.init.parsing_templates.error"), err)
	}
}

func HandleEtag(etag string, w http.ResponseWriter, r *http.Request) bool {
	if et := r.Header.Get(model.HEADER_ETAG_CLIENT); len(etag) > 0 {
		if et == etag {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}

	return false
}
