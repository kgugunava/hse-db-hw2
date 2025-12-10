package main

import (
    "github.com/kgugunava/database/db"
    "github.com/kgugunava/gui"
    "github.com/kgugunava/database/recorder"
    "github.com/kgugunava/database/scanner"
    "github.com/kgugunava/database/models"
)

func main() {
    scannerInstance := scanner.NewScanner()

    recorderInstance := recorder.NewRecorder(scannerInstance)
	
    database := &db.Db{
        Scanner:     scannerInstance,
        Recorder:    recorderInstance,
		FilePath:    "input.jsonl",
        IdIndex:     make(map[int]int64),
        NameIndex:   make(map[string]map[int]bool),
        GpaIndex:    make(map[float64]map[int]bool),
        ActiveIndex: make(map[bool]map[int]bool),
        RecordInfo:  make(map[int]models.RecordInfo),
    }

    database.LoadIndex("input.jsonl")

    guiInstance := gui.NewGUI(database)
    guiInstance.Run()
}