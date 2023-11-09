// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/geth/accountcmd.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/console"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	"github.com/urfave/cli/v2"
)

var AccountCommand = &cli.Command{
	Name:     "account",
	Usage:    "Manage accounts",
	Category: "ACCOUNT COMMANDS",
	Description: `
Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.

It supports interactive mode, when you are prompted for password as well as
non-interactive mode where passwords are supplied via a given password file.
Non-interactive mode is only meant for scripted use on test networks or known
safe environments.

Make sure you remember the password you gave when creating a new account (with
either new or import). Without it you are not able to unlock your account.

Note that exporting your key in unencrypted format is NOT supported.

Keys are stored under <DATADIR>/keystore.
It is safe to transfer the entire directory or the individual keys therein
between klay nodes by simply copying.

Make sure you backup your keys regularly.`,
	Before: beforeAccountCmd,
	Subcommands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "Print summary of existing accounts",
			Action: accountList,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.KeyStoreDirFlag,
			},
			Description: `
Print a short summary of all accounts`,
		},
		{
			Name:   "new",
			Usage:  "Create a new account",
			Action: accountCreate,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.KeyStoreDirFlag,
				utils.PasswordFileFlag,
				utils.LightKDFFlag,
			},
			Description: `
Creates a new account and prints the address.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.

For non-interactive use the passphrase can be specified with the --password flag:

Note, this is meant to be used for testing only, it is a bad idea to save your
password to file or expose in any other way.
`,
		},
		{
			Name:      "update",
			Usage:     "Update an existing account",
			Action:    accountUpdate,
			ArgsUsage: "<address>",
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.KeyStoreDirFlag,
				utils.LightKDFFlag,
			},
			Description: `
Update an existing account.

The account is saved in the newest version in encrypted format, you are prompted
for a passphrase to unlock the account and another to save the updated file.

This same command can therefore be used to migrate an account of a deprecated
format to the newest format or change the password for an account.

For non-interactive use the passphrase can be specified with the --password flag:

Since only one password can be given, only format update can be performed,
changing your password is only possible interactively.
`,
		},
		{
			Name:   "import",
			Usage:  "Import a private key into a new account",
			Action: accountImport,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.KeyStoreDirFlag,
				utils.PasswordFileFlag,
				utils.LightKDFFlag,
			},
			ArgsUsage: "<keyFile>",
			Description: `
Imports an unencrypted private key from <keyfile> and creates a new account.
Prints the address.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.

For non-interactive use the passphrase can be specified with the --password flag:

Note, as you can directly copy your encrypted accounts to another klay instance,
this import mechanism is not needed when you transfer an account between
nodes.
`,
		},
		{
			Name:   "bls-info",
			Usage:  "Calculate BLS public key info",
			Action: accountBlsInfo,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.NodeKeyFileFlag,
				utils.NodeKeyHexFlag,
				utils.BlsNodeKeyFileFlag,
				utils.BlsNodeKeyHexFlag,
			},
			Description: `
Calculate BLS public key info (the public key and proof-of-possession)
then saves to bls-publicinfo-NODEID.json.

EXAMPLES

# From the nodekey file at the default location (DATADIR/nodekey)
kcn account bls-info

# From the nodekey file at a custom location
kcn account bls-info --nodekey /data/nodekey

# From the custom BLS key file at a custom location
kcn account bls-info --bls-nodekey DATADIR/bls-nodekey
`,
		},
		{
			Name:   "bls-import",
			Usage:  "Import a BLS private key from an EIP-2335 keystore JSON",
			Action: accountBlsImport,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.BlsNodeKeystoreFileFlag,
				utils.PasswordFileFlag,
			},
			Description: `
Decrypt an EIP-2335 keystore JSON and save the secret key to default location DATADIR/bls-nodekey.

EXAMPLES

kcn account bls-decrypt --bls-nodekeystore bls-keystore.json
`,
		},
		{
			Name:   "bls-export",
			Usage:  "Export a BLS private key to an EIP-2335 keystore JSON",
			Action: accountBlsExport,
			Flags: []cli.Flag{
				utils.DataDirFlag,
				utils.PasswordFileFlag,
				utils.LightKDFFlag,
			},
			Description: `
Export the BLS private key from the default location DATADIR/bls-nodekey
to an EIP-2335 keystore JSON bls-keystore-NODEID.json.
`,
		},
	},
}

func beforeAccountCmd(ctx *cli.Context) error {
	// Silence INFO logs from MakeConfigNode() or SetNodeConfig()
	// Account commands are almost independent from the regular node operation,
	// so INFO logs about networking or chain config are not necessary.
	if glogger, err := debug.GetGlogger(); err == nil {
		log.ChangeGlobalLogLevel(glogger, log.Lvl(log.LvlError))
	}
	return nil
}

func accountList(ctx *cli.Context) error {
	stack, _ := utils.MakeConfigNode(ctx)
	var index int
	for _, wallet := range stack.AccountManager().Wallets() {
		for _, account := range wallet.Accounts() {
			fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
			index++
		}
	}
	return nil
}

// tries unlocking the specified account a few times.
func UnlockAccount(ctx *cli.Context, ks *keystore.KeyStore, address string, i int, passwords []string) (accounts.Account, string) {
	account, err := utils.MakeAddress(ks, address)
	if err != nil {
		log.Fatalf("Could not list accounts: %v", err)
	}
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trials+1, 3)
		password := getPassPhrase(prompt, false, i, passwords)
		err = ks.Unlock(account, password)
		if err == nil {
			logger.Info("Unlocked account", "address", account.Address.Hex())
			return account, password
		}
		if err, ok := err.(*keystore.AmbiguousAddrError); ok {
			logger.Info("Unlocked account", "address", account.Address.Hex())
			return ambiguousAddrRecovery(ks, err, password), password
		}
		if err != keystore.ErrDecrypt {
			// No need to prompt again if the error is not decryption-related.
			break
		}
	}
	// All trials expended to unlock account, bail out
	log.Fatalf("Failed to unlock account %s (%v)", address, err)

	return accounts.Account{}, ""
}

// getPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool, i int, passwords []string) string {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if i < len(passwords) {
			return passwords[i]
		}
		return passwords[len(passwords)-1]
	}
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		log.Fatalf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			log.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			log.Fatalf("Passphrases do not match")
		}
	}
	return password
}

func ambiguousAddrRecovery(ks *keystore.KeyStore, err *keystore.AmbiguousAddrError, auth string) accounts.Account {
	fmt.Printf("Multiple key files exist for address %x:\n", err.Addr)
	for _, a := range err.Matches {
		fmt.Println("  ", a.URL)
	}
	fmt.Println("Testing your passphrase against all of them...")
	var match *accounts.Account
	for _, a := range err.Matches {
		if err := ks.Unlock(a, auth); err == nil {
			match = &a
			break
		}
	}
	if match == nil {
		log.Fatalf("None of the listed files could be unlocked.")
	}
	fmt.Printf("Your passphrase unlocked %s\n", match.URL)
	fmt.Println("In order to avoid this warning, you need to remove the following duplicate key files:")
	for _, a := range err.Matches {
		if a != *match {
			fmt.Println("  ", a.URL)
		}
	}
	return *match
}

// accountCreate creates a new account into the keystore defined by the CLI flags.
func accountCreate(ctx *cli.Context) error {
	cfg := utils.KlayConfig{Node: utils.DefaultNodeConfig()}
	// Load config file.
	if file := ctx.String(utils.ConfigFileFlag.Name); file != "" {
		if err := utils.LoadConfig(file, &cfg); err != nil {
			log.Fatalf("%v", err)
		}
	}
	cfg.SetNodeConfig(ctx)
	scryptN, scryptP, keydir, err := cfg.Node.AccountConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	password := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

	address, err := keystore.StoreKey(keydir, password, scryptN, scryptP)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", address)
	return nil
}

// accountUpdate transitions an account from a previous format to the current
// one, also providing the possibility to change the pass-phrase.
func accountUpdate(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		log.Fatalf("No accounts specified to update")
	}
	stack, _ := utils.MakeConfigNode(ctx)
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	for _, addr := range ctx.Args().Slice() {
		account, oldPassword := UnlockAccount(ctx, ks, addr, 0, nil)
		newPassword := getPassPhrase("Please give a new password. Do not forget this password.", true, 0, nil)
		if err := ks.Update(account, oldPassword, newPassword); err != nil {
			log.Fatalf("Could not update the account: %v", err)
		}
	}
	return nil
}

func accountImport(ctx *cli.Context) error {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		log.Fatalf("keyfile must be given as argument")
	}
	key, err := crypto.LoadECDSA(keyfile)
	if err != nil {
		log.Fatalf("Failed to load the private key: %v", err)
	}
	stack, _ := utils.MakeConfigNode(ctx)
	passphrase := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	acct, err := ks.ImportECDSA(key, passphrase)
	if err != nil {
		log.Fatalf("Could not create the account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", acct.Address)
	if _acct, err := ks.Find(acct); err == nil {
		fmt.Println("Your account is imported at", _acct.URL.Path)
	}
	return nil
}

func accountBlsInfo(ctx *cli.Context) error {
	utils.CheckExclusive(ctx,
		utils.BlsNodeKeyFileFlag,
		utils.BlsNodeKeyHexFlag,
		utils.BlsNodeKeystoreFileFlag,
	)

	// Parse CLI arguments in the same way as running a node.
	var (
		_, cfg  = utils.MakeConfigNode(ctx)
		nodeKey = cfg.Node.NodeKey()
		blsPriv = cfg.Node.BlsNodeKey()
	)
	if ctx.IsSet(utils.BlsNodeKeystoreFileFlag.Name) {
		if key, err := loadBlsNodeKeystore(ctx); err != nil {
			return err
		} else {
			blsPriv = key
		}
	}

	// Calculate filename from node address.
	nodeAddr := crypto.PubkeyToAddress(nodeKey.PublicKey)
	name := fmt.Sprintf("bls-publicinfo-%s.json", nodeAddr.Hex())

	// Calculate public key info.
	pub := blsPriv.PublicKey().Marshal()
	pop := bls.PopProve(blsPriv).Marshal()
	info := map[string]string{
		"address": nodeAddr.Hex(),
		"pub":     hex.EncodeToString(pub),
		"pop":     hex.EncodeToString(pop),
	}
	infoJson, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	infoJson = append(infoJson, '\n')

	writeFile(name, infoJson, 0o644) // Ordinary non-secret text file permission.
	return nil
}

func accountBlsImport(ctx *cli.Context) error {
	if !ctx.IsSet(utils.BlsNodeKeystoreFileFlag.Name) {
		return errors.New("no BLS keystore file specified")
	}
	blsPriv, err := loadBlsNodeKeystore(ctx)
	if err != nil {
		return err
	}
	blsPrivStr := []byte(hex.EncodeToString(blsPriv.Marshal()))

	// Parse CLI arguments in the same way as running a node.
	_, cfg := utils.MakeConfigNode(ctx)
	path := cfg.Node.ResolvePath("bls-nodekey")

	fmt.Printf("Importing BLS key: pub=%x\n", blsPriv.PublicKey().Marshal())
	writeFile(path, blsPrivStr, 0o400) // Secret file permission.
	return nil
}

func accountBlsExport(ctx *cli.Context) error {
	var (
		_, cfg                 = utils.MakeConfigNode(ctx)
		nodeKey                = cfg.Node.NodeKey()
		path                   = cfg.Node.ResolvePath("bls-nodekey")
		scryptN, scryptP, _, _ = cfg.Node.AccountConfig()
	)

	blsPriv, err := bls.LoadKey(path)
	if err != nil {
		log.Fatalf("Cannot load BLS key at '%s': %v", path, err)
	}

	storedPasswords := utils.MakePasswordList(ctx)
	password := getPassPhrase("Please give a new password. Do not forget this password.",
		true, 0, storedPasswords)

	plainKeystore := keystore.NewKeyEIP2335(blsPriv)
	encryptedJson, err := keystore.EncryptKeyEIP2335(plainKeystore, password, scryptN, scryptP)
	if err != nil {
		return err
	}

	nodeAddr := crypto.PubkeyToAddress(nodeKey.PublicKey)
	name := fmt.Sprintf("bls-keystore-%s.json", nodeAddr.Hex())

	fmt.Printf("Exporting BLS key: pub=%x\n", blsPriv.PublicKey().Marshal())
	writeFile(name, encryptedJson, 0o400) // Secret file permission.
	return nil
}

func loadBlsNodeKeystore(ctx *cli.Context) (bls.SecretKey, error) {
	file := ctx.String(utils.BlsNodeKeystoreFileFlag.Name)
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	storedPasswords := utils.MakePasswordList(ctx)
	password := getPassPhrase("Enter the password",
		false, 0, storedPasswords)

	plainKeystore, err := keystore.DecryptKeyEIP2335(content, password)
	if err != nil {
		return nil, err
	}
	return plainKeystore.SecretKey, nil
}

func writeFile(path string, content []byte, perm fs.FileMode) {
	if _, err := os.Stat(path); err == nil {
		log.Fatalf("file '%s' already exists", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatalf("cannot write '%s': %v", path, err)
	}
	if err := os.WriteFile(path, content, perm); err != nil {
		log.Fatalf("cannot write '%s': %v", path, err)
	}
	fmt.Printf("Successfully wrote '%s'\n", path)
}
