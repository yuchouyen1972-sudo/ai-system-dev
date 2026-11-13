#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""完整交易系統回測：LiveTMFEngine 決策邏輯（live_75k_final.py:84-107）+ 真實模型 + 真實15m。
與實盤一致：特徵函數直接 import live 版；prob>=0.62、var<=0.01、EV>=8、1%風控、
保證金1.5x、日-150/週-300熔斷、MD 10%、持有中不加碼、bar收盤進場、後續bar先碰停損再停利。
假設（實盤碼內未定義）：每日/週 loss點數按自然日/ISO週重置；無max_hold（抱到出場或資料尾）。
注意：實盤邏輯下 direction 恆為+1（prob>=0.62>0.5），即只做多，本回測忠實保留。
"""
import sys
sys.path.insert(0, r"C:\Users\yp86n\OneDrive\文件\Default Project")
import warnings
warnings.filterwarnings('ignore')
import pickle
import os
import pandas as pd, numpy as np

import live_75k_final as live

PARQ = r"C:\Users\yp86n\OneDrive\文件\Default Project\mtxf-ai-trader\data_store\TXF_15m_2006_20260901.parquet"
MODEL_PATH = os.getenv("SYS_MODEL")  # 指定則用該模型檔，否則用實盤引擎內建模型
OUT_T = os.getenv("SYS_OUT_T", r"C:\Users\yp86n\OneDrive\文件\Default Project\sys_trades.csv")
OUT_E = os.getenv("SYS_OUT_E", r"C:\Users\yp86n\OneDrive\文件\Default Project\sys_equity.csv")

PROB_THR, VAR_THR = 0.62, 0.01
STOP_MIN, STOP_MAX, TAKE_MULT, COST_PTS = 40, 80, 1.8, 8
WARM = 200

print("載入資料...", flush=True)
df = pd.read_parquet(PARQ).sort_index().between_time("08:45", "13:45")
print(f"{len(df)}根 {df.index[0]}->{df.index[-1]}", flush=True)

print("建特徵（與實盤同函數）...", flush=True)
feats = live.build_features_15m(df)
if MODEL_PATH:
    with open(MODEL_PATH, "rb") as f:
        _md = pickle.load(f)
    feature_cols, _models = _md["features"], _md["models"]
    print(f"使用外部模型 {MODEL_PATH} cfg={_md.get('cfg')}", flush=True)
else:
    feature_cols, _models = live.feature_cols, live.xgb_model["models"]
X = feats[feature_cols].fillna(0).values.astype(np.float32)

print(f"批次預測（{len(_models)}模型）...", flush=True)
P = np.array([m.predict_proba(X)[:, 1] for m in _models])
mean_p, var_p = P.mean(axis=0), P.var(axis=0)

close = df['close'].values; high = df['high'].values; low = df['low'].values
dates = df.index
close_std = pd.Series(close).rolling(14).std().values
_up = np.where(close > df['open'].values, 1, np.where(close < df['open'].values, -1, 0))
_s = pd.Series(_up)
cc_arr = (((_s.groupby((_s != _s.shift()).cumsum()).cumcount() + 1)
           * _s.replace(0, np.nan).fillna(1)).abs().astype(int).values)
CONSEC_MIN = 3  # 與實盤引擎同值（2021-24選、2025-26驗證）

equity, peak = 75000.0, 75000.0
max_dd = 0.0
daily_pts, weekly_pts = 0.0, 0.0
cur_day, cur_week = dates[WARM].date(), dates[WARM].isocalendar().week
trades, curve = [], []
cnt = {"gate": 0, "risk": 0, "var": 0, "prob": 0, "ev": 0, "signal": 0}

i = WARM
n = len(df)
open_pos = None
while i < n:
    d, w = dates[i].date(), dates[i].isocalendar().week
    if d != cur_day: daily_pts, cur_day = 0.0, d
    if w != cur_week: weekly_pts, cur_week = 0.0, w
    if open_pos is not None:
        e, stop_px, take_px, entry_i = open_pos
        exited = None
        for j in range(i, n):
            if low[j] <= stop_px: exited = (j, stop_px, -(e - stop_px)); break
            if high[j] >= take_px: exited = (j, take_px, take_px - e); break
        if exited is None:
            j, exit_px, pts = n - 1, close[-1], close[-1] - e
            reason = "資料尾強制平倉"
        else:
            j, exit_px, pts = exited
            reason = "停損" if exit_px == stop_px else "停利"
        cost = 35 * 2 + e * 10 * 0.00002 * 2
        money = pts * 10 - cost
        equity += money; peak = max(peak, equity); max_dd = max(max_dd, peak - equity)
        daily_pts += pts; weekly_pts += pts
        trades.append({"entry_time": dates[entry_i], "exit_time": dates[j], "direction": 1,
                       "entry": e, "exit": exit_px, "stop": int(e - stop_px), "take": int(take_px - e),
                       "pnl_pts": round(float(pts), 1), "pnl_money": round(float(money), 1),
                       "prob": round(float(mean_p[entry_i]), 4), "reason": reason})
        curve.append((dates[j], equity))
        open_pos = None
        i = j + 1
        continue
    # 無持倉：走實盤 decide 同順序過濾（碎盤門控先行，與引擎一致）
    if cc_arr[i] < CONSEC_MIN:
        cnt["gate"] += 1; i += 1; continue
    stop_points = int(np.clip(round(close_std[i] * 1.5), STOP_MIN, STOP_MAX))
    take_points = int(stop_points * TAKE_MULT)
    lots = int((equity * 0.01) // (stop_points * 10))
    lots = min(lots, 1)
    if equity < 35050 * 1.5 or daily_pts <= -150 or weekly_pts <= -300 or (peak - equity) / peak >= 0.10 or lots == 0:
        cnt["risk"] += 1; i += 1; continue
    if var_p[i] > VAR_THR: cnt["var"] += 1; i += 1; continue
    if mean_p[i] < PROB_THR: cnt["prob"] += 1; i += 1; continue
    ev = mean_p[i] * take_points - (1 - mean_p[i]) * stop_points - COST_PTS
    if ev < COST_PTS: cnt["ev"] += 1; i += 1; continue
    cnt["signal"] += 1
    e = close[i]
    open_pos = (e, e - stop_points, e + take_points, i)
    i += 1

tdf = pd.DataFrame(trades)
tdf.to_csv(OUT_T, index=False, encoding="utf-8-sig")
pd.DataFrame(curve, columns=["datetime", "equity"]).to_csv(OUT_E, index=False, encoding="utf-8-sig")
print(f"過濾統計 gate={cnt['gate']} risk={cnt['risk']} var={cnt['var']} prob={cnt['prob']} ev={cnt['ev']} 進場={cnt['signal']}", flush=True)

if len(tdf) == 0:
    print("無交易"); sys.exit()
wins = tdf[tdf.pnl_money > 0]; losses = tdf[tdf.pnl_money <= 0]
pf = wins.pnl_money.sum() / abs(losses.pnl_money.sum()) if len(losses) and losses.pnl_money.sum() != 0 else float('inf')
ret = (equity - 75000) / 75000
days = (dates[-1] - dates[0]).days or 1
ann = (1 + ret) ** (365 / days) - 1 if ret > -1 else -1
rets = tdf.pnl_money / 75000
sharpe = rets.mean() / rets.std() * np.sqrt(252 * 6) if rets.std() != 0 else 0
print("=" * 60, flush=True)
print(f"交易 {len(tdf)}筆 勝率 {(tdf.pnl_money > 0).mean() * 100:.1f}% PF {pf:.2f}", flush=True)
print(f"權益 75000->{equity:.0f} 淨利 {equity - 75000:.0f} ({ret * 100:.1f}%) 年化 {ann * 100:.1f}%", flush=True)
print(f"最大回撤 {max_dd:.0f}元({max_dd / peak * 100:.1f}%) Sharpe {sharpe:.2f}", flush=True)
print(f"平均獲利 {wins.pnl_money.mean():.0f} 平均虧損 {losses.pnl_money.mean():.0f} 期望值 {tdf.pnl_money.mean():.0f}元/筆", flush=True)
print(f"停損 {int((tdf.reason == '停損').sum())} 停利 {int((tdf.reason == '停利').sum())}", flush=True)
print("=" * 60, flush=True)
tdf["yr"] = pd.to_datetime(tdf["entry_time"]).dt.year
g = tdf.groupby("yr")["pnl_money"].agg(["count", "sum"])
g["wr"] = tdf.groupby("yr")["pnl_money"].apply(lambda x: (x > 0).mean() * 100).round(1)
print(g.round(0).to_string(), flush=True)
