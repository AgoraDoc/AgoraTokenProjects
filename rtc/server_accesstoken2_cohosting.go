package main

import (
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
)

type rtc_int_token_struct struct{
	Uid_rtc_int uint32 `json:"uid"`
	Channel_name string `json:"ChannelName"`
	Role uint32 `json:"role"`
}

var rtc_token string
var int_uid uint32
var channel_name string

var role_num uint32
var role rtctokenbuilder.Role

// Use RtcTokenBuilder to generate an RTC token.
func generateRtcToken(int_uid uint32, channelName string, role rtctokenbuilder.Role){

	appID := "<Your App ID>"
	appCertificate := "<Your App Cert>"
	// Number of seconds after which the token expires.
	// For demonstration purposes the expiry time is set to 40 seconds. This shows you the automatic token renew actions of the client.
	tokenExpireTimeInSeconds := uint32(40)
	joinChannelPrivilegeExpireTimeInSeconds := uint32(40)
	pubAudioPrivilegeExpireTimeInSeconds := uint32(40)
	pubVideoPrivilegeExpireTimeInSeconds := uint32(40)
    pubDataStreamPrivilegeExpireTimeInSeconds := uint32(40)

	// result, err := rtctokenbuilder.BuildTokenWithUid(appID, appCertificate, channelName, int_uid, role, expireTimeInSeconds)
	result, err := rtctokenbuilder.BuildTokenWithUidAndPrivilege(appID, appCertificate, channelName, int_uid, role, tokenExpireTimeInSeconds,
		joinChannelPrivilegeExpireTimeInSeconds, pubAudioPrivilegeExpireTimeInSeconds, pubVideoPrivilegeExpireTimeInSeconds, pubDataStreamPrivilegeExpireTimeInSeconds)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with uid: %s\n", result)
		fmt.Printf("uid is %d\n", int_uid )
		fmt.Printf("ChannelName is %s\n", channelName)
		fmt.Printf("Role is %d\n", role)
	}
	rtc_token = result
}


func rtcTokenHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
	w.Header().Set("Access-Control-Allow-Headers", "*");

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
		return
	}


	var t_int rtc_int_token_struct
	var unmarshalErr *json.UnmarshalTypeError
	int_decoder := json.NewDecoder(r.Body)
	int_err := int_decoder.Decode(&t_int)
	if (int_err == nil) {

				int_uid = t_int.Uid_rtc_int
				channel_name = t_int.Channel_name
				role_num = t_int.Role
				switch role_num {
					case 1:
						role = rtctokenbuilder.RolePublisher
					case 2:
						role = rtctokenbuilder.RoleSubscriber
				}
	}
	if (int_err != nil) {

		if errors.As(int_err, &unmarshalErr){
				errorResponse(w, "Bad request. Wrong type provided for field " + unmarshalErr.Value  + unmarshalErr.Field + unmarshalErr.Struct, http.StatusBadRequest)
				} else {
				errorResponse(w, "Bad request.", http.StatusBadRequest)
			}
		return
	}

	generateRtcToken(int_uid, channel_name, role)
	errorResponse(w, rtc_token, http.StatusOK)
	log.Println(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}

func main(){
	// Handling routes
	// RTC token from RTC num uid
	http.HandleFunc("/fetch_rtc_token", rtcTokenHandler)
	fmt.Printf("Starting server at port 8082\n")

	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}