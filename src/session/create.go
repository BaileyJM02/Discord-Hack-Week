package session

import(
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/cli"
)

// Create guild/session function
func Create(client *discordgo.Session, name string, user *discordgo.User, options cli.Options) (session *Session, error error) {
  session = &Session{
    Name:name,
    Owner:user,
  }

  guild, error := client.GuildCreate(session.Name)

  session.Guild = guild

  if error != nil {
    return
  }

  // todo: init guild creation scripts

  return
}
