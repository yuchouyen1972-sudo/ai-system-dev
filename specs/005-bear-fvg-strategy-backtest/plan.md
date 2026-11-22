# Implementation Plan: BEAR FVG 可執行策略回測器

**Branch**: `005-bear-fvg-strategy-backtest` | **Date**: 2026-11-22 | **Spec**: [spec.md](./spec.md)

## Summary

以 004 已驗證的偵測器（U10 BEAR FVG）與成本模組為基礎，新增離線策略回測器：單口滾動互斥、缺口上緣停損、風報比停利、遠景上限、跨盤強平、滑價 0/1/2 情境；全史執行產出逐筆交易＋權益曲統計＋情境矩陣＋分年切片＋事件研究對照，verdict 以最嚴情境為準。

## Technical Context

**Language/Version**: Python 3.11（mtxf-system venv，uv run）
**Primary Dependencies**: pandas / numpy / pyarrow（唯讀）——沿用 004，零新依賴
**Storage**: 1m parquet 唯讀；輸出 docs/validation_005/（CSV+MD）
**Testing**: 專案自製 check()/FAILS 腳本範式（合成案例可手算複核）＋既有 27 項回歸
**Target Platform**: Windows 11 x64 CLI
**Performance Goals**: 全史 < 5 分鐘（FR-010）
**Constraints**: 1m K only；輸出僅終端/檔案；不改動既有行為（FR-009）；常數命名註來源
**Scale/Scope**: 639,323 根；1 新回測腳本＋1 新測試檔＋pro_terminal 零改動

## Constitution Check

| # | 原則 | 狀態 |
|---|------|------|
| 1 | 規格優先 | ✓ 本 spec 005 |
| 2 | 可維護性 | ✓ docstring＋型別＋常數來源（沿 idea_detectors 範式） |
| 3 | 測試優先 | ✓ 測試先行（合成案例手算複核）＋歷史數據驗證＋equity curve（本功能本體） |
| 4 | 安全與隱私 | ✓ 小檔案輸出 |
| 5 | SDD 流程 | ✓ specify→plan→tasks→implement→converge |

**Post-Phase-1 re-check**: 無違規（GATE PASS）——純新增 CLI 研究工具，零跨層耦合。

## Project Structure

```text
mtxf-system/
├── visualization/idea_detectors.py   # 既有（唯讀使用：detect_fvg/_atr/cost_pts）
├── scripts/
│   ├── strategy_bear_fvg.py          # 新增：回測器＋全史執行＋報告
│   └── test_strategy_bear_fvg.py     # 新增：合成案例測試（先寫先 FAIL）
└── docs/validation_005/              # 輸出（trades CSV＋summary MD）
```

**Structure Decision**: 回測器放 scripts/（離線研究工具，非 GUI/引擎元件——四層鐵則：顯示輔助層不干預決策引擎，本工具連 GUI 都不碰）。contracts/ 略（內部工具）。

## Complexity Tracking

> 無憲法違規，留空。
