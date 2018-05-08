package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/anthdm/consenter/pkg/consensus"
	"github.com/anthdm/consenter/pkg/consensus/solo"
	"github.com/anthdm/consenter/pkg/network"
	"github.com/urfave/cli"
)

var (
	errMissingPrivateKey = errors.New("missing private key for consensus node")
)

func main() {
	ctl := cli.NewApp()
	ctl.Name = "concenter"
	ctl.Usage = "Pluggable blockchain consensus simulation framework"
	ctl.Commands = []cli.Command{
		newNodeCommand(),
	}
	ctl.Run(os.Args)
}

func newNodeCommand() cli.Command {
	return cli.Command{
		Name:   "node",
		Usage:  "Start a single consenter node",
		Action: startServer,
		Flags: []cli.Flag{
			cli.IntFlag{Name: "tcp"},
			cli.StringFlag{Name: "seed"},
			cli.BoolFlag{Name: "consensus"},
			cli.StringFlag{Name: "privkey"},
			cli.StringFlag{Name: "engine"},
		},
	}
}

func startServer(ctx *cli.Context) error {
	var (
		isConsensusNode = ctx.Bool("consensus")
		privKey         *ecdsa.PrivateKey
		err             error
		engine          consensus.Engine
	)
	if isConsensusNode {
		if len(ctx.String("privkey")) == 0 {
			return cli.NewExitError(errMissingPrivateKey, 1)
		}
		// For simulation we dont need the right (bitcoin) elliptic curve points.
		// The default P256 curve that comes with the Go stdlib will do.
		privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if len(ctx.String("engine")) == 0 {
			return cli.NewExitError("engine cannot be empty if running a consensus node", 1)
		}
		switch ctx.String("engine") {
		case "solo":
			// block generation 15 seconds.
			engine = solo.NewEngine(15 * time.Second)
		default:
			return cli.NewExitError("invalid engine option", 1)
		}
	}
	cfg := network.ServerConfig{
		ListenAddr:     ctx.Int("tcp"),
		DialTimeout:    3 * time.Second,
		BootstrapNodes: parseSeeds(ctx.String("seed")),
		Consensus:      isConsensusNode,
		PrivateKey:     privKey,
	}
	srv := network.NewServer(cfg, engine)
	return cli.NewExitError(srv.Start(), 1)
}

func parseSeeds(str string) []string {
	if len(str) == 0 {
		return nil
	}
	return strings.Split(str, ",")
}
