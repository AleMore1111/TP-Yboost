package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Decrypts data using AES-CFB.
func decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Extract the nonce from the encrypted data
	nonce := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	// Use CFB stream decryption
	stream := cipher.NewCFBDecrypter(block, nonce)
	decrypted := make([]byte, len(encryptedData))
	stream.XORKeyStream(decrypted, encryptedData)

	return decrypted, nil
}

func main() {
	// Define the folder path
	folderPath := `C:\Users\alexis\Desktop\malware`

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Println("Error: The folder does not exist at this location:", folderPath)
		return
	}

	// Prompt the user to enter the AES key
	fmt.Print("Enter the AES key (hex-encoded): ")
	reader := bufio.NewReader(os.Stdin)
	keyHex, _ := reader.ReadString('\n')
	keyHex = strings.TrimSpace(keyHex) // Remove any whitespace characters

	// Decode the hex-encoded key
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Println("Error decoding the key:", err)
		return
	}

	fmt.Println("Key decoded successfully")

	// Walk through all files in the folder and decrypt them
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println("Processing file:", path)

			// Read the encrypted file content
			encryptedData, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading the file:", err)
				return err
			}

			// Decrypt the file content
			decryptedData, err := decrypt(encryptedData, key)
			if err != nil {
				fmt.Println("Error decrypting the data:", err)
				return err
			}

			// Write the decrypted data back to the file
			if err := ioutil.WriteFile(path, decryptedData, 0644); err != nil {
				fmt.Println("Error writing the decrypted data to the file:", err)
				return err
			}

			fmt.Println("File decrypted successfully:", path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the folder:", err)
		return
	}

	fmt.Println("Folder decrypted successfully.")
}
