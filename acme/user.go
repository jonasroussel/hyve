package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/jonasroussel/proxbee/config"
)

var LetsEncryptUser User

type User struct {
	Registration *registration.Resource
	PrivateKey   ecdsa.PrivateKey
}

func (u User) GetEmail() string {
	return ""
}
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u User) GetPrivateKey() crypto.PrivateKey {
	return u.PrivateKey
}

func LoadUser() error {
	userDir := config.DATA_DIR + "/user"

	rawReg, err := os.ReadFile(userDir + "/registration.json")
	if os.IsNotExist(err) {
		user, err := CreateAccount(userDir)
		if err != nil {
			return err
		}

		LetsEncryptUser = *user

		return nil
	} else if err != nil {
		return err
	} else {
		user := User{}

		err = json.Unmarshal(rawReg, &user.Registration)
		if err != nil {
			return err
		}

		rawKey, err := os.ReadFile(userDir + "/private.key")
		if err != nil {
			return err
		}
		block, _ := pem.Decode(rawKey)
		if block == nil {
			return errors.New("private key is not PEM encoded")
		}

		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		user.PrivateKey = *key

		LetsEncryptUser = user

		return nil
	}
}

func CreateAccount(dataDir string) (*User, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	user := User{PrivateKey: *privateKey}
	config := lego.NewConfig(user)

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.Registration = reg

	privateKeyDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	pkOutFile, err := os.Create(dataDir + "/private.key")
	if err != nil {
		return nil, err
	}
	defer pkOutFile.Close()

	err = pem.Encode(pkOutFile, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	if err != nil {
		return nil, err
	}

	regJSON, err := json.Marshal(*reg)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(dataDir+"/registration.json", regJSON, 0620)
	if err != nil {
		panic(err)
	}

	return &user, nil
}
