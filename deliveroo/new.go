package deliveroo

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"strings"
	"sync"
)

const QueryUserAuth = `SELECT "user".auth_token FROM "user" WHERE "user".wii_id = $1 LIMIT 1`
const AuthKey = "3F4528472B4B6250655368566D597133"
const AuthIV = "3475377821412544"

type Deliveroo struct {
	mutex     sync.RWMutex
	auth      string
	longitude float64
	latitude  float64
	userId    string
}

func NewDeliveroo(db *pgxpool.Pool, r *http.Request) (*Deliveroo, error) {
	var encryptedAuth string
	row := db.QueryRow(context.Background(), QueryUserAuth, r.Header.Get("X-WiiID"))
	err := row.Scan(&encryptedAuth)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(AuthKey))
	if err != nil {
		return nil, err
	}

	auth := make([]byte, len(encryptedAuth))
	cipherText, err := hex.DecodeString(encryptedAuth)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(AuthIV))
	mode.CryptBlocks(auth, cipherText)
	auth = bytes.Trim(auth, "\x00")
	auth, err = pkcs7Unpad(auth, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	d := Deliveroo{
		mutex:     sync.RWMutex{},
		auth:      string(auth),
		longitude: 0,
		latitude:  0,
		userId:    "",
	}

	err = d.GetUserID()
	if err != nil {
		return nil, err
	}

	err = d.GetUserAddress()
	if err != nil {
		return nil, err
	}

	return &d, err
}

func (d *Deliveroo) GetUserID() error {
	token, err := base64.StdEncoding.DecodeString(strings.Replace(d.auth, "Basic ", "", -1))
	if err != nil {
		return err
	}

	d.userId = strings.Split(string(token), ":")[0]

	return nil
}
