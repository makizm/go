package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"

	// "github.com/manyminds/api2go/examples/model"
	// "github.com/manyminds/api2go/examples/resource"
	// "github.com/manyminds/api2go/examples/storage"
	"github.com/makizm/ssmartthings/lan-device-handler/routing/model"
	"github.com/makizm/ssmartthings/lan-device-handler/routing/resource"
	"github.com/makizm/ssmartthings/lan-device-handler/routing/storage"
)

// APIServer gin engine
// once initialize start it with .Run()
func APIServer(mode string) *gin.Engine {

	var engine *gin.Engine

	// Set Gin running mode
	if mode == "prod" {
		if gin.Mode() != "release" {
			// Switch to release mode
			gin.SetMode(gin.ReleaseMode)
		}

		// Create production ready instance
		engine = gin.New()
	} else {
		if gin.Mode() != "debug" {
			// Switch to debug mode
			gin.SetMode(gin.DebugMode)
		}

		// Create debug instance with logging and middleware
		engine = gin.Default()
	}

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("/"),
		gingonic.New(engine),
	)

	deviceStorage := storage.NewDeviceStorage()

	// Create some devices, experimental remove later
	d := &model.Device{
		ID:   " ",
		Name: "Device one",
		Type: "Switch",
	}
	deviceStorage.Insert(*d)
	d = &model.Device{
		ID:   " ",
		Name: "Device two",
		Type: "Switch",
	}
	deviceStorage.Insert(*d)
	d = &model.Device{
		ID:   " ",
		Name: "Device three",
		Type: "Switch",
	}
	deviceStorage.Insert(*d)
	d = nil

	api.AddResource(model.Device{}, resource.DeviceResource{DeviceStorage: deviceStorage})

	// Static content
	engine.GET("/status", func(c *gin.Context) {
		c.String(200, "OK")
	})

	return engine
}
