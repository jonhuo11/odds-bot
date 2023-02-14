/*
TODO eventually make a store_actual.go which uses postgres or sqlite
*/

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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

var (
	walletDataFileName string = "wallet_data.txt"
)

func newMockStore() (*mockStore, error) {
	s := &mockStore{
		wallets:     make(map[string]*Wallet),
		odds:        make(map[string]*OddsModel),
		oddsOptions: make(map[string]*OddsOptionModel),
		oddsBets:    make(map[string]*OddsBetModel),
	}

	// load wallet data
	f, err := os.OpenFile(walletDataFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id := scanner.Text()
		if id == "" {
			continue
		}
		scanner.Scan()
		if bal, err := strconv.Atoi(scanner.Text()); err == nil {
			s.wallets[id] = &Wallet{
				ownerId: id,
				balance: bal,
			}
		} else {
			return nil, err
		}
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	out := "\nLoaded wallet balances:\n"
	for _, w := range s.wallets {
		out += fmt.Sprintf("- user %v, balance %v\n", w.ownerId, w.balance)
	}
	log.Println(out)

	// make a test game with odds to bet on
	s.setOdds(OddsModel{
		id:      "0",
		name:    "test",
		ownerId: "0",
		started: true,
		winner:  "",
	})
	s.setOddsOpt(OddsOptionModel{
		id:        "opt0",
		gameId:    "0",
		name:      "testopt0",
		moneyline: 150,
	})
	s.setOddsOpt(OddsOptionModel{
		id:        "opt1",
		gameId:    "0",
		name:      "testopt1",
		moneyline: -200,
	})

	return s, nil
}

func (s *mockStore) cleanup() {
	// write wallet data to walletdata.txt
	f, err := os.OpenFile(walletDataFileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	for _, w := range s.wallets {
		f.WriteString(fmt.Sprintf("%v\n%v\n", w.ownerId, w.balance))
	}
}

func (s *mockStore) getWallet(ownerId string) (Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, e := s.wallets[ownerId]
	if !e {
		return Wallet{}, sql.ErrNoRows
	}
	return *w, nil
}

func (s *mockStore) setWallet(w Wallet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wallets[w.ownerId] = &w
	return nil
}

func (s *mockStore) updateWalletDelta(ownerId string, amt int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, e := s.wallets[ownerId]
	if !e {
		return sql.ErrNoRows
	}
	s.wallets[ownerId].balance += amt
	return nil
}

func (s *mockStore) getOdds(ownerId string, gameName string) (OddsModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, o := range s.odds {
		if o.name == gameName && o.ownerId == ownerId {
			return *o, nil
		}
	}
	return OddsModel{}, sql.ErrNoRows
}

func (s *mockStore) getOddsFromId(gameId string) (OddsModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, e := s.odds[gameId]
	if !e {
		return OddsModel{}, sql.ErrNoRows
	}
	return *o, nil
}

func (s *mockStore) setOdds(odds OddsModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.odds[odds.id] = &odds
	return nil
}

func (s *mockStore) setOddsOpt(opt OddsOptionModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.oddsOptions[opt.id] = &opt
	return nil
}

func (s *mockStore) getOddsOptFromId(optionId string) (OddsOptionModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, e := s.oddsOptions[optionId]
	if !e {
		return OddsOptionModel{}, sql.ErrNoRows
	}
	return *o, nil
}

func (s *mockStore) getOddsOptFromGameIdAndName(gameId string, optName string) (OddsOptionModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, o := range s.oddsOptions {
		if o.gameId == gameId && o.name == optName {
			return *o, nil
		}
	}
	return OddsOptionModel{}, sql.ErrNoRows
}

// list of all options belonging to a game
func (s *mockStore) getOddsOptsForGame(gameId string) ([]OddsOptionModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]OddsOptionModel, 0)
	for _, o := range s.oddsOptions {
		if o.gameId == gameId {
			out = append(out, *o)
		}
	}
	return out, nil
}

func (s *mockStore) setBet(bet OddsBetModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.oddsBets[bet.id] = &bet
	return nil
}

func (s *mockStore) getBetsForUser(id string) ([]OddsBetModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]OddsBetModel, 0)
	for _, b := range s.oddsBets {
		if b.ownerId == id {
			out = append(out, *b)
		}
	}
	return out, nil
}

func (s *mockStore) getBetsForGame(gameId string) ([]OddsBetModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]OddsBetModel, 0)
	for _, b := range s.oddsBets {
		o, e := s.oddsOptions[b.optionId]
		if e {
			if o.gameId == gameId {
				out = append(out, *b)
			}
		} else {
			return nil, sql.ErrNoRows
		}
	}
	return out, nil
}

func (s *mockStore) calculateWinnings(betId string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, e := s.oddsBets[betId]
	if e {
		if o, e2 := s.oddsOptions[b.optionId]; e2 {
			return calculateWinnings(b.amount, o.moneyline), nil
		}
	}
	return 0, sql.ErrNoRows
}
