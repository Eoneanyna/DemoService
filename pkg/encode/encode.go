package encode

import (
	"github.com/go-kratos/kratos/v2/errors"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
	"strings"
)

type Response struct {
	Code    int         `json:"statusCode" form:"statusCode"`
	Message string      `json:"message" form:"message"`
	Data    interface{} `json:"data" form:"data"`
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	reply := Response{
		Code: int(se.Code),
		Message: se.Message,
		Data: nil,
	}

	codec, _ := http2.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType(codec.Name()))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	reply := Response{
		Code: 200,
		Data: v,
		Message: "success",
	}

	codec, _ := http2.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType(codec.Name()))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return nil
}

func contentType(subtype string) string {
	return strings.Join([]string{"application", subtype}, "/")
}