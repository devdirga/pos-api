package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
)

// branch staging
// commit 1
// commit 2

// feature/dev-staging1
// commit1

// feature/two
// commit 1
// commit 2

// fix/staging-1
// commit 1
// commit 2

type Configuration struct {
	EncryptionKey string
	DbHost        string
	DbPort        string
	DbName        string
	DbUsername    string
	DbPassword    string
}

var configuration Configuration

type ModelBase struct {
	Id        int       `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	ModelBase        // replaces gorm.Model
	Email     string `gorm:"not null; unique"`
	Password  string `gorm:"not null" json:"-"`
	Name      string `gorm:"not null; type:varchar(100)"` // unique_index
}

type Authentication struct {
	ModelBase
	User   User `gorm:"foreignkey:UserID; not null"`
	UserID int
	Token  string `gorm:"type:varchar(200); not null"`
}

type CustomHTTPSuccess struct {
	Data string `json:"data"`
}

type ErrorType struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CustomHTTPError struct {
	Error ErrorType `json:"error"`
}

var GormDB *gorm.DB

func init() {
	configuration = Configuration{
		EncryptionKey: "F61L8L7CUCGN0NK6336I8TFP9Y2ZOS43",
		DbHost:        "127.0.0.1",
		DbPort:        "5432",
		DbName:        "majoo",
		DbUsername:    "postgres",
		DbPassword:    "admin123",
	}
}

func GetConfig() Configuration {
	return configuration
}

func ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return
}

func EncryptString(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func DecryptString(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
