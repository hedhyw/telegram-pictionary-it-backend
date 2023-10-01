// Package telegram contains the common functions and structures for working
// with Telegram messenger such as decoding, verification, etc.
package telegram

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"

	"github.com/gorilla/schema"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/samber/lo"
)

// Decoder helps to decode initData metadata. See `DecodeInitData` for more details.
type Decoder struct {
	botToken      string
	schemaDecoder *schema.Decoder
}

// NewDecoder creates a new telegram.Decoder. BotToken is required.
// See `DecodeInitData` for more details.
func NewDecoder(botToken string) *Decoder {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)

	return &Decoder{
		botToken:      botToken,
		schemaDecoder: schemaDecoder,
	}
}

// DecodeInitData decodes URL-encoded "initData" and verifies its integrity.
// It compares the received hash parameter with the hexadecimal representation
// of the HMAC-SHA-256 signature of the data-check-string with the secret key,
// which is the HMAC-SHA-256 signature of the bot's token with the constant
// string WebAppData used as a key.
func (d *Decoder) DecodeInitData(initDataRaw string) (*InitDataMeta, error) {
	metaValues, err := url.ParseQuery(initDataRaw)
	if err != nil {
		return nil, semerr.NewBadRequestError(fmt.Errorf("parsing query: %w", err))
	}

	err = verifyHash(metaValues, d.botToken)
	if err != nil {
		return nil, fmt.Errorf("verifying: %w", err)
	}

	var meta InitDataMeta
	err = d.schemaDecoder.Decode(&meta, metaValues)
	if err != nil {
		return nil, fmt.Errorf("decoding meta: %w", err)
	}

	meta.Raw = metaValues

	return &meta, nil
}

// verifyHash verifies the hash of URL-encoded WebAppData.
// values should include "hash" key.
func verifyHash(values url.Values, botToken string) error {
	const (
		keyHash = "hash"

		labelData = "WebAppData"
	)

	if botToken == "" {
		return semerr.NewBadRequestError(semerr.Error("empty bot token"))
	}

	actualHash := values.Get(keyHash)
	if actualHash == "" {
		return semerr.NewBadRequestError(semerr.Error("empty hash"))
	}

	dataCheckKeys := lo.Keys(values)
	sort.Strings(dataCheckKeys)

	var dataCheckBuf bytes.Buffer

	for i, key := range dataCheckKeys {
		if key == keyHash {
			continue
		}

		if i > 0 {
			dataCheckBuf.WriteByte('\n')
		}

		dataCheckBuf.WriteString(key)
		dataCheckBuf.WriteByte('=')
		dataCheckBuf.WriteString(values.Get(key))
	}

	secretKey := calculateHashSum(botToken, []byte(labelData))
	expectedHash := hex.EncodeToString(
		calculateHashSum(dataCheckBuf.String(), secretKey),
	)

	if expectedHash != actualHash {
		return semerr.NewBadRequestError(semerr.Error("invalid hash"))
	}

	return nil
}

// calculateHashSum generates a HMAC-SHA256 hash for provided value
// using the given secret.
func calculateHashSum(value string, secret []byte) []byte {
	hasher := hmac.New(sha256.New, secret)

	hasher.Write([]byte(value))

	return hasher.Sum(nil)
}
