package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kanimaru/godeconz/http"
	"reflect"
	"testing"
)

func TestLightSettingCallback_getButtonsForLights(t *testing.T) {
	type args struct {
		light http.LightResponseState
	}
	tests := []struct {
		name    string
		args    args
		want    []tgbotapi.InlineKeyboardButton
		wantErr bool
	}{
		{
			name: "No Light Reachable",
			args: args{
				light: http.LightResponseState{
					State: http.LightResponseStateDetail{
						Reachable: toPtr(false),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "On lights should have off button",
			args: args{
				light: http.LightResponseState{
					State: http.LightResponseStateDetail{
						Reachable: toPtr(true),
						On:        toPtr(true),
					},
				},
			},
			want: []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Off", "off"),
			},
			wantErr: false,
		},
		{
			name: "Off lights should have on button",
			args: args{
				light: http.LightResponseState{
					State: http.LightResponseStateDetail{
						Reachable: toPtr(true),
						On:        toPtr(false),
					},
				},
			},
			want: []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("On", "on"),
			},
			wantErr: false,
		},
		{
			name: "Temp lights should have temp button",
			args: args{
				light: http.LightResponseState{
					Ctmin: toPtr(10),
					Ctmax: toPtr(100),
					State: http.LightResponseStateDetail{
						Reachable: toPtr(true),
					},
				},
			},
			want: []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("On", "on"),
				tgbotapi.NewInlineKeyboardButtonData("Temperature", "temp"),
			},
			wantErr: false,
		},
		{
			name: "Color lights should have color button",
			args: args{
				light: http.LightResponseState{
					Hascolor: toPtr(true),
					State: http.LightResponseStateDetail{
						Reachable: toPtr(true),
					},
				},
			},
			want: []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("On", "on"),
				tgbotapi.NewInlineKeyboardButtonData("Color", "color"),
			},
			wantErr: false,
		},
		{
			name: "Lights with Brightness should have Brightness button",
			args: args{
				light: http.LightResponseState{
					State: http.LightResponseStateDetail{
						Reachable: toPtr(true),
						Bri:       toPtr(1),
					},
				},
			},
			want: []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("On", "on"),
				tgbotapi.NewInlineKeyboardButtonData("Brightness", "bright"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LightSettingCallback{}

			got, err := l.getButtonsForLights([]http.LightResponseState{tt.args.light})

			if (err != nil) != tt.wantErr {
				t.Errorf("getButtonsForLights() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getButtonsForLights() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
