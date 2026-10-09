package handlers

import (
	"context"
	"encoding/csv"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/CreateLab/laritmo/internal/models"
)

const (
	tableCacheTTL = 5 * time.Minute
	maxTableRows  = 500
	maxTableCols  = 40
)

var sheetIDRe = regexp.MustCompile(`^/spreadsheets/d/([^/]+)`)

type tableCacheEntry struct {
	table   [][]string
	expires time.Time
}

// gradeSheetTable забирает данные из открытой Google-таблицы в виде CSV,
// чтобы показать только заполненные ячейки, а не весь лист Google.
// Закрытая таблица (401/403/редирект на вход) даёт nil.
type gradeSheetTable struct {
	client *http.Client
	mu     sync.Mutex
	cache  map[string]tableCacheEntry
}

func newGradeSheetTable() *gradeSheetTable {
	return &gradeSheetTable{
		client: &http.Client{
			Timeout: 5 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		cache: make(map[string]tableCacheEntry),
	}
}

// csvURL строит адрес выгрузки данных листа; "" если ссылка не на Google Sheets.
func csvURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() != "docs.google.com" {
		return ""
	}
	m := sheetIDRe.FindStringSubmatch(u.Path)
	if m == nil || m[1] == "e" {
		return ""
	}
	gid := u.Query().Get("gid")
	if gid == "" && u.Fragment != "" {
		if f, err := url.ParseQuery(u.Fragment); err == nil {
			gid = f.Get("gid")
		}
	}
	res := "https://docs.google.com/spreadsheets/d/" + m[1] + "/gviz/tq?tqx=out:csv&headers=1"
	if gid != "" {
		res += "&gid=" + url.QueryEscape(gid)
	}
	return res
}

func (t *gradeSheetTable) fetch(ctx context.Context, src string) [][]string {
	t.mu.Lock()
	if c, ok := t.cache[src]; ok && time.Now().Before(c.expires) {
		t.mu.Unlock()
		return c.table
	}
	t.mu.Unlock()

	table := t.download(ctx, src)

	t.mu.Lock()
	t.cache[src] = tableCacheEntry{table: table, expires: time.Now().Add(tableCacheTTL)}
	t.mu.Unlock()
	return table
}

func (t *gradeSheetTable) download(ctx context.Context, src string) [][]string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return nil
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	r := csv.NewReader(resp.Body)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil || len(rows) == 0 {
		return nil
	}
	if len(rows) > maxTableRows {
		rows = rows[:maxTableRows]
	}
	for i, row := range rows {
		if len(row) > maxTableCols {
			rows[i] = row[:maxTableCols]
		}
	}
	return rows
}

type gradeSheetResponse struct {
	models.GradeSheet
	Table [][]string `json:"table,omitempty"`
}

func (t *gradeSheetTable) decorate(ctx context.Context, sheets []models.GradeSheet) []gradeSheetResponse {
	out := make([]gradeSheetResponse, len(sheets))
	var wg sync.WaitGroup
	for i, s := range sheets {
		out[i].GradeSheet = s
		src := csvURL(s.SheetURL)
		if src == "" {
			continue
		}
		wg.Add(1)
		go func(i int, src string) {
			defer wg.Done()
			out[i].Table = t.fetch(ctx, src)
		}(i, src)
	}
	wg.Wait()
	return out
}
