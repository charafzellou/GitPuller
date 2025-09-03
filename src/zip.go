package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func zipRepositoriesFolder(pseudo string) {
	// Set input_folder to the name of the folder you wish to zip
	input_folder := fmt.Sprintf("../assets/%s", pseudo)
	output_zip := fmt.Sprintf("../assets/%s.zip", pseudo)

	// Create a new ZIP file
	zipFile, err := os.Create(output_zip)
	if err != nil {
		log.Fatalf("Error creating ZIP file: %v\n", err)
	}
	defer zipFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the folder and add files to the ZIP
	err = filepath.Walk(input_folder, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("filepath.Walk() failed : %s", err)
		}

		// Create a new file header
		fileHeader, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			log.Fatalf("zip.FileInfoHeader() failed : %s", err)
		}

		// Set the name of the file within the ZIP archive
		// fileHeader.Name, err = filepath.Rel(input_folder, filePath)
		fileHeader.Name = filepath.ToSlash(filePath[len(input_folder):])
		if err != nil {
			log.Fatalf("filepath.Rel() failed : %s", err)
		}

		// Create a writer for the file within the ZIP archive
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			log.Fatalf("zipWriter.CreateHeader() failed : %s", err)
		}

		// check if the file is not a directory
		if !fileInfo.IsDir() {
			// Open and copy the file's contents to the ZIP archive
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatalf("os.Open() failed : %s", err)
			}
			defer file.Close()

			// Copy contents of the file to the ZIP file
			_, err = io.Copy(writer, file)
			if err != nil {
				log.Fatalf("io.Copy() failed : %s", err)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error zipping folder : %v\n", err)
		return
	}

	fmt.Printf("Folder '%s' zipped successfully to '%s'\n", input_folder, output_zip)
}
