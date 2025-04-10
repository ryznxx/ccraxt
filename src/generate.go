package src

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

func randomSize(min, max int64) int64 {
	return min + randInt64(max-min)
}

func randInt64(n int64) int64 {
	if n <= 0 {
		panic("Invalid range")
	}
	num, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		panic(err)
	}
	return num.Int64()
}

func RandomFileName() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	name := make([]byte, 8)
	for i := range name {
		name[i] = chars[randInt64(int64(len(chars)))]
	}
	return fmt.Sprintf("%s.dll", string(name))
}

func CreateRandomFile(fileName string) float32 {
	fileSize := randomSize(20*1024*1024, 50*1024*1024) // Min 20MB, Max 50MB

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return 0
	}
	defer file.Close()

	buffer := make([]byte, 4096) // 4KB buffer
	var written int64

	for written < fileSize {
		toWrite := fileSize - written
		if toWrite > int64(len(buffer)) {
			toWrite = int64(len(buffer))
		}
		_, err := file.Write(buffer[:toWrite])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return 0
		}
		written += toWrite
	}

	return float32(fileSize) // Tetap return dalam bytes
}
