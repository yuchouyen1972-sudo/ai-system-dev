# Tasks: 四維度趨勢分析器（FourDimensionTrendAnalyzer）

**Input**: Design documents from `/specs/001-fourdim-trend-analyzer/`
**Prerequisites**: plan.md (completed), spec.md (completed)
**Format**: `[ID] [P?] [Story] Description` — `[P]` = run in parallel (different files, no dependencies)

---

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 [P] [Setup] 建立 `advisor/` 目錄與 `advisor/__init__.py`（參考模組層入口，與 Elliott 相位模組 `elliott_phase.py` 同層）

---

## Phase 2: Foundational (Blocking Prerequisites)

**CRITICAL**: 以下 T002–T005 皆為同一檔案 `advisor/fourdim_analyzer.py` 的不同區塊，必須按順序完成，不可平行（避免 edit 衝突）。

- [ ] T002 [US1 P1] 實作 `advisor/fourdim_analyzer.py`：`FourDimensionTrendAnalyzer` 類別、`_calculate_indicators`（MACD/body_ratio/FVG）、`d1_direction`、`d2_stage`、`d3_energy`、`d4_level`、`analyze(df, timeframe="1m")` 回傳 dict（含 `trend_score` 0-100、`decision`）
- [ ] T003 [US2 P1] 修正 FVG 填補判斷：缺口 `[high[i-2], low[i]]` 若 i+1 或 i+2 的 `low` 跌破 `high[i-2]`，該 FVG 不計入有效未填補數（取代原兩半段 max 對比語意）
- [ ] T004 [US2 P1] 修正頂背離判斷：以「價格創高的當下高點」對應的 DIF 值與前波 DIF 峰值比較，而非兩半段 max 直接對比
- [ ] T005 [US1 P1] 加 score=None 語意：資料不足（< lookback*2 根）或零波動（max_dif=0）時，`trend_score=None`、`decision="資料不足/觀望"`，不拋例外、不回 0 分

---

## Phase 3: User Story 1 - 四維度趨勢參考評分 (P1 / MVP)

- [ ] T006 [US1 P1] 驗證 `analyze()` 回傳結構完整（d1_direction/d2_stage/d3_energy/d4_level/trend_score/decision），`trend_score` 為整數 0-100 或 None

---

## Phase 4: User Story 2 - 語意缺陷修正 (P1 / 品質門檻)

- [ ] T007 [US2 P1] 頂背離測試通過：價格創新高但 DIF 該高點已下降的序列，`d2_stage` 不得為「健康主升段」
- [ ] T008 [US2 P1] FVG 填補測試通過：缺口被 i+1/i+2 low 回測填補後，有效未填補 FVG 計數遞減

---

## Phase 5: User Story 3 - 併入參考層，forecast_engine 零改 (P2)

- [ ] T009 [P] [US3 P2] 驗證 `advisor/fourdim_analyzer.py` 匯出與 `is_decision_gate=False` 旗標存在；確認 `forecast_engine.py` 零 diff（不改任何決策閘門條件）

---

## Phase 6: 測試與驗證 (所有 story)

- [ ] T010 [P] [所有 story] 撰寫 `tests/test_fourdim_analyzer.py`：強多/強空/震盪/頂背離/FVG 填補/資料不足/零波動單元測試
- [ ] T011 [所有 story] 跑 `pytest tests/test_fourdim_analyzer.py` 並確認全綠（verdict pass）

---

## Phase 7: Convergence

- [ ] T012 [US3 P1] 定義 `is_decision_gate=False` 旗標於 `advisor/fourdim_analyzer.py`；完成 forecast_engine 參考層註冊（advisor 側引用，forecast_engine 原件零改）per FR-007/US3 (partial)
- [ ] T013 [US1 P1] 補強測試：強多頭斷言 `trend_score >= 60` 且 decision 含「可尋找進場點」；強空頭斷言 `trend_score <= 40` per SC-001/SC-002 (partial)

---

**依賴順序**：T001 → T002 → T003 → T004 → T005 → T006 → T007/T008 → T009 → T010 → T011 → T012 → T013。T001、T009、T010 可與前序不同檔案平行；T002–T005 同一檔案必須序列；T012 應在 T009 之後。
