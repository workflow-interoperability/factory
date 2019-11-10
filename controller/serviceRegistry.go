package controller

import (
	"context"

	"github.com/workflow-interoperability/factory/grpc"
	"github.com/workflow-interoperability/factory/lib"
	"github.com/workflow-interoperability/factory/model"
	"github.com/zeebe-io/zeebe/clients/go/pb"
)

type ServiceRegistry struct{}

func (s ServiceRegistry) ListDefinitions(ctx context.Context, in *grpc.ServiceRegistryListDefinitionsRq) (*grpc.ServiceRegistryListDefinitionsRs, error) {
	var ret []*grpc.ServiceRegistryListDefinitionsRsSequence
	factorys, err := model.ListFactories()
	if err != nil {
		return &grpc.ServiceRegistryListDefinitionsRs{}, err
	}
	for _, v := range factorys {
		tmp := grpc.ServiceRegistryListDefinitionsRsSequence{
			DefinitionKey: v.Key,
			Description:   v.Definition,
			Name:          v.Name,
		}
		ret = append(ret, &tmp)
	}
	return &grpc.ServiceRegistryListDefinitionsRs{
		Sequence: ret,
	}, nil
}
func (s ServiceRegistry) NewDefinition(ctx context.Context, in *grpc.ServiceRegistryNewDefinitionRq) (*grpc.ServiceRegistryNewDefinitionRs, error) {
	key := lib.GenerateXID()
	// deploy it first
	client, err := connectZeebe()
	if err != nil {
		return &grpc.ServiceRegistryNewDefinitionRs{}, err
	}
	_, err = (*client).NewDeployWorkflowCommand().AddResource([]byte(in.Definition), key, pb.WorkflowRequestObject_BPMN).Send()
	if err != nil {
		return &grpc.ServiceRegistryNewDefinitionRs{}, err
	}
	// add record to database
	factory := model.Factory{
		Key:        key,
		Language:   in.ProcessLanguage,
		Definition: in.Definition,
	}
	return &grpc.ServiceRegistryNewDefinitionRs{
		Definition: factory.Definition,
	}, model.InsertFactory(factory)
}
func (s ServiceRegistry) GetProperties(ctx context.Context, in *grpc.ServiceRegistryGetPropertiesRq) (*grpc.ServiceRegistryGetPropertiesRs, error) {
	sR, err := model.GetServiceRegistry(in.Key)
	if err != nil {
		return &grpc.ServiceRegistryGetPropertiesRs{}, err
	}
	return &grpc.ServiceRegistryGetPropertiesRs{
		Key:         sR.Key,
		Name:        sR.Name,
		Description: sR.Description,
		Version:     sR.Version,
		Status:      sR.Status,
	}, nil
}
func (s ServiceRegistry) SetProperties(ctx context.Context, in *grpc.ServiceRegistryGetPropertiesRs) (*grpc.ServiceRegistryGetPropertiesRq, error) {
	sR := model.ServiceRegistry{
		Key:         in.Key,
		Name:        in.Name,
		Description: in.Description,
		Version:     in.Version,
		Status:      in.Status,
	}
	return &grpc.ServiceRegistryGetPropertiesRq{}, model.InsertServiceRegistry(sR)
}
