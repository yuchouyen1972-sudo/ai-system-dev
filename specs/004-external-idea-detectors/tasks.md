# Tasks: 004-external-idea-detectors

**Input**: Design documents from `/specs/004-external-idea-detectors/`（plan / research / data-model / quickstart）

**Prerequisites**: plan.md ✓、spec.md ✓、research.md ✓（D1~D8）、data-model.md ✓、contracts/（略——內部工具）

**狀態**：T001~T014 全數完成；converge 二輪＝**converged**（2026-11-22，F1/US3 已解）

**Tests**: 憲法「測試優先」明文要求——本清單含測試任務，且**測試任務先於實作任務**（先 FAIL 後 PASS）。

**Organization**: 依 spec.md 三個 User Story 分階段；路徑根＝`C:\Users\yp86n\AppData\Local\mtxf-system`（實作 repo）；規格 repo＝`C:\Users\yp86n\OneDrive\文件\Default Project`。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 不同檔案、無未完成依賴，可平行
- **[Story]**: 對應 spec.md 之 US1/US2/US3

---

## Phase 1: Setup

- [x] T001 建立 `visualization/idea_detectors.py` 骨架：模組級 docstring（構想來源＝兩套外部腳本之 killzone/FVG/sweep-reclaim；鐵則聲明＝1m K only、顯示輔助層不干預決策引擎、輸出僅終端/檔案）；全部常數依 research D6 命名並註來源（KILLZONE_TW=(08:45,09:30)、KILLZONE_US=(21:30,23:30)、OPENING_RANGE_MIN=15、LIQ_POOL_LOOKBACK=20、FVG_MIN_POINTS=3.0、FVG_ATR_MULT=0.25、ATR_PERIOD=14、HORIZONS=(15,30,60)、POINT_VALUE_TWD=10、FEE_RT_TWD_RANGE=(30,50)、TAX_RATE=0.00002）；三個類別簽名（型別標註、回傳 dataclass/TypedDict），方法體暫以 NotImplementedError 佔位

---

## Phase 2: Foundational

- [x] T002 [P] 建立 `scripts/test_idea_detectors.py` 測試骨架：沿用 `scripts/test_interaction.py` 的 check()/FAILS 範式（offscreen 不需要——本檔純邏輯測試）；合成資料產生器：①含「週四~週一」跨週末的 1m 片段 ②三根 K 向上缺口情境 ③掃蕩前低收回情境 ④小於下限的雜訊缺口情境

**Checkpoint**: 骨架與測試台就緒，User Story 實作可開始

---

## Phase 3: User Story 1 — 三個新偵測器獨立運作 (Priority: P1) ★ MVP

**Goal**: U9/U10/U11 各自可獨立呼叫並回傳 data-model.md 定義之結構化結果

**Independent Test**: `uv run python scripts/test_idea_detectors.py` 中 US1 區段逐項 PASS（合成資料即可，不需全史）

### Tests for User Story 1（先寫，先 FAIL）

- [x] T003 [US1] 於 `scripts/test_idea_detectors.py` 撰寫 US1 單元測試：①U9：週五 15:00 後與週六 00:00~04:59 棒之 trading_day 皆為「次週一」；週六日日盤時段 session=CLOSED 且無 killzone；②U10：情境②形成 BULL 缺口且 top/bottom/width 正確、情境④不標記；③U11：情境③標記 SWEEP_LOW_RECLAIM 且 pool_ref 等於前 20 根最低、當根不入池（pool 以 shift(1) 起算）；④回傳結構欄位完整性（對照 data-model.md 欄位表）。**執行確認 FAIL**（NotImplementedError）

### Implementation for User Story 1

- [x] T004 [US1] 實作 `U9SessionContext`（`visualization/idea_detectors.py`）：FR-003 規則歸屬（research D1 向量化，禁 +9h 近似）；session/killzone/開盤區間/bias 依 data-model.md 欄位規則
- [x] T005 [US1] 實作 `U10FvgDetector`：三根 K 結構、gap_min=max(FVG_MIN_POINTS, ATR14×FVG_ATR_MULT)、窗口內回補追蹤（research D2 無狀態窗口語義）
- [x] T006 [US1] 實作 `U11SweepReclaimDetector`：LIQ_POOL_LOOKBACK=20 池值以 shift(1) 起算；兩方向事件
- [x] T007 [US1] 執行 `uv run python scripts/test_idea_detectors.py` → US1 區段全部轉 PASS（修至綠）

**Checkpoint**: 三偵測器獨立可用（MVP 可驗收）

---

## Phase 4: User Story 2 — 全歷史誠實驗證與可交易性判定 (Priority: P2)

**Goal**: 全史 639,323 根事件研究報告（CSV+MD）＋成本區間雙值 verdict

**Independent Test**: `uv run python scripts/validate_ideas_004.py` 跑畢並產出報告；數字可由獨立抽樣複核

### Tests for User Story 2（先寫，先 FAIL）

- [x] T008 [US2] 於 `scripts/test_idea_detectors.py` 撰寫整合斷言：①時序偽象=0——以小樣本跑統計管線，斷言進場根集合 ∩ 出場根集合 = ∅（SC-003）；②交易日歸屬——對全史抽樣（每季首日）斷言週五夜盤棒 100% 歸屬次週一（SC-004）。**執行確認 FAIL**（統計管線未存在）

### Implementation for User Story 2

- [x] T009 [US2] 實作 `scripts/validate_ideas_004.py`：讀 `data_store/TMF_1m.parquet`（ts=epoch 秒橋接，research D7）；三偵測器全史掃描；事件研究（research D5：次根開盤進場、H∈HORIZONS 遠景收盤出場——D4 根除同根偽象）；CostModel 雙值（D3）；killzone 內/外切層對照；累積期望曲線 maxDD；含 params_snapshot（FR-011）
- [x] T010 [US2] 產出 `docs/validation_004/ideas_summary.csv`＋`docs/validation_004/validation_report_004.md`（verdict 三分法：淨期望下限>0 可交易／上限<0 不可交易／跨零待調台——FR-008 僅檔案與終端）；終端印摘要；記錄耗時（SC-001 <5 分）；執行 T008 斷言轉 PASS

**Checkpoint**: 誠實判定出爐——**若判定「不可交易」，如實記錄，不得以調參或樂觀成本掩飾**（SC-006）

---

## Phase 5: User Story 3 — 終端圖表呈現 (Priority: P3)

**Goal**: pro_terminal 三個疊圖開關（killzone 底色／缺口區塊＋回補狀態／掃蕩標記）

**Independent Test**: `--selftest` SELFTEST_OK＋互動測試 15 項 ALL PASS＋人工回放目視

### Implementation for User Story 3

- [x] T011 [US3] `visualization/pro_terminal.py` 最小掛接：沿 2971-2979 既有樣式（paint 前批量呼叫、回傳 dict 供標註）；三個圖層各為獨立 pyqtgraph item、可開關；不動既有 U1-U8 圖層與任何既有行為（FR-010）；資料不足或 CLOSED 時段不繪製
- [x] T012 [US3] 回歸驗證：`uv run python visualization/pro_terminal.py --selftest` → SELFTEST_OK；`uv run python scripts/test_interaction.py` → 15 檢查項 ALL PASS（SC-005）

**Checkpoint**: 三個 User Story 皆獨立可用

---

## Phase 6: Polish & Cross-Cutting

- [x] T013 執行 quickstart.md 全部步驟留檔（測試/全史報告/selftest）；HANDOFF F-87 入帳（mtxf-system repo commit：idea_detectors.py＋2 腳本＋pro_terminal 掛接＋validation_004 報告；Default Project repo commit：tasks.md 勾選狀態）；準備 `/speckit.converge` 收斂驗證

---

## Dependencies & Execution Order

### Phase Dependencies
- Setup（T001）→ Foundational（T002）→ US1（T003~T007）→ US2（T008~T010）→ US3（T011~T012）→ Polish（T013）
- US2 與 US3 **相互獨立**，皆依賴 US1 完成（T007 checkpoint 後可平行）

### Within US1
- T003 測試先行 → T004/T005/T006 同檔**序列**實作（避衝突）→ T007 轉綠

### Parallel Opportunities
- T001 與 T002 不同檔可平行（[P] 標 T002）
- US1 完成後：US2（CLI 統計）與 US3（GUI 掛接）可由不同人平行
- T004~T006 屬同檔序列，不標 [P]（避免同檔衝突）

---

## Implementation Strategy

### MVP First（US1 only）
1. T001~T002 骨架 → 2. T003 測試 FAIL → 3. T004~T006 實作 → 4. T007 ALL PASS → **STOP 驗收 US1**

### 誠實驗證紀律（本功能之核心價值）
- T010 的 verdict 出爐**先於**任何調參討論；「待調參」須在報告列出現參數與替代參數建議，調參後重跑 T008~T010 全鏈
- 禁止任何以 <4 點成本、>H=60 遠景、或樣本切割方式美化數字的行為（FR-004/FR-005/FR-011）

### Notes
- 每完成一任務或邏輯組即 commit（mtxf-system repo；specs 勾選狀態同步 Default Project repo）
- 所有終端輸出 UTF-8（`io.TextIOWrapper` 樣式）；不建立任何網頁產物

---

## Phase 7: Convergence（2026-11-22 converge 產出）

**Purpose**: 收斂盤點 11 FR＋6 SC＋10 AC＋8 研究決策＋5 憲法原則——9/10 AC 過關，餘 1 項 US3 互動驗證 partial

- [x] T014 [US3] 補 US3 疊圖互動自動化驗證：①offscreen 下三開關（KZ/FVG/掃蕩）`setChecked` 切換後 plot items 正確增減且不影響既有圖層（per US3/AC1-2，partial）②paint 含疊圖之單幀耗時抽樣（per US3/AC3「無可感知延遲」量化佐證，partial）③執行 quickstart.md §3 人工步驟並留檔輸出（per plan: 2971-2979 掛接決策）
