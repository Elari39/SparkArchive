package content

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSplitFrontMatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantFront string
		wantBody  string
		wantErr   bool
	}{
		{
			name:      "标准 front-matter",
			input:     "---\nslug: a\n---\n正文\n",
			wantFront: "slug: a",
			wantBody:  "正文\n",
		},
		{
			name:     "无 front-matter",
			input:    "只有正文",
			wantBody: "只有正文",
		},
		{
			name:    "未闭合",
			input:   "---\nslug: a\n正文",
			wantErr: true,
		},
		{
			name: "CRLF 归一化",
			// 前后换行统一归一化为 LF，正文保留其后的换行
			input:     "---\r\nslug: a\r\n---\r\n正文\r\n",
			wantFront: "slug: a",
			wantBody:  "正文\n",
		},
		{
			name:      "BOM 前缀",
			input:     "\ufeff---\nslug: a\n---\n正文\n",
			wantFront: "slug: a",
			wantBody:  "正文\n",
		},
		{
			name:      "分隔符行尾空格",
			input:     "--- \nslug: a\n---\t\n正文\n",
			wantFront: "slug: a",
			wantBody:  "正文\n",
		},
		{
			name: "正文内 --- 不误判",
			// 闭合并未出现，返回未闭合错误而非把正文中的分隔符当结尾
			input:   "---\nslug: a\n\n---not-a-delim\n正文",
			wantErr: true,
		},
		{
			name: "首行非分隔符但有 --- 前缀",
			// ---- 不是分隔符，整体视作正文
			input:    "----\n正文\n---\n",
			wantBody: "----\n正文\n---\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			front, body, err := splitFrontMatter(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatal("期望报错，实际为 nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("意外错误: %v", err)
			}
			if front != tc.wantFront {
				t.Errorf("front = %q, 期望 %q", front, tc.wantFront)
			}
			if body != tc.wantBody {
				t.Errorf("body = %q, 期望 %q", body, tc.wantBody)
			}
		})
	}
}

func TestSplitSections(t *testing.T) {
	t.Parallel()

	md := "## 生平\n\n生平内容。\n\n## 政治贡献\n\n政治内容。\n\n## 思想体系\n\n思想内容。\n\n## 争议与评价\n\n争议内容。\n"
	got := splitSections(md)
	if len(got) != 4 {
		t.Fatalf("章节数 = %d, 期望 4", len(got))
	}
	wantKinds := []string{"bio", "politics", "thought", "controversy"}
	for i, want := range wantKinds {
		if got[i].Kind != want {
			t.Errorf("第 %d 节 Kind = %q, 期望 %q", i, got[i].Kind, want)
		}
	}
	if !strings.Contains(got[1].BodyMD, "政治内容") {
		t.Errorf("政治贡献正文未正确提取: %q", got[1].BodyMD)
	}
}

func TestKindFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title string
		want  string
	}{
		{"政治贡献", "politics"},
		{"经济贡献", "economy"},
		{"文化贡献", "culture"},
		{"思想体系", "thought"},
		{"争议与评价", "controversy"},
		{"其他随机标题", "other"},
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			got, _ := kindFor(tc.title)
			if got != tc.want {
				t.Errorf("kindFor(%q) = %q, 期望 %q", tc.title, got, tc.want)
			}
		})
	}
}

// TestValidateBoundaries 覆盖边界校验：空 slug/name、自环边、空档案。
func TestValidateBoundaries(t *testing.T) {
	t.Parallel()

	base := func() *Archive {
		return &Archive{
			People: []Person{{Slug: "a", Name: "甲"}},
		}
	}

	tests := []struct {
		name    string
		mutate  func(*Archive)
		wantErr bool
	}{
		{
			name:    "空档案",
			mutate:  func(a *Archive) { a.People = nil },
			wantErr: true,
		},
		{
			name:    "人物缺 slug",
			mutate:  func(a *Archive) { a.People[0].Slug = "" },
			wantErr: true,
		},
		{
			name:    "人物缺 name",
			mutate:  func(a *Archive) { a.People[0].Name = "" },
			wantErr: true,
		},
		{
			name: "术语关联自身",
			mutate: func(a *Archive) {
				a.Terms = []Term{{Slug: "t", Term: "词", Related: []string{"t"}}}
			},
			wantErr: true,
		},
		{
			name: "关系边自环",
			mutate: func(a *Archive) {
				a.Edges = []Edge{{From: "a", To: "a"}}
			},
			wantErr: true,
		},
		{
			name: "著作缺 person",
			mutate: func(a *Archive) {
				a.Works = []Work{{Slug: "w", Title: "作", Person: ""}}
			},
			wantErr: true,
		},
		{
			name:   "合法最小档案",
			mutate: func(*Archive) {},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := base()
			tc.mutate(a)
			err := a.validate()
			if tc.wantErr && err == nil {
				t.Fatal("期望报错，实际为 nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("意外错误: %v", err)
			}
		})
	}
}

// TestLoadRealContent 校验仓库内真实内容的引用完整性。
func TestLoadRealContent(t *testing.T) {
	t.Parallel()

	dir := filepath.Join("..", "..", "..", "content")
	a, err := Load(dir)
	if err != nil {
		t.Fatalf("加载真实内容失败: %v", err)
	}
	if len(a.People) != 15 {
		t.Errorf("人物数 = %d, 期望 15", len(a.People))
	}
	if len(a.Works) == 0 || len(a.Events) == 0 || len(a.Terms) == 0 || len(a.Edges) == 0 {
		t.Fatalf("内容不完整: works=%d events=%d terms=%d edges=%d",
			len(a.Works), len(a.Events), len(a.Terms), len(a.Edges))
	}

	// 每个人物都必须具备四域贡献与争议章节
	for _, p := range a.People {
		seen := map[string]bool{}
		for _, s := range p.Sections {
			seen[s.Kind] = true
		}
		for _, k := range []string{"politics", "economy", "culture", "thought", "controversy"} {
			if !seen[k] {
				t.Errorf("人物 %s 缺少 %s 章节", p.Slug, k)
			}
		}
	}
}
