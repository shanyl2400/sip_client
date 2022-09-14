package devices

import (
	"fmt"
	"sipsimclient/errors"
	"sipsimclient/log"
	"sipsimclient/model"
	"sipsimclient/repository"
	"sync"
	"time"

	"github.com/alexeyco/simpletable"
)

type DeviceManager interface {
	Init()
	List() DeviceList
	Get(name string) (Device, error)
	Add(req AddDeviceRequest) error

	Update(req UpdateDeviceRequest) error

	Connect(name string) error
	Disconnect(name string) error
	Remove(name string) error
	Send(name string, msg Message) error
	DoSend(name string, msgType MessageType, val map[string]string) error

	Logs(name string, theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error)
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
	devices     map[string]Device
	deviceMutex sync.Mutex
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
	dm.deviceMutex.Lock()
	defer dm.deviceMutex.Unlock()
	device, exists := dm.devices[name]
	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device, nil
}
func (dm *deviceMapManager) Add(req AddDeviceRequest) error {
	//check name duplicate
	dm.deviceMutex.Lock()
	defer dm.deviceMutex.Unlock()

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
func (dm *deviceMapManager) Update(req UpdateDeviceRequest) error {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[req.Name]
	dm.deviceMutex.Unlock()
	if !exists {
		return errors.ErrDeviceNotExists
	}
	err := device.Update(req.Password, req.Protocol)
	if err != nil {
		return err
	}

	dm.deviceMutex.Lock()
	dm.devices[req.Name] = device
	dm.deviceMutex.Unlock()

	return nil
}

func (dm *deviceMapManager) Connect(name string) error {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[name]
	dm.deviceMutex.Unlock()

	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Connect()
}
func (dm *deviceMapManager) Disconnect(name string) error {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[name]
	dm.deviceMutex.Unlock()

	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Disconnect()
}
func (dm *deviceMapManager) Remove(name string) error {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[name]
	dm.deviceMutex.Unlock()

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
		log.Warnf("data persistence failed, err: %v", err)
	}

	return nil
}

func (dm *deviceMapManager) DoSend(name string, msgType MessageType, val map[string]string) error {
	device, err := dm.Get(name)
	if err != nil {
		log.Warnf("no such device name: %v, err: %v", name, err)
		return err
	}
	msg := createMessage(device, msgType, val)
	if msg == nil {
		log.Infof("invlaid message type: %v", msgType)
		return errors.ErrInvalidMessageType
	}
	return dm.Send(name, msg)
}

func (dm *deviceMapManager) Send(name string, msg Message) error {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[name]
	dm.deviceMutex.Unlock()
	if !exists {
		return errors.ErrDeviceNotExists
	}
	return device.Send(msg)
}

func (dm *deviceMapManager) Logs(name string, theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error) {
	dm.deviceMutex.Lock()
	device, exists := dm.devices[name]
	dm.deviceMutex.Unlock()

	if !exists {
		return nil, errors.ErrDeviceNotExists
	}
	return device.Logs(theme, start, end)
}

type AddDeviceRequest struct {
	Name     string
	Password string
	Protocol NetProtocol
}
type UpdateDeviceRequest struct {
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
