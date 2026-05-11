package server

import (
	"github.com/aszaychik/prima-inti-api/internal/brand"
	"github.com/aszaychik/prima-inti-api/internal/category"
	"github.com/aszaychik/prima-inti-api/internal/company"
	"github.com/aszaychik/prima-inti-api/internal/user"
)

type Handlers struct {
	User     *user.Handler
	Company  *company.Handler
	Category *category.Handler
	Brand    *brand.Handler
}
