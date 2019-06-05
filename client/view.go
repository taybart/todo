package main

import (
	"github.com/taybart/log"
	"github.com/zserge/webview"
	"io"
	"net/http"
	"net/url"
)

// Counter is a simple example of automatic Go-to-JS data binding
type Counter struct {
	Value int `json:"value"`
}

// Add increases the value of a counter by n
func (c *Counter) Add(n int) {
	c.Value = c.Value + int(n)
}

// Reset sets the value of a counter back to zero
func (c *Counter) Reset() {
	c.Value = 0
}

var uiFrameworkName = "ReactJS+Babel"

func loadUIFramework(w webview.WebView) {
	// Inject React and Babel
	w.Eval(string(MustAsset("js/react/vendor/babel.min.js")))
	w.Eval(string(MustAsset("js/react/vendor/preact.min.js")))

	// Inject our app code
	w.Eval(fmt.Sprintf(`(function(){
		var n=document.createElement('script');
		n.setAttribute('type', 'text/babel');
		n.appendChild(document.createTextNode("%s"));
		document.body.appendChild(n);
	})()`, template.JSEscapeString(string(MustAsset("js/react/app.jsx")))))

	// Process our code with Babel
	w.Eval(`Babel.transformScriptTags()`)
}

func view(ch chan string) error {
	const myHTML = `<!doctype html><html><body>Hello World!</body></html>`
	w := webview.New(webview.Settings{
		URL: `data:text/html,` + url.PathEscape(myHTML),
	})
	defer w.Exit()

	w.Dispatch(func() {
		// Inject controller
		w.Bind("counter", &Counter{})

		// Inject CSS
		w.InjectCSS(string(MustAsset("js/styles.css")))

		// Inject web UI framework and app UI code
		loadUIFramework(w)
	})
	w.Run()
	w.Run()
	log.Print("test")
	// webview.Open("Hello", "http://"+ln.Addr().String(), 400, 300, false)
	return nil
}
