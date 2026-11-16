# Kali Linux 功能 vs ARTEX AI gap analysis

來源：kali.org/tools/ 官方工具頁（擷取於 2026-11-16），與 ARTEX AI v0.9 原始碼（artex-global）比對。
本報告為純檔案輸出，非網頁/瀏覽器呈現。

---

## 1. Kali Linux 完整功能分類（官方 16 大類）

Kali Linux 以 MITRE ATT&CK 風格組織工具，共 16 大類，下表列出各類別與代表工具：

| # | 類別 | 代表工具（節錄） |
|---|------|----------------|
| 1 | Reconnaissance（偵察） | metagoofil, spiderfoot, emailharvester, sherlock, amass, nmap, theHarvester, dnsrecon, dirb, gobuster, ffuf, burpsuite, nuclei, zaproxy, bettercap, kismet, hackrf_info |
| 2 | Resource Development（資源開發） | code-oss, donut, pyinstaller, olevba, afl-fuzz, gdb, ghidra, radare2, rizin, cutter, apktool, jadx-gui, msfvenom, searchsploit |
| 3 | Initial Access（初始存取） | gophish, setoolkit, metasploit-framework, sqlmap, commix |
| 4 | Execution（執行） | armitage, metasploit-framework, beef-xss, xsser, nishang, powersploit |
| 5 | Persistence（永續化） | laudanum, phpggc, webacoo, webshells, weevely, backdoor-factory |
| 6 | Privilege Escalation（權限提升） | lynis, peass, linpeas, winpeas, unix-privesc-check, bloodyad |
| 7 | Defense Evasion（防禦規避） | mimikatz, evil-winrm, impacket-scripts, netexec, steghide, outguess, donut, shellter, veil, ccrypt |
| 8 | Credential Access（憑證存取） | mimikatz, rubeus, samdump2, hashid, cewl, crunch, hydra, john, hashcat, aircrack-ng, wifite, kerberoast, trufflehog, gitxray |
| 9 | Discovery（探索） | amass, masscan, nmap, sslscan, sslyze, wireshark, scapy, tcpdump, arping, fierce, hping3, snmp-check, bloodhound, sharphound |
| 10 | Lateral Movement（橫向移動） | crackmapexec, evil-winrm, impacket-scripts, impacket-psexec, rdesktop, mimikatz, netexec, smbmap |
| 11 | Collection（收集） | ettercap, evilginx2, mitmproxy, mitm6, ssldump, sslsplit, wifipumpkin3, httrack |
| 12 | Command and Control（C2） | metasploit-framework, powershell-empire, chisel, ligolo-ng, sshuttle, netcat, socat, havoc, hoaxshell, villain, adaptix |
| 13 | Exfiltration（外洩） | netcat, impacket-smbserver, goshs, raven |
| 14 | Impact（影響） | dhcpig, goldeneye, mdk3, siege, slowhttptest, t50, thc-ssl-dos, scapy |
| 15 | Forensics（鑑識） | autopsy, binwalk, foremost, photorec, scalpel, testdisk, dc3dd, dcfldd, guymager, pdfid, pdf-parser, yara, volatility |
| 16 | Services and Other Tools（報告/實驗室/系統服務） | cherrytree, dradis, faraday, maltego, obsidian, pipal, juice-shop, dvwa, defectdojo, gvm, thehive, kali-tweaks, shell-gpt, tailscale |

> 註：官方頁面共收錄數百個工具；上表僅列代表工具，完整清單請見 kali.org/tools/。

---

## 2. ARTEX AI v0.9 現有功能（從原始碼分析）

ARTEX 是「AI 自動滲透測試協作平台」，定位為平台/協作層，透過「技能（skill）」調用外部工具，而非直接內建全部 Kali 工具。

### 2.1 技能系統（skill.yaml）
- `server/skill_yaml.go`：SkillYAML 結構與 5 個 API 端點（Load / LoadAll / Validate / Write / GetSchema / ScanDir / ToLegacy / LegacyTo / GetDirInfo）
- skill.yaml 欄位：name, display_name, description, version, category, author, license, tags, i18n, compatibility, dependencies, allowed_tools, parameters, entrypoint, outputs, examples, documentation, maintainers, status, test, hooks, permissions
- entrypoint 類型：script / binary / docker / mcp
- output 類型：json / text / file / asset / finding
- 技能分類（category）：**recon / exploit / post-exploit / utility / custom**（僅 5 類）
- 前端：`web/src/app/(main)/system/skills/page.tsx` + `SkillYAMLDisplay.tsx`，支援技能庫檔案樹、編輯、上傳 zip、MCP 連結、可見性控制、呼叫統計、skill.yaml 詳情檢視

### 2.2 任務（Task）
- Task 狀態：created / queued / running / paused / done / failed /timeout
- Task Category、Task Template、Task Archive（封存/還原）
- 探索圖：TaskNode（task/begin/goal/intent/fact/finding/hint/digest）、Edge（spawns/derived_from/yields/proves/covers）

### 2.3 資產（Asset）
- AssetNode：company / domain / ip / port / service / site / endpoint / parameter / tech / credential / data
- 新統一資產模型：root_domain / ip / subdomain / app / service / endpoint
- Company、ScopeRule（domain / ip / cidr / icp / keyword）、覆蓋度圖（Coverage Graph）

### 2.4 發現（Finding）
- severity：critical / high / medium / low
- status：pending / in_progress / confirmed / resolved / fixed / false_positive / ignored / duplicate / risk_accepted
- FindingAssetTree（按資產組織）、FindingStats、FindingTrafficBinding（流量證據綁定）

### 2.5 對話/活動（Conversation / Activity）
- Conversation（主 agent / planner / worker）、Session
- Activity：tool_use / tool_result / text / thinking / result / user / intent / round / usage / llm_switch / llm_failover / intercept_request
- LLM 記錄（LLMRecordItem / LLMRecordDetail）

### 2.6 LLM 與 Agent 管理
- LLMProfile（format / model / base_url / proxy / thinking_type / reasoning_effort / retry）
- LLMPool（故障轉移、breaker、intent 重跑）、LLM usage metering
- Agent（builtin / custom）、AgentTrigger（定時/發現/目標/超時/工具/任務觸發）

### 2.7 MCP 與工具
- MCPServer（stdio / http / sse）、MCPTool、技能與 MCP 連結
- Tool（builtin / shell / command / script / http）、InterceptRule（allow / deny / ask）
- AssetInterceptRule（block / allow）、JudgeConfig（模型兜底審批）

### 2.8 流量（Traffic）
- TrafficExchange（請求/回應記錄）、TrafficHost、TrafficDetail
- 流量證據快照（TrafficEvidenceSnapshot）

### 2.9 通知（Notify）
- NotificationChannel（kind / realtime-digest / filter）、NotificationDelivery
- 漏洞 IM 推送（notify_enabled / digest_interval_min / public_base_url）

### 2.10 設定（Settings）
- traffic_capture、web_search（ddgs/brave-free/tavily/deepseek）、global_proxy、python_interpreter、workers、task_concurrency
- notify、output_language（en/ko/zh/es/ja）

### 2.11 其他
- 更新檢查（GET /api/update/check，GitHub 版本比較）、日誌（/api/logs）、認證（/api/auth/*）

---

## 3. Gap 比較：Kali 16 大類 vs ARTEX 現有功能

| # | Kali 類別 | ARTEX 對應 | 狀態 | 缺口說明 |
|---|-----------|-----------|------|----------|
| 1 | Reconnaissance | 部分 | ⚠️ 部分 | 有 asset 發現、flow 捕獲；無內建 nmap/amass/theHarvester/gobuster/burpsuite/zaproxy 等偵察工具包裹 |
| 2 | Resource Development | 無 | ❌ 缺口 | 無 payload 生成（msfvenom/donut）、反編譯（jadx/apktool）、模糊測試（afl-fuzz）等技能 |
| 3 | Initial Access | 無 / exploit 部分 | ❌ 缺口 | 無 SQLi（sqlmap）、社交工程（gophish/setoolkit）、命令注入（commix）等初始存取技能 |
| 4 | Execution | 有 | ✅ | 技能 entrypoint 支援 script/binary/docker/mcp，可執行外部工具 |
| 5 | Persistence | 無 | ❌ 缺口 | 無 Webshell/後門/持久化技能（weevely/backdoor-factory 等） |
| 6 | Privilege Escalation | 無 | ❌ 缺口 | 無權限提升掃描（linpeas/winpeas/bloodyad）技能 |
| 7 | Defense Evasion | 無 | ❌ 缺口 | 無防禦規避技能（steganography/mimikatz/shellter/veil 等） |
| 8 | Credential Access | 部分 | ⚠️ 部分 | 有「credential」資產類型；無密碼破解（john/hashcat/hydra）、WiFi 破解（aircrack-ng/wifite）技能 |
| 9 | Discovery | 部分 | ⚠️ 部分 | 有資產發現與 endpoint 探索；無 nmap/masscan/wireshark/sslscan/bloodhound 等探索工具直接整合 |
| 10 | Lateral Movement | 無 | ❌ 缺口 | 無橫向移動技能（impacket-psexec/evil-winrm/rdesktop/smbmap 等） |
| 11 | Collection | 部分 | ⚠️ 部分 | 有流量捕獲與證據綁定；無 MITM（ettercap/evilginx2/mitmproxy）與流量擷取技能 |
| 12 | Command and Control | 部分 | ⚠️ 部分 | 有 MCP（stdio/http/sse）；無內建 C2（metasploit/powershell-empire/chisel/sshuttle 等） |
| 13 | Exfiltration | 無 | ❌ 缺口 | 無外洩測試技能（netcat/impacket-smbserver/goshs/raven 等） |
| 14 | Impact | 無 | ❌ 缺口 | 無 DoS/壓力測試技能（siege/slowhttptest/thc-ssl-dos/goldeneye 等） |
| 15 | Forensics | 部分 | ⚠️ 部分 | 有證據快照與流量證據；無鑑識工具（autopsy/binwalk/foremost/scalpel/yara）直接整合 |
| 16 | Services and Other Tools | 部分 | ⚠️ 部分 | 有 findings/report 輸出；無 Dradis/Faraday/Maltego 等報告平台與弱點環境（juice-shop/dvwa）整合 |

**統計**：16 類中，完整對應 1 類（Execution）、部分對應 6 類（Recon/Credential/Discovery/Collection/C2/Forensics/Services）、無對應 9 類（Resource Dev/Initial Access/Persistence/Privilege Escalation/Defense Evasion/Lateral Movement/Exfiltration/Impact）。

---

## 4. 核心缺口總結

1. **技能分類體系過於精簡**：ARTEX 僅有 recon/exploit/post-exploit/utility/custom 5 類，無法涵蓋 Kali 的 16 大功能領域（尤其缺 resource-development, initial-access, persistence, privilege-escalation, defense-evasion, lateral-movement, exfiltration, impact）。
2. **缺乏內建 Kali 工具包裹技能**：未提供 nmap、metasploit、burpsuite、sqlmap、john、hashcat、wireshark、binwalk、autopsy 等主流 Kali 工具的 skill.yaml 模板與入口點。
3. **前端無針對各 Kali 類別的儀表板**：skills 頁僅為通用技能庫，沒有依照 Kali 功能領域（偵察/弱點/Web/密碼/無線/鑑識等）的專門工作流與卡片。
4. **缺乏 Kali 工具安裝/管理流程**：無從 Kali 來源自動安裝、版本比對、工具可用性檢查的平台功能。

---

## 5. 補齊建議（待使用者確認後實作）

### 5.1 後端（Go）
- 擴充 `skill_yaml.go` 的 `Category` 驗證白名單，加入 Kali 16 大類（recon, resource-dev, initial-access, execution, persistence, priv-esc, defense-evasion, credential-access, discovery, lateral-movement, collection, c2, exfiltration, impact, forensics, services）。
- 新增 `server/kali_tools.go`：Kali 工具目錄資料來源（靜態清單，可從官方同步）、工具安裝/版本檢查、技能模板產生（依工具類別產生 skill.yaml）。
- 新增 API：`GET /api/kali/tools`（列舉 Kali 工具）、`POST /api/kali/tools/{tool}/install`（安裝為技能）、`GET /api/kali/categories`（類別總覽）。

### 5.2 前端（Next.js）
- 在 `skills` 頁新增「Kali 工具目錄」分頁，依 16 大類分組展示工具，支援一鍵安裝為 ARTEX 技能。
- 新增 `KaliCategoryDashboard` 依類別彙總呼叫統計與發現數。
- `SkillYAMLDisplay` 新增 Kali 工具來源標籤（kali.org 對應連結）。

### 5.3 skill.yaml 模板（範例）
為主流 Kali 工具產生 skill.yaml 模板，例如：nmap-recon、metasploit-exploit、burpsuite-web-scan、sqlmap-sqli、john-password-crack、wireshark-sniff、autopsy-forensics 等。

---

## 6. 下一步（需使用者確認）

上述 5.x 屬於「新功能／新模組實作」，根據 FACTORY.md 路由規則，需先經過 **Spec Kit 規格層**（`/speckit.specify` → `/speckit.plan` → `/speckit.tasks`），再進 **ai-dev-team pipeline**（Planner → Coder → Reviewer）實作，最後驗證（測試／lint）才算交付。

請確認：
1. 是否接受以上 gap analysis 結果？
2. 是否授權進入 Spec Kit 規格層，開始撰寫「Kali 功能補齊」規格與計畫？

報告已儲存：`D:\101.中信銀備份資料\T108.AI模擬網路防衛\artex-global\docs\KALI_GAP_ANALYSIS.md`
