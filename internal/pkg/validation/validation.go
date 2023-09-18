package validation

import (
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/fr5270937/notifications_service/internal/model"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/business_errors"
)

func ValidateCreateClientRequest(req model.CreateClientRequest) error {
	if err := ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return err
	}
	if err := ValidateTag(req.Tag); err != nil {
		return err
	}
	if err := ValidateMobileOperatorCode(req.MobileOperatorCode); err != nil {
		return err
	}

	return nil
}

func ValidateUpdateClientRequest(id string, req model.UpdateClientRequest) error {
	if err := ValidateID(id); err != nil {
		return err
	}
	if err := ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return err
	}
	if err := ValidateTag(req.Tag); err != nil {
		return err
	}
	if err := ValidateMobileOperatorCode(req.MobileOperatorCode); err != nil {
		return err
	}

	return nil
}

func ValidateDeleteClientRequest(id string) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	return nil
}

func ValidateCreateMailingRequest(req model.CreateMailingRequest) error {
	if err := ValidateText(req.Text); err != nil {
		return err
	}
	if err := ValidateFilter(req.Filter); err != nil {
		return err
	}

	return nil
}

func ValidateUpdateMailingRequest(id string, req model.UpdateMailingRequest) error {
	if err := ValidateID(id); err != nil {
		return err
	}
	if err := ValidateText(req.Text); err != nil {
		return err
	}
	if err := ValidateFilter(req.Filter); err != nil {
		return err
	}

	return nil
}

func ValidateDeleteMailingRequest(id string) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	return nil
}

func ValidateMobileOperatorCode(mobileOperatorCode int32) error {
	if mobileOperatorCode < 900 || mobileOperatorCode > 1000 {
		return business_errors.ErrOperatorCodeNotInInterval
	}

	return nil
}

func ValidateTag(tag string) error {
	if len(tag) > 100 {
		return business_errors.ErrTagLenMoreHundred
	}

	return nil
}

func ValidateID(ID string) error {
	_, err := uuid.Parse(ID)
	if err != nil {
		return business_errors.ErrInvalidID
	}

	return nil
}

func ValidateText(str string) error {
	if len(str) > 1000 {
		return business_errors.ErrTextLenMoreThousand
	}

	return nil
}

func ValidateFilter(filter string) error {
	if len(filter) > 100 {
		return business_errors.ErrFilterLenMoreThousand
	}
	if len(filter) > 0 && isDigit(filter[0]) {
		mobileOperatorCode, err := strconv.Atoi(filter)
		if err != nil {
			return business_errors.ErrOperatorCodeMustConsistOfDigits
		}
		if err = ValidateMobileOperatorCode(int32(mobileOperatorCode)); err != nil {
			return err
		}
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 0 {
		switch {
		case phoneNumber[0] != '7':
			return business_errors.ErrNumberFirstDigitDontSeven
		case len(phoneNumber) != 11:
			return business_errors.ErrNumberLenMoreEleven
		}
	}
	for idx := range phoneNumber {
		if !isDigit(phoneNumber[idx]) {
			return business_errors.ErrNumberMustConsistOfDigits
		}
	}

	return nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
