package discord

import (
	"fmt"
	"log"
	"strings"
	"time"

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

var reminderCommand = &discordgo.ApplicationCommand{
	Name:        "reminder",
	Description: "Manage reminders",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add a new reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "What to remind you about",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "day",
					Description: "Day (1-31)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
					MaxValue:    31,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "month",
					Description: "Month (1-12)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
					MaxValue:    12,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "year",
					Description: "Year (e.g., 2024)",
					Required:    true,
					MinValue:    func() *float64 { v := 2024.0; return &v }(),
					MaxValue:    2030,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List all your reminders",
		},
	},
}

func RegisterSavingsCommands(sess *discordgo.Session, st *store.Store) error {
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		data := i.ApplicationCommandData()
		switch data.Name {
		case "savings":
			handleSavingsCommand(s, i, st)
		case "reminder":
			handleReminderCommand(s, i, st)
		}
	})

	_, err := sess.ApplicationCommandBulkOverwrite(sess.State.User.ID, "", []*discordgo.ApplicationCommand{
		savingsCommand,
		reminderCommand,
	})
	if err != nil {
		return fmt.Errorf("register commands: %w", err)
	}
	log.Println("Registered /savings and /reminder commands")
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

func handleReminderCommand(s *discordgo.Session, i *discordgo.InteractionCreate, st *store.Store) {
	data := i.ApplicationCommandData()
	if data.Name != "reminder" {
		return
	}
	if len(data.Options) == 0 {
		respond(s, i, "Usage: `/reminder add <description> <day> <month> <year>` or `/reminder list`")
		return
	}

	sub := data.Options[0]
	switch sub.Name {
	case "add":
		handleAddReminder(s, i, st, sub)
	case "list":
		handleListReminders(s, i, st)
	default:
		respond(s, i, "Unknown subcommand. Use `add` or `list`.")
	}
}

func handleAddReminder(s *discordgo.Session, i *discordgo.InteractionCreate, st *store.Store, sub *discordgo.ApplicationCommandInteractionDataOption) {
	if len(sub.Options) != 4 {
		respond(s, i, "Please provide description, day, month, and year.")
		return
	}

	description := sub.Options[0].StringValue()
	day := int(sub.Options[1].IntValue())
	month := int(sub.Options[2].IntValue())
	year := int(sub.Options[3].IntValue())

	// Validate date
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if date.Day() != day || date.Month() != time.Month(month) || date.Year() != year {
		respond(s, i, "Invalid date. Please check your day, month, and year values.")
		return
	}

	// Check if date is in the past
	loc, _ := time.LoadLocation("Australia/Sydney")
	now := time.Now().In(loc)
	dateInSydney := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	
	if dateInSydney.Before(now.Truncate(24 * time.Hour)) {
		respond(s, i, "Cannot create reminders for dates in the past.")
		return
	}

	// Format date as DD-MM-YYYY (matching your existing format)
	dateStr := fmt.Sprintf("%02d-%02d-%04d", day, month, year)

	// Create reminder
	st.CreateReminder(description, dateStr)

	respond(s, i, fmt.Sprintf("✅ Reminder created!\n**Description:** %s\n**Date:** %s", description, dateStr))
}

func handleListReminders(s *discordgo.Session, i *discordgo.InteractionCreate, st *store.Store) {
	reminders, err := st.GetReminders()
	if err != nil {
		log.Printf("GetReminders failed: %v", err)
		respond(s, i, "Failed to retrieve reminders.")
		return
	}

	if len(reminders) == 0 {
		respond(s, i, "📝 You have no reminders.")
		return
	}

	var builder strings.Builder
	builder.WriteString("📝 **Your Reminders:**\n\n")

	for _, reminder := range reminders {
		discordStatus := "⏰ Discord pending"
		if reminder.Discord {
			discordStatus = "✅ Sent to Discord"
		}
		gcalStatus := "⏰ Calendar pending"
		if reminder.Gcal {
			gcalStatus = "✅ Added to Google Calendar"
		}

		builder.WriteString(fmt.Sprintf("**#%d** - %s\n", reminder.ID, reminder.Date))
		builder.WriteString(fmt.Sprintf("📄 %s\n", reminder.Description))
		builder.WriteString(fmt.Sprintf("📊 %s · %s\n\n", discordStatus, gcalStatus))
	}

	content := builder.String()
	// Discord has a 2000 character limit for messages
	if len(content) > 1900 {
		content = content[:1900] + "\n... (truncated)"
	}

	respond(s, i, content)
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
