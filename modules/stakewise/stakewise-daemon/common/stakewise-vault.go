package swcommon

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/eth"
	batch "github.com/rocket-pool/batch-query"
)

const (
	stakewiseVaultAbiString string = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"assets","type":"uint256"}],"name":"CheckpointCreated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"address","name":"feeRecipient","type":"address"}],"name":"FeeRecipientUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"receiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"assets","type":"uint256"}],"name":"FeeSharesMinted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"address","name":"keysManager","type":"address"}],"name":"KeysManagerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":false,"internalType":"string","name":"metadataIpfsHash","type":"string"}],"name":"MetadataUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"}],"name":"ValidatorRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"bytes32","name":"validatorsRoot","type":"bytes32"}],"name":"ValidatorsRootUpdated","type":"event"},{"inputs":[],"name":"admin","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"capacity","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"shares","type":"uint256"}],"name":"convertToAssets","outputs":[{"internalType":"uint256","name":"assets","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"assets","type":"uint256"}],"name":"convertToShares","outputs":[{"internalType":"uint256","name":"shares","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feePercent","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeRecipient","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"getShares","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"isStateUpdateRequired","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"keysManager","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"queuedShares","outputs":[{"internalType":"uint128","name":"","type":"uint128"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"validatorsRegistryRoot","type":"bytes32"},{"internalType":"uint256","name":"deadline","type":"uint256"},{"internalType":"bytes","name":"validators","type":"bytes"},{"internalType":"bytes","name":"signatures","type":"bytes"},{"internalType":"string","name":"exitSignaturesIpfsHash","type":"string"}],"internalType":"struct IKeeperValidators.ApprovalParams","name":"keeperParams","type":"tuple"},{"internalType":"bytes32[]","name":"proof","type":"bytes32[]"}],"name":"registerValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"validatorsRegistryRoot","type":"bytes32"},{"internalType":"uint256","name":"deadline","type":"uint256"},{"internalType":"bytes","name":"validators","type":"bytes"},{"internalType":"bytes","name":"signatures","type":"bytes"},{"internalType":"string","name":"exitSignaturesIpfsHash","type":"string"}],"internalType":"struct IKeeperValidators.ApprovalParams","name":"keeperParams","type":"tuple"},{"internalType":"uint256[]","name":"indexes","type":"uint256[]"},{"internalType":"bool[]","name":"proofFlags","type":"bool[]"},{"internalType":"bytes32[]","name":"proof","type":"bytes32[]"}],"name":"registerValidators","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_feeRecipient","type":"address"}],"name":"setFeeRecipient","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_keysManager","type":"address"}],"name":"setKeysManager","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"metadataIpfsHash","type":"string"}],"name":"setMetadata","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"_validatorsRoot","type":"bytes32"}],"name":"setValidatorsRoot","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"totalAssets","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalShares","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"bytes32","name":"rewardsRoot","type":"bytes32"},{"internalType":"int160","name":"reward","type":"int160"},{"internalType":"uint160","name":"unlockedMevReward","type":"uint160"},{"internalType":"bytes32[]","name":"proof","type":"bytes32[]"}],"internalType":"struct IKeeperRewards.HarvestParams","name":"harvestParams","type":"tuple"}],"name":"updateState","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"validatorIndex","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"validatorsRoot","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"withdrawableAssets","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
)

// ABI cache
var stakewiseVaultAbi abi.ABI
var stakewiseVaultOnce sync.Once

// Binding for Stakewise vaults
type StakewiseVault struct {
	contract *eth.Contract
	txMgr    *eth.TransactionManager
}

// Create a new Stakewise vault instance
func NewStakewiseVault(address common.Address, ec eth.IExecutionClient, txMgr *eth.TransactionManager) (*StakewiseVault, error) {
	// Parse the ABI
	var err error
	stakewiseVaultOnce.Do(func() {
		var parsedAbi abi.ABI
		parsedAbi, err = abi.JSON(strings.NewReader(stakewiseVaultAbiString))
		if err == nil {
			stakewiseVaultAbi = parsedAbi
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing Stakewise Vault ABI: %w", err)
	}

	// Create the contract
	contract := &eth.Contract{
		ContractImpl: bind.NewBoundContract(address, stakewiseVaultAbi, ec, ec, ec),
		Address:      address,
		ABI:          &stakewiseVaultAbi,
	}

	return &StakewiseVault{
		contract: contract,
		txMgr:    txMgr,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the current validators root in the contracts
func (c *StakewiseVault) GetValidatorsRoot(mc *batch.MultiCaller, out *common.Hash) {
	eth.AddCallToMulicaller(mc, c.contract, out, "validatorsRoot")
}

// ====================
// === Transactions ===
// ====================

// Set the validator deposit data root for the vault
func (c *StakewiseVault) SetDepositDataRoot(dataRoot common.Hash, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.contract, "setValidatorsRoot", opts, dataRoot)
}
