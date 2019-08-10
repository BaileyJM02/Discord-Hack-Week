package session

import (
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/util"
  "go.uber.org/goleak"
  "testing"
)

func TestMain(m *testing.M) {
  goleak.VerifyTestMain(m)
}

func TestCreate(t *testing.T) {

  client := new(discordgo.Session)
  client.GuildCreate = func(name string) (st *discordgo.Guild, err error) {
    util.GetSugaredLogger().Info("Correct call to GuildCreate")
    return
  }

  _, err  := Create(client, "test", "", "", nil)
  if err != nil {
    t.Fail()
  }
}