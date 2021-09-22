package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"plutus/config"
	"time"
)

// namespacesHandler handles the REST API endpoint /namespaces
func (a *APIServer) namespacesHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Handling request for namespaces")

	namespaces, err := config.GetNamespaces()

	if err != nil {
		respondWithError(w, 500, err)
		return
	}

	respondWithJSON(w, 200, map[string][]string{"namespaces": namespaces})
}

// lastRefreshHandler handles the REST API endpoint /lastRefresh
func (a *APIServer) lastRefreshHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Handling request for lastRefresh")
	duration := time.Since(a.lastRefreshed)
	respondWithJSON(w, 200, map[string]int{"passedTime": int(math.Ceil(duration.Minutes()))})
}

// namespacesHandler handles the REST API endpoint /ui
// redirects to plutus-ui
func (a *APIServer) uiHandler(w http.ResponseWriter, req *http.Request) {

	var redirectURL string
	if req.URL.RawQuery != "" {
		redirectURL = fmt.Sprintf("%s&%s", a.UIAddress, req.URL.RawQuery)
	} else {
		redirectURL = fmt.Sprintf("%s", a.UIAddress)
	}
	http.Redirect(w, req, redirectURL, http.StatusSeeOther)
	return
}

// pathHandler handles the REST API endpoint /path
func (a *APIServer) pathHandler(w http.ResponseWriter, r *http.Request) {
	_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithErrorMsg(w, 400, "Bad request. Problem reading body of the request.")
		return
	}

	var body pathRequestBody
	err = json.Unmarshal(_body, &body)
	if err != nil {
		respondWithErrorMsg(w, 400, "Bad request. Please ensure that the path field is specified in the body.")
		return
	}

	path := body.Path
	namespace := body.Namespace
	a.logger.Info("Handling request for path ", path)
	payload := make(map[string]interface{})

	usersSet, err := a.redisClient.QueryUsersFromVaultPath(namespace, path)
	if err != nil {
		respondWithError(w, 500, err)
		return
	}

	payload["users"] = usersSet.AsSlice()
	respondWithJSON(w, 200, payload)
}

// userHandler handles the REST API endpoint /user
func (a *APIServer) userHandler(w http.ResponseWriter, r *http.Request) {
	_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithErrorMsg(w, 400, "Bad request. Problem reading body of the request.")
		return
	}

	var body userRequestBody
	err = json.Unmarshal(_body, &body)
	if err != nil {
		respondWithErrorMsg(w, 400, "Bad request. Please ensure that the path field is specified in the body.")
		return
	}
	username := body.Username
	namespace := body.Namespace

	a.logger.Info("Handling request for user ", username)
	payload := make(map[string]interface{})

	policies, err := a.redisClient.QueryPoliciesFromUsername(namespace, username)
	if err != nil {
		respondWithError(w, 500, err)
		return
	}
	payload["policies"] = policies.AsSlice()

	roles, err := a.redisClient.QueryRolesFromUsername(namespace, username)
	if err != nil {
		respondWithError(w, 500, err)
		return
	}
	payload["roles"] = roles

	groups, err := a.redisClient.QueryGroupsFromUsername(namespace, username)
	if err != nil {
		respondWithError(w, 500, err)
		return
	}
	payload["groups"] = groups

	paths, err := a.redisClient.QueryPathsFromUsername(namespace, username)
	if err != nil {
		respondWithError(w, 500, err)
		return
	}
	payload["paths"] = paths.AsSlice()
	respondWithJSON(w, 200, payload)
}
