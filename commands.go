package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	storeDao store

	dmPerm = false

	commands = []discordgo.ApplicationCommand{

		{
			Name:         "odds",
			Description:  "manage a game of odds",
			DMPermission: &dmPerm,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "new",
					Description: "start a new game of odds",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "name of what you are betting on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},

				{
					Name:        "del",
					Description: "deletes a game of odds",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "name of the odds game to delete",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},

				{
					Name:        "delchoice",
					Description: "deletes a bet choice from a game of odds",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "game",
							Description: "game to delete option from",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "choice",
							Description: "name of bet choice to delete",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},

				{
					Name:        "add",
					Description: "add a choice to your odds game",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "game",
							Description: "the odds game this option will be added to",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "choice",
							Description: "name of the new option to bet on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "moneyline",
							Description: "what % of your bet you get back if you win; neg = favourite, positive = underdog",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "start",
					Description: "let the betting commence",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "game",
							Description: "the name of the odds game to start",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "info",
					Description: "information about this odds game",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "game",
							Description: "the name of the odds game to get info for",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommand,
					// TODO add option for getting info about other people's games
				},
			},
		},

		{
			Name:         "bet",
			DMPermission: &dmPerm,
			Description:  "manage bets on an odds game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "make",
					Description: "make a bet",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "the id of the odds game to bet on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "option",
							Description: "the name of the option to bet on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "amount",
							Description: "the amount to bet",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "del",
					Description: "delete a bet",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "the id of the odds game betted on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "option",
							Description: "the name of the option to betted on",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},

		{
			Name:         "wallet",
			DMPermission: &dmPerm,
			Description:  "info about your wallet",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "get",
					Description: "get your wallet info",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},

		{
			Name:        "ping",
			Description: "health check",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"odds":   oddsHandler,
		"wallet": walletHandler,
		"ping":   pingHandler,
		"bet":    betHandler,
	}
)

func makeOptsMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	m := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		m[opt.Name] = opt
	}
	return m
}
