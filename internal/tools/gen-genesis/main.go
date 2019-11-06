package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/logger"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/config"
	db "github.com/tendermint/tm-db"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randompassword() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var (
	validators = flag.String("validators", "engine", "list of validator names separated with a comma")
	chainid    = flag.String("chain-id", "mesg-chain", "chain id")
	path       = flag.String("path", ".genesis/", "genesis folder path")
)

const (
	tendermintPath = "tendermint"
	cosmosPath     = "cosmos"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	validatorNames := strings.Split(*validators, ",")
	passwords := make([]string, len(validatorNames))
	for i := 0; i < len(validatorNames); i++ {
		passwords[i] = randompassword()
	}

	// init app factory
	db, err := db.NewGoLevelDB("app", filepath.Join(*path, cosmosPath))
	if err != nil {
		logrus.Fatalln(err)
	}
	appFactory := cosmos.NewAppFactory(logger.TendermintLogger(), db)

	// register the backend modules to the app factory.
	enginesdk.NewBackend(appFactory)

	// init cosmos app
	app, err := cosmos.NewApp(appFactory)
	if err != nil {
		logrus.Fatalln(err)
	}

	// init keybase
	kb, err := cosmos.NewKeybase(filepath.Join(*path, cosmosPath))
	if err != nil {
		logrus.Fatalln(err)
	}

	// create validators
	vals := []cosmos.GenesisValidator{}
	peers := []string{}
	for i, valName := range validatorNames {
		cfg := config.DefaultConfig()
		cfg.SetRoot(filepath.Join(filepath.Join(*path, tendermintPath), valName))
		if err := os.MkdirAll(filepath.Dir(cfg.GenesisFile()), 0755); err != nil {
			logrus.Fatalln(err)
		}
		if err := os.MkdirAll(filepath.Join(cfg.DBDir()), 0755); err != nil {
			logrus.Fatalln(err)
		}
		mnemonic, err := kb.NewMnemonic()
		if err != nil {
			logrus.Fatalln(err)
		}
		acc, err := kb.CreateAccount(valName, mnemonic, "", passwords[i], 0, 0)
		if err != nil {
			logrus.Fatalln(err)
		}
		genVal, err := cosmos.NewGenesisValidator(kb,
			valName,
			passwords[i],
			cfg.PrivValidatorKeyFile(),
			cfg.PrivValidatorStateFile(),
			cfg.NodeKeyFile(),
		)
		if err != nil {
			logrus.Fatalln(err)
		}
		vals = append(vals, genVal)
		peer := fmt.Sprintf("%s@%s:26656", genVal.NodeID, genVal.Name)
		logrus.WithFields(map[string]interface{}{
			"name":     genVal.Name,
			"address":  acc.GetAddress().String,
			"password": genVal.Password,
			"mnemonic": mnemonic,
			"nodeID":   genVal.NodeID,
			"peer":     peer,
		}).Infof("Validator #%d\n", i+1)
		peers = append(peers, peer)
	}

	// generate and save genesis
	_, err = cosmos.GenGenesis(app.Cdc(), kb, app.DefaultGenesis(), *chainid, filepath.Join(*path, "genesis.json"), vals)
	if err != nil {
		logrus.Fatalln(err)
	}

	// save peers list
	if err := ioutil.WriteFile(filepath.Join(*path, "peers.txt"), []byte(strings.Join(peers, ",")), 0644); err != nil {
		log.Fatalln("error during writing peers file:", err)
	}

	logrus.Infof("genesis created with success in folder %q\n", *path)
}
