package auth

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/CleverseAcademy/cd-compose-deployment/api/dto"
	"github.com/CleverseAcademy/cd-compose-deployment/api/services"
	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var PublicKey *ecdsa.PublicKey

func init() {
	pKey, err := jwt.ParseECPublicKeyFromPEM(config.AppConfig.PublicKeyPEMBytes)
	if err != nil {
		panic(err)
	}

	PublicKey = pKey
}

type IArgsCreateSignatureVerificationMiddleware struct {
	services.IService
}

func SignatureVerificationMiddleware(args IArgsCreateSignatureVerificationMiddleware) fiber.Handler {
	return func(c *fiber.Ctx) error {
		signature := c.Get("Authorization")
		data, err := jwt.ParseWithClaims(signature, &dto.SignatureClaims{}, func(t *jwt.Token) (interface{}, error) {
			return PublicKey, nil
		})
		if err != nil {
			return errors.Wrap(err, "SignatureVerification")
		}

		claims, ok := data.Claims.(*dto.SignatureClaims)
		if !ok {
			return fiber.NewError(fiber.StatusPreconditionFailed, fmt.Sprintf("Given claims is of type %s, not SignatureClaims", reflect.TypeOf(data).String()))
		}

		if claims.ExpiresAt == nil || claims.NotBefore == nil {
			return fiber.NewError(fiber.StatusPreconditionRequired, "exp and nbf must be defined")
		}

		if claims.ExpiresAt.Sub(claims.NotBefore.Time) > config.AppConfig.TokenWindow {
			return fiber.NewError(fiber.StatusPreconditionFailed, "lifetime of token is too long")
		}

		checksumHex := fmt.Sprintf("%x", sha256.Sum256(c.Body()))
		if checksumHex != claims.PayloadChecksum {
			return fiber.NewError(fiber.StatusPreconditionFailed, "checksum mismatch")
		}

		request := new(dto.DeployImageDto)

		err = c.BodyParser(request)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}

		expectedJti, err := args.GetNextJTI(request.Service)
		if err != nil {
			fmt.Println(err)
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		if expectedJti != claims.ID {
			return fiber.NewError(fiber.StatusFailedDependency, fmt.Sprintf("JTI mismatched (get an updated one by GET %s)", constants.PathGetDeploymentJTI))
		}

		return c.Next()
	}
}
