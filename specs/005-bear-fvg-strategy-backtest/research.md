# Research: 005-bear-fvg-strategy-backtest

**Date**: 2026-11-22｜**狀態**：設計決策全部定案（D1-D6），無 [NEEDS CLARIFICATION]

## D1 停損法：缺口上緣＋buffer（主）vs ATR 倍數（對照組）

- **Decision**: 主停損＝缺口上緣（top）＋ SL_BUFFER_ATR(0.25)×ATR14；下限 SL_MIN_PTS(10)。對照組＝entry＋SL_ATR_MULT(1.0)×ATR（兩組並列報告）。
- **Rationale**: 構想語義——空頭失衡（BEAR FVG）被回補即失效，停損掛在結構位有邏輯；純 ATR 倍數無結構意義，作對照驗證結構位的價值。
- **Alternatives rejected**: 純固定點數（無波動適應）；追蹤停損（v1 不做，控制複雜度，揭露）。

## D2 滑價方向：一律「往不利方向」

- **Decision**: 做空進場成交價＝open − slip（賣更低＝未來回補利潤少）；回補成交價＝觸發價 + slip（買更高＝虧更多）。停利回補同理＋slip。
- **Rationale**: 統一「不利方向」定義，避免多空滑價符號錯置（外部腳本曾無滑價模型）。

## D3 跨盤強平：13:40/04:55 收盤價平倉

- **Decision**: 持倉期間遇到 13:40（日盤尾）或 04:55（夜盤尾）棒 → 以該根收盤價平倉（reason=EOD_FLAT）；該根若同時觸 SL/TP → 先 SL/TP（結構優先於時間）。強平棒不存在（資料缺口）→ 以時點前最後一根收盤價平倉。
- **Rationale**: 休市段不可持倉（保證金/跳空風險不可控）；外部腳本同慣例且棒存在已實測（04:55/13:40 皆在資料中）。

## D4 統計口徑

- **Decision**: PF＝總獲利點數/|總虧損點數|；MDD 依逐筆淨點數累積權益曲線（close-of-trade equity）；平均持有根數＝exit_i − entry_i。
- **Rationale**: 業界標準＋與 004 事件研究可對照（同為 trade-level equity）。

## D5 對照歸因框架

- **Decision**: 報告並列「事件研究 dedup H=15（raw +27.0）」與策略期望值；差異歸因三項：停損截斷（右尾被砍）、強平提前、滑價成本——各項以「移除該規則的重跑差值」量化（停損→純遠景；強平→不強平）。
- **Rationale**: 防止「回測器寫壞了但看起來合理」——歸因可驗證才是誠實對照。

## D6 訊號與互斥的介面

- **Decision**: 回測器吃 detect_fvg 全事件序列（過濾 BEAR）＋行索引；互斥沿用 004 select_nonoverlap 語義但改「動態出場後才可再進場」（非固定 H）——在回測迴圈內直接實作（持倉中忽略新訊號）。
- **Rationale**: 固定 H 版本無法表達停損/強平的變長持有；迴圈內互斥是單口滾動的天然形式。

## Data Model（摘要）

- **StrategyTrade**: entry_ts/exit_ts、reason∈{SL,TP,HORIZON,EOD_FLAT}、entry_px/exit_px、raw_pts/cost_pts/net_pts、hold_bars
- **EquitySummary**: trades/win_rate/pf/expectancy/mdd/avg_hold/verdict/params_snapshot
- **SlippageMatrix**: slip∈{0,1,2} × cost∈{lo,hi} → expectancy 淨值（6 格）
- 常數：SL_BUFFER_ATR=0.25、SL_MIN_PTS=10、SL_ATR_MULT=1.0（對照組）、TP_R_MULT=1.5、H_MAX=60、SLIPPTS=(0,1,2)、強平時點 13:40/04:55（皆註來源）
