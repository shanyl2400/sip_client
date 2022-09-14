package api

import (
	"net/http"
	"sipsimclient/model"
	"sipsimclient/service/devices"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) listUsers(c *gin.Context) {
	users, err := s.userService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
func (s *Server) deleteUser(c *gin.Context) {
	userName := c.Param("user_name")
	if userName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "user_name is empty"})
		return
	}
	err := s.userService.Delete(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}
func (s *Server) updateUserPassword(c *gin.Context) {
	req := new(model.UpdatePasswordRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	name := c.MustGet(contextUserNameKey).(string)
	err = s.userService.UpdatePassword(name, req.Password, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}
func (s *Server) resetUserPassword(c *gin.Context) {
	userName := c.Param("user_name")
	if userName == "" {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: "user_name is empty"})
		return
	}
	err := s.userService.ResetPassword(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}
func (s *Server) registerUser(c *gin.Context) {
	req := new(model.RegisterUserRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	err = s.userService.Register(req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: "success"})
}

func (s *Server) login(c *gin.Context) {
	//TODO: implement it
	userName := c.PostForm("username")
	password := c.PostForm("password")

	token, err := s.userService.Login(userName, password)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.HttpErrorResponse{Message: token})
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

func (s *Server) updateDevice(c *gin.Context) {
	req := new(model.UpdateDeviceRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			model.HttpErrorResponse{Message: err.Error()})
		return
	}
	err = s.deviceManager.Update(devices.UpdateDeviceRequest{
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
