package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	iconResource, _ := fyne.LoadResourceFromPath("")
	a.SetIcon(iconResource)
	w := a.NewWindow("")
	w.SetMaster()
	showHomeScreen(w)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func makeNav(w fyne.Window) fyne.CanvasObject {
	tree := widget.NewTreeWithStrings(menuItems)
	tree.OnSelected = func(id string) {
		if id == "Request Order" {
			showUploadScreen(w)
		}
		if id == "Home" {
			showHomeScreen(w)
		}
	}
	return container.NewBorder(nil, nil, nil, nil, tree)
}

func showHomeScreen(w fyne.Window) {
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("Home Screen")
	intro.Wrapping = fyne.TextWrapWord
	homeData := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)

	split := container.NewHSplit(makeNav(w), homeData)
	split.Offset = 0
	w.SetContent(split)
}

func showUploadScreen(w fyne.Window) {
	content := container.NewVBox(
		widget.NewLabel("Upload Excel/CSV"),
		widget.NewButton("Upload", func() {
			uploadFile(w)
		}),
	)

	split := container.NewHSplit(makeNav(w), content)
	split.Offset = 0
	w.SetContent(split)
}

func uploadFile(w fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			defer reader.Close()
			fileName := reader.URI().Name()
			ext := getFileExtension(fileName)
			if ext != ".csv" && ext != ".xls" && ext != ".xlsx" {
				dialog.ShowError(errors.New("unsupported file format"), w)
				return
			}

			// Read the content of the uploaded file
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Convert to JSON
			jsonData, err := convertToJSON(data, ext)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Make API request
			err = makeAPIRequest(jsonData)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Optionally, keep the file for history
			// err = saveFileForHistory(fileName, data)
			// if err != nil {
			// 	log.Println("Error saving file for history:", err)
			// }

			dialog.ShowInformation("Success", "File uploaded and processed successfully", w)
		}
	}, w)
}

func getFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func convertToJSON(data []byte, ext string) ([]byte, error) {
	if ext == ".csv" {
		return csvToJSON(data)
	} else if ext == ".xls" || ext == ".xlsx" {
		return excelToJSON(data)
	}
	return nil, errors.New("unsupported file format")
}

func csvToJSON(data []byte) ([]byte, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return json.Marshal(records)
}

func excelToJSON(data []byte) ([]byte, error) {
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}

	var rows [][]string
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			var cells []string
			for _, cell := range row.Cells {
				cells = append(cells, cell.String())
			}
			rows = append(rows, cells)
		}
	}

	return json.Marshal(rows)
}

func makeAPIRequest(data []byte) error {
	fmt.Println("Sending data to API:", string(data))
	// Implement API request logic here
	return nil
}

func saveFileForHistory(fileName string, data []byte) error {
	// Implement logic to save the uploaded file for history
	// This is just a placeholder implementation
	filePath := filepath.Join("history", fileName)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	fmt.Println("File saved for history:", filePath)
	return nil
}

var menuItems = map[string][]string{
	"": {"Home", "Request Order"},
}
