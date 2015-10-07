package main

import (
	"fmt"
	"os"

	"github.com/stumpyfr/udger"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:\n\tgo run main.go ./udgerdb.dat \"Opera/9.50 (Nintendo DSi; Opera/507; U; en-US)\"")
		os.Exit(0)
	}

	u, err := udger.New(os.Args[1])
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	fmt.Println(len(u.Browsers), "browsers loaded")
	fmt.Println(len(u.OS), "OS loaded")
	fmt.Println(len(u.Devices), "device types loaded")
	fmt.Println("")

	ua, err := u.Lookup(os.Args[2])
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	fmt.Printf("%+v\n", ua)
}
