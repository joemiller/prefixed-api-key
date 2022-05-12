package apikey

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerator_defaults(t *testing.T) {
	gen, _ := NewGenerator("mycompany")
	assert.Equal(t, "mycompany", gen.keyPrefix)
	assert.Equal(t, "", gen.shortTokenPrefix)
	assert.Equal(t, defaultShortTokenLength, gen.shortTokenLength)
	assert.Equal(t, defaultLongTokenLength, gen.longTokenLength)
}

func TestGenerateAPIKey_defaults(t *testing.T) {
	gen, _ := NewGenerator("mycompany")
	k, _ := gen.GenerateAPIKey()
	spew.Dump(k)
	assert.Contains(t, k.Token(), "mycompany_")
	assert.Equal(t, defaultShortTokenLength, len(k.ShortToken()))
	assert.Equal(t, defaultLongTokenLength, len(k.LongToken()))
}

func TestNewGenerator_withOptions(t *testing.T) {
	gen, _ := NewGenerator("mycompany",
		WithShortTokenPrefix("foo"),
		WithShortTokenLength(9),
		WithLongTokenLength(25),
	)
	k, _ := gen.GenerateAPIKey()
	assert.Contains(t, k.Token(), "mycompany_")
	assert.Contains(t, k.Token(), "_foo")
	assert.Equal(t, 9, len(k.ShortToken()))
	assert.Equal(t, 25, len(k.LongToken()))
}

func TestParseAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectedKey Key
		shouldErr   bool
	}{
		{
			name:      "invalid key",
			token:     "foo",
			shouldErr: true,
		},
		{
			name:  "valid key",
			token: "mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp",
			expectedKey: Key{
				token:         "mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp",
				shortToken:    "5TJMbnP3thd",
				longToken:     "DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp",
				longTokenHash: "dec53496c2cd35bbb44270f9f896299a920b08b3c49cbf6ccd0c2cbb8efdcb83",
			},
			shouldErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			key, err := ParseAPIKey(tc.token)

			if tc.shouldErr {
				assert.Error(t, err)
			}
			if !tc.shouldErr {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedKey, key)
		})
	}
}

func TestCheckAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		hash      string
		shouldErr bool
		expected  bool
	}{
		{
			name:      "correct hash",
			token:     "mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp",
			hash:      "dec53496c2cd35bbb44270f9f896299a920b08b3c49cbf6ccd0c2cbb8efdcb83",
			shouldErr: false,
			expected:  true,
		},
		{
			name:      "incorrect hash",
			token:     "mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp",
			hash:      "incorrect",
			shouldErr: false,
			expected:  false,
		},
		{
			name:      "invalid token",
			token:     "1234",
			hash:      "",
			shouldErr: true,
			expected:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := CheckAPIKey(tc.token, tc.hash)

			if tc.shouldErr {
				assert.Error(t, err)
			}
			if !tc.shouldErr {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, ok)
		})
	}
}

func TestPadStart(t *testing.T) {

	tests := []struct {
		name     string
		s        string
		len      int
		pad      string
		expected string
	}{
		{
			name:     "empty string",
			s:        "",
			len:      8,
			pad:      "0",
			expected: "00000000",
		},
		{
			name:     "padded string",
			s:        "1234567",
			len:      8,
			pad:      "0",
			expected: "01234567",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := padStart(tc.s, tc.len, tc.pad)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
