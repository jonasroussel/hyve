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
	"github.com/jonasroussel/proxbee/tools"
)

var ActiveUser User

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

func LoadOrCreateUser() error {
	userDir := tools.Env.DataDir + "/user"

	user, err := loadUser(userDir)
	if os.IsNotExist(err) {
		user, err = createAccount(userDir)
	}

	if err != nil {
		return err
	}

	ActiveUser = *user

	return nil
}

func loadUser(userDir string) (*User, error) {
	user := User{}

	// Load the registration info

	rawReg, err := os.ReadFile(userDir + "/registration.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawReg, &user.Registration)
	if err != nil {
		return nil, err
	}

	// Load the private key

	rawKey, err := os.ReadFile(userDir + "/private.key")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(rawKey)
	if block == nil {
		return nil, errors.New("private key is not PEM encoded")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	user.PrivateKey = *key

	return &user, nil
}

func createAccount(userDir string) (*User, error) {
	// Create the user

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

	// Save his private key

	privateKeyDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	privKeyFile, err := os.Create(userDir + "/private.key")
	if err != nil {
		return nil, err
	}
	defer privKeyFile.Close()

	err = pem.Encode(privKeyFile, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	if err != nil {
		return nil, err
	}

	// Save his registration info

	regJSON, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(userDir+"/registration.json", regJSON, 0620)
	if err != nil {
		panic(err)
	}

	return &user, nil
}
