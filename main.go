package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	file, _ := os.Create(filename)
	defer file.Close()
	file.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("welcome to file encryptooooor")
	fmt.Println("Enter e to encrypt data to a new file, Enter ef to encrypt an existing file,")
	fmt.Println("Enter d to decrypt a file, Enter dns to decrypt a file and save decrypted data into a new file")
	fmt.Println("press control + c to exit program")
	fmt.Println("-----------------------------------------------")

	scanner.Scan()
	userChoice := scanner.Text()

	var filename string
	var data string
	var password string

	if userChoice == "" {
		println("no input please restart program")

	}

	if userChoice == "e" {

		//get new filename
		fmt.Println("Enter file name")
		fmt.Scan(&filename)

		//get data
		fmt.Println("Enter data to encrypt")
		fmt.Scan(&data)

		//get password
		fmt.Println("Enter password")
		fmt.Scan(&password)

		//encrypt file
		encryptFile(filename, []byte(data), password)

		println("data encryted")

	}

	if userChoice == "ef" {
		var existingfilename string
		//get existing filename
		println("Enter Existing Filename to decrypt")
		fmt.Scan(&existingfilename)

		existingdata, _ := ioutil.ReadFile(existingfilename)

		println("enter new filename")
		fmt.Scan(&filename)

		//get file password
		println("Enter password to decrypt")
		fmt.Scan(&password)

		//encrypt existing file
		encryptFile(filename, existingdata, password)

		println("file encrypted")
	}

	if userChoice == "d" {
		//get filename
		println("Enter Filename to decrypt")
		fmt.Scan(&filename)

		//get file password
		println("Enter password to decrypt")
		fmt.Scan(&password)

		//decrypt file
		decrypted := decryptFile(filename, password)

		//print decrypted data in string
		println("decryted file data :", string(decrypted))
	}

	if userChoice == "dns" {
		//get filename
		println("Enter Filename to decrypt")
		fmt.Scan(&filename)

		var newfilename string
		println("Enter New filename to save decrypted file data as seperated file")
		fmt.Scan(&newfilename)

		//get file password
		println("Enter password to decrypt existing file")
		fmt.Scan(&password)

		//decrypt file
		decrypted := decryptFile(filename, password)

		//save decrypted data in a new file
		file, _ := os.Create(newfilename)
		defer file.Close()
		file.Write(decrypted)

	}

}
