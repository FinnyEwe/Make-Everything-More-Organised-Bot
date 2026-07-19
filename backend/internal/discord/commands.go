package discord

import (
	"fmt"
	"log"

	"backend/internal/store"

	"github.com/bwmarrin/discordgo"
)

var minAmount = 0.01

var savingsCommand = &discordgo.ApplicationCommand{
	Name:        "savings",
	Description: "Edit the savings amount",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add to savings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "Amount to add",
					Required:    true,
					MinValue:    &minAmount,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "sub",
			Description: "Subtract from savings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "Amount to subtract",
					Required:    true,
					MinValue:    &minAmount,
				},
			},
		},
	},
}

func RegisterSavingsCommands(sess *discordgo.Session, st *store.Store) error {
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		handleSavingsCommand(s, i, st)
	})

	_, err := sess.ApplicationCommandBulkOverwrite(sess.State.User.ID, "", []*discordgo.ApplicationCommand{savingsCommand})
	if err != nil {
		return fmt.Errorf("register /savings: %w", err)
	}
	log.Println("Registered /savings command")
	return nil
}

func handleSavingsCommand(s *discordgo.Session, i *discordgo.InteractionCreate, st *store.Store) {
	data := i.ApplicationCommandData()
	if data.Name != "savings" {
		return
	}
	if len(data.Options) == 0 {
		respond(s, i, "Usage: `/savings add <amount>` or `/savings sub <amount>`")
		return
	}

	sub := data.Options[0]
	amount := sub.Options[0].FloatValue()

	var operand store.Operand
	switch sub.Name {
	case "add":
		operand = store.OperandAdd
	case "sub":
		operand = store.OperandSub
	default:
		respond(s, i, "Unknown subcommand. Use `add` or `sub`.")
		return
	}

	savings, err := st.UpdateSavings(operand, amount)
	if err != nil {
		log.Printf("UpdateSavings failed: %v", err)
		respond(s, i, "Failed to update savings.")
		return
	}

	verb := "Added"
	if operand == store.OperandSub {
		verb = "Subtracted"
	}
	respond(s, i, fmt.Sprintf("%s `$%.2f`. Savings is now `$%.2f`.", verb, amount, savings.Amount))
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		log.Printf("InteractionRespond failed: %v", err)
	}
}
