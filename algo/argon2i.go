package algo

import (
	"encoding/base64"

	"github.com/arjprd/crypt-service/driver"
	"golang.org/x/crypto/argon2"
)

const (
	ALGO_NAME_ARGON2I       = "argon2i"
	DEFAULT_ARGON2I_TIME    = 1
	DEFAULT_ARGON2I_MEMORY  = 64 * 1024
	DEFAULT_ARGON2I_THREADS = 1
	DEFAULT_ARGON2I_KEYLEN  = 32
)

type Argon2i struct {
	salt    []byte
	time    uint
	memory  uint32
	threads uint8
	keyLen  uint32
	c       *driver.Config
}

func NewArgon2iHash(salt string, time uint, memory uint32, threads uint8, keyLen uint32, c *driver.Config) HashAlgorithm {
	return &Argon2i{
		salt:    []byte(salt),
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		c:       c,
	}
}

func (a *Argon2i) Generate(password string) (string, error) {
	hash := argon2.Key([]byte(password), a.salt, uint32(a.time), a.memory, a.threads, a.keyLen)
	base64Hash := base64.StdEncoding.EncodeToString(hash)
	a.c.Logger().Info("argon2i hash generated %s", base64Hash)
	return base64Hash, nil
}

func (a *Argon2i) Verify(hash string, password string) bool {
	generatedhash, err := a.Generate(password)
	if err != nil {
		a.c.Logger().Error("argon2i hash generation failed: %+v", err)
		return false
	}
	if generatedhash != hash {
		a.c.Logger().Error("argon2i hash and password mismatch")
		return false
	}
	return true
}
