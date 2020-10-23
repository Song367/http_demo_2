package Services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func MyLimit(){

}

func GetUserInfoRequest(c context.Context,request *http.Request,u interface{})error{
	req:=u.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(req.UserId)
	return nil
}

func GetUserInfoResponse(c context.Context,res *http.Response) (response interface{}, err error){
	if res.StatusCode>400{
		return nil,errors.New("no data")
	}
	var UserRes UserResponse
	err=json.NewDecoder(res.Body).Decode(&UserRes)
	if err!=nil{
		return nil, err
	}
	return UserRes,nil
}