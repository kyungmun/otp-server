package repository

import (
	"fmt"
	"log"

	"github.com/kyungmun/otp-server/models"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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

func (r *UserRepository) GetAll(page, pageSize int) (*[]models.User, error) {
	Users := &[]models.User{}

	if (page <= 0) && (pageSize <= 0) {
		err := r.db.Find(Users).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Scopes(r.paginate(page, pageSize)).Find(Users).Error
		if err != nil {
			return nil, err
		}
	}

	return Users, nil
}

func (r *UserRepository) UpdateRecord(User *models.User) (*models.User, error) {

	fmt.Println(User)

	//grom 에서는 구조체로 저장시 0 값인 필드는 업데이트 하지 않는다., 맵으로 변환하여 저장해야함.
	//Put 메소드로 Update 에서는 전체 필드를 넘겨 받아서 처리한다. 단, 값이 숫자 필드가 0인 값이 있다면 맵으로 변환해서 처리 해야함.
	//err := r.db.Model(&User).Where("m_id = ?", User.UserID).Updates(User).Error

	//맵으로 변경해야할 값을 넘기면 0 값도 저장 됨.
	//err := r.db.Model(&User).Where("m_id = ?", User.UserID).Updates(jsonData).Error

	//select 해서 변경할 필드를 지정 가능. 이렇게하면 0 도 저장 됨.
	err := r.db.Model(&User).Select("*").Where("username = ?", User.Username).Updates(User).Error
	if err != nil {
		return nil, err
	}

	return User, nil
}

func (r *UserRepository) PatchRecord(username string, jsonData map[string]interface{}) (*models.User, error) {

	fmt.Println(jsonData)

	User := &models.User{}

	//patch update 일때는 맵으로 전달된 필드만 업데이트 한다.
	err := r.db.Model(&User).Where("username = ?", username).Updates(jsonData).Error
	if err != nil {
		return nil, err
	}

	err2 := r.db.Where("username = ?", username).First(User).Error
	if err2 != nil {
		log.Printf("%v", err2)
		return nil, err2
	}

	return User, nil
}

func (r *UserRepository) DeleteByID(username string) error {

	User := &models.User{}

	// result := r.db.Where("otp_id = ?", otp_id).First(User)
	// if result.RowsAffected == 0 { // returns count of records found
	// 	fmt.Printf("count : 0")
	// 	return fiber.NewError(fiber.StatusNotFound, "No Record Found")
	// }

	err := r.db.Where("username = ?", username).Delete(User).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(username string) (*models.User, error) {

	User := &models.User{}

	err := r.db.Where("username = ?", username).First(User).Error
	if err != nil {
		//log.Printf("%v", err)
		return nil, err
	}

	return User, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {

	User := &models.User{}

	err := r.db.Where("email = ?", email).First(User).Error
	if err != nil {
		//log.Printf("%v", err)
		return nil, err
	}

	return User, nil
}

func (r *UserRepository) Create(User *models.User) (*models.User, error) {

	err := r.db.Create(&User).Error

	if err != nil {
		return nil, err
	}

	return User, nil
}
