package content

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// sectionKindMap 把二级标题映射到固定的领域键，保证前端可按政治/经济/文化/思想分组。
var sectionKindMap = []struct {
	keys  []string
	kind  string
	title string
	order int
}{
	{[]string{"政治贡献"}, "politics", "政治贡献", 1},
	{[]string{"经济贡献"}, "economy", "经济贡献", 2},
	{[]string{"文化贡献"}, "culture", "文化贡献", 3},
	{[]string{"思想体系"}, "thought", "思想体系", 4},
	{[]string{"争议与评价", "争议"}, "controversy", "争议与评价", 5},
	{[]string{"简介", "生平"}, "bio", "简介", 0},
}

// kindFor 依据章节标题判定领域键；未匹配的归为 other。
func kindFor(title string) (string, int) {
	for _, m := range sectionKindMap {
		for _, k := range m.keys {
			if strings.Contains(title, k) {
				return m.kind, m.order
			}
		}
	}
	return "other", 90
}

// splitSections 把 Markdown 正文按 "## " 二级标题切分为若干章节。
// 第一个标题之前的内容（若有）作为 preamble 归入 bio 章节。
func splitSections(md string) []Section {
	lines := strings.Split(md, "\n")
	type raw struct {
		title string
		body  []string
	}
	var raws []raw
	cur := -1 // 用下标而非指针：append 可能重新分配底层数组
	for _, ln := range lines {
		if t, ok := strings.CutPrefix(ln, "## "); ok {
			raws = append(raws, raw{title: strings.TrimSpace(t)})
			cur = len(raws) - 1
			continue
		}
		if cur < 0 {
			continue
		}
		raws[cur].body = append(raws[cur].body, ln)
	}

	out := make([]Section, 0, len(raws))
	for _, r := range raws {
		kind, ord := kindFor(r.title)
		out = append(out, Section{
			Kind:   kind,
			Title:  r.title,
			BodyMD: strings.TrimSpace(strings.Join(r.body, "\n")),
			Ord:    ord,
		})
	}
	slices.SortStableFunc(out, func(a, b Section) int { return a.Ord - b.Ord })
	return out
}

// parsePersonFile 解析 content/people/*.md：YAML front-matter + Markdown 正文。
func parsePersonFile(path string) (Person, error) {
	var p Person
	raw, err := os.ReadFile(path)
	if err != nil {
		return p, fmt.Errorf("读取人物档案 %s: %w", path, err)
	}
	fm, body, err := splitFrontMatter(string(raw))
	if err != nil {
		return p, fmt.Errorf("解析人物档案 %s: %w", path, err)
	}
	if err := yaml.Unmarshal([]byte(fm), &p); err != nil {
		return p, fmt.Errorf("解析人物 front-matter %s: %w", path, err)
	}
	if p.Slug == "" {
		p.Slug = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	p.Sections = splitSections(body)
	return p, nil
}

// splitFrontMatter 切出 --- 包裹的 YAML 头与其余正文。
// 分隔符必须独占一行（允许行尾空格），开头允许 UTF-8 BOM。
func splitFrontMatter(s string) (front, body string, err error) {
	s = strings.TrimPrefix(s, "\ufeff")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if !isFrontMatterDelim(s) {
		return "", s, nil // 无 front-matter，全部视作正文
	}
	rest := s[strings.IndexByte(s, '\n')+1:] // 跳过首行分隔符
	idx := findClosingDelim(rest)
	if idx < 0 {
		return "", "", fmt.Errorf("front-matter 未闭合")
	}
	front = strings.TrimSuffix(rest[:idx], "\n")
	body = rest[idx+len("\n---"):]
	// 丢弃分隔符行剩余的空白与换行。
	if nl := strings.IndexByte(body, '\n'); nl >= 0 {
		body = body[nl+1:]
	} else {
		body = ""
	}
	return front, body, nil
}

// isFrontMatterDelim 判断 s 的首行是否为 front-matter 起始分隔符。
func isFrontMatterDelim(s string) bool {
	line, _, _ := strings.Cut(s, "\n")
	return strings.TrimRight(line, " \t") == "---"
}

// findClosingDelim 找到独占一行的 --- 分隔符起点（其前一个字符为 \n）。
func findClosingDelim(s string) int {
	for i := 0; ; {
		j := strings.Index(s[i:], "\n---")
		if j < 0 {
			return -1
		}
		j += i
		lineEnd := strings.IndexByte(s[j+1:], '\n')
		tail := ""
		if lineEnd < 0 {
			tail = s[j+1:]
		} else {
			tail = s[j+1 : j+1+lineEnd]
		}
		if strings.TrimRight(tail, " \t") == "---" {
			return j
		}
		i = j + 1
	}
}

func readYAML[T any](path string) (T, error) {
	var zero T
	raw, err := os.ReadFile(path)
	if err != nil {
		return zero, fmt.Errorf("读取 %s: %w", path, err)
	}
	var v T
	if err := yaml.Unmarshal(raw, &v); err != nil {
		return zero, fmt.Errorf("解析 %s: %w", path, err)
	}
	return v, nil
}

// Load 读取 contentDir 下的全部内容并组装成 Archive。
func Load(contentDir string) (*Archive, error) {
	a := &Archive{}

	site, err := readYAML[Site](filepath.Join(contentDir, "site.yaml"))
	if err != nil {
		return nil, err
	}
	a.Site = site

	// 人物：按文件名排序保证顺序稳定
	personPaths, err := filepath.Glob(filepath.Join(contentDir, "people", "*.md"))
	if err != nil {
		return nil, fmt.Errorf("查找人物档案: %w", err)
	}
	slices.Sort(personPaths)
	for _, p := range personPaths {
		person, err := parsePersonFile(p)
		if err != nil {
			return nil, err
		}
		a.People = append(a.People, person)
	}
	if len(a.People) == 0 {
		return nil, fmt.Errorf("在 %s 下未找到任何人物档案", contentDir)
	}

	works, err := readYAML[WorksFile](filepath.Join(contentDir, "works.yaml"))
	if err != nil {
		return nil, err
	}
	a.Works = works.Works

	events, err := readYAML[EventsFile](filepath.Join(contentDir, "events", "timeline.yaml"))
	if err != nil {
		return nil, err
	}
	a.Events = events.Events

	terms, err := readYAML[TermsFile](filepath.Join(contentDir, "terms.yaml"))
	if err != nil {
		return nil, err
	}
	a.Terms = terms.Terms

	relations, err := readYAML[RelationsFile](filepath.Join(contentDir, "relations.yaml"))
	if err != nil {
		return nil, err
	}
	a.Edges = relations.Edges

	// 事件按年份稳定排序
	slices.SortStableFunc(a.Events, func(x, y Event) int { return x.Year - y.Year })

	if err := a.validate(); err != nil {
		return nil, err
	}
	return a, nil
}

// validate 检查引用完整性，尽早暴露内容错误而不是等到运行时。
func (a *Archive) validate() error {
	people := make(map[string]bool, len(a.People))
	for _, p := range a.People {
		if p.Slug == "" {
			return fmt.Errorf("人物缺少 slug（name=%q）", p.Name)
		}
		if p.Name == "" {
			return fmt.Errorf("人物缺少 name: %s", p.Slug)
		}
		if people[p.Slug] {
			return fmt.Errorf("人物 slug 重复: %s", p.Slug)
		}
		people[p.Slug] = true
	}
	if len(people) == 0 {
		return fmt.Errorf("未加载到任何人物档案")
	}

	works := make(map[string]bool, len(a.Works))
	for _, w := range a.Works {
		if w.Slug == "" {
			return fmt.Errorf("著作缺少 slug: %s", w.Title)
		}
		if w.Title == "" {
			return fmt.Errorf("著作缺少 title: %s", w.Slug)
		}
		if works[w.Slug] {
			return fmt.Errorf("著作 slug 重复: %s", w.Slug)
		}
		works[w.Slug] = true
		if !people[w.Person] {
			return fmt.Errorf("著作 %s 引用了不存在的人物 %s", w.Slug, w.Person)
		}
	}

	terms := make(map[string]bool, len(a.Terms))
	for _, t := range a.Terms {
		if t.Slug == "" {
			return fmt.Errorf("术语缺少 slug（term=%q）", t.Term)
		}
		if t.Term == "" {
			return fmt.Errorf("术语缺少 term: %s", t.Slug)
		}
		if terms[t.Slug] {
			return fmt.Errorf("术语 slug 重复: %s", t.Slug)
		}
		terms[t.Slug] = true
	}
	for _, t := range a.Terms {
		for _, r := range t.Related {
			if r == t.Slug {
				return fmt.Errorf("术语 %s 关联了自身", t.Slug)
			}
			if !terms[r] {
				return fmt.Errorf("术语 %s 关联了不存在的术语 %s", t.Slug, r)
			}
		}
		for _, s := range t.Sources {
			if !works[s] {
				return fmt.Errorf("术语 %s 引用了不存在的著作 %s", t.Slug, s)
			}
		}
		for _, p := range t.People {
			if !people[p] {
				return fmt.Errorf("术语 %s 引用了不存在的人物 %s", t.Slug, p)
			}
		}
	}

	for _, e := range a.Events {
		if e.Title == "" {
			return fmt.Errorf("事件缺少 title（year=%d）", e.Year)
		}
		for _, p := range e.People {
			if !people[p] {
				return fmt.Errorf("事件 %s 引用了不存在的人物 %s", e.Title, p)
			}
		}
	}

	for _, ed := range a.Edges {
		if ed.From == ed.To {
			return fmt.Errorf("关系边指向自身: %s", ed.From)
		}
		if !people[ed.From] || !people[ed.To] {
			return fmt.Errorf("关系边引用了不存在的人物: %s -> %s", ed.From, ed.To)
		}
	}
	return nil
}
