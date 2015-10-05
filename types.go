package udger

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Udger struct {
	db           *sql.DB
	RexBrowsers  []RexData
	RexDevices   []RexData
	RexOS        []RexData
	Browsers     map[int]Browser
	OS           map[int]OS
	Devices      map[int]Device
	browserTypes map[int]string
	browserOS    map[int]int
}

type Browser struct {
	Name    string `json:"name"`
	Engine  string `json:"engine"`
	typ     int
	Type    string `json:"type"`
	Company string `json:"company"`
	Icon    string `json:"icon"`
}

type RexData struct {
	Id    int
	Regex string
}

type Info struct {
	Browser Browser `json:"browser"`
	OS      OS      `json:"os"`
	Device  Device  `json:"device"`
}

type OS struct {
	Name    string `json:"name"`
	Family  string `json:"family"`
	Icon    string `json:"icon"`
	Company string `json:"compny"`
}

type Device struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
