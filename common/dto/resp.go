package dto

import "tic3001-go-server/common/constant"

type ResponseDto struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponseDto(code int, msg string, data interface{}) ResponseDto {
	return ResponseDto{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func GetSuccessRespDto(data interface{}) ResponseDto {
	if data == nil {
		data = ""
	}
	return ResponseDto{
		Code: constant.RespCodeSuccess,
		Msg:  constant.RespMsgSuccess,
		Data: data,
	}
}

func GetServerErrorRespDto() ResponseDto {
	return ResponseDto{
		Code: constant.RespCodeServerError,
		Msg:  constant.RespMsgServerError,
		Data: "",
	}
}

func GetClientParamErrorRespDto(msg string) ResponseDto {
	return ResponseDto{
		Code: constant.RespCodeClientParamError,
		Msg:  msg,
		Data: "",
	}
}
