# Tasks: TMF 參數配置驗證

**Input**: Design documents from `specs/.specify/memory/`

**Prerequisites**: plan.md (required), spec.md (required)

**Tests**: 單元測試 (四件套驗證)

**Organization**: 單一功能，不分多個 user story

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可平行執行 (不同檔案，無相依)
- **[Story]**: 對應 user story (US1, US2)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 確認現有環境可執行 tmf_engine.py

- [ ] T001 執行 `python tmf_engine.py` 確認現有主程式正常運作

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 新增 validate_config 函數與 TMFDecisionEngine 整合

**⚠️ CRITICAL**: 無 user story 可開始，此階段為唯一實作

- [ ] T002 [US1] 在 tmf_engine.py 新增 `validate_config(spec, account) -> dict` 函數，實作 FR-001 至 FR-006 驗證項目
- [ ] T003 [US1] 在 TMFDecisionEngine.__init__ 呼叫 validate_config，附於 `self.config_validation` (FR-002)
- [ ] T004 [US1] 新增 `tests/test_config_validation.py` 單元測試，涵蓋正常與異常參數案例 (FR-003 至 FR-007)

---

## Phase 3: User Story 1 - 啟動時自動驗證參數 (P1) 🎯 MVP

**Goal**: TMFDecisionEngine 初始化時自動驗證參數，結果可透過 config_validation 取得

**Independent Test**: 建立不同參數的 TMFSpec/Account，初始化 TMFDecisionEngine，檢查 config_validation.valid 與 errors

### Tests for User Story 1

- [ ] T005 [US1] 測試正常參數：valid=True，errors 為空
- [ ] T006 [US1] 測試保證金不足：valid=False，errors 包含保證金相關
- [ ] T007 [US1] 測試 tick_value=0：valid=False
- [ ] T008 [US1] 測試 min_stop_points > max_stop_points：valid=False
- [ ] T009 [US1] 測試 max_risk_pct > 0.05：valid=True 但 warnings 非空

---

## Phase 4: User Story 2 - 按需回傳結構化驗證報告 (P2)

**Goal**: validate_config 回傳結構化 dict，供外部檢視

**Independent Test**: 直接呼叫 validate_config()，斷言回傳結構

### Tests for User Story 2

- [ ] T010 [US2] 測試回傳 dict 包含 valid/warnings/errors/checks 四個欄位
- [ ] T011 [US2] 測試每個 check item 包含 name/passed/detail 三個欄位

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 四件套驗證與 collect

- [ ] T012 [P] 執行 `python -m pytest tests/test_config_validation.py -q` 確認測試通過
- [ ] T013 [P] 執行 `ruff check tmf_engine.py tests/test_config_validation.py`
- [ ] T014 [P] 執行 `ruff format --check tmf_engine.py tests/test_config_validation.py`
- [ ] T015 [P] 執行 `mypy tmf_engine.py`
- [ ] T016 [US1] `/speckit.converge` 收斂驗證，確認 Converged

## Dependencies & Execution Order

- **Setup (T001)**: 無相依，可立即執行
- **Foundational (T002-T004)**: 依賴 Setup，區塊所有後續
- **User Story 1 (T005-T009)**: 依賴 Foundational，可獨立測試
- **User Story 2 (T010-T011)**: 依賴 Foundational，可獨立測試
- **Polish (T012-T016)**: 依賴所有上述完成

## Parallel Opportunities

- T005-T011 各測試案例可平行撰寫 (不同測試函數)
- T012-T015 四件套可平行執行 (不同命令)

## Notes

- [P] 任務可平行
- 每個測試先寫、先失敗，再實作讓其通過 (TDD)
- Commit 後再執行 converge