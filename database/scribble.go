package database

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/nanobox-io/golang-scribble"

	"github.com/nanopack/portal/config"
)

type (
	ScribbleDatabase struct {
		scribbleDb *scribble.Driver
	}
)

func (s *ScribbleDatabase) Init() error {
	u, err := url.Parse(config.DatabaseConnection)
	if err != nil {
		return err
	}
	dir := u.Path
	db, err := scribble.New(dir, nil)
	if err != nil {
		return err
	}

	s.scribbleDb = db
	return nil
}

func (s ScribbleDatabase) GetServices() ([]Service, error) {
	var services []Service
	values, err := s.scribbleDb.ReadAll("services")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			err = NoServiceError
		}
		return nil, err
	}
	for i := range values {
		var service Service
		if err = json.Unmarshal([]byte(values[i]), &service); err != nil {
			return nil, fmt.Errorf("Bad JSON syntax received in body")
		}
		services = append(services, service)
	}
	return services, nil
}

func (s ScribbleDatabase) GetService(id string) (*Service, error) {
	service := Service{}
	err := s.scribbleDb.Read("services", id, &service)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			err = NoServiceError
		}
		return nil, err
	}
	return &service, nil
}

func (s ScribbleDatabase) SetServices(services []Service) error {
	// s.scribbleDb.Delete("services", "")
	for i := range services {
		err := s.scribbleDb.Write("services", services[i].Id, services[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ScribbleDatabase) SetService(service *Service) error {
	return s.scribbleDb.Write("services", service.Id, *service)
}

func (s ScribbleDatabase) DeleteService(id string) error {
	return s.scribbleDb.Delete("services", id)
}

func (s ScribbleDatabase) SetServer(svcId string, server *Server) error {
	service, err := s.GetService(svcId)
	if err != nil {
		return err
	}
	service.Servers = append(service.Servers, *server)

	return s.scribbleDb.Write("services", service.Id, service)
}

func (s ScribbleDatabase) SetServers(svcId string, servers []Server) error {
	service, err := s.GetService(svcId)
	if err != nil {
		return err
	}

	// pretty simple, reset all servers
	service.Servers = servers

	return s.scribbleDb.Write("services", service.Id, service)
}

func (s ScribbleDatabase) DeleteServer(svcId, srvId string) error {
	service, err := s.GetService(svcId)
	if err != nil {
		// // if read was successful, but found no
		// if strings.Contains(err.Error(), "found") {
		// 	return nil
		// }
		return nil
	}
	for _, srv := range service.Servers {
		if srv.Id == srvId {
			// todo: empty or a = append(a[:i], a[i+1:]...)
			srv = Server{}
		}
	}

	return s.scribbleDb.Write("services", service.Id, service)
}

func (s ScribbleDatabase) GetServer(svcId, srvId string) (*Server, error) {
	service := Service{}
	err := s.scribbleDb.Read("services", svcId, &service)
	if err != nil {
		return nil, err
	}

	for _, srv := range service.Servers {
		if srv.Id == "srvId" {
			return &srv, nil
		}
	}
	return nil, NoServerError
}
