package shipping

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotFound       = errors.New("not found")
	InternalServerErr = errors.New("internal server error")

	Validate = validator.New()
)

type PackRepository interface {
	GetByProductID(ctx context.Context, productID uint64) ([]uint64, error)
}

type PackConfig struct {
	Count int64  `json:"number_of_packs"`
	Size  uint64 `json:"pack_size"`
}

func (pc PackConfig) String() string {
	return fmt.Sprintf("%d x %d", pc.Count, pc.Size)
}
