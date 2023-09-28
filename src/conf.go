package main

import (
	"fmt"
	"log"
	"os"
)

func setVarEnvs() (string, string, string) {
	mode, mode_set := os.LookupEnv("MODE")
	if !mode_set {
		log.Fatal("MODE environment variable not set")
	}

	pseudo, pseudo_set := os.LookupEnv("PSEUDONAME")
	if !pseudo_set {
		log.Fatal("PSEUDONAME environment variable not set")
	}
	token, token_set := os.LookupEnv("PERSONAL_ACCESS_TOKEN")
	if !token_set {
		log.Println("PERSONAL_ACCESS_TOKEN environment variable not set")
	}

	switch mode {
	case "users":
		fmt.Printf("You are pulling repositories for a USER under this configuration.\n")
		fmt.Printf("Username to pull: %s\n", pseudo)
		fmt.Printf("Personal Token used: %s\n", token)
	case "orgs":
		fmt.Printf("You are pulling repositories for an ORGANIZATION under this configuration.\n")
		fmt.Printf("Organization to pull: %s\n", pseudo)
		fmt.Printf("Personal Token used: %s\n", token)
	default:
		log.Fatal("MODE environment variable not set to either 'users' or 'orgs'")
	}
	return mode, pseudo, token
}
