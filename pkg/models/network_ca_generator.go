package models

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/slackhq/nebula/cert"
)

// Helpers
func (n *Network) getArgon2Parameters() *cert.Argon2Parameters {
	return cert.NewArgon2Parameters(
		uint32(n.ArgonMemory),
		uint8(n.ArgonParallelism),
		uint32(n.ArgonIterations),
	)
}

func (n *Network) getIPs() []*net.IPNet {
	var ipNets []*net.IPNet
	for _, ip := range n.IPs {
		_, ipNet, _ := net.ParseCIDR(ip)
		ipNets = append(ipNets, ipNet)
	}

	return ipNets
}

func (n *Network) getSubnets() []*net.IPNet {
	var subnetNets []*net.IPNet
	for _, subnet := range n.Subnets {
		_, subnetNet, _ := net.ParseCIDR(subnet)
		subnetNets = append(subnetNets, subnetNet)
	}

	return subnetNets
}

// Generate new certificate authority
func (n *Network) NewCA() (*Certificate, error) {
	var curve cert.Curve
	var pub, rawPriv []byte
	var err error

	switch n.Curve {
	case "25519", "X25519", "Curve25519", "CURVE25519":
		curve = cert.Curve_CURVE25519
		pub, rawPriv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("error generating Ed25519 keys: %s", err)
		}
	case "P256":
		curve = cert.Curve_P256
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("error generating ECDSA keys: %s", err)
		}
		eKey, err := key.ECDH()
		if err != nil {
			return nil, fmt.Errorf("error converting ECDSA key: %s", err)
		}
		rawPriv = eKey.Bytes()
		pub = eKey.PublicKey().Bytes()
	default:
		return nil, fmt.Errorf("unsupported curve: %s", n.Curve)
	}

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      n.Name,
			Groups:    n.Groups,
			Ips:       n.getIPs(),
			Subnets:   n.getSubnets(),
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(n.Duration),
			PublicKey: pub,
			IsCA:      true,
			Curve:     curve,
		},
	}

	err = nc.Sign(curve, rawPriv)
	if err != nil {
		return nil, fmt.Errorf("error signing certificate: %s", err)
	}

	var keyBytes []byte
	if n.Encrypt {
		passphrase := []byte(n.Passphrase)
		kdfParams := n.getArgon2Parameters()

		keyBytes, err = cert.EncryptAndMarshalSigningPrivateKey(curve, rawPriv, passphrase, kdfParams)
		if err != nil {
			return nil, fmt.Errorf("error encrypting key: %s", err)
		}
	} else {
		keyBytes = cert.MarshalSigningPrivateKey(curve, rawPriv)
	}

	certBytes, err := nc.MarshalToPEM()
	if err != nil {
		return nil, fmt.Errorf("error marshalling certificate: %s", err)
	}

	pubBytes := cert.MarshalPublicKey(curve, pub)

	return &Certificate{
		IsCA:       true,
		ID:         uuid.New(),
		NotBefore:  nc.Details.NotBefore,
		NotAfter:   nc.Details.NotAfter,
		Key:        keyBytes,
		Pub:        pubBytes,
		Crt:        certBytes,
		Passphrase: n.Passphrase,
		IPs:        n.IPs,
		Groups:     n.Groups,
		Subnets:    n.Subnets,
	}, nil
}
