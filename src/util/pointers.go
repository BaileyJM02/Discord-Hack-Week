package util

import (
  "github.com/bwmarrin/discordgo"
)


func RemoveUserPointer(user *discordgo.User) discordgo.User {
	return discordgo.User{
		ID: user.ID,
    Email: user.Email,
    Username: user.Username,
    Avatar: user.Avatar,
    Locale: user.Locale,
    Discriminator: user.Discriminator,
    Token: user.Token,
    Verified: user.Verified,
    MFAEnabled: user.MFAEnabled,
    Bot: user.Bot,
	}
} 