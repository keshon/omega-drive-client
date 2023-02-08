package rcd

import (
	"app/src/conf"
	"app/src/state"
	"app/src/utils"
	"fmt"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// General rclone response
type input struct {
	Name string `json:"name,omitempty"`
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Input  input  `json:"input,omitempty"`
	Path   string `json:"path,omitempty"`
	Status int    `json:"status,omitempty"`
}

// List of VFS
type VfsesResponse struct {
	VFSES []string `json:"vfses,omitempty"`
}

// Current syncing
type Transferring struct {
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	Bytes      int64   `json:"bytes"`
	Percentage int     `json:"percentage"`
	Speed      float64 `json:"speed,omitempty"`
	SpeedAvg   float64 `json:"speedAvg,omitempty"`
	//ETA        int     `json:"eta,omitempty"`
}
type CoreStatsResponse struct {
	Each []Transferring `json:"transferring,omitempty"`
}

// History transfers
type Transferred struct {
	Name      string `json:"name,omitempty"`
	Size      string `json:"size,omitempty"`
	Bytes     string `json:"bytes,omitempty"`
	Checked   bool   `json:"checked"`
	Timestamp string `json:"timestamp,omitempty"`
	Error     string `json:"error,omitempty"`
	JobID     int    `json:"jobid,omitempty"`
}

type CoreTransferedResponse struct {
	Each []Transferred `json:"transferred,omitempty"` //map[string]interface{}
}

func apiConfigCreate(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `config/create?` + param
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error creating config for via 'config/create': " + resp.Error)
	}
	fmt.Println(resp)
	return resp
}

func apiMountMount(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `mount/mount?` + param
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error mounting via 'mount/mount': " + resp.Error)
	}

	return resp
}

func apiMountUnmount(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `mount/unmount?mountPoint=` + param
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error unmounting letter via 'mount/unmount': " + resp.Error)
	}
	return resp
}

func apiMountUnmountall() Response {
	var resp Response
	path := "mount/unmountall"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error unmounting letter via 'mount/unmount': " + resp.Error)
	}
	return resp
}

func apiVfsList() VfsesResponse {
	var resp VfsesResponse
	path := "vfs/list"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	if len(resp.VFSES) == 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error getting vfs list via 'vfs/list'")
	}
	return resp
}

func apiVfsRefresh(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := "vfs/refresh?fs=" + param
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error refreshing cache directory via 'vfs/refresh': " + resp.Error)
	}
	return resp
}

func apiFscacheClear() Response {
	var resp Response
	path := "fscache/clear"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error clearing cache via 'fscache/clear': " + resp.Error)
	}
	return resp
}

func apiCoreQuit() Response {
	var resp Response
	path := "core/quit"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Error quitting Rclone via 'core/quit': " + resp.Error)
	}
	return resp
}

func apiCoreTransfered() CoreTransferedResponse {
	var resp CoreTransferedResponse
	path := "core/transferred"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

func apiCoreStats() CoreStatsResponse {
	var resp CoreStatsResponse
	path := "core/stats"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

/*
func apiJobList() JobIDs {
	var resp JobIDs
	path := "job/list"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

func apiJobStatus(param int) JobStatusResponse {
	if param == 0 {
		return JobStatusResponse{}
	}
	var resp JobStatusResponse
	path := `job/status?jobid=` + strconv.Itoa(param)
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	if len(resp.Error) > 0 {
		log.Println("[apiJobStatus] Error. Details:")
		log.Println(resp)
	}
	return resp
}
*/
