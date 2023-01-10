package mqtt

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const TestFile = `
path:
  "/group1/light1/":
    name: Btn1
    topic: /test/1
    typ: single
  "/*/*/":
    name: Btn2
    topic: /{GroupID}/{LightID}
    typ: toggle
    toggleValues:
    - name: On
      topic: /{GroupID}/{LightID}
      message: 
       on: true
    - name: Off
      topic: /{GroupID}/{LightID}
      message: 
       on: false
  "/group1/light3/":
    name: Btn3
    typ: multi
    multiValues: 
      - name: Val1
        topic: /group1/light3
        typ: single
      - name: Val2
        topic: /group1/light3
        typ: single
`

func TestLoadConfig(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *AdditionalButtonConfig
		wantErr bool
	}{
		{
			name: "Load File",
			args: args{strings.NewReader(TestFile)},
			want: &AdditionalButtonConfig{
				Path: map[string]ButtonConfig{
					"/group1/light1/": {
						Name:  "Btn1",
						Topic: "/test/1",
						Typ:   TypeSingle,
					},
					"/*/*/": {
						Name:  "Btn2",
						Topic: "/{GroupID}/{LightID}",
						Typ:   TypeToggle,
						ToggleValues: []ButtonConfig{
							{
								Name:  "On",
								Topic: "/{GroupID}/{LightID}",
								Message: map[string]interface{}{
									"on": true,
								},
							},
							{
								Name:  "Off",
								Topic: "/{GroupID}/{LightID}",
								Message: map[string]interface{}{
									"on": false,
								},
							},
						},
					},
					"/group1/light3/": {
						Name: "Btn3",
						Typ:  TypeMulti,
						MultiValues: []ButtonConfig{
							{
								Name:  "Val1",
								Topic: "/group1/light3",
								Typ:   TypeSingle,
							},
							{
								Name:  "Val2",
								Topic: "/group1/light3",
								Typ:   TypeSingle,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() got = \n%+v, want \n%+v", got, tt.want)
			}
		})
	}
}
