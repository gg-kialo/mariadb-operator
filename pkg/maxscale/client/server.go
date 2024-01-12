package client

import (
	"context"
	"fmt"

	mdbhttp "github.com/mariadb-operator/mariadb-operator/pkg/http"
)

type ServerParameters struct {
	Address  string `json:"address"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type ServerAttributes struct {
	State      string           `json:"state,omitempty"`
	Parameters ServerParameters `json:"parameters"`
}

type ServerClient struct {
	client *mdbhttp.Client
}

func (s *ServerClient) List(ctx context.Context) ([]Data[ServerAttributes], error) {
	var list List[ServerAttributes]
	res, err := s.client.Get(ctx, "servers", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}
	if err := handleResponse(res, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}

func (s *ServerClient) Create(ctx context.Context, name string, params ServerParameters) error {
	object := &Object[ServerAttributes]{
		Data: Data[ServerAttributes]{
			ID:   name,
			Type: ObjectTypeServers,
			Attributes: ServerAttributes{
				Parameters: params,
			},
		},
	}
	res, err := s.client.Post(ctx, "servers", object, nil)
	if err != nil {
		return err
	}
	return handleResponse(res, nil)
}
