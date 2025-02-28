package business

import (
	craveModel "crave/shared/model"
	"regexp"
)

type PersonBusiness struct {
}

func NewPersonBusiness() *PersonBusiness {
	return &PersonBusiness{}
}

func (biz *PersonBusiness) GetFilterFlag() craveModel.Filter {
	return craveModel.PERSON
}

func (biz *PersonBusiness) Apply(html *string) craveModel.Filter {
	if biz.containsProfile(html) {
		return 0
	}
	return biz.GetFilterFlag()
}

func (biz *PersonBusiness) containsProfile(html *string) bool {
	//pattern := `<strong[^>]*\\s*>\\s*(본명|출생|거주지|본관|학력|부모|친인척|형제자매|자녀|병역|별명|종교|서명|플랫폼)\\s*</strong>`
	pattern := `<strong[^>]*>\s*(본명|출생|거주지|본관|학력|부모|친인척|형제자매|자녀|병역|별명|종교|서명|플랫폼)\s*</strong>`
	re := regexp.MustCompile(pattern)
	return re.MatchString(*html)
}

//룩삼 <-> 김진효 필터
