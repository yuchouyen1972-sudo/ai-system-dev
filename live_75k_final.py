import sys, pickle, numpy as np, pandas as pd
sys.path.insert(0, r"C:\Users\yp86n\OneDrive\文件\Default Project")
import tmf_engine as engine

INITIAL_EQUITY = 75000
MAX_LOTS = 1
MAX_RISK_PCT = 0.01
STOP_MIN = 40; STOP_MAX = 80
TAKE_MULT = 1.8
DYNAMIC_THRESHOLD = 75000 * 1.5
PROB_THR = 0.62
VAR_THR = 0.01
CONSEC_MIN = 3  # 碎盤門控：進場當根連續同向根數不足則觀望（2021-24選、2025-26驗證有效）

with open(r"C:\Users\yp86n\OneDrive\文件\Default Project\xgb_improved.pkl", "rb") as f:
    xgb_model = pickle.load(f)
feature_cols = xgb_model["features"]

def build_features_15m(df):
    c, h, l, v, o = df['close'], df['high'], df['low'], df['volume'], df['open']
    f = pd.DataFrame(index=df.index)
    f['ret1']=c.pct_change(1); f['ret5']=c.pct_change(5); f['ret20']=c.pct_change(20)
    f['ret60']=c.pct_change(60)
    f['ma20']=c.rolling(20).mean(); f['ma60']=c.rolling(60).mean(); f['ma120']=c.rolling(120).mean()
    f['bias20']=c/f['ma20']-1; f['bias60']=c/f['ma60']-1; f['bias120']=c/f['ma120']-1
    f['trend_strength']=(f['ma20']-f['ma60'])/f['ma60']
    f['trend_strength_60_120']=(f['ma60']-f['ma120'])/f['ma120']
    pc=c.shift(1)
    tr=pd.concat([h-l,(h-pc).abs(),(l-pc).abs()],axis=1).max(axis=1)
    f['atr14']=tr.rolling(14).mean(); f['atr_ratio']=f['atr14']/c
    f['atr21']=tr.rolling(21).mean(); f['atr21_ratio']=f['atr21']/c
    delta=c.diff()
    gain=delta.where(delta>0,0).rolling(14).mean()
    loss=-delta.where(delta<0,0).rolling(14).mean()
    rs=gain/loss.replace(0,np.nan)
    f['rsi']=(100-(100/(1+rs))).fillna(50); f['rsi_5']=f['rsi'].rolling(5).mean()
    f['rsi_10']=f['rsi'].rolling(10).mean(); f['rsi_ma_diff']=f['rsi']-f['rsi_5']
    f['std20']=c.rolling(20).std(); f['std20_ratio']=f['std20']/c
    f['std21']=c.rolling(21).std(); f['vol_regime']=f['std20']/f['std21']
    f['bb_upper']=f['ma20']+2*f['std20']; f['bb_lower']=f['ma20']-2*f['std20']
    f['bb_pos']=(c-f['bb_lower'])/(f['bb_upper']-f['bb_lower']).replace(0,np.nan)
    f['bb_width']=(f['bb_upper']-f['bb_lower'])/c
    vma20=v.rolling(20).mean()
    f['vol_ratio']=v/vma20.replace(0,np.nan); f['vol_ratio']=f['vol_ratio'].fillna(1).clip(0,5)
    f['vol_accel']=f['vol_ratio']/f['vol_ratio'].shift(5).replace(0,np.nan)
    f['vol_accel']=f['vol_accel'].fillna(1).clip(0,3)
    f['hour']=df.index.hour; f['dow']=df.index.dayofweek; f['minute']=df.index.minute
    f['intraday_ret']=c/o-1; f['capo']=(c-o)/tr.replace(0,np.nan)
    f['body']=abs(c-o)/tr.replace(0,np.nan)
    f['upper_shadow']=(h-np.maximum(c,o))/tr.replace(0,np.nan)
    f['lower_shadow']=(np.minimum(c,o)-l)/tr.replace(0,np.nan)
    f['consec']=(c>o).astype(int).groupby((c<=o).cumsum()).cumcount()
    f['consec']=np.where(c>o,f['consec']+1,-f['consec']-1)
    f['hi_lo_ratio']=(h-l)/c
    f['mom_diverge']=f['rsi']-f['ret20']*100
    f['mom_diverge5']=f['rsi_5']-f['ret5']*100
    roll_max=c.rolling(20).max(); roll_min=c.rolling(20).min()
    f['pct_rank']=(c-roll_min)/(roll_max-roll_min+1e-8)
    f['vwap_proxy']=((c*h)+(c*l))/2
    f['vwap_dev']=(c-f['vwap_proxy'])/f['atr14'].replace(0,np.nan)
    from scipy.stats import skew, kurtosis as scipy_kurt
    f['price_skew']=c.rolling(20).apply(lambda x: skew(x) if len(x)>2 else 0, raw=False)
    f['price_kurt']=c.rolling(20).apply(lambda x: scipy_kurt(x, fisher=True) if len(x)>2 else 0, raw=False)
    f['range_20']=(h-l).rolling(20).mean()/c
    f['range_5']=(h-l).rolling(5).mean()/c
    return f

class LiveTMFEngine:
    def __init__(self, equity=INITIAL_EQUITY):
        self.equity = equity; self.peak = equity
        self.daily_loss_pts = 0; self.weekly_loss_pts = 0
        self.open_position = False; self.position_side = 0
        self.position_entry = 0.0; self.position_stop = 0.0; self.position_target = 0.0; self.position_qty = 0
    def risk_check(self, stop_points):
        risk_money = self.equity * MAX_RISK_PCT
        loss_per_lot = stop_points * 10
        lots = int(risk_money // loss_per_lot)
        lots = min(lots, MAX_LOTS)
        if self.equity < 35050 * 1.5: return 0
        if self.daily_loss_pts <= -150: return 0
        if self.weekly_loss_pts <= -300: return 0
        dd = (self.peak - self.equity) / self.peak
        if dd >= 0.10: return 0
        return lots
    def decide(self, df_15m, current_price):
        feats = build_features_15m(df_15m)
        X = feats[feature_cols].fillna(0).values[-1:].reshape(1, -1)
        probs = np.array([m.predict_proba(X)[:, 1] for m in xgb_model["models"]])
        mean_prob = float(probs.mean()); var = float(probs.var())
        cl = df_15m['close'].values; op = df_15m['open'].values
        _s = 1 if cl[-1] > op[-1] else (-1 if cl[-1] < op[-1] else 0)
        _cc = 0
        for k in range(len(cl) - 1, -1, -1):
            _sk = 1 if cl[k] > op[k] else (-1 if cl[k] < op[k] else 0)
            if _sk == _s: _cc += 1
            else: break
        if _cc < CONSEC_MIN:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var,
                    "stop": 0, "take": 0, "reason": f"碎盤門控 consec={_cc}<{CONSEC_MIN}",
                    "optimized": self.equity >= DYNAMIC_THRESHOLD}
        atr_val = df_15m['close'].rolling(14).std().iloc[-1] if len(df_15m) >= 14 else 50
        stop_points = int(np.clip(round(atr_val * 1.5), STOP_MIN, STOP_MAX))
        take_points = int(stop_points * TAKE_MULT)
        lots = self.risk_check(stop_points)
        is_optimized = self.equity >= DYNAMIC_THRESHOLD
        if lots == 0:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": "風控拒絕", "optimized": is_optimized}
        if var > VAR_THR:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": "MC Dropout 不確定性過高", "optimized": is_optimized}
        if mean_prob < PROB_THR:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": f"信心不足 {mean_prob:.2f}<{PROB_THR}", "optimized": is_optimized}
        cost_points = 8
        expected = mean_prob * take_points - (1 - mean_prob) * stop_points - cost_points
        if expected < cost_points:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": f"期望值不足 EV={expected:.1f}", "optimized": is_optimized}
        direction = 1 if mean_prob > 0.5 else -1
        if self.open_position:
            return {"action": "HOLD", "lots": 0, "direction": 0, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": "已持倉 1 口無加碼", "optimized": is_optimized}
        return {"action": "BUY" if direction == 1 else "SELL", "lots": lots, "direction": direction, "prob": mean_prob, "var": var, "stop": stop_points, "take": take_points, "reason": f"XGB OK EV={expected:.1f}點", "optimized": is_optimized}
    def execute(self, action_dict, current_price):
        if action_dict["action"] == "HOLD":
            return {"order": "NONE", "status": "OBSERVE", "reason": action_dict["reason"]}
        qty = action_dict["lots"]; direction = action_dict["direction"]
        entry = current_price
        stop = entry - direction * action_dict["stop"]
        target = entry + direction * action_dict["take"]
        self.open_position = True; self.position_side = direction
        self.position_entry = entry; self.position_stop = stop; self.position_target = target; self.position_qty = qty
        return {"order": "BUY" if direction == 1 else "SELL", "qty": qty, "entry": entry, "stop": stop, "target": target, "reason": action_dict["reason"], "optimized": action_dict.get("optimized", False)}
    def rearm(self):
        """熔斷重啟：日/週虧損歸零、峰值重設為當前權益。
        必須人工確認盤面恢復後才能呼叫，禁止自動呼叫。"""
        self.daily_loss_pts = 0; self.weekly_loss_pts = 0; self.peak = self.equity
    def on_tick(self, price):
        if not self.open_position:
            return None
        if self.position_side > 0:
            if price <= self.position_stop:
                return {"action": "FLATTEN", "reason": "停損", "pnl_points": int((self.position_stop - self.position_entry)//10)}
            if price >= self.position_target:
                return {"action": "TAKE_PROFIT", "reason": "停利", "pnl_points": int((self.position_target - self.position_entry)//10)}
        else:
            if price >= self.position_stop:
                return {"action": "FLATTEN", "reason": "停損", "pnl_points": int((self.position_stop - self.position_entry)//10)}
            if price <= self.position_target:
                return {"action": "TAKE_PROFIT", "reason": "停利", "pnl_points": int((self.position_entry - self.position_target)//10)}
        return {"action": "HOLD", "monitor": True, "equity": self.equity}

if __name__ == "__main__":
    feature_cols = xgb_model["features"]
    engine_live = LiveTMFEngine()
    df_15m = pd.DataFrame({
        "open": [23000, 23050, 23100, 23080, 23100, 23150, 23100, 23050],
        "close": [23050, 23100, 23080, 23150, 23200, 23100, 23050, 23120],
        "high": [23050, 23100, 23150, 23120, 23200, 23250, 23150, 23120],
        "low": [22950, 23000, 23050, 23000, 23100, 23150, 23000, 23100],
        "volume": [3000, 2800, 3200, 2500, 2900, 3100, 2800, 2600]
    }, index=pd.date_range("2026-09-01 08:45", periods=8, freq="15min"))

    decision = engine_live.decide(df_15m.iloc[-4:], df_15m['close'].iloc[-1])
    print("決策:", decision)
    order = engine_live.execute(decision, df_15m['close'].iloc[-1])
    print("下單:", order)
    print(f"\n75k 實盤引擎已啟動 (XGBoost 7 模型 DART MC Dropout)。參數: 初始:{INITIAL_EQUITY} 口數:{MAX_LOTS} 風險:{MAX_RISK_PCT*100:.0f}% 停損:{STOP_MIN}-{STOP_MAX} 風險報酬:1:{TAKE_MULT} 閾值:{PROB_THR} MC_Var_THR:{VAR_THR} 動態門檻:{DYNAMIC_THRESHOLD}")
