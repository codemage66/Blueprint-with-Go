// Package flash adds the flashes to the view template.
package flash

import (
	"fmt"
	"net/http"

	"github.com/blue-jay/blueprint/lib/env"

	flashlib "github.com/blue-jay/core/flash"
	"github.com/blue-jay/core/view"
)

// Service represents the services required for this controller.
type Service struct {
	env.Service
}

// Modify adds the flashes to the view.
func (s Service) Modify(w http.ResponseWriter, r *http.Request, v *view.Info) {
	sess, _ := s.Sess.Instance(r)

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]flashlib.Info, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case flashlib.Info:
				v.Vars["flashes"].([]flashlib.Info)[i] = f.(flashlib.Info)
			default:
				v.Vars["flashes"].([]flashlib.Info)[i] = flashlib.Info{fmt.Sprint(f), flashlib.Standard}
			}

		}
		sess.Save(r, w)
	}
}
