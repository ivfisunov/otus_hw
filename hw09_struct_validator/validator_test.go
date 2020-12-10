package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$|len:8"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5|regexp:^\\d.\\d.\\d$"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	testsPass := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "111111111111111111111111111111111111",
				Name:   "Ivan",
				Age:    40,
				Email:  "xx@ya.ru",
				Role:   "admin",
				Phones: []string{"32423425678", "01234567890"},
			},
			nil,
		},
		{
			App{Version: "1.0.0"},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "hello",
			},
			nil,
		},
		{
			Token{
				Header:    []byte{5, 10},
				Payload:   []byte{5, 6},
				Signature: []byte{6, 6},
			},
			nil,
		},
	}

	for i, tt := range testsPass {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}

	testsFail := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "1",
				Name:   "Ivan",
				Age:    90,
				Email:  "xx@ya.ru",
				Role:   "provider",
				Phones: []string{"32423425678", "01234567890"},
			},
			nil,
		},
		{
			App{Version: "r.r.r"},
			nil,
		},
		{
			Response{
				Code: 201,
				Body: "hello",
			},
			nil,
		},
	}

	for i, tt := range testsFail {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Error(t, err)
		})
	}
}
