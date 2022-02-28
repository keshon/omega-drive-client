package rclone

type Params map[string]interface{}

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

// Completed transfers
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

// Current transfers
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

/*
type JobIDs struct {
	Each []int `json:"jobids,omitempty"`
}

type JobStatusResponse struct {
	ID        int       `json:"id,omitempty"`
	Group     string    `json:"group,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Error     string    `json:"error,omitempty"`
	Finished  bool      `json:"finished,omitempty"`
	Success   bool      `json:"success,omitempty"`
	Duration  float64   `json:"duration,omitempty"`
	Output    Params    `json:"output,omitempty"`
}
*/
