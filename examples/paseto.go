package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oarkflow/paseto"

	"github.com/oarkflow/pkg/key"
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
	var secret = "OdR4DlWhZk6osDd0qXLdVT88lHOvj143"
	token := `KSmyIbv3l7rmKsNBvLtd0hNSQjh1PZdXXOkSPtFkpfq2jlHgqw_B2iZUve-HucnOKBPqSXzkdmji46VxgjZLCe89pogAfg`
	/*token, err := key.Generate(secret, []byte(`1`))
	if err != nil {
		panic(err)
	}
	fmt.Println(token)*/
	fmt.Println(key.Validate(secret, token))

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
