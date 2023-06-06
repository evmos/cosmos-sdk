package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
)

var (
	SetWithdrawerAddressMsg        = sdk.MsgTypeURL(&distributiontypes.MsgSetWithdrawAddress{})
	WithdrawDelegatorRewardMsg     = sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{})
	WithdrawValidatorCommissionMsg = sdk.MsgTypeURL(&distributiontypes.MsgWithdrawValidatorCommission{})
)

func TestAuthzAuthorizations(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	testCases := []struct {
		name       string
		msgTypeUrl string
		msg        sdk.Msg
		expUpdated distributiontypes.DistributionAuthorization
		allowed    []string
		expectErr  bool
	}{
		{
			"fail - set withdrawer address not in allowed list",
			SetWithdrawerAddressMsg,
			&distributiontypes.MsgSetWithdrawAddress{
				DelegatorAddress: "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				WithdrawAddress:  "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf48",
			},
			distributiontypes.DistributionAuthorization{
				MessageType: SetWithdrawerAddressMsg,
				AllowedList: []string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"},
			},
			[]string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"},
			true,
		},
		{
			"fail - withdraw validator commission address not in allowed list",
			WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgWithdrawValidatorCommission{
				ValidatorAddress: "cosmosvaloper1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			},
			distributiontypes.DistributionAuthorization{
				MessageType: WithdrawValidatorCommissionMsg,
				AllowedList: []string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"},
			},
			[]string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"},
			true,
		},
		{
			"fail - withdraw delegator rewards address not in allowed list",
			WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				ValidatorAddress: "cosmosvaloper1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			},
			distributiontypes.DistributionAuthorization{
				MessageType: WithdrawValidatorCommissionMsg,
				AllowedList: []string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf42"},
			},
			[]string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf42"},
			true,
		},
		{
			"success - set withdrawer address in allowed list",
			WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgSetWithdrawAddress{
				DelegatorAddress: "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				WithdrawAddress:  "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf46",
			},
			distributiontypes.DistributionAuthorization{
				MessageType: WithdrawValidatorCommissionMsg,
				AllowedList: []string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf46"},
			},
			[]string{"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf46"},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			distAuth := distributiontypes.NewDistributionAuthorization(tc.msgTypeUrl, tc.allowed)
			resp, err := distAuth.Accept(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if resp.Updated != nil {
					require.Equal(t, tc.expUpdated.String(), resp.Updated.String())
				}
			}
		})
	}
}
