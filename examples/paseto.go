package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/oarkflow/pkg/paseto"
)

var data10 = []byte(`
{
	"type": "object",
	"properties": {
		"em": {
			"type": "object",
			"properties": {
				"code": {
					"type": "string"
				},
				"encounter_uid": {
					"type": "integer"
				},
				"work_item_uid": {
					"type": "integer"
				},
				"billing_provider": {
					"type": "string"
				},
				"resident_provider": {
					"type": "string"
				}
			}
		},
		"cpt": {
			"type" : "array",
			"items" : {
				"type": "object",
				"properties": {
					"code": {
						"type": "string"
					},
					"encounter_uid": {
						"type": "integer"
					},
					"work_item_uid": {
						"type": "integer"
					},
					"billing_provider": {
						"type": "string"
					},
					"resident_provider": {
						"type": "string"
					}
				}
			}
		}
	}
}
`)

func main() {
	start := time.Now()
	var secret = "OdR4DlWhZk6osDd0qXLdVT88lHOvj14K"
	token, err := GenerateApiKey(secret, data10)
	if err != nil {
		panic(err)
	}
	validatedKey := ValidateApiKey(secret, token)
	if validatedKey.Error != nil {
		panic(validatedKey.Error)
	}
	fmt.Println(time.Since(start))
}

func GenerateApiKey(secret string, payload []byte) (string, error) {
	claims := paseto.CustomClaim(payload)
	pv4 := paseto.NewPV4Local()
	symK, err := paseto.NewSymmetricKey([]byte(secret), paseto.Version4)
	if err != nil {
		return "", err
	}

	return pv4.Encrypt(symK, claims)
}

func ValidateApiKey(secret, token string) ValidApiKey {
	pv4 := paseto.NewPV4Local()
	symK, err := paseto.NewSymmetricKey([]byte(secret), paseto.Version4)
	if err != nil {
		return ValidApiKey{Error: err, Valid: false}
	}
	tk := pv4.Decrypt(token, symK)
	if tk.Err() != nil {
		return ValidApiKey{Error: tk.Err(), Valid: false}
	}

	if tk.HasFooter() {
		return ValidApiKey{Error: errors.New("footer was not passed to the library"), Valid: false}
	}

	var cc paseto.CustomClaim
	if err := tk.ScanClaims(&cc); err != nil {
		return ValidApiKey{Error: err, Valid: false}
	}
	return ValidApiKey{Payload: cc, Valid: true}
}

type ValidApiKey struct {
	Payload paseto.CustomClaim
	Valid   bool
	Error   error
}

func (v *ValidApiKey) Unmarshal(result interface{}) error {
	return json.Unmarshal(v.Payload, &result)
}
