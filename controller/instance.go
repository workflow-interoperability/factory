package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/workflow-interoperability/factory/grpc"
	"github.com/workflow-interoperability/factory/model"
)

type Instance struct{}

func (Instance) GetProperties(ctx context.Context, in *grpc.InstanceGetPropertiesRq) (*grpc.InstanceGetPropertiesRs, error) {
	instance, err := model.GetInstance(in.InstanceKey)
	if err != nil {
		return &grpc.InstanceGetPropertiesRs{}, err
	}
	return &grpc.InstanceGetPropertiesRs{
		InstanceKey: instance.Key,
		State:       grpc.State(instance.State),
		Name:        instance.Name,
		Subject:     instance.Subject,
		Description: instance.Description,
		FactoryKey:  instance.FactoryKey,
		Observers:   strings.Split(instance.Observers, "|"),
		ContextData: instance.ContextData,
		ResultData:  instance.ResultData,
	}, nil
}

func (i Instance) SetProperties(ctx context.Context, in *grpc.InstanceSetPropertiesRq) (*grpc.InstanceGetPropertiesRs, error) {
	instance := model.Instance{
		Key:         in.InstanceKey,
		Subject:     in.Subject,
		Description: in.Description,
		ContextData: in.ContextData,
	}
	if err := model.UpdateInstance(instance); err != nil {
		return &grpc.InstanceGetPropertiesRs{}, err
	}
	return i.GetProperties(ctx, &grpc.InstanceGetPropertiesRq{
		InstanceKey: in.InstanceKey,
	})
}

func (Instance) Subscribe(ctx context.Context, in *grpc.InstanceSubscribeRq) (*grpc.InstanceSubscribeRs, error) {
	instance, err := model.GetInstance(in.InstanceKey)
	if err != nil {
		return &grpc.InstanceSubscribeRs{}, err
	}
	observers := strings.Split(instance.Observers, "|")
	for _, v := range observers {
		if v == in.ObserverKey {
			return &grpc.InstanceSubscribeRs{}, err
		}
	}
	observers = append(observers, in.ObserverKey)
	instance.Observers = strings.Join(observers, "|")
	return &grpc.InstanceSubscribeRs{}, model.AddInstanceObserver(in.InstanceKey, instance.Observers)
}

func (Instance) Unsubscribe(ctx context.Context, in *grpc.InstanceSubscribeRq) (*grpc.InstanceSubscribeRs, error) {
	instance, err := model.GetInstance(in.InstanceKey)
	if err != nil {
		return &grpc.InstanceSubscribeRs{}, err
	}
	observers := strings.Split(instance.Observers, "|")
	for i, v := range observers {
		if v == in.ObserverKey {
			observers = append(observers[:i], observers[i+1:]...)
		}
	}
	instance.Observers = strings.Join(observers, "|")
	return &grpc.InstanceSubscribeRs{}, model.AddInstanceObserver(in.InstanceKey, instance.Observers)
}

// TODO
func (Instance) ListActivities(ctx context.Context, in *grpc.InstanceListActivitiesRq) (*grpc.InstanceListActivitiesRs, error) {
	return &grpc.InstanceListActivitiesRs{}, errors.New("not implemented yet")
}

// TODO
func (Instance) ChangeState(ctx context.Context, in *grpc.InstanceChangeStateRq) (*grpc.InstanceChangeStateRq, error) {
	instance := model.Instance{
		Key:   in.InstanceKey,
		State: int(in.State),
	}
	return &grpc.InstanceChangeStateRq{}, model.UpdateInstance(instance)
}
