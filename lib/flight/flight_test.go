// Package flight_test
package flight_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/blue-jay/blueprint/lib/flight"
	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/jsonconfig"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
	"github.com/blue-jay/core/xsrf"
)

// Info contains the application settings.
type Info struct {
	Session  session.Info  `json:"Session"`
	Template view.Template `json:"Template"`
	View     view.Info     `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// TestRace tests for race conditions.
func TestRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			// Configuration
			config := &Info{}

			// Load the configuration file
			err := jsonconfig.Load("../../env.json", config)
			if err != nil {
				t.Error("Could not load: env.json")
			}

			// Set up the session cookie store
			session.SetConfig(config.Session)

			// Set up the views
			config.View.SetTemplates(config.Template.Root, config.Template.Children)

			// Store the view in flight
			flight.StoreConfig(
				asset.Info{},
				form.Info{},
				config.View,
				xsrf.Info{},
				nil)

			// Test the context retrieval
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost/foo", nil)
			c := flight.Context(w, r)

			c.Asset.Folder = "test"
			log.Println(c.Asset.Folder)

			c.View.BaseURI = "monkey"
			log.Println(c.View.BaseURI)
		}()
	}
}
