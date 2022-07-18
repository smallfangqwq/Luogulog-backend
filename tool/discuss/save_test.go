package discuss

import (
	"fmt"
	"luogulog/declare"
	"testing"

	"github.com/spf13/viper"
)

func TestGetDiscussReply(t *testing.T) {
	ConfingData := viper.New()
	ConfingData.SetConfigFile("../../config.yaml")
	var Configs declare.Config
	ConfingData.ReadInConfig()
	ConfingData.Unmarshal(&Configs)
	val, _ := GetDiscussReply(1,244341, Configs.Request)
	fmt.Print(val[0].Content)
}