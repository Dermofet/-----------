package presenter

import (
	"fmt"
	"go-test-grpc-http/internal/api/http/view"
	"go-test-grpc-http/internal/entity"
)

type tokenPresenter struct {
}

func NewTokenPresenter() *tokenPresenter {
	return &tokenPresenter{}
}

func (t *tokenPresenter) ToTokenView(token *entity.Token) (*view.TokenView, error) {
	res, err := token.String()
	if err != nil {
		return nil, fmt.Errorf("can't make token view: %w", err)
	}
	return &view.TokenView{
		Token: res,
	}, nil
}
