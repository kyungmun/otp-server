package repository

import (
	"fmt"
	"log"

	"github.com/kyungmun/otp-server/models"
	"gorm.io/gorm"
)

func NewOtpRepository(db *gorm.DB) (*OtpRepository, error) {
	return &OtpRepository{db: db}, nil
}

type OtpRepository struct {
	db *gorm.DB
}

func (r *OtpRepository) paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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

func (r *OtpRepository) GetAll(page, pageSize int) (*[]models.Otp, error) {
	Otps := &[]models.Otp{}

	if (page <= 0) && (pageSize <= 0) {
		err := r.db.Find(Otps).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Scopes(r.paginate(page, pageSize)).Find(Otps).Error
		if err != nil {
			return nil, err
		}
	}

	return Otps, nil
}

func (r *OtpRepository) UpdateRecord(Otp *models.Otp) (*models.Otp, error) {

	fmt.Println(Otp)

	//grom 에서는 구조체로 저장시 0 값인 필드는 업데이트 하지 않는다., 맵으로 변환하여 저장해야함.
	//Put 메소드로 Update 에서는 전체 필드를 넘겨 받아서 처리한다. 단, 값이 숫자 필드가 0인 값이 있다면 맵으로 변환해서 처리 해야함.
	//err := r.db.Model(&Otp).Where("m_id = ?", Otp.UserID).Updates(Otp).Error

	//맵으로 변경해야할 값을 넘기면 0 값도 저장 됨.
	//err := r.db.Model(&Otp).Where("m_id = ?", Otp.UserID).Updates(jsonData).Error

	//select 해서 변경할 필드를 지정 가능. 이렇게하면 0 도 저장 됨.
	err := r.db.Model(&Otp).Select("*").Where("otp_id = ?", Otp.OtpID).Updates(Otp).Error
	if err != nil {
		return nil, err
	}

	return Otp, nil
}

func (r *OtpRepository) PatchRecord(otp_id string, jsonData map[string]interface{}) (*models.Otp, error) {

	fmt.Println(jsonData)

	Otp := &models.Otp{}

	//patch update 일때는 맵으로 전달된 필드만 업데이트 한다.
	err := r.db.Model(&Otp).Where("otp_id = ?", otp_id).Updates(jsonData).Error
	if err != nil {
		return nil, err
	}

	err2 := r.db.Where("otp_id = ?", otp_id).First(Otp).Error
	if err2 != nil {
		log.Printf("%v", err2)
		return nil, err2
	}

	return Otp, nil
}

func (r *OtpRepository) DeleteByID(otp_id string) error {

	Otp := &models.Otp{}

	// result := r.db.Where("otp_id = ?", otp_id).First(Otp)
	// if result.RowsAffected == 0 { // returns count of records found
	// 	fmt.Printf("count : 0")
	// 	return fiber.NewError(fiber.StatusNotFound, "No Record Found")
	// }

	err := r.db.Where("otp_id = ?", otp_id).Delete(Otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OtpRepository) GetByID(otp_id string) (*models.Otp, error) {

	Otp := &models.Otp{}

	err := r.db.Where("otp_id = ?", otp_id).First(Otp).Error
	if err != nil {
		//log.Printf("%v", err)
		return nil, err
	}

	return Otp, nil
}

func (r *OtpRepository) Create(Otp *models.Otp) (*models.Otp, error) {

	err := r.db.Create(&Otp).Error

	if err != nil {
		return nil, err
	}

	return Otp, nil
}
