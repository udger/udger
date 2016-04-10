package udger_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/udger/udger"
)

func TestInvalidDbName(t *testing.T) {
	Convey("load invalid path", t, func() {
		udger, err := udger.New("./toto.dat")
		So(err, ShouldNotBeNil)
		So(udger, ShouldBeNil)
	})
}

func TestValidDbName(t *testing.T) {
	Convey("load valid path", t, func() {
		udger, err := udger.New("./udgerdb_v3.dat")
		So(err, ShouldBeNil)
		So(udger, ShouldNotBeNil)

		Convey("test memory database", func() {
			So(len(udger.Browsers), ShouldBeGreaterThan, 0)
			So(len(udger.Devices), ShouldBeGreaterThan, 0)
			So(len(udger.OS), ShouldBeGreaterThan, 0)

			Convey("test lookup MAC", func() {
				info, err := udger.Lookup("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2575.0 Safari/537.36")
				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)

				Convey("test lookup info", func() {
					So(info.OS.Company, ShouldResemble, "Apple Computer, Inc.")
					So(info.OS.Family, ShouldResemble, "OS X")
					So(info.OS.Icon, ShouldResemble, "macosx.png")
					So(info.OS.Name, ShouldResemble, "OS X 10.11 El Capitan")

					So(info.Device.Name, ShouldResemble, "Personal computer")
					So(info.Device.Icon, ShouldResemble, "desktop.png")

					So(info.Browser.Company, ShouldResemble, "Google Inc.")
					So(info.Browser.Engine, ShouldResemble, "WebKit/Blink")
					So(info.Browser.Family, ShouldResemble, "Chrome")
					So(info.Browser.Icon, ShouldResemble, "chrome.png")
					So(info.Browser.Name, ShouldResemble, "Chrome 49.0.2575.0")
					So(info.Browser.Type, ShouldResemble, "Browser")
					So(info.Browser.Version, ShouldResemble, "49.0.2575.0")
				})
			})

			Convey("test lookup Win", func() {
				info, err := udger.Lookup("Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0)")
				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)

				Convey("test lookup info", func() {
					So(info.OS.Company, ShouldResemble, "Microsoft Corporation.")
					So(info.OS.Family, ShouldResemble, "Windows")
					So(info.OS.Icon, ShouldResemble, "windows-7.png")
					So(info.OS.Name, ShouldResemble, "Windows 7")

					So(info.Device.Name, ShouldResemble, "Personal computer")
					So(info.Device.Icon, ShouldResemble, "desktop.png")

					So(info.Browser.Company, ShouldResemble, "Microsoft Corporation.")
					So(info.Browser.Engine, ShouldResemble, "Trident")
					So(info.Browser.Family, ShouldResemble, "IE")
					So(info.Browser.Icon, ShouldResemble, "msie.png")
					So(info.Browser.Name, ShouldResemble, "IE 8.0")
					So(info.Browser.Type, ShouldResemble, "Browser")
					So(info.Browser.Version, ShouldResemble, "8.0")
				})
			})

			Convey("test Nintendo DS", func() {
				info, err := udger.Lookup("Opera/9.50 (Nintendo DSi; Opera/507; U; en-US)")
				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)

				Convey("test lookup info", func() {
					So(info.OS.Company, ShouldResemble, "Nintendo of America Inc.")
					So(info.OS.Family, ShouldResemble, "Nintendo")
					So(info.OS.Icon, ShouldResemble, "nintendoDS.png")
					So(info.OS.Name, ShouldResemble, "Nintendo DS")

					So(info.Device.Name, ShouldResemble, "Game console")
					So(info.Device.Icon, ShouldResemble, "console.png")

					So(info.Browser.Company, ShouldResemble, "Opera Software ASA.")
					So(info.Browser.Engine, ShouldResemble, "Presto/Blink")
					So(info.Browser.Family, ShouldResemble, "Opera")
					So(info.Browser.Icon, ShouldResemble, "opera.png")
					So(info.Browser.Name, ShouldResemble, "Opera 9.50")
					So(info.Browser.Type, ShouldResemble, "Browser")
					So(info.Browser.Version, ShouldResemble, "9.50")
				})
			})
		})
	})
}
