# Implementation Plan: TMF 參數配置驗證

**Branch**: `[001-tmf-config-validation]` | **Date**: 2026-11-13 | **Spec**: `specs/.specify/memory/spec.md`

**Input**: Feature specification from `specs/.specify/memory/spec.md`

## Summary

為 tmf_engine.py 新增 `validate_config(spec, account)` 純函數，回傳結構化驗證報告 (valid/warnings/errors/checks)，並在 TMFDecisionEngine.__init__ 時自動執行並附於 `engine.config_validation`。驗證項目包含保證金足夠、風險比例上限、停損範圍、tick_value 非零、peak_equity 異常等。

## Technical Context

**Language/Version**: Python 3.11

**Primary Dependencies**: Python stdlib dataclasses, typing (tmf_engine 已有 numpy/pandas，但驗證函數僅需 stdlib)

**Storage**: N/A (純記憶體驗證，不寫檔)

**Testing**: pytest（現行四件套）

**Target Platform**: Linux server / Windows dev machine

**Project Type**: library (tmf_engine.py 作為模組被匯入)

**Performance Goals**: 驗證函數執行 < 1ms

**Constraints**: 不修改現有 TMFSpec / Account 欄位定義；不破壞現有 `python tmf_engine.py` 主程式行為

**Scale/Scope**: 單一檔案 tmf_engine.py 新增約 30-40 行程式碼，單一測試檔案

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- 規格優先：✅ 已有 spec.md
- 可維護性：✅ 函數名稱明確、欄位結構化、有 docstring
- 測試優先：✅ 將建立對應單元測試
- 安全與隱私：✅ 不涉及憑證、不寫檔、不外送
- SDD 流程遵守：✅ constitution → specify → plan → tasks → implement → converge

## Project Structure

### Documentation (this feature)

```text
specs/.specify/memory/
├── spec.md              # 規格 (已建立)
├── plan.md              # 計畫 (本檔)
└── tasks.md             # 任務清單 (下一步)
```

### Source Code (repository root)

```text
tmf_engine.py            # 新增 validate_config 函數 + TMFDecisionEngine.__init__ 整合
tests/test_config_validation.py   # 新增單元測試 (四件套驗證)
```

**Structure Decision**: 單一專案結構，驗證函數直接加入 tmf_engine.py，測試放在根目錄 tests/。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | N/A | N/A |