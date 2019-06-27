package session

import(
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/cli"
  "github.com/finione/Discord-Hack-Week/src/util"
  "go.uber.org/zap"
  "fmt"
)

var logger *zap.SugaredLogger

func init() {
  logger = util.GetSugaredLogger()
}

// Create guild/session function
func Create(client *discordgo.Session, name, userID, returnChannelID string, options *cli.Options) (*Session, error) {
  session := &Session{
    Name: name,
    OwnerID: userID,
  }

  // Create guild
  guild, err := client.GuildCreate(session.Name)

  if err != nil {
    logger.Errorf("Err: %v\n",err)
    return &Session{}, err
  }

  // Assign guild to session struct
  session.Guild = guild

  /** TODO: Move to OAuth2 (MAYBE) **/
  // err := client.GuildMemberAdd(user.Token, session.Guild.ID, user.ID, "", emptyStringMap, false, false)

  // Create the welcome channel
  channel, err := client.GuildChannelCreate(session.Guild.ID, "Welcome", discordgo.ChannelTypeGuildText)
  if err != nil {
    logger.Errorf("Err: %v\n", err)
    return &Session{}, err
  }
  // Create an invite to the welcome channel
  invite, err := client.ChannelInviteCreate(channel.ID, discordgo.Invite{Temporary: true, MaxAge: 86400, MaxUses: 4})
  if err != nil {
    logger.Errorf("Err: %v\n", err)
    return &Session{}, err
  }

  // Get the user to join the server
  logger.Infof("Please click this link to join your new server: https://discord.gg/%v", invite.Code)
  if returnChannelID != "" {
    client.ChannelMessageSend(returnChannelID, fmt.Sprintf("Please click this link to join your new server: <https://discord.gg/%v>", invite.Code))
  }

  // Send welcome message
  client.ChannelMessageSend(channel.ID, fmt.Sprint("**Welcome to your new guild!**\n\nI've since left, but I hope you already feel at home.... feel free to change anything."))

  // Server layouts
  // TODO: Add Channels depending on type
  switch options.ServerType {
    case "Bot & Support": {
      // welcome, _ := client.GuildChannelCreate(session.Guild.ID, "Welcome", discordgo.ChannelTypeGuildText)
      // rules, _ := client.GuildChannelCreate(session.Guild.ID, "Rules", discordgo.ChannelTypeGuildText)
      
    }
    case "Support": {}
    case "Fun": {}
    case "Project": {}
    case  "Product / Service": {}
    default: {}
  }

  
  return session, nil
}
