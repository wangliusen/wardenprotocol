package cli

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/warden-protocol/wardenprotocol/warden/x/warden/types/v1beta2"
)

// NewTxCmd returns a root CLI command handler for x/warden transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        v1beta2.ModuleName,
		Short:                      "Warden transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		FulfillKeyRequestTxCmd(),
		RejectKeyRequestTxCmd(),
	)

	return txCmd
}

func FulfillKeyRequestTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fulfill-key-request [request-id] [public-key-data]",
		Short: "Fulfill a key request providing the public key.",
		Long: `Fulfill a key request providing the public key.
The sender of this transaction must be a party of the Keychain for the request.
The public key must be a base64 encoded string.`,
		Example: fmt.Sprintf("%s tx warden fulfill-key-request 1234 aGV5dGhlcmU=", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reqId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pk, err := base64.StdEncoding.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := &v1beta2.MsgUpdateKeyRequest{
				Creator:   clientCtx.GetFromAddress().String(),
				Status:    v1beta2.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
				RequestId: reqId,
				Result:    v1beta2.NewMsgUpdateKeyRequestKey(pk),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RejectKeyRequestTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-key-request [request-id] [reason]",
		Short: "Reject a key request providing the reason.",
		Long: `Reject a key request providing a reason.
The sender of this transaction must be a party of the Keychain for the request.`,
		Example: fmt.Sprintf("%s tx warden reject-key-request 1234 'something happened'", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reqId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &v1beta2.MsgUpdateKeyRequest{
				Creator:   clientCtx.GetFromAddress().String(),
				Status:    v1beta2.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED,
				RequestId: reqId,
				Result:    v1beta2.NewMsgUpdateKeyRequestReject(args[1]),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
