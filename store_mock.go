/*
TODO eventually make a store_actual.go which uses postgres or sqlite
*/

package main

import (
	"sync"
)

type mockStore struct {
	mu sync.Mutex

	// maps represent tables in sql
	wallets     map[string]*Wallet
	odds        map[string]*OddsModel
	oddsOptions map[string]*OddsOptionModel
	oddsBets    map[string]*OddsBetModel
}

func newMockStore() (*mockStore, error) {
	s := &mockStore{
		wallets:     make(map[string]*Wallet),
		odds:        make(map[string]*OddsModel),
		oddsOptions: make(map[string]*OddsOptionModel),
		oddsBets:    make(map[string]*OddsBetModel),
	}

	return s, nil
}

func (s *mockStore) getWallet(ownerId string) (Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) setWallet(w Wallet) error {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) updateWalletDelta(ownerId string, amt int) error {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getOdds(ownerId string, gameName string) (OddsModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getOddsFromId(gameId string) (OddsModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) setOdds(odds OddsModel) error {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) setOddsOpt(opt OddsOptionModel) error {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getOddsOptFromId(optionId string) (OddsOptionModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getOddsOptFromGameIdAndName(gameId string, optName string) (OddsOptionModel, error) {
	panic("not implemented") // TODO: Implement
}

// list of all options belonging to a game
func (s *mockStore) getOddsOptsForGame(gameId string) ([]OddsOptionModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) setBet(bet OddsBetModel) error {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getBetsForUser(id string) ([]OddsBetModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) getBetsForGame(gameId string) ([]OddsBetModel, error) {
	panic("not implemented") // TODO: Implement
}

func (s *mockStore) calculateWinnings(betId string) (int, error) {
	panic("not implemented") // TODO: Implement
}
