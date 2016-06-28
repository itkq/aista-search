package plugin

import (
	"html/template"
	"time"
)

func FormattedTime() template.FuncMap {
	f := make(template.FuncMap)

	f["FormatTime"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}

	return f
}
