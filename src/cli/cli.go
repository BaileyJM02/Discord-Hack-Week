package cli

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
)

// Options struct for server creation
type Options struct {
	Token string
	UserID string
	ServerName string
	ServerType string
}

// Start function to start CLI
func Start() (*Options, error) {
	// make options map
	options := make(map[string]string)

	// define all validations
	validateToken := func(input string) error {
		r, _ := regexp.Compile(`([N|M][a-zA-Z\d-_]{23}[.][a-zA-Z\d-_]{6}[.][a-zA-Z\d-_]{27})`)
		valid := r.MatchString(input)

		if valid != true {
			return errors.New("Invalid token, please use a bot token that can be found at https://discordapp.com/developers/applications/")
		}
		return nil
	}

	validateID := func(input string) error {
		// Is it 18?
		if len(input) < 18 {
			return errors.New("Invalid ID, please enter your user ID")
		}
		return nil
	}

	validateString := func(input string) error {
		if input == "" {
			return errors.New("Invalid string, please type a string")
		}
		return nil
	}

	// define first prompt - asking for bot token
	prompt := promptui.Prompt{
		Label:    "Token",
		Validate: validateToken,
	}

	// start first prompt
	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return &Options{}, err
	}

	// pass the result to the map
	options["token"] = result

	// define user id prompt
	prompt = promptui.Prompt{
		Label:    "User ID",
		Validate: validateID,
	}

	result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return &Options{}, err
	}

	// pass the result to the map
	options["userID"] = result

	// define server types
	serverTypes := []string{"Bot & Support", "Support", "Fun", "Project", "Product / Service"}
	index := -1

	// populate prompt select
	for index < 0 {
		prompt := promptui.Select{
			Label:    "What server type do you want to create?",
			Items:    serverTypes,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			serverTypes = append(serverTypes, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return &Options{}, err
	}

	// pass the result to the map
	options["serverType"] = result

	// define server name prompt
	prompt = promptui.Prompt{
		Label:    "Server Name",
		Validate: validateString,
	}

	result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return &Options{}, err
	}

	// pass the result to the map
	options["serverName"] = result


	// return map with error nil
	return &Options{
		Token: options["token"],
		UserID: options["userID"],
		ServerName: options["serverName"],
		ServerType: options["serverType"],
	}, nil
}