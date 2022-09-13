package devices

import (
	"fmt"
	"sipsimclient/errors"
	"sipsimclient/model"
	"sipsimclient/repository"

	"github.com/alexeyco/simpletable"
)

type DeviceManager interface {
	Init()
	List() DeviceList
	Get(name string) (Device, error)
	Add(req AddDeviceRequest) error

	Connect(name string) error
	Disconnect(name string) error
	Remove(name string) error
	Send(name string, msg Message) error

	Logs(name string, theme model.Theme) ([]string, error)
}

type DeviceList []Device

func (devices DeviceList) String() string {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "PROTO"},
			{Align: simpletable.AlignCenter, Text: "ADDR"},
			{Align: simpletable.AlignCenter, Text: "STATE"},
		},
	}
	for i := range devices {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i)},
			{Text: devices[i].Name()},
			{Text: string(devices[i].Protocol())},
			{Text: devices[i].Address()},
			{Text: string(devices[i].State())},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompactLite)
	return table.String()
}

type deviceMapManager struct {
	devices map[string]Device
}

func (dm *deviceMapManager) Init() {
}

func (dm *deviceMapManager) List() DeviceList {
	//TODO: implement it
	ans := make([]Device, 0, len(dm.devices))
	for _, v := range dm.devices {
		ans = append(ans, v)
	}
	return ans
}
func (dm *deviceMapManager) Get(name string) (Device, error) {
	device, exists := dm.devices[name]
	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device, nil
}
func (dm *deviceMapManager) Add(req AddDeviceRequest) error {
	//check name duplicate
	_, exists := dm.devices[req.Name]
	if exists {
		return errors.ErrDuplicateDeviceName
	}
	device, err := NewDevice(req)
	if err != nil {
		return err
	}
	dm.devices[req.Name] = device

	return nil
}

func (dm *deviceMapManager) Connect(name string) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Connect()
}
func (dm *deviceMapManager) Disconnect(name string) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Disconnect()
}
func (dm *deviceMapManager) Remove(name string) error {
	device, exists := dm.devices[name]
	if !exists {
		//no such name
		return nil
	}
	if device.State() == DeviceStateOnline {
		err := device.Disconnect()
		if err != nil {
			return err
		}
	}
	delete(dm.devices, name)

	//release logger
	ReleaseLogger(name)

	//data persistence
	err := repository.GetDeviceRepository().Delete(name)
	if err != nil {
		fmt.Println("data persistence failed, err:", err)
	}

	return nil
}
func (dm *deviceMapManager) Send(name string, msg Message) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Send(msg)
}

func (dm *deviceMapManager) Logs(name string, theme model.Theme) ([]string, error) {
	device, exists := dm.devices[name]
	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device.Logs(theme)
}

type AddDeviceRequest struct {
	Name     string
	Password string
	Protocol NetProtocol
}

func NewDeviceManager() DeviceManager {
	persistenceManager := &devicePersistenceManager{
		manager: &deviceMapManager{
			devices: map[string]Device{},
		},
	}
	return persistenceManager
}
