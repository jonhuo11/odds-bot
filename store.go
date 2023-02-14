package main

type store interface {
	cleanup() // dumps data to disk

	getWallet(ownerId string) (Wallet, error)
	setWallet(w Wallet) error
	updateWalletDelta(ownerId string, amt int) error

	getOdds(ownerId string, gameName string) (OddsModel, error)
	getOddsFromId(gameId string) (OddsModel, error)
	setOdds(odds OddsModel) error
	setOddsStarted(gameId string, started bool) error

	setOddsOpt(opt OddsOptionModel) error
	getOddsOptFromId(optionId string) (OddsOptionModel, error)
	getOddsOptFromGameIdAndName(gameId, optName string) (OddsOptionModel, error)
	// list of all options belonging to a game
	getOddsOptsForGame(gameId string) ([]OddsOptionModel, error)

	setBet(bet OddsBetModel) error
	getBetsForUser(id string) ([]OddsBetModel, error)
	getBetsForGame(gameId string) ([]OddsBetModel, error)

	calculateWinnings(betId string) (int, error)
}

// models exactly represent rows in sql
type Wallet struct {
	ownerId string
	balance int
}

type OddsModel struct {
	id      string
	name    string
	ownerId string
	started bool
	winner  string
}

type OddsOptionModel struct {
	id        string
	gameId    string
	name      string
	moneyline int
}

type OddsBetModel struct {
	id       string
	ownerId  string
	optionId string
	amount   int
}
