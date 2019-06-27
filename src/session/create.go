package session

import(
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/cli"
)

func Create(client *discordgo.Session, name string, user *discordgo.User, cli.Options) (session *Session, error error) {
  session = &Session{
    Name:name,
    Owner:user,
  }
  session.Name = name
  session.Owner = user
  guild, error := client.GuildCreate(session.Name)

  session.Guild = guild

  if error != nil {
    return
  }

  // todo: init guild creation scripts

  return
}
