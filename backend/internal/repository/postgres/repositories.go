package postgres

import (
	"pos-system/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	db          *pgxpool.Pool
	auth        *AuthRepo
	products    *ProductRepo
	inventories *InventoryRepo
	sales       *SaleRepo
	users       *UserRepo
	dashboard   *DashboardRepo
	system      *SystemRepo
	audit       *AuditRepo
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	r := &Repositories{db: db}
	r.auth = &AuthRepo{db: db}
	r.products = &ProductRepo{db: db}
	r.inventories = &InventoryRepo{db: db}
	r.sales = &SaleRepo{db: db}
	r.users = &UserRepo{db: db}
	r.dashboard = &DashboardRepo{db: db}
	r.system = &SystemRepo{db: db}
	r.audit = &AuditRepo{db: db}
	return r
}

func (r *Repositories) Auth() repository.AuthRepository             { return r.auth }
func (r *Repositories) Products() repository.ProductRepository      { return r.products }
func (r *Repositories) Inventories() repository.InventoryRepository { return r.inventories }
func (r *Repositories) Sales() repository.SaleRepository            { return r.sales }
func (r *Repositories) Users() repository.UserRepository            { return r.users }
func (r *Repositories) Dashboard() repository.DashboardRepository   { return r.dashboard }
func (r *Repositories) System() repository.SystemRepository         { return r.system }
func (r *Repositories) Audit() repository.AuditRepository           { return r.audit }
