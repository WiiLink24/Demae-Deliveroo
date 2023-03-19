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

var AuthKey = []byte{0x3F, 0x45, 0x28, 0x47, 0x2B, 0x4B, 0x62, 0x50, 0x65, 0x53, 0x68, 0x56, 0x6D, 0x59, 0x71, 0x33}
var AuthIV = []byte{0x54, 0x68, 0x57, 0x6D, 0x5A, 0x71, 0x34, 0x74, 0x37, 0x77, 0x21, 0x7A, 0x25, 0x43, 0x2A, 0x46}

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
