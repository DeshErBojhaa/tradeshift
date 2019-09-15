package api

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewResponse(t *testing.T) {
	type args struct {
		statusCode int
		mediaType  string
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "OKJSON",
			args: args{
				statusCode: 200,
				mediaType:  "application/json",
			},
			want: &Response{
				statusCode: 200,
				mediaType:  "application/json",
				data:       map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			got := NewResponse(tt.args.statusCode, tt.args.mediaType)
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestResponse_Data(t *testing.T) {
	type fields struct {
		statusCode int
		mediaType  string
		data       map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "HappyPath",
			fields: fields{
				statusCode: 200,
				mediaType:  "application/json",
				data:       map[string]interface{}{},
			},
			args: args{
				key:   "csv",
				value: "name,email,title\nLuke,luke@tatooine.planet,Jedi\n",
			},
			want: &Response{
				statusCode: 200,
				mediaType:  "application/json",
				data: map[string]interface{}{
					"csv": "name,email,title\nLuke,luke@tatooine.planet,Jedi\n",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			r := &Response{
				statusCode: tt.fields.statusCode,
				mediaType:  tt.fields.mediaType,
				data:       tt.fields.data,
			}

			got := r.Data(tt.args.key, tt.args.value)
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestResponse_marshalJSON(t *testing.T) {
	type thing struct {
		Name    string `json:"name"`
		Purpose string `json:"purpose"`
	}

	type fields struct {
		statusCode int
		mediaType  string
		data       map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "HappyPath",
			fields: fields{
				statusCode: 200,
				mediaType:  "application/json",
				data: map[string]interface{}{
					"thing": thing{
						Name:    "rock",
						Purpose: "Crush scissors.",
					},
				},
			},
			want: []byte(`{"thing":{"name":"rock","purpose":"Crush scissors."}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			r := &Response{
				statusCode: tt.fields.statusCode,
				mediaType:  tt.fields.mediaType,
				data:       tt.fields.data,
			}
			got := r.marshalJSON()
			g.Expect(got).To(Equal(tt.want))
		})
	}
}
