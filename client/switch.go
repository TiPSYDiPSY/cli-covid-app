package client

import (
	"fmt"
	"os"

	"github.com/guptarohit/asciigraph"
)

type Switch struct {
	countryName  string
	commands     map[string]func() func(string, []CountryData) error
	commandsHelp map[string]string
}

func NewSwitch(countryName string) Switch {
	s := Switch{countryName: countryName}
	s.commands = map[string]func() func(string, []CountryData) error{
		"confirmed": s.confirmed,
		"recovered": s.recovered,
		"deaths":    s.death,
		"active":    s.active,
	}

	s.commandsHelp = map[string]string{
		"confirmed": "Data about confirmed cases in <country>",
		"recovered": "Data about recovered cases in <country>",
		"deaths":    "Data about deaths cases in <country>",
		"active":    "Data about active cases in <country>",
	}

	return s
}

func (s Switch) Switch() error {
	httpClient := NewHTTPClient(s.countryName)
	data, err := httpClient.GetData()
	if err != nil {
		return err
	}

	cmdName := os.Args[len(os.Args)-1]
	cmd, ok := s.commands[cmdName]
	if !ok {
		return fmt.Errorf("invalid command '%s'", cmdName)
	}
	return cmd()(cmdName, data)
}

func (s Switch) Help() {
	fmt.Println("Usage of covid-cli-app:\n-country <name> <commands> [<args>]")
	for name := range s.commands {
		fmt.Printf("\t\t%-10s --help\n", name)
	}
	fmt.Println("Run 'app COMMAND --help' for more information on a command.")
}

func (s Switch) confirmed() func(string, []CountryData) error {
	return func(cmd string, countriesData []CountryData) error {
		if err := s.checkArgs(cmd); err != nil {
			return err
		}
		var data []float64
		for _, countryData := range countriesData {
			data = append(data, countryData.Confirmed/10000)
		}

		graph := asciigraph.Plot(data)
		fmt.Println(graph)
		return nil
	}
}

func (s Switch) recovered() func(string, []CountryData) error {
	return func(cmd string, countriesData []CountryData) error {
		if err := s.checkArgs(cmd); err != nil {
			return err
		}

		var data []float64
		for _, countryData := range countriesData {
			data = append(data, countryData.Recovered/10000)
		}

		graph := asciigraph.Plot(data)
		fmt.Println(graph)
		return nil
	}
}

func (s Switch) death() func(string, []CountryData) error {
	return func(cmd string, countriesData []CountryData) error {
		if err := s.checkArgs(cmd); err != nil {
			return err
		}
		var data []float64
		for _, countryData := range countriesData {
			data = append(data, countryData.Deaths/10000)
		}

		graph := asciigraph.Plot(data)
		fmt.Println(graph)
		return nil
	}
}

func (s Switch) active() func(string, []CountryData) error {
	return func(cmd string, countriesData []CountryData) error {
		if err := s.checkArgs(cmd); err != nil {
			return err
		}
		var data []float64
		for _, countryData := range countriesData {
			data = append(data, countryData.Active/10000)
		}

		graph := asciigraph.Plot(data)
		fmt.Println(graph)
		return nil
	}
}

func (s Switch) CommandHelp() error {
	cmdHelp, ok := s.commandsHelp[os.Args[1]]
	if !ok {
		return fmt.Errorf("invalid command '%s'", os.Args[1])
	}
	fmt.Printf("Usage of %s: %s\n", os.Args[1], cmdHelp)
	return nil
}

func (s Switch) checkArgs(cmd string) error {
	if len(os.Args) != 4 {
		return fmt.Errorf("%s expects 4 arg(s), %d provided", cmd, len(os.Args))
	}
	return nil
}
