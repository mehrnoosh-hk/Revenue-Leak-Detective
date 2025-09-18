package handlers

import (
	"log/slog"
	"net/http"
	"rdl-api/config"
)

// VersionResponse represents the version information response
type VersionResponse struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildDate   string `json:"build_date"`
	GoVersion   string `json:"go_version,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// VersionHandler returns version information
func VersionHandler(buildInfo *config.BuildInfoConfig, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := VersionResponse{
			Version:   buildInfo.GIT_TAG,
			Commit:    buildInfo.GIT_COMMIT_HASH,
			BuildDate: buildInfo.BUILD_TIMESTAMP,
		}

		WriteJSONSuccessResponse(r.Context(), w, slog.Default(), response)
	}

}
