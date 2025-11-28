package content

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkiling/statemachine"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

type state struct {
	status statemachine.Status
	step   string
}

type runnerCommon interface {
	Complete(ctx context.Context, stateID uuid.UUID) (state state, executeErr error, err error)
	GetStateByID(ctx context.Context, stateID uuid.UUID) (state, error)
	TargetDeliveryStatus() DeliveryStatus
	ToDeliveryStatus(status statemachine.Status) DeliveryStatus
	RunnerType() runners.Type
}

type deliveryRunner struct {
	runner TVShowDeliveryState
}

func (d deliveryRunner) RunnerType() runners.Type {
	return runners.TVShowDelivery
}

func (d deliveryRunner) TargetDeliveryStatus() DeliveryStatus {
	return DeliveryStatusInProgress
}

func (d deliveryRunner) ToDeliveryStatus(status statemachine.Status) DeliveryStatus {
	switch status {
	case statemachine.CompletedStatus:
		return DeliveryStatusDelivered
	case statemachine.FailedStatus:
		return DeliveryStatusFailed
	default:
		return DeliveryStatusInProgress
	}
}

func (d deliveryRunner) Complete(ctx context.Context, stateID uuid.UUID) (st state, executeErr error, err error) {
	res, err1, err2 := d.runner.Complete(ctx, stateID)

	if res != nil {
		st = state{
			status: res.Status,
			step:   string(res.Step),
		}
	}
	return st, err1, err2
}

func (d deliveryRunner) GetStateByID(ctx context.Context, stateID uuid.UUID) (state, error) {
	res, err := d.runner.GetStateByID(ctx, stateID)
	if err != nil {
		return state{}, err
	}
	if res == nil {
		return state{}, fmt.Errorf("failed to complete state")
	}
	return state{
		status: res.Status,
		step:   string(res.Step),
	}, nil
}

type deleteRunner struct {
	runner TVShowDeleteState
}

func (d deleteRunner) RunnerType() runners.Type {
	return runners.TVShowDelete
}

func (d deleteRunner) TargetDeliveryStatus() DeliveryStatus {
	return DeliveryStatusDeleting
}

func (d deleteRunner) ToDeliveryStatus(status statemachine.Status) DeliveryStatus {
	switch status {
	case statemachine.CompletedStatus:
		return DeliveryStatusDeleted
	case statemachine.FailedStatus:
		return DeliveryStatusFailed
	default:
		return DeliveryStatusDeleting
	}
}

func (d deleteRunner) Complete(ctx context.Context, stateID uuid.UUID) (st state, executeErr error, err error) {
	res, err1, err2 := d.runner.Complete(ctx, stateID)
	if res == nil {
		return state{}, nil, fmt.Errorf("failed to complete state")
	}
	return state{
		status: res.Status,
		step:   string(res.Step),
	}, err1, err2
}

func (d deleteRunner) GetStateByID(ctx context.Context, stateID uuid.UUID) (state, error) {
	res, err := d.runner.GetStateByID(ctx, stateID)
	if err != nil {
		return state{}, err
	}
	if res == nil {
		return state{}, fmt.Errorf("failed to complete state")
	}
	return state{
		status: res.Status,
		step:   string(res.Step),
	}, nil
}
