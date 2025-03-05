package business

import (
	craveModel "crave/shared/model"
	"regexp"
)

type NameBusiness struct {
}

func NewNameBusiness() *NameBusiness {
	return &NameBusiness{}
}

func (biz *NameBusiness) GetFilterFlag() craveModel.Filter {
	return craveModel.NAME
}

func (biz *NameBusiness) Apply(name *string) craveModel.Filter {
	if biz.containsDateFormat(name) || biz.containsSpecialChars(name) || biz.containsChineseCharacter(name) || biz.containsWord(name) {
		return biz.GetFilterFlag()
	}
	return 0
}

func (biz *NameBusiness) containsWord(name *string) bool {
	pattern := `\b(분류|출생|대한민국|미국|중국)\b`
	re := regexp.MustCompile(pattern)
	return re.MatchString(*name)
}

func (biz *NameBusiness) containsChineseCharacter(name *string) bool {
	pattern := `[\p{Han}]`
	re := regexp.MustCompile(pattern)
	return re.MatchString(*name)
}

func (biz *NameBusiness) containsSpecialChars(name *string) bool {
	pattern := `[/!?@.,:&[]`
	re := regexp.MustCompile(pattern)
	return re.MatchString(*name)
}

func (biz *NameBusiness) containsDateFormat(name *string) bool {
	//pattern := `\d+(?:년|년대|년도|\d+월\s*\d+일)`
	//pattern := `\d+(?:년|년대|년도)|\d{1,2}월\s*\d{1,2}일|\d{4}년\s*\d{1,2}월\s*\d{1,2}일`
	pattern := `\d+(?:년|년대|년도)|\d{1,2}월(?:\s*\d{1,2}일)?|\d{4}년\s*\d{1,2}월(?:\s*\d{1,2}일)?`
	re := regexp.MustCompile(pattern)

	return re.MatchString(*name)
}

//룩삼 <-> 김진효 필터
