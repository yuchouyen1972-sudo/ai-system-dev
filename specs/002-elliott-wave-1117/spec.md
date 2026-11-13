# Feature Specification: Elliott Wave 11-17 整合（空頭鏡射＋ZigZag 0.5%）

**Feature Branch**: `[002-elliott-wave-1117]`

**Created**: 2026-11-13

**Status**: Draft

**Input**: User description: "Elliott Wave 11-17 整合：補空頭鏡射（validate_bearish_impulse 三大鐵律對稱）＋zigzag 調參 0.3-0.5%（與 Elliott 相位模組 elliott_phase.py 同層，advisor 參考層，不進決策閘門）"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 空頭推動浪鏡射驗證 (Priority: P1)

身為 advisor 管線，需要對空頭推動浪結構進行鏡射驗證（三大鐵律：W2 不得超過 W0、W3 不得最短、W4 與 W1 不得重疊），輸出驗證結果與 Fibonacci 指標。

**Why this priority**: 現有 elliott_phase.py 僅多頭推動浪驗證，空頭側缺失，導致 11-17 空頭信號無對稱驗證。

**Independent Test**: 給定 6 個樞紐價格（起點 HIGH），呼叫 `validate_bearish_impulse(p)`，斷言回傳 `(is_valid, violations, fib_metrics)`，對違反鐵律的序列回傳 violations 非空。

**Acceptance Scenarios**:

1. **Given** 符合三大鐵律的空頭推動浪序列，**When** 呼叫 `validate_bearish_impulse`，**Then** `is_valid=True`，violations 為空。
2. **Given** W2 反彈高點超過 W0 起點的序列，**When** 呼叫 `validate_bearish_impulse`，**Then** `is_valid=False`，violations 包含「鐵律一」。

---

### User Story 2 - ZigZag 0.5% 參數調整與確認時點 (Priority: P1)

身為 advisor 管線，需要 ZigZag 樞紐偏差參數從 3%（5m 日內過大）調整為 0.5%（elliott_phase.py BASE_DEV），且樞紐只在其反轉被確認的棒上可見（無未來函數）。

**Why this priority**: 3% ZigZag 對 5m 日內資料太大，會漏掉短周期結構；0.5% 為 zigzag_tune 掃描結果，需作為預設參數。

**Independent Test**: 建立 ZigZag 實例，參數 0.005，輸入價格序列，斷言樞紐數量與位置符合 0.5% 偏差，且無未來函數（樞紐不依賴未來價格）。

**Acceptance Scenarios**:

1. **Given** zigzag_deviation_pct=0.005，**When** 執行樞紐偵測，**Then** BASE_DEV=0.005，樞紐僅在反轉確認後出現。
2. **Given** 同一序列以 3% 與 0.5% 執行，**When** 比較樞紐數量，**Then** 0.5% 產出更多短周期樞紐。

---

### User Story 3 - 併入 advisor 參考層，forecast_engine 零改 (Priority: P2)

將 Elliott 相位模組與空頭鏡射函數作為 advisor 參考層（與 fourdim_analyzer.py 同層），輸出評分與標籤；嚴禁進入決策閘門。

**Why this priority**: 與 fourdim_analyzer.py 同樣鐵律，advisor 純參考層，不改 forecast_engine 原件。

**Independent Test**: 檢查 advisor 模組匯出，斷言 Elliott 相關函數存在且 `is_decision_gate=False`；forecast_engine 零 diff。

**Acceptance Scenarios**:

1. **Given** advisor 模組載入，**When** 檢查模組，**Then** Elliott 函數存在且 is_decision_gate 為 False。
2. **Given** 現有決策引擎測試，**When** 執行，**Then** 全數通過、零變更。

---

### Edge Cases

- 樞紐數量少於 6 個時，`validate_bearish_impulse` 應回傳 `is_valid=False` 並註明資料不足。
- ZigZag 價格序列全相同（零波動）時，不應拋例外，應回傳空樞紐列表。
- 非 TMF 商品（不同點值）時，評分僅依賴價格結構，不硬編碼商品參數。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST 提供 `validate_bearish_impulse(p)` 函數，輸入 6 個樞紐，回傳 `(is_valid: bool, violations: list, fib_metrics: dict)`。
- **FR-002**: System MUST 執行三大鐵律檢查：W2 高點不得超過 W0；W3 不得為 1/3/5 中最短；W4 高點不得與 W1 低點重疊。
- **FR-003**: System MUST ZigZag 樞紐偏差參數預設為 0.5%（BASE_DEV=0.005），與 elliott_phase.py 一致，不使用 3%。
- **FR-004**: System MUST ZigZag 樞紐僅在反轉被確認的棒上可見（無未來函數）。
- **FR-005**: System MUST Elliott 相位模組置於 advisor 側（與 fourdim_analyzer.py 同層），forecast_engine 原件 MUST NOT 修改，具 `is_decision_gate=False` 旗標。
- **FR-006**: System MUST 保留原 `elliott_wave_engine.py` 完整不變（使用者原件完整保留）。

### Key Entities *(include if feature involves data)*

- **BearishImpulseResult**: tuple (is_valid, violations, fib_metrics)。
- **ZigZagPivot**: 樞紐點，含價格、時間索引、類型（HIGH/LOW）。
- **ModuleRegistration**: advisor 模組註冊項，含 is_decision_gate 旗標。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 符合三大鐵律的空頭序列，`validate_bearish_impulse` 回傳 `is_valid=True`。
- **SC-002**: 違反任一鐵律的序列，回傳 `is_valid=False` 且 violations 非空。
- **SC-003**: ZigZag 參數 BASE_DEV=0.005（非 0.03）。
- **SC-004**: forecast_engine 載入後決策閘門條件與載入前完全一致（現有測試全綠、零 diff）。
- **SC-005**: 資料不足（< 6 樞紐）或零波動輸入不拋例外。

## Assumptions

- 輸入價格序列由 advisor/tick 管線提供，欄位為 OHLCV。
- 點值、手續費、稅等商品參數不進入本模組（由決策層負責）。
- 本模組輸出為「參考標籤」，嚴禁作為交易依據（鐵律）；決策閘門不引用其輸出。
- 原 `elliott_wave_engine.py` 與 `elliott_phase.py` 完整保留於 mtxf-local 根目錄。