package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/therecipe/qt/widgets"
	"golang.org/x/net/html/charset"
	"mvdan.cc/xurls/v2"
)

const workersCount = 10

type outputTextArea struct {
	mu sync.Mutex
	w  *widgets.QTextEdit
}

var output outputTextArea

func checkLink(links chan string, onlyWithOffers bool) {
	for link := range links {
		resp, err := http.Get(link)
		if err == nil && resp.StatusCode == 200 {
			reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
			if err != nil {
				continue
			}
			body, err := ioutil.ReadAll(reader)
			if err != nil ||
				strings.Index(string(body), "://inmak.com") == -1 ||
				(onlyWithOffers && strings.Index(string(body), "<title>ИнМАК") == -1) {
				continue
			}
			output.mu.Lock()
			output.w.SetPlainText(output.w.ToPlainText() + link + "\n")
			output.w.VerticalScrollBar().SetValue(output.w.VerticalScrollBar().Maximum())
			output.mu.Unlock()
		}
	}
}
func main() {

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 600*400
	// and sets the title
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(600, 400)
	window.SetWindowTitle("Проверка ссылок на vizitka-inmak")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	// create a text edit
	// with a custom placeholder text
	// and add it to the central widgets layout
	input := widgets.NewQTextEdit(nil)
	input.SetPlaceholderText("Вставьте список ссылок ...")

	output.w = widgets.NewQTextEdit(nil)
	output.w.SetPlaceholderText("Здесь появятся рабочие ссылки ...")
	output.w.SetReadOnly(true)

	widget.Layout().AddWidget(input)
	widget.Layout().AddWidget(output.w)

	checkbox := widgets.NewQCheckBox2("Только визитки с коммерческими предложениями", nil)
	widget.Layout().AddWidget(checkbox)

	// create a button
	// connect the clicked signal
	// and add it to the central widgets layout
	button := widgets.NewQPushButton2("Проверить ссылки", nil)
	button.ConnectClicked(func(bool) {
		text := input.ToPlainText()
		onlyWithOffers := checkbox.IsChecked()
		button.SetEnabled(false)
		output.w.SetPlainText("")

		rxStrict := xurls.Strict()
		go func() {
			links := make(chan string)
			var wg sync.WaitGroup
			wg.Add(workersCount)
			for i := 1; i <= workersCount; i++ {
				go func() {
					defer wg.Done()
					checkLink(links, onlyWithOffers)
				}()
			}
			for _, l := range rxStrict.FindAllString(text, -1) {
				links <- l
			}
			close(links)
			wg.Wait()
			lines := strings.Split(strings.TrimSpace(output.w.ToPlainText()), "\n")
			print(lines)
			sort.Strings(lines)
			output.w.SetPlainText(strings.Join(lines, "\n") + "\n")
			output.w.VerticalScrollBar().SetValue(output.w.VerticalScrollBar().Maximum())
			button.SetEnabled(true)
		}()
	})
	widget.Layout().AddWidget(button)

	// make the window visible
	window.ShowMaximized()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}
