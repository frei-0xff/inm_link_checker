package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/therecipe/qt/widgets"
	"mvdan.cc/xurls/v2"
)

func main() {

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(600, 400)
	window.SetWindowTitle("Проверка ссылок vizitka-inmak")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	// create a line edit
	// with a custom placeholder text
	// and add it to the central widgets layout
	// input := widgets.NewQLineEdit(nil)
	input := widgets.NewQTextEdit(nil)
	input.SetPlaceholderText("Вставьте список ссылок ...")

	output := widgets.NewQTextEdit(nil)
	output.SetPlaceholderText("Здесь появятся рабочие ссылки ...")
	output.SetReadOnly(true)

	widget.Layout().AddWidget(input)
	widget.Layout().AddWidget(output)

	links := make(chan string)
	// create a button
	// connect the clicked signal
	// and add it to the central widgets layout
	button := widgets.NewQPushButton2("Проверить!", nil)
	button.ConnectClicked(func(bool) {
		text := input.ToPlainText()
		button.SetEnabled(false)
		output.SetText("")

		rxStrict := xurls.Strict()
		go func() {
			for _, l := range rxStrict.FindAllString(text, -1) {
				links <- l
			}
			button.SetEnabled(true)
		}()
	})
	widget.Layout().AddWidget(button)

	// make the window visible
	window.ShowMaximized()

	var outputMu sync.Mutex
	checkLink := func() {
		for link := range links {
			resp, err := http.Get(link)
			if err == nil && resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil || strings.Index(string(body), "://inmak.com") == -1 {
					continue
				}
				outputMu.Lock()
				output.SetText(output.ToPlainText() + link + "\n")
				output.VerticalScrollBar().SetValue(output.VerticalScrollBar().Maximum())
				outputMu.Unlock()
			}
		}
	}
	go checkLink()
	go checkLink()
	go checkLink()
	go checkLink()
	go checkLink()
	go checkLink()
	go checkLink()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}
