package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func EncryptAes(key []byte, clearSourceFile string, encryptedDestinationFile string) error {
	if len(key) != 32 {
		return fmt.Errorf("error in AES key, must be 32 bytes but %v bytes found", len(key))
	}

	inFile, err := os.Open(clearSourceFile)
	if err != nil {
		return errors.New("Error opening clear source file for reading: " + err.Error())
	}
	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("Error creating cypher block: " + err.Error())
	}

	iv := make([]byte, block.BlockSize()) // Never use more than 2^32 random nonces with a given key because of the risk of repeat.
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return errors.New("Error creating IV: " + err.Error())
	}

	outFile, err := os.OpenFile(encryptedDestinationFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("Error opening encrpt destination file for writing: " + err.Error())
	}
	defer outFile.Close()

	buf := make([]byte, 1024) // The buffer size must be multiple of 16 bytes.
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			outFile.Write(buf[:n]) // Write into file.
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.New("Error encrypting file: " + err.Error())
		}
	}

	outFile.Write(iv)

	return nil
}

func DecryptAes(key []byte, encryptedSourceFile string, decryptedDestinationFile string) error {
	inFile, err := os.Open(encryptedSourceFile)
	if err != nil {
		return errors.New("Error opening encrypted source file for reading: " + err.Error())
	}
	defer inFile.Close()

	block, err := aes.NewCipher([]byte(key)) // The key should be 16 bytes (AES-128), 24 bytes (AES-192) or 32 bytes (AES-256)
	if err != nil {
		return errors.New("Error creating cipher block: " + err.Error())
	}

	fi, err := inFile.Stat()
	if err != nil {
		return errors.New("Error getting source file stats: " + err.Error())
	}

	iv := make([]byte, block.BlockSize()) // Never use more than 2^32 random nonces with a given key because of the risk of repeat.
	msgLen := fi.Size() - int64(len(iv))
	if _, err := inFile.ReadAt(iv, msgLen); err != nil {
		return errors.New("Error reading IV: " + err.Error())
	}

	outFile, err := os.OpenFile(decryptedDestinationFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("Error opening decrypt destination file: " + err.Error())
	}
	defer outFile.Close()

	buf := make([]byte, 1024) // The buffer size must be multiple of 16 bytes.
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			if n > int(msgLen) { // The last bytes are the IV, don't belong the original message
				n = int(msgLen)
			}

			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			outFile.Write(buf[:n]) // Write into file
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.New("Error decrypting file: " + err.Error())
		}
	}

	return nil
}
