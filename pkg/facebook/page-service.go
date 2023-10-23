package facebook

type PageService struct {
	accessToken string
}

func NewPageService(accessToken string) *PageService {
	return &PageService{accessToken}
}
