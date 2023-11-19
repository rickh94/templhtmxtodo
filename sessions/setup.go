package sessions

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"templtodo3/config"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

var SessionManager *scs.SessionManager

func Init() {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.AppConfig.RedisURI)
		},
	}

	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour
	SessionManager.Store = redisstore.New(pool)
}

func PutEncryptedString(c context.Context, key string, value string) {
	SessionManager.Put(c, key, encrypt(value))
}

func GetEncryptedString(c context.Context, key string) string {
	if SessionManager.GetString(c, key) == "" {
		return ""
	}
	return decrypt(SessionManager.GetString(c, key))
}

func encrypt(plaintext string) string {
	aes, err := aes.NewCipher([]byte(config.AppConfig.SecretKey[0:32]))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext)
}

func decrypt(ciphertext string) string {
	aes, err := aes.NewCipher([]byte(config.AppConfig.SecretKey[0:32]))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}
