// the udger package allow you to load in memory and lookup the user agent database to extract value from the provided user agent
package udger

import (
	"database/sql"
	"errors"
	//"fmt"
	"os"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	_ "github.com/mattn/go-sqlite3"
)

// create a new instance of Udger and load all the database in memory to allow fast lookup
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

// lookup one user agent and return a Info struct who contains all the metadata possible for the UA.
func (this *Udger) Lookup(ua string) (*Info, error) {
	info := &Info{}

	browserId, version, err := this.findData(ua, this.rexBrowsers, true)
	if err != nil {
		return nil, err
	}

	info.Browser = this.Browsers[browserId]
	info.Browser.Name = info.Browser.Family + " " + version
	info.Browser.Version = version
	info.Browser.Type = this.browserTypes[info.Browser.typ]

	if val, ok := this.browserOS[browserId]; ok {
		info.OS = this.OS[val]
	} else {
		osId, _, err := this.findData(ua, this.rexOS, false)
		if err != nil {
			return nil, err
		}
		info.OS = this.OS[osId]
	}

	deviceId, _, err := this.findData(ua, this.rexDevices, false)
	if err != nil {
		return nil, err
	}
	if val, ok := this.Devices[deviceId]; ok {
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

func (this *Udger) cleanRegex(r string) string {
	if strings.HasSuffix(r, "/si") {
		r = r[:len(r)-3]
	}
	if strings.HasPrefix(r, "/") {
		r = r[1:]
	}

	return r
}

func (this *Udger) findData(ua string, data []rexData, withVersion bool) (idx int, value string, err error) {
	for i := 0; i < len(data); i++ {
		data[i].Regex = this.cleanRegex(data[i].Regex)
		r, err := pcre.Compile(data[i].Regex, pcre.CASELESS)
		if err != nil {
			return -1, "", errors.New(err.String())
		}
		matcher := r.MatcherString(ua, 0)
		match := matcher.MatchString(ua, 0)
		if !match {
			continue
		}

		if withVersion && matcher.Present(1) {
			return data[i].ID, matcher.GroupString(1), nil
		}

		return data[i].ID, "", nil
	}

	return -1, "", nil
}

func (this *Udger) init() error {
	rows, err := this.db.Query("SELECT client_id, regstring FROM udger_client_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		this.rexBrowsers = append(this.rexBrowsers, d)
	}
	rows.Close()

	rows, err = this.db.Query("SELECT deviceclass_id, regstring FROM udger_deviceclass_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		this.rexDevices = append(this.rexDevices, d)
	}
	rows.Close()

	rows, err = this.db.Query("SELECT os_id, regstring FROM udger_os_regex ORDER by sequence ASC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d rexData
		rows.Scan(&d.ID, &d.Regex)
		this.rexOS = append(this.rexOS, d)
	}
	rows.Close()

	rows, err = this.db.Query("SELECT id, class_id, name,engine,vendor,icon FROM udger_client_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Browser
		var id int
		rows.Scan(&id, &d.typ, &d.Family, &d.Engine, &d.Company, &d.Icon)
		this.Browsers[id] = d
	}
	rows.Close()

	rows, err = this.db.Query("SELECT id, name, family, vendor, icon FROM udger_os_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d OS
		var id int
		rows.Scan(&id, &d.Name, &d.Family, &d.Company, &d.Icon)
		this.OS[id] = d
	}
	rows.Close()

	rows, err = this.db.Query("SELECT id, name, icon FROM udger_deviceclass_list")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Device
		var id int
		rows.Scan(&id, &d.Name, &d.Icon)
		this.Devices[id] = d
	}
	rows.Close()

	rows, err = this.db.Query("SELECT id, client_classification FROM udger_client_class")
	if err != nil {
		return err
	}
	for rows.Next() {
		var d string
		var id int
		rows.Scan(&id, &d)
		this.browserTypes[id] = d
	}
	rows.Close()

	rows, err = this.db.Query("SELECT client_id, os_id FROM udger_client_os_relation")
	if err != nil {
		return err
	}
	for rows.Next() {
		var browser int
		var os int
		rows.Scan(&browser, &os)
		this.browserOS[browser] = os
	}
	rows.Close()

	return nil
}
