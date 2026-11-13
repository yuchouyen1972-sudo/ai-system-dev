import sys
sys.path.insert(0, r"C:\Users\yp86n\OneDrive\文件\Default Project")
import importlib.util
spec = importlib.util.spec_from_file_location("mybacktest", r"C:\Users\yp86n\OneDrive\文件\Default Project\backtest.py")
btmod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(btmod)
import tmf_engine as engine
import pandas as pd, numpy as np, pickle

# 1. 載入完整 15m
df = pd.read_parquet(r"C:\Users\yp86n\OneDrive\文件\Default Project\mtxf-ai-trader\data_store\TXF_15m_2006_20260901.parquet")
df = df.sort_index().between_time("08:45","13:45")
print(f"15m {len(df)} 根 {df.index.min()} -> {df.index.max()}")

# 2. 預先建特徵（一次）
def build_features(d):
    c,h,l,v = d['close'],d['high'],d['low'],d['volume']
    ret1=c.pct_change(1); ret5=c.pct_change(5); ret20=c.pct_change(20)
    ma20=c.rolling(20).mean(); ma60=c.rolling(60).mean(); ma120=c.rolling(120).mean()
    bias20=c/ma20-1; bias60=c/ma60-1; bias120=c/ma120-1
    pc=c.shift(1); tr=pd.concat([h-l,(h-pc).abs(),(l-pc).abs()],axis=1).max(axis=1)
    atr14=tr.rolling(14).mean(); atr_ratio=atr14/c
    delta=c.diff(); gain=delta.where(delta>0,0).rolling(14).mean(); loss=-delta.where(delta<0,0).rolling(14).mean()
    rs=gain/loss.replace(0,np.nan); rsi=(100-(100/(1+rs))).fillna(50); rsi_5=rsi.rolling(5).mean()
    std20=c.rolling(20).std(); bb_up=ma20+2*std20; bb_low=ma20-2*std20
    bb_pos=(c-bb_low)/(bb_up-bb_low).replace(0,np.nan)
    vol_ratio=v/v.rolling(20).mean().replace(0,np.nan); vol_ratio=vol_ratio.fillna(1).clip(0,5)
    return pd.DataFrame({
        'ret1':ret1,'ret5':ret5,'ret20':ret20,
        'bias20':bias20,'bias60':bias60,'bias120':bias120,
        'atr_ratio':atr_ratio,'rsi':rsi,'rsi_5':rsi_5,
        'bb_pos':bb_pos,'vol_ratio':vol_ratio,
        'hour':d.index.hour,'dow':d.index.dayofweek,'std20':std20
    })

feat_df = build_features(df)
feature_cols = ['ret1','ret5','ret20','bias20','bias60','bias120','atr_ratio','rsi','rsi_5','bb_pos','vol_ratio','hour','dow','std20']
feat_df = feat_df.fillna(0)

# 3. 預先載入 XGBoost 模型並預測全部
with open(r"C:\Users\yp86n\OneDrive\文件\Default Project\xgb_15m_full.pkl", "rb") as f:
    model_data = pickle.load(f)
models = model_data["models"]
print("預先計算 5 模型預測...")
# 分批預測以控制記憶體 (108k 根 x 18 特徵)
batch_size = 5000
mean_probs = np.zeros(len(feat_df))
var_probs = np.zeros(len(feat_df))
for start in range(0, len(feat_df), batch_size):
    end = min(start+batch_size, len(feat_df))
    X = feat_df.iloc[start:end][feature_cols].values
    probs = np.array([m.predict_proba(X)[:,1] for m in models])
    mean_probs[start:end] = probs.mean(axis=0)
    var_probs[start:end] = probs.var(axis=0)
    if start % 20000 == 0:
        print(f"  已預測 {end} / {len(feat_df)}")

df['mean_prob'] = mean_probs
df['mean_var'] = var_probs
print("預測完成，已附加到 df")

# 4. 快速 Walk-Forward (用預計算的 mean_prob/var，不再建特徵)
bt = btmod.WalkForwardBacktest(initial_equity=75000)

def wf(df, tb=20000, te=5000, st=5000):
    s=[]; start=0; wid=0
    while start+tb+te <= len(df):
        s.append((wid, df.iloc[start:start+tb], df.iloc[start+tb:start+tb+te]))
        start+=st; wid+=1
    return s

splits = wf(df)
print(f"Walk-Forward {len(splits)} 窗口")

Trade = btmod.Trade
equity=75000; peak=75000; max_dd=0; trades=[]; equity_curve=[]

for wid, train, test in splits:
    acct = engine.Account(equity=equity, peak_equity=peak)
    # 快速決策: 直接用預計算的 prob/var，不再建特徵
    # 但要模擬原始邏輯: 每 20 根決策一次
    for i in range(60, len(test), 20):
        # 取當前點的預測值
        idx_global = len(train) + i  # 在完整 df 中的索引
        # 但我使用的是 sub_df 切分，需要映射回完整 df
        # 簡化: 直接取 sub_df 的 index 對應完整 df
        ts_idx = test.index[i]
        # 在完整 df 找對應位置
        pos_in_df = df.index.get_loc(ts_idx)
        mean_prob = float(df['mean_prob'].iloc[pos_in_df])
        var = float(df['mean_var'].iloc[pos_in_df])
        
        # 風控
        entr_price = test.iloc[i]['close']
        # 停損用 ATR（簡化）
        atr_val = df['close'].rolling(14).std().iloc[pos_in_df] if pos_in_df >= 14 else 50
        stop_points = int(np.clip(round(atr_val*1.5), 40, 80))
        take_points = int(stop_points*1.8)
        lots = 1 if (equity >= 35050*1.5 and eq >= 75000 and equity >= 75000*0.9) else int(equity*0.01//(stop_points*10))
        lots = max(0, min(lots, 1))
        
        # 閾值
        thr = 0.58
        if lots == 0 or var > 0.015 or mean_prob < thr:
            # HOLD
            continue
        # 期望值
        cost_points = 8
        expected = mean_prob*take_points - (1-mean_prob)*stop_points - cost_points
        if expected < cost_points: continue
        
        direction = 1 if mean_prob > 0.5 else -1
        entry_price = entr_price
        stop_price = entry_price - direction * stop_points
        take_price = entry_price + direction * take_points
        
        # 模擬退出 (簡化為 20 根後或碰到停損/停利)
        exit_idx = min(i + 20, len(test)-1)
        exit_price = test.iloc[exit_idx]['close']
        # 檢查停損/停利在期間
        for j in range(i+1, min(i+21, len(test))):
            bar = test.iloc[j]
            if direction == 1:
                if bar['low'] <= stop_price: exit_price = stop_price; exit_idx = j; break
                if bar['high'] >= take_price: exit_price = take_price; exit_idx = j; break
            else:
                if bar['high'] >= stop_price: exit_price = stop_price; exit_idx = j; break
                if bar['low'] <= take_price: exit_price = take_price; exit_idx = j; break
        
        pnl_points = (exit_price - entry_price) * direction
        pnl_money = pnl_points * 10 - (35*2 + entry_price*10*0.00002*2)
        equity += pnl_money
        peak = max(peak, equity)
        max_dd = max(max_dd, peak - equity)
        
        # 更新每週/日熔斷 (簡化)
        # 不精確但可用
        
        trades.append(Trade(test.index[i], test.index[exit_idx], direction, entry_price, exit_price, stop_points, take_points, int(pnl_points), pnl_money, mean_prob, var, f"XGB precomputed EF={expected:.1f}", wid))
        equity_curve.append((test.index[exit_idx], equity))
        
        if equity < 75000*0.9:  # 快速熔斷檢查
            break
    
    if wid % 3 == 0:
        print(f"Window {wid} done, equity={equity:.0f}, trades={len(trades)}")

# 結果
if trades:
    tdf = pd.DataFrame([t.__dict__ for t in trades])
    wins = tdf[tdf['pnl_money']>0]; losses = tdf[tdf['pnl_money']<=0]
    pf = wins['pnl_money'].sum()/abs(losses['pnl_money'].sum()) if len(losses) and losses['pnl_money'].sum()!=0 else float('inf')
    ret = (equity - 75000) / 75000
    sharpe = tdf['pnl_money'].mean() / tdf['pnl_money'].std() * np.sqrt(252*6) if tdf['pnl_money'].std()>0 else 0
    print(f"\n=== XGBacktest 結果 ===")
    print(f"交易 {len(tdf)}  PF {pf:.2f}  回報 {ret*100:.1f}%  回撤 {max_dd/peak*100:.1f}%  Sharpe {sharpe:.2f}")
    print(f"結束權益: {equity:.0f}")
    tdf.to_csv("backtest_trades.csv", index=False, encoding="utf-8-sig")
    pd.DataFrame(equity_curve, columns=["datetime","equity"]).to_csv("equity_curve.csv", index=False, encoding="utf-8-sig")
    print("已輸出 backtest_trades.csv, equity_curve.csv")
else:
    print("無交易產生")
