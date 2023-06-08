package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
)

var (
	validatorAddr   = "cosmosvaloper1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"
	delegatorAddr   = "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"
	randomAddr      = "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf42"
	withdrawerAddr  = "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf48"
	withdrawerAddr2 = "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf46"
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
			distributiontypes.SetWithdrawerAddressMsg,
			&distributiontypes.MsgSetWithdrawAddress{
				DelegatorAddress: delegatorAddr,
				WithdrawAddress:  withdrawerAddr,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.SetWithdrawerAddressMsg,
				AllowedList: []string{delegatorAddr},
			},
			[]string{delegatorAddr},
			true,
		},
		{
			"fail - withdraw validator commission address not in allowed list",
			distributiontypes.WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgWithdrawValidatorCommission{
				ValidatorAddress: validatorAddr,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.WithdrawValidatorCommissionMsg,
				AllowedList: []string{delegatorAddr},
			},
			[]string{delegatorAddr},
			true,
		},
		{
			"fail - withdraw delegator rewards address not in allowed list",
			distributiontypes.WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: delegatorAddr,
				ValidatorAddress: validatorAddr,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.WithdrawValidatorCommissionMsg,
				AllowedList: []string{randomAddr},
			},
			[]string{randomAddr},
			true,
		},
		{
			"success - set withdrawer address in allowed list",
			distributiontypes.WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgSetWithdrawAddress{
				DelegatorAddress: delegatorAddr,
				WithdrawAddress:  withdrawerAddr2,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.WithdrawValidatorCommissionMsg,
				AllowedList: []string{withdrawerAddr2},
			},
			[]string{withdrawerAddr2},
			false,
		},
		{
			"success - withdraw delegator rewards address in allowed list",
			distributiontypes.WithdrawDelegatorRewardMsg,
			&distributiontypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: delegatorAddr,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.WithdrawDelegatorRewardMsg,
				AllowedList: []string{delegatorAddr},
			},
			[]string{delegatorAddr},
			false,
		},
		{
			"success - withdraw validator commission address in allowed list",
			distributiontypes.WithdrawValidatorCommissionMsg,
			&distributiontypes.MsgWithdrawValidatorCommission{
				ValidatorAddress: validatorAddr,
			},
			distributiontypes.DistributionAuthorization{
				MessageType: distributiontypes.WithdrawValidatorCommissionMsg,
				AllowedList: []string{validatorAddr},
			},
			[]string{validatorAddr},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			distAuth := distributiontypes.NewDistributionAuthorization(tc.msgTypeUrl, tc.allowed...)
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
