# AI 軟體工廠工作 SOP（全體 AI 必讀）

> 目的：讓本工作區所有 AI（主代理／各 session、ai-dev-team 管線角色、ai-software-factory 階段代理人）以**最高效率、最低成本、最正確、可互查**的方式分工協作。
> 文件層級：`AGENTS.md`（人類原則，最高）→ `FACTORY.md`（路由規則）→ **本檔（執行細節、計畫、時程、互查）** → 各倉 tests/README（技術細節）。
> 建立：2026-10-07。每次重大狀態變更後，主代理須同步更新本檔與 `.omo/plans/ai-factory-continuation-plan.md`。

---

## 1. 角色分工表

| 角色 | 職責 | 輸入 | 輸出 | 驗證責任 |
|---|---|---|---|---|
| 主代理（當前 session 的 AI） | 任務接收、路由判斷、整合落盤、驗證、狀態同步、回報 | 使用者指令、本檔、FACTORY.md | 程式碼/文件、commit（push 需逐次授權）、回報 | 跑測試/lint、查 CI、同步計畫檔 |
| ai-dev-team pipeline | 任務級草稿：Planner→Coder→Reviewer（失敗含 debugger） | 任務描述字串 | `pipeline_runs/` 下各 stage 檔＋評分（SCORE/VERDICT） | 草稿＝未經驗證素材，**主代理測試通過才算數** |
| ai-software-factory | 專案級全流程：規劃→實作→審查→除錯→測試＋KPI/deadline/仲裁 | `"<需求>" -n <名> -d <時數>` | 專案目錄＋checkpoint | 同上＋`--kpi-report` 稽核 |
| 本機 Ollama（127.0.0.1:11434） | 無 key 時唯一算力；有 key 時仍可被選為備援 | HTTP（OpenAI 相容 `/v1/chat/completions`） | 模型回應 | doctor 量測（S1） |
| GitHub Actions CI | push 後最終防線 | commit push | job 結果（綠/紅） | 主代理查證（截圖/DOM，run 頁「地位」） |

**鐵律**：
- 管線產出一律視為未經驗證素材 → 必須過測試／lint 才落盤交付。
- `push` 逐次取得使用者授權；commit 訊息英文 SEMANTIC 單行。
- 三次連續失敗 → 停手、回報、換策略（不硬改）。

---

## 2. 任務協作主流程（SOP）

1. **路由**：依 `FACTORY.md` 判斷——新功能／架構變更／錯誤修復／構想評估先經 Spec Kit 規格層（`specs/`），再進管線 vs 例外直做（單行修改、測試、git、文件、問答＝直做）。
2. **執行**：
   - 進管線 → 下指令 → 收草稿 → 主代理比對專案既有慣例 → 落盤。
   - 無 `OPENROUTER_API_KEY` → 兩套工具自動停用遠端、走本機（黃字警告屬正常，不必等待指示）。
3. **驗證清單**（兩倉都要）：
   - Spec Kit 規格層：`/speckit.converge` 回報 `Converged`（SDD）／`/speckit-bug-test` verdict=`verified`（Bug）／`/speckit-assess-decide` 決策=`go`（Assessment）
   - devteam：`python -m pytest tests/ -q`、`ruff check`、`ruff format --check`、`mypy src`
   - factory：`python -m pytest tests/ -q`、`ruff check`（**基線 224，不得新增**）
4. **提交**：`set GIT_MASTER=1&& git add ... && git commit -m "..."` → push（**先取得授權**）→ `git ls-remote origin main` 對 head。
5. **CI 查證**：BrowserSkill（daemon + Chrome 擴充）看 run 頁，「地位：成功」＋job 數對齊才結案。
6. **狀態同步**：更新 `.omo/plans/ai-factory-continuation-plan.md` 與本檔相關章節。
7. **回報**：繁體中文、無 emoji、附證據表（測試數字、commit、CI）。

---

## 3. 無 key／本機模式運作規則（含已知基線）

- 無 key 時：`RouterConfig(disabled_models=全部遠端)`；caller 按 `provider` 分流（`ollama` → `OLLAMA_BASE_URL`，預設 `http://127.0.0.1:11434/v1`，dummy bearer；遠端 → 呼叫時 raise RuntimeError 由 chain 接住）。
- 逾時：**本機分支 600s**（冷載入+慢生成緩衝）；遠端 devteam 180s、factory `config.timeout` 300s（env `TIMEOUT`）。
- fallback 鏈：16 條遠端鏈尾端全接本地；三大本地（deepseek-r1:32b／qwen3.8:27b／laguna-xs-2.1）→ `qwen2.5-coder:7b`；健康路由會在**失敗一次後**事前跳過（S1 目標：變成不燒那一次）。
- **本機硬體與速度基線（2026-10-07 實測，調優前）**：
  - GPU：RTX 3060 Ti 8GB（driver 617.14，可用約 6GB）；`OLLAMA_*` 調優變數原本**皆未設定**。
  - `qwen2.5-coder:7b`：秒級回應 → **本機唯一生產級主力**。
  - `qwen3.8:27b`：600 token＋冷載入 = 383s（finish=length）→ 全量輸出必爆 600s 預算。
  - `deepseek-r1:32b`：暖機 trivial 11-19s；600 token >390s 超時（r1 思考＋~1.5 tok/s）→ **本機不可用於實戰輸出**。
  - 無 key 完整 pipeline 實證：891.7s、三階段全本機、SCORE 9 PASS。
- 無 key 避開機制：`OLLAMA_DISABLE_MODELS`（逗號分隔 id，併入兩側 disable 集，預設空＝不變）＋ `ai-dev-team doctor` 建議清單（S1）。
- **✅ 事故注記（2026-10-07，已復原）**：被中斷的調參指令遺留了 3 個 HKCU 變數（`OLLAMA_KEEP_ALIVE=30m`、`OLLAMA_KV_CACHE_TYPE=q8_0`、`OLLAMA_FLASH_ATTENTION=1`）並使 Ollama 於 12:18 重啟；新 process 讀到 `OLLAMA_MODELS=D:\101.中信儲備份資料\T102.AI\ollama-models`（**當時為空庫，blobs=0**），而 8 個真模型共 **137.6GB 在預設 `C:\Users\yp86n\.ollama\models`** → `ollama list` 為空、管線本機呼叫全失敗；doctor 可正確偵測（4/4 MISSING、exit 1）。**復原＝方案 b：robocopy 137.6GB 至 D:（44 blobs／8 manifests 核對一致）＋重啟** → `ollama list` 8 模型全回、serve log 確認三調參變數全生效、doctor 重跑 4/4 present（7b 76.7 tok/s、laguna 34.1 OK；qwen3.8/deepseek 2-3 tok/s SLOW → 建議列入 `OLLAMA_DISABLE_MODELS`）。

---

## 4. 工作計畫（改善清單 S1-S11）

| ID | 項目 | 狀態 | 所屬階段 |
|---|---|---|---|
| S1 | 速度感知預檢：`ai-dev-team doctor`（Ollama 可達、模型對表、tok/s 量測、慢模型建議）＋ `OLLAMA_DISABLE_MODELS` 事前避開 | **✅ 完成**（doctor Phase B；disable 接線 Phase C） | Phase B+C |
| S2 | 失敗顯示修復：失敗不把 prompt 當成果印、duration 實記、cp950 `?` 降級 | **✅ 完成**（Phase D） | Phase D |
| S3 | Ollama 調參落地（KV q8_0、flash attention、keep-alive 30m）＋前後 tok/s 對比 | 偵察完成、待放行 | Phase E |
| S4 | 鏈上每次嘗試（model/耗時/錯誤）寫入 stage 檔 | 待做 | 下一批 |
| S5 | 工廠全流程無 key 真跑一次（KPI/deadline/仲裁實證） | 待做 | 下一批 |
| S6 | `@pytest.mark.e2e` 真 HTTP 測試（無 Ollama 即 skip） | 待做 | 下一批 |
| S7 | 逾時分層（7b 120s／27b 600s） | 待做 | 下一批 |
| S8 | 排程自動跑（定時 discover／KPI 稽核） | 需求待使用者定義 | 待定 |
| S9 | openrouter_proxy 純 OpenRouter、list-models 標題未標本機 | 待做 | 低優先 |
| S10 | discover.yml 空清單警告 | 等下週排程自然驗證 | 被動 |
| S11 | FACTORY.md 常駐規則對新 session 生效 | 待使用者重啟 opencode | 使用者動作 |

---

## 5. 工作時程

| 階段 | 內容 | 預估 | 目標日 | 狀態 |
|---|---|---|---|---|
| Phase A | 本 SOP 建立＋FACTORY.md 指引 | 8-10 分 | 2026-10-07 | **✅ 完成** |
| Phase B | S1：doctor 命令＋測試 | 15-20 分 | 2026-10-07 | **✅ 完成**（142 passed+2 skipped；ruff/format/mypy 綠；實跑 exit 1 正常；計量修正 cbcb539 後 151+2skip） |
| Phase C | S1：`OLLAMA_DISABLE_MODELS` 兩側接線＋測試 | 10-15 分 | 2026-10-07 | **✅ 完成**（devteam `1d905f8`＋factory `f7a1b60`；148 passed+2skip／48 passed；已 push，CI 綠） |
| Phase D | S2：失敗顯示修復＋測試 | 15-25 分 | 2026-10-07 | **✅ 完成**（devteam `dc27bd4`＋factory `ec71260`；151 passed+2skip／49 passed、ruff 224=基線；已 push，CI #33/#7 成功） |
| Phase E | S3：Ollama 調參（重啟約 5-10 秒斷線）＋前後量測 | 8-12 分 | 放行後即日 | 待放行（**三變數已入 HKCU 且復原後確認生效；E 僅剩正式前後對比量測，見 §3 事故注記**） |
| Phase F | 全量驗收（兩倉四件套）＋計畫檔同步 | 3-5 分 | 隨 E 後 | 待放行 |
| 下一批 | S4-S7 | 視範圍 | 另行排程 | — |

合計 B-F 約 50-75 分鐘。**push 一律另請授權。**

---

## 6. 可互查驗證矩陣

| 產出物 | 誰產出 | 交叉驗證方式 | 證據位置 |
|---|---|---|---|
| 規格文件（SDD／Bug／Assessment） | Spec Kit 規格層 | `/speckit.converge`＝`Converged` 或 verdict／決策欄位確認 | `specs/.specify/memory/`、`specs/.specify/bugs/`、`specs/.specify/assessments/` |
| 程式碼 | 主代理／管線 | 四件套（pytest/ruff/format/mypy）＋ CI | 倉內 `tests/`、GitHub Actions run |
| 管線草稿 | devteam pipeline | 主代理比對慣例＋測試後才落盤 | `pipeline_runs/<時間戳>/` |
| 選型與路由 | router | `test_router`／`test_local_*` 單測＋doctor 實測 | `tests/`、終端輸出 |
| 逾時與 fallback 行為 | caller/workflow | 單測鎖契約（timeout=600/180/300、鏈到 7b） | `test_local_models.py`、`test_local_routing.py` |
| 硬體速度主張 | 主代理 | 前後對比量化（同 prompt、同參數） | 本檔 §3 基線表 |
| 專案狀態真相 | 主代理 | **三向一致**：計畫檔 ↔ `git log` ↔ CI | `.omo/plans/ai-factory-continuation-plan.md`、git、Actions |
| 規格層狀態 | Spec Kit 規格層 | `specs/.specify/memory/` 與 `.omo/plans/` 互不覆蓋，進度一致 | `specs/.specify/memory/`、`.omo/plans/` |
| 上線 | 主代理 | 使用者逐次授權＋`ls-remote` head 對比 | 終端紀錄 |
| 模型註冊表 vs 實機 | 主代理 | `/api/tags`、`/api/show` ctx 對 `ModelSpec` | doctor（S1） |

---

## 7. 成本規則

- **時間成本（本機實測序）**：7b（秒級）≪ 27b（分鐘級）< 33b ≈ 32b-r1（十分鐘級／不可用）。無 key 模式預設資源＝時間。
- **API 成本（有 key）**：每次 pipeline 為真實 API 開銷＋數分鐘延遲 → 瑣碎任務走 FACTORY.md 例外清單直做。
- **診斷成本**：先 doctor／量測再動手；禁止無數據猜測效能根因（本輪教訓：curl 三層才定位到 32b 速度問題）。
- 憑證：key 只留 secrets 與本機，絕不入任何檔案。

---

## 8. 快速指令集

```bash
# 任務級（Planner→Coder→Reviewer，回傳草稿）
ai-dev-team pipeline "<任務描述>"
ai-dev-team quick "<prompt>"          # 純選型，不下呼叫
ai-dev-team doctor                    # S1 完成後可用：預檢＋速度量測

# 專案級全流程
ai-software-factory "<需求>" -n <專案名> -d <交期小時>
ai-software-factory --resume
ai-software-factory --kpi-report

# 驗收
python -m pytest tests/ -q            # 於各倉根目錄
python -m ruff check . && python -m ruff format --check .   # devteam
python -m mypy src                    # devteam
```

> 對應 README／細節：`ai_dev_team/README.md`、`ai_software_factory/README.md`；路由全文見 `FACTORY.md`。

---

## 9. 工作指導原則（所有 AI 開工前必讀，使用者 2026-10-14 頒布）

### 9.1 開工前計畫書（MD 先行）
- 2 步以上任務，開工前先寫計畫書 `.omo/plans/<任務名>.md`，寫到「換一個 AI 來也完全看得懂」才算數：
  - 目標（一句話驗收標準）、範圍（做什麼／不做什麼）、步驟分解（每步原子、可獨立驗證）
  - 每步預估時程（分鐘）＋全程總預估
  - 稽核計畫（§9.3）
- 計畫書存檔後才開工；重大任務照既有規則先過 Momus 審。
- 每份計畫書必附兩圖（mermaid，與計畫書同目錄）：甘特圖（預估 vs 實際）＋Optima 生產路徑圖
  （各 AI／角色分工、關鍵路徑、哪段可並行）；完工後補實際數與誤差表，作為下次估時輸入。

### 9.2 計時與超時鐵律
- 開工＝計時開始（記錄 wall-clock 起點）。
- 任一階段實際耗時超過該階段預估 **30 分鐘以上** → 立即暫停手邊工作，找原因（卡在哪、哪個假設錯了），回報使用者；不空轉、不硬撐、不悶頭換招。
- 既有「三次連續失敗→停手回報換策略」維持。

### 9.3 稽核計畫（換手可驗）
- 每個產出物都要有可重跑的驗證方式：**指令＋預期結果＋證據位置**，寫進計畫書。
- 程式碼：pytest／ruff／format／mypy 數字；文件：回讀確認；狀態主張：三向一致（計畫檔↔git↔CI／終端）。

### 9.4 分階段執行＋逐段回報
- 照計畫書階段順序，一次一段；每段完成當下立刻回報（做了什麼、證據、實際耗時 vs 預估），再進下一段。
- 範圍變更 → 先改計畫書再動手。

### 9.5 背景任務殭屍防制（保護 OpenCode API）
- 發射即登記：每個 background task（`bg_…`）記用途＋所屬階段。
- 用完即個別 cancel；禁用 `background_cancel(all=true)`。
- 同一目的不重複發射；不等 poll，等 system-reminder；段收尾前清點未關 task。
- 超時判死：單一任務超過該段預估＋30 分無進展 → 視為殭屍 → 個別 cancel＋記錄原因＋換策略。
- 失敗即停：殭屍任務禁止重試風暴燒 API 額度。
- OS 層：臨時腳本跑完即刪；非我方行程只列不殺。
