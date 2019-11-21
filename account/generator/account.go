package generator

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/status-go/extkeys"
	statusproto "github.com/status-im/status-protocol-go/types"
)

type account struct {
	privateKey  *ecdsa.PrivateKey
	extendedKey *extkeys.ExtendedKey
}

func (a *account) toAccountInfo() AccountInfo {
	publicKeyHex := statusproto.EncodeHex(crypto.FromECDSAPub(&a.privateKey.PublicKey))
	addressHex := crypto.PubkeyToAddress(a.privateKey.PublicKey).Hex()

	return AccountInfo{
		PublicKey: publicKeyHex,
		Address:   addressHex,
	}
}

func (a *account) toIdentifiedAccountInfo(id string) IdentifiedAccountInfo {
	info := a.toAccountInfo()
	keyUID := sha256.Sum256(crypto.FromECDSAPub(&a.privateKey.PublicKey))
	keyUIDHex := statusproto.EncodeHex(keyUID[:])
	return IdentifiedAccountInfo{
		AccountInfo: info,
		ID:          id,
		KeyUID:      keyUIDHex,
	}
}

func (a *account) toGeneratedAccountInfo(id string, mnemonic string) GeneratedAccountInfo {
	idInfo := a.toIdentifiedAccountInfo(id)
	return GeneratedAccountInfo{
		IdentifiedAccountInfo: idInfo,
		Mnemonic:              mnemonic,
	}
}

// AccountInfo contains a PublicKey and an Address of an account.
type AccountInfo struct {
	PublicKey string `json:"publicKey"`
	Address   string `json:"address"`
}

// IdentifiedAccountInfo contains AccountInfo and the ID of an account.
type IdentifiedAccountInfo struct {
	AccountInfo
	ID     string `json:"id"`
	KeyUID string `json:"keyUid"`
}

// GeneratedAccountInfo contains IdentifiedAccountInfo and the mnemonic of an account.
type GeneratedAccountInfo struct {
	IdentifiedAccountInfo
	Mnemonic string `json:"mnemonic"`
}

func (a GeneratedAccountInfo) toGeneratedAndDerived(derived map[string]AccountInfo) GeneratedAndDerivedAccountInfo {
	return GeneratedAndDerivedAccountInfo{
		GeneratedAccountInfo: a,
		Derived:              derived,
	}
}

// GeneratedAndDerivedAccountInfo contains GeneratedAccountInfo and derived AccountInfo mapped by derivation path.
type GeneratedAndDerivedAccountInfo struct {
	GeneratedAccountInfo
	Derived map[string]AccountInfo `json:"derived"`
}
