package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/oarkflow/pkg/otp"
	"github.com/oarkflow/pkg/otp/totp"
)

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	key, err := totp.GenerateWithOpts(
		totp.WithIssuer("edelberg"),
		totp.WithAccountName("s.baniya.np@gmail.com"),
		totp.WithGenDigits(otp.DigitsSix),
		totp.WithGenPeriod(60),
	)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	os.WriteFile("qr-code.png", buf.Bytes(), 0644)

	fmt.Println("Validating TOTP...")
	fmt.Println("Scan the generated QR code and enter the generated number below")
	passcode := promptForPasscode()
	valid, err := totp.ValidateWithOpts(passcode, key.Secret(), totp.WithAlgorithm(otp.AlgorithmSHA1))

	if err != nil {
		panic(err)
	}
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
