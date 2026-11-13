"""
微型台指 TMF 決策引擎骨架 - NT$75,000 專用
規格: 乘數 10元/點, 1跳1點=10元, 原始保證金35050, 維持26900 (依期交所公告為準，2026/08/12更新)
原則: 單筆風險1% + MC Dropout 不確定性過濾 + 日/週熔斷
"""

from dataclasses import dataclass

import numpy as np
import pandas as pd


# ============== 1. 合約與帳戶設定 ==============
@dataclass
class TMFSpec:
    multiplier: int = 10  # 10元/點
    tick_point: int = 1  # 最小升降1點
    tick_value: int = 10  # 1點=10元
    original_margin: int = 35050  # 原始保證金 (依期交所公告 2026/08/12)
    maintenance_margin: int = 26900  # 維持保證金 (依期交所公告 2026/08/12)
    tax_rate: float = 0.00002  # 期交稅 2/100000
    commission_per_side: int = 35  # 依期貨商, 單邊35元為例
    min_stop_points: int = 40
    max_stop_points: int = 80


@dataclass
class Account:
    equity: float = 75000  # 權益數
    available: float = 75000
    max_risk_pct: float = 0.01  # 單筆1% =750元
    daily_loss_limit_points: int = 150  # 日熔斷 -150點 = -1500元
    weekly_loss_limit_points: int = 300
    daily_loss_points: int = 0
    weekly_loss_points: int = 0
    max_drawdown_pct: float = 0.10
    peak_equity: float = 75000

    def reset_daily(self):
        """換日重置日虧損點數"""
        self.daily_loss_points = 0

    def reset_weekly(self):
        """換週重置週虧損點數"""
        self.weekly_loss_points = 0

    def rebaseline(self):
        """熔斷重啟：以當前權益為新高點。必須人工確認盤面恢復後才能呼叫，
        禁止在程式內自動呼叫（否則熔斷形同虛設）。"""
        self.peak_equity = self.equity


SPEC = TMFSpec()


def validate_config(spec: TMFSpec, account: Account) -> dict:
    """驗證 TMFSpec 與 Account 參數合理性，回傳結構化報告。

    回傳 dict 結構:
        {
            "valid": bool,          # 整體是否通過
            "warnings": list[str],  # 警告（不阻斷但需注意）
            "errors": list[str],    # 錯誤（阻斷，valid=False）
            "checks": list[dict],   # 每個檢查項目 {name, passed, detail}
        }
    """
    warnings: list = []
    errors: list = []
    checks: list = []

    def _check(name: str, passed: bool, detail: str):
        checks.append({"name": name, "passed": bool(passed), "detail": detail})
        if not passed:
            errors.append(f"{name}: {detail}")

    def _warn(name: str, detail: str):
        warnings.append(f"{name}: {detail}")
        checks.append({"name": name, "passed": True, "detail": detail})

    # FR-006: tick_value 必須大於 0，否則 position_size 會除以零
    _check("tick_value_nonzero", spec.tick_value > 0, f"tick_value={spec.tick_value}")

    # FR-005: min_stop_points 不得大於 max_stop_points
    _check(
        "stop_range_valid",
        spec.min_stop_points <= spec.max_stop_points,
        f"min={spec.min_stop_points} max={spec.max_stop_points}",
    )

    # FR-003: 權益必須 >= original_margin * 1.5
    required_margin = spec.original_margin * 1.5
    _check(
        "margin_sufficient",
        account.equity >= required_margin,
        f"equity={account.equity} required={required_margin}",
    )

    # FR-004: max_risk_pct 超過 5% 視為過高風險（警告層級）
    if account.max_risk_pct > 0.05:
        _warn("risk_ratio_high", f"max_risk_pct={account.max_risk_pct} > 5%")
    else:
        checks.append(
            {
                "name": "risk_ratio_ok",
                "passed": True,
                "detail": f"max_risk_pct={account.max_risk_pct}",
            }
        )

    # peak_equity 不應小於 equity（異常狀態）
    if account.peak_equity < account.equity:
        _warn(
            "peak_equity_anomaly",
            f"peak_equity={account.peak_equity} < equity={account.equity}",
        )

    # daily_loss_limit_points 與 weekly_loss_limit_points 應為正數
    _check(
        "daily_loss_limit_positive",
        account.daily_loss_limit_points > 0,
        f"daily_loss_limit_points={account.daily_loss_limit_points}",
    )
    _check(
        "weekly_loss_limit_positive",
        account.weekly_loss_limit_points > 0,
        f"weekly_loss_limit_points={account.weekly_loss_limit_points}",
    )

    valid = len(errors) == 0

    return {
        "valid": valid,
        "warnings": warnings,
        "errors": errors,
        "checks": checks,
    }


# ============== 2. 風控計算器 ==============
class RiskManager:
    def __init__(self, spec: TMFSpec, account: Account):
        self.spec = spec
        self.account = account

    def position_size(self, stop_points: int) -> int:
        """依1%風險計算可做口數, 75k小資金只允許0或1口"""
        if (
            stop_points < self.spec.min_stop_points
            or stop_points > self.spec.max_stop_points
        ):
            return 0
        risk_amount = self.account.equity * self.account.max_risk_pct  # 750
        loss_per_lot = stop_points * self.spec.tick_value  # e.g. 50*10=500
        lots = int(risk_amount // loss_per_lot)
        return min(lots, 1)  # 75k硬上限1口

    def margin_check(self, lots: int = 1) -> bool:
        """檢查保證金是否足夠, 需保留維持率>150%"""
        required = self.spec.original_margin * lots
        # 建議保留至少 1.5倍原始保證金
        return self.account.equity >= required * 1.5

    def daily_circuit_breaker(self) -> bool:
        return self.account.daily_loss_points <= -self.account.daily_loss_limit_points

    def weekly_circuit_breaker(self) -> bool:
        return self.account.weekly_loss_points <= -self.account.weekly_loss_limit_points

    def drawdown_check(self) -> bool:
        dd = (self.account.peak_equity - self.account.equity) / self.account.peak_equity
        return dd >= self.account.max_drawdown_pct

    def calc_stop_by_atr(self, atr: float, atr_mult: float = 1.5) -> int:
        points = int(round(atr * atr_mult))
        return int(
            np.clip(points, self.spec.min_stop_points, self.spec.max_stop_points)
        )

    def update_pnl(self, pnl_points: int):
        pnl_money = pnl_points * self.spec.tick_value
        # 扣除成本
        cost = self.spec.commission_per_side * 2 + abs(pnl_points * 0.0)  # 稅另計
        # 期交稅依契約價值: 23000*10*0.00002=4.6/邊
        tax = 5 * 2  # 簡化
        net = pnl_money - cost - tax if pnl_points != 0 else 0
        self.account.equity += net
        self.account.daily_loss_points += pnl_points
        self.account.weekly_loss_points += pnl_points
        self.account.peak_equity = max(self.account.peak_equity, self.account.equity)
        return net


# ============== 3. 訊號層 - MC Dropout 集成 (示意) ==============
class MCDropoutEnsemble:
    """
    用MC Dropout估計不確定性. 實務需在PyTorch模型推理時開啟dropout
    此處以模擬演示: 採樣30次, 取mean/var
    """

    def __init__(
        self,
        n_samples: int = 30,
        prob_threshold: float = 0.58,
        var_threshold: float = 0.015,
    ):
        self.n_samples = n_samples
        self.prob_threshold = prob_threshold
        self.var_threshold = var_threshold

    def predict_with_uncertainty(
        self, model_fn, features: np.ndarray
    ) -> tuple[float, float]:
        # model_fn: 單次forward回傳機率 0~1
        samples = [model_fn(features) for _ in range(self.n_samples)]
        mean_prob = float(np.mean(samples))
        var = float(np.var(samples))
        return mean_prob, var

    def should_trade(self, mean_prob: float, var: float) -> bool:
        if var > self.var_threshold:
            return False
        if mean_prob < self.prob_threshold and mean_prob > (1 - self.prob_threshold):
            return False  # 信心不足
        return True


# 模擬模型 (實際替換為你的LSTM/XGB)
def dummy_model_fn(features: np.ndarray) -> float:
    # 假設最後一維是動能因子
    score = np.tanh(features[-1]) * 0.3 + 0.5 + np.random.normal(0, 0.07)
    return float(np.clip(score, 0.05, 0.95))


# ============== 4. 決策引擎本體 ==============
@dataclass
class Signal:
    direction: int  # 1多 -1空 0觀望
    prob: float
    uncertainty: float
    stop_points: int
    take_points: int
    reason: str


class TMFDecisionEngine:
    def __init__(self, account: Account = None):
        self.spec = SPEC
        self.account = account or Account()
        self.config_validation = validate_config(self.spec, self.account)
        self.risk = RiskManager(self.spec, self.account)
        self.ensemble = MCDropoutEnsemble()

    def build_features(self, df: pd.DataFrame) -> np.ndarray:
        """範例特徵: 需替換為你的特徵工廠"""
        close = df["close"]
        ema20 = close.ewm(span=20).mean().iloc[-1]
        ema60 = close.ewm(span=60).mean().iloc[-1]
        # ATR
        high_low = df["high"] - df["low"]
        atr = high_low.rolling(14).mean().iloc[-1]
        momentum = (close.iloc[-1] - close.iloc[-1 - 10]) / close.iloc[-1 - 10]
        return np.array([ema20, ema60, atr, momentum])

    def decide(self, df: pd.DataFrame) -> Signal:
        # 0. 熔斷檢查
        if self.risk.daily_circuit_breaker():
            return Signal(0, 0.5, 0, 0, 0, "日熔斷 -150點, 今日停止交易")
        if self.risk.weekly_circuit_breaker():
            return Signal(0, 0.5, 0, 0, 0, "週熔斷 -300點, 本週停止")
        if self.risk.drawdown_check():
            return Signal(0, 0.5, 0, 0, 0, "最大回撤10% 暫停實盤")
        if not self.risk.margin_check(1):
            return Signal(0, 0.5, 0, 0, 0, f"權益{self.account.equity}不足1.5倍保證金")

        # 1. 特徵與ATR停損
        feats = self.build_features(df)
        atr = feats[2] if not np.isnan(feats[2]) else 50
        stop_points = self.risk.calc_stop_by_atr(atr, 1.5)
        take_points = int(stop_points * 1.8)  # RR 1:1.8

        # 2. MC Dropout 不確定性
        mean_prob, var = self.ensemble.predict_with_uncertainty(dummy_model_fn, feats)

        # 3. 風險報酬與口數檢查
        lots = self.risk.position_size(stop_points)
        if lots == 0:
            return Signal(
                0,
                mean_prob,
                var,
                stop_points,
                take_points,
                f"停損{stop_points}點超出40-80限制或風險不足",
            )

        if not self.ensemble.should_trade(mean_prob, var):
            return Signal(
                0,
                mean_prob,
                var,
                stop_points,
                take_points,
                f"信心不足 prob={mean_prob:.2f} var={var:.4f}",
            )

        # 4. 方向判定 + 期望值檢查（方向感：空方用反面機率，否則空單永遠被擋）
        # 成本約 80元/趟 = 8點, 需期望值>2倍成本
        cost_points = 8
        direction = 1 if mean_prob > 0.5 else -1
        p_side = mean_prob if direction == 1 else 1 - mean_prob
        expected_points = (
            p_side * take_points - (1 - p_side) * stop_points - cost_points
        )
        if expected_points < cost_points:
            return Signal(
                0,
                mean_prob,
                var,
                stop_points,
                take_points,
                f"期望值{expected_points:.1f}點不足",
            )

        return Signal(
            direction,
            mean_prob,
            var,
            stop_points,
            take_points,
            f"OK lots={lots} EV={expected_points:.1f}點",
        )

    def execute_order(self, signal: Signal, price: float) -> dict:
        if signal.direction == 0:
            return {"action": "HOLD", "reason": signal.reason}
        lots = self.risk.position_size(signal.stop_points)
        stop_price = price - signal.direction * signal.stop_points
        take_price = price + signal.direction * signal.take_points
        return {
            "action": "BUY" if signal.direction == 1 else "SELL",
            "lots": lots,
            "price": price,
            "stop": stop_price,
            "take": take_price,
            "prob": signal.prob,
            "var": signal.uncertainty,
            "reason": signal.reason,
        }


# ============== 5. 回測與使用範例 ==============
if __name__ == "__main__":
    # 產生假K線測試
    np.random.seed(42)
    dates = pd.date_range("2024-01-01", periods=500, freq="5min")
    price = 23000 + np.cumsum(np.random.randn(500) * 8)
    df = pd.DataFrame(
        {
            "close": price,
            "high": price + np.random.rand(500) * 5,
            "low": price - np.random.rand(500) * 5,
        }
    )

    engine = TMFDecisionEngine(Account(equity=75000))
    print(
        f"初始權益: {engine.account.equity} 保證金檢查1口: {engine.risk.margin_check(1)}"
    )
    print(
        f"單筆風險1%={engine.account.equity * 0.01}元, 50點停損口數={engine.risk.position_size(50)}"
    )

    # 模擬連續5次決策
    for i in range(5):
        sub_df = df.iloc[i * 100 : (i + 1) * 100]
        sig = engine.decide(sub_df)
        order = engine.execute_order(sig, price=sub_df["close"].iloc[-1])
        print(
            f"\n[{i + 1}] Signal dir={sig.direction} prob={sig.prob:.3f} var={sig.uncertainty:.4f} stop={sig.stop_points} take={sig.take_points}"
        )
        print(f"    -> {order}")
        # 模擬損益
        if order["action"] != "HOLD":
            pnl = np.random.choice([sig.take_points, -sig.stop_points], p=[0.45, 0.55])
            net = engine.risk.update_pnl(pnl)
            print(
                f"    模擬結算 {pnl}點 淨損益 {net:.0f}元 權益 {engine.account.equity:.0f}"
            )

    print(
        f"\n最終權益 {engine.account.equity:.0f} 日虧{engine.account.daily_loss_points}點 週虧{engine.account.weekly_loss_points}點"
    )
