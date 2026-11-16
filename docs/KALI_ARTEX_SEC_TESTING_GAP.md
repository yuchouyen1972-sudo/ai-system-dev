# Kali Linux vs ARTEX AI 資安測試功能 gap 報告

**產出日期**：2026-11-16  
**來源**：kali.org/tools/（官方 16 大類）、ARTEX AI v0.9 原始碼（artex-global）  
**狀態**：純分析報告，非網頁輸出，符合 AGENTS.md 第零級規則  
**可接手性**：本檔供任何後續 AI 直接接手實作補齊工作

---

## 1. Kali Linux 資安測試功能全景（16 大類，官方 kali.org/tools/）

| # | 類別 | 代表工具 | 測試目的 |
|---|------|----------|----------|
| 1 | Reconnaissance（偵察） | nmap, amass, theHarvester, spiderfoot, sherlock, dnsrecon, gobuster, ffuf, whatweb | 網路偵察、資產發現、OSINT、Web 掃描 |
| 2 | Resource Development（資源開發） | msfvenom, donut, apktool, jadx, ghidra, rizin, searchsploit, pyinstaller | payload 生成、反編譯、模糊測試、漏洞資料庫 |
| 3 | Initial Access（初始存取） | sqlmap, commix, gophish, setoolkit | SQLi、命令注入、社交工程、釣魚 |
| 4 | Execution（執行） | armitage, beef-xss, xsser, nishang | 遠端程式執行、瀏覽器漏洞利用 |
| 5 | Persistence（持久化） | weevely, webshells, backdoor-factory, cymothoa | Webshell、後門、持久化驻留 |
| 6 | Privilege Escalation（權限提升） | linpeas, winpeas, bloodyad, unix-privesc-check, lynis | 權限提升掃描與利用 |
| 7 | Defense Evasion（防禦規避） | mimikatz, evil-winrm, steghide, veil, donut, ccrypt | 防禦規避、隱寫、免殺、憑證竊取 |
| 8 | Credential Access（憑證存取） | hashcat, john the ripper, hydra, aircrack-ng, wifite, trufflehog, kerberoast | 密碼破解、WiFi 破解、憑證外洩、Kerberoasting |
| 9 | Discovery（探索） | masscan, sslscan, wireshark, bloodhound, fierce, hping3, snmp-check | 服務發現、SSL/TLS 分析、網路封包、Active Directory |
| 10 | Lateral Movement（橫向移動） | impacket-psexec, evil-winrm, rdesktop, smbmap, crackmapexec | 橫向移動、遠端桌面、SMB、WinRM |
| 11 | Collection（收集） | ettercap, mitmproxy, evilginx2, sslsplit, httrack | 中間人攻擊、流量擷取、網站鏡像 |
| 12 | Command and Control（C2） | metasploit, chisel, sshuttle, netcat, socat, havoc, powershell-empire | C2 通道、隧道、代理、反連 |
| 13 | Exfiltration（外洩） | netcat, impacket-smbserver, goshs, raven | 資料外洩測試 |
| 14 | Impact（影響/DoS） | siege, slowhttptest, goldeneye, t50, thc-ssl-dos, dhcpig | DoS/壓力測試、資源耗盡 |
| 15 | Forensics（取證） | autopsy, volatility, binwalk, foremost, scalpel, sleuth kit, yara, dc3dd, guymager | 磁碟取證、記憶體取證、檔案救援、惡意碼分析 |
| 16 | Services & Other（服務與其他） | dradis, faraday, maltego, gvm, juice-shop, dvwa, cherrytree, pipal | 報告平台、弱點環境、案件管理 |

> 官方頁面共收錄 600+ 工具；上表僅列代表工具，完整清單請見 kali.org/tools/。

---

## 2. ARTEX AI v0.9 現有資安測試能力（從原始碼與文件分析）

### 2.1 內建程式（Docker 映像，12 個）
nmap, curl, wget, ripgrep, jq, git, dig, netcat, whois, Node.js 22, Playwright + Chromium, Playwright MCP

### 2.2 核心架構
- **雙圖結構**：資產圖（共用，root_domain/subdomain/ip/service/app/endpoint）+ 探索圖（每任務，goal/intent/fact/finding/hint）
- **多代理自動化**：planner（規劃者）+ workers（工作者，預設 3 個），事件驅動迴圈
- **錄製 MITM 代理**（:8788）：所有工具呼叫流量記錄為證據
- **攔截閘門**：規則 → LLM 法官 → 人工核准，三層閘門
- **發現重跑**：retester 代理重新驗證發現，自動標記 Fixed/Still reproducible
- **資產覆蓋度圖**：力導向圖，顯示 rendered/untested 資產
- **MCP / 技能 / 記憶 / 工作區**：擴充介面，支援遠端工具伺服器與可複用技能
- **單一 Go 二進制**：Next.js UI 內嵌（go:embed），僅需 PostgreSQL 作為外部依賴
- **支援語言**：en / ko / zh / es / ja

### 2.3 現有技能（skill.yaml）支援
- 5 個 API 端點：Load / LoadAll / Validate / Write / GetSchema / ScanDir / ToLegacy / LegacyTo / GetDirInfo
- skill.yaml 欄位：name, display_name, description, version, category, author, license, tags, i18n, compatibility, dependencies, allowed_tools, parameters, entrypoint, outputs, examples, documentation, maintainers, status, test, hooks, permissions
- entrypoint 類型：script / binary / docker / mcp
- 技能分類（原有）：recon / exploit / post-exploit / utility / custom（僅 5 類）
- 前端：skills 頁技能庫檔案樹、編輯、上傳 zip、MCP 連結、可見性控制、呼叫統計、skill.yaml 詳情檢視

---

## 3. Gap 比較：Kali 16 大類 vs ARTEX AI 現有功能

| # | Kali 類別 | ARTEX 對應 | 狀態 | 缺失工具／能力 |
|---|-----------|-----------|------|---------------|
| 1 | Reconnaissance | ⚠️ 部分 | 有 nmap（內建）、asset 發現、flow 捕獲 | 無 amass/theHarvester/gobuster/burpsuite/zaproxy 等偵察工具直接整合為技能 |
| 2 | Resource Development | ❌ 缺口 | 無 | msfvenom, donut, apktool, jadx, ghidra, rizin, searchsploit, pyinstaller |
| 3 | Initial Access | ❌ 缺口 | 無 | sqlmap, commix, gophish, setoolkit |
| 4 | Execution | ✅ 有 | 技能 entrypoint 支援 script/binary/docker/mcp | — |
| 5 | Persistence | ❌ 缺口 | 無 | weevely, webshells, backdoor-factory, cymothoa |
| 6 | Privilege Escalation | ❌ 缺口 | 無 | linpeas, winpeas, bloodyad, unix-privesc-check |
| 7 | Defense Evasion | ❌ 缺口 | 無 | mimikatz, evil-winrm, steghide, veil, donut, ccrypt |
| 8 | Credential Access | ⚠️ 部分 | 有「credential」資產類型 | 無 hashcat/john/hydra/aircrack-ng/wifite/trufflehog 等直接整合 |
| 9 | Discovery | ⚠️ 部分 | 有資產發現與 endpoint 探索 | 無 masscan/wireshark/sslscan/bloodhound 等直接整合 |
| 10 | Lateral Movement | ❌ 缺口 | 無 | impacket-psexec/evil-winrm/rdesktop/smbmap/crackmapexec |
| 11 | Collection | ⚠️ 部分 | 有錄製 MITM 代理與流量證據 | 無 ettercap/sslsplit/sslsniff 等直接整合 |
| 12 | Command and Control | ⚠️ 部分 | 有 MCP（stdio/http/sse） | 無 metasploit/chisel/sshuttle/netcat/socat 等 C2 直接整合 |
| 13 | Exfiltration | ❌ 缺口 | 無 | netcat/impacket-smbserver/goshs/raven |
| 14 | Impact | ❌ 缺口 | 無 | siege/slowhttptest/goldeneye/t50/thc-ssl-dos/dhcpig |
| 15 | Forensics | ⚠️ 部分 | 有證據快照與流量證據 | 無 autopsy/volatility/binwalk/foremost/scalpel/yara/dc3dd/guymager |
| 16 | Services & Other | ⚠️ 部分 | 有 findings/report 輸出 | 無 dradis/faraday/maltego 整合，無 juice-shop/dvwa 弱點環境 |

**統計**：16 類中，完整對應 1 類（Execution）、部分對應 6 類（Recon/Credential/Discovery/Collection/C2/Forensics/Services）、無對應 9 類（Resource Dev/Initial Access/Persistence/Privilege Escalation/Defense Evasion/Lateral Movement/Exfiltration/Impact）。

---

## 4. 已實作補齊（本次工作）

| 項目 | 檔案 | 說明 |
|------|------|------|
| skill.yaml Category 擴充 | `server/skill_yaml.go` | 加入 14 個 Kali 類別（resource-dev/initial-access/persistence/priv-esc/defense-evasion/credential-access/discovery/lateral-movement/collection/c2/exfiltration/impact/forensics/services） |
| Kali 工具後端 | `server/kali_tools.go`（新） | 約 80 工具靜態清單（涵蓋 16 類別）、`GET /api/kali/categories`、`GET /api/kali/tools`、`POST /api/kali/tools/{name}/install`、`GET /api/kali/stats`（per-category call/finding 統計）、skill.yaml 自動產生與驗證、skill_usage ledger 記錄 |
| 路由註冊 | `server/server.go` | 註冊三個 `/api/kali/*` 端點 |
| Kali 工具目錄頁 | `web/src/app/(main)/system/kali/page.tsx`（新）| 16 類別工具展示、一鍵安裝為技能、kali.org 連結 |
| Kali 分類儀表板 | `web/src/app/(main)/system/kali-dashboard/page.tsx`（新）| tool_count/skill_count/call_count/finding_count/coverage% 彙總 |
| 技能頁連結 | `web/src/app/(main)/system/skills/page.tsx` | header 新增「Kali Tools Directory」連結 |
| API 客戶端 | `web/src/lib/api.ts` | 新增 `kaliCategories/kaliTools/kaliToolInstall/kaliStats` |
| 型別定義 | `web/src/lib/types.ts` | 新增 `KaliCategoryInfo/KaliCategoryStats/KaliToolInfo/KaliStatsResponse` |
| Go 單元測試 | `server/kali_tools_test.go`（新）| 4 測試：類別清單、工具清單、分類過濾、安裝產生 skill.yaml |
| Spec Kit 規格層 | `specs/001-kali-gap-remediation/` | spec.md/plan.md/tasks.md/checklist/feature.json，Phase 7 Convergence 已 append 5 個缺口任務（T022-T026） |

---

## 5. 尚未補齊的缺類別（建議另開工作項）

以下 9 個 Kali 類別尚無 ARTEX 對應技能，建議依優先序另開 Spec Kit 規格層工作：

| 優先序 | Kali 類別 | 建議補齊工具（範例） |
|--------|-----------|---------------------|
| P1 | Initial Access | sqlmap, commix, gophish, setoolkit |
| P1 | Credential Access | hashcat, john, hydra, aircrack-ng, wifite, trufflehog |
| P1 | Privilege Escalation | linpeas, winpeas, bloodyad |
| P2 | Resource Development | msfvenom, donut, apktool, jadx, ghidra, rizin |
| P2 | Defense Evasion | mimikatz, evil-winrm, steghide, veil, donut |
| P2 | Lateral Movement | impacket-psexec, evil-winrm, rdesktop, smbmap |
| P3 | Exfiltration | netcat, impacket-smbserver, goshs, raven |
| P3 | Impact | siege, slowhttptest, goldeneye, t50, thc-ssl-dos |
| P3 | Forensics | autopsy, volatility, binwalk, foremost, scalpel, sleuth kit |

---

## 6. 驗證結果（本次工作）

| 驗證項目 | 結果 |
|----------|------|
| `go build ./server/...` | ✅ 通過 |
| `go vet ./server/...` | ✅ 通過 |
| `go test ./server/... -run TestKali` | ✅ 4/4 通過 |
| `npm run build` (Next.js) | ✅ 通過（TypeScript 無錯誤，路由 `/system/kali`、`/system/kali-dashboard` 已編譯） |
| `docker compose -f docker-compose.local.yml config` | ✅ 配置有效 |
| `artex.exe` 建構與啟動 | ✅ 44 MB，成功啟動並監聽 127.0.0.1:8787 |
| `docker compose up -d --build` | ⚠️ 未執行（本機 Windows Docker Desktop daemon 未啟動，需 Linux/Docker 環境） |
| `git commit` / `git push` | ⚠️ 未執行（需 remote URL） |

---

## 7. 下一步建議（供後續 AI 接手）

1. **在 Linux/Docker 環境執行 `docker compose up -d --build`**，確認 PostgreSQL 容器與 ARTEX 容器正常啟動（T023，本次因環境限制未執行）。
2. **執行 `/speckit.converge` 最終收斂驗證**，確認程式碼與 spec/plan/tasks 一致（Phase 7 已 append T022-T026，T022/T024/T025/T026 已完成，T023 需 Docker 環境）。
3. **提供 remote URL** 後執行 `git remote add origin <url>` 並 `git push -u origin master`，將本次變更推送至遠端。
4. **另開 Spec Kit 規格層工作**，針對第 5 節列出的 9 個缺類別（特別是 P1 的 Initial Access、Credential Access、Privilege Escalation）編寫 skill.yaml 模板與後端整合。
5. **考慮在 Docker 映像中內建 nmap（已有）、sqlmap、nikto、wapiti、nuclei 等常用掃描器**，並增加 `/api/kali/tools/{tool}/run` 執行端點，將 Kali 工具輸出寫入 ARTEX 發現（finding）與流量證據。

---

## 8. 相關檔案清單

**新建檔案：**
- `specs/001-kali-gap-remediation/spec.md`
- `specs/001-kali-gap-remediation/plan.md`
- `specs/001-kali-gap-remediation/tasks.md`
- `specs/001-kali-gap-remediation/checklists/requirements.md`
- `specs/.specify/feature.json`
- `server/kali_tools.go`
- `server/kali_tools_test.go`
- `web/src/app/(main)/system/kali/page.tsx`
- `web/src/app/(main)/system/kali-dashboard/page.tsx`
- `docs/KALI_ARTEX_SEC_TESTING_GAP.md`（本檔）
- `docs/KALI_GAP_ANALYSIS.md`（先前建立）
- `docs/KALI_REMEDIATION_PLAN.md`（先前建立，已更新狀態）

**修改檔案：**
- `server/skill_yaml.go`（Category 白名單擴充）
- `server/server.go`（`/api/kali/*` 路由註冊）
- `web/src/app/(main)/system/kali/page.tsx`（JSX 片語修正）
- `web/src/app/(main)/system/kali-dashboard/page.tsx`（stats API 串接）
- `web/src/app/(main)/system/skills/page.tsx`（Kali 連結 + PlayIcon import）
- `web/src/lib/api.ts`（Kali API 方法 + 型別 import）
- `web/src/lib/types.ts`（Kali 型別定義）

---

> 本報告純檔案輸出，未使用網頁/瀏覽器呈現（符合 AGENTS.md 第零級規則）。任何後續 AI 可直接以本檔為基礎接手補齊工作。
