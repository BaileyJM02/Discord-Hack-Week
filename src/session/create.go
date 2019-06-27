package session

import(
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/cli"
  "github.com/finione/Discord-Hack-Week/src/util"
  "go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
  logger := util.GetSugaredLogger()
}

// Create guild/session function
func Create(client *discordgo.Session, name string, user *discordgo.User, options cli.Options) (session *Session, error error) {
  session = &Session{
    Name:name,
    Owner:user,
  }

  guild, error := client.GuildCreate(session.Name)

  session.Guild = guild

  if error != nil {
    logger.Errorf("Err: %v\n", error)
    return
  }
  /** TODO: Move to OAuth2 **/
  // err := client.GuildMemberAdd(user.Token, session.Guild.ID, user.ID, "", emptyStringMap, false, false)

  channel, error := client.GuildChannelCreate(session.Guild.ID, "Welcome", discordgo.ChannelTypeGuildText)
  if error != nil {
    logger.Errorf("Err: %v\n", error)
    return
  }
  invite, error := client.ChannelInviteCreate(channel.ID, discordgo.Invite{Temporary: true, MaxAge: 900000000000, MaxUses: 4})
  if error != nil {
    logger.Errorf("Err: %v\n", error)
  }

  logger.Infof("Please click this link to join your new server: https://discord.gg/%v", invite.Code)

  
  // todo: init guild creation scripts

  return
}
