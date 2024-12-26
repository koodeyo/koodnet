package models

import (
	"fmt"
	"math"
	"net"
	"slices"
	"strings"
)

// Validators
func (n *Network) validate() error {
	if strings.TrimSpace(n.Name) == "" {
		return NewValidationError("name cannot be empty")
	}

	if n.Duration <= 0 {
		return NewValidationError("duration must be greater than 0")
	}

	if n.Encrypt {
		if err := n.validateEncryption(); err != nil {
			return err
		}
	}

	if len(n.IPs) == 0 {
		return NewValidationError("at least one IP is required")
	}

	if err := n.validateIPs(); err != nil {
		return err
	}

	if err := n.validateSubnets(); err != nil {
		return err
	}

	if err := n.validateGroups(); err != nil {
		return err
	}

	return nil
}

func (n *Network) validateEncryption() error {
	if len(n.Passphrase) == 0 {
		return NewValidationError("passphrase is required when encryption is enabled")
	}

	validCurves := []string{"25519", "X25519", "Curve25519", "CURVE25519", "P256"}
	if !slices.Contains(validCurves, n.Curve) {
		return NewValidationError("invalid curve; valid options are '25519' or 'P256'")
	}

	if n.ArgonMemory <= 0 || n.ArgonMemory > math.MaxUint32 {
		return fmt.Errorf("argon_memory must be greater than 0 and no more than %d KiB", uint32(math.MaxUint32))
	}

	if n.ArgonParallelism <= 0 || n.ArgonParallelism > math.MaxUint8 {
		return fmt.Errorf("argon_parallelism must be greater than 0 and no more than %d", math.MaxUint8)
	}

	if n.ArgonIterations <= 0 || n.ArgonIterations > math.MaxUint32 {
		return fmt.Errorf("argon_iterations must be greater than 0 and no more than %d", uint32(math.MaxUint32))
	}

	return nil
}

func (n *Network) validateIPs() error {
	for _, ip := range n.IPs {
		if _, _, err := net.ParseCIDR(ip); err != nil {
			return NewValidationError("invalid IP: " + ip)
		}
	}
	return nil
}

func (n *Network) validateSubnets() error {
	if len(n.Subnets) > 0 {
		for _, subnet := range n.Subnets {
			if _, _, err := net.ParseCIDR(subnet); err != nil {
				return NewValidationError("invalid subnet: " + subnet)
			}
		}
	}
	return nil
}

func (n *Network) validateGroups() error {
	if len(n.Groups) > 0 {
		for _, group := range n.Groups {
			if len(strings.TrimSpace(group)) == 0 {
				return NewValidationError("group names cannot be empty")
			}
		}
	}
	return nil
}
