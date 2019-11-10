package controller

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/workflow-interoperability/factory/grpc"
	"github.com/workflow-interoperability/factory/model"
)

// Factory implement factory grpc interface
type Factory struct{}

func (Factory) GetProperties(ctx context.Context, in *grpc.FactoryGetPropertiesRq) (*grpc.FactoryGetPropertiesRs, error) {
	factory, err := model.GetFactory(in.Key)
	if err != nil {
		return &grpc.FactoryGetPropertiesRs{}, err
	}
	return &grpc.FactoryGetPropertiesRs{
		Key:         factory.Key,
		Name:        factory.Name,
		Subject:     factory.Subject,
		Description: factory.Description,
	}, nil
}

func (Factory) SetProperties(ctx context.Context, in *grpc.FactoryGetPropertiesRs) (*grpc.FactoryGetPropertiesRq, error) {
	factory := model.Factory{
		Subject:     in.Subject,
		Description: in.Description,
		Name:        in.Name,
	}
	return &grpc.FactoryGetPropertiesRq{}, model.UpdateFactory(in.Key, factory)
}

func (Factory) CreateInstance(ctx context.Context, in *grpc.FactoryCreateInstanceRq) (*grpc.FactoryCreateInstanceRs, error) {
	// decode data
	var variables map[string]interface{}
	err := json.Unmarshal([]byte(in.ContextData), &variables)
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	// ask workflow engine to start a instance
	client, err := connectZeebe()
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	// get factory message
	factory, err := model.GetFactory(in.FactoryKey)
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	req, err := (*client).NewCreateInstanceCommand().BPMNProcessId(factory.Name).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	res, err := req.Send()
	if err != nil {
		return &grpc.FactoryCreateInstanceRs{}, err
	}
	// add to database
	instance := model.Instance{
		Key:         strconv.Itoa(int(res.GetWorkflowInstanceKey())),
		Name:        in.Name,
		Subject:     in.Subject,
		Description: in.Description,
		ContextData: in.ContextData,
	}
	return &grpc.FactoryCreateInstanceRs{
		InstanceKey: strconv.Itoa(int(res.GetWorkflowInstanceKey())),
	}, model.InsertInstance(instance)
}

func (Factory) ListInstances(ctx context.Context, in *grpc.FactoryListInstancesRq) (*grpc.FactoryListInstancesRs, error) {
	instances, err := model.ListInstances()
	if err != nil {
		return &grpc.FactoryListInstancesRs{}, err
	}
	ret := make([]*grpc.Instance, 0)
	for _, v := range instances {
		tmp := grpc.Instance{
			InstanceKey: v.Key,
			Name:        v.Name,
			Subject:     v.Subject,
		}
		ret = append(ret, &tmp)
	}
	return &grpc.FactoryListInstancesRs{
		Sequence: ret,
	}, nil
}

func (Factory) GetDefinition(ctx context.Context, in *grpc.FactoryGetDefinitionRq) (*grpc.FactoryGetDefinitionRs, error) {
	factory, err := model.GetFactory(in.Key)
	if err != nil {
		return &grpc.FactoryGetDefinitionRs{}, err
	}
	return &grpc.FactoryGetDefinitionRs{
		Definition: factory.Definition,
	}, nil
}

func (Factory) SetDefinition(ctx context.Context, in *grpc.FactorySetDefinitionRq) (*grpc.FactorySetDefinitionRs, error) {
	return &grpc.FactorySetDefinitionRs{}, model.UpdateFactory(in.Key, model.Factory{
		Definition: in.Definition,
	})
}
