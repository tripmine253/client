// Copyright 2020 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package client

import (
	"errors"
	"fmt"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"golang.org/x/net/context"
)

const (
	defaultTeamSearchLimit = 3
)

type CmdTeamSearch struct {
	libkb.Contextified

	query string
	limit int
}

func newCmdTeamSearch(cl *libcmdline.CommandLine, g *libkb.GlobalContext) cli.Command {
	return cli.Command{
		Name:         "search",
		ArgumentHelp: "<query>",
		Usage:        "Search for open teams on Keybase.",
		Action: func(c *cli.Context) {
			cmd := NewCmdTeamSearchRunner(g)
			cl.ChooseCommand(cmd, "search", c)
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "limit",
				Value: defaultTeamSearchLimit,
				Usage: fmt.Sprintf(
					"How many teams to return at most (default %d)",
					defaultTeamSearchLimit),
			},
		},
	}
}

func NewCmdTeamSearchRunner(g *libkb.GlobalContext) *CmdTeamSearch {
	return &CmdTeamSearch{Contextified: libkb.NewContextified(g)}
}

func (c *CmdTeamSearch) ParseArgv(ctx *cli.Context) (err error) {
	if len(ctx.Args()) != 1 {
		return errors.New("usage: keybase team search <query>")
	}
	c.query = ctx.Args().Get(0)
	c.limit = ctx.Int("limit")
	return nil
}

func (c *CmdTeamSearch) Run() error {
	ctx, ctxCancel := context.WithCancel(context.TODO())
	defer ctxCancel()
	ctx = libkb.WithLogTag(ctx, "CTS")

	cli, err := GetTeamSearchClient(c.G())
	if err != nil {
		return err
	}

	res, err := cli.TeamSearch(ctx, keybase1.TeamSearchArg{
		Query: c.query,
		Limit: c.limit,
	})
	if err != nil {
		return err
	}

	ui := c.G().UI.GetTerminalUI()
	if len(res.Results) == 0 {
		ui.Printf("No results found.\n")
		return nil
	}

	for _, team := range res.Results {
		ui.Printf("%s (%d members)\n", team.Name, team.MemberCount)
		if team.Description != nil {
			ui.Printf("\t%s\n", *team.Description)
		}
		if !team.InTeam {
			ui.Printf("\tYou can join this open team with `keybase team request-access %s`\n", team.Name)
		}
	}

	return nil
}

func (c *CmdTeamSearch) GetUsage() libkb.Usage {
	return libkb.Usage{
		Config:    true,
		API:       true,
		KbKeyring: true,
	}
}