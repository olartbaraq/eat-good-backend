/*
 *   Copyright (c) 2023 Mubaraq Akanbi
 *   All rights reserved.
 *   Created by Mubaraq Akanbi
 */
 package api

 import (
	 "log"
	 "net/http"
	 "net/url"
	 "strconv"
	 "strings"
	 "sync"
	 "time"
	 "unicode"
	 "unicode/utf8"
 
	 "github.com/go-playground/validator/v10"
 )
 
 // ValidatePassword checks if the password meets the specified criteria.
 var ValidatePassword validator.Func = func(fl validator.FieldLevel) bool {
	 password := fl.Field().Interface().(string)
 
	 // Check if the password is at least 8 characters long
	 if utf8.RuneCountInString(password) < 8 {
		 return false
	 }
 
	 // Check if the password contains at least one digit and one symbol
	 hasDigit := false
	 hasSymbol := false
	 hasUpper := false
	 hasLower := false
	 for _, char := range password {
		 if unicode.IsNumber(char) {
			 hasDigit = true
		 }
		 if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			 hasSymbol = true
		 }
		 if unicode.IsUpper(char) {
			 hasUpper = true
		 }
		 if unicode.IsLower(char) {
			 hasLower = true
		 }
	 }
 
	 //fmt.Println("Validating password:", password)
 
	 return hasDigit && hasSymbol && hasUpper && hasLower
 }
 
 // ImageURLValidation is a custom validator function to check if the URL points to an image.
 var ImageURLValidation validator.Func = func(fl validator.FieldLevel) bool {
 
	 urlStrArray := fl.Field().Interface().([]string)
 
	 // Create a channel to receive results from goroutines
	 imgCh := make(chan bool, len(urlStrArray))
	 //defer close(imgCh)
 
	 // Use a WaitGroup to wait for all goroutines to finish
	 var wg sync.WaitGroup
 
	 for _, urlStr := range urlStrArray {
		 wg.Add(1)
		 go isImageURL(urlStr, imgCh, &wg)
	 }
 
	 // Close the channel when all goroutines finish
	 go func() {
		 wg.Wait()
		 close(imgCh)
	 }()
 
	 // Check the results from the channel
	 for isImage := range imgCh {
		 if !isImage {
			 return false
		 }
	 }
 
	 return true
 }
 
 func isImageURL(urlStr string, ch chan<- bool, wg *sync.WaitGroup) {
 
	 defer wg.Done()
 
	 u, err := url.Parse(urlStr)
	 if err != nil || u.Scheme == "" || u.Host == "" {
		 ch <- false
		 return
	 }
 
	 client := http.Client{
		 Timeout: 5 * time.Second, // Set a timeout for the HTTP request
	 }
 
	 resp, err := client.Get(u.String())
	 if err != nil {
		 // Handle different types of errors (e.g., network error, timeout)
		 // Later
		 ch <- false
		 return
	 }
	 defer resp.Body.Close()
 
	 // Check if the content type indicates an image
	 contentType := resp.Header.Get("Content-Type")
	 isImage := strings.HasPrefix(contentType, "image/")
 
	 // Check if the image is not more than 500KB
	 isNotMoreThan500Kb := false
	 contentLength := resp.Header.Get("Content-Length")
	 if len(contentLength) > 0 {
		 length, err := strconv.ParseInt(contentLength, 10, 64)
		 if err == nil && length <= 500*1024 { // 500 KB in bytes
			 isNotMoreThan500Kb = true
		 }
	 }
 
	 ch <- isImage && isNotMoreThan500Kb
 }
 
 var PriceValidation validator.Func = func(fl validator.FieldLevel) bool {
	 price := fl.Field().Interface().(string)
 
	 priceFloat, err := strconv.ParseFloat(price, 64)
	 if err != nil {
		 log.Fatal(err.Error())
	 }
 
	 if priceFloat < 0 {
		 return false
	 }
	 return true
 }
 