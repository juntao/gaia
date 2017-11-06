package stake

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/state"
)

// nolint
var (
	// Keys for store prefixes
	CandidatesPubKeysKey = []byte{0x01} // key for all candidates' pubkeys
	ParamKey             = []byte{0x02} // key for global parameters relating to staking

	// Key prefixes
	CandidateKeyPrefix     = []byte{0x03} // prefix for each key to a candidate
	DelegatorBondKeyPrefix = []byte{0x05} // prefix for each key to a delegator's bond
)

// GetCandidateKey - get the key for the candidate with pubKey
func GetCandidateKey(pubKey crypto.PubKey) []byte {
	return append(CandidateKeyPrefix, pubKey.Bytes()...)
}

func getDelegatorBondKey(delegator sdk.Actor, candidate crypto.PubKey) []byte {
	bondBytes := append(wire.BinaryBytes(&delegator), candidate.Bytes()...)
	return append(DelegatorBondKeyPrefix, bondBytes...)
}

/////////////////////////////////////////////////////////////////////////////////

// Get the active list of all the candidate pubKeys and owners
func loadCandidatesPubKeys(store state.SimpleDB) (pubKeys []crypto.PubKey) {
	bytes := store.Get(CandidatesPubKeysKey)
	if bytes == nil {
		return
	}
	err := wire.ReadBinaryBytes(bytes, &pubKeys)
	if err != nil {
		panic(err)
	}
	return
}
func saveCandidatesPubKeys(store state.SimpleDB, pubKeys []crypto.PubKey) {
	b := wire.BinaryBytes(pubKeys)
	store.Set(CandidatesPubKeysKey, b)
}

// LoadCandidate - loads the pubKey bond set
// TODO ultimately this function should be made unexported... being used right now
// for patchwork of tick functionality therefor much easier if exported until
// the new SDK is created
func LoadCandidate(store state.SimpleDB, pubKey crypto.PubKey) *Candidate {
	b := store.Get(GetCandidateKey(pubKey))
	if b == nil {
		return nil
	}
	candidate := new(Candidate)
	err := wire.ReadBinaryBytes(b, candidate)
	if err != nil {
		panic(err) // This error should never occure big problem if does
	}
	return candidate
}

func saveCandidate(store state.SimpleDB, candidate *Candidate) {

	if !store.Has(GetCandidateKey(candidate.PubKey)) {
		// TODO to be replaced with iteration in the multistore?
		pks := loadCandidatesPubKeys(store)
		saveCandidatesPubKeys(store, append(pks, candidate.PubKey))
	}

	b := wire.BinaryBytes(*candidate)
	store.Set(GetCandidateKey(candidate.PubKey), b)
}

func removeCandidate(store state.SimpleDB, pubKey crypto.PubKey) {
	store.Remove(GetCandidateKey(pubKey))

	// TODO to be replaced with iteration in the multistore?
	pks := loadCandidatesPubKeys(store)
	for i := range pks {
		if pks[i].Equals(pubKey) {
			saveCandidatesPubKeys(store,
				append(pks[:i], pks[i+1:]...))
			break
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////

func loadDelegatorBond(store state.SimpleDB,
	delegator sdk.Actor, candidate crypto.PubKey) *DelegatorBond {

	delegatorBytes := store.Get(getDelegatorBondKey(delegator, candidate))
	if delegatorBytes == nil {
		return nil
	}

	bond := new(DelegatorBond)
	err := wire.ReadBinaryBytes(delegatorBytes, bond)
	if err != nil {
		panic(err)
	}
	return bond
}

func saveDelegatorBond(store state.SimpleDB, delegator sdk.Actor, bond *DelegatorBond) {

	b := wire.BinaryBytes(*bond)
	store.Set(getDelegatorBondKey(delegator, bond.PubKey), b)
}

func removeDelegatorBond(store state.SimpleDB, delegator sdk.Actor, candidate crypto.PubKey) {
	store.Remove(getDelegatorBondKey(delegator, candidate))
}

/////////////////////////////////////////////////////////////////////////////////

// load/save the global staking params
func loadParams(store state.SimpleDB) (params Params) {
	b := store.Get(ParamKey)
	if b == nil {
		return defaultParams()
	}

	err := wire.ReadBinaryBytes(b, &params)
	if err != nil {
		panic(err) // This error should never occure big problem if does
	}

	return
}
func saveParams(store state.SimpleDB, params Params) {
	b := wire.BinaryBytes(params)
	store.Set(ParamKey, b)
}
