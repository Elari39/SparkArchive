-- 星火档案馆 数据库结构
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE person (
  id          INTEGER PRIMARY KEY,
  slug        TEXT NOT NULL UNIQUE,
  name        TEXT NOT NULL,
  name_en     TEXT NOT NULL DEFAULT '',
  birth       TEXT NOT NULL DEFAULT '',
  death       TEXT NOT NULL DEFAULT '',
  birth_place TEXT NOT NULL DEFAULT '',
  nationality TEXT NOT NULL DEFAULT '',
  epithet     TEXT NOT NULL DEFAULT '',
  portrait_key TEXT NOT NULL DEFAULT '',
  summary     TEXT NOT NULL DEFAULT '',
  thesis      TEXT NOT NULL DEFAULT '',
  ord         INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE person_section (
  id        INTEGER PRIMARY KEY,
  person_id INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  kind      TEXT NOT NULL,
  title     TEXT NOT NULL,
  body_md   TEXT NOT NULL,
  ord       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_section_person ON person_section(person_id);

CREATE TABLE person_timeline (
  id        INTEGER PRIMARY KEY,
  person_id INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  year      INTEGER NOT NULL,
  date_text TEXT NOT NULL DEFAULT '',
  title     TEXT NOT NULL,
  body      TEXT NOT NULL DEFAULT '',
  ord       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_timeline_person ON person_timeline(person_id);

CREATE TABLE work (
  id         INTEGER PRIMARY KEY,
  slug       TEXT NOT NULL UNIQUE,
  person_id  INTEGER REFERENCES person(id) ON DELETE SET NULL,
  title      TEXT NOT NULL,
  year       INTEGER NOT NULL,
  category   TEXT NOT NULL DEFAULT '',
  summary    TEXT NOT NULL DEFAULT '',
  source     TEXT NOT NULL DEFAULT '',
  source_url TEXT NOT NULL DEFAULT ''
);
-- 按人物查著作（人物详情页侧栏、/api/works?person=）
CREATE INDEX idx_work_person ON work(person_id);

CREATE TABLE work_excerpt (
  id      INTEGER PRIMARY KEY,
  work_id INTEGER NOT NULL REFERENCES work(id) ON DELETE CASCADE,
  text    TEXT NOT NULL,
  note    TEXT NOT NULL DEFAULT '',
  ord     INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_excerpt_work ON work_excerpt(work_id);

CREATE TABLE event (
  id        INTEGER PRIMARY KEY,
  year      INTEGER NOT NULL,
  date_text TEXT NOT NULL DEFAULT '',
  title     TEXT NOT NULL,
  body      TEXT NOT NULL DEFAULT '',
  location  TEXT NOT NULL DEFAULT '',
  category  TEXT NOT NULL DEFAULT ''
);
-- 年表按年份排序 + 按范畴筛选（/api/events?category=）
CREATE INDEX idx_event_year ON event(year);
CREATE INDEX idx_event_category ON event(category);

CREATE TABLE event_person (
  event_id  INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE,
  person_id INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  PRIMARY KEY (event_id, person_id)
);
-- 复合主键已覆盖 event_id 前缀；按 person_id 反查事件需单独索引
CREATE INDEX idx_event_person_person ON event_person(person_id);

CREATE TABLE term (
  id         INTEGER PRIMARY KEY,
  slug       TEXT NOT NULL UNIQUE,
  term       TEXT NOT NULL,
  aliases    TEXT NOT NULL DEFAULT '',
  definition TEXT NOT NULL DEFAULT '',
  body       TEXT NOT NULL DEFAULT ''
);

CREATE TABLE term_person (
  term_id   INTEGER NOT NULL REFERENCES term(id) ON DELETE CASCADE,
  person_id INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  PRIMARY KEY (term_id, person_id)
);
-- 复合主键已覆盖 term_id 前缀；按 person_id 反查术语需单独索引
CREATE INDEX idx_term_person_person ON term_person(person_id);

CREATE TABLE term_related (
  term_id    INTEGER NOT NULL REFERENCES term(id) ON DELETE CASCADE,
  related_id INTEGER NOT NULL REFERENCES term(id) ON DELETE CASCADE,
  PRIMARY KEY (term_id, related_id)
);
-- 复合主键已覆盖 term_id 前缀；按 related_id 反查需单独索引
CREATE INDEX idx_term_related_related ON term_related(related_id);

CREATE TABLE term_work (
  term_id INTEGER NOT NULL REFERENCES term(id) ON DELETE CASCADE,
  work_id INTEGER NOT NULL REFERENCES work(id) ON DELETE CASCADE,
  PRIMARY KEY (term_id, work_id)
);
-- 复合主键已覆盖 term_id 前缀；按 work_id 反查需单独索引
CREATE INDEX idx_term_work_work ON term_work(work_id);

CREATE TABLE relation (
  id        INTEGER PRIMARY KEY,
  from_id   INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  to_id     INTEGER NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  kind      TEXT NOT NULL DEFAULT '',
  label     TEXT NOT NULL DEFAULT '',
  note      TEXT NOT NULL DEFAULT ''
);
-- 关系图 JOIN 两端人物
CREATE INDEX idx_relation_from ON relation(from_id);
CREATE INDEX idx_relation_to ON relation(to_id);

-- 全文索引：中文必须用 trigram 分词器。
-- 默认的 unicode61 不切分 CJK，整段会被当作单一 token，导致中文检索失效。
CREATE VIRTUAL TABLE search_fts USING fts5(
  title,
  body,
  kind UNINDEXED,
  ref UNINDEXED,
  person UNINDEXED,
  tokenize='trigram'
);
