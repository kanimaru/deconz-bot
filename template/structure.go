package template

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/PerformLine/go-stockutil/log"
	"github.com/PerformLine/go-stockutil/sliceutil"
	"strings"
)

type Element struct {
	Name   string `xml:"name,attr"`
	Data   string `xml:"data,attr,omitempty"`
	Use    string `xml:"use,attr,omitempty"`
	Parent *Element
}

type Preset struct {
	XMLName xml.Name `xml:"preset"`
	Button  []Button `xml:"button"`
	Row     []Row    `xml:"row"`
}

type Button struct {
	Element
	Label   string  `xml:"label,attr"`
	OnClick *string `xml:"onClick,attr,omitempty"`
	Text    string  `xml:",chardata"`
	View    *View   `xml:"view,omitempty"`
}

type Row struct {
	Element
	Text   string    `xml:",chardata"`
	Button []*Button `xml:"button"`
}

type View struct {
	Element
	XMLName      xml.Name `xml:"view"`
	Text         string   `xml:",chardata"`
	Row          []*Row   `xml:"row"`
	refButtonMap map[string]*Button
	Data         interface{}
}

func (b *Button) GetId() string {
	sb := strings.Builder{}
	var cur = &b.Element

	for {
		sb.WriteString(cur.Name)
		cur = cur.Parent
		if cur == nil {
			break
		} else {
			sb.WriteRune('<')
		}
	}
	log.Infof("%v: %v", b.Label, sb.String())
	sum := md5.Sum([]byte(sb.String()))
	return fmt.Sprintf("%x", sum)
}

func (r *Row) copy(row *Row) {
	if r.Element.Name == "" {
		r.Element.Name = row.Element.Name
	}
	for _, button := range row.Button {
		if sliceutil.Contains(r.Button, button) {
			continue
		}
		copyButton := *button
		r.Button = append(r.Button, &copyButton)
	}
}

func (b *Button) copy(button *Button) {
	if b.View == nil {
		b.View = button.View
	}
	if b.Label == "" {
		b.Label = button.Label
	}
	if b.Data == "" {
		b.Data = button.Data
	}
	if b.OnClick == nil {
		b.OnClick = button.OnClick
	}
	if b.Element.Name == "" {
		b.Element.Name = button.Element.Name
	}
}

func (v View) FindButton(id string) *Button {
	return v.refButtonMap[id]
}
