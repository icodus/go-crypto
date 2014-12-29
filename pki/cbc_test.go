package pki

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CBC_Bijection(t *testing.T) {
	t.Parallel()

	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.Nil(t, err)

	pkS3FS := PrivateKey(*rsaKey)
	pubS3FS := PublicKey(rsaKey.PublicKey)

	var blobIn []byte
	var blobOut []byte
	var encryptionBlock EncryptionBlock

	// NO PADDING
	blobIn = make([]byte, 96)
	_, err = rand.Reader.Read(blobIn)
	assert.Nil(t, err)

	encryptionBlock, err = pubS3FS.EncryptCBC(blobIn)
	assert.Nil(t, err)

	blobOut, err = pkS3FS.DecryptCBC(encryptionBlock)
	assert.Nil(t, err)

	assert.Equal(t, blobIn, blobOut)

	// SOME PADDING
	blobIn = make([]byte, 97)
	_, err = rand.Reader.Read(blobIn)
	assert.Nil(t, err)

	encryptionBlock, err = pubS3FS.EncryptCBC(blobIn)
	assert.Nil(t, err)

	blobOut, err = pkS3FS.DecryptCBC(encryptionBlock)
	assert.Nil(t, err)

	assert.Equal(t, blobIn, blobOut)

	// < aes.BlockSize
	blobIn = make([]byte, aes.BlockSize-2)
	_, err = rand.Reader.Read(blobIn)
	assert.Nil(t, err)

	encryptionBlock, err = pubS3FS.EncryptCBC(blobIn)
	assert.Nil(t, err)

	blobOut, err = pkS3FS.DecryptCBC(encryptionBlock)
	assert.Nil(t, err)

	assert.Equal(t, blobIn, blobOut)
}
