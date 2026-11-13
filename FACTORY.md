# AI 工廠預設工作模式（Always-On）

> **全體 AI 先讀 `SOP.md`（同目錄）**：角色分工、協作主流程、工作計畫與時程（S1-S11）、可互查驗證矩陣、成本規則。本檔只管「任務走哪條路」。

本工作區的預設工作方式：凡涉及開發任務，優先使用本區的兩套自有工具，而非全部由主代理直接手寫。此規則對本工作區內所有 OpenCode session／agent 生效。

## 路由規則

### 走 Spec Kit 規格層（前置，新增）

以下類型的任務，開工時先執行 Spec Kit 規格層，產出規格後才進入開發管線：

- 新功能／新模組的實作 → `/speckit.specify` → `/speckit.plan` → `/speckit.tasks`
- 錯誤修復（Bug） → `/speckit-bug-assess` → `/speckit-bug-fix` → `/speckit-bug-test`
- 構想評估 → `/speckit-assess-intake` → … → `/speckit-assess-decide`

規格層產出存於 `specs/.specify/memory/`，規格確認（converge / verified / go）後，再進入 ai-dev-team / ai-software-factory 管線實作。

例外：單行級修改、錯字、格式化、跑測試／lint、git 操作、檔案搬運、計畫檔／文件／設定檔維護、回答問題、研究、解釋類，可直接走現行例外清單，不需經過 Spec Kit 規格層。

### 走 ai-dev-team（任務級，預設主力）

以下類型的任務，若已通過 Spec Kit 規格層，則直接從草稿整合開始；否則先經過規格層。開工時先用 ai-dev-team 取得草稿，再由主代理整合、驗證、落盤：

- 新功能／新模組的實作
- 重構
- 程式碼審查
- 除錯（pipeline 含 debugger 角色）

指令：

```bash
ai-dev-team pipeline "<任務描述>"      # Planner → Coder → Reviewer 協作管線，回傳草稿
ai-dev-team quick "<prompt>"           # 單次智慧選型（快）
ai-dev-team explain "<prompt>"         # 解釋模型選型理由
```

整合流程：pipeline 產出 → 主代理比對專案既有慣例 → 落盤 → 跑測試／lint 驗證 → 才算完成。管線草稿一律視為未經驗證素材。

### 走 ai-software-factory（整體專案級）

使用者要求「跑工廠／整體生成一個專案／自動化全流程」時：

```bash
ai-software-factory "<需求描述>" -n <專案名> -d <交期小時>   # 全流程：規劃→實作→審查→除錯→測試
ai-software-factory --resume                                 # 斷點續傳
ai-software-factory --kpi-report                             # 事後 KPI 稽核
```

### 例外（直接做，不進管線）

- 單行級修改、錯字、格式化
- 跑測試／lint、git 操作、檔案搬運
- 計畫檔／文件／設定檔維護
- 回答問題、研究、解釋類（無程式碼可生成，管線無意義）

例外的判準：進管線的成本（數分鐘延遲＋API 開銷）明顯大於任務本身時，直接做。拿不準就進管線。

### 前置檢查

- `OPENROUTER_API_KEY` 未設定時管線照跑：兩套工具會自動停用遠端 OpenRouter 模型、改走本機 Ollama（`OLLAMA_BASE_URL`，預設 `http://127.0.0.1:11434/v1`）並印出警告；此情況下直接執行，回報時說明本次全用本機模型即可，不必等待指示。key 絕不寫入任何檔案。
- 管線產出必須通過測試／lint 驗證才算交付。

## 成本意識

每次 pipeline 都是真實 API 開銷與數分鐘延遲。預設是「開發類任務先走工廠」，但琐碎操作優先套用例外清單，避免為改一行字跑整條管線。
