package reader

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const LinesPerFile = "200000"

func GetCurrDir() string {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	return path
}

func GetFileList(filename string) []string {
	var path = GetCurrDir()

	var originalPath = filepath.Join(path, filename)
	var tmpPath = filepath.Join(path, "tmp", filename)

	_, err := os.OpenFile(tmpPath+"aa", os.O_RDONLY, 0644)

	if errors.Is(err, os.ErrNotExist) {
		_, err := exec.Command("split", "-l", LinesPerFile, originalPath, tmpPath).Output()

		if err != nil {
			log.Println(err)
		}
	}

	var pattern = tmpPath + "*"
	files, _ := filepath.Glob(pattern)

	return files
}

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.Comma = ' '

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
