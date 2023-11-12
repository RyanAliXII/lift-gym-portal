package browser

import (
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)




var once sync.Once;
var browser * Browser

type Browser struct {
	b * rod.Browser
	launcher * launcher.Launcher
	page * rod.Page
}
func(b * Browser) Goto (url string)(*rod.Page,error){
	err := b.page.Navigate(url)
	return b.page, err
}
func(b * Browser)GetPage()*rod.Page{
	return b.page
}
func (b * Browser)GetBrowser()*rod.Browser{
	return b.b
}
func (b * Browser)GetLauncher()*rod.Browser{
	return b.b
}

func NewBrowser() (*Browser, error){
	var err error = nil
	once.Do(func ()  {
		path, _ := launcher.LookPath()
		if(path == ""){
			panic("Chromium binary not found.")
		}
		launcher := launcher.New().Headless(true).Bin(path).Leakless(true)
		url, err := launcher.Launch()
		if err != nil {
			return
		}
		b := rod.New().ControlURL(url)
		err = b.Connect()
		if err != nil {
			return
		}
		page, err := b.Page(proto.TargetCreateTarget{})
		if err != nil {
			return 
		}
		browser = &Browser{
			b: b,
			page: page,
			launcher: launcher,
		}
	})
	return browser, err
}