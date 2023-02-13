package controller

// 각 기능 HTTP API 요청을 Fiber 프레임워크를 사용해서 서비스와 연결 설정

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/kyungmun/otp-server/repository"
	"github.com/kyungmun/otp-server/service"
)

type FiberHendler struct {
	App *fiber.App
}

func NewFiber() *FiberHendler {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	return &FiberHendler{App: app}
}

func (f *FiberHendler) Listen(port string) {
	err := f.App.Listen(port)
	if err != nil {
		log.Panic("service start fail")
	}
}

func (f *FiberHendler) SetupRoutes3(svc *service.ServiceRepository) {
	//OTP 정보 관리
	otpRepo, err := repository.NewOtpRepository(svc.DB)
	if err != nil {
		log.Fatal("could not otp repository create")
	}

	otpService, err := service.NewOtpServices(otpRepo)
	if err != nil {
		log.Fatal("could not otp services create")
	}

	otpController := NewOtpController(otpService)

	otpController.SetupRoutes(f.App)

}

func (f *FiberHendler) SetupRoutes2(repo *repository.OtpRepository) {
	//OTP 정보 관리

	//서비스에 레포지토리 주입
	otpService, err := service.NewOtpServices(repo)
	if err != nil {
		log.Fatal("could not otp services create")
	}

	//컨트롤러에 서비스 주입
	otpController := NewOtpController(otpService)

	//컨트롤러 라우팅 설정
	otpController.SetupRoutes(f.App)
}

func (f *FiberHendler) SetupOtpRoutes(svc *service.OtpServices) {
	//OTP 정보 관리

	//컨트롤러에 서비스 주입
	otpController := NewOtpController(svc)

	//컨트롤러 라우팅 설정
	otpController.SetupRoutes(f.App)
}

func (f *FiberHendler) SetupUserRoutes(svc *service.UserServices) {
	//OTP 정보 관리

	//컨트롤러에 서비스 주입
	userController := NewUserController(svc)

	//컨트롤러 라우팅 설정
	userController.SetupRoutes(f.App)
}

func (f *FiberHendler) SetupAuthRoutes(svc *service.UserServices) {
	//OTP 정보 관리

	//컨트롤러에 서비스 주입
	authController := NewAuthController(svc)

	//컨트롤러 라우팅 설정
	authController.SetupRoutes(f.App)
}
