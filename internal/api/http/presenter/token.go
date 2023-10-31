package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type tokenPresenter struct {
}

func NewTokenPresenter() *tokenPresenter {
	return &tokenPresenter{}
}

func (p *tokenPresenter) ToTokenView(token *entity.Token) (*view.TokenView, error) {
	token_string, err := token.String()
	if err != nil {
		return nil, err
	}

	return &view.TokenView{
		Token: token_string,
	}, nil
}
