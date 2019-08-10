package session

import "github.com/bwmarrin/discordgo"

/**
we call a setup session a session for simplicity
a session extends a guild object, providing all the functionality a basic guild object has
**not to be mistaken with discordgo.Session (find a better name?)**
*/
type Session struct {
  Guild *discordgo.Guild
  // the user creating the session
  OwnerID string
  Name string
}
