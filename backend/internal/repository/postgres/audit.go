package postgres

import (
	"context"

	"pos-system/backend/internal/app/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct{ db *pgxpool.Pool }

func (r *AuditRepo) List(ctx context.Context, action, entityType, query string, limit int) ([]domain.AuditLog, error) {
	if limit <= 0 || limit > 300 {
		limit = 150
	}
	rows, err := r.db.Query(ctx, `
		SELECT al.id, al.user_id, COALESCE(u.name, ''), COALESCE(u.email, ''), al.action, al.entity_type,
			al.entity_id, COALESCE(al.old_data::text, ''), COALESCE(al.new_data::text, ''),
			COALESCE(al.ip_address::text, ''), al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id=al.user_id
		WHERE al.deleted_at IS NULL
			AND ($1='' OR al.action=$1)
			AND ($2='' OR al.entity_type=$2)
			AND (
				$3=''
				OR al.action ILIKE '%'||$3||'%'
				OR al.entity_type ILIKE '%'||$3||'%'
				OR al.entity_id::text ILIKE '%'||$3||'%'
				OR u.name ILIKE '%'||$3||'%'
				OR u.email ILIKE '%'||$3||'%'
			)
		ORDER BY al.created_at DESC
		LIMIT $4`, action, entityType, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]domain.AuditLog, 0)
	for rows.Next() {
		var log domain.AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.UserName, &log.UserEmail, &log.Action, &log.EntityType,
			&log.EntityID, &log.OldData, &log.NewData, &log.IPAddress, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}
