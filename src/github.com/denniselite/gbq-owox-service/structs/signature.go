package structs

import (
	"encoding/json"
)

type SendRequest struct {
}

type EmptyResponse struct {
}

func (res *EmptyResponse) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}