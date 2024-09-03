package stores

const (
	CREATEQUERY  = "INSERT INTO templates (tenant_id, id, name, description, content, service, universal, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)"
	GETBYIDQUERY = "SELECT tenant_id, id, name, description, content, service, universal, created_at, updated_at FROM templates WHERE id = $1 AND tenant_id = $2"
	DELETEQUERY  = "DELETE FROM templates where id=$1 AND tenant_id = $2"
	GETQUERY     = "SELECT tenant_id, id, name, description, content, service, universal, created_at, updated_at FROM templates WHERE tenant_id = $1 AND "
	COUNTQUERY   = "SELECT count(*) FROM templates WHERE tenant_id = $1"
)
