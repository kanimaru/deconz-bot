package template

import (
	"encoding/xml"
	"github.com/PerformLine/go-stockutil/log"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Engine interface {
	// Apply Fills the given template (by name) and returns the generic structure of the view
	Apply(name string, data interface{}) (*View, error)
}

type engine struct {
	templates *template.Template

	// presets for buttons
	buttons map[string]Button
	// preset for rows
	rows map[string]Row

	processedViews []*View
}

func CreateEngineByDir(dir string) Engine {
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

func CreateEngineByFiles(files ...string) Engine {
	e := &engine{
		templates:      template.New("root"),
		buttons:        make(map[string]Button),
		rows:           make(map[string]Row),
		processedViews: make([]*View, 100),
	}

	tmpl, err := e.templates.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	e.templates = tmpl
	e.findPresets(files)
	return e
}

// findPresets looking through all templates for preset information and saves them for later use.
func (e *engine) findPresets(paths []string) {
	for _, path := range paths {
		preset, err := e.parsePreset(path)
		if err != nil {
			continue
		}
		for _, row := range preset.Row {
			e.rows[row.Name] = row
		}
		for _, button := range preset.Button {
			e.buttons[button.Name] = button
		}
	}
}

// Apply Fills the given template (by name) and returns the generic structure of the view
func (e *engine) Apply(name string, data interface{}) (*View, error) {
	view, err := e.parseView(name, data)
	if err != nil {
		return nil, err
	}
	e.inflate(view)
	// clear processed views
	e.processedViews = e.processedViews[:0]
	return view, nil
}

// parsePreset uses the given template content and parses the structure from the file
func (e *engine) parsePreset(name string) (*Preset, error) {
	file, err := os.Open(name)
	if err != nil {
		log.Debugf("Can't open preset %w cause: %w", name, err)
		return nil, err
	}
	decoder := xml.NewDecoder(file)
	preset := Preset{}
	err = decoder.Decode(&preset)
	if err != nil {
		return nil, err
	}
	return &preset, nil
}

// parseView uses the given template content and parses the structure from the file
func (e *engine) parseView(name string, data interface{}) (*View, error) {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		err := e.templates.ExecuteTemplate(writer, name, data)
		if err != nil {
			panic(err)
		}
	}()

	decoder := xml.NewDecoder(reader)
	structure := View{}
	err := decoder.Decode(&structure)
	if err != nil {
		return nil, err
	}
	return &structure, nil
}

// processed checks if the view was already processed
func (e *engine) processed(view *View) bool {
	for _, processedView := range e.processedViews {
		if view == processedView {
			return true
		}
	}
	return false
}

// inflate is looking for the use of preset information within the structure and replaces it with the real version also
// adds the references to the parents
func (e *engine) inflate(view *View) {
	if view.Name == "" {
		view.Name = generateName()
	}
	view.refButtonMap = make(map[string]*Button)
	view.Text = strings.TrimSpace(view.Text)

	for _, row := range view.Row {
		if row.Use != "" {
			rowPreset, ok := e.rows[row.Use]
			if !ok {
				continue
			}
			row.copy(&rowPreset)
		}
		if row.Name == "" {
			row.Name = generateName()
		}
		row.Parent = &view.Element
		row.Text = strings.TrimSpace(row.Text)

		for _, button := range row.Button {
			if button.Use != "" {
				buttonPreset, ok := e.buttons[button.Use]
				if !ok {
					continue
				}
				button.copy(&buttonPreset)
			}
			if button.Name == "" {
				button.Name = generateName()
			}
			button.Parent = &row.Element
			view.refButtonMap[button.GetId()] = button
			button.Text = strings.TrimSpace(button.Text)

			if button.View != nil && !e.processed(button.View) {
				e.processedViews = append(e.processedViews, button.View)
				button.View.Element.Parent = &button.Element
				e.inflate(button.View)
			}
		}
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateName() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
