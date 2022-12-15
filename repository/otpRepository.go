package repository

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kyungmun/otp-server/models"
	"gorm.io/gorm"
)

func NewOtpRepository(db *gorm.DB) (*OtpRepository, error) {
	return &OtpRepository{db: db}, nil
}

type OtpRepository struct {
	db *gorm.DB
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 2
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *OtpRepository) GetAll(page, pageSize int) (*[]models.OtpRegistry, error) {
	otpRegistrys := &[]models.OtpRegistry{}

	if (page <= 0) && (pageSize <= 0) {
		err := r.db.Find(otpRegistrys).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Scopes(Paginate(page, pageSize)).Find(otpRegistrys).Error
		if err != nil {
			return nil, err
		}
	}

	return otpRegistrys, nil
}

func (r *OtpRepository) UpdateRecord(otpRegistry *models.OtpRegistry) (*models.OtpRegistry, error) {

	fmt.Println(otpRegistry)

	//grom 에서는 구조체로 저장시 0 값인 필드는 업데이트 하지 않는다., 맵으로 변환하여 저장해야함.
	//Put 메소드로 Update 에서는 전체 필드를 넘겨 받아서 처리한다. 단, 값이 숫자 필드가 0인 값이 있다면 맵으로 변환해서 처리 해야함.
	//err := r.db.Model(&otpRegistry).Where("m_id = ?", otpRegistry.UserID).Updates(otpRegistry).Error

	//맵으로 변경해야할 값을 넘기면 0 값도 저장 됨.
	//err := r.db.Model(&otpRegistry).Where("m_id = ?", otpRegistry.UserID).Updates(jsonData).Error

	//select 해서 변경할 필드를 지정 가능. 이렇게하면 0 도 저장 됨.
	err := r.db.Model(&otpRegistry).Select("*").Where("otp_id = ?", otpRegistry.OtpID).Updates(otpRegistry).Error
	if err != nil {
		return nil, err
	}

	return otpRegistry, nil
}

func (r *OtpRepository) PatchRecord(otp_id string, jsonData map[string]interface{}) (*models.OtpRegistry, error) {

	fmt.Println(jsonData)

	otpRegistry := &models.OtpRegistry{}

	//patch update 일때는 맵으로 전달된 필드만 업데이트 한다.
	err := r.db.Model(&otpRegistry).Where("otp_id = ?", otp_id).Updates(jsonData).Error
	if err != nil {
		return nil, err
	}

	err2 := r.db.Where("otp_id = ?", otp_id).First(otpRegistry).Error
	if err2 != nil {
		log.Printf("%v", err2)
		return nil, err2
	}

	return otpRegistry, nil
}

func (r *OtpRepository) DeleteByID(otp_id string) error {

	otpRegistry := &models.OtpRegistry{}

	result := r.db.Where("otp_id = ?", otp_id).First(otpRegistry)
	if result.RowsAffected == 0 { // returns count of records found
		fmt.Printf("count : 0")
		return fiber.NewError(fiber.StatusNotFound, "No Record Found")
	}

	err := r.db.Where("otp_id = ?", otp_id).Delete(otpRegistry).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OtpRepository) GetByID(otp_id string) (*models.OtpRegistry, error) {
	//fmt.Println(">> ", m.DBEngine)

	otpRegistry := &models.OtpRegistry{}
	//id := context.Params("id")

	err := r.db.Where("otp_id = ?", otp_id).First(otpRegistry).Error
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return otpRegistry, nil
}

func (r *OtpRepository) Create(otpRegistry *models.OtpRegistry) (*models.OtpRegistry, error) {

	err := r.db.Create(&otpRegistry).Error

	if err != nil {
		return nil, err
	}

	return otpRegistry, nil
}
