package db

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"encoding/json"
	"io"

	"github.com/kgugunava/database/scanner"
	"github.com/kgugunava/database/recorder"
	"github.com/kgugunava/database/models"
)

type Db struct {
	Scanner *scanner.Scanner
	Recorder *recorder.Recorder
	FilePath string
	IdIndex map[int]int64 // id - offset
	NameIndex map[string]map[int]bool
	GpaIndex map[float64]map[int]bool
	ActiveIndex map[bool]map[int]bool // map[int]bool = id: true
	RecordInfo map[int]models.RecordInfo
}

func (db *Db) CreateBackup(backupDir string) error {
    dbPath := db.FilePath

    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        return fmt.Errorf("DB file does not exist: %s", dbPath)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s/backup_%s.jsonl", backupDir, timestamp)

    src, err := os.Open(dbPath)
    if err != nil {
        return fmt.Errorf("error opening DB file: %w", err)
    }
    defer src.Close()

    dst, err := os.Create(backupPath)
    if err != nil {
        return fmt.Errorf("error creating backup file: %w", err)
    }
    defer dst.Close()

    scanner := bufio.NewScanner(src)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue
        }

        var student models.Student
        err := json.Unmarshal([]byte(line), &student)
        if err != nil {
            continue
        }

        if _, exists := db.IdIndex[student.Id]; exists {
            _, err = dst.Write(append([]byte(line), '\n'))
            if err != nil {
                return fmt.Errorf("error writing to backup file: %w", err)
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading DB file: %w", err)
    }

    fmt.Printf("Backup created: %s\n", backupPath)
    return nil
}

func (db *Db) LoadIndex(dbFilePath string) error {
    db.IdIndex = make(map[int]int64)
    db.NameIndex = make(map[string]map[int]bool)
    db.GpaIndex = make(map[float64]map[int]bool)
    db.ActiveIndex = make(map[bool]map[int]bool)
    db.RecordInfo = make(map[int]models.RecordInfo)

    file, err := os.Open(dbFilePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var offset int64 = 0

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            offset++
            continue
        }

        var student models.Student
        err := json.Unmarshal([]byte(line), &student)
        if err != nil {
            offset += int64(len(line)) + 1
            continue
        }

        db.IdIndex[student.Id] = offset

        if db.NameIndex[student.Name] == nil {
            db.NameIndex[student.Name] = make(map[int]bool)
        }
        db.NameIndex[student.Name][student.Id] = true

        if db.GpaIndex[student.Gpa] == nil {
            db.GpaIndex[student.Gpa] = make(map[int]bool)
        }
        db.GpaIndex[student.Gpa][student.Id] = true

        if db.ActiveIndex[student.Active] == nil {
            db.ActiveIndex[student.Active] = make(map[int]bool)
        }
        db.ActiveIndex[student.Active][student.Id] = true

        db.RecordInfo[student.Id] = models.RecordInfo{
            Name:   student.Name,
            Gpa:    student.Gpa,
            Active: student.Active,
        }

        offset += int64(len(line)) + 1
    }

    return scanner.Err()
}

func (db *Db) RestoreFromBackup(backupPath string) error {
    src, err := os.Open(backupPath)
    if err != nil {
        return fmt.Errorf("error opening backup file: %w", err)
    }
    defer src.Close()

    dst, err := os.Create(db.FilePath)
    if err != nil {
        return fmt.Errorf("error creating DB file: %w", err)
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    if err != nil {
        return fmt.Errorf("error copying backup to DB file: %w", err)
    }

    err = db.LoadIndex(db.FilePath)
    if err != nil {
        return fmt.Errorf("error rebuilding indexes: %w", err)
    }

    fmt.Printf("Database restored from: %s\n", backupPath)
    return nil
}