package services

import (
	"context"
	"errors"
	"time"

	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func SeedAdminUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repositories.GetUserByEmail(ctx, "admin@dexra.ai")
	if err == nil {
		return nil // User already exists
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	admin := &models.User{
		Email:        "admin@dexra.ai",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	return repositories.CreateUser(ctx, admin)
}

func Login(email, password string) (string, string, *models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := repositories.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	accessToken, err := generateToken(user.ID.Hex(), time.Hour*24)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := generateToken(user.ID.Hex(), time.Hour*24*7)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func GoogleLogin(ctx context.Context, idToken string) (string, string, *models.User, error) {
	// Verify the ID token
	payload, err := idtoken.Validate(ctx, idToken, config.AppConfig.GoogleClientID)
	if err != nil {
		return "", "", nil, errors.New("invalid google token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return "", "", nil, errors.New("email not provided by google")
	}

	// Check if user exists
	user, err := repositories.GetUserByEmail(ctx, email)
	if err != nil {
		// User doesn't exist, create them
		user = &models.User{
			Email:        email,
			PasswordHash: "", // No password for Google OAuth users
			CreatedAt:    time.Now(),
		}
		if err := repositories.CreateUser(ctx, user); err != nil {
			return "", "", nil, errors.New("failed to create user account")
		}
		// Retrieve again to get the generated MongoDB ID
		user, _ = repositories.GetUserByEmail(ctx, email)
	}

	accessToken, err := generateToken(user.ID.Hex(), time.Hour*24)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := generateToken(user.ID.Hex(), time.Hour*24*7)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func GoogleLoginCodeFlow(ctx context.Context, oauthConfig *oauth2.Config, code string) (string, string, *models.User, error) {
	// Exchange code for token
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", nil, errors.New("failed to exchange token: " + err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", "", nil, errors.New("no id_token field in oauth2 token")
	}

	// Verify the ID token
	payload, err := idtoken.Validate(ctx, rawIDToken, config.AppConfig.GoogleClientID)
	if err != nil {
		return "", "", nil, errors.New("invalid google token: " + err.Error())
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return "", "", nil, errors.New("email not provided by google")
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	// Check if user exists
	user, err := repositories.GetUserByEmail(ctx, email)
	if err != nil {
		// User doesn't exist, create them
		user = &models.User{
			Email:        email,
			Name:         name,
			Picture:      picture,
			PasswordHash: "", // No password for Google OAuth users
			CreatedAt:    time.Now(),
		}
		if err := repositories.CreateUser(ctx, user); err != nil {
			return "", "", nil, errors.New("failed to create user account")
		}
		// Retrieve again to get the generated MongoDB ID
		user, _ = repositories.GetUserByEmail(ctx, email)
	} else {
		// Optionally update the name/picture if it changed, but we can just use the existing for now
		// Or we can save it to the DB. For brevity, we just return the existing user.
		// However, it's nice to keep their latest profile picture synced:
		if user.Picture != picture || user.Name != name {
			user.Name = name
			user.Picture = picture
			repositories.UpdateUser(ctx, user) // we'll need to create this method
		}
	}

	accessToken, err := generateToken(user.ID.Hex(), time.Hour*24)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := generateToken(user.ID.Hex(), time.Hour*24*7)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func generateToken(userID string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
