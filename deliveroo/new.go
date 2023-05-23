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

const QueryUserAuth = `SELECT "user".auth_token, "user".roo_uid FROM "user" WHERE "user".wii_id = $1 LIMIT 1`

var AuthKey = []byte{0x3F, 0x45, 0x28, 0x47, 0x2B, 0x4B, 0x62, 0x50, 0x65, 0x53, 0x68, 0x56, 0x6D, 0x59, 0x71, 0x33}
var AuthIV = []byte{0x54, 0x68, 0x57, 0x6D, 0x5A, 0x71, 0x34, 0x74, 0x37, 0x77, 0x21, 0x7A, 0x25, 0x43, 0x2A, 0x46}

type Deliveroo struct {
	mutex        sync.RWMutex
	auth         string
	longitude    float64
	latitude     float64
	userId       string
	uID          string
	response     string
	responseCode int
}

func (d *Deliveroo) SetResponse(response string) {
	d.response = response
}

func (d *Deliveroo) Response() string {
	return d.response
}

func (d *Deliveroo) ResponseCode() int {
	return d.responseCode
}

func NewDeliveroo(db *pgxpool.Pool, r *http.Request) (*Deliveroo, error) {
	var encryptedAuth string
	var rooUID string
	row := db.QueryRow(context.Background(), QueryUserAuth, r.Header.Get("X-WiiID"))
	err := row.Scan(&encryptedAuth, &rooUID)
	if err != nil {
		return &Deliveroo{response: "Database error", responseCode: http.StatusInternalServerError}, err
	}

	block, err := aes.NewCipher(AuthKey)
	if err != nil {
		return &Deliveroo{response: "AES error", responseCode: http.StatusInternalServerError}, err
	}

	auth := make([]byte, len(encryptedAuth))
	cipherText, err := hex.DecodeString(encryptedAuth)
	if err != nil {
		return &Deliveroo{response: "hex.DecodeString error", responseCode: http.StatusInternalServerError}, err
	}

	mode := cipher.NewCBCDecrypter(block, AuthIV)
	mode.CryptBlocks(auth, cipherText)
	auth = bytes.Trim(auth, "\x00")
	auth, err = pkcs7Unpad(auth, aes.BlockSize)
	if err != nil {
		return &Deliveroo{response: "Invalid pkcs7 padding: problem with bot?", responseCode: http.StatusInternalServerError}, err
	}

	d := Deliveroo{
		mutex:        sync.RWMutex{},
		auth:         string(auth),
		longitude:    0,
		latitude:     0,
		userId:       "",
		uID:          rooUID,
		response:     "",
		responseCode: http.StatusInternalServerError,
	}

	err = d.GetUserID()
	if err != nil {
		return &Deliveroo{response: "Failed to get user id", responseCode: http.StatusUnauthorized}, err
	}

	err = d.GetUserAddress()
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
