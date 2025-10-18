package hw09structvalidator

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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"regexp:\\d+|len:5"`
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
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "test",
				Name:   "test",
				Age:    34,
				Email:  "test@test.com",
				Role:   "test",
				Phones: []string{"12345"},
				meta:   nil,
			},
			expectedErr: fmt.Errorf(
				"field: ID | err: ошибка превышения длины строки\n" +
					"field: Phones | err: ошибка превышения длины строки\n"),
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "test",
				Age:    34,
				Email:  "test@test.com",
				Role:   "test",
				Phones: []string{"12345678901", "88005555555"},
				meta:   nil,
			},
			expectedErr: fmt.Errorf(""),
		},
		{
			in: App{
				Version: "12323232323",
			},
			expectedErr: fmt.Errorf(
				"field: Version | err: ошибка превышения длины строки\n", //nolint:revive
			),
		},
		{
			in: App{
				Version: "sasha",
			},
			expectedErr: fmt.Errorf(
				"field: Version | err: строка не соответствует регулярному выражению\n", //nolint:revive
			),
		},
		{
			in: App{
				Version: "123123dasda",
			},
			expectedErr: fmt.Errorf(
				"field: Version | err: ошибка превышения длины строки\n", //nolint:revive
			),
		},
		{
			in: App{
				Version: "55555",
			},
			expectedErr: fmt.Errorf(""),
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: fmt.Errorf(""),
		},
		{
			in: Response{
				Body: "anytext",
				Code: 503,
			},
			expectedErr: fmt.Errorf(
				"field: Code | err: число не входит в указанное множество\n"),
		},
		{
			in: Response{
				Body: "anytext",
				Code: 500,
			},
			expectedErr: fmt.Errorf(""),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.EqualError(t, Validate(tt.in), tt.expectedErr.Error())
			// _ = tt
		})
	}
}
