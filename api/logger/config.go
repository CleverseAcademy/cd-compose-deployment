package logger

import (
	"io"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pkg/errors"
)

func NewAccessLogMiddleware(o io.Writer) fiber.Handler {
	return logger.New(logger.Config{
		Output:     o,
		Format:     "${time} ${method} ${path} ${status} ${token_chksm} ${body}",
		TimeFormat: time.RFC3339,
		CustomTags: map[string]logger.LogFunc{
			"token_chksm": func(output logger.Buffer, c *fiber.Ctx, _ *logger.Data, _ string) (int, error) {
				token := c.Get("Authorization")

				chksm, err := utils.Base64EncodedSha256(token)
				if err != nil {
					return 0, errors.Wrap(err, "CustomTags.token_chksm")
				}

				return output.WriteString(chksm)
			},
		},
	})
}
