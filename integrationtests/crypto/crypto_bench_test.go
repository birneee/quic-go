package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/chacha20poly1305"
	"runtime"
	"sync"
	"testing"
)

var aeadFactories = map[string]func(t assert.TestingT) (cipher.AEAD, [12]byte){
	"aes":    randomAES,
	"chacha": randomChaCha,
}

func randomAES(t assert.TestingT) (cipher.AEAD, [12]byte) {
	var key [16]byte
	_, err := rand.Read(key[:])
	assert.NoError(t, err)
	var nonce [12]byte
	_, err = rand.Read(nonce[:])
	assert.NoError(t, err)
	block, err := aes.NewCipher(key[:])
	assert.NoError(t, err)
	aesgcm, err := cipher.NewGCM(block)
	assert.NoError(t, err)
	return aesgcm, nonce
}

func randomChaCha(t assert.TestingT) (cipher.AEAD, [12]byte) {
	var key [32]byte
	_, err := rand.Read(key[:])
	assert.NoError(t, err)
	var nonce [12]byte
	_, err = rand.Read(nonce[:])
	assert.NoError(t, err)
	aead, err := chacha20poly1305.New(key[:])
	assert.NoError(t, err)
	return aead, nonce
}

func BenchmarkSeal(b *testing.B) {
	aead, nonce := randomAES(b)
	var buf [1516]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aead.Seal(buf[:], nonce[:], buf[:1500], nil)
	}
	b.StopTimer()
	b.ReportMetric(float64(1500*b.N)/b.Elapsed().Seconds()/1e9*8, "Gbps")
}

func BenchmarkSealParallel(b *testing.B) {
	for cipherName, aeadFactory := range aeadFactories {
		aead, nonce := aeadFactory(b)
		var plaintext [1500]byte
		b.Run(fmt.Sprintf("%s", cipherName), func(b *testing.B) {
			for workers := 1; workers <= runtime.NumCPU(); workers++ {
				b.Run(fmt.Sprintf("%d", workers), func(b *testing.B) {
					jobsPerWorker := b.N
					wg := sync.WaitGroup{}
					wg.Add(workers)
					b.ResetTimer()
					for i := 0; i < workers; i++ {
						go func() {
							var ciphertext [1516]byte
							for j := 0; j < jobsPerWorker; j++ {
								aead.Seal(ciphertext[:], nonce[:], plaintext[:], nil)
							}
							wg.Done()
						}()
					}
					wg.Wait()
					b.StopTimer()
					b.ReportMetric(float64(1500*b.N*workers)/b.Elapsed().Seconds()/1e9*8, "Gbps")
				})
			}
		})
	}
}

func BenchmarkOpen(b *testing.B) {
	aead, nonce := randomAES(b)
	var buf [1500]byte
	_, err := rand.Read(buf[:])
	assert.NoError(b, err)
	ciphertext := aead.Seal(nil, nonce[:], buf[:], nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := aead.Open(buf[:], nonce[:], ciphertext, nil)
		assert.NoError(b, err)
	}
	b.StopTimer()

	b.ReportMetric(float64(1500*b.N)/b.Elapsed().Seconds()/1e9*8, "Gbps")
}

func BenchmarkOpenParallel(b *testing.B) {
	aead, nonce := randomAES(b)
	var ciphertextBuf [1516]byte
	_, err := rand.Read(ciphertextBuf[:])
	assert.NoError(b, err)
	ciphertext := aead.Seal(nil, nonce[:], ciphertextBuf[:1500], nil)
	for workers := 1; workers <= runtime.NumCPU(); workers++ {
		b.Run(fmt.Sprintf("%d", workers), func(b *testing.B) {
			jobsPerWorker := b.N
			wg := sync.WaitGroup{}
			wg.Add(workers)
			b.ResetTimer()
			for i := 0; i < workers; i++ {
				go func() {
					var plaintextBuf [1500]byte
					for j := 0; j < jobsPerWorker; j++ {
						_, err := aead.Open(plaintextBuf[:], nonce[:], ciphertext, nil)
						assert.NoError(b, err)
					}
					wg.Done()
				}()
			}
			wg.Wait()
			b.StopTimer()

			b.ReportMetric(float64(b.N*1500*workers)/b.Elapsed().Seconds()/1e9*8, "Gbps")
		})
	}
}
