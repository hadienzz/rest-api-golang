package token

import "github.com/gofiber/fiber/v2"

func SetAuthToken(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		// IMPORTANT (dev on localhost):
		// - "Secure: true" + SameSite=None membuat cookie DITOLAK
		//   oleh browser ketika berjalan di http://localhost.
		// - Karena frontend (Next.js) dan backend sama-sama di "localhost"
		//   hanya beda port, mereka tetap dianggap same-site, jadi Lax cukup.
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}
