package main

import (
	"database/sql"
	"sync"

	"github.com/google/uuid"
)

type mockStore struct {
	mu      sync.Mutex
	wallets map[string]int
	odds    map[string]map[string]*Odds // owner to odds name to odds object
	oddsIds map[string]*Odds
}

func newMockStore() (*mockStore, error) {
	return &mockStore{
		wallets: make(map[string]int),
		odds:    make(map[string]map[string]*Odds),
		oddsIds: make(map[string]*Odds),
	}, nil
}

func (s *mockStore) getWallet(id string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, e := s.wallets[id]
	if !e {
		return 0, sql.ErrNoRows
	}
	return v, nil
}

func (s *mockStore) setWallet(id string, amt int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wallets[id] = amt
	return nil
}

func (s *mockStore) setWalletDelta(id string, amt int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, e := s.wallets[id]
	if !e {
		return sql.ErrNoRows
	}
	s.wallets[id] += amt
	return nil
}

func (s *mockStore) getOdds(owner, name string) (*Odds, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, e := s.odds[owner][name]
	if !e {
		return nil, sql.ErrNoRows
	}
	return v, nil
}

func (s *mockStore) getOddsFromId(id string) (*Odds, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, e := s.oddsIds[id]
	if !e {
		return nil, sql.ErrNoRows
	}
	return v, nil
}

func (s *mockStore) setOdds(owner string, odds Odds) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, e := s.odds[owner]; !e {
		s.odds[owner] = make(map[string]*Odds)
	}

	s.odds[owner][odds.name] = &odds
	if s.odds[owner][odds.name].options == nil {
		s.odds[owner][odds.name].options = make(map[string]*OddsOption)
	}

	id := uuid.New().String()
	if odds.id != "" {
		id = odds.id
	}
	s.oddsIds[id] = s.odds[owner][odds.name]
	s.oddsIds[id].id = id

	s.oddsIds[id].owner = owner

	return nil
}

func (s *mockStore) setOddsOpt(owner string, gamename string, oddsOpt OddsOption) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, e := s.odds[owner][gamename]; !e {
		return sql.ErrNoRows
	}

	s.odds[owner][gamename].options[oddsOpt.name] = &oddsOpt

	return nil
}

func (s *mockStore) delOdds(owner string, gamename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, e := s.odds[owner][gamename]; !e {
		return sql.ErrNoRows
	}
	delete(s.oddsIds, s.odds[owner][gamename].id)
	delete(s.odds[owner], gamename)
	return nil
}

func (s *mockStore) delOddsOpt(owner, gamename, optionname string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, e := s.odds[owner][gamename]; !e {
		return sql.ErrNoRows
	}
	delete(s.odds[owner][gamename].options, optionname)
	return nil
}

func (s *mockStore) getOddsOpt(owner string, gamename string, optname string) (*OddsOption, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, e := s.odds[owner][gamename]; !e {
		return nil, sql.ErrNoRows
	}
	o, e := s.odds[owner][gamename].options[optname]
	if !e {
		return nil, sql.ErrNoRows
	}
	return o, nil
}

func (s *mockStore) setBet(better string, gameid string, optname string, amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if odds, e := s.oddsIds[gameid]; e {
		if odds.bets == nil {
			odds.bets = make(map[string]*OddsBet)
		}
		odds.bets[better] = &OddsBet{
			betterid: better,
			amt:      amount,
			option:   odds.options[optname],
		}
		return nil
	}
	return sql.ErrNoRows
}
