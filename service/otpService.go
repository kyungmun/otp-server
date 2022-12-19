package service

// OTP 비즈니스 로직만 처리하는 서비스

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/kyungmun/otp-server/models"
	"github.com/kyungmun/otp-server/repository"
	"github.com/xlzd/gotp"
)

func NewOtpServices(r *repository.OtpRepository) (*OtpServices, error) {
	return &OtpServices{repo: r}, nil
}

type OtpServices struct {
	repo *repository.OtpRepository
}

func (s *OtpServices) GetAll(page, pageSize int) (*[]models.OtpRegistry, error) {
	otpRegistrys, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}
	return otpRegistrys, nil
}

func (s *OtpServices) GetRecordByID(id string) (*models.OtpRegistry, error) {

	otpRegistry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return otpRegistry, nil
}

func (s *OtpServices) OtpVerify(otp_id string, otp_num string) (bool, string) {

	otpRegistry, err := s.repo.GetByID(otp_id)
	if err != nil {
		return false, ""
	}

	fmt.Printf("%v\n", otpRegistry)

	hasher := &gotp.Hasher{
		HashName: "SHA1",
		Digest:   sha1.New,
	}

	if otpRegistry.Algorithem == "SHA256" {
		hasher = &gotp.Hasher{
			HashName: otpRegistry.Algorithem,
			Digest:   sha256.New,
		}
	} else if otpRegistry.Algorithem == "SHA512" {
		hasher = &gotp.Hasher{
			HashName: otpRegistry.Algorithem,
			Digest:   sha512.New,
		}
	}

	otp := gotp.NewTOTP(otpRegistry.SecretKey, otpRegistry.Digit, otpRegistry.Cycle, hasher)
	otpValue := otp.Now()
	fmt.Printf("current otp : %s\n", otpValue)

	return otp.Verify(otp_num, time.Now().Unix()), otpValue
}

func (s *OtpServices) UpdateRecord(otpRegistry *models.OtpRegistry) (*models.OtpRegistry, error) {

	otpRegistryNew, err := s.repo.UpdateRecord(otpRegistry)

	if err != nil {
		return nil, err
	}

	return otpRegistryNew, nil
}

func (s *OtpServices) PatchRecord(userId string, jsonData map[string]interface{}) (*models.OtpRegistry, error) {

	otpRegistryNew, err := s.repo.PatchRecord(userId, jsonData)

	if err != nil {
		return nil, err
	}

	return otpRegistryNew, nil
}

func (s *OtpServices) DeleteRecord(otp_id string) error {

	err := s.repo.DeleteByID(otp_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *OtpServices) CreateRecord(otpRegistry *models.OtpRegistry) (*models.OtpRegistry, error) {

	//secret key random (10byte 사용해야 base32 인코딩하면 패딩없이 16자리 나옴)
	keylength := 10
	otpRegistry.SecretKey = gotp.RandomSecret(keylength)
	otpRegistryNew, err := s.repo.Create(otpRegistry)

	if err != nil {
		return nil, err
	}

	return otpRegistryNew, nil
}
