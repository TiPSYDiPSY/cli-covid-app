package main

import (
	"cli-covid-app/client"
	"flag"
	"fmt"
	"os"
)

var (
	helpFLag    = flag.Bool("help", false, "displays a list of commands")
	countryFlag = flag.String("country", "latvia", "name of country for which the data will be collected")
)

func main() {
	flag.Parse()

	s := client.NewSwitch(*countryFlag)
	if *helpFLag || len(os.Args) == 1 {
		s.Help()
		return
	}

	if len(os.Args) == 3 && os.Args[2] == "--help" {
		if err := s.CommandHelp(); err != nil {
			fmt.Printf("cmd arguments error: %v\n", err)
			os.Exit(2)
		}
		return
	}

	err := s.Switch()
	if err != nil {
		fmt.Printf("cmd switch error: %v\n", err)
		os.Exit(2)
	}
}
