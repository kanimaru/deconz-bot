package template

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/PerformLine/go-stockutil/log"
	"strings"
)

type Element struct {
	Name   string `xml:"name,attr"`
	parent *Element
}

type Structure struct {
	XMLName xml.Name `xml:"bot"`
	Text    string   `xml:",chardata"`
	Button  []Button `xml:"button"`
	Row     []Row    `xml:"row"`
	View    *View    `xml:"view,omitempty"`
}

type Button struct {
	Element
	Use     string  `xml:"use,attr"`
	Label   string  `xml:"label,attr"`
	OnClick *string `xml:"onClick,attr,omitempty"`
	Data    string  `xml:"data,attr"`
	Text    string  `xml:",chardata"`
	View    *View   `xml:"view,omitempty"`
}

type Row struct {
	Element
	Text   string   `xml:",chardata"`
	Button []Button `xml:"button"`
	Use    string   `xml:"use,attr"`
}

type View struct {
	Element
	Text string `xml:",chardata"`
	Name string `xml:"name,attr"`
	Row  []Row  `xml:"row"`
}

func (b *Button) getId() string {
	sb := strings.Builder{}
	var cur = &b.Element

	for cur.parent != nil {
		sb.WriteString(cur.Name)
		sb.WriteRune('<')
		cur = cur.parent
	}
	log.Infof(sb.String())
	sum := md5.Sum([]byte(sb.String()))
	return fmt.Sprintf("%x", sum)
}

func (r *Row) copy(row Row) {
	for _, button := range row.Button {
		r.Button = append(r.Button, button)
	}
}

func (b *Button) copy(button Button) {
	if button.View != nil {
		b.View = button.View
	}
	if button.Label != "" {
		b.Label = button.Label
	}
	if button.Data != "" {
		b.Data = button.Data
	}
	if button.OnClick != nil {
		b.OnClick = button.OnClick
	}
}
