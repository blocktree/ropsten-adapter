/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwtester

import (
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/blocktree/quorum-adapter/quorum"
	"testing"
)

func TestCallSmartContractABI(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WJyrmJL67qb6LsxtkaEkBrZoVtg83Ecj1y"
	accountID := "Ff7w6KZZ9UWuVgVwVkvA9aegGxVTdKNr5jC2LDv142gs"

	contract := openwallet.SmartContract{
		Address: "0x471483d78a0940bcfd38efb4a8cc017e285bd417",
		Symbol:  "TESTETH",
	}
	contract.SetABI(quorum.ERC20_ABI_JSON)

	callParam := []string{
		"approve",
		"0x19a4b5d6ea319a5d5ad1d4cc00a5e2e28cac5ec3",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}

	result, err := tm.CallSmartContractABI(testApp, walletID, accountID, &contract, callParam)
	if err != nil {
		t.Errorf("CallSmartContractABI failed, unexpected error: %v", err)
		return
	}
	log.Infof("result: %+v", result)
	//0x19a4b5d6ea319a5d5ad1d4cc00a5e2e28cac5ec3
}

func TestSmartContractTransaction(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WJyrmJL67qb6LsxtkaEkBrZoVtg83Ecj1y"
	accountID := "Ff7w6KZZ9UWuVgVwVkvA9aegGxVTdKNr5jC2LDv142gs"

	contract := openwallet.SmartContract{
		Address: "0x471483d78a0940bcfd38efb4a8cc017e285bd417",
		Symbol:  "TESTETH",
	}
	contract.SetABI(quorum.ERC20_ABI_JSON)
	callParam := []string{
		"approve",
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}

	rawTx, err := tm.CreateSmartContractTransaction(testApp, walletID, accountID, "", "", &contract, callParam)
	if err != nil {
		t.Errorf("CreateSmartContractTransaction failed, unexpected error: %v", err)
		return
	}
	//log.Infof("rawTx: %+v", rawTx)

	_, err = tm.SignSmartContractTransaction(testApp, walletID, accountID, "12345678", rawTx)
	if err != nil {
		t.Errorf("SignSmartContractTransaction failed, unexpected error: %v", err)
		return
	}
	rawTx.AwaitResult = true
	rawTx.AwaitTimeout = 10
	tx, err := tm.SubmitSmartContractTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		t.Errorf("SubmitSmartContractTransaction failed, unexpected error: %v", err)
		return
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	for i, event := range tx.Events {
		log.Std.Notice("data.Events[%d]: %+v", i, event)
	}
}
