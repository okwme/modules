package cli

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okwme/modules/incubator/faucet/internal/types"
)

// GetTxCmd return faucet sub-command for tx
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	faucetTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "faucet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetTxCmd.AddCommand(flags.PostCommands(
		GetCmdMint(cdc),
		GetCmdMintFor(cdc),
		GetCmdInitial(cdc),
		GetPublishKey(cdc),
	)...)

	return faucetTxCmd
}

// GetCmdWithdraw is the CLI command for mining coin
func GetCmdMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint",
		Short: "mint coin to sender address",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			denom := args[0]

			msg := types.NewMsgMint(cliCtx.GetFromAddress(), cliCtx.GetFromAddress(), denom)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdWithdraw is the CLI command for mining coin
func GetCmdMintFor(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mintfor [address] [denom]",
		Short: "mint coin for new address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address, _ := sdk.AccAddressFromBech32(args[0])

			denom := args[1]
			msg := types.NewMsgMint(cliCtx.GetFromAddress(), address, denom)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetPublishKey(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "Publish current account as an public faucet. Do NOT add many coins in this account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			kb, errkb := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			if errkb != nil {
				return errkb
			}

			// check local key
			armor, err := kb.Export(cliCtx.GetFromName())
			if err != nil {
				return err
			}

			msg := types.NewMsgFaucetKey(cliCtx.GetFromAddress(), armor)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdInitial(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize mint key for faucet",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			kb, errkb := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			if errkb != nil {
				return errkb
			}

			// check local key
			_, err := kb.Get(types.ModuleName)
			if err == nil {
				return errors.New("faucet existed")
			}

			// fetch from chain
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/key", types.ModuleName), nil)
			if err != nil {
				return nil
			}
			var rkey types.FaucetKey
			cdc.MustUnmarshalJSON(res, &rkey)

			if len(rkey.Armor) == 0 {
				return errors.New("Faucet key has not published")
			}
			// import to keybase
			kb.Import(types.ModuleName, rkey.Armor)
			return nil

		},
	}
}
