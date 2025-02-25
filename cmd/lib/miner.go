package lib

import (
	"crave/shared/configuration"
)

func Start(container configuration.IContainer) error {
	if err := startApiServer(container); err != nil {
		return err
	}
	if err := startGrpcServer(container); err != nil {
		return err
	}
	return nil
}
func startApiServer(container configuration.IContainer) error {
	container.DefineRoute()
	return nil
}

func startGrpcServer(container configuration.IContainer) error {
	container.DefineGrpc()
	return nil
}
