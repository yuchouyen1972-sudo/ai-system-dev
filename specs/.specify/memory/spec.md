# Feature Specification: tmf_1m_addons 偵測器整合進 tmf_1m_engine 即時 TICK 引擎

**Feature Branch**: `[003-addon-detectors-integration]`

**Created**: 2026-11-17

**Status**: Draft

**Input**: User description: "將 tmf_1m_addons.py 的 12 個偵測器與共振過濾層整合進 tmf_1m_engine.py 即時 TICK 引擎，作為進場 OR 條件與進場前過濾層，不改變 TICK socket 監聽與 K 線彙總功能"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 進場訊號擴充 (Priority: P1)

[Describe this user journey in plain language]

即時 TICK 引擎在每根 1m K 線收盤時，除了現有的 KD/MACD/VR-OBV 五種型態進場訊號外，額外偵測 12 種缺失功能偵測器（EMA 缺口、趨勢線突破、四步反轉、大漲四大特徵、水平阻力位突破、真假突破、誘空狀態機、四維度趨勢、突破回踩、逃命型態、橫久必跌、斐波那契簡化版）。若任一偵測器觸發買進訊號，即可進場做多（OR 關係），不減少現有交易。

**Why this priority**: 擴充進場訊號來源，提升對不同市場型態的捕捉能力，同時保留現有績效基準（勝率 50.72%、PF 1.05、MDD 51,110）。

**Independent Test**: 可透過 `generate_tmf_forecast(df_1m)` 與各偵測器的 `detect_latest_signal(df)` 獨立呼叫，驗證訊號正確性，無需完整回測資料。

**Acceptance Scenarios**:

1. **Given** 1m K 線資料符合四步反轉條件，**When** 呼叫 `FourStepReversal().analyze(df)`，**Then** 回傳 `signal == "LONG"` 或 `"SHORT"`
2. **Given** 1m K 線資料符合斐波那契黃金區域條件，**When** 呼叫 `FibonacciStrategyDetector().detect_latest_signal(df)`，**Then** 回傳 `is_buy_signal == True`
3. **Given** 1m K 線資料出現逃命型態（墓碑線/射擊之星+次根確認），**When** 呼叫 `TopEscapeDetector().detect_latest_signal(df)`，**Then** 回傳 `confirmed_escape_signal == True`
4. **Given** 現有 KD/MACD/VR-OBV 訊號與新偵測器訊號同時觸發，**When** `generate_trading_signal(df)`，**Then** `action == 1` 且 `reason_code` 為 90/96/97 之一

### User Story 2 - 進場前共振過濾層 (Priority: P1)

[Describe this user journey in plain language]

進場前經過 `check_resonance(df, action)` 過濾層，對多單進行四步反轉方向一致性、大漲四大特徵（min 3 項）、四維度趨勢（多頭方向＋健康階段＋強能量＋級別充足）、真假突破（假突破/誘多否決）、逃命型態（否決多單）、橫久必跌（盤久必跌否決）等過濾。score >= 6 才通過，3-5 為 LOW 信心但允許進場，< 3 直接拒絕。

**Why this priority**: 作為安全網，避免在明顯反轉或盤整破位時進場，不影響現有績效但增加風控層。

**Independent Test**: 可透過 `check_resonance(df, action=1)` 獨立呼叫，驗證 pass/score/reasons 輸出。

**Acceptance Scenarios**:

1. **Given** 四維度趨勢顯示震盪（SIDEWAYS_CHOP），**When** `check_resonance(df, action=1)`，**Then** `pass == False`，reasons 包含「四維度未共振」
2. **Given** 真假突破偵測到假突破/誘多，**When** `check_resonance(df, action=1)`，**Then** `pass == False`，reasons 包含「真假突破：偵測到假突破/誘多」
3. **Given** 逃命型態確認訊號觸發，**When** `check_resonance(df, action=1)`，**Then** `pass == False`，reasons 包含「逃命型態」
4. **Given** 所有過濾條件通過，score >= 6，**When** `check_resonance(df, action=1)`，**Then** `pass == True`

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST 將 `tmf_1m_addons.py` 的 12 個偵測器類別整合進 `tmf_1m_engine.py` 的 `generate_trading_signal` 函數
- **FR-002**: System MUST 將突破回踩（`BreakoutPullbackDetector`）作為 OR 條件進場，reason_code = 96
- **FR-003**: System MUST 將斐波那契策略（`FibonacciStrategyDetector`）作為 OR 條件進場，reason_code = 97
- **FR-004**: System MUST 將斐波那契簡化版（`Fib1MinAnalyzer`）作為獨立偵測器，不直接進場，僅供終端輸出
- **FR-005**: System MUST 將 `check_resonance` 作為進場前過濾層，score < 3 拒絕進場，3-5 為 LOW 信心但允許，>= 6 通過
- **FR-006**: System MUST 不改變 TICK socket 監聽（127.0.0.1:8798）與 K 線彙總功能
- **FR-007**: System MUST 語法檢查通過（`py_compile`）
- **FR-008**: System MUST 全量回測績效與現有基準相近（勝率 50.72%、PF 1.05、MDD 51,110 範圍內）

### Key Entities *(include if feature involves data)*

- **BreakoutPullbackDetector**: 偵測突破回踩訊號，回傳 signal/box_high/stop_loss/note
- **FibonacciStrategyDetector**: 偵測斐波那契黃金區域買入訊號，回傳 is_buy_signal/zone/stop_loss_ref/current_price
- **Fib1MinAnalyzer**: 簡化版斐波那契分析，回傳 in_gold_zone/is_overpriced/is_below_gold_zone
- **TopEscapeDetector**: 偵測墓碑線/上吊線/射擊之星逃命訊號，回傳 confirmed_escape_signal/signal_description
- **SidewaysDropDetector**: 偵測橫久必跌破位訊號，回傳 breakdown_signal/bear_warning_score/box_low
- **check_resonance**: 整合過濾層，回傳 (pass: bool, score: int, reasons: list)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `python -m py_compile tmf_1m_engine.py` 語法檢查通過
- **SC-002**: 全量回測交易次數與基準相近（±5% 內）
- **SC-003**: 全量回測勝率與 PF 與基準相近（±0.05 內）
- **SC-004**: TICK socket 監聽 127.0.0.1:8798 功能保持不變
- **SC-005**: 12 個偵測器皆可獨立呼叫 `detect_latest_signal(df)` 回傳正確結構

## Assumptions

- `tmf_1m_addons.py` 已建立且包含 12 個偵測器類別與 `check_resonance` 函數
- `tmf_1m_engine.py` 為即時 TICK 版引擎，具有 socket 監聽與 K 線彙總功能
- 整合後的 `tmf_1m_engine.py` 同時具備即時交易與回測功能
- 不修改現有 KD(14,3,3)、MACD(8,34,9)、VR(34)+OBV 進場邏輯
