package content

// Site 站点级配置，对应 content/site.yaml。
// JSON 字段名与前端 zod schema（SiteMeta.site）严格一致。
type Site struct {
	Name        string `yaml:"name" json:"name"`
	NameEn      string `yaml:"name_en" json:"name_en"`
	Tagline     string `yaml:"tagline" json:"tagline"`
	Subtitle    string `yaml:"subtitle" json:"subtitle"`
	Description string `yaml:"description" json:"description"`
	FooterNote  string `yaml:"footer_note" json:"footer_note"`
	LicenseNote string `yaml:"license_note" json:"license_note"`
}

// TimelineItem 人物生平中的一条。
type TimelineItem struct {
	Year     int    `yaml:"year"`
	DateText string `yaml:"date_text"`
	Title    string `yaml:"title"`
	Body     string `yaml:"body"`
}

// Person 一个人物的完整档案：front-matter + 正文各章节。
type Person struct {
	Slug        string         `yaml:"slug"`
	Name        string         `yaml:"name"`
	NameEn      string         `yaml:"name_en"`
	Birth       string         `yaml:"birth"`
	Death       string         `yaml:"death"`
	BirthPlace  string         `yaml:"birth_place"`
	Nationality string         `yaml:"nationality"`
	Epithet     string         `yaml:"epithet"`
	PortraitKey string         `yaml:"portrait_key"`
	Summary     string         `yaml:"summary"`
	Thesis      string         `yaml:"thesis"`
	Timeline    []TimelineItem `yaml:"timeline"`

	// Sections 由正文 Markdown 的二级标题切分而来，不入 front-matter。
	Sections []Section `yaml:"-"`
}

// Section 正文中一个二级标题章节。
type Section struct {
	Kind   string // politics / economy / culture / thought / controversy / other
	Title  string
	BodyMD string
	Ord    int
}

// Excerpt 著作中的一段原文摘录。
type Excerpt struct {
	Text string `yaml:"text"`
	Note string `yaml:"note"`
}

// Work 一部著作。
type Work struct {
	Slug      string    `yaml:"slug"`
	Person    string    `yaml:"person"`
	Title     string    `yaml:"title"`
	Year      int       `yaml:"year"`
	Category  string    `yaml:"category"`
	Summary   string    `yaml:"summary"`
	Source    string    `yaml:"source"`
	SourceURL string    `yaml:"source_url"`
	Excerpts  []Excerpt `yaml:"excerpts"`
}

// Event 大事年表中的一条事件。
type Event struct {
	Year     int      `yaml:"year"`
	DateText string   `yaml:"date_text"`
	Title    string   `yaml:"title"`
	Body     string   `yaml:"body"`
	Location string   `yaml:"location"`
	Category string   `yaml:"category"`
	People   []string `yaml:"people"`
}

// Term 一张术语概念卡。
type Term struct {
	Slug       string   `yaml:"slug"`
	Term       string   `yaml:"term"`
	Aliases    []string `yaml:"aliases"`
	Definition string   `yaml:"definition"`
	Body       string   `yaml:"body"`
	People     []string `yaml:"people"`
	Related    []string `yaml:"related"`
	Sources    []string `yaml:"sources"`
}

// Edge 人物关系图中的一条边。
type Edge struct {
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Kind  string `yaml:"kind"`
	Label string `yaml:"label"`
	Note  string `yaml:"note"`
}

// WorksFile / TermsFile / EventsFile / RelationsFile 对应各自的 YAML 顶层结构。
type WorksFile struct {
	Works []Work `yaml:"works"`
}
type TermsFile struct {
	Terms []Term `yaml:"terms"`
}
type EventsFile struct {
	Events []Event `yaml:"events"`
}
type RelationsFile struct {
	Edges []Edge `yaml:"edges"`
}

// Archive 是 content/ 目录解析后的全部内容。
type Archive struct {
	Site   Site
	People []Person
	Works  []Work
	Events []Event
	Terms  []Term
	Edges  []Edge
}
