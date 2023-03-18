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

func (s *OtpServices) GetAll(page, pageSize int) (*[]models.Otp, error) {
	Otps, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}
	return Otps, nil
}

func (s *OtpServices) GetRecordByID(id string) (*models.Otp, error) {

	Otp, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return Otp, nil
}

func (s *OtpServices) OtpVerify(otp_id string, otp_num string) (bool, string) {

	Otp, err := s.repo.GetByID(otp_id)
	if err != nil {
		return false, ""
	}

	fmt.Printf("%v\n", Otp)

	hasher := &gotp.Hasher{
		HashName: "SHA1",
		Digest:   sha1.New,
	}

	if Otp.Algorithms == "SHA256" {
		hasher = &gotp.Hasher{
			HashName: Otp.Algorithms,
			Digest:   sha256.New,
		}
	} else if Otp.Algorithms == "SHA512" {
		hasher = &gotp.Hasher{
			HashName: Otp.Algorithms,
			Digest:   sha512.New,
		}
	}

	otp := gotp.NewTOTP(Otp.SecretKey, Otp.Digit, Otp.Cycle, hasher)
	otpValue := otp.Now()
	fmt.Printf("current otp : %s\n", otpValue)

	return otp.Verify(otp_num, time.Now().Unix()), otpValue
}

func (s *OtpServices) UpdateRecord(Otp *models.Otp) (*models.Otp, error) {

	OtpNew, err := s.repo.UpdateRecord(Otp)

	if err != nil {
		return nil, err
	}

	return OtpNew, nil
}

func (s *OtpServices) PatchRecord(userId string, jsonData map[string]interface{}) (*models.Otp, error) {

	OtpNew, err := s.repo.PatchRecord(userId, jsonData)

	if err != nil {
		return nil, err
	}

	return OtpNew, nil
}

func (s *OtpServices) DeleteRecord(otp_id string) error {

	err := s.repo.DeleteByID(otp_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *OtpServices) CreateRecord(Otp *models.Otp) (*models.Otp, error) {

	//secret key random
	//10byte 사용하면 base32 인코딩시 패딩없이 16자리 나옴
	//20byte 사용하면 base32 인코딩시 패딩없이 32자리 나옴
	keylength := 20
	Otp.SecretKey = gotp.RandomSecret(keylength)
	OtpNew, err := s.repo.Create(Otp)

	if err != nil {
		return nil, err
	}

	return OtpNew, nil
}
