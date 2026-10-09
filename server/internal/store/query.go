package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrNotFound 表示请求的条目不存在。
var ErrNotFound = errors.New("未找到")

// PersonBrief 列表用的人物摘要。
type PersonBrief struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
	Birth       string `json:"birth"`
	Death       string `json:"death"`
	BirthPlace  string `json:"birth_place"`
	Nationality string `json:"nationality"`
	Epithet     string `json:"epithet"`
	PortraitKey string `json:"portrait_key"`
	Summary     string `json:"summary"`
	Thesis      string `json:"thesis"`
}

// Section 人物的一个正文章节。
type Section struct {
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	BodyMD string `json:"body_md"`
	Ord    int    `json:"ord"`
}

// TimelineItem 生平年表条目。
type TimelineItem struct {
	Year     int    `json:"year"`
	DateText string `json:"date_text"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

// NameRef 一个引用条目（slug + 显示名）。
type NameRef struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// TermRef 指向术语的引用，显示字段为 term。
type TermRef struct {
	Slug string `json:"slug"`
	Term string `json:"term"`
}

// TitleRef 指向著作的引用，显示字段为 title。
type TitleRef struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// PersonDetail 人物详情。
type PersonDetail struct {
	PersonBrief
	Timeline []TimelineItem `json:"timeline"`
	Sections []Section      `json:"sections"`
	Works    []WorkRef      `json:"works"`
	Terms    []TermRef      `json:"terms"`
}

// WorkRef 著作引用。
type WorkRef struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

// Persons 返回全部人物，按 ord 排序。
func (s *Store) Persons(ctx context.Context) ([]PersonBrief, error) {
	return queryRows(ctx, s.db,
		`SELECT slug, name, name_en, birth, death, birth_place,
		 nationality, epithet, portrait_key, summary, thesis FROM person ORDER BY ord`,
		func(rows *sql.Rows) (PersonBrief, error) {
			var p PersonBrief
			err := rows.Scan(&p.Slug, &p.Name, &p.NameEn, &p.Birth, &p.Death, &p.BirthPlace,
				&p.Nationality, &p.Epithet, &p.PortraitKey, &p.Summary, &p.Thesis)
			if err != nil {
				return p, fmt.Errorf("扫描人物: %w", err)
			}
			return p, nil
		})
}

// Person 按 slug 返回人物详情。
func (s *Store) Person(ctx context.Context, slug string) (*PersonDetail, error) {
	var p PersonDetail
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT id, slug, name, name_en, birth, death, birth_place,
		nationality, epithet, portrait_key, summary, thesis FROM person WHERE slug = ?`, slug).
		Scan(&id, &p.Slug, &p.Name, &p.NameEn, &p.Birth, &p.Death, &p.BirthPlace,
			&p.Nationality, &p.Epithet, &p.PortraitKey, &p.Summary, &p.Thesis)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询人物 %s: %w", slug, err)
	}

	if p.Timeline, err = s.personTimeline(ctx, id); err != nil {
		return nil, err
	}
	if p.Sections, err = s.personSections(ctx, id); err != nil {
		return nil, err
	}
	if p.Works, err = s.personWorks(ctx, id); err != nil {
		return nil, err
	}
	if p.Terms, err = s.personTerms(ctx, id); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) personTimeline(ctx context.Context, id int64) ([]TimelineItem, error) {
	return queryRows(ctx, s.db, `SELECT year, date_text, title, body FROM person_timeline
		WHERE person_id = ? ORDER BY ord`,
		func(rows *sql.Rows) (TimelineItem, error) {
			var t TimelineItem
			err := rows.Scan(&t.Year, &t.DateText, &t.Title, &t.Body)
			if err != nil {
				return t, fmt.Errorf("扫描年表: %w", err)
			}
			return t, nil
		}, id)
}

func (s *Store) personSections(ctx context.Context, id int64) ([]Section, error) {
	return queryRows(ctx, s.db, `SELECT kind, title, body_md, ord FROM person_section
		WHERE person_id = ? ORDER BY ord`,
		func(rows *sql.Rows) (Section, error) {
			var sec Section
			err := rows.Scan(&sec.Kind, &sec.Title, &sec.BodyMD, &sec.Ord)
			if err != nil {
				return sec, fmt.Errorf("扫描章节: %w", err)
			}
			return sec, nil
		}, id)
}

func (s *Store) personWorks(ctx context.Context, id int64) ([]WorkRef, error) {
	return queryRows(ctx, s.db, `SELECT slug, title, year FROM work
		WHERE person_id = ? ORDER BY year`,
		func(rows *sql.Rows) (WorkRef, error) {
			var w WorkRef
			err := rows.Scan(&w.Slug, &w.Title, &w.Year)
			if err != nil {
				return w, fmt.Errorf("扫描相关著作: %w", err)
			}
			return w, nil
		}, id)
}

func (s *Store) personTerms(ctx context.Context, id int64) ([]TermRef, error) {
	return queryRows(ctx, s.db, `SELECT t.slug, t.term FROM term t
		JOIN term_person tp ON tp.term_id = t.id WHERE tp.person_id = ? ORDER BY t.id`,
		func(rows *sql.Rows) (TermRef, error) {
			var r TermRef
			err := rows.Scan(&r.Slug, &r.Term)
			if err != nil {
				return r, fmt.Errorf("扫描相关术语: %w", err)
			}
			return r, nil
		}, id)
}

// WorkBrief 列表用的著作摘要。
type WorkBrief struct {
	Slug       string `json:"slug"`
	Person     string `json:"person"`
	PersonName string `json:"person_name"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	Category   string `json:"category"`
	Summary    string `json:"summary"`
}

// Excerpt 著作摘录。
type Excerpt struct {
	Text string `json:"text"`
	Note string `json:"note"`
}

// WorkDetail 著作详情。
type WorkDetail struct {
	WorkBrief
	Source    string    `json:"source"`
	SourceURL string    `json:"source_url"`
	Excerpts  []Excerpt `json:"excerpts"`
}

// Works 返回著作列表，可按人物与标题筛选。
func (s *Store) Works(ctx context.Context, person, q string) ([]WorkBrief, error) {
	sqlText := `SELECT w.slug, COALESCE(p.slug,''), COALESCE(p.name,''), w.title, w.year,
		w.category, w.summary FROM work w LEFT JOIN person p ON p.id = w.person_id WHERE 1=1`
	var args []any
	if person != "" {
		sqlText += " AND p.slug = ?"
		args = append(args, person)
	}
	if q != "" {
		sqlText += " AND w.title LIKE ? ESCAPE '\\'"
		args = append(args, "%"+escapeLike(q)+"%")
	}
	sqlText += " ORDER BY w.year, w.id"

	return queryRows(ctx, s.db, sqlText,
		func(rows *sql.Rows) (WorkBrief, error) {
			var w WorkBrief
			err := rows.Scan(&w.Slug, &w.Person, &w.PersonName, &w.Title, &w.Year,
				&w.Category, &w.Summary)
			if err != nil {
				return w, fmt.Errorf("扫描著作: %w", err)
			}
			return w, nil
		}, args...)
}

// Work 按 slug 返回著作详情。
func (s *Store) Work(ctx context.Context, slug string) (*WorkDetail, error) {
	var w WorkDetail
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT w.id, w.slug, COALESCE(p.slug,''), COALESCE(p.name,''),
		w.title, w.year, w.category, w.summary, w.source, w.source_url
		FROM work w LEFT JOIN person p ON p.id = w.person_id WHERE w.slug = ?`, slug).
		Scan(&id, &w.Slug, &w.Person, &w.PersonName, &w.Title, &w.Year, &w.Category,
			&w.Summary, &w.Source, &w.SourceURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询著作 %s: %w", slug, err)
	}

	excerpts, err := queryRows(ctx, s.db,
		`SELECT text, note FROM work_excerpt WHERE work_id = ? ORDER BY ord`,
		func(rows *sql.Rows) (Excerpt, error) {
			var e Excerpt
			if err := rows.Scan(&e.Text, &e.Note); err != nil {
				return e, fmt.Errorf("扫描摘录: %w", err)
			}
			return e, nil
		}, id)
	if err != nil {
		return nil, err
	}
	w.Excerpts = excerpts
	return &w, nil
}

// Event 一条大事年表事件。
type Event struct {
	ID       int64     `json:"id"`
	Year     int       `json:"year"`
	DateText string    `json:"date_text"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Location string    `json:"location"`
	Category string    `json:"category"`
	People   []NameRef `json:"people"`
}

// Events 返回事件列表，可按人物与范畴筛选。
func (s *Store) Events(ctx context.Context, person, category string) ([]Event, error) {
	sqlText := `SELECT DISTINCT e.id, e.year, e.date_text, e.title, e.body, e.location, e.category
		FROM event e`
	var args []any
	if person != "" {
		sqlText += " JOIN event_person ep ON ep.event_id = e.id JOIN person p ON p.id = ep.person_id"
	}
	sqlText += " WHERE 1=1"
	if person != "" {
		sqlText += " AND p.slug = ?"
		args = append(args, person)
	}
	if category != "" {
		sqlText += " AND e.category = ?"
		args = append(args, category)
	}
	sqlText += " ORDER BY e.year, e.id"

	out, err := queryRows(ctx, s.db, sqlText,
		func(rows *sql.Rows) (Event, error) {
			var e Event
			if err := rows.Scan(&e.ID, &e.Year, &e.DateText, &e.Title, &e.Body,
				&e.Location, &e.Category); err != nil {
				return e, fmt.Errorf("扫描事件: %w", err)
			}
			e.People = []NameRef{}
			return e, nil
		}, args...)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return out, nil
	}

	ids := make([]int64, len(out))
	for i := range out {
		ids[i] = out[i].ID
	}
	// 一次性取回全部关联人物，避免 N+1 查询
	peopleByEvent, err := s.eventPeople(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range out {
		if ps, ok := peopleByEvent[out[i].ID]; ok {
			out[i].People = ps
		}
	}
	return out, nil
}

func (s *Store) eventPeople(ctx context.Context, ids []int64) (map[int64][]NameRef, error) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `SELECT ep.event_id, p.slug, p.name
		FROM event_person ep JOIN person p ON p.id = ep.person_id
		WHERE ep.event_id IN (`+placeholders+`) ORDER BY p.ord`, args...)
	if err != nil {
		return nil, fmt.Errorf("查询事件关联人物: %w", err)
	}
	defer rows.Close()

	m := map[int64][]NameRef{}
	for rows.Next() {
		var eid int64
		var r NameRef
		if err := rows.Scan(&eid, &r.Slug, &r.Name); err != nil {
			return nil, fmt.Errorf("扫描事件人物: %w", err)
		}
		m[eid] = append(m[eid], r)
	}
	return m, rows.Err()
}

// TermBrief 列表用的术语摘要。
type TermBrief struct {
	Slug       string   `json:"slug"`
	Term       string   `json:"term"`
	Aliases    []string `json:"aliases"`
	Definition string   `json:"definition"`
}

// TermDetail 术语详情。
type TermDetail struct {
	TermBrief
	Body    string     `json:"body"`
	People  []NameRef  `json:"people"`
	Related []TermRef  `json:"related"`
	Sources []TitleRef `json:"sources"`
}

// splitAliases 把入库时用顿号连接的字串还原为切片。
func splitAliases(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}
	return strings.Split(s, "、")
}

// Terms 返回术语列表，可按名称筛选。
func (s *Store) Terms(ctx context.Context, q string) ([]TermBrief, error) {
	sqlText := "SELECT slug, term, aliases, definition FROM term WHERE 1=1"
	var args []any
	if q != "" {
		sqlText += " AND (term LIKE ? ESCAPE '\\' OR aliases LIKE ? ESCAPE '\\')"
		pat := "%" + escapeLike(q) + "%"
		args = append(args, pat, pat)
	}
	sqlText += " ORDER BY id"

	return queryRows(ctx, s.db, sqlText,
		func(rows *sql.Rows) (TermBrief, error) {
			var t TermBrief
			var aliases string
			err := rows.Scan(&t.Slug, &t.Term, &aliases, &t.Definition)
			if err != nil {
				return t, fmt.Errorf("扫描术语: %w", err)
			}
			t.Aliases = splitAliases(aliases)
			return t, nil
		}, args...)
}

// Term 按 slug 返回术语详情。
func (s *Store) Term(ctx context.Context, slug string) (*TermDetail, error) {
	var t TermDetail
	var id int64
	var aliases string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, slug, term, aliases, definition, body FROM term WHERE slug = ?", slug).
		Scan(&id, &t.Slug, &t.Term, &aliases, &t.Definition, &t.Body)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询术语 %s: %w", slug, err)
	}
	t.Aliases = splitAliases(aliases)

	var err2 error
	if t.People, err2 = s.nameRefs(ctx, `SELECT p.slug, p.name FROM term_person tp
		JOIN person p ON p.id = tp.person_id WHERE tp.term_id = ? ORDER BY p.ord`, id); err2 != nil {
		return nil, err2
	}
	if t.Related, err2 = s.termRefs(ctx, `SELECT t2.slug, t2.term FROM term_related tr
		JOIN term t2 ON t2.id = tr.related_id WHERE tr.term_id = ? ORDER BY t2.id`, id); err2 != nil {
		return nil, err2
	}
	if t.Sources, err2 = s.titleRefs(ctx, `SELECT w.slug, w.title FROM term_work tw
		JOIN work w ON w.id = tw.work_id WHERE tw.term_id = ? ORDER BY w.year`, id); err2 != nil {
		return nil, err2
	}
	return &t, nil
}

func (s *Store) nameRefs(ctx context.Context, query string, id int64) ([]NameRef, error) {
	return queryRows(ctx, s.db, query,
		func(rows *sql.Rows) (NameRef, error) {
			var r NameRef
			err := rows.Scan(&r.Slug, &r.Name)
			if err != nil {
				return r, fmt.Errorf("扫描术语关联人物: %w", err)
			}
			return r, nil
		}, id)
}

func (s *Store) termRefs(ctx context.Context, query string, id int64) ([]TermRef, error) {
	return queryRows(ctx, s.db, query,
		func(rows *sql.Rows) (TermRef, error) {
			var r TermRef
			err := rows.Scan(&r.Slug, &r.Term)
			if err != nil {
				return r, fmt.Errorf("扫描术语关联: %w", err)
			}
			return r, nil
		}, id)
}

func (s *Store) titleRefs(ctx context.Context, query string, id int64) ([]TitleRef, error) {
	return queryRows(ctx, s.db, query,
		func(rows *sql.Rows) (TitleRef, error) {
			var r TitleRef
			err := rows.Scan(&r.Slug, &r.Title)
			if err != nil {
				return r, fmt.Errorf("扫描术语关联著作: %w", err)
			}
			return r, nil
		}, id)
}

// GraphNode 关系图节点。
type GraphNode struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Epithet string `json:"epithet"`
}

// GraphEdge 关系图边。
type GraphEdge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Kind  string `json:"kind"`
	Label string `json:"label"`
	Note  string `json:"note"`
}

// Graph 返回关系图的全部节点与边。
func (s *Store) Graph(ctx context.Context) ([]GraphNode, []GraphEdge, error) {
	nodes, err := queryRows(ctx, s.db, "SELECT slug, name, epithet FROM person ORDER BY ord",
		func(rows *sql.Rows) (GraphNode, error) {
			var n GraphNode
			if err := rows.Scan(&n.Slug, &n.Name, &n.Epithet); err != nil {
				return n, fmt.Errorf("扫描节点: %w", err)
			}
			return n, nil
		})
	if err != nil {
		return nil, nil, err
	}

	edges, err := queryRows(ctx, s.db, `SELECT pf.slug, pt.slug, r.kind, r.label, r.note
		FROM relation r JOIN person pf ON pf.id = r.from_id JOIN person pt ON pt.id = r.to_id ORDER BY r.id`,
		func(rows *sql.Rows) (GraphEdge, error) {
			var e GraphEdge
			if err := rows.Scan(&e.From, &e.To, &e.Kind, &e.Label, &e.Note); err != nil {
				return e, fmt.Errorf("扫描关系边: %w", err)
			}
			return e, nil
		})
	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}

// Counts 各类条目数量，用于站点元信息。
type Counts struct {
	People int `json:"people"`
	Works  int `json:"works"`
	Events int `json:"events"`
	Terms  int `json:"terms"`
}

// Counts 统计各类条目数量。
//
// 运行时数据库以 immutable 只读打开，条目数在整个进程生命周期内不变，
// 故用 sync.Once 缓存首次统计结果，避免 /api/meta 每次请求都跑 4 条 COUNT(*)。
func (s *Store) Counts(ctx context.Context) (Counts, error) {
	s.countsOnce.Do(func() {
		var c Counts
		queries := []struct {
			sql string
			dst *int
		}{
			{"SELECT COUNT(*) FROM person", &c.People},
			{"SELECT COUNT(*) FROM work", &c.Works},
			{"SELECT COUNT(*) FROM event", &c.Events},
			{"SELECT COUNT(*) FROM term", &c.Terms},
		}
		for _, q := range queries {
			if err := s.db.QueryRowContext(ctx, q.sql).Scan(q.dst); err != nil {
				s.countsErr = fmt.Errorf("统计条目: %w", err)
				return
			}
		}
		s.counts = c
	})
	return s.counts, s.countsErr
}
