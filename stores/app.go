package stores

import (
	"fmt"
	"strconv"

	"honnef.co/go/js/dom"

	"strings"

	"time"

	"github.com/dave/flux"
	"github.com/gopherjs/gopherjs/js"
)

type App struct {
	Dispatcher flux.DispatcherInterface
	Watcher    flux.WatcherInterface
	Notifier   flux.NotifierInterface

	Archive    *ArchiveStore
	Editor     *EditorStore
	Connection *ConnectionStore
	Local      *LocalStore
	Scanner    *ScannerStore
	Compile    *CompileStore
	Share      *ShareStore
	Request    *RequestStore
	Deploy     *DeployStore
	Page       *PageStore
	Source     *SourceStore
	History    *HistoryStore
}

func (a *App) Init() {

	n := flux.NewNotifier()
	a.Notifier = n
	a.Watcher = n

	a.Archive = NewArchiveStore(a)
	a.Editor = NewEditorStore(a)
	a.Connection = NewConnectionStore(a)
	a.Local = NewLocalStore(a)
	a.Scanner = NewScannerStore(a)
	a.Compile = NewCompileStore(a)
	a.Share = NewShareStore(a)
	a.Request = NewRequestStore(a)
	a.Deploy = NewDeployStore(a)
	a.Page = NewPageStore(a)
	a.Source = NewSourceStore(a)
	a.History = NewHistoryStore(a)

	a.Dispatcher = flux.NewDispatcher(
		// Notifier:
		a.Notifier,
		// Stores:
		a.Archive,
		a.Editor,
		a.Connection,
		a.Local,
		a.Scanner,
		a.Compile,
		a.Share,
		a.Request,
		a.Deploy,
		a.Page,
		a.Source,
		a.History,
	)
}

func (a *App) Dispatch(action flux.ActionInterface) chan struct{} {
	return a.Dispatcher.Dispatch(action)
}

func (a *App) Watch(key interface{}, f func(done chan struct{})) {
	a.Watcher.Watch(key, f)
}

func (a *App) Delete(key interface{}) {
	a.Watcher.Delete(key)
}

func (a *App) Fail(err error) {
	// TODO: improve this
	js.Global.Call("alert", err.Error())
}

func (a *App) Debug(message ...interface{}) {
	js.Global.Get("console").Call("log", message...)
}

var lastLog *struct{}

// LogHide hides the message after 2 seconds
func (a *App) LogHide(args ...interface{}) {
	a.Log(args...)
	if len(args) > 0 {
		// clear message after 2 sec if not changed
		before := lastLog
		go func() {
			<-time.After(time.Second * 2)
			if before == lastLog {
				a.Log()
			}
		}()
	}
}

func (a *App) Log(args ...interface{}) {
	m := dom.GetWindow().Document().GetElementByID("message")
	var message string
	if len(args) > 0 {
		message = strings.TrimSuffix(fmt.Sprintln(args...), "\n")
	}
	if m.InnerHTML() != message {
		if message != "" {
			js.Global.Get("console").Call("log", "Status", strconv.Quote(message))
		}
		requestAnimationFrame()
		m.SetInnerHTML(message)
		requestAnimationFrame()
		lastLog = &struct{}{}
	}
}

func (a *App) Logf(format string, args ...interface{}) {
	a.Log(fmt.Sprintf(format, args...))
}

func (a *App) LogHidef(format string, args ...interface{}) {
	a.LogHide(fmt.Sprintf(format, args...))
}

func requestAnimationFrame() {
	c := make(chan struct{})
	js.Global.Call("requestAnimationFrame", func() { close(c) })
	<-c
}
