package plugin

import (
	"aista-search/db"
	"html/template"
)

func EpisodeStatus() template.FuncMap {
	f := make(template.FuncMap)

	f["EpisodeStatus"] = func(status uint) string {
		ep := db.Episode{Status: status}
		return ep.GetStatus()
	}

	return f
}
