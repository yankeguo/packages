// Copyright (c) The Cryptorium Authors.
// Licensed under the MIT License.

package hardhat

import (
	"context"
	"math/big"
	"os"
	"path"
	"strings"

	contraget "github.com/cryptoriums/contraget/pkg/cli"
	"github.com/cryptoriums/packages/ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

const DefaultUrl = "ws://127.0.0.1:8545"

func ReplaceContract(ctx context.Context, nodeURL string, contractPath string, contractName string, contractAddrToReplace common.Address) error {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	cfg := &contraget.Cli{
		Path:        contractPath,
		Name:        "contract",
		ObjectsDst:  "tmp",
		DownloadDst: "tmp",
	}

	if err := contraget.Run(cfg); err != nil {
		return errors.Wrap(err, "generating the contract bin file")
	}

	_bin, err := os.ReadFile(path.Join(cfg.ObjectsDst, contractName+".bin"))
	if err != nil {
		return errors.Wrap(err, "reading the bin file")
	}
	bin := string(_bin)
	// For solidity contracts
	startDeployBin := strings.LastIndex(bin, "60806040523480156")
	// For vyper contracts.
	if startDeployBin == -1 {
		startDeployBin = strings.LastIndex(bin, "600436")
	}

	if startDeployBin == -1 {
		return errors.New("start index of runtime Bytecode not found in the generated binary file")
	}

	err = rpcClient.CallContext(ctx, nil, "hardhat_setCode", contractAddrToReplace, "0x"+bin[startDeployBin:])
	if err != nil {
		return errors.Wrap(err, "hardhat_setCode call")
	}

	return nil
}

func Mine(ctx context.Context, nodeURL string) error {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	err = rpcClient.CallContext(ctx, nil, "evm_mine")
	if err != nil {
		return errors.Wrap(err, "calling evm_mine")
	}

	return nil
}

func SetNextBlockTimestamp(ctx context.Context, nodeURL string, ts int64) error {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	err = rpcClient.CallContext(ctx, nil, "evm_setNextBlockTimestamp", big.NewInt(ts))
	if err != nil {
		return errors.Wrap(err, "calling evm_setNextBlockTimestamp")
	}

	return nil
}

func TxWithImpersonateAccount(ctx context.Context, nodeURL string, from common.Address, to common.Address, abiJ string, funcName string, args []interface{}) (string, error) {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return "", errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	err = rpcClient.CallContext(ctx, nil, "hardhat_impersonateAccount", from)
	if err != nil {
		return "", errors.Wrap(err, "calling hardhat_impersonateAccount")
	}

	abiParsed, err := abi.JSON(strings.NewReader(abiJ))
	if err != nil {
		return "", errors.Wrap(err, "parsing the abi")
	}
	data, err := abiParsed.Pack(funcName, args...)
	if err != nil {
		return "", errors.Wrap(err, "packing the args")
	}
	optsT := ethereum.SendTransactionOpts{
		From: from,
		To:   to,
		Data: hexutil.Encode(data),
	}
	var txHash string
	err = rpcClient.CallContext(ctx, &txHash, "eth_sendTransaction", optsT)
	if err != nil {
		return "", errors.Wrap(err, "calling eth_sendTransaction")
	}

	return txHash, nil
}

func SetBalance(ctx context.Context, nodeURL string, of common.Address, amnt *big.Int) error {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	err = rpcClient.CallContext(ctx, nil, "hardhat_setBalance", of, hexutil.EncodeBig(amnt))
	if err != nil {
		return errors.Wrap(err, "calling hardhat_setBalance")
	}

	return nil
}

func SetStorageAt(ctx context.Context, nodeURL string, addr common.Address, idx string, val string) error {
	rpcClient, err := rpc.DialContext(ctx, nodeURL)
	if err != nil {
		return errors.Wrap(err, "creating rpc client")
	}
	defer rpcClient.Close()

	err = rpcClient.CallContext(ctx, nil, "hardhat_setStorageAt", addr.Hex(), idx, val)
	if err != nil {
		return errors.Wrap(err, "calling hardhat_setStorageAt")
	}

	return nil
}
