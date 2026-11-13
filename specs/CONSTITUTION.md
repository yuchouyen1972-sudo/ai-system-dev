# 專案憲法 (Constitution)

本專案採用 Spec-Driven Development (SDD) 流程。所有功能開發、錯誤修復與構想評估必須遵守下列原則。

## 核心原則

### 1. 規格優先 (Spec-First)
任何程式碼變更必須先有對應的規格文件（`specs/.specify/memory/` 或 `specs/` 下的 Markdown 規格）。未經規格審核即實作的程式碼不得合併進主分支。

### 2. 可維護性 (Maintainability)
- 所有模組必須有模組級說明文件（docstring / README）。
- 公開 API 必須有類型標註與輸入輸出說明。
- 避免魔法數字；所有常數必須命名並加註來源。

### 3. 測試優先 (Testing-First)
- 新功能必須伴隨對應的單元測試與整合測試。
- 回測/交易邏輯變更必須在歷史數據上驗證，並比較前後 equity curve。
- 測試必須在 CI / 本地環境中可重複執行。

### 4. 安全與隱私 (Security & Privacy)
- `.opencode/` 及所有憑證、金鑰、token 必須加入 `.gitignore`，絕不納入版本控制。
- 模型權重（`.pkl`）、原始數據（`.csv`）等大型二進位檔不納入版本控制，僅保留譜名清單。
- 日誌檔（`.log`）不納入版本控制。

### 5. SDD 流程遵守
每個功能必須依序經過：
1. `/speckit.constitution` — 確認專案原則
2. `/speckit.specify` — 撰寫規格
3. `/speckit.plan` — 建立實作計畫
4. `/speckit.tasks` — 產出任務清單
5. `/speckit.implement` — 執行實作
6. `/speckit.converge` — 收斂驗證

## 適用範圍

本憲法適用於 `specs/` 子目錄下的所有 SDD 產物，以及與之對應的專案根目錄源碼（如 `tmf_engine.py`、`backtest.py`、`train_*.py` 等）。

## 修訂

本憲法可透過 `/speckit.specify` 重新執行來更新，變更必須經過專案維護者審核。

---
*建立於 2026-11-13，以 Spec Kit SDD 流程初始化。*