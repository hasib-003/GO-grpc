package entity

type PaginationRequest struct {
	Limit int `query:"limit" validate:"required,gte=1,lte=100"`
	Page  int `query:"page" validate:"required,gte=0"`
}

func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}

	offset := (p.Page - 1) * p.Limit
	return offset
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}

	return p.Limit
}
