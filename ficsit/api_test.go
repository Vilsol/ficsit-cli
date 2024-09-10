package ficsit

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/MarvinJWendt/testza"

	"github.com/satisfactorymodding/ficsit-cli/cfg"
)

var client graphql.Client

func init() {
	cfg.SetDefaults()
	client = InitAPI()
}

func TestMods(t *testing.T) {
	response, err := Mods(context.Background(), client, ModFilter{})
	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, response)
	testza.AssertNotNil(t, response.Mods)
	testza.AssertNotNil(t, response.Mods.Mods)
	testza.AssertNotZero(t, response.Mods.Count)
	testza.AssertNotZero(t, len(response.Mods.Mods))
}
