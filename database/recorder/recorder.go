package recorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"

	"github.com/xuri/excelize/v2"

	"github.com/kgugunava/database/models"
	"github.com/kgugunava/database/scanner"
)

type Recorder struct {
	Scanner *scanner.Scanner
}

func NewRecorder(scanner *scanner.Scanner) *Recorder {
	return &Recorder{Scanner: scanner}
}

func (r *Recorder) MakeNewRecord(data models.Student) models.Record {
	return models.Record{
		Id:      1,
		Student: data,
	}
}

func (r *Recorder) MakeRecordsFromList(data []models.Student) []models.Record {
	var records []models.Record
	for _, value := range data {
		records = append(records, r.MakeNewRecord(value))
	}
	return records
}

func (r *Recorder) AddNewRecord(record models.Record, dbFilePath string, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
	if _, exists := idIndex[record.Student.Id]; exists {
		err := errors.New("record with this id already exists")
		log.Fatal("Record with this id already exists: ", err)
		return err
	}

	file, err := os.OpenFile(dbFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error while openning file for add new record: ", err)
		return err
	}
	defer file.Close()

	data, err := json.Marshal(record.Student)
	if err != nil {
		log.Fatal("Error while marshalling json new record: ", err)
		return err
	}

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		log.Fatal("Error while writing new record to file: ", err)
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("Error while getting file info: ", err)
		return err
	}

	offset := fileInfo.Size() - int64(len(data)) - 1
	idIndex[record.Student.Id] = offset
	if nameIndex[record.Student.Name] == nil {
		nameIndex[record.Student.Name] = make(map[int]bool)
	}
	nameIndex[record.Student.Name][record.Student.Id] = true

	if gpaIndex[record.Student.Gpa] == nil {
		gpaIndex[record.Student.Gpa] = make(map[int]bool)
	}
	gpaIndex[record.Student.Gpa][record.Student.Id] = true

	if activeIndex[record.Student.Active] == nil {
		activeIndex[record.Student.Active] = make(map[int]bool)
	}
	activeIndex[record.Student.Active][record.Student.Id] = true

	recordInfo[record.Student.Id] = models.RecordInfo{
		Name:   record.Student.Name,
		Gpa:    record.Student.Gpa,
		Active: record.Student.Active,
	}

	return nil
}

func (r *Recorder) AddNewRecordsFromList(records []models.Record, dbFilePath string, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) {
	for _, record := range records {
		r.AddNewRecord(record, dbFilePath, idIndex, nameIndex, gpaIndex, activeIndex, recordInfo)
	}
}

func (r *Recorder) DeleteRecordById(id int, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
	info, exists := recordInfo[id]
	if !exists {
		return fmt.Errorf("record with ID %d does not exist", id)
	}

	delete(idIndex, id)

	delete(nameIndex[info.Name], id)
	if len(nameIndex[info.Name]) == 0 {
		delete(nameIndex, info.Name)
	}

	delete(gpaIndex[info.Gpa], id)
	if len(gpaIndex[info.Gpa]) == 0 {
		delete(gpaIndex, info.Gpa)
	}

	delete(activeIndex[info.Active], id)
	if len(activeIndex[info.Active]) == 0 {
		delete(activeIndex, info.Active)
	}

	delete(recordInfo, id)

	return nil
}

func (r *Recorder) DeleteRecordByName(name string, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
	ids, exists := nameIndex[name]
	if !exists {
		return fmt.Errorf("no records found with name: %s", name)
	}

	for id := range ids {
		delete(idIndex, id)

		info := recordInfo[id]
		delete(gpaIndex[info.Gpa], id)
		if len(gpaIndex[info.Gpa]) == 0 {
			delete(gpaIndex, info.Gpa)
		}

		delete(activeIndex[info.Active], id)
		if len(activeIndex[info.Active]) == 0 {
			delete(activeIndex, info.Active)
		}

		delete(recordInfo, id)
	}

	delete(nameIndex, name)

	return nil
}

func (r *Recorder) DeleteRecordByGpa(gpa float64, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
    ids, exists := gpaIndex[gpa]
    if !exists {
        return fmt.Errorf("no records found with GPA: %f", gpa)
    }

    for id := range ids {
        delete(idIndex, id)

        info := recordInfo[id]
        delete(nameIndex[info.Name], id)
        if len(nameIndex[info.Name]) == 0 {
            delete(nameIndex, info.Name)
        }

        delete(activeIndex[info.Active], id)
        if len(activeIndex[info.Active]) == 0 {
            delete(activeIndex, info.Active)
        }

        delete(recordInfo, id)
    }

    delete(gpaIndex, gpa)

    return nil
}


func (r *Recorder) DeleteRecordByActive(active bool, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
    ids, exists := activeIndex[active]
    if !exists {
        return fmt.Errorf("no records found with active: %t", active)
    }

    for id := range ids {
        delete(idIndex, id)

        info := recordInfo[id]
        delete(nameIndex[info.Name], id)
        if len(nameIndex[info.Name]) == 0 {
            delete(nameIndex, info.Name)
        }

        delete(gpaIndex[info.Gpa], id)
        if len(gpaIndex[info.Gpa]) == 0 {
            delete(gpaIndex, info.Gpa)
        }

        delete(recordInfo, id)
    }

    delete(activeIndex, active)

    return nil
}

func (r *Recorder) EditRecord(newRecord models.Record, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
    id := newRecord.Student.Id

    oldInfo, exists := recordInfo[id]
    if !exists {
        return fmt.Errorf("record with ID %d does not exist", id)
    }

    delete(nameIndex[oldInfo.Name], id)
    if len(nameIndex[oldInfo.Name]) == 0 {
        delete(nameIndex, oldInfo.Name)
    }

    delete(gpaIndex[oldInfo.Gpa], id)
    if len(gpaIndex[oldInfo.Gpa]) == 0 {
        delete(gpaIndex, oldInfo.Gpa)
    }

    delete(activeIndex[oldInfo.Active], id)
    if len(activeIndex[oldInfo.Active]) == 0 {
        delete(activeIndex, oldInfo.Active)
    }

    newInfo := models.RecordInfo{
        Name:   newRecord.Student.Name,
        Gpa:    newRecord.Student.Gpa,
        Active: newRecord.Student.Active,
    }

    if nameIndex[newInfo.Name] == nil {
        nameIndex[newInfo.Name] = make(map[int]bool)
    }
    nameIndex[newInfo.Name][id] = true

    if gpaIndex[newInfo.Gpa] == nil {
        gpaIndex[newInfo.Gpa] = make(map[int]bool)
    }
    gpaIndex[newInfo.Gpa][id] = true

    if activeIndex[newInfo.Active] == nil {
        activeIndex[newInfo.Active] = make(map[int]bool)
    }
    activeIndex[newInfo.Active][id] = true

    recordInfo[id] = newInfo

    return nil
}

func (r *Recorder) FindById(id int, idIndex map[int]int64, dbFilePath string) (*models.Student, error) {
    offset, exists := idIndex[id]
    if !exists {
        return nil, fmt.Errorf("record with ID %d not found", id)
    }

    file, err := os.Open(dbFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    _, err = file.Seek(offset, 0)
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    line := scanner.Text()

    var student models.Student
    err = json.Unmarshal([]byte(line), &student)
    if err != nil {
        return nil, err
    }

    return &student, nil
}

func (r *Recorder) FindByName(name string, nameIndex map[string]map[int]bool, idIndex map[int]int64, dbFilePath string) ([]models.Student, error) {
    ids, exists := nameIndex[name]
    if !exists {
        return nil, fmt.Errorf("no records found with name: %s", name)
    }

    var results []models.Student
    file, err := os.Open(dbFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    for id := range ids {
        offset, exists := idIndex[id]
        if !exists {
            continue
        }

        _, err = file.Seek(offset, 0)
        if err != nil {
            continue
        }

        scanner := bufio.NewScanner(file)
        scanner.Scan()
        line := scanner.Text()

        var student models.Student
        err := json.Unmarshal([]byte(line), &student)
        if err != nil {
            continue
        }

        results = append(results, student)
    }

    return results, nil
}

func (r *Recorder) FindByGpa(gpa float64, gpaIndex map[float64]map[int]bool, idIndex map[int]int64, dbFilePath string) ([]models.Student, error) {
    ids, exists := gpaIndex[gpa]
    if !exists {
        return nil, fmt.Errorf("no records found with GPA: %f", gpa)
    }

    var results []models.Student
    file, err := os.Open(dbFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    for id := range ids {
        offset, exists := idIndex[id]
        if !exists {
            continue
        }

        _, err = file.Seek(offset, 0)
        if err != nil {
            continue
        }

        scanner := bufio.NewScanner(file)
        scanner.Scan()
        line := scanner.Text()

        var student models.Student
        err := json.Unmarshal([]byte(line), &student)
        if err != nil {
            continue
        }

        results = append(results, student)
    }

    return results, nil
}

func (r *Recorder) FindByActive(active bool, activeIndex map[bool]map[int]bool, idIndex map[int]int64, dbFilePath string) ([]models.Student, error) {
    ids, exists := activeIndex[active]
    if !exists {
        return nil, fmt.Errorf("no records found with active: %t", active)
    }

    var results []models.Student
    file, err := os.Open(dbFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    for id := range ids {
        offset, exists := idIndex[id]
        if !exists {
            continue
        }

        _, err = file.Seek(offset, 0)
        if err != nil {
            continue
        }

        scanner := bufio.NewScanner(file)
        scanner.Scan()
        line := scanner.Text()

        var student models.Student
        err := json.Unmarshal([]byte(line), &student)
        if err != nil {
            continue
        }

        results = append(results, student)
    }

    return results, nil
}

func (r *Recorder) ImportFromXLSX(xlsxPath string, dbFilePath string, idIndex map[int]int64, nameIndex map[string]map[int]bool, gpaIndex map[float64]map[int]bool, activeIndex map[bool]map[int]bool, recordInfo map[int]models.RecordInfo) error {
    f, err := excelize.OpenFile(xlsxPath)
    if err != nil {
        return fmt.Errorf("error opening XLSX file: %w", err)
    }
    defer f.Close()

    rows, err := f.GetRows("Sheet1")
    if err != nil {
        return fmt.Errorf("error reading sheet: %w", err)
    }
    if len(rows) == 0 {
        return fmt.Errorf("XLSX file is empty")
    }
    rows = rows[1:]

    for _, row := range rows {
        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            continue
        }

        name := row[1]
        gpa, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            continue
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
			continue
        }

        student := models.Student{
            Id:     id,
            Name:   name,
            Gpa:    gpa,
            Active: active,
        }
		record := models.Record{
            Id:      student.Id,
            Student: student,
        }

        err = r.AddNewRecord(record, dbFilePath, idIndex, nameIndex, gpaIndex, activeIndex, recordInfo)
        if err != nil {
            return fmt.Errorf("error adding record from XLSX: %w", err)
        }
    }

    fmt.Println("Import from XLSX completed")
    return nil
}