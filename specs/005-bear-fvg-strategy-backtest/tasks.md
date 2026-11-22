# Tasks: 005-bear-fvg-strategy-backtest-bear-fvg-strategy-backtest

**狀態**：T001~T008 全數完成（2026-11-22 17:29~18:1x）；verdict＝可執行（最嚴情境）＋安慰劑對照通過。

**Tests**: 憲法測試優先——T002 測試先行，先 FAIL 後 PASS。

## Phase 1: Setup
- [x] T001 建立 `scripts/strategy_bear_fvg.py` 骨架：docstring（來源=004 驗證結論；鐵則聲明）＋常數區（D1-D6 全參數命名註來源）＋`run_backtest(...)`/`build_report(...)` 簽名（NotImplementedError 佔位）

## Phase 2: User Story 1（P1）
- [x] T002 [US1] 先寫合成測試 `scripts/test_strategy_bear_fvg.py`：①BEAR 缺口+下跌→空單成立、進場=次根開盤−slip；②反彈觸停損→SL＋不利滑價；③跨 13:40→EOD_FLAT；④持倉中新訊號跳過（互斥）；⑤停損停利同根→保守 SL。**執行確認 FAIL**
- [x] T003 [US1] 實作 `run_backtest(ts, o, h, l, c, events, slip, ...)`：迴圈式單口滾動（D6）、SL=缺口上緣+buffer（D1）、TP=1.5R、H_MAX、強平（D3）、不利滑價（D2）
- [x] T004 [US1] 執行測試至 ALL PASS（SC-003）
## Phase 3: User Story 2（P2）
- [x] T005 [US2] 滑價×成本矩陣：3 slip × 2 cost 六格 expectancy＋verdict（最嚴格為準，FR-008）
## Phase 4: User Story 3（P3）
- [x] T006 [US3] 全史執行＋報告：逐筆 CSV＋summary MD（權益統計＋分年切片＋事件研究對照＋D5 歸因三項）＋耗時（SC-001/002/004/005）
- [x] T007 [US3] 回歸驗證：既有測試台 27 項 ALL PASS＋確認 idea_detectors/pro_terminal 零改動（SC-006、FR-009）
## Phase 5: Polish
- [x] T008 HANDOFF F-90 入帳＋WORKPLAN_ACTIVE 更新＋雙 repo commit＋converge 準備
