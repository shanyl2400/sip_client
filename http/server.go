package http

import (
	"fmt"
	"net/http"
	"sipsimclient/config"
	"sipsimclient/devices"
	"sipsimclient/model"
	"time"

	"github.com/gin-gonic/gin"
)

func Start() {
	fmt.Printf("start http service :%v\n", config.Get().HttpServerPort)
	server := NewServer(config.Get().HttpServerPort)
	server.Serve()
}

type Server struct {
	deviceManager devices.DeviceManager
	port          int
}

func (s *Server) Serve() {
	router := gin.Default()
	s.route(router)
	router.Run(fmt.Sprintf(":%v", s.port))
}
func (s *Server) route(r *gin.Engine) {
	//user
	r.POST("/login", s.login)

	//device
	r.GET("/devices", s.mustLogin, s.listDevices)
	r.POST("/device", s.mustLogin, s.addDevice)
	r.PUT("/device/connect/:device_name", s.mustLogin, s.connectDevice)
	r.PUT("/device/disconnect/:device_name", s.mustLogin, s.disconnectDevice)
	r.DELETE("/device/:device_name", s.mustLogin, s.deleteDevice)
	r.POST("/device/send/:device_name", s.mustLogin, s.sendMessage)
	r.GET("/device/logs/:device_name", s.mustLogin, s.deviceLogs)
}

func (s *Server) login(c *gin.Context) {
	//TODO: implement it
}

func (s *Server) mustLogin(c *gin.Context) {
	//TODO: implement it
}

func (s *Server) listDevices(c *gin.Context) {
	deviceList := s.deviceManager.List()
	ans := make([]*model.DeviceBrief, len(deviceList))
	for i := range deviceList {
		ans[i] = &model.DeviceBrief{
			Name:     deviceList[i].Name(),
			Protocol: string(deviceList[i].Protocol()),
			Address:  deviceList[i].Address(),
			State:    string(deviceList[i].State()),
		}
	}
	c.JSON(http.StatusOK, ans)
}

func (s *Server) addDevice(c *gin.Context) {
	req := new(model.AddDeviceRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	err = s.deviceManager.Add(devices.AddDeviceRequest{
		Name:     req.Name,
		Password: req.Password,
		Protocol: devices.NetProtocol(req.Protocol),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func (s *Server) deleteDevice(c *gin.Context) {
	deviceName := c.Param("device_name")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "device_name is empty"})
		return
	}
	err := s.deviceManager.Remove(deviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func (s *Server) connectDevice(c *gin.Context) {
	deviceName := c.Param("device_name")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "device_name is empty"})
		return
	}
	err := s.deviceManager.Connect(deviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func (s *Server) disconnectDevice(c *gin.Context) {
	deviceName := c.Param("device_name")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "device_name is empty"})
		return
	}
	err := s.deviceManager.Disconnect(deviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func (s *Server) deviceLogs(c *gin.Context) {
	deviceName := c.Param("device_name")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "device_name is empty"})
		return
	}
	theme := c.Query("theme")
	if theme == "" {
		theme = model.ThemeAll
	}

	startAtStr := c.Query("start_at")
	endAtStr := c.Query("end_at")

	var err error
	var start, end time.Time
	end, err = time.Parse("2006-01-02 15:04:05", endAtStr)
	if endAtStr == "" || err != nil {
		end = time.Now()
	}
	start, err = time.Parse("2006-01-02 15:04:05", startAtStr)
	if startAtStr == "" || err != nil {
		start = end.Add(-time.Minute * 5)
	}

	logs, err := s.deviceManager.Logs(deviceName, model.Theme(theme), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (s *Server) sendMessage(c *gin.Context) {
	deviceName := c.Param("device_name")
	if deviceName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "device_name is empty"})
		return
	}
	req := new(model.SendMessageRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	err = s.deviceManager.DoSend(deviceName, devices.MessageType(req.MessageType), req.Values)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func NewServer(port int) *Server {
	manager := devices.NewDeviceManager()
	manager.Init()
	return &Server{
		port:          port,
		deviceManager: manager,
	}
}
