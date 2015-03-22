package main

import (
	"fmt"

	. "github.com/tendermint/tendermint/common"
	. "github.com/tendermint/tendermint/vm"
	"github.com/tendermint/tendermint/vm/sha3"
)

type FakeAppState struct {
	accounts map[string]*Account
	storage  map[string]Word
	logs     []*Log
}

func (fas *FakeAppState) GetAccount(addr Word) (*Account, error) {
	account := fas.accounts[addr.String()]
	if account != nil {
		return account, nil
	} else {
		return nil, Errorf("Invalid account addr: %v", addr)
	}
}

func (fas *FakeAppState) UpdateAccount(account *Account) error {
	_, ok := fas.accounts[account.Address.String()]
	if !ok {
		return Errorf("Invalid account addr: %v", account.Address.String())
	} else {
		// Nothing to do
		return nil
	}
}

func (fas *FakeAppState) DeleteAccount(account *Account) error {
	_, ok := fas.accounts[account.Address.String()]
	if !ok {
		return Errorf("Invalid account addr: %v", account.Address.String())
	} else {
		// Delete account
		delete(fas.accounts, account.Address.String())
		return nil
	}
}

func (fas *FakeAppState) CreateAccount(creator *Account) (*Account, error) {
	addr := createAddress(creator)
	account := fas.accounts[addr.String()]
	if account == nil {
		return &Account{
			Address:     addr,
			Balance:     0,
			Code:        nil,
			Nonce:       0,
			StorageRoot: Zero,
		}, nil
	} else {
		return nil, Errorf("Invalid account addr: %v", addr)
	}
}

func (fas *FakeAppState) GetStorage(addr Word, key Word) (Word, error) {
	_, ok := fas.accounts[addr.String()]
	if !ok {
		return Zero, Errorf("Invalid account addr: %v", addr)
	}

	value, ok := fas.storage[addr.String()+key.String()]
	if ok {
		return value, nil
	} else {
		return Zero, nil
	}
}

func (fas *FakeAppState) SetStorage(addr Word, key Word, value Word) (bool, error) {
	_, ok := fas.accounts[addr.String()]
	if !ok {
		return false, Errorf("Invalid account addr: %v", addr)
	}

	_, ok = fas.storage[addr.String()+key.String()]
	fas.storage[addr.String()+key.String()] = value
	return ok, nil
}

func (fas *FakeAppState) AddLog(log *Log) {
	fas.logs = append(fas.logs, log)
}

func main() {
	appState := &FakeAppState{
		accounts: make(map[string]*Account),
		storage:  make(map[string]Word),
		logs:     nil,
	}
	params := Params{
		BlockHeight: 0,
		BlockHash:   Zero,
		BlockTime:   0,
		GasLimit:    0,
	}
	ourVm := NewVM(appState, params, Zero)

	// Create accounts
	account1 := &Account{
		Address: Uint64ToWord(100),
	}
	account2 := &Account{
		Address: Uint64ToWord(101),
	}

	var gas uint64 = 1000
	output, err := ourVm.Call(account1, account2, []byte{0x5B, 0x60, 0x00, 0x56}, []byte{}, 0, &gas)
	fmt.Printf("Output: %v Error: %v\n", output, err)
}

// Creates a 20 byte address and bumps the nonce.
func createAddress(creator *Account) Word {
	nonce := creator.Nonce
	creator.Nonce += 1
	temp := make([]byte, 32+8)
	copy(temp, creator.Address[:])
	PutUint64(temp[32:], nonce)
	return RightPadWord(sha3.Sha3(temp)[:20])
}