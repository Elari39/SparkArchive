package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sparkarchive/server/internal/content"
)

// Ingest 把整份 Archive 写入空数据库，并建立全文索引。
// 整个过程在单个事务中完成：要么全部成功，要么完全回滚。
func (s *Store) Ingest(ctx context.Context, a *content.Archive) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚是空操作

	personIDs := make(map[string]int64, len(a.People))
	for i, p := range a.People {
		id, err := insertPerson(ctx, tx, p, i)
		if err != nil {
			return err
		}
		personIDs[p.Slug] = id
		if err := insertPersonDetail(ctx, tx, id, p); err != nil {
			return err
		}
	}

	workIDs := make(map[string]int64, len(a.Works))
	for i, w := range a.Works {
		pid := personIDs[w.Person]
		id, err := insertWork(ctx, tx, w, pid, i)
		if err != nil {
			return err
		}
		workIDs[w.Slug] = id
		for j, e := range w.Excerpts {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO work_excerpt (work_id, text, note, ord) VALUES (?, ?, ?, ?)`,
				id, e.Text, e.Note, j); err != nil {
				return fmt.Errorf("写入摘录 %s#%d: %w", w.Slug, j, err)
			}
		}
	}

	for _, e := range a.Events {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO event (year, date_text, title, body, location, category) VALUES (?, ?, ?, ?, ?, ?)`,
			e.Year, e.DateText, e.Title, e.Body, e.Location, e.Category)
		if err != nil {
			return fmt.Errorf("写入事件 %s: %w", e.Title, err)
		}
		eid, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取事件 ID: %w", err)
		}
		for _, slug := range e.People {
			pid, ok := personIDs[slug]
			if !ok {
				return fmt.Errorf("事件 %q 引用了未知人物 %q", e.Title, slug)
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO event_person (event_id, person_id) VALUES (?, ?)`, eid, pid); err != nil {
				return fmt.Errorf("关联事件人物: %w", err)
			}
		}
	}

	termIDs := make(map[string]int64, len(a.Terms))
	for _, t := range a.Terms {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO term (slug, term, aliases, definition, body) VALUES (?, ?, ?, ?, ?)`,
			t.Slug, t.Term, strings.Join(t.Aliases, "、"), t.Definition, t.Body)
		if err != nil {
			return fmt.Errorf("写入术语 %s: %w", t.Slug, err)
		}
		tid, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取术语 ID: %w", err)
		}
		termIDs[t.Slug] = tid
		for _, slug := range t.People {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO term_person (term_id, person_id) VALUES (?, ?)`, tid, personIDs[slug]); err != nil {
				return fmt.Errorf("关联术语人物: %w", err)
			}
		}
		for _, slug := range t.Sources {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO term_work (term_id, work_id) VALUES (?, ?)`, tid, workIDs[slug]); err != nil {
				return fmt.Errorf("关联术语著作: %w", err)
			}
		}
	}
	// 关联术语需等所有术语 ID 就绪
	for _, t := range a.Terms {
		for _, slug := range t.Related {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO term_related (term_id, related_id) VALUES (?, ?)`,
				termIDs[t.Slug], termIDs[slug]); err != nil {
				return fmt.Errorf("关联术语 %s -> %s: %w", t.Slug, slug, err)
			}
		}
	}

	for _, e := range a.Edges {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO relation (from_id, to_id, kind, label, note) VALUES (?, ?, ?, ?, ?)`,
			personIDs[e.From], personIDs[e.To], e.Kind, e.Label, e.Note); err != nil {
			return fmt.Errorf("写入关系 %s->%s: %w", e.From, e.To, err)
		}
	}

	if err := buildFTS(ctx, tx, a, personIDs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务: %w", err)
	}
	return nil
}

func insertPerson(ctx context.Context, tx *sql.Tx, p content.Person, ord int) (int64, error) {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO person (slug, name, name_en, birth, death, birth_place, nationality,
		 epithet, portrait_key, summary, thesis, ord) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Slug, p.Name, p.NameEn, p.Birth, p.Death, p.BirthPlace,
		p.Nationality, p.Epithet, p.PortraitKey, p.Summary, p.Thesis, ord)
	if err != nil {
		return 0, fmt.Errorf("写入人物 %s: %w", p.Slug, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取人物 ID %s: %w", p.Slug, err)
	}
	return id, nil
}

func insertPersonDetail(ctx context.Context, tx *sql.Tx, id int64, p content.Person) error {
	for i, t := range p.Timeline {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO person_timeline (person_id, year, date_text, title, body, ord) VALUES (?, ?, ?, ?, ?, ?)`,
			id, t.Year, t.DateText, t.Title, t.Body, i); err != nil {
			return fmt.Errorf("写入 %s 年表: %w", p.Slug, err)
		}
	}
	for _, s := range p.Sections {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO person_section (person_id, kind, title, body_md, ord) VALUES (?, ?, ?, ?, ?)`,
			id, s.Kind, s.Title, s.BodyMD, s.Ord); err != nil {
			return fmt.Errorf("写入 %s 章节 %s: %w", p.Slug, s.Title, err)
		}
	}
	return nil
}

func insertWork(ctx context.Context, tx *sql.Tx, w content.Work, personID int64, ord int) (int64, error) {
	var pid any
	if personID != 0 {
		pid = personID
	}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO work (slug, person_id, title, year, category, summary, source, source_url)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		w.Slug, pid, w.Title, w.Year, w.Category, w.Summary, w.Source, w.SourceURL)
	if err != nil {
		return 0, fmt.Errorf("写入著作 %s: %w", w.Slug, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取著作 ID %s: %w", w.Slug, err)
	}
	return id, nil
}

// buildFTS 把可检索的文本灌入 FTS5 索引。
// body 中同时包含人名/出处等辅助词，以便跨类型检索命中。
func buildFTS(ctx context.Context, tx *sql.Tx, a *content.Archive, personIDs map[string]int64) error {
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO search_fts (title, body, kind, ref, person) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("准备 FTS 写入: %w", err)
	}
	defer stmt.Close()

	add := func(title, body, kind, ref, person string) error {
		if _, err := stmt.ExecContext(ctx, title, body, kind, ref, person); err != nil {
			return fmt.Errorf("写入 FTS (%s/%s): %w", kind, ref, err)
		}
		return nil
	}

	personName := make(map[string]string, len(a.People))
	for _, p := range a.People {
		personName[p.Slug] = p.Name
	}

	for _, p := range a.People {
		body := strings.Join([]string{p.Summary, p.Epithet, p.Thesis, p.BirthPlace, p.Nationality}, " ")
		if err := add(p.Name, body, "person", p.Slug, p.Name); err != nil {
			return err
		}
		for _, s := range p.Sections {
			// 贡献/思想章节单独建索引，检索结果可直达具体章节
			if s.Kind == "bio" || s.Kind == "other" {
				continue
			}
			if err := add(p.Name+" · "+s.Title, s.BodyMD, "person_section", p.Slug, p.Name); err != nil {
				return err
			}
		}
		for _, t := range p.Timeline {
			if err := add(fmt.Sprintf("%s %d %s", p.Name, t.Year, t.Title), t.Body,
				"person_section", p.Slug, p.Name); err != nil {
				return err
			}
		}
	}

	for _, w := range a.Works {
		name := personName[w.Person]
		var sb strings.Builder
		sb.WriteString(w.Summary)
		sb.WriteString(" ")
		sb.WriteString(w.Person + " " + name)
		for _, e := range w.Excerpts {
			sb.WriteString(" ")
			sb.WriteString(e.Text)
		}
		if err := add(w.Title, sb.String(), "work", w.Slug, name); err != nil {
			return err
		}
	}

	for _, e := range a.Events {
		var names []string
		for _, slug := range e.People {
			names = append(names, personName[slug])
		}
		body := e.Body + " " + e.Location + " " + e.Category + " " + strings.Join(names, " ")
		if err := add(e.Title, body, "event", fmt.Sprint(e.Year), strings.Join(names, "、")); err != nil {
			return err
		}
	}

	for _, t := range a.Terms {
		var names []string
		for _, slug := range t.People {
			names = append(names, personName[slug])
		}
		body := t.Definition + " " + t.Body + " " + strings.Join(t.Aliases, " ") + " " + strings.Join(names, " ")
		if err := add(t.Term, body, "term", t.Slug, strings.Join(names, "、")); err != nil {
			return err
		}
	}

	return nil
}
