// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package gommi

import (
	"log/slog"
	"net/http"
)

// package dash provides a simple dashboard for the job controller
func RouteLanding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobsData := TDLanding{
		Title: "opentsg-studio in Minikube",
		Ptr:   "", // relative path to the root folder nothing
	}

	//render the Jobs
	err := tpl.landing.ExecuteTemplate(w, "page", jobsData)
	if err != nil {
		slog.Error("jobsMain template render error", "err", err)
	}

}
