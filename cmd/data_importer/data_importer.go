package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	_db "github.com/vladwithcode/lex_app/internal/db"
)

type importedCase struct {
	CaseId   string `json:"case_id"`
	CaseType string `json:"case_type"`
}

func main() {
	db, err := _db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	var (
		filePath string
	)
	flag.StringVar(
		&filePath,
		"f",
		"./data.json",
		"Path to the JSON file containing the data to import",
	)
	flag.Parse()

	data, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer data.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(data)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var importedCases []importedCase
	err = json.Unmarshal(buf.Bytes(), &importedCases)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	successful := 0
	for _, c := range importedCases {
		log.Printf("Importing case %s:%s", c.CaseId, c.CaseType)

		newCase, err := _db.NewCase(c.CaseId, c.CaseType)
		if err != nil {
			log.Printf("Invalid case: %v", err)
			continue
		}

		err = _db.InsertCase(context.Background(), db, newCase)
		if err != nil {
			log.Printf("Error inserting case: %s:%s\n  %v", newCase.CaseId, newCase.CaseType, err)
			continue
		}
		successful++
	}
	log.Printf("Successfully imported %d/%d cases", successful, len(importedCases))
}
