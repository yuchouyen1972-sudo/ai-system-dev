# Implementation Plan: tmf_1m_addons 偵測器整合進 tmf_1m_engine 即時 TICK 引擎

**Branch**: `[003-addon-detectors-integration]` | **Date**: 2026-11-17 | **Spec**: `specs/.specify/memory/spec.md`

**Input**: Feature specification from `specs/003-addon-detectors-integration/spec.md`

## Summary

將 tmf_1m_addons.py 的 12 個偵測器與共振過濾層整合進 tmf_1m_engine.py 即時 TICK 引擎，作為進場 OR 條件與進場前過濾層，不改變 TICK socket 監聽與 K 線彙總功能。

## Technical Context

**Language/Version**: Python 3.11

**Primary Dependencies**: numpy, pandas, scipy.signal.argrelextrema (現行 tmf_1m_engine.py 已有)，tmf_1m_addons.py（新增偵測器模組）

**Storage**: N/A (純記憶體偵測，不寫檔)

**Testing**: 單元測試 + 全量回測 (TMF_1m.parquet)

**Target Platform**: Linux server / Windows dev machine

**Project Type**: library + script (tmf_1m_engine.py 作為模組被匯入，亦可直接執行回測)

**Performance Goals**: 每根 K 線偵測時間 < 50ms，不影響即時交易延遲

**Constraints**: 不修改現有 KD(14,3,3)、MACD(8,34,9)、VR(34)+OBV 進場邏輯；不修改 TICK socket 監聽與 K 線彙總功能；新增功能採 OR 關係，不減少現有交易

**Scale/Scope**: 
- tmf_1m_engine.py 新增約 80-100 行程式碼（導入、偵測器實例化、OR 條件、over濾層呼叫）
- tmf_1m_addons.py 已存在（12 個偵測器 + check_resonance）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- 規格優先：✅ 已有 spec.md
- 可維護性：✅ 偵測器類別獨立、職責單一、有 docstring
- 測試優先：✅ 將建立對應單元測試與全量回測驗證
- 安全與隱私：✅ 不涉及憑證、不寫檔、不外送、不修改現有進場邏輯
- SDD 流程遵守：✅ constitution → specify → plan → tasks → implement → converge

## Project Structure

### Documentation (this feature)

```text
specs/003-addon-detectors-integration/
├── spec.md              # 規格 (已建立)
├── plan.md              # 計畫 (本檔)
└── tasks.md             # 任務清單 (下一步)
specs/.specify/memory/
├── spec.md              # 規格副本 (指向本 spec)
├── plan.md              # 計畫副本 (指向本 plan)
└── tasks.md             # 任務清單副本 (指向本 tasks)
```

### Source Code (repository root)

```text
C:\Users\yp86n\AppData\Local\mtxf-local\scripts\tmf_1m_engine.py    # 新增 12 個偵測器整合 + check_resonance 過濾層
C:\Users\yp86n\AppData\Local\mtxf-local\scripts\tmf_1m_addons.py    # 已存在 12 個偵測器 + check_resonance
```

**Structure Decision**: 單一專案結構，偵測器邏輯在 tmf_1m_addons.py，整合邏輯在 tmf_1m_engine.py，測試透過全量回測驗證。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | N/A | N/A |
