package devices

import (
	"fmt"
	"sipsimclient/errors"

	"github.com/alexeyco/simpletable"
)

type DeviceManager interface {
	List() DeviceList
	Get(name string) (Device, error)
	Add(req AddDeviceRequest) error

	Connect(name string) error
	Disconnect(name string) error
	Remove(name string) error
	Send(name string, msg Message) error

	Logs(name string) ([]string, error)
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

type DeviceMapManager struct {
	devices map[string]Device
}

func (dm *DeviceMapManager) List() DeviceList {
	//TODO: implement it
	ans := make([]Device, 0, len(dm.devices))
	for _, v := range dm.devices {
		ans = append(ans, v)
	}
	return ans
}
func (dm *DeviceMapManager) Get(name string) (Device, error) {
	device, exists := dm.devices[name]
	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device, nil
}
func (dm *DeviceMapManager) Add(req AddDeviceRequest) error {
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

func (dm *DeviceMapManager) Connect(name string) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Connect()
}
func (dm *DeviceMapManager) Disconnect(name string) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Disconnect()
}
func (dm *DeviceMapManager) Remove(name string) error {
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
	return nil
}
func (dm *DeviceMapManager) Send(name string, msg Message) error {
	device, exists := dm.devices[name]
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Send(msg)
}

func (dm *DeviceMapManager) Logs(name string) ([]string, error) {
	device, exists := dm.devices[name]
	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device.Logs()
}

type AddDeviceRequest struct {
	Name     string
	Password string
	Protocol NetProtocol
}

func NewDeviceManager() DeviceManager {
	return &DeviceMapManager{
		devices: map[string]Device{},
	}
}
