# Feature Specification: TMF 參數配置驗證 (Config Validation)

**Feature Branch**: `[001-tmf-config-validation]`

**Created**: 2026-11-13

**Status**: Draft

**Input**: User description: "為 tmf_engine.py 新增 TMF 參數配置驗證功能，在程式啟動或請求時檢查 TMFSpec 與 Account 參數的合理性（保證金足夠、風險比例、停損範圍、熔斷點數、最大回撤等），並回傳結構化驗證報告。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 啟動時自動驗證參數 (Priority: P1)

[Describe this user journey in plain language]
TMFDecisionEngine 初始化時，自動檢查 TMFSpec 與 Account 參數是否符合業務規則（如權益是否足夠 1.5 倍保證金、單筆風險是否不超過 5%、停損範圍是否在 40-80 點等），若發現不合理參數，在 Signal.reason 中標記並回傳警告。

**Why this priority**: 防止因參數設定錯誤導致交易引擎以不當風險啟動，是風控的第一道防線。

**Independent Test**: 可透過建立不同參數的 TMFSpec/Account 實例，呼叫驗證函數，檢查回傳的 valid / warnings / errors 是否正確，無需完整回測資料。

**Acceptance Scenarios**:

1. **Given** TMFSpec(original_margin=16100), Account(equity=75000), **When** TMFDecisionEngine 初始化， **Then** 驗證通過，valid=True，無錯誤。
2. **Given** Account(equity=10000) 小於 1.5 倍原始保證金 24150， **When** 初始化， **Then** valid=False，errors 包含保證金不足。
3. **Given** TMFSpec(min_stop_points=40, max_stop_points=80), **When** 呼叫驗證檢查停損範圍， **Then** 若 min > max 則回傳錯誤。

### User Story 2 - 按需回傳結構化驗證報告 (Priority: P2)

[Describe this user journey in plain language]
提供獨立函數 validate_config(spec, account)，回傳 dict 包含 valid(bool)、warnings(list)、errors(list)、checks(list of {name, passed, detail})，供外部程式或 logs 檢視。

**Why this priority**: 讓監控程式、日誌系統、或人工審核能以結構化方式取得參數合理性結果，不需解析文字訊息。

**Independent Test**: 直接呼叫 validate_config()，斷言回傳 dict 結構與欄位正確。

**Acceptance Scenarios**:

1. **Given** 正常參數， **When** validate_config()， **Then** 回傳 dict，valid=True，errors 為空 list。
2. **Given** 異常參數， **When** validate_config()， **Then** 回傳 dict，valid=False，errors 非空，每個 error 有 name 與 detail。

---

### Edge Cases

- TMFSpec.tick_value = 0 導致 position_size 除以零，應回傳錯誤。
- Account.max_risk_pct > 0.1 (10%) 視為過高風險，應回傳警告而非錯誤。
- daily_loss_limit_points 與 weekly_loss_limit_points 為正數但邏輯上應為負值比較，應檢查。
- peak_equity < equity（異常狀態），應回傳警告。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST 提供 `validate_config(spec: TMFSpec, account: Account) -> dict` 函數，回傳結構化驗證結果。
- **FR-002**: System MUST 在 TMFDecisionEngine.__init__ 時自動呼叫 validate_config，將結果附於 engine.config_validation 屬性。
- **FR-003**: System MUST 檢查權益是否 >= original_margin * 1.5，不足則 valid=False。
- **FR-004**: System MUST 檢查 max_risk_pct 是否在 0 < pct <= 0.05 範圍，超過 0.05 則為警告。
- **FR-005**: System MUST 檢查 min_stop_points <= max_stop_points，違反則 valid=False。
- **FR-006**: System MUST 檢查 tick_value > 0，等於 0 則 valid=False。
- **FR-007**: System MUST 回傳 dict 結構包含：valid (bool), warnings (list), errors (list), checks (list of dict)。
- **FR-008**: System MUST 每個 check item 包含 name (str), passed (bool), detail (str)。

### Key Entities *(include if feature involves data)*

- **ValidationResult**: dict 結構，包含 valid/warnings/errors/checks 欄位。
- **CheckItem**: dict 結構，包含 name/passed/detail 欄位。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 正常參數下 validate_config 回傳 valid=True，durations < 1ms。
- **SC-002**: 異常參數（如保證金不足）下 valid=False，errors 非空。
- **SC-003**: TMFDecisionEngine 初始化後，engine.config_validation 為 dict 且非 None。
- **SC-004**: 現有 `python tmf_engine.py` 主程式仍可正常執行，不產生例外。

## Assumptions

- 現有 tmf_engine.py 中的 TMFSpec 與 Account dataclass 結構不變。
- 驗證函數為純函數，不修改 spec 或 account 狀態。
- 不依賴外部資料庫或網路，完全在記憶體中執行。