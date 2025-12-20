package follow

import (
	"go-fiber-api/internal/common/response"
	"go-fiber-api/internal/util/token"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FollowController interface {
	FollowMerchant(c *fiber.Ctx) error
	UnfollowMerchant(c *fiber.Ctx) error
	GetMerchantFollowStatus(c *fiber.Ctx) error
}

type followController struct {
	service FollowService
}

func NewFollowController(service FollowService) FollowController {
	return &followController{
		service: service,
	}
}

func (h *followController) FollowMerchant(c *fiber.Ctx) error {
	var req FollowRequest
	id := c.Params("id")

	parsedId, err := uuid.Parse(id)

	if parsedId == uuid.Nil || err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant ID")
	}

	claims := c.Locals("user_id").(*token.CustomClaims).UserID
	req.MerchantID = parsedId
	req.UserID = claims

	followResp, err := h.service.FollowMerchant(&req)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "merchant followed successfully", followResp)
}

func (h *followController) UnfollowMerchant(c *fiber.Ctx) error {
	var req FollowRequest
	id := c.Params("id")

	parsedId, err := uuid.Parse(id)

	if parsedId == uuid.Nil || err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant ID")
	}
	claims := c.Locals("user_id").(*token.CustomClaims).UserID
	req.MerchantID = parsedId
	req.UserID = claims

	unfollowed, err := h.service.UnfollowMerchant(&req)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, "merchant unfollowed successfully", unfollowed)
}

func (h *followController) GetMerchantFollowStatus(c *fiber.Ctx) error {
	var followReq FollowRequest

	id := c.Params("id")

	user_id := c.Locals("user_id").(*token.CustomClaims).UserID
	parsedId, err := uuid.Parse(id)

	if parsedId == uuid.Nil || err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant ID")
	}
	followReq.UserID = user_id
	followReq.MerchantID = parsedId

	followResp, err := h.service.GetMerchantFollowStatus(&followReq)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	log.Println(followResp)

	return response.Success(c, "follow status retrieved successfully", followResp)

}
