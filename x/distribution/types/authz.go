package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"golang.org/x/exp/slices"
)

var _ authz.Authorization = &DistributionAuthorization{}

// NewDistributionAuthorization creates a new DistributionAuthorization.
func NewDistributionAuthorization(msgType string, allowed []string) *DistributionAuthorization {
	return &DistributionAuthorization{
		MessageType: msgType,
		AllowedList: allowed,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (m *DistributionAuthorization) MsgTypeURL() string {
	return m.MessageType
}

// Accept implements Authorization.Accept. It checks, that the
// withdrawer for MsgSetWithdrawAddress,
// validator for MsgWithdrawValidatorCommission,
// the delegator address for MsgWithdrawDelegatorReward
// is in the allowed list. If these conditions are met, the AcceptResponse is returned.
func (m *DistributionAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	switch msg := msg.(type) {
	case *MsgSetWithdrawAddress:
		if !slices.Contains(m.AllowedList, msg.WithdrawAddress) {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("address is not in the allowed list")
		}
	case *MsgWithdrawValidatorCommission:
		if !slices.Contains(m.AllowedList, msg.ValidatorAddress) {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("address is not in the allowed list")
		}
	case *MsgWithdrawDelegatorReward:
		if !slices.Contains(m.AllowedList, msg.DelegatorAddress) {
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
