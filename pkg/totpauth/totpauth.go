package totpauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func Auth(key string, timestart uint64, timestep uint64) string {
	return totp(key, timestart, timestep)
}

func totp(key string, timestart uint64, timestep uint64) string {
	counter := (uint64(time.Now().Unix()) - timestart) / timestep
	// counter := uint64(float64(unixTime-t0) / float64(timeStep)) //prior method
	return hotp(key, counter)
}

func hotp(key string, counter uint64) string {
	hash := computeHmacSha1(key, counter)
	otpCode := truncate(hash, 6)
	return otpCode
}

func truncate(hash []byte, digitCount int) string {
	// Calculate the offset from the last byte of the hash
	offset := uintptr(hash[len(hash)-1] & 0xf)

	// Extract 4 bytes from the hash starting from the calculated offset
	extractedValue := int(hash[offset])<<24 | int(hash[offset+1])<<16 | int(hash[offset+2])<<8 | int(hash[offset+3])

	// Ensure the extracted value is positive by BITWISE AND operation
	extractedPositiveValue := int64(extractedValue & 0x7fffffff)

	otp := int(int64(extractedPositiveValue) % int64(math.Pow10(digitCount)))

	// Format the OTP with leading zeros to fit the digit count
	formattedOTP := fmt.Sprintf(fmt.Sprintf("%%0%dd", digitCount), otp)

	return formattedOTP
}

func computeHmacSha1(key string, msg uint64) []byte {
	upperedLetterKey := strings.ToUpper(key)
	keyByte, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(upperedLetterKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding key: %v\n", err)
	}

	msgByte := make([]byte, 8)
	binary.BigEndian.PutUint64(msgByte, msg)

	hash := hmac.New(sha1.New, keyByte)
	hash.Write(msgByte)

	return hash.Sum(nil)
}
