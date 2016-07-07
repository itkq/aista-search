package plugin

import (
	"html/template"
	"strings"
)

func ImagePath() template.FuncMap {
	f := make(template.FuncMap)

	f["ImagePath"] = func(s string) string {
		return strings.Replace(s, "./img", "/img", 1)
	}

	return f
}

func ThumbnailPath() template.FuncMap {
	f := make(template.FuncMap)

	f["ThumbnailPath"] = func(s string) string {
		return strings.Replace(s, "./img", "/img/thumb", 1)
	}

	return f
}
