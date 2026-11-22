# Data Model: 004-external-idea-detectors

**Date**: 2026-11-22｜上游：[spec.md](./spec.md) Key Entities｜研究依據：[research.md](./research.md)

## 實體

### U9SessionContext（每根 1m K 的情境標註）

| 欄位 | 型別 | 規則/驗證 |
|------|------|-----------|
| ts | int64 epoch 秒 | 唯一鍵；1m 對齊 |
| session | enum {DAY, NIGHT, CLOSED} | DAY=[08:45,13:44]；NIGHT=[15:00,23:59]∪[00:00,04:59]；其餘 CLOSED |
| trading_day | date | FR-003 規則歸屬；SC-004 抽驗（週五夜盤→週一） |
| in_killzone_tw | bool | [08:45,09:30]（KILLZONE_TW） |
| in_killzone_us | bool | [21:30,23:30]（KILLZONE_US） |
| or_high / or_low | float \| None | 交易日日盤開盤區間（前 OPENING_RANGE_MIN 根）；日盤首 15 根內即時形成 |
| tday_open | float \| None | 該交易日首根 open |
| bias | enum {+1, 0, -1} | close vs tday_open 與已形成區間中線（+1：雙雙之上；-1：雙雙之下；否則 0） |

**狀態轉移**：無（純標註）。**不變式**：週六/週日棒 → 不得產生 trading_day 標註為週六日（SC-004 相關斷言）。

### FvgEvent（U10 缺口事件）

| 欄位 | 型別 | 規則/驗證 |
|------|------|-----------|
| formed_ts | int64 | 第三根（當根）時間 |
| direction | enum {BULL, BEAR} | BULL：low[i] − high[i−2] > gap_min；BEAR：low[i−2] − high[i] > gap_min |
| top / bottom | float | BULL：top=low[i]、bottom=high[i−2]；BEAR 對稱 |
| width_pts | float | top − bottom；須 > FVG_MIN_POINTS 與 ATR×FVG_ATR_MULT 之較大者 |
| filled_ts | int64 \| None | 其後任一根價格觸及缺口另一緣（BULL：low ≤ bottom；BEAR：high ≥ top）即回補；全史掃描時可跨根追蹤 |

**狀態轉移**：formed → (filled | unfilled-at-window-end)。

### SweepReclaimEvent（U11 掃蕩收回事件）

| 欄位 | 型別 | 規則/驗證 |
|------|------|-----------|
| ts | int64 | 觸發根 |
| direction | enum {SWEEP_LOW_RECLAIM, SWEEP_HIGH_RECLAIM} | 前 LIQ_POOL_LOOKBACK 根（不含當根）之池值：low < pool_low 且 close > pool_low；或 high > pool_high 且 close < pool_high |
| pool_ref | float | 觸發時的池值 |

**不變式**：池值以 shift(1) 起算——當根不計入自己的池（無前視）。

### CostModel（FR-004）

| 常數 | 值 | 來源 |
|------|-----|------|
| POINT_VALUE_TWD | 10 | 微型台指契約乘數 |
| FEE_RT_TWD_RANGE | (30, 50) | 券商手續費區間 |
| TAX_RATE | 0.00002 | 期交稅率（期交所規定） |
| cost_pts(price) | (fee + tax×price×10)/10 | 上下限雙值呈列 |

### ValidationReport（FR-008 / SC-002）

| 欄位 | 規則 |
|------|------|
| detector / params_snapshot | U9-U11 與常數快照（FR-011 參數敏感性揭露） |
| horizon | ∈ {15, 30, 60}（research D5） |
| events / win_rate / raw_pts_mean | 事件研究統計 |
| net_pts_mean_low / net_pts_mean_high | 扣成本區間下限/上限後淨期望（雙值） |
| cum_expectancy_maxdd | 累積期望曲線最大回撤（點） |
| killzone_split | 事件在 killzone 內/外對照（U9×U10/U11 交互） |
| verdict | 可交易（淨期望下限 > 0）／不可交易（上限 < 0）／待調參（跨零） |

## 關係

- U9 為切層維度（每根 K），U10/U11 為事件（稀疏）——事件透過 ts 對應到 U9 的 in_killzone/session/trading_day。
- 全部實體僅由 1m OHLCV 序列構成（FR-002）；不引入任何外部指標依賴。
