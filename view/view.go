// source: https://github.com/josephspurrier/gowebapp/blob/master/vendor/app/shared/server/server.go
package view

import (
	"aista-search/session"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	childTemplates     []string
	rootTemplate       string
	templateCollection = make(map[string]*template.Template)
	pluginCollection   = make(template.FuncMap)
	mutexPlugins       sync.RWMutex
	mutex              sync.RWMutex
	sessionName        string
	viewInfo           View
)

// Template root and children
type Template struct {
	Root     string   `json:"Root"`
	Children []string `json:"Children"`
}

type View struct {
	BaseURI   string
	Caching   bool
	Folder    string
	Name      string
	Extension string
	Vars      map[string]interface{}
	context   *gin.Context
}

type Flash struct {
	Message string
	Class   string
}

func init() {
	// Magic goes here to allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// Source: http://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(Flash{})
}

func Configure() {
	viewInfo = View{
		BaseURI:   "/",
		Folder:    "templates",
		Extension: "tmpl",
	}
	rootTemplate = "base"
}

// LoadTemplates will set the root and child templates
func LoadTemplates(rootTemp string, childTemps []string) {
	rootTemplate = rootTemp
	childTemplates = childTemps
}

// LoadPlugins will combine all template.FuncMaps into one map and then set the
// plugins for the templates
// If a func already exists, it is rewritten, there is no error
func LoadPlugins(fms ...template.FuncMap) {
	// Final FuncMap
	fm := make(template.FuncMap)

	// Loop through the maps
	for _, m := range fms {
		// Loop through each key and value
		for k, v := range m {
			fm[k] = v
		}
	}

	// Load the plugins
	mutexPlugins.Lock()
	pluginCollection = fm
	mutexPlugins.Unlock()
}

func New(c *gin.Context) *View {
	v := &View{}
	v.Vars = make(map[string]interface{})
	v.Folder = viewInfo.Folder
	v.Extension = viewInfo.Extension
	v.context = c

	return v
}

func (v *View) Render() {
	c := v.context

	// Get the template collection from cache
	mutex.RLock()
	tc, ok := templateCollection[v.Name]
	mutex.RUnlock()

	// Get the plugin collection
	mutexPlugins.RLock()
	pc := pluginCollection
	mutexPlugins.RUnlock()

	// If the template collection is not cached or caching is disabled
	if !ok || !viewInfo.Caching {

		// List of template names
		var templateList []string
		templateList = append(templateList, rootTemplate)
		templateList = append(templateList, v.Name)
		templateList = append(templateList, childTemplates...)

		// Loop through each template and test the full path
		for i, name := range templateList {
			// Get the absolute path of the root template
			path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
			if err != nil {
				http.Error(c.Writer, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			templateList[i] = path
		}

		// Determine if there is an error in the template syntax
		templates, err := template.New(v.Name).Funcs(pc).ParseFiles(templateList...)

		if err != nil {
			http.Error(c.Writer, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Cache the template collection
		mutex.Lock()
		templateCollection[v.Name] = templates
		mutex.Unlock()

		// Save the template collection
		tc = templates
	}

	session := session.Instance(c.Request)

	if flashes := session.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]Flash, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				v.Vars["flashes"].([]Flash)[i] = f.(Flash)
			default:
				v.Vars["flashes"].([]Flash)[i] = Flash{f.(string), "alert"}
			}

		}
		session.Save(c.Request, c.Writer)
	}

	err := tc.ExecuteTemplate(c.Writer, rootTemplate+"."+v.Extension, v.Vars)

	if err != nil {
		http.Error(c.Writer, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}
}
