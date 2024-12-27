package models

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/slackhq/nebula/cert"
	"golang.org/x/crypto/curve25519"
)

func (h *Host) getSubnets() []*net.IPNet {
	var subnetNets []*net.IPNet
	for _, subnet := range h.Subnets {
		_, subnetNet, _ := net.ParseCIDR(subnet)
		subnetNets = append(subnetNets, subnetNet)
	}

	return subnetNets
}

func (h *Host) NewCert(ca Certificate) (*Certificate, error) {
	var err error
	var curve cert.Curve
	var caKey []byte

	caKey, _, curve, err = cert.UnmarshalSigningPrivateKey(ca.Key)
	if err == cert.ErrPrivateKeyEncrypted {
		// Convert the passphrase from string to byte slice
		passphrase := []byte(ca.Passphrase)

		if len(passphrase) == 0 {
			return nil, fmt.Errorf("cannot open encrypted ca-key without passphrase")
		}

		curve, caKey, _, err = cert.DecryptAndUnmarshalSigningPrivateKey(passphrase, ca.Key)
		if err != nil {
			return nil, fmt.Errorf("error while parsing encrypted ca-key: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error while parsing ca-key: %s", err)
	}

	caCert, _, err := cert.UnmarshalNebulaCertificateFromPEM(ca.Crt)
	if err != nil {
		return nil, fmt.Errorf("error while parsing ca-crt: %s", err)
	}

	if err := caCert.VerifyPrivateKey(curve, caKey); err != nil {
		return nil, fmt.Errorf("refusing to sign, root certificate does not match private key")
	}

	issuer, err := caCert.Sha256Sum()
	if err != nil {
		return nil, fmt.Errorf("error while getting -ca-crt fingerprint: %s", err)
	}

	if caCert.Expired(time.Now()) {
		return nil, fmt.Errorf("ca certificate is expired")
	}

	ip, ipNet, err := net.ParseCIDR(h.IP)
	if err != nil {
		return nil, fmt.Errorf("invalid ip definition: %s", err)
	}
	if ip.To4() == nil {
		return nil, fmt.Errorf("invalid ip definition: can only be ipv4, have %s", h.IP)
	}
	ipNet.IP = ip

	var pub, rawPriv []byte
	if len(h.InPub) > 0 {
		var pubCurve cert.Curve
		pub, _, pubCurve, err = cert.UnmarshalPublicKey(h.InPub)
		if err != nil {
			return nil, fmt.Errorf("error while parsing in-pub: %s", err)
		}
		if pubCurve != curve {
			return nil, fmt.Errorf("curve of in-pub does not match ca")
		}
	} else {
		pub, rawPriv = newKeypair(curve)
	}

	// Calculate the duration
	duration := time.Until(caCert.Details.NotAfter) - time.Hour

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      h.Name,
			Ips:       []*net.IPNet{ipNet},
			Groups:    h.Groups,
			Subnets:   h.getSubnets(),
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(duration),
			PublicKey: pub,
			IsCA:      false,
			Issuer:    issuer,
			Curve:     curve,
		},
	}

	if err := nc.CheckRootConstrains(caCert); err != nil {
		return nil, fmt.Errorf("refusing to sign, root certificate constraints violated: %s", err)
	}

	err = nc.Sign(curve, caKey)
	if err != nil {
		return nil, fmt.Errorf("error while signing: %s", err)
	}

	keyBytes := cert.MarshalPrivateKey(curve, rawPriv)

	certBytes, err := nc.MarshalToPEM()
	if err != nil {
		return nil, fmt.Errorf("error while marshalling certificate: %s", err)
	}

	pubBytes := cert.MarshalPublicKey(curve, pub)

	return &Certificate{
		IsCA:       false,
		ID:         uuid.New(),
		NotBefore:  nc.Details.NotBefore,
		NotAfter:   nc.Details.NotAfter,
		Passphrase: ca.Passphrase,
		Key:        keyBytes,
		Pub:        pubBytes,
		Crt:        certBytes,
	}, nil
}

// newKeypair generates a new keypair based on the specified curve
func newKeypair(curve cert.Curve) ([]byte, []byte) {
	switch curve {
	case cert.Curve_CURVE25519:
		return x25519Keypair()
	case cert.Curve_P256:
		return p256Keypair()
	default:
		return nil, nil
	}
}

func x25519Keypair() ([]byte, []byte) {
	privkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, privkey); err != nil {
		panic(err)
	}

	pubkey, err := curve25519.X25519(privkey, curve25519.Basepoint)
	if err != nil {
		panic(err)
	}

	return pubkey, privkey
}

func p256Keypair() ([]byte, []byte) {
	privkey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey := privkey.PublicKey()
	return pubkey.Bytes(), privkey.Bytes()
}
