package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/akamensky/base58"
)

const (
	defaultShortTokenLength = 8
	defaultLongTokenLength  = 24
)

type Option func(*Generator) error

type Key struct {
	token         string
	shortToken    string
	longToken     string
	longTokenHash string
}

type Generator struct {
	keyPrefix        string
	shortTokenPrefix string
	shortTokenLength int
	longTokenLength  int
}

func WithShortTokenPrefix(key string) Option {
	return func(g *Generator) error {
		g.shortTokenPrefix = key
		return nil
	}
}

func WithShortTokenLength(n int) Option {
	return func(g *Generator) error {
		g.shortTokenLength = n
		return nil
	}
}

func WithLongTokenLength(n int) Option {
	return func(g *Generator) error {
		g.longTokenLength = n
		return nil
	}
}

func NewGenerator(keyPrefix string, opts ...Option) (*Generator, error) {
	// defaults
	g := &Generator{
		keyPrefix:        keyPrefix,
		shortTokenPrefix: "",
		shortTokenLength: defaultShortTokenLength,
		longTokenLength:  defaultLongTokenLength,
	}

	for _, f := range opts {
		err := f(g)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g *Generator) GenerateAPIKey() (Key, error) {
	key := Key{}

	shortTokenBytes := make([]byte, g.shortTokenLength)
	_, err := rand.Read(shortTokenBytes)
	if err != nil {
		return key, err
	}
	longTokenBytes := make([]byte, g.longTokenLength)
	_, err = rand.Read(longTokenBytes)
	if err != nil {
		return key, err
	}

	shortToken := padStart(base58.Encode(shortTokenBytes), g.shortTokenLength, "0")
	shortToken = (g.shortTokenPrefix + shortToken)
	key.shortToken = shortToken[0:g.shortTokenLength]

	key.longToken = padStart(base58.Encode(longTokenBytes), g.longTokenLength, "0")
	key.longToken = key.longToken[0:g.longTokenLength]

	key.longTokenHash = hashLongToken(key.longToken)

	key.token = fmt.Sprintf("%s_%s_%s", g.keyPrefix, key.shortToken, key.longToken)

	return key, nil
}

func ParseAPIKey(token string) (Key, error) {
	key := Key{}
	parts := strings.Split(token, "_")
	if len(parts) != 3 {
		return key, fmt.Errorf("invalid token")
	}
	key.token = token
	key.shortToken = parts[1]
	key.longToken = parts[2]
	key.longTokenHash = hashLongToken(key.longToken)
	return key, nil
}

func (k Key) Token() string {
	return k.token
}

func (k Key) ShortToken() string {
	return k.shortToken
}

func (k Key) LongToken() string {
	return k.longToken
}

func (k Key) LongTokenHash() string {
	return k.longTokenHash
}

func CheckAPIKey(token string, expectedLongTokenHash string) (bool, error) {
	key, err := ParseAPIKey(token)
	if err != nil {
		return false, err
	}
	return key.LongTokenHash() == expectedLongTokenHash, nil
}

func padStart(str string, len int, pad string) string {
	if n := utf8.RuneCountInString(str); n < len {
		str = strings.Repeat(pad, len-n) + string(str)
	}
	return str
}

func hashLongToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
