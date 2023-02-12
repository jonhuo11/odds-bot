package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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
		newOdds := Odds{}
		if v, ok := om["name"]; ok {
			newOdds.name = v.StringValue()
		} else {
			res = "An error occurred"
			break
		}
		_, err := storeDao.getOdds(i.Member.User.ID, newOdds.name)
		if err != nil {
			if err == sql.ErrNoRows {
				if err = storeDao.setOdds(i.Member.User.ID, newOdds); err != nil {
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
		newOddsOpt := OddsOption{}
		if gameName, ok := om["game"]; ok {
			_, err := storeDao.getOdds(i.Member.User.ID, gameName.StringValue())
			if err != nil {
				if err == sql.ErrNoRows {
					res = "No odds game with this name was found"
					break
				}
				res = "An error occurred"
				break
			}

			if choice, ok := om["choice"]; ok {
				if moneyline, ok := om["moneyline"]; ok {
					newOddsOpt.name = choice.StringValue()
					newOddsOpt.moneyline = int(moneyline.IntValue())
					storeDao.setOddsOpt(i.Member.User.ID, gameName.StringValue(), newOddsOpt)

					res = fmt.Sprintf("Added option %v (%v) to odds game %v", newOddsOpt.name, newOddsOpt.moneyline, gameName.StringValue())
					break
				}
			}
		}
		res = "An error occurred"

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

	case "start":
		res = "An error occurred"

	case "info":
		options = options[0].Options
		om := makeOptsMap(options)
		if game, ok := om["game"]; ok {
			o, err := storeDao.getOdds(i.Member.User.ID, game.StringValue())
			if err != nil {
				if err == sql.ErrNoRows {
					res = "No odds game with this name under your account was found"
					break
				}
				res = "An error occurred"
				break
			}
			opts := ""
			for i := range o.options {
				opt := o.options[i]
				opts += "\n- " + opt.name + fmt.Sprintf(" (%v)", opt.moneyline)
			}
			var winner = "Undecided"
			if o.winner != "" {
				winner = o.winner
			}
			res = fmt.Sprintf(
				"**Game name:** %v\n**Created by:** <@%v>\n**Options:** %v\n**Winner:** %v\n",
				o.name, i.Member.User.ID, opts, winner,
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
				storeDao.setWallet(i.Member.User.ID, 100) // new users get 100 coins
				res = fmt.Sprintf("New user detected, have %v counts on the house", amt)
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

// bet
func betHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//om := makeOptsMap(i.ApplicationCommandData().Options)
	res := ""
	// TODO
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
