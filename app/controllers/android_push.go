package controllers

import (
	"github.com/clifford-un/gorush/gorush"
)

func sendAndroidNotification(token string, message string, title string) {
	gorush.PushConf.Android.APIKey = ADD HERE API KEY
	//gorush.PushConf.Android.APIKey = app.CliFirebaseApiKey

	gorush.PushConf.Android.Enabled = true
	req := gorush.PushNotification{
		Platform: gorush.PlatFormAndroid,
		Message:  message,
		Title:    title,
	}
	req.Tokens = []string{token}

	gorush.PushToAndroid(req)
}
