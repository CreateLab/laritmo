package handlers

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/CreateLab/laritmo/internal/models"
)

const embedCheckTTL = 5 * time.Minute

var sheetIDRe = regexp.MustCompile(`^/spreadsheets/d/([^/]+)`)

type embedCacheEntry struct {
	ok      bool
	expires time.Time
}

// gradeSheetEmbedder решает, можно ли показать таблицу Google Sheets в iframe.
// Закрытая таблица (401/403/редирект на вход) не встраивается.
type gradeSheetEmbedder struct {
	client *http.Client
	mu     sync.Mutex
	cache  map[string]embedCacheEntry
}

func newGradeSheetEmbedder() *gradeSheetEmbedder {
	return &gradeSheetEmbedder{
		client: &http.Client{
			Timeout: 4 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		cache: make(map[string]embedCacheEntry),
	}
}

// embedURL превращает ссылку на таблицу в адрес для iframe; "" если ссылка не подходит.
func embedURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() != "docs.google.com" {
		return ""
	}
	q := u.Query()
	if strings.Contains(u.Path, "/pubhtml") {
		q.Set("widget", "true")
		q.Set("headers", "false")
		u.RawQuery = q.Encode()
		u.Fragment = ""
		return u.String()
	}
	m := sheetIDRe.FindStringSubmatch(u.Path)
	if m == nil || m[1] == "e" {
		return ""
	}
	gid := q.Get("gid")
	if gid == "" && u.Fragment != "" {
		if f, err := url.ParseQuery(u.Fragment); err == nil {
			gid = f.Get("gid")
		}
	}
	res := "https://docs.google.com/spreadsheets/d/" + m[1] + "/preview"
	if gid != "" {
		res += "?gid=" + url.QueryEscape(gid)
	}
	return res
}

func (e *gradeSheetEmbedder) isOpen(ctx context.Context, embed string) bool {
	e.mu.Lock()
	if c, ok := e.cache[embed]; ok && time.Now().Before(c.expires) {
		e.mu.Unlock()
		return c.ok
	}
	e.mu.Unlock()

	ok := false
	if req, err := http.NewRequestWithContext(ctx, http.MethodGet, embed, nil); err == nil {
		if resp, err := e.client.Do(req); err == nil {
			resp.Body.Close()
			ok = resp.StatusCode == http.StatusOK
		}
	}

	e.mu.Lock()
	e.cache[embed] = embedCacheEntry{ok: ok, expires: time.Now().Add(embedCheckTTL)}
	e.mu.Unlock()
	return ok
}

type gradeSheetResponse struct {
	models.GradeSheet
	EmbedURL string `json:"embed_url,omitempty"`
}

func (e *gradeSheetEmbedder) decorate(ctx context.Context, sheets []models.GradeSheet) []gradeSheetResponse {
	out := make([]gradeSheetResponse, len(sheets))
	var wg sync.WaitGroup
	for i, s := range sheets {
		out[i].GradeSheet = s
		embed := embedURL(s.SheetURL)
		if embed == "" {
			continue
		}
		wg.Add(1)
		go func(i int, embed string) {
			defer wg.Done()
			if e.isOpen(ctx, embed) {
				out[i].EmbedURL = embed
			}
		}(i, embed)
	}
	wg.Wait()
	return out
}
