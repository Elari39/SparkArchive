#!/usr/bin/env pwsh
# 星火档案馆 端到端冒烟测试
# 用法: pwsh -File scripts/smoke.ps1 [-BaseUrl http://localhost:12026]
#
# 两点注意事项：
# 1) 统一用 Invoke-WebRequest + ConvertFrom-Json 而非 Invoke-RestMethod。
#    PowerShell 7 的 Invoke-RestMethod 对 JSON 数组有展开（unrolling）行为，
#    会导致 @(Invoke-RestMethod ...).Count 恒为 1，从而产生假失败。
# 2) 默认地址用 127.0.0.1 而非 localhost。
#    Docker Desktop 的端口代理只监听 IPv4，但 localhost 会优先解析到 ::1，
#    .NET 的 HttpClient 会先尝试 ::1 并等待约 21 秒超时才回退到 IPv4，
#    使整轮测试耗时从数秒膨胀到十几分钟。
param(
    [string]$BaseUrl = "http://127.0.0.1:12026"
)

$ErrorActionPreference = "Stop"
$script:fail = 0
$script:pass = 0

function Check {
    param([string]$Name, [scriptblock]$Test)
    try {
        if (& $Test) {
            Write-Host "  [PASS] $Name" -ForegroundColor Green
            $script:pass++
        } else {
            Write-Host "  [FAIL] $Name" -ForegroundColor Red
            $script:fail++
        }
    } catch {
        Write-Host "  [FAIL] $Name -- $_" -ForegroundColor Red
        $script:fail++
    }
}

# 取回并解析 JSON（数组不会丢失长度信息）
function Api {
    param([string]$Path)
    $resp = Invoke-WebRequest "$BaseUrl$Path" -UseBasicParsing
    if ($resp.Content.Trim() -eq "") { return @() }
    return ($resp.Content | ConvertFrom-Json)
}

function Esc { param([string]$s) [uri]::EscapeDataString($s) }

Write-Host ""
Write-Host "星火档案馆 冒烟测试 -> $BaseUrl" -ForegroundColor Cyan

Write-Host ""
Write-Host "API 端点" -ForegroundColor Yellow
Check "GET /api/health 返回 200" {
    (Invoke-WebRequest "$BaseUrl/api/health" -UseBasicParsing).StatusCode -eq 200
}
Check "GET /api/meta 人物数为 15" {
    $m = Api "/api/meta"
    $m.counts.people -eq 15 -and $m.counts.works -gt 0 -and $m.counts.events -gt 0 -and $m.counts.terms -gt 0
}
Check "GET /api/people 返回 15 人" {
    $p = Api "/api/people"
    $p.Count -eq 15
}
Check "GET /api/people/mao 含四域贡献与争议章节" {
    $p = Api "/api/people/mao"
    $kinds = $p.sections.kind
    ($kinds -contains "politics") -and ($kinds -contains "economy") -and
    ($kinds -contains "culture") -and ($kinds -contains "thought") -and
    ($kinds -contains "controversy")
}
Check "GET /api/people/marx 含生平年表" {
    $p = Api "/api/people/marx"
    $p.timeline.Count -gt 5
}
Check "十五个人物都含全部四域贡献" {
    $ok = $true
    foreach ($slug in @("marx", "engels", "plekhanov", "zetkin", "katayama", "lenin", "luxemburg", "kollontai", "li-dazhao", "ho-chi-minh", "gramsci", "mao", "castro", "guevara", "sankara")) {
        $k = (Api "/api/people/$slug").sections.kind
        foreach ($need in @("politics", "economy", "culture", "thought")) {
            if ($k -notcontains $need) { $ok = $false; Write-Host "      $slug 缺少 $need" -ForegroundColor DarkYellow }
        }
    }
    $ok
}
Check "GET /api/works 非空" {
    (Api "/api/works").Count -gt 0
}
Check "GET /api/works/imperialism 含摘录" {
    (Api "/api/works/imperialism").excerpts.Count -gt 0
}
Check "GET /api/events 非空" {
    (Api "/api/events").Count -gt 0
}
Check "GET /api/events?person=mao 筛选生效" {
    $all = (Api "/api/events").Count
    $f = (Api "/api/events?person=mao").Count
    $f -gt 0 -and $f -lt $all
}
Check "GET /api/events?category=理论 筛选生效" {
    $all = (Api "/api/events").Count
    $f = (Api "/api/events?category=$([uri]::EscapeDataString('理论'))").Count
    $f -gt 0 -and $f -lt $all
}
Check "GET /api/terms 非空" {
    (Api "/api/terms").Count -gt 0
}
Check "GET /api/terms/proletariat 含关联与出处" {
    $t = Api "/api/terms/proletariat"
    $t.related.Count -gt 0 -and $t.sources.Count -gt 0 -and $t.people.Count -gt 0
}
Check "GET /api/relations 有 15 节点与边" {
    $g = Api "/api/relations"
    $g.nodes.Count -eq 15 -and $g.edges.Count -gt 0
}
Check "GET /api/works?person=lenin 筛选生效" {
    $all = (Api "/api/works").Count
    $f = (Api "/api/works?person=lenin").Count
    $f -gt 0 -and $f -lt $all
}

Write-Host ""
Write-Host "中文全文检索 (FTS5 trigram + LIKE 降级)" -ForegroundColor Yellow
Check "长词 '无产阶级' 跨类型命中" {
    $r = Api "/api/search?q=$(Esc '无产阶级')"
    $kinds = @($r.items.kind | Select-Object -Unique)
    $r.total -gt 3 -and $kinds.Count -ge 2
}
Check "长词 '十月革命' 命中" {
    (Api "/api/search?q=$(Esc '十月革命')").total -gt 0
}
Check "长词 '剩余价值' 命中" {
    (Api "/api/search?q=$(Esc '剩余价值')").total -gt 0
}
Check "短词 '列宁' 走 LIKE 降级仍命中" {
    (Api "/api/search?q=$(Esc '列宁')").total -gt 0
}
Check "短词 '游击' 命中" {
    (Api "/api/search?q=$(Esc '游击')").total -gt 0
}
Check "短词 '古巴' 命中" {
    (Api "/api/search?q=$(Esc '古巴')").total -gt 0
}
Check "短词 '群众' 命中" {
    (Api "/api/search?q=$(Esc '群众')").total -gt 0
}
Check "检索结果含摘要文本" {
    $r = Api "/api/search?q=$(Esc '无产阶级')"
    -not [string]::IsNullOrWhiteSpace($r.items[0].snippet)
}
Check "类型过滤 type=term 生效" {
    $r = Api "/api/search?q=$(Esc '无产阶级')&type=term"
    @($r.items.kind | Select-Object -Unique) -eq "term"
}
Check "分页 limit 生效" {
    $r = Api "/api/search?q=$(Esc '无产阶级')&limit=2"
    $r.items.Count -le 2
}
Check "空查询返回 0 且不报错" {
    (Api "/api/search?q=").total -eq 0
}
Check "无处可匹配的关键词返回 0" {
    (Api "/api/search?q=zzzznonexistentzzz").total -eq 0
}
Check "SQL 注入串不导致 500 且数据完好" {
    $before = (Api "/api/people").Count
    $r = Api "/api/search?q=$(Esc "'; DROP TABLE person;--")"
    $after = (Api "/api/people").Count
    $r.total -ge 0 -and $before -eq 15 -and $after -eq 15
}
Check "FTS 语法字符不导致 500" {
    $ok = $true
    foreach ($q in @('"', 'NEAR(', '* AND ()', '%', '_')) {
        try { Api "/api/search?q=$(Esc $q)" | Out-Null } catch { $ok = $false }
    }
    $ok
}

Write-Host ""
Write-Host "错误处理" -ForegroundColor Yellow
Check "不存在的人物返回 404" {
    try { Invoke-WebRequest "$BaseUrl/api/people/nobody" -UseBasicParsing | Out-Null; $false }
    catch { $_.Exception.Response.StatusCode.value__ -eq 404 }
}
Check "不存在的著作返回 404" {
    try { Invoke-WebRequest "$BaseUrl/api/works/nobody" -UseBasicParsing | Out-Null; $false }
    catch { $_.Exception.Response.StatusCode.value__ -eq 404 }
}
Check "不存在的术语返回 404" {
    try { Invoke-WebRequest "$BaseUrl/api/terms/nobody" -UseBasicParsing | Out-Null; $false }
    catch { $_.Exception.Response.StatusCode.value__ -eq 404 }
}

Write-Host ""
Write-Host "前端路由 (SPA 回退)" -ForegroundColor Yellow
foreach ($route in @("/", "/people", "/people/mao", "/works", "/works/imperialism", "/timeline", "/glossary", "/glossary/proletariat", "/graph", "/search", "/about")) {
    Check "GET $route 返回应用外壳" {
        $r = Invoke-WebRequest "$BaseUrl$route" -UseBasicParsing
        $r.StatusCode -eq 200 -and $r.Content -match 'id="app"'
    }
}
Check "未知路径回退到 SPA" {
    $r = Invoke-WebRequest "$BaseUrl/no/such/route" -UseBasicParsing
    $r.StatusCode -eq 200 -and $r.Content -match 'id="app"'
}
Check "静态资源带长期缓存头" {
    $html = (Invoke-WebRequest "$BaseUrl/" -UseBasicParsing).Content
    $m = [regex]::Match($html, '/assets/[^"]+\.js')
    if (-not $m.Success) { return $false }
    $asset = Invoke-WebRequest "$BaseUrl$($m.Value)" -UseBasicParsing
    $asset.Headers["Cache-Control"] -match "immutable"
}

Write-Host ""
Write-Host "人物肖像资产 (历史照片)" -ForegroundColor Yellow
# 15 位人物现已全部配图，故逐个断言双格式可用；来源与许可见 public/assets/portraits/CREDITS.md。
foreach ($slug in @("marx", "engels", "plekhanov", "zetkin", "katayama", "lenin", "luxemburg", "kollontai", "li-dazhao", "ho-chi-minh", "gramsci", "mao", "castro", "guevara", "sankara")) {
    Check "GET /assets/portraits/$slug.{webp,jpg} 双格式可用" {
        $w = Invoke-WebRequest "$BaseUrl/assets/portraits/$slug.webp" -UseBasicParsing
        $j = Invoke-WebRequest "$BaseUrl/assets/portraits/$slug.jpg" -UseBasicParsing
        $w.StatusCode -eq 200 -and $w.Headers["Content-Type"] -match "image/webp" -and
        $j.StatusCode -eq 200 -and $j.Headers["Content-Type"] -match "image/jpeg"
    }
}
Check "肖像资产命中 /assets/ 长期缓存 (immutable)" {
    $r = Invoke-WebRequest "$BaseUrl/assets/portraits/marx.webp" -UseBasicParsing
    $r.Headers["Cache-Control"] -match "immutable"
}
Check "GET /favicon.svg 返回 200" {
    (Invoke-WebRequest "$BaseUrl/favicon.svg" -UseBasicParsing).StatusCode -eq 200
}
Check "GET /og-cover.png 返回 200 且为 image/png" {
    $r = Invoke-WebRequest "$BaseUrl/og-cover.png" -UseBasicParsing
    $r.StatusCode -eq 200 -and $r.Headers["Content-Type"] -match "image/png"
}

Write-Host ""
Write-Host "----------------------------------------"
if ($script:fail -eq 0) {
    Write-Host "全部通过: $($script:pass) 项" -ForegroundColor Green
    exit 0
} else {
    Write-Host "通过 $($script:pass) 项, 失败 $($script:fail) 项" -ForegroundColor Red
    exit 1
}
