package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/GoogleCloudPlatform/golang-samples/run/helloworld/db/sqlc"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/utils"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type User struct {
	server *Server
}

type UpdateUserParams struct {
	ID        string     `json:"id" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Phone     string    `json:"phone" binding:"required,len=11"`
	Address   string    `json:"address" binding:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserPasswordParams struct {
	ID        string     `json:"id" binding:"required"`
	Password  string    `json:"password" binding:"required,min=8" validate:"passwordStrength"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCodeInput struct {
	UserID string  `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `form:"email"`
}

type UserResponse struct {
	ID         string     `json:"id"`
	Lastname   string    `json:"lastname"`
	Firstname  string    `json:"firstname"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	Email      string    `json:"email"`
	IsLoggedIn bool      `json:"isLoggedIn"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DeleteUserParam struct {
	ID string `json:"id"`
}

func (u User) router(server *Server) {
	u.server = server
	serverGroup := server.router.Group("/users")
	serverGroup.PUT("/update", u.updateUser, AuthenticatedMiddleware())
	serverGroup.PUT("/update_password", u.updatePassword, AuthenticatedMiddleware())
	serverGroup.DELETE("/deactivate", u.deleteUser, AuthenticatedMiddleware())
	serverGroup.GET("/profile", u.userProfile, AuthenticatedMiddleware())
	serverGroup.GET("/get_email", u.getUserEmail)
	serverGroup.GET("/send_code_to_user", u.sendCodetoUser)
	serverGroup.POST("/verify_code", u.verifyCode)
}

//var VerificationCodes = make(map[int64]VerificationCode)

type VerificationResponse struct {
	UserID        string
	GeneratedCode string
	ExpiresAt     time.Duration
	Email         string
}

func extractTokenFromRequest(ctx *gin.Context) (string, error) {
	// Extract the token from the Authorization header
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("unauthorized request")
	}

	// Expecting the header to be in the format "Bearer <token>"
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 && strings.ToLower(headerParts[0]) != "bearer" {
		return "", errors.New("invalid token format")
	}

	return headerParts[1], nil
}

func returnIdRole(tokenString string) (string, string, error) {

	if tokenString == "" {
		return "", "", errors.New("unauthorized: Missing or invalid token")
	}

	userId, role, err := tokenManager.VerifyToken(tokenString)

	if err != nil {
		return "", "", errors.New("failed to verify token")
	}

	return userId, role, nil
}



func (u *User) deleteUser(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		ctx.Abort()
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	id := DeleteUserParam{}

	if userId != id.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	if err := ctx.ShouldBindJSON(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	err = u.server.queries.DeleteUser(context.Background(), id.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "user deactivated sucessfully",
	})
}

func (u *User) updateUser(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	user := UpdateUserParams{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if userId != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	arg := db.UpdateUserParams{
		ID:        user.ID,
		Email:     strings.ToLower(user.Email),
		Phone:     user.Phone,
		Address:   user.Address,
		UpdatedAt: time.Now(),
	}

	userToUpdate, err := u.server.queries.UpdateUser(context.Background(), arg)

	if err != nil {
		handleCreateUserError(ctx, err)
		return
	}

	userResponse := UserResponse{
		ID:        userToUpdate.ID,
		Lastname:  userToUpdate.Lastname,
		Firstname: userToUpdate.Firstname,
		Email:     userToUpdate.Email,
		Phone:     userToUpdate.Phone,
		Address:   userToUpdate.Address,
		CreatedAt: userToUpdate.CreatedAt,
		UpdatedAt: userToUpdate.UpdatedAt,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "user updated successfully",
		"data":    userResponse,
	})
}

func (u *User) userProfile(ctx *gin.Context) {
	value, exist := ctx.Get("id")

	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  exist,
			"message": "Unauthorized",
		})
		return
	}

	userId, ok := value.(string)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  exist,
			"message": "Issue Encountered, try again later",
		})
		return
	}

	user, err := u.server.queries.GetUserById(context.Background(), userId)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "Unauthorized",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Issue Encountered, try again later",
		})
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user fetched successfully",
		"data":    userResponse,
	})
}

func (u *User) getUserEmail(ctx *gin.Context) {

	user := ForgotPasswordInput{}

	if err := ctx.ShouldBindQuery(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(user.Email) == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "no input entered",
		})
		return
	}

	userEmail, err := u.server.queries.GetUserByEmail(context.Background(), strings.ToLower(user.Email))

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"Error":      err.Error(),
			"message":    "The requested user with the specified email does not exist.",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userResponse := UserResponse{
		ID:        userEmail.ID,
		Lastname:  userEmail.Lastname,
		Firstname: userEmail.Firstname,
		Email:     userEmail.Email,
		Phone:     userEmail.Phone,
		Address:   userEmail.Address,
		CreatedAt: userEmail.CreatedAt,
		UpdatedAt: userEmail.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "user retrieved successfully",
		"data":       userResponse,
	})
}

func (u *User) sendCodetoUser(ctx *gin.Context) {
	// Bind User Input for validation

	user := ForgotPasswordInput{}

	if err := ctx.ShouldBindQuery(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(user.Email) == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "no input entered",
		})
		return
	}

	userGot, err := u.server.queries.GetUserByEmail(context.Background(), strings.ToLower(user.Email))

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"Error":      err.Error(),
			"message":    "The requested user with the specified email does not exist.",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// GENERATE THE CODE AND STORE
	codeChan := make(chan string)

	go func(c chan string) {
		//Generate a 4 digit random code
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		code := rng.Intn(9000) + 1000
		c <- fmt.Sprintf("%04d", code)

	}(codeChan)

	returnedCode := <-codeChan

	stringUserId := fmt.Sprintf("%v", userGot.ID)

	timeout := 10 * time.Minute

	//fmt.Println("Did we get here?")

	err = Rdb.Set(ctx, stringUserId, returnedCode, timeout).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"Error":      err.Error(),
		})
		return
	}

	// TODO: Send generated code to the user email address
	var wg sync.WaitGroup

	errorChan := make(chan error)

	wg.Add(1)

	//fmt.Println("About to enter send email goroutine")

	go func(userEmail, code string, e chan<- error) {
		defer wg.Done()

		//fmt.Println("About to read html")
		// filereader, err := os.ReadFile("verification.html")
		// if err != nil {
		// 	e <- err
		// 	ctx.JSON(http.StatusInternalServerError, gin.H{
		// 		"statusCode": http.StatusInternalServerError,
		// 		"Error":      err.Error(),
		// 	})
		// 	ctx.Abort()
		// 	return
		// }

		// messagetoSend := string(filereader)
		// _ = messagetoSend

		//fmt.Println("File converted")

		newmessage := fmt.Sprintf("Hi %v,\n\nWe've received your request for a single-use code to use with your Ra'Nkan account.\n\nYour verification code is: %v,\n\nIf you didn't request this code, you can safely ignore this email. Someone else might have typed your email address by mistake.\nThanks,\nThe Ra'Nkan account team\n", userEmail, code)
		sender := u.server.config2.GoogleUsername
		password := u.server.config2.GooglePassword
		smtpHost := "smtp.gmail.com"
		smtpPort := 587

		message := gomail.NewMessage()
		message.SetHeader("From", sender)
		message.SetHeader("To", userEmail)
		message.SetHeader("Subject", "Verification Code")
		message.SetBody("text/plain", newmessage)
		//message.AddAlternative("text/html", messagetoSend)
		message.Embed("rankan.png")

		// Set up the email server configuration
		dialer := gomail.NewDialer(smtpHost, smtpPort, sender, password)

		//fmt.Println("we got to dialer")

		// Send the email
		if err := dialer.DialAndSend(message); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"Error":      err.Error(),
			})
			e <- err
			return
		}

		//fmt.Println("we sent the mail")

		e <- nil

	}(userGot.Email, returnedCode, errorChan)

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	errVal := <-errorChan

	if errVal != nil {
		ctx.Abort()
		return
	}

	//fmt.Println("Email goroutine ended")

	coderesponse := VerificationResponse{
		UserID:        userGot.ID,
		GeneratedCode: returnedCode,
		ExpiresAt:     timeout,
		Email:         userGot.Email,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "code sent to user successfully",
		"anyError":   errVal,
		"data":       coderesponse,
	})
}

func (u *User) verifyCode(ctx *gin.Context) {

	codeInput := UserCodeInput{}

	if err := ctx.ShouldBindJSON(&codeInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	stringUserId := fmt.Sprintf("%+v", codeInput.UserID)
	storedCode, err := Rdb.Get(ctx, stringUserId).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"error":      err.Error(),
			"message":    "Key does not exist",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"Error":      err.Error(),
		})
		return
	}

	timeout := 10 * time.Minute

	expirationTime := time.Now().Add(timeout)

	// Check if the code expiry is less than 10 min
	if time.Now().After(expirationTime) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "code expired",
		})
		return
	}

	if codeInput.Code != storedCode {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "Invalid verification code",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "code verification successful",
	})

	Rdb.Del(ctx, stringUserId)
}

func (u *User) updatePassword(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	user := UpdateUserPasswordParams{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		stringErr := string(err.Error())
		if strings.Contains(stringErr, "passwordStrength") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Error": `
						"Password must be minimum of 8 characters",
						"Password must be contain at least a number",
						"Password must be contain at least a symbol",
						"Password must be contain a upper case letter"
						`,
			})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if userId != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
	}

	userToUpdatePassword, err := u.server.queries.UpdateUserPassword(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	userResponse := UserResponse{
		ID:        userToUpdatePassword.ID,
		Lastname:  userToUpdatePassword.Lastname,
		Firstname: userToUpdatePassword.Firstname,
		Email:     userToUpdatePassword.Email,
		Phone:     userToUpdatePassword.Phone,
		Address:   userToUpdatePassword.Address,
		CreatedAt: userToUpdatePassword.CreatedAt,
		UpdatedAt: userToUpdatePassword.UpdatedAt,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "password updated successfully",
		"data":    userResponse,
	})
}
