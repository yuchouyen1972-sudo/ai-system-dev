# Implementation Plan: 外部策略構想吸收——U9 場景框架／U10 FVG／U11 掃蕩收回

**Branch**: `004-external-idea-detectors` | **Date**: 2026-11-22 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-external-idea-detectors/spec.md`

## Summary

將兩套實測不可交易的外部 1m 策略腳本中之可取構想（killzone 時間窗、FVG 缺口、流動性掃蕩收回）重實作為三個獨立偵測器（U9/U10/U11），沿用既有 `ChannelPatternDetector` 的「批量無狀態」API 範式；以向量化批次統計在完整歷史 639,323 根上量化每個構想的原始優勢，並以實證成本區間（4.4~6.4 點）給出誠實可交易性判定。外部腳本程式碼不採用（其回測器有已證實的時序偽象與樂觀成本缺陷）。

## Technical Context

**Language/Version**: Python 3.11（mtxf-system venv，uv run）

**Primary Dependencies**: pandas（Arrow-backed str）、numpy、pyarrow（parquet 唯讀）；GUI 層 PySide6＋pyqtgraph（既有）

**Storage**: 1m parquet store 唯讀（`mtxf-system/data_store/TMF_1m.parquet`，ts=epoch 秒；`mtxf-local/...`，ts=datetime64——批次統計兩店皆可，擇一並記錄）

**Testing**: 專案自製腳本範式（同 `scripts/test_interaction.py` 的 check()/FAILS 樣式，offscreen GUI 測試），全數終端輸出、本地可重複

**Target Platform**: Windows 11 x64 桌面（pro_terminal）＋CLI 腳本

**Project Type**: desktop-app（PySide6）＋CLI 分析腳本（單一 repo：mtxf-system）

**Performance Goals**: 三偵測器全歷史批次統計 < 5 分鐘（FR-007）；GUI 疊圖不使重畫超出既有幀预算

**Constraints**: 1m K+tick only（禁非 1m 重採樣決策）；四層架構——偵測器屬顯示輔助層，不干預決策引擎；輸出僅終端/檔案；名稱禁 wawa；所有常數命名並註來源

**Scale/Scope**: 639,323 根 × 3 偵測器；1 個新模組、2 個新腳本、pro_terminal 最小掛接（3 個疊圖開關）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | 憲法原則 | 通過/依據 |
|---|----------|-----------|
| 1 | 規格優先 | ✓ spec.md 已提交（Default Project repo `30edd98`），本計畫僅實作其 FR-001~FR-011 |
| 2 | 可維護性 | ✓ 新模組含模組級 docstring；公開 API 全部型別標註；所有常數命名並註來源（FR-009） |
| 3 | 測試優先 | ✓ 單元＋整合測試（scripts/test_idea_detectors.py）；回測/驗證邏輯在歷史數據（639,323 根）上執行並產出權益曲線（cumulative expectancy）與前後對照（原始 vs 扣成本） |
| 4 | 安全與隱私 | ✓ 產出僅小體積 CSV/MD；parquet/log 不入規格 repo |
| 5 | SDD 流程 | ✓ specify→plan（本件）→tasks→implement→converge |

**Post-Phase-1 re-check**: 無違規——data-model 全部落在單一新模組＋2 腳本，未引入新依賴、新儲存或跨層耦合。GATE PASS。

## Project Structure

### Documentation (this feature)

```text
specs/004-external-idea-detectors/
├── plan.md              # 本件
├── research.md          # Phase 0：D1~D8 決策
├── data-model.md        # Phase 1：實體與驗證規則
├── quickstart.md        # Phase 1：端到端驗證指南
├── checklists/requirements.md
└── tasks.md             # Phase 2（/speckit.tasks 產出）
```

### Source Code (repository root: mtxf-system)

```text
visualization/
├── pattern_detector.py     # 既有 ChannelPatternDetector（不動，FR-010）
├── idea_detectors.py       # 新增：U9SessionContext / U10FvgDetector / U11SweepReclaimDetector
└── pro_terminal.py         # 最小掛接：疊圖開關×3（僅顯示層，沿 2971-2979 掛接樣式）

scripts/
├── validate_ideas_004.py   # 新增：全歷史批次統計 → docs/validation_004/（CSV+MD）
└── test_idea_detectors.py  # 新增：單元＋整合＋回歸測試（含 SC-003/SC-004/SC-005）

docs/
└── validation_004/         # 批次統計報告輸出（CSV＋Markdown）
```

**Structure Decision**: 偵測器放獨立新模組 `idea_detectors.py` 而非擴充 `pattern_detector.py`——FR-010 要求不動既有檔案；pro_terminal 掛接點維持既有「paint 前批量 update、回傳 dict 供標註」樣式（2971-2979）。 contracts/ 略過——本功能為內部工具，無外部介面。

## Complexity Tracking

> 無憲法違規，本表留空。
