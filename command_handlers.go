package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// moneyline is converted to percentage
// moneyline = 100 ==> reward = bet
// moneyline = -100 ==> reward = bet
// moneyline = -150 ==> reward = bet * (-1/moneyline/100)
// moneyline be
func calculateWinnings(bet, moneyline int) int {
	ml := float64(moneyline)
	b := float64(bet)
	if ml <= -100 {
		m := (ml / 100) * -1
		return int(b * (1 / m))
	}
	if ml >= 100 {
		m := ml / 100
		return int(b * m)
	}
	return 0 // you cannot have a moneyline between -100 and 100
}

// odds new name
// odds add game choice moneyline
// odds info game
// odds start game
// odds del name
// odds delchoice game choice
func oddsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	res := ""
	switch options[0].Name {
	case "new":
		options = options[0].Options
		om := makeOptsMap(options)
		newOdds := OddsModel{
			id:      uuid.NewString(),
			ownerId: i.Member.User.ID,
			started: false,
		}
		if v, ok := om["name"]; ok {
			newOdds.name = v.StringValue()
		} else {
			res = "An error occurred"
			break
		}
		// check for existing game with this name, do not allow duplicates
		_, err := storeDao.getOdds(i.Member.User.ID, newOdds.name)
		if err != nil {
			if err == sql.ErrNoRows {
				if err = storeDao.setOdds(newOdds); err != nil {
					res = "An error occurred"
					break
				}
				res = fmt.Sprintf("Created new odds game called %v", newOdds.name)
				break
			}
			res = "An error occurred"
			break
		}
		res = "There is already an odds game with this name"

	case "add":
		options = options[0].Options
		om := makeOptsMap(options)
		newOddsOpt := OddsOptionModel{
			id: uuid.NewString(),
		}
		if gameName, ok := om["game"]; ok {
			game, err := storeDao.getOdds(i.Member.User.ID, gameName.StringValue())
			if err != nil {
				if err == sql.ErrNoRows {
					res = "No odds game with this name was found"
					break
				}
				res = "An error occurred"
				break
			}
			newOddsOpt.gameId = game.id

			// TODO check for existing choice with this name, do not allow duplicates

			if choice, ok := om["choice"]; ok {
				if moneyline, ok := om["moneyline"]; ok {
					// moneyline must not be between -100 and 100
					moneylineInt := moneyline.IntValue()
					if moneylineInt > -100 && moneylineInt < 100 {
						res = "Moneyline cannot be between -100 and 100"
						break
					}

					newOddsOpt.name = choice.StringValue()
					newOddsOpt.moneyline = int(moneylineInt)
					storeDao.setOddsOpt(newOddsOpt)

					res = fmt.Sprintf("Added option %v (%v) to odds game %v", newOddsOpt.name, newOddsOpt.moneyline, gameName.StringValue())
					break
				}
			}
		}
		res = "An error occurred"

	/*
		case "del":
			// odds del name
			options = options[0].Options
			om := makeOptsMap(options)
			if gamename, ok := om["name"]; ok {
				if err := storeDao.delOdds(i.Member.User.ID, gamename.StringValue()); err != nil {
					res = err.Error()
					break
				}
				res = "Deleted odds game " + gamename.StringValue()
				break
			}
			res = "An error occurred"
	*/
	/*
		case "delchoice":
			// odds delchoice game choice
			options = options[0].Options
			om := makeOptsMap(options)

			if gamename, ok := om["game"]; ok {
				if choice, ok := om["choice"]; ok {
					if err := storeDao.delOddsOpt(i.Member.User.ID, gamename.StringValue(), choice.StringValue()); err != nil {
						if err == sql.ErrNoRows {
							res = "No game with this name found"
						}
						res = err.Error()
						break
					}
					res = "Deleted option " + choice.StringValue() + " from odds game " + gamename.StringValue()
					break
				}
			}

			res = "An error occurred"
	*/
	case "start":
		// allows users to start betting on this odds game
		res = "An error occurred"

	case "info":
		options = options[0].Options
		om := makeOptsMap(options)
		if game, ok := om["game"]; ok {

			// get odds, options, and bets
			oddsGame, err := storeDao.getOdds(i.Member.User.ID, game.StringValue())
			if err != nil { // try to find this game
				if err == sql.ErrNoRows {
					res = "No odds game with this name under your account was found"
					break
				}
				res = "An error occurred"
				break
			}
			gameOptions, err1 := storeDao.getOddsOptsForGame(oddsGame.id)
			gameBetters, err2 := storeDao.getBetsForGame(oddsGame.id)
			if err1 != nil || err2 != nil {
				res = "An error occurred"
				break
			}

			opts := ""
			for i := range gameOptions {
				opt := gameOptions[i]
				opts += "\n- " + opt.name + fmt.Sprintf(" (%v)", opt.moneyline)
			}
			betters := ""
			for i := range gameBetters {
				b := gameBetters[i]
				betters += "\n- " + "<@" + b.ownerId + ">"
			}
			var winner = "Undecided"
			if oddsGame.winner != "" {
				winner = oddsGame.winner
			}
			res = fmt.Sprintf(
				"**Game name:** %v\n**Game ID:** %v\n**Created by:** <@%v>\n**Options:** %v\n**Betters:** %v\n**Winner:** %v\n",
				oddsGame.name, oddsGame.id, i.Member.User.ID, opts, betters, winner,
			)
			break
		}
		res = "An error occurred"

	default:
		res = "An error occurred"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res,
		},
	})
}

func walletHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	res := ""
	switch options[0].Name {
	case "get":
		amt, err := storeDao.getWallet(i.Member.User.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err = storeDao.setWallet(Wallet{
					ownerId: i.Member.User.ID,
					balance: 1000,
				}); err != nil {
					res = err.Error()
					break
				}
				// new users get 1000 coins
				// TODO load this number from settings.json?
				res = fmt.Sprintf("New user detected, have %v coins on the house", 1000)
			} else {
				res = err.Error()
			}
		} else {
			res = fmt.Sprintf("You have %v coins", amt)
		}
	default:
		res = "An error occurred"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res,
		},
	})
}

// bet make id option amount
// bet del id option
// bet list
func betHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	res := ""
	switch options[0].Name {
	case "list":
		// list all active bets belonging to yourself
		list, err := storeDao.getBetsForUser(i.Member.User.ID)
		if err != nil {
			res = err.Error()
			break
		}
		for i := range list {
			betItem := list[i]
			// get the option and game that was bet on
			betOption, err1 := storeDao.getOddsOptFromId(betItem.optionId)
			betGame, err2 := storeDao.getOddsFromId(betOption.gameId)
			winnings, err3 := storeDao.calculateWinnings(betItem.id)
			if err1 != nil || err2 != nil || err3 != nil {
				res = "An error occurred"
				break
			}
			res += fmt.Sprintf(
				"\n- %v coin bet on the outcome '%v' in game '%v', with a potential payout of %v",
				betItem.amount, betOption.name, betGame.name, winnings,
			)
		}
		if res == "" {
			res = "No active bets"
		}
	case "make":
		om := makeOptsMap(options[0].Options)
		if id, ok := om["id"]; ok {
			if opt := om["option"]; ok {
				if amt := om["amount"]; ok {
					// bet amount must be positive
					if amt.IntValue() < 0 {
						res = "You cannot place a negative bet"
						break
					}
					// check if option exists, and if user has enough coins
					var w Wallet
					var err error
					var newUserMsg string
					if w, err = storeDao.getWallet(i.Member.User.ID); err != nil {
						if err == sql.ErrNoRows {
							newUserMsg = "New user detected, granting 1000 coins on the house. "
							// grant 1000 coins on house
							if err := storeDao.setWallet(Wallet{
								ownerId: i.Member.User.ID,
								balance: 1000,
							}); err != nil {
								res = err.Error()
								break
							}
							w.balance = 1000
						} else {
							res = err.Error()
							break
						}
					}
					if w.balance < int(amt.IntValue()) {
						res = "You do not have enough coins to place this bet"
						break
					}

					// get the odds from id
					var moneyline int
					var optid string
					if odds, err := storeDao.getOddsFromId(id.StringValue()); err != nil {
						if err == sql.ErrNoRows {
							res = "No odds game with this id exists"
							break
						}
						res = err.Error()
						break
					} else {
						if oddsopt, err := storeDao.getOddsOptFromGameIdAndName(odds.id, opt.StringValue()); err != nil {
							if err == sql.ErrNoRows {
								res = "No option with this name exists under this game id"
								break
							}
							res = err.Error()
							break
						} else {
							moneyline = oddsopt.moneyline
							optid = oddsopt.id
						}
					}
					// set the bet
					if err := storeDao.setBet(OddsBetModel{
						id:       uuid.NewString(),
						ownerId:  i.Member.User.ID,
						optionId: optid,
						amount:   int(amt.IntValue()),
					}); err != nil {
						res = "An error occurred when setting the bet: " + err.Error()
						break
					}
					// subtract from wallet
					if err := storeDao.updateWalletDelta(i.Member.User.ID, int(-1*amt.IntValue())); err != nil {
						res = err.Error()
						break
					}
					res = fmt.Sprintf(
						"%vMade a %v coin bet on '%v', your new coin balance is %v. If you win this bet, you will gain %v coins!",
						newUserMsg, amt.IntValue(), opt.StringValue(), int(w.balance-int(amt.IntValue())), calculateWinnings(int(amt.IntValue()), moneyline),
					)
					break
				}
			}
		}
		res = "An error occurred"
	default:
		res = "An error occurred"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res,
		},
	})
}

func pingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
}
