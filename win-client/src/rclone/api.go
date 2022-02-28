package rclone

import (
	"app/src/conf"
	"app/src/utils"
	"log"
	"net/http"
)

func APIConfigCreate(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `config/create?` + param
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIConfigCreate] Error creating config for via 'config/create'. Details:")
		log.Println(resp)
	}
	return resp
}

func APIMountMount(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `mount/mount?` + param
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIMountMount] Error mounting via 'mount/mount'. Details:")
		log.Println(resp)
	}
	return resp
}

func APIMountUnmount(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := `mount/unmount?mountPoint=` + param
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIMountUnmount] Error unmounting letter via 'mount/unmount'. Details:")
		log.Println(resp)
	}
	return resp
}

func APIMountUnmountall() Response {
	var resp Response
	path := "mount/unmountall"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIMountUnmountall] Error unmounting letter via 'mount/unmount'. Details:")
		log.Println(resp)
	}
	return resp
}

func APIVfsList() VfsesResponse {
	var resp VfsesResponse
	path := "vfs/list"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	if len(resp.VFSES) == 0 {
		log.Println("[APIVfsList] Error getting vfs list via 'vfs/list'. Details:")
		log.Println(resp)

	}
	return resp
}

func APIVfsRefresh(param string) Response {
	if param == "" {
		return Response{}
	}
	var resp Response
	path := "vfs/refresh?fs=" + param
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIVfsRefresh] Error refreshing cache directory via 'vfs/refresh'. Details:")
		log.Println(resp)
	}
	return resp
}

func APIFscacheClear() Response {
	var resp Response
	path := "fscache/clear"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[APIFscacheClear] Error clearing cache via 'fscache/clear'. Details:")
		log.Println(resp)
	}
	return resp
}

func APICoreQuit() Response {
	var resp Response
	path := "core/quit"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)
	if len(resp.Error) > 0 {
		log.Println("[Quit] Error quitting Rclone via 'core/quit'. Details:")
		log.Println(resp)
	}
	return resp
}

func APICoreTransfered() CoreTransferedResponse {
	var resp CoreTransferedResponse
	path := "core/transferred"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

func APICoreStats() CoreStatsResponse {
	var resp CoreStatsResponse
	path := "core/stats"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

/*
func APIJobList() JobIDs {
	var resp JobIDs
	path := "job/list"
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	return resp
}

func APIJobStatus(param int) JobStatusResponse {
	if param == 0 {
		return JobStatusResponse{}
	}
	var resp JobStatusResponse
	path := `job/status?jobid=` + strconv.Itoa(param)
	utils.Request(conf.RcAuthEncoded, http.MethodPost, conf.RcHost+path, nil, &resp)

	if len(resp.Error) > 0 {
		log.Println("[APIJobStatus] Error. Details:")
		log.Println(resp)
	}
	return resp
}
*/
