package tendermint

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"encoding/json"

	ttypes "gitlab.33.cn/chain33/chain33/consensus/drivers/tendermint/types"
	"gitlab.33.cn/chain33/chain33/types"
)

// database keys
var (
	stateKey = []byte("stateKey")
)

//-----------------------------------------------------------------------------

// State is a short description of the latest committed block of the Tendermint consensus.
// It keeps all information necessary to validate new blocks,
// including the last validator set and the consensus params.
// All fields are exposed so the struct can be easily serialized,
// but none of them should be mutated directly.
// Instead, use state.Copy() or state.NextState(...).
// NOTE: not goroutine-safe.
type State struct {
	// Immutable
	ChainID string

	// LastBlockHeight=0 at genesis (ie. block(H=0) does not exist)
	LastBlockHeight  int64
	LastBlockTotalTx int64
	LastBlockID      ttypes.BlockID
	LastBlockTime    int64

	// LastValidators is used to validate block.LastCommit.
	// Validators are persisted to the database separately every time they change,
	// so we can query for historical validator sets.
	// Note that if s.LastBlockHeight causes a valset change,
	// we set s.LastHeightValidatorsChanged = s.LastBlockHeight + 1
	Validators                  *ttypes.ValidatorSet
	LastValidators              *ttypes.ValidatorSet
	LastHeightValidatorsChanged int64

	// Consensus parameters used for validating blocks.
	// Changes returned by EndBlock and updated after Commit.
	ConsensusParams                  ttypes.ConsensusParams
	LastHeightConsensusParamsChanged int64

	// Merkle root of the results from executing prev block
	LastResultsHash []byte

	// The latest AppHash we've received from calling abci.Commit()
	AppHash []byte
}

// Copy makes a copy of the State for mutating.
func (s State) Copy() State {
	return State{
		ChainID: s.ChainID,

		LastBlockHeight:  s.LastBlockHeight,
		LastBlockTotalTx: s.LastBlockTotalTx,
		LastBlockID:      s.LastBlockID,
		LastBlockTime:    s.LastBlockTime,

		Validators:                  s.Validators.Copy(),
		LastValidators:              s.LastValidators.Copy(),
		LastHeightValidatorsChanged: s.LastHeightValidatorsChanged,

		ConsensusParams:                  s.ConsensusParams,
		LastHeightConsensusParamsChanged: s.LastHeightConsensusParamsChanged,

		AppHash: s.AppHash,

		LastResultsHash: s.LastResultsHash,
	}
}

// Equals returns true if the States are identical.
func (s State) Equals(s2 State) bool {
	return bytes.Equal(s.Bytes(), s2.Bytes())
}

// Bytes serializes the State using go-wire.
func (s State) Bytes() []byte {
	sbytes, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Error reading GenesisDoc: %v", err)
		return nil
	}
	return sbytes
}

// IsEmpty returns true if the State is equal to the empty State.
func (s State) IsEmpty() bool {
	return s.Validators == nil // XXX can't compare to Empty
}

// GetValidators returns the last and current validator sets.
func (s State) GetValidators() (last *ttypes.ValidatorSet, current *ttypes.ValidatorSet) {
	return s.LastValidators, s.Validators
}

//------------------------------------------------------------------------
// Create a block from the latest state

// MakeBlock builds a block with the given txs and commit from the current state.
func (s State) MakeBlock(height int64, round int64, Txs []*types.Transaction, commit *types.TendermintCommit) *ttypes.TendermintBlock {
	// build base block
	block := ttypes.MakeBlock(height, round, Txs, commit)

	// fill header with state data
	block.Header.ChainID = s.ChainID
	block.Header.TotalTxs = s.LastBlockTotalTx + block.Header.NumTxs
	block.Header.LastBlockID = &s.LastBlockID.BlockID
	block.Header.ValidatorsHash = s.Validators.Hash()
	block.Header.AppHash = s.AppHash
	block.Header.ConsensusHash = s.ConsensusParams.Hash()
	block.Header.LastResultsHash = s.LastResultsHash

	return block
}

//------------------------------------------------------------------------
// Genesis

// MakeGenesisStateFromFile reads and unmarshals state from the given
// file.
//
// Used during replay and in tests.
func MakeGenesisStateFromFile(genDocFile string) (State, error) {
	genDoc, err := MakeGenesisDocFromFile(genDocFile)
	if err != nil {
		return State{}, err
	}
	return MakeGenesisState(genDoc)
}

// MakeGenesisDocFromFile reads and unmarshals genesis doc from the given file.
func MakeGenesisDocFromFile(genDocFile string) (*ttypes.GenesisDoc, error) {
	genDocJSON, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read GenesisDoc file: %v", err)
	}
	genDoc, err := ttypes.GenesisDocFromJSON(genDocJSON)
	if err != nil {
		return nil, fmt.Errorf("Error reading GenesisDoc: %v", err)
	}
	return genDoc, nil
}

// MakeGenesisState creates state from ttypes.GenesisDoc.
func MakeGenesisState(genDoc *ttypes.GenesisDoc) (State, error) {
	err := genDoc.ValidateAndComplete()
	if err != nil {
		return State{}, fmt.Errorf("Error in genesis file: %v", err)
	}

	// Make validators slice
	validators := make([]*ttypes.Validator, len(genDoc.Validators))
	for i, val := range genDoc.Validators {
		pubKey, err := ttypes.PubKeyFromString(val.PubKey.Data)
		if err != nil {
			return State{}, fmt.Errorf("Error validate[i] in genesis file: %v", i, err)
		}

		// Make validator
		validators[i] = &ttypes.Validator{
			Address:     ttypes.GenAddressByPubKey(pubKey),
			PubKey:      pubKey.Bytes(),
			VotingPower: val.Power,
		}
	}

	return State{

		ChainID: genDoc.ChainID,

		LastBlockHeight: 0,
		LastBlockID:     ttypes.BlockID{},
		LastBlockTime:   genDoc.GenesisTime.UnixNano(),

		Validators:                  ttypes.NewValidatorSet(validators),
		LastValidators:              ttypes.NewValidatorSet(nil),
		LastHeightValidatorsChanged: 1,

		ConsensusParams:                  *genDoc.ConsensusParams,
		LastHeightConsensusParamsChanged: 1,

		AppHash: genDoc.AppHash,
	}, nil
}
