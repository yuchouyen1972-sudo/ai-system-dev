# Tasks: tmf_1m_addons 偵測器整合進 tmf_1m_engine 即時 TICK 引擎

**Input**: Design documents from `specs/003-addon-detectors-integration/`

**Prerequisites**: plan.md (required), spec.md (required)

**Tests**: 單元測試 + 全量回測 (TMF_1m.parquet)

**Organization**: 單一功能，分四個階段

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可平行執行 (不同檔案，無相依)
- **[Story]**: 對應 user story (US1, US2)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 確認現有環境可執行 tmf_1m_engine.py 與 tmf_1m_addons.py

- [ ] T001 執行 `python tmf_1m_engine.py` 確認現有 TICK 引擎正常運作（回測模式）
- [ ] T002 執行 `python -m py_compile tmf_1m_addons.py` 確認偵測器模組語法正確

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 將 tmf_1m_addons.py 的偵測器匯入 tmf_1m_engine.py，新增 OR 條件進場邏輯

**⚠️ CRITICAL**: 無 user story 可開始，此階段為唯一實作

- [ ] T003 [US1] 在 tmf_1m_engine.py 新增 `from scripts.tmf_1m_addons import ...` 匯入 12 個偵測器類別與 check_resonance
- [ ] T004 [US1] 在 `generate_trading_signal` 新增 BreakoutPullbackDetector 偵測，作為 OR 條件 (reason_code=96)
- [ ] T005 [US1] 在 `generate_trading_signal` 新增 FibonacciStrategyDetector 偵測，作為 OR 條件 (reason_code=97)
- [ ] T006 [US1] 在 `generate_trading_signal` 進場訊號確認後插入 check_resonance 過濾層 (score<3 拒絕，3-5 LOW 信心)

## Phase 3: User Story 1 - 進場訊號擴充 (P1) 🎯 MVP

**Goal**: 進場訊號包含原有 KD/MACD/VR-OBV 與新增的 12 個偵測器 OR 條件

**Independent Test**: 單獨呼叫各偵測器的 detect_latest_signal(df)，驗證回傳結構正確

### Tests for User Story 1

- [ ] T007 [US1] 測試 FourStepReversal 偵測器回傳 signal 欄位
- [ ] T008 [US1] 測試 FibonacciStrategyDetector 偵測器回傳 is_buy_signal 欄位
- [ ] T009 [US1] 測試 TopEscapeDetector 偵測器回傳 confirmed_escape_signal 欄位
- [ ] T010 [US1] 測試 SidewaysDropDetector 偵測器回傳 breakdown_signal 欄位
- [ ] T011 [US1] 全量回測 `python tmf_1m_engine.py`，確認交易次數與基準相近 (±5% 內)

## Phase 4: User Story 2 - 進場前共振過濾層 (P1)

**Goal**: check_resonance 過濾層正常運作，score<3 拒絕進場

**Independent Test**: 直接呼叫 check_resonance(df, action=1)，斷言回傳 (pass, score, reasons) 結構

### Tests for User Story 2

- [ ] T012 [US2] 測試 check_resonance 回傳 (pass: bool, score: int, reasons: list) 結構
- [ ] T013 [US2] 測試四維度趨勢未共振時 pass == False
- [ ] T014 [US2] 測試假突破/誘多時多單被否決
- [ ] T015 [US2] 測試逃命型態確認時多單被否決

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 四件套驗證與 collect

- [ ] T016 [P] 執行 `python -m py_compile tmf_1m_engine.py` 語法檢查通過
- [ ] T017 [P] 全量回測勝率與 PF 與基準相近 (±0.05 內)
- [ ] T018 [P] TICK socket 監聽 127.0.0.1:8798 功能保持不變
- [ ] T019 [US1] `/speckit.converge` 收斂驗證，確認 Converged

## Dependencies & Execution Order

- **Setup (T001-T002)**: 無相依，可立即執行
- **Foundational (T003-T006)**: 依賴 Setup，區塊所有後續
- **User Story 1 (T007-T011)**: 依賴 Foundational，可獨立測試
- **User Story 2 (T012-T015)**: 依賴 Foundational，可獨立測試
- **Polish (T016-T019)**: 依賴所有上述完成

## Parallel Opportunities

- T007-T015 各測試案例可平行撰寫 (不同偵測器)
- T016-T018 四件套可平行執行 (不同命令)

## Notes

- [P] 任務可平行
- 每個測試先寫、先失敗，再實作讓其通過 (TDD)
- Commit 後再執行 converge
