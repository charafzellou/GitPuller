package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func writeCSV(pseudo string, repositories []Repository, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create folder for CSV files
	err := os.MkdirAll(fmt.Sprintf("../assets/%s", pseudo), 0755)
	if err != nil {
		log.Fatalf("Error creating folder for CSV files: %s", err.Error())
	}

	// Create CSV file
	filename := fmt.Sprintf("../assets/%s/%s.csv", pseudo, time.Now().UTC().Format("2006-01-02-15-04-05"))
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %s", err.Error())
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row to CSV
	header := []string{"Name", "URL", "Private"}
	err = writer.Write(header)
	if err != nil {
		log.Fatalf("Error writing CSV header: %s", err.Error())
	}

	// Write repository details to CSV
	for _, repo := range repositories {
		row := []string{repo.Name, repo.URL, strconv.FormatBool(repo.Private)}
		err = writer.Write(row)
		if err != nil {
			log.Fatalf("Error writing CSV row for %s: %s", repo.Name, err.Error())
		}
	}
	fmt.Printf("CSV file written successfully: %s\n", file.Name())
}
