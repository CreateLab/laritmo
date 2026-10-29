package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// dumpTables are the tables included in the dump. Order is preserved in the output.
var dumpTables = []string{"courses", "lectures", "labs", "grade_sheets", "exam_questions", "users"}

// dumpExcludedColumns are never written to the dump (secrets).
var dumpExcludedColumns = map[string]map[string]bool{
	"users": {"password_hash": true},
}

type DumpHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewDumpHandler(db *sql.DB, logger *slog.Logger) *DumpHandler {
	return &DumpHandler{db: db, logger: logger}
}

// Dump godoc
// @Summary Full database dump (admin)
// @Description Returns all tables as a JSON file. Password hashes are excluded.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {file} file
// @Router /admin/dump [get]
func (h *DumpHandler) Dump(c *gin.Context) {
	ctx := c.Request.Context()
	tables := make(map[string][]map[string]any, len(dumpTables))

	for _, table := range dumpTables {
		rows, err := h.dumpTable(c, table)
		if err != nil {
			h.logger.ErrorContext(ctx, "Failed to dump table", "table", table, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dump database"})
			return
		}
		tables[table] = rows
	}

	now := time.Now()
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="laritmo-dump-%s.json"`, now.Format("2006-01-02-150405")))
	c.JSON(http.StatusOK, gin.H{
		"exported_at": now.UTC().Format(time.RFC3339),
		"tables":      tables,
	})
	h.logger.InfoContext(ctx, "Database dump exported", "user", c.GetString("username"))
}

func (h *DumpHandler) dumpTable(c *gin.Context, table string) ([]map[string]any, error) {
	// table comes from the fixed dumpTables list, never from user input
	rows, err := h.db.QueryContext(c.Request.Context(), "SELECT * FROM `"+table+"`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, 0)
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any, len(cols))
		for i, col := range cols {
			if dumpExcludedColumns[table][col] {
				continue
			}
			if b, ok := vals[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = vals[i]
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
