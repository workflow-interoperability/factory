package controller

import (
	"context"
	"time"

	"github.com/workflow-interoperability/factory/grpc"
	"github.com/workflow-interoperability/factory/model"
)

type Activity struct{}

func (Activity) GetProperties(ctx context.Context, in *grpc.ActivityGetPropertiesRq) (*grpc.ActivityGetPropertiesRs, error) {
	activity, err := model.GetActivity(in.ActivityKey)
	if err != nil {
		return &grpc.ActivityGetPropertiesRs{}, err
	}
	return &grpc.ActivityGetPropertiesRs{
		Key:            activity.Key,
		State:          int32(activity.State),
		Name:           activity.Name,
		Description:    activity.Description,
		ValidStates:    activity.ValidStates,
		InstanceKey:    activity.InstanceKey,
		RemoteInstance: activity.RemoteInstance,
		StartedDate:    activity.StartedDate,
		DueDate:        activity.DueDate,
	}, nil
}

func (a Activity) SetProperties(ctx context.Context, in *grpc.ActivitySetPropertiesRs) (*grpc.ActivityGetPropertiesRs, error) {
	activity := model.Activity{
		Key:          in.Key,
		State:        int(in.State),
		Name:         in.Name,
		Description:  in.Description,
		LastModified: time.Now().Unix(),
		// TODO: handle `Data`
	}
	if err := model.UpdateActivity(in.Key, activity); err != nil {
		return &grpc.ActivityGetPropertiesRs{}, err
	}
	return a.GetProperties(ctx, &grpc.ActivityGetPropertiesRq{
		ActivityKey: in.Key,
	})
}
func (Activity) CompleteActivity(ctx context.Context, in *grpc.ActivityCompleteActivityRq) (*grpc.ActivityCompleteActivityRs, error) {
	// TODO: handle `option`
	activity, err := model.GetActivity(in.Key)
	if err != nil {
		return &grpc.ActivityCompleteActivityRs{}, err
	}
	activity.DueDate = time.Now().Unix()
	return &grpc.ActivityCompleteActivityRs{}, model.UpdateActivity(in.Key, activity)
}
