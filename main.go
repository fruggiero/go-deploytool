package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	// Define flags
	token := flag.String("token", "", "override token")
	feedIdentifier := flag.String("feed", "", "feed identifier")
	serverAddress := flag.String("vaultaddr", "", "vault address")

	// Parse flags
	flag.Parse()

	// Check flags
	if *serverAddress == "" {
		*serverAddress = os.Getenv("VAULT_ADDR")
	}

	if *serverAddress == "" {
		fmt.Println("No vault address set")
		fmt.Println("Ensure that VAULT_ADDR environment variable is present, or pass -vaultaddr flag")
		flag.Usage()
	}

	if *token == "" {
		*token = os.Getenv("VAULT_ID_TOKEN")
	}

	if *token == "" {
		fmt.Println("No token set")
		fmt.Println("Ensure that VAULT_ID_TOKEN environment variable is present, or pass -token flag")
		flag.Usage()
		return
	}

	if *feedIdentifier == "" {
		fmt.Println("Error: the flag -feed is mandatory.")
		flag.Usage()
		return
	}

	// Get feed address
	feedAddress, err := getFeedAddress(*feedIdentifier)
	if err != nil {
		fmt.Printf("Key '%s' not found.\n", *feedIdentifier)
		return
	} else {
		fmt.Printf("Using address for feed '%s': %s\n", *feedIdentifier, feedAddress)
	}

	config := vault.DefaultConfig()
	config.Address = *serverAddress

	if config.Address == "" {
		fmt.Println("Environment variable VAULT_ADDR not found or empty")
		return
	}

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken(*token)
}

func getFeedAddress(feedName string) (string, error) {
	// Let's first read the `config.json` file
	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	config, ok := payload["Config"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("key %s not found", "Config")
	}
	feeds, ok := config["Feeds"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("key %s not found", "Feeds")
	}
	feed, ok := feeds[feedName]
	if !ok {
		return "", fmt.Errorf("key %s not found", feedName)
	}

	return feed.(string), nil
}
