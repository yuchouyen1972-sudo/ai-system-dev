# Quickstart: 004-external-idea-detectors 驗證指南

**前置**：`C:\Users\yp86n\AppData\Local\mtxf-system`；`data_store/TMF_1m.parquet`（639,323 根）；venv 就緒。

## 1) 單元＋整合＋回歸測試（SC-003/SC-004/SC-005）

```powershell
cd C:\Users\yp86n\AppData\Local\mtxf-system
& C:\Users\yp86n\.local\bin\uv.exe run python scripts\test_idea_detectors.py
```

**期望輸出**：逐項 `[PASS]`，末行 `結果：ALL PASS`。涵蓋：
- U10：合成三根 K 形成缺口 → 方向/上下緣/寬度正確；雜訊缺口（< 下限）不標記
- U11：合成掃蕩收回 → 方向/池值正確；當根不入池（無前視）
- U9：合成含週五夜盤片段 → 週五 15:00 後與週六凌晨棒之 trading_day 皆為週一；補班六語義（不生成）
- 時序偽象斷言：進場根集合 ∩ 出場根集合 = ∅（SC-003）
- 回歸：既有 `scripts/test_interaction.py` 15 檢查項維持 ALL PASS（SC-005）

## 2) 全歷史誠實驗證（SC-001/SC-002/SC-006）

```powershell
& C:\Users\yp86n\.local\bin\uv.exe run python scripts\validate_ideas_004.py
```

**期望輸出**：
- 終端印出各偵測器摘要（筆數/勝率/原始點數/扣成本淨期望/verdict）
- 產出 `docs/validation_004/ideas_summary.csv` 與 `validation_report_004.md`
- 耗時 < 5 分鐘（SC-001）；verdict 依成本區間雙值判定（SC-006）

## 3) GUI 疊圖人工驗收（US3）

```powershell
& C:\Users\yp86n\.local\bin\uv.exe run python visualization\pro_terminal.py --selftest
```

**期望**：`SELFTEST_OK`。人工步驟（桌面）：進入回放 → 開啟三個疊圖開關 → 拖曳時光機至含缺口/掃蕩事件時段 → 標註正確呈現且切換畫格無可感知延遲。

## 4) 設計對照

- 偵測器 API 範式對照：`visualization/pattern_detector.py`（批量無狀態 `update()` 回傳 dict）
- 掛接點對照：`visualization/pro_terminal.py` 2971-2979
- 憲法對照：`specs/CONSTITUTION.md`（測試優先：歷史數據驗證＋equity curve 對照 → 本指南第 2 步之累積期望曲線）
