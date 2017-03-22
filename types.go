package udger

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// Udger contains the data and exposes the Lookup(ua string) function
type Udger struct {
	db           *sql.DB
	rexBrowsers  []rexData
	rexDevices   []rexData
	rexOS        []rexData
	browserTypes map[int]string
	browserOS    map[int]int
	Browsers     map[int]Browser
	OS           map[int]OS
	Devices      map[int]Device
}

// Info is the struct returned by the Lookup(ua string) function, contains everything about the UA
type Info struct {
	Browser Browser `json:"browser"`
	OS      OS      `json:"os"`
	Device  Device  `json:"device"`
}

// Browser contains information about the browser type, engine and off course it's name
type Browser struct {
	Name    string `json:"name"`
	Family  string `json:"family"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
	typ     int
	Type    string `json:"type"`
	Company string `json:"company"`
	Icon    string `json:"icon"`
}

type rexData struct {
	ID            int
	Regex         string
	RegexCompiled pcre.Regexp
}

// OS contains all the information about the operating system
type OS struct {
	Name    string `json:"name"`
	Family  string `json:"family"`
	Icon    string `json:"icon"`
	Company string `json:"company"`
}

// Device contains all the information about the device type
type Device struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
