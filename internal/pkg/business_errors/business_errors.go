package business_errors

import "errors"

var (
	// ErrMailingNotFound .
	ErrMailingNotFound = errors.New("рассылка не найдена")
	// ErrClientNotFound .
	ErrClientNotFound = errors.New("клиент не найден")
	// ErrInvalidID .
	ErrInvalidID = errors.New("ID не соответствует правилам UUID")
	// ErrNumberFirstDigitDontSeven .
	ErrNumberFirstDigitDontSeven = errors.New("невалидный номер телефона(первая цифра не равна 7)")
	// ErrNumberLenMoreEleven .
	ErrNumberLenMoreEleven = errors.New("невалидный номер телефона(количество цифр не равно 11)")
	// ErrNumberMustConsistOfDigits .
	ErrNumberMustConsistOfDigits = errors.New("номер должен состоять из цифр")
	// ErrOperatorCodeMustConsistOfDigits .
	ErrOperatorCodeMustConsistOfDigits = errors.New("мобильный код оператора должен состоять из цифр")
	// ErrOperatorCodeNotInInterval .
	ErrOperatorCodeNotInInterval = errors.New("мобильный код оператора должен быть от 900 до 1000")
	// ErrTagLenMoreHundred .
	ErrTagLenMoreHundred = errors.New("тег не может быть больше 100 символов")
	// ErrTextLenMoreThousand .
	ErrTextLenMoreThousand = errors.New("текст рассылки не может быть больше 1000 символов")
	// ErrFilterLenMoreThousand .
	ErrFilterLenMoreThousand = errors.New("фильтр не может быть больше 100 символов")
	// ErrMailingTimeIsUp .
	ErrMailingTimeIsUp = errors.New("время рассылки вышло")
)
