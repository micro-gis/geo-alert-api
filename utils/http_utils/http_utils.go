package http_utils

import (
	"encoding/json"
	"github.com/micro-gis/oauth-go/oauth"
	"github.com/micro-gis/utils/rest_errors"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, err rest_errors.RestErr) {
	ResponseJson(w, err.Status(), err)
}

func AuthenticateRequest(r *http.Request, forceAuth bool, forceSameUserId int64 ) rest_errors.RestErr {
	if err := oauth.AuthenticateRequest(r); err != nil {
		return err
	}
	// For forcing authentication
	if forceAuth {
		if callerId := oauth.GetCallerId(r); callerId == 0 {
			err := rest_errors.NewRestError("Authentication required", http.StatusUnauthorized, "unauthorized", nil)
			return err
		}
	}

	if forceSameUserId != 0 {
		if callerId := oauth.GetCallerId(r); callerId != forceSameUserId {
			err := rest_errors.NewRestError("Authentication required", http.StatusUnauthorized, "unauthorized", nil)
			return err
		}
	}
	return nil
}
