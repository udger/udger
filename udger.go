// Package udger package allow you to load in memory and lookup the user agent database to extract value from the provided user agent
package udger

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	_ "github.com/mattn/go-sqlite3"
)

// New creates a new instance of Udger and load all the database in memory to allow fast lookup
// you need to pass the sqlite database in parameter
func New(dbPath string) (*Udger, error) {
	u := &Udger{
		Browsers:     make(map[int]Browser),
		OS:           make(map[int]OS),
		Devices:      make(map[int]Device),
		browserTypes: make(map[int]string),
		browserOS:    make(map[int]int),
	}
	var err error

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, err
	}

	u.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer u.db.Close()

	err = u.init()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Lookup one user agent and return a Info struct who contains all the metadata possible for the UA.
func (udger *Udger) Lookup(ua string) (*Info, error) {
	info := &Info{}

	browserID, version, err := udger.findDataWithVersion(ua, udger.rexBrowsers, true)
	if err != nil {
		return nil, err
	}

	info.Browser = udger.Browsers[browserID]
	info.Browser.Name = info.Browser.Family + " " + version
	info.Browser.Version = version
	info.Browser.Type = udger.browserTypes[info.Browser.typ]

	if val, ok := udger.browserOS[browserID]; ok {
		info.OS = udger.OS[val]
	} else {
		osID, _, err := udger.findData(ua, udger.rexOS, false)
		if err != nil {
			return nil, err
		}
		info.OS = udger.OS[osID]
	}

	deviceID, _, err := udger.findData(ua, udger.rexDevices, false)
	if err != nil {
		return nil, err
	}
	if val, ok := udger.Devices[deviceID]; ok {
		info.Device = val
	} else if info.Browser.typ == 3 { // if browser is mobile, we can guess its a mobile
		info.Device = Device{
			Name: "Smartphone",
			Icon: "phone.png",
		}
	} else if info.Browser.typ == 5 || info.Browser.typ == 10 || info.Browser.typ == 20 || info.Browser.typ == 50 {
		info.Device = Device{
			Name: "Other",
			Icon: "other.png",
		}
	} else {
		//nothing so personal computer
		info.Device = Device{
			Name: "Personal computer",
			Icon: "desktop.png",
		}
	}

	return info, nil
}

func (udger *Udger) cleanRegex(r string) string {
	if strings.HasSuffix(r, "/si") {
		r = r[:len(r)-3]
	}
	if strings.HasPrefix(r, "/") {
		r = r[1:]
	}

	return r
}

func (udger *Udger) findDataWithVersion(ua string, data []rexData, withVersion bool) (idx int, value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			idx, value, err = udger.findData(ua, data, false)
		}
	}()

	idx, value, err = udger.findData(ua, data, withVersion)

	return idx, value, err
}

func (udger *Udger) findData(ua string, data []rexData, withVersion bool) (idx int, value string, err error) {
	for i := 0; i < len(data); i++ {
		r := data[i].RegexCompiled
		matcher := r.MatcherString(ua, 0)
		if !matcher.MatchString(ua, 0) {
			continue
		}

		if withVersion && matcher.Present(1) {
			return data[i].ID, matcher.GroupString(1), nil
		}

		return data[i].ID, "", nil
	}

	return -1, "", nil
}

func (udger *Udger) init() error {
	rows, err := udger.db.Query("SELECT client_id, regstring FROM udger_client_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		d.Regex = udger.cleanRegex(d.Regex)
		r, err := pcre.Compile(d.Regex, pcre.CASELESS)
		if err != nil {
			return errors.New(err.String())
		}
		d.RegexCompiled = r
		udger.rexBrowsers = append(udger.rexBrowsers, d)
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT deviceclass_id, regstring FROM udger_deviceclass_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		d.Regex = udger.cleanRegex(d.Regex)
		r, err := pcre.Compile(d.Regex, pcre.CASELESS)
		if err != nil {
			return errors.New(err.String())
		}
		d.RegexCompiled = r
		udger.rexDevices = append(udger.rexDevices, d)
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT os_id, regstring FROM udger_os_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		d.Regex = udger.cleanRegex(d.Regex)
		r, err := pcre.Compile(d.Regex, pcre.CASELESS)
		if err != nil {
			return errors.New(err.String())
		}
		d.RegexCompiled = r
		udger.rexOS = append(udger.rexOS, d)
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT id, class_id, name,engine,vendor,icon FROM udger_client_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Browser
		var id int
		rows.Scan(&id, &d.typ, &d.Family, &d.Engine, &d.Company, &d.Icon)
		udger.Browsers[id] = d
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT id, name, family, vendor, icon FROM udger_os_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d OS
		var id int
		rows.Scan(&id, &d.Name, &d.Family, &d.Company, &d.Icon)
		udger.OS[id] = d
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT id, name, icon FROM udger_deviceclass_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Device
		var id int
		rows.Scan(&id, &d.Name, &d.Icon)
		udger.Devices[id] = d
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT id, client_classification FROM udger_client_class")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d string
		var id int
		rows.Scan(&id, &d)
		udger.browserTypes[id] = d
	}
	rows.Close()

	rows, err = udger.db.Query("SELECT client_id, os_id FROM udger_client_os_relation")
	if err != nil {
		return err
	}
	for rows.Next() {
		var browser int
		var os int
		rows.Scan(&browser, &os)
		udger.browserOS[browser] = os
	}
	rows.Close()

	return nil
}
