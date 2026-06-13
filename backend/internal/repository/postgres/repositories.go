package postgres

import (
	"pos-system/backend/internal/app/port"

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

func (r *Repositories) Auth() port.AuthRepository             { return r.auth }
func (r *Repositories) Products() port.ProductRepository      { return r.products }
func (r *Repositories) Inventories() port.InventoryRepository { return r.inventories }
func (r *Repositories) Sales() port.SaleRepository            { return r.sales }
func (r *Repositories) Users() port.UserRepository            { return r.users }
func (r *Repositories) Dashboard() port.DashboardRepository   { return r.dashboard }
func (r *Repositories) System() port.SystemRepository         { return r.system }
func (r *Repositories) Audit() port.AuditRepository           { return r.audit }
