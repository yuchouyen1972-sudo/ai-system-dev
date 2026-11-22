# Feature Specification: TAIFEX 行事曆整合（共享休市日模組）

**Branch**: `006-taifex-calendar` | **Created**: 2026-11-22 | **Status**: Draft

**Input**: 掃描器 v3 與偵測器 U9 的已知限制「國定假日未內建」——以官方來源取得休市日資料，建立共享行事曆模組，消除假日假缺/假標註。

## Research（資料源與規則）

- **D1 資料源**：期交所 2026 行事曆 PDF（`file/taifex/CHINESE/4/2026Calendar.pdf`）＋官方公告（端午 6/19、勞動節 5/1、國慶補假 10/9、光復補假 10/26、颱風 7/10/10/23）＋英文 Holiday Schedule PDF——**多源交叉一致**；櫃買中心與兩第三方佐證。2027 開休市公告 PDF 已知來源（台期交字第1150003597號函），PDF 解析列後續。
- **D2 休市規則**：國定假日全日休（含盤後夜盤）；**夜盤存在條件＝次日曆日為交易日**（端午 6/19 休 → 6/18 晚盤休）；**封關日夜盤例外**（2/11 晚盤跨 2/12 05:00，雖 2/12 休市）——以 JSON `night_exceptions` 標註。
- **D3 範圍**：2026 全年（11 組假日＋2 颱風日）；2025 以前歷史資料之假日已由資料本身反映（假日本無棒，偵測器連續性條件免疫），不回填；2027 待 PDF 解析（揭露）。
- **D4 介面**：`taifex_calendar.py` 共享模組——`is_trading_day(d)`、`next_trading_day(d)`、`prev_trading_day(d)`、`trading_day_map()`（日曆日→交易日查表，供偵測器歸屬）；JSON 資料檔獨立於程式碼（更新假日不動程式）。

## Requirements

- **FR-001**: System MUST 提供共享行事曆模組（週末＋假日＋颱風臨時休市），JSON 資料檔驅動。
- **FR-002**: 掃描器 `check_data_gaps.py` 的營業日判定 MUST 改用共享模組（假日不生成期望；假日前一晚亦不生成，夜盤例外除外）。
- **FR-003**: 偵測器 `idea_detectors.py` 的交易日歸屬 MUST 改用共享模組查表（假日跨越正確歸屬至次一交易日）。
- **FR-004**: System MUST 保持既有測試全綠（掃描器 v3 語義不變、004/005 行為不回歸）。
- **FR-005**: 封關日夜盤例外 MUST 以資料檔欄位支援（night_exceptions）。
- **FR-006**: 所有輸出僅終端/檔案；常數與資料來源註明（多源交叉清單入 JSON provenance）。

## Success Criteria

- **SC-001**: 端午情境合成測試：6/18 晚、6/19 全日不生成期望；6/22 日盤正常（掃描器）。
- **SC-002**: 封關夜例外：2/11 晚間生成期望（掃描器）。
- **SC-003**: 偵測器歸屬：假日次日歸屬正確（6/19→6/22 型）；既有 US1/US2 測試不回歸。
- **SC-004**: 既有測試台 27 項＋005 測試＋互動 15 項 ALL PASS。

## Tasks

- [ ] T001 [P] 建立 `config/taifex_holidays.json`（2026 休市 11 組＋颱風 2 日＋night_exceptions＋provenance 多源清單）
- [ ] T002 [P] 建立 `visualization/taifex_calendar.py`：載入 JSON、is_trading_day/next/prev/trading_day_map（2024-2027 查表）
- [ ] T003 掃描器接線（FR-002）＋合成端午/封關夜測試（SC-001/002）
- [ ] T004 偵測器接線（FR-003）＋歸屬測試（SC-003）
- [ ] T005 回歸驗證（SC-004：測試台＋005＋互動）＋HANDOFF F-91＋雙 repo commit
