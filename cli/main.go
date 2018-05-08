package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"os"
	"strings"
	"time"

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
		},
	}
}

func startServer(ctx *cli.Context) error {
	var (
		isConsensusNode = ctx.Bool("consensus")
		privKey         *ecdsa.PrivateKey
		err             error
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
	}
	cfg := network.ServerConfig{
		ListenAddr:     ctx.Int("tcp"),
		DialTimeout:    3 * time.Second,
		BootstrapNodes: parseSeeds(ctx.String("seed")),
		Consensus:      isConsensusNode,
		PrivateKey:     privKey,
	}
	srv := network.NewServer(cfg)
	return cli.NewExitError(srv.Start(), 1)
}

func parseSeeds(str string) []string {
	if len(str) == 0 {
		return nil
	}
	return strings.Split(str, ",")
}
