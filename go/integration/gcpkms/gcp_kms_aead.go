// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

// Package gcpkms provides for the GCP Cloud KMS integration
package gcpkms

import (
	"encoding/base64"

	"google.golang.org/api/cloudkms/v1"

	"github.com/google/tink/go/tink"
)

// GCPAead represents a GCP KMS service to a particular URI.
type GCPAead struct {
	keyURI string
	kms    cloudkms.Service
}

var _ tink.AEAD = (*GCPAead)(nil)

// NewGCPAead returns a new GCP KMS service.
func NewGCPAead(keyURI string, kms *cloudkms.Service) *GCPAead {
	return &GCPAead{
		keyURI: keyURI,
		kms:    *kms,
	}
}

// Encrypt AEAD encrypts the plaintext data and uses addtionaldata from authentication.
func (a *GCPAead) Encrypt(plaintext, additionalData []byte) ([]byte, error) {

	req := &cloudkms.EncryptRequest{
		Plaintext:                   base64.StdEncoding.EncodeToString(plaintext),
		AdditionalAuthenticatedData: base64.StdEncoding.EncodeToString(additionalData),
	}
	resp, err := a.kms.Projects.Locations.KeyRings.CryptoKeys.Encrypt(a.keyURI, req).Do()
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Ciphertext)
}

// Decrypt AEAD decrypts the data and verified the additional data.
func (a *GCPAead) Decrypt(ciphertext, additionalData []byte) ([]byte, error) {

	req := &cloudkms.DecryptRequest{
		Ciphertext:                  base64.StdEncoding.EncodeToString(ciphertext),
		AdditionalAuthenticatedData: base64.StdEncoding.EncodeToString(additionalData),
	}
	resp, err := a.kms.Projects.Locations.KeyRings.CryptoKeys.Decrypt(a.keyURI, req).Do()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}
