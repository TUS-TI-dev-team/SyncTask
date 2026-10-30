package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogoutInput(t *testing.T) {
	tests := []struct {
		name 		string
		input		LogoutInput
		wantbool	bool
	}{
		{
			name: "正常系: uuid の形式,バージョンが正しい1-> true",
			input: LogoutInput{SessionID: "1b50f79f-816e-4ec8-b9ec-24e573a0c8ff"},
			wantbool: true,
		},
		{
			name: "正常系: uuid の形式,バージョンが正しい2 -> true",
			input: LogoutInput{SessionID: "e00d6f4a-31aa-41a9-b9e9-967f76903dfd"},
			wantbool: true,
		},
		{
			name: "異常系: uuid の形式が正しくバージョンが不正1 -> false",
			input: LogoutInput{SessionID: "01a18b8e-512c-788f-b0c0-0c36922e7b5"},
			wantbool: false,
		},
		{
			name: "異常系: uuid の形式が正しくバージョンが不正2 -> false",
			input: LogoutInput{SessionID: "01a18b8e-512c-7493-b28b-3a77d1fd444c"},
			wantbool: false,
		},
		{
			name: "異常系: uuid の形式,バージョンともに不正1 -> false",
			input: LogoutInput{SessionID: "eb4cb619-94cfa463d-ad05a6a5d04c1f46f"},
			wantbool: false,
		},
		{
			name: "異常系: uuid の形式,バージョンともに不正2 -> false",
			input: LogoutInput{SessionID: "eb4cb619-f,/a.wr,.,21,31."},
			wantbool: false,
		},
		{
			name: "異常系: uuid の形式,バージョンともに不正3 -> false",
			input: LogoutInput{SessionID: "eb4cb619-94cfa463d-ad05a6awgezsgoibj;inlr.zjgs.awvjcｘｑamvlwCQX+WPiejowiht5d04c1f46f"},
			wantbool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := tt.input.Validate()
			assert.Equal(t, tt.wantbool, tf)
		})
	}
}
