package scanner

import (
	"log"
	"os"
	"bufio"
	"encoding/json"

	"github.com/kgugunava/database/models"
)

type Scanner struct {
	
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) ReadFileInList(fileName string) [][]byte {
	var res [][]byte

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error while openning file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Bytes())
	}

	return res
}

func (s *Scanner) ParseJson(data [][]byte) []models.Student {
	var res []models.Student
	for _, value := range data {
		var curStudent models.Student
		if err := json.Unmarshal(value, &curStudent); err != nil {
			log.Fatal("Error while parsing json: ", err)
		}
		res = append(res, curStudent)
	}
	return res
}