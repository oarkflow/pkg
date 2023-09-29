package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/enc"
	"github.com/oarkflow/pkg/str"
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

var data = `
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
`

var secret = "OdR4DlWhZk6osDd0qXLdVT88lHOvj14K"

func encTestWithCompressed() {
	compressed := str.ToCompressedString(data10)
	t, e := enc.Encrypt(compressed, secret)
	if e != nil {
		panic(e)
	}
	c, e := enc.Decrypt(t, secret)
	if e != nil {
		panic(e)
	}
	bt, err := str.FromCompressedString(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bt), t)
}

func encTest() {
	start := time.Now()
	t, e := enc.Encrypt(data, secret)
	if e != nil {
		panic(e)
	}
	_, e = enc.Decrypt(t, secret)
	if e != nil {
		panic(e)
	}
	fmt.Println(time.Since(start))
	// fmt.Println(c, t)
}

func main() {
	encTest()
}
