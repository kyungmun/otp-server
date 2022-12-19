package models

import "gorm.io/gorm"

type OtpRegistry struct {
	gorm.Model
	OtpID      string `gorm:"column:otp_id; primaryKey; type:varchar(32); not null" json:"otp_id" validate:"required"`
	SecretKey  string `gorm:"column:secret_key; type:varchar(32); not null" json:"secret_key"`
	Algorithem string `gorm:"column:algorithem; type:varchar(8); not null" json:"algorithem" validate:"required"`
	Digit      int    `gorm:"column:digit; type:int; not null" json:"digit" validate:"required"`
	Cycle      int    `gorm:"column:cycle; type:int; not null" json:"cycle" validate:"required"`
}

// TableName overrides the table name used by User to `tablename`
// 자동으로 생성되는 테이블명은 구조체명을 기준으로 스네이크 스타일 및 끝에 s 문자가 자동으로 붙기에 지정한 테이블명을 사용하려면 아래 메소드 구현 필요.
//func (OtpRegistrys) TableName() string {
//	return "otp_registries"
//}

func MigrateOtpRegistrys(db *gorm.DB) error {
	err := db.AutoMigrate(&OtpRegistry{})
	return err
}

// 데이터 저장용 인터페이스
type OtpRepository interface {
	GetIndex() ([]*OtpRegistry, error)
	GetByID(otp_id string) (*OtpRegistry, error)
	Fetch(offset, limit int) ([]*OtpRegistry, error)
	Create(otpRegistry *OtpRegistry) (*OtpRegistry, error)
	Update(otp_id string, otpRegistry *OtpRegistry) (*OtpRegistry, error)
	Delete(otp_id string) error
}

// 데이터 비즈니스로직 처리용 인터페이스
type OtpRegistryUseCase interface {
	GetByID(otp_id string) (*OtpRegistry, error)
	Fetch(offset, limit int) ([]*OtpRegistry, error)
	Create(otpRegistry *OtpRegistry) (*OtpRegistry, error)
	Update(otp_id string, otpRegistry *OtpRegistry) (*OtpRegistry, error)
	Delete(otp_id string) error
}
