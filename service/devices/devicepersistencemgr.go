package devices

import (
	"sipsimclient/log"
	"sipsimclient/model"
	"sipsimclient/repository"
	"time"
)

type devicePersistenceManager struct {
	manager DeviceManager
}

func (d *devicePersistenceManager) Init() {
	//load from bolt & create them
	devices, err := repository.GetDeviceRepository().All()
	if err != nil {
		log.Warn("can't read devices from repository")
		return
	}

	for i := range devices {
		err := d.manager.Add(AddDeviceRequest{
			Name:     devices[i].Name,
			Password: devices[i].Password,
			Protocol: NetProtocol(devices[i].Protocol),
		})
		if err != nil {
			log.Warnf("create device from repository failed, name: %v", devices[i].Name)
			return
		}

		err = d.manager.Connect(devices[i].Name)
		if err != nil {
			log.Warnf("create connection with device from repository failed, name: %v", devices[i].Name)
			return
		}
	}
}

func (d *devicePersistenceManager) List() DeviceList {
	return d.manager.List()
}
func (d *devicePersistenceManager) Get(name string) (Device, error) {
	return d.manager.Get(name)
}
func (d *devicePersistenceManager) Add(req AddDeviceRequest) error {
	err := d.manager.Add(req)
	if err != nil {
		return err
	}
	//data persistence
	err = repository.GetDeviceRepository().Put(&repository.Device{
		Name:     req.Name,
		Password: req.Password,
		Protocol: string(req.Protocol),
	})
	if err != nil {
		log.Warnf("data persistence failed, err: %v", err)
	}
	return nil
}

func (d *devicePersistenceManager) Update(req UpdateDeviceRequest) error {
	err := d.manager.Update(req)
	if err != nil {
		return err
	}
	err = repository.GetDeviceRepository().Put(&repository.Device{
		Name:     req.Name,
		Password: req.Password,
		Protocol: string(req.Protocol),
	})
	if err != nil {
		log.Warnf("%v data persistence failed, err: %v", req.Name, err)
		return err
	}
	return nil
}

func (d *devicePersistenceManager) Connect(name string) error {
	return d.manager.Connect(name)
}
func (d *devicePersistenceManager) Disconnect(name string) error {
	return d.manager.Disconnect(name)
}
func (d *devicePersistenceManager) Remove(name string) error {
	err := d.manager.Remove(name)
	if err != nil {
		return err
	}
	//data persistence
	err = repository.GetDeviceRepository().Delete(name)
	if err != nil {
		log.Warnf("data persistence failed, err: %v", err)
	}
	return nil
}

func (d *devicePersistenceManager) Send(name string, msg Message) error {
	return d.manager.Send(name, msg)
}
func (d *devicePersistenceManager) DoSend(name string, msgType MessageType, val map[string]string) error {
	return d.manager.DoSend(name, msgType, val)
}

func (d *devicePersistenceManager) Logs(name string, theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error) {
	return d.manager.Logs(name, theme, start, end)
}
