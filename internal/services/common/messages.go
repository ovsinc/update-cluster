package common

type CommonResp struct {
	OK     bool `json:"ok"`
	Errors map[string]string
}

type HelloRequest struct {
	Msg string `json:"msg"`
}

func (hr HelloRequest) Sanitize() HelloRequest {
	if Policy == nil {
		return hr
	}

	return HelloRequest{
		Msg: Policy.Sanitize(hr.Msg),
	}
}

type HelloResponse struct {
	CommonResp
	Msg string `json:"msg"`
}
