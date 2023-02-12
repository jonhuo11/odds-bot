package main

type store interface {
	getWallet(id string) (int, error)
	setWallet(id string, amt int) error
	setWalletDelta(id string, amt int) error

	getOdds(owner string, name string) (*Odds, error)
	setOdds(owner string, odds Odds) error
	delOdds(owner string, gamename string) error
	setOddsOpt(owner string, gamename string, opt OddsOption) error
	getOddsOpt(owner string, gamename string, optname string) (*OddsOption, error)
	delOddsOpt(owner, gamename, optionname string) error
}

type Odds struct {
	name    string
	id      string
	options map[string]*OddsOption
	bets    map[string]*OddsBet // discord id --> bet --> option
	winner  string
}

type OddsOption struct {
	name      string
	moneyline int
}

type OddsBet struct {
	betterid string
	amt      int
	option   *OddsOption
}
