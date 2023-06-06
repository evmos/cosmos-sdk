package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authz.Authorization = &DistributionAuthorization{}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (m *DistributionAuthorization) MsgTypeURL() string {
	return m.MessageType
}

// Accept implements Authorization.Accept. It checks, that the
// withdrawer for MsgSetWithdrawAddress,
// validator for MsgWithdrawValidatorCommission
// the delegator address for MsgWithdrawDelegatorReward
// is in the allowed list. If these conditions are met, the AcceptResponse is returned.
func (m *DistributionAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	switch msg := msg.(type) {
	case *MsgSetWithdrawAddress:
		if !checkAddressInList(msg.WithdrawAddress, m.AllowedList) {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("address is not in the allowed list")
		}
	case *MsgWithdrawValidatorCommission:
		if !checkAddressInList(msg.ValidatorAddress, m.AllowedList) {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("address is not in the allowed list")
		}
	case *MsgWithdrawDelegatorReward:
		if !checkAddressInList(msg.DelegatorAddress, m.AllowedList) {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("address is not in the allowed list")
		}
	default:
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidRequest.Wrap("unknown msg type")
	}

	return authz.AcceptResponse{
		Accept: true,
		Delete: false,
		Updated: &DistributionAuthorization{
			AllowedList: m.AllowedList,
			MessageType: m.MessageType,
		},
	}, nil
}

// ValidateBasic performs a stateless validation of the fields.
func (m *DistributionAuthorization) ValidateBasic() error {
	if len(m.AllowedList) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("allowed list cannot be empty")
	}

	return nil
}

// checkAddressInList checks if the given address is in the given list.
// If the list is empty, it returns true.
func checkAddressInList(address string, list []string) bool {
	for _, addr := range list {
		if addr == address {
			return true
		}
	}

	return false
}
