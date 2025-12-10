package gui

import (
    "fmt"
    "strconv"
    "time"
    "os"

    "fyne.io/fyne/v2"    
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/dialog"

    "github.com/kgugunava/database/db"
    "github.com/kgugunava/database/models"
)

type GUI struct {
    App    fyne.App
    Window fyne.Window
    DB     *db.Db

    list       *widget.List
    idEntry    *widget.Entry
    nameEntry  *widget.Entry
    gpaEntry   *widget.Entry
    activeEntry *widget.Entry
}

func NewGUI(database *db.Db) *GUI {
    a := app.New()
    w := a.NewWindow("Student Database")
    w.Resize(fyne.NewSize(1200, 800))

    gui := &GUI{
        App:    a,
        Window: w,
        DB:     database,
    }

    gui.setupUI()
    return gui
}

func (g *GUI) setupUI() {
    // поля для ввода
    g.idEntry = widget.NewEntry()
    g.idEntry.SetPlaceHolder("ID")

    g.nameEntry = widget.NewEntry()
    g.nameEntry.SetPlaceHolder("Name")

    g.gpaEntry = widget.NewEntry()
    g.gpaEntry.SetPlaceHolder("GPA")

    g.activeEntry = widget.NewEntry()
    g.activeEntry.SetPlaceHolder("Active (true/false)")

    // ТАБЛИЦА 
    g.list = widget.NewList(
        func() int { return len(g.DB.IdIndex) },
        func() fyne.CanvasObject {
            return container.NewHBox(
                widget.NewLabel("ID: "),
                widget.NewLabel("Name: "),
                widget.NewLabel("GPA: "),
                widget.NewLabel("Active: "),
            )
        },
        func(id widget.ListItemID, obj fyne.CanvasObject) {
        },
    )

    //  ДОБАВЛЕНИЕ 
    addForm := widget.NewForm(
        &widget.FormItem{Text: "ID", Widget: g.idEntry},
        &widget.FormItem{Text: "Name", Widget: g.nameEntry},
        &widget.FormItem{Text: "GPA", Widget: g.gpaEntry},
        &widget.FormItem{Text: "Active", Widget: g.activeEntry},
    )

    addBtn := widget.NewButton("Add Student", g.addStudent)

    // УДАЛЕНИЕ 
    deleteIdEntry := widget.NewEntry()
    deleteIdEntry.SetPlaceHolder("ID to delete")
    deleteNameEntry := widget.NewEntry()
    deleteNameEntry.SetPlaceHolder("Name to delete")
    deleteGpaEntry := widget.NewEntry()
    deleteGpaEntry.SetPlaceHolder("GPA to delete")
    deleteActiveEntry := widget.NewEntry()
    deleteActiveEntry.SetPlaceHolder("Active to delete (true/false)")

    deleteByIdBtn := widget.NewButton("Delete by ID", func() {
        id, err := strconv.Atoi(deleteIdEntry.Text)
        if err != nil {
            g.showNotification("Invalid ID")
            return
        }
        g.deleteStudentByIdWithId(id)
    })

    deleteByNameBtn := widget.NewButton("Delete by Name", func() {
        g.deleteStudentByNameWithName(deleteNameEntry.Text)
    })

    deleteByGpaBtn := widget.NewButton("Delete by GPA", func() {
        gpa, err := strconv.ParseFloat(deleteGpaEntry.Text, 64)
        if err != nil {
            g.showNotification("Invalid GPA")
            return
        }
        g.deleteStudentByGpaWithGpa(gpa)
    })

    deleteByActiveBtn := widget.NewButton("Delete by Active", func() {
        active, err := strconv.ParseBool(deleteActiveEntry.Text)
        if err != nil {
            g.showNotification("Invalid Active value")
            return
        }
        g.deleteStudentByActiveWithActive(active)
    })

    deleteForm := widget.NewForm(
        &widget.FormItem{Text: "ID", Widget: deleteIdEntry},
        &widget.FormItem{Text: "Name", Widget: deleteNameEntry},
        &widget.FormItem{Text: "GPA", Widget: deleteGpaEntry},
        &widget.FormItem{Text: "Active", Widget: deleteActiveEntry},
    )

    // ПОИСК 
    searchIdEntry := widget.NewEntry()
    searchIdEntry.SetPlaceHolder("ID to search")
    searchNameEntry := widget.NewEntry()
    searchNameEntry.SetPlaceHolder("Name to search")
    searchGpaEntry := widget.NewEntry()
    searchGpaEntry.SetPlaceHolder("GPA to search")
    searchActiveEntry := widget.NewEntry()
    searchActiveEntry.SetPlaceHolder("Active to search (true/false)")

    searchByIdBtn := widget.NewButton("Search by ID", func() {
        id, err := strconv.Atoi(searchIdEntry.Text)
        if err != nil {
            g.showNotification("Invalid ID")
            return
        }
        g.searchStudentByIdWithId(id)
    })

    searchByNameBtn := widget.NewButton("Search by Name", func() {
        g.searchStudentByNameWithName(searchNameEntry.Text)
    })

    searchByGpaBtn := widget.NewButton("Search by GPA", func() {
        gpa, err := strconv.ParseFloat(searchGpaEntry.Text, 64)
        if err != nil {
            g.showNotification("Invalid GPA")
            return
        }
        g.searchStudentByGpaWithGpa(gpa)
    })

    searchByActiveBtn := widget.NewButton("Search by Active", func() {
        active, err := strconv.ParseBool(searchActiveEntry.Text)
        if err != nil {
            g.showNotification("Invalid Active value")
            return
        }
        g.searchStudentByActiveWithActive(active)
    })

    searchForm := widget.NewForm(
        &widget.FormItem{Text: "ID", Widget: searchIdEntry},
        &widget.FormItem{Text: "Name", Widget: searchNameEntry},
        &widget.FormItem{Text: "GPA", Widget: searchGpaEntry},
        &widget.FormItem{Text: "Active", Widget: searchActiveEntry},
    )

    // ОСНОВНОЕ МЕНЮ
    backupBtn := widget.NewButton("Create Backup", g.createBackup)
    restoreBtn := widget.NewButton("Restore from Backup", g.restoreFromBackup)
    importBtn := widget.NewButton("Import from XLSX", g.importFromXLSX)

    // КОНТЕНТ 
    content := container.NewVBox(
        widget.NewCard("Add Student", "", container.NewVBox(addForm, addBtn)),
        widget.NewCard("Delete Student", "", container.NewVBox(deleteForm, 
            container.NewHBox(deleteByIdBtn, deleteByNameBtn, deleteByGpaBtn, deleteByActiveBtn))),
        widget.NewCard("Search Student", "", container.NewVBox(searchForm,
            container.NewHBox(searchByIdBtn, searchByNameBtn, searchByGpaBtn, searchByActiveBtn))),
        container.NewHBox(backupBtn, restoreBtn, importBtn),
        widget.NewLabel("Records:"),
        g.list,
    )

    g.Window.SetContent(content)
}

func (g *GUI) addStudent() {
    id, err := strconv.Atoi(g.idEntry.Text)
    if err != nil {
        g.showNotification("Invalid ID")
        return
    }

    gpa, err := strconv.ParseFloat(g.gpaEntry.Text, 64)
    if err != nil {
        g.showNotification("Invalid GPA")
        return
    }

    active, err := strconv.ParseBool(g.activeEntry.Text)
    if err != nil {
        g.showNotification("Invalid Active value")
        return
    }

    student := models.Student{
        Id:     id,
        Name:   g.nameEntry.Text,
        Gpa:    gpa,
        Active: active,
    }

    record := models.Record{
        Id:      student.Id,
        Student: student,
    }

    err = g.DB.Recorder.AddNewRecord(
        record,
        g.DB.FilePath,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error adding record: " + err.Error())
        return
    }

    g.showNotification("Student added successfully")
    g.list.Refresh()
    g.clearInputs()
}

// УДАЛЕНИЕ 

func (g *GUI) deleteStudentById() {
    id, err := strconv.Atoi(g.idEntry.Text)
    if err != nil {
        g.showNotification("Invalid ID")
        return
    }

    err = g.DB.Recorder.DeleteRecordById(
        id,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting record: " + err.Error())
        return
    }

    g.showNotification("Student deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByName() {
    err := g.DB.Recorder.DeleteRecordByName(
        g.nameEntry.Text,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByGpa() {
    gpa, err := strconv.ParseFloat(g.gpaEntry.Text, 64)
    if err != nil {
        g.showNotification("Invalid GPA")
        return
    }

    err = g.DB.Recorder.DeleteRecordByGpa(
        gpa,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByActive() {
    active, err := strconv.ParseBool(g.activeEntry.Text)
    if err != nil {
        g.showNotification("Invalid Active value")
        return
    }

    err = g.DB.Recorder.DeleteRecordByActive(
        active,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

// ПОИСК 

func (g *GUI) searchStudentById() {
    id, err := strconv.Atoi(g.idEntry.Text)
    if err != nil {
        g.showNotification("Invalid ID")
        return
    }

    student, err := g.DB.Recorder.FindById(id, g.DB.IdIndex, g.DB.FilePath)
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    g.showNotification(fmt.Sprintf("Found student: %s (ID: %d)", student.Name, student.Id))
}

func (g *GUI) searchStudentByName() {
    results, err := g.DB.Recorder.FindByName(
        g.nameEntry.Text,
        g.DB.NameIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    g.showNotification(fmt.Sprintf("Found %d records by name", len(results)))
}

func (g *GUI) searchStudentByGpa() {
    gpa, err := strconv.ParseFloat(g.gpaEntry.Text, 64)
    if err != nil {
        g.showNotification("Invalid GPA")
        return
    }

    results, err := g.DB.Recorder.FindByGpa(
        gpa,
        g.DB.GpaIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    g.showNotification(fmt.Sprintf("Found %d records by GPA", len(results)))
}

func (g *GUI) searchStudentByActive() {
    active, err := strconv.ParseBool(g.activeEntry.Text)
    if err != nil {
        g.showNotification("Invalid Active value")
        return
    }

    results, err := g.DB.Recorder.FindByActive(
        active,
        g.DB.ActiveIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    g.showNotification(fmt.Sprintf("Found %d records by Active", len(results)))
}

// BACKUP и XLSX 

func (g *GUI) createBackup() {
    err := os.MkdirAll("./backups", 0755)
    if err != nil {
        g.showNotification("Error creating backup directory: " + err.Error())
        return
    }

    err = g.DB.CreateBackup("./backups")
    if err != nil {
        g.showNotification("Error creating backup: " + err.Error())
        return
    }

    g.showNotification("Backup created successfully")
}

func (g *GUI) restoreFromBackup() {
    dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil || reader == nil {
            return
        }
        defer reader.Close()

        backupPath := reader.URI().Path()

        err = g.DB.RestoreFromBackup(backupPath)
        if err != nil {
            g.showNotification("Error restoring from backup: " + err.Error())
            return
        }

        g.showNotification("Database restored successfully")
        g.list.Refresh()
    }, g.Window)
}

func (g *GUI) importFromXLSX() {
    dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil || reader == nil {
            return
        }
        defer reader.Close()

        xlsxPath := reader.URI().Path()

        err = g.DB.Recorder.ImportFromXLSX(
            xlsxPath,
            g.DB.FilePath,
            g.DB.IdIndex,
            g.DB.NameIndex,
            g.DB.GpaIndex,
            g.DB.ActiveIndex,
            g.DB.RecordInfo,
        )
        if err != nil {
            g.showNotification("Error importing from XLSX: " + err.Error())
            return
        }

        g.showNotification("Import completed successfully")
        g.list.Refresh()
    }, g.Window)
}

func (g *GUI) showNotification(message string) {
    dialog := widget.NewModalPopUp(widget.NewLabel(message), g.Window.Canvas())
    dialog.Show()
    
    go func() {
        time.Sleep(5 * time.Second)
        dialog.Hide()
    }()
}

func (g *GUI) clearInputs() {
    g.idEntry.SetText("")
    g.nameEntry.SetText("")
    g.gpaEntry.SetText("")
    g.activeEntry.SetText("")
}

func (g *GUI) Run() {
    g.Window.ShowAndRun()
}

// УДАЛЕНИЕ С ОТДЕЛЬНЫМИ ПАРАМЕТРАМИ 
func (g *GUI) deleteStudentByIdWithId(id int) {
    err := g.DB.Recorder.DeleteRecordById(
        id,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting record: " + err.Error())
        return
    }

    g.showNotification("Student deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByNameWithName(name string) {
    err := g.DB.Recorder.DeleteRecordByName(
        name,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByGpaWithGpa(gpa float64) {
    err := g.DB.Recorder.DeleteRecordByGpa(
        gpa,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

func (g *GUI) deleteStudentByActiveWithActive(active bool) {
    err := g.DB.Recorder.DeleteRecordByActive(
        active,
        g.DB.IdIndex,
        g.DB.NameIndex,
        g.DB.GpaIndex,
        g.DB.ActiveIndex,
        g.DB.RecordInfo,
    )
    if err != nil {
        g.showNotification("Error deleting records: " + err.Error())
        return
    }

    g.showNotification("Students deleted successfully")
    g.list.Refresh()
}

//  ПОИСК С ОТДЕЛЬНЫМИ ПАРАМЕТРАМИ 
func (g *GUI) searchStudentByIdWithId(id int) {
    student, err := g.DB.Recorder.FindById(id, g.DB.IdIndex, g.DB.FilePath)
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    message := fmt.Sprintf("Found student:\nID: %d\nName: %s\nGPA: %.2f\nActive: %t", 
        student.Id, student.Name, student.Gpa, student.Active)
    g.showNotification(message)
}

func (g *GUI) searchStudentByNameWithName(name string) {
    results, err := g.DB.Recorder.FindByName(
        name,
        g.DB.NameIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    if len(results) == 0 {
        g.showNotification("No students found with name: " + name)
        return
    }

    message := fmt.Sprintf("Found %d students with name '%s':\n", len(results), name)
    for _, student := range results {
        message += fmt.Sprintf("\nID: %d, Name: %s, GPA: %.2f, Active: %t", 
            student.Id, student.Name, student.Gpa, student.Active)
    }

    g.showNotification(message)
}

func (g *GUI) searchStudentByGpaWithGpa(gpa float64) {
    results, err := g.DB.Recorder.FindByGpa(
        gpa,
        g.DB.GpaIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    if len(results) == 0 {
        g.showNotification("No students found with GPA: " + fmt.Sprintf("%.2f", gpa))
        return
    }

    message := fmt.Sprintf("Found %d students with GPA %.2f:\n", len(results), gpa)
    for _, student := range results {
        message += fmt.Sprintf("\nID: %d, Name: %s, GPA: %.2f, Active: %t", 
            student.Id, student.Name, student.Gpa, student.Active)
    }

    g.showNotification(message)
}

func (g *GUI) searchStudentByActiveWithActive(active bool) {
    results, err := g.DB.Recorder.FindByActive(
        active,
        g.DB.ActiveIndex,
        g.DB.IdIndex,
        g.DB.FilePath,
    )
    if err != nil {
        g.showNotification("Error searching: " + err.Error())
        return
    }

    if len(results) == 0 {
        g.showNotification("No students found with Active: " + fmt.Sprintf("%t", active))
        return
    }

    message := fmt.Sprintf("Found %d students with Active %t:\n", len(results), active)
    for _, student := range results {
        message += fmt.Sprintf("\nID: %d, Name: %s, GPA: %.2f, Active: %t", 
            student.Id, student.Name, student.Gpa, student.Active)
    }

    g.showNotification(message)
}