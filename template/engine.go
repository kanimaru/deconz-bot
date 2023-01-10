package template

import (
	"encoding/xml"
	"github.com/PerformLine/go-stockutil/log"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type Engine struct {
	templates *template.Template

	// Unique button mapping for back references. Every button in every view needs a unique ID for the button
	refButtonMap map[string]*Button
	// presets for buttons
	buttons map[string]Button
	// preset for rows
	rows map[string]Row

	processedViews []*View
}

func CreateEngineByDir(dir string) *Engine {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get working directory.")
	}
	workingDir := filepath.Join(wd, dir)
	files, err := os.Open(workingDir)
	if err != nil {
		log.Fatalf("Can't open %v", workingDir)
	}
	defer files.Close()

	fileInfos, err := files.ReadDir(-1)
	if err != nil {
		log.Fatalf("Can't read directory %v", workingDir)
	}
	fileNames := make([]string, len(fileInfos))
	for i, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		fileNames[i] = dir + fileInfo.Name()
	}
	return CreateEngineByFiles(fileNames...)
}

func CreateEngineByFiles(files ...string) *Engine {
	engine := Engine{
		templates: template.New("root"),
		buttons:   map[string]Button{},
		rows:      map[string]Row{},
	}

	tmpl, err := engine.templates.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	engine.templates = tmpl
	engine.findPresets()
	return &engine
}

// findPresets looking through all templates for preset information and saves them for later use.
func (e *Engine) findPresets() {
	for _, t := range e.templates.Templates() {
		structure, err := e.Apply(t.Name(), nil)
		if err != nil {
			continue
		}
		for _, row := range structure.Row {
			e.rows[row.Name] = row
		}
		for _, button := range structure.Button {
			e.buttons[button.Name] = button
		}
	}
}

// Apply Fills the given template (by name) and returns the generic structure of the view
func (e *Engine) Apply(name string, data interface{}) (*Structure, error) {
	structure, err := e.parseStructure(name, data)
	if err != nil {
		return nil, err
	}
	e.inflate(structure.View)
	e.resolveButtonRefs(structure.View)
	return structure, nil
}

// parseStructure uses the given template content and parses the structure from the file
func (e *Engine) parseStructure(name string, data interface{}) (*Structure, error) {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		err := e.templates.ExecuteTemplate(writer, name, data)
		if err != nil {
			panic(err)
		}
	}()

	decoder := xml.NewDecoder(reader)
	structure := Structure{}
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &structure, nil
}

// resolveButtonRefs checks all buttons and generates unique IDs for the buttons
func (e *Engine) resolveButtonRefs(view *View) {
	e.processedViews = append(e.processedViews, view)
	for _, row := range view.Row {
		for _, button := range row.Button {
			id := button.getId()
			e.refButtonMap[id] = &button
			if button.View != nil && e.notProcessed(button.View) {
				e.resolveButtonRefs(view)
			}
		}
	}
}

// notProcessed checks if the view was already processed
func (e *Engine) notProcessed(view *View) bool {
	for _, processedView := range e.processedViews {
		if view == processedView {
			return true
		}
	}
	return false
}

func (e *Engine) findById(id string) *Button {
	return e.refButtonMap[id]
}

// inflate is looking for the use of preset information within the structure and replaces it with the real version also
// adds the references to the parents
func (e *Engine) inflate(view *View) {
	for i, row := range view.Row {
		if row.Use != "" {
			rowPreset, ok := e.rows[row.Use]
			if !ok {
				continue
			}
			rowPreset.parent = &view.Element
			rowPreset.copy(row)
			view.Row[i] = rowPreset
		}

		for i, button := range row.Button {
			if button.Use != "" {
				buttonPreset, ok := e.buttons[button.Name]
				if !ok {
					continue
				}
				buttonPreset.parent = &view.Element
				buttonPreset.copy(button)
				row.Button[i] = buttonPreset
			}
		}
	}
}
