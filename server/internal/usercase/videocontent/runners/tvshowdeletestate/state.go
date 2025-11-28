package tvshowdeletestate

import (
	"github.com/kkiling/statemachine"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

type CreateState = statemachine.CreateState[TVShowDeleteData, runners.Metadata, StepDelete]
type State = statemachine.State[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type]
type Step = statemachine.Step[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type]
type StepRegistration = statemachine.StepRegistration[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type]
type StepContext = statemachine.StepContext[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type]
type StepResult = statemachine.StepResult[TVShowDeleteData, StepDelete]
type StateMachineService = statemachine.StateMachine[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type, CreateOptions]

func NewState(contentDelivery ContentDeleted, stateMachineStorage statemachine.Storage) *StateMachineService {
	return statemachine.NewService[TVShowDeleteData, runners.FailData, runners.Metadata, StepDelete, runners.Type, CreateOptions](
		statemachine.Config{},
		stateMachineStorage,
		NewTaskRunner(contentDelivery),
	)
}
