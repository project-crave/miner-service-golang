package page

type IClient interface {
	GetHtml(url string) (*string, error)
}
