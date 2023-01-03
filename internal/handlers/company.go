package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
)

type CompaniesHandler struct {
	CompanyManager managers.ICompanyManager
}

func GetCompaniesHandler(companyManager managers.ICompanyManager) CompaniesHandler {
	return CompaniesHandler{
		CompanyManager: companyManager,
	}
}

func (handler *CompaniesHandler) Create(ctx *fiber.Ctx) (err error) {

	var form forms.CompanyCreateForm
	if err = ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusBadRequest,
			Data:      nil,
		})
	}

	err = utils.GetCustomValidator().Validate(form)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   utils.FormValidationError,
			Error:     true,
			ErrorCode: http.StatusBadRequest,
			Data:      utils.CustomValidatorErr(err),
		})
	}

	company, err := handler.CompanyManager.Create(form)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    nil,
		})
	}

	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   "",
		Error:     false,
		ErrorCode: http.StatusOK,
		Data:      company,
	})
}

func (handler *CompaniesHandler) Login(ctx *fiber.Ctx) (err error) {
	var form forms.LoginForm
	if err = ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusBadRequest,
			Data:      nil,
		})
	}

	if err = utils.GetCustomValidator().Validate(form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusBadRequest,
			Data:      utils.CustomValidatorErr(err),
		})
	}

	ts, err := handler.CompanyManager.Login(form)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    nil,
		})
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message: "",
		Error:   false,
		Data:    tokens,
	})
}

func (handler *CompaniesHandler) Logout(ctx *fiber.Ctx) (err error) {
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	err = handler.CompanyManager.DeleteTokens(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   "",
		Error:     true,
		ErrorCode: http.StatusOK,
		Data:      nil,
	})
}

func (handler *CompaniesHandler) Refresh(ctx *fiber.Ctx) error {
	mapToken := map[string]string{}
	if err := ctx.BodyParser(&mapToken); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnprocessableEntity,
			Data:      nil,
		})
	}
	refreshToken := mapToken["refresh_token"]
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(models.ResponseModel{
				Message:   err.Error(),
				Error:     true,
				ErrorCode: http.StatusUnprocessableEntity,
				Data:      nil,
			})
		}
		companyID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["company_id"]), 10, 64)
		if err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(models.ResponseModel{
				Message:   err.Error(),
				Error:     true,
				ErrorCode: http.StatusUnprocessableEntity,
				Data:      nil,
			})
		}
		//Delete the previous Refresh Token
		err = handler.CompanyManager.DeleteAuth(refreshUuid)
		if err != nil { //if any goes wrong
			return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
				Message:   err.Error(),
				Error:     true,
				ErrorCode: http.StatusUnauthorized,
				Data:      nil,
			})
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := handler.CompanyManager.CreateToken(uint(companyID), claims)
		if createErr != nil {
			return ctx.Status(http.StatusForbidden).JSON(models.ResponseModel{
				Message:   err.Error(),
				Error:     true,
				ErrorCode: http.StatusForbidden,
				Data:      nil,
			})
		}
		//save the tokens metadata to redis
		saveErr := handler.CompanyManager.CreateAuth(uint(companyID), ts)
		if saveErr != nil {
			return ctx.Status(http.StatusForbidden).JSON(models.ResponseModel{
				Message:   err.Error(),
				Error:     true,
				ErrorCode: http.StatusForbidden,
				Data:      nil,
			})
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return ctx.Status(http.StatusCreated).JSON(models.ResponseModel{
			Message:   "",
			Error:     false,
			ErrorCode: http.StatusCreated,
			Data:      tokens,
		})
	}
	return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
		Message:   err.Error(),
		Error:     true,
		ErrorCode: http.StatusUnauthorized,
		Data:      nil,
	})
}

func ExtractTokenMetadata(ctx *fiber.Ctx) (*models.AccessDetails, error) {
	token, err := VerifyToken(ctx)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		companyID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["company_id"]), 10, 32)
		if err != nil {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			CompanyID:  uint(companyID),
		}, nil
	}
	return nil, err
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func VerifyToken(ctx *fiber.Ctx) (*jwt.Token, error) {
	token_string := ExtractToken(ctx)
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(ctx *fiber.Ctx) string {
	bear_token := ctx.Request().Header.Peek("Authorization")
	strArr := strings.Split(bytes.NewBuffer(bear_token).String(), " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func TokenValid(ctx *fiber.Ctx) error {
	token, err := VerifyToken(ctx)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}
