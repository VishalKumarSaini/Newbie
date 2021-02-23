package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"Newbie/db"
	model "Newbie/models"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	otp    string = "0000"
	trials        = 0
)

// SignUpHandler ...
func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	collection, err := db.GetDBCollection("user")

	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: user.Phone}}).Decode(&result)

	if err == nil {
		res.Result = "Phone Number already Registered!!"
		json.NewEncoder(w).Encode(res)
		return
	}

	if err != nil {
		if err.Error() == "mongo: no documents in result" {

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			accountSid := "ACfef2e7e1b15a56f0a07375161ba7e773"
			authToken := "142cdc89a35cc2c42f45b3ebe3ea6a46"
			urlStr := "https://api.twilio.com/2010-04-01/Accounts/ACfef2e7e1b15a56f0a07375161ba7e773/Messages.json"

			max := 9999
			min := 1000
			rand.Seed(time.Now().UnixNano())
			otp = strconv.Itoa(rand.Intn(max-min+1) + min)

			msgData := url.Values{}
			msgData.Set("To", "+917018132601")
			msgData.Set("From", "+14105933183")
			msgData.Set("Body", otp)
			msgDataReader := *strings.NewReader(msgData.Encode())

			client := &http.Client{}
			req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
			req.SetBasicAuth(accountSid, authToken)
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, _ := client.Do(req)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {

				var data map[string]interface{}
				decoder := json.NewDecoder(resp.Body)
				err := decoder.Decode(&data)

				if err == nil {
					fmt.Println(data["sid"])
				}
			} else {
				fmt.Println(resp.Status)
			}
			res.Result = "Phone Authentication Required!"
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	/*res.Error = err.Error()
	json.NewEncoder(w).Encode(res)
	return*/
}

//SignUpAuthHandler ...
func SignUpAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userOtp model.OtpContainer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userOtp)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if userOtp.OtpEntered == otp {
		fmt.Println("The signUp authentication is successful!")
		res.Result = "The signUp authentication is successful!"
		json.NewEncoder(w).Encode(res)
		otp = "0000"

	} else {
		if trials < 3 {
			json.NewEncoder(w).Encode("Please enter correct OTP")
			trials++
			json.NewEncoder(w).Encode("Trials left")
			json.NewEncoder(w).Encode(3 - trials)
		}
		if trials == 3 {
			json.NewEncoder(w).Encode("No more trials left")
			collection, err := db.GetDBCollection("user")
			if err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			_, err = collection.DeleteOne(context.TODO(), bson.M{"phone": userOtp.Number})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Bad OTP !")
			res.Error = "OTP Did not Match!"
			json.NewEncoder(w).Encode(res)
			otp = "0000"
			trials = 0
		}

	}
}

//login handler

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var login model.Login
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &login)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, err := db.GetDBCollection("user")
	var result model.Login
	err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: login.Contact}}).Decode(&result)

	if err != nil {
		res.Result = "Not Registered!"
		json.NewEncoder(w).Encode(res)
		return
	}

	if err == nil {
		res.Result = "Welcome Buddy,Enter Otp!"
		json.NewEncoder(w).Encode(res)
		accountSid := "ACfef2e7e1b15a56f0a07375161ba7e773"
		authToken := "142cdc89a35cc2c42f45b3ebe3ea6a46"
		urlStr := "https://api.twilio.com/2010-04-01/Accounts/ACfef2e7e1b15a56f0a07375161ba7e773/Messages.json"

		max := 9999
		min := 1000
		rand.Seed(time.Now().UnixNano())
		otp = strconv.Itoa(rand.Intn(max-min+1) + min)

		msgData := url.Values{}
		msgData.Set("To", "+917018132601")
		msgData.Set("From", "+14105933183")
		msgData.Set("Body", otp)
		msgDataReader := *strings.NewReader(msgData.Encode())

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := client.Do(req)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {

			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)

			if err == nil {
				fmt.Println(data["sid"])
			}
		} else {
			fmt.Println(resp.Status)
		}

	}
}

//login Auth
func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userOtp model.OtpContainer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userOtp)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if userOtp.OtpEntered == otp {
		fmt.Println("The Login authentication is successful!")
		res.Result = "The Login authentication is successful!"
		json.NewEncoder(w).Encode(res)
		otp = "0000"
	} else {
		if trials < 3 {
			json.NewEncoder(w).Encode("Please enter correct OTP")
			trials++
			json.NewEncoder(w).Encode("Trials left")
			json.NewEncoder(w).Encode(3 - trials)
			//trials++
		}
		if trials == 3 {
			json.NewEncoder(w).Encode("No more trials left")

			if err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			fmt.Println("Bad OTP !")
			res.Error = "OTP Did not Match!"
			json.NewEncoder(w).Encode(res)
			otp = "0000"
			trials = 0

		}

	}
}
