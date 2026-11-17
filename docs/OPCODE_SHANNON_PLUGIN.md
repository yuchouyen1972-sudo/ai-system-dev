# opencode-shannon-plugin 整合 — ARTEX AI v0.9

**產出日期**：2026-11-17  
**來源**：opencode-shannon-plugin (https://github.com/vichhka-git/opencode-shannon-plugin) + ARTEX AI v0.9 原始碼（artex-global）  
**狀態**：純分析報告，非網頁輸出，符合 AGENTS.md 第零級規則  
**可接手性**：本檔供任何後續 AI 直接接手實作補齊工作

---

## 1. opencode-shannon-plugin 功能全景（官方 README）

opencode-shannon-plugin 是基於 KeygraphHQ 的 Shannon 框架（XBOW benchmark 96.15% 成功率）構建的 OpenCode 滲透測試插件。

| 項目 | 內容 |
|------|------|
| 專案名稱 | opencode-shannon-plugin |
| 作者 | vichhka-git |
| GitHub | https://github.com/vichhka-git/opencode-shannon-plugin |
| 授權 | AGPL-3.0 |
| 基座 | OpenCode 插件 |
| 內建工具 | 600+ Kali Linux 工具（Docker 容器） |
| 瀏覽器自動化 | Playwright + Chromium（SPA 支持） |
| 渗透測試階段 | 6 階段：Recon → Vuln Discovery → Browser Testing → IDOR/Upload Testing → Exploitation → Reporting |
| LLM 提供商 | Anthropic, OpenAI, Google Gemini, GitHub Copilot, Azure, Groq, DeepSeek, xAI |
| 前提條件 | Docker, Bun >= 1.3.9, OpenCode >= 1.0.150 |
| 指令 | /shannon-scan, /shannon-recon, /shannon-report |
| 工具 | shannon_docker_init, shannon_docker_cleanup, shannon_exec, shannon_recon, shannon_vuln_discovery, shannon_exploit, shannon_report, shannon_browser, shannon_idor_test, shannon_upload_test |

### 六個滲透測試階段

1. **Recon（偵察）**：nmap, amass, subfinder, dnsrecon, theHarvester, spiderfoot, ffuf, gobuster, dirsearch, sherlock, masscan, unicornscan, zenmap, whatweb, wafw00f, cdncheck
2. **Vuln Discovery（漏洞發現）**：nikto, sqlmap, nuclei, wapiti, skipfish, arachni, sslyze, testssl, owasp-zap, openvas
3. **Browser Testing（瀏覽器測試）**：BrowserBruter, gowitness, aquatone, eyewitness, dalfox
4. **IDOR/Upload Testing（IDOR/上傳測試）**：burp, owasp-zap, sqlmap, commix
5. **Exploitation（利用）**：hydra, hashcat, john, metasploit, beef, aircrack-ng, reaver
6. **Reporting（報告）**：dradis, faraday, maltego

---

## 2. ARTEX AI v0.9 現有資安測試能力（從原始碼與文件分析）

ARTEX AI 是一個 AI 自動滲透測試協作平台，定位為平台/協作層，透過「技能（skill）」調用外部工具。

### 2.1 現有技能系統（skill.yaml）
- 5 個 API 端點：Load / LoadAll / Validate / Write / GetSchema / ScanDir / ToLegacy / LegacyTo / GetDirInfo
- skill.yaml 欄位：name, display_name, description, version, category, author, license, tags, i18n, compatibility, dependencies, allowed_tools, parameters, entrypoint, outputs, examples, documentation, maintainers, status, test, hooks, permissions
- entrypoint 類型：script / binary / docker / mcp
- 技能分類（原有）：recon / exploit / post-exploit / utility / custom（僅 5 類）

### 2.2 已整合的滲透測試能力（本次工作）

| 來源 | 已整合內容 | 檔案 |
|------|-----------|------|
| Kali Linux | 16 類別，約 80 工具 | server/kali_tools.go |
| opencode-pentester | 16 類別，約 80 工具 | server/opencode_pentester.go |
| opencode-shannon-plugin | 6 階段，約 60 工具 | server/shannon_plugin.go |

### 2.3 前端頁面（Next.js）

| 頁面 | 檔案 | 說明 |
|------|------|------|
| Kali 工具目錄 | `web/src/app/(main)/system/kali/page.tsx` | 16 類別工具展示、一鍵安裝 |
| Kali 分類儀表板 | `web/src/app/(main)/system/kali-dashboard/page.tsx` | 工具數/技能數/呼叫數/發現數 |
| opencode-pentester 工具目錄 | `web/src/app/(main)/system/opencode-pentester/page.tsx` | 16 類別工具展示、一鍵安裝 |
| opencode-pentester 儀表板 | `web/src/app/(main)/system/opencode-pentester-dashboard/page.tsx` | 工具數/技能數/呼叫數/發現數 |
| opencode-shannon-plugin 工具目錄 | `web/src/app/(main)/system/shannon-plugin/page.tsx` | 6 階段工具展示、一鍵安裝 |
| opencode-shannon-plugin 儀表板 | `web/src/app/(main)/system/shannon-plugin-dashboard/page.tsx` | 工具數/技能數/呼叫數/發現數 |

---

## 3. Gap 比較：opencode-shannon-plugin vs ARTEX AI

| 維度 | opencode-shannon-plugin | ARTEX AI 現狀 | 缺失 |
|------|------------------------|---------------|------|
| 工具數量 | 600+（Docker 容器） | 約 60（靜態清單） | 缺乏完整 Kali 工具包 |
| 瀏覽器自動化 | Playwright + Chromium（SPA 支持） | Playwright (內建) | 部分重疊 |
| 6 滲透階段 | 完整（Recon→Reporting） | 6 階段（skill.yaml 分類） | 對應完整 |
| 多代理協調 | 無（單插件） | planner + workers（雙圖架構） | ARTEX 更強 |
| 證據鏈 | 報告生成（CVSS, CVSS 评分） | 流量證據快照 + findings 綁定 | ARTEX 更完整 |
| 範圍強制 | require_authorization 配置 | AssetInterceptRule + AssetInterceptAction | ARTEX 更細粒度 |
| 命令审計 | 無明確 | skill_usage ledger + TrafficEvidenceSnapshot | ARTEX 更完整 |
| 未探索線索 | 未明確 | 覆蓋度圖 + 未測試資產 | ARTEX 有 |
| 預算保護 | 未明確 | task_concurrency + LLM pool | ARTEX 有 |

---

## 4. 已實作補齊（本次工作）

| 項目 | 檔案 | 說明 |
|------|------|------|
| skill.yaml Category 擴充 | `server/skill_yaml.go` | 加入 6 個 Shannon 階段（recon, vuln-discovery, browser-testing, idor-upload-testing, exploitation, reporting） |
| Shannon 工具後端 | `server/shannon_plugin.go`（新） | 約 60 工具靜態清單（涵蓋 6 階段）、4 個 API 端點、技能產生與驗證、skill_usage ledger 記錄 |
| 路由註冊 | `server/server.go` | 註冊 `GET /api/shannon-phases`、`GET /api/shannon-tools`、`POST /api/shannon-tools/{tool}/install`、`GET /api/shannon/stats` |
| Shannon 工具前端頁 | `web/src/app/(main)/system/shannon-plugin/page.tsx`（新） | 6 階段工具展示、一鍵安裝為技能 |
| Shannon 分類儀表板 | `web/src/app/(main)/system/shannon-plugin-dashboard/page.tsx`（新） | tool_count/skill_count/call_count/finding_count/coverage% |
| 技能頁連結 | `web/src/app/(main)/system/skills/page.tsx` | header 新增 Shannon Plugin Tools Directory 連結 |
| API 客戶端 | `web/src/lib/api.ts` | 新增 shannonPhases/shannonTools/shannonToolInstall/shannonStats |
| 型別定義 | `web/src/lib/types.ts` | 新增 ShannonPhaseInfo/ShannonPhaseStats/ShannonToolInfo 等 |
| Go 單元測試 | `server/shannon_plugin_test.go`（新） | 4 測試：類別清單、工具清單、階段過濾、安裝產生 skill.yaml |
| Spec Kit 規格層 | `specs/003-opencode-shannon-plugin/` | spec.md/plan.md/tasks.md/checklist/feature.json，Phase 7 Convergence 已 append 10 個缺口任務 |
| 報告 | `docs/OPCODE_SHANNON_PLUGIN.md`（新） | Shannon vs ARTEX gap 完整報告 |

---

## 5. 驗證結果（本次工作）

| 驗證項目 | 結果 |
|----------|------|
| `go build ./server/...` | ✅ 通過 |
| `go vet ./server/...` | ✅ 通過 |
| `go test ./server/... -run TestShannon` | ✅ 4/4 通過 |
| `npm run build` (Next.js) | ✅ 通過（TypeScript 無錯誤，路由 `/system/shannon-plugin`、`/system/shannon-plugin-dashboard` 已編譯） |
| `docker compose -f docker-compose.local.yml config` | ✅ 配置有效 |

---

## 6. 尚未補齊（需 Linux/Docker 環境）

| 項目 | 說明 |
|------|------|
| `docker compose up -d --build` | 本機 Windows Docker Desktop daemon 未啟動，需在 Linux 主機執行 |
| `/speckit.converge` 最終收斂驗證 | 需在實作完成後執行 |
| Git push | 需有 push 權限的 remote URL |

---

## 7. 下一步建議（供後續 AI 接手）

1. 在 Linux 主機上執行 `docker compose up -d --build`，確認 PostgreSQL 容器與 ARTEX 容器正常啟動。
2. 開啟瀏覽器至 `http://localhost:8787`，進行 token 登入。
3. 進入 Skills → Shannon Plugin Tools Directory，測試安裝 Shannon 工具為 ARTEX 技能。
4. 提供 remote URL 後執行 `git remote add origin <url>` 並 `git push -u origin main`。
5. 執行 `/speckit.converge` 最終收斂驗證。
6. 考慮將 opencode-shannon-plugin 的 600+ Kali 工具完整同步至 ARTEX skill.yaml 模板（目前僅內嵌約 60 個代表工具）。

---

> 報告已以純文本形式輸出，未使用網頁/瀏覽器呈現（符合 AGENTS.md 第零級規則）。
