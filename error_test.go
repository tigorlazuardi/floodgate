package floodgate

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newGateError(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want GateError
	}{
		{
			name: "success",
			args: args{name: "some service", err: errors.New("test")},
			want: GateError{
				Name:      "some service",
				OriginErr: errors.New("test"),
				Err:       "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newGateError(tt.args.name, tt.args.err))
		})
	}
}

func TestGateError_MarshalJSON(t *testing.T) {
	type fields struct {
		Name      string
		OriginErr error
		Err       string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "empty json marshal error variable",
			fields: fields{
				Name:      "some service",
				OriginErr: errors.New("test"),
				Err:       "test",
			},
			want: []byte(`{"cause":"test","error":"test","name":"some service"}`),
		},
		{
			name: "json marshal not empty error variable",
			fields: fields{
				Name:      "some service",
				OriginErr: &strconv.NumError{Func: "test", Num: "0", Err: errors.New("test")},
				Err:       (&strconv.NumError{Func: "test", Num: "0", Err: errors.New("test")}).Error(),
			},
			want: []byte(`{"cause":"{\"Func\":\"test\",\"Num\":\"0\",\"Err\":{}}","error":"strconv.test: parsing \"0\": test","name":"some service"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ge := GateError{
				Name:      tt.fields.Name,
				OriginErr: tt.fields.OriginErr,
				Err:       tt.fields.Err,
			}
			got, err := ge.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, string(tt.want), string(got))
		})
	}
}
