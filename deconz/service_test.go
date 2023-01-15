package deconz

import (
	"encoding/json"
	"github.com/kanimaru/godeconz"
	"github.com/kanimaru/godeconz/http"
	"testing"
)

func TestService_SetLightState(t *testing.T) {
	t.Parallel()
	type args struct {
		state LightState
	}
	tests := []struct {
		name     string
		args     args
		expected http.LightRequestState
	}{
		{
			name: "Turn off on bri 0",
			args: args{
				state: LightState{
					Brightness: toPtr(uint8(0)),
				},
			},
			expected: http.LightRequestState{
				Bri: toPtr(uint8(0)),
				On:  toPtr(false),
			},
		},
		{
			name: "Turn On on bri < 0",
			args: args{
				state: LightState{
					Brightness: toPtr(uint8(50)),
				},
			},
			expected: http.LightRequestState{
				Bri: toPtr(uint8(50)),
				On:  toPtr(true),
			},
		},
		{
			name: "Turn Red with right HSV",
			args: args{
				state: LightState{
					Color: "#FF0000",
				},
			},
			expected: http.LightRequestState{
				Hue: toPtr(uint16(0)),
				Sat: toPtr(uint8(255)),
				Bri: toPtr(uint8(255)),
			},
		},
		{
			name: "Turn Blue with right HSV",
			args: args{
				state: LightState{
					Color: "#0000AA",
				},
			},
			expected: http.LightRequestState{
				Hue: toPtr(uint16(43690)),
				Sat: toPtr(uint8(255)),
				Bri: toPtr(uint8(170)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient, httpClient := CreateMockClient()

			d := service[any]{
				client: mockClient,
			}
			d.SetLightState(tt.args.state, "test_light")
			got, _ := json.Marshal(httpClient.lastData)
			expected, _ := json.Marshal(tt.expected)
			if string(got) != string(expected) {
				t.Errorf("Expected %s got %s", expected, got)
			}
		})
	}
}

func CreateMockClient() (*http.Client[any], *MockAdapter) {
	adapter := MockAdapter{}
	client := http.CreateClient[any](&adapter, godeconz.Settings{
		Address:      "",
		HttpProtocol: "",
		ApiKey:       "",
	})
	return &client, &adapter
}

type MockAdapter struct {
	lastData interface{}
}

func (m *MockAdapter) Get(path string, container interface{}) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockAdapter) Post(path string, data interface{}, container interface{}) (any, error) {
	m.lastData = data
	return nil, nil
}

func (m *MockAdapter) Put(path string, data interface{}, container interface{}) (any, error) {
	m.lastData = data
	return nil, nil
}

func (m *MockAdapter) Delete(path string, container interface{}) (any, error) {
	//TODO implement me
	panic("implement me")
}

func toPtr[T any](val T) *T {
	return &val
}
