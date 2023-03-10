package controller

//개인기록 요청받아 서비스에서 결과를 받아와 응답주는 처리

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kyungmun/otp-server/middleware"
	"github.com/kyungmun/otp-server/models"
	"github.com/kyungmun/otp-server/service"
)

type OtpController struct {
	svc *service.OtpServices
}

func NewOtpController(s *service.OtpServices) *OtpController {
	return &OtpController{svc: s}
}

func (c *OtpController) GetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		pageSize = 0
	}
	personalRecords, err := c.svc.GetAll(page, pageSize)
	if err != nil {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": err.Error(),
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "fiber engine, otp registry all get successfully",
		"count":   len(*personalRecords),
		"data":    personalRecords,
		"result":  true,
	})
	return nil
}

func (c *OtpController) GetRecordByID(ctx *fiber.Ctx) error {
	otp_id := ctx.Params("otp_id")
	if otp_id == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "otp_id cannot be empty",
			"result":  false,
		})
		return nil
	}

	otpRegistry, err := c.svc.GetRecordByID(otp_id)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
			"otp_id":  otp_id,
			"result":  false,
		})
		return nil
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "record id gotten successfully",
		"data":    otpRegistry,
	})
	return nil
}

func (c *OtpController) OtpVerify(ctx *fiber.Ctx) error {
	otp_id := ctx.Params("otp_id")
	otp_num := ctx.Query("otp_num", "0")

	if otp_num == "0" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "required query name otp_num",
			"otp_id":  otp_id,
			"result":  false,
		})
		return nil
	}

	_, err := c.svc.GetRecordByID(otp_id)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
			"otp_id":  otp_id,
			"result":  false,
		})
		return nil
	}

	//OTP 검증
	result, otp_real_num := c.svc.OtpVerify(otp_id, otp_num)

	if result {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message":       "otp verify successfully",
			"otp_id":        otp_id,
			"otp_input_num": otp_num,
			"otp_real_num":  otp_real_num,
			"result":        true,
		})
	} else {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message":       "otp verify fail",
			"otp_id":        otp_id,
			"otp_input_num": otp_num,
			"otp_real_num":  otp_real_num,
			"result":        false,
		})
	}

	return nil
}

func (c *OtpController) UpdateRecord(ctx *fiber.Ctx) error {
	otpRegistry := &models.OtpRegistry{}

	err := ctx.BodyParser(&otpRegistry)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": err.Error(),
			"result":  false,
		})
		return err
	}
	fmt.Println(otpRegistry)

	otpRegistryNew, err := c.svc.UpdateRecord(otpRegistry)

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "otp registry update has been successfully",
		"data":    otpRegistryNew,
		"result":  true,
	})

	return nil
}

func (c *OtpController) PatchRecord(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	if userId == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
			"result":  false,
		})
		return nil
	}

	var jsonMap map[string]interface{}
	err := ctx.BodyParser(&jsonMap)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	fmt.Println(jsonMap)

	otpRegistryNew, err := c.svc.PatchRecord(userId, jsonMap)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update otp registry",
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "otp registry update has been successfully",
		"data":    otpRegistryNew,
		"result":  true,
	})

	return nil
}

func (c *OtpController) DeleteRecord(ctx *fiber.Ctx) error {
	otp_id := ctx.Params("otp_id")
	log.Printf("param otp_id : %s", otp_id)
	if otp_id == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "otp_id cannot be empty",
			"result":  false,
		})
		return nil
	}

	_, err := c.svc.GetRecordByID(otp_id)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
			"otp_id":  otp_id,
			"result":  false,
		})
		return nil
	}

	err = c.svc.DeleteRecord(otp_id)
	if err != nil {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": err.Error(),
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": fmt.Sprintf("otp id [%s] delete successfully", otp_id),
		"result":  true,
	})

	return nil
}

func (c *OtpController) CreateRecord(ctx *fiber.Ctx) error {
	otpRegistry := &models.OtpRegistry{}

	err := ctx.BodyParser(&otpRegistry)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "request failed",
			"result":  false,
		})
		return err
	}

	otpRecord, _ := c.svc.GetRecordByID(otpRegistry.OtpID)
	if otpRecord != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "exists otp_id",
			"result":  false,
		})
		return err
	}

	validator := validator.New()
	err = validator.Struct(otpRegistry)

	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": err,
			"result":  false,
		})
		return err
	}

	otpRegistryNew, err := c.svc.CreateRecord(otpRegistry)

	if err != nil {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "could not create otp registry",
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "otp registry Create has been successfully",
		"data":    otpRegistryNew,
		"result":  true,
	})

	return nil
}

func (c *OtpController) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/otp", middleware.Protected(), c.CreateRecord)
	api.Get("/otp", c.GetAll)
	api.Get("/otp/:otp_id", c.GetRecordByID)
	api.Get("/otp/verify/:otp_id", c.OtpVerify)
	api.Put("/otp/:otp_id", middleware.Protected(), c.UpdateRecord)
	api.Patch("/otp/:otp_id", middleware.Protected(), c.PatchRecord)
	api.Delete("otp/:otp_id", middleware.Protected(), c.DeleteRecord)
}
