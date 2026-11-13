#!/usr/bin/env python3
"""
實時監控報表引擎 - 直接連接 live_75k_final.py 的 LiveTMFEngine
功能:
  1. 讀取引擎當前狀態 (權益、持倉、風控)
  2. 計算未實現損益 (若有持倉)
  3. 產生即時報表 (Markdown)
  4. 可接入真實行情源 (預留接口)
"""
import sys, os
sys.path.insert(0, r"C:\Users\yp86n\OneDrive\文件\Default Project")
import importlib.util
spec = importlib.util.spec_from_file_location("live_engine", r"C:\Users\yp86n\OneDrive\文件\Default Project\live_75k_final.py")
engine_mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(engine_mod)

import pandas as pd, numpy as np

class LiveMonitorReport:
    def __init__(self, engine=None, initial_equity=75000):
        self.engine = engine or engine_mod.LiveTMFEngine(equity=initial_equity)
        self.initial_equity = initial_equity
        self.current_price = None
        self.report_file = "live_monitor_report.md"

    def set_current_price(self, price: float):
        self.current_price = price

    def get_engine_state(self):
        """讀取引擎即時狀態"""
        e = self.engine
        state = {
            "權益": e.equity,
            "峰值權益": e.peak,
            "初始權益": self.initial_equity,
            "報酬 %": ((e.equity - self.initial_equity) / self.initial_equity) * 100,
            "持倉狀態": "持有" if e.open_position else "無持倉",
            "持倉方向": 1 if e.position_side > 0 else (-1 if e.position_side < 0 else 0),
            "持倉數量": e.position_qty,
            "持倉方向中文": "多單" if e.position_side > 0 else ("空單" if e.position_side < 0 else "無"),
            "入場價": e.position_entry if e.open_position else 0,
            "停損價": e.position_stop if e.open_position else 0,
            "目標價": e.position_target if e.open_position else 0,
            "日損點數": e.daily_loss_pts,
            "週損點數": e.weekly_loss_pts,
            "回撤 %": ((e.peak - e.equity) / e.peak * 100) if e.peak > 0 else 0,
            "風險門檻": 35050 * 1.5,
            "動態優化門檻": 75000 * 1.5,
        }
        return state

    def compute_unrealized_pnl(self, price: float = None):
        """計算未實現損益（若有持倉且提供當前價格）"""
        price = price or self.current_price
        state = self.get_engine_state()
        if not self.engine.open_position or price is None:
            return 0, 0
        direction = state["持倉方向"]
        entry = state["入場價"]
        qty = state["持倉數量"]
        # 微型台指: 1 點 = 10 元
        unrealized_points = (price - entry) * direction
        unrealized_money = unrealized_points * qty * 10 - (35 * 2)  # 扣手續費估算
        return unrealized_points, unrealized_money

    def generate_report(self, price: float = None, output="live_monitor_report.md"):
        self.set_current_price(price)
        state = self.get_engine_state()
        upnl_pts, upnl_money = self.compute_unrealized_pnl(price)

        md = f"""=== 實時交易監控報表 ===
生成時間: {pd.Timestamp.now().strftime('%Y-%m-%d %H:%M:%S')}
引擎來源: live_75k_final.py (LiveTMFEngine)

=== 引擎參數 ===
初始權益: NT$ {self.initial_equity:,}
風險比例: 1% (NT$ {self.initial_equity*0.01:,.0f})
最大口數: 1
標的規格: 微型台指 (TMF) — 乘數 10 元/點
閾值: 0.62 | MC 變異閾值: 0.01
動態優化門檻: NT$ {75000*1.5:,}

=== 當前狀態 ===
權益: NT$ {state['權益']:,.0f}
峰值權益: NT$ {state['峰值權益']:,.0f}
報酬率: {state['報酬 %']:.1f}%
最大回撤: {state['回撤 %']:.2f}%

持倉狀態: {state['持倉狀態']}
持倉方向: {state['持倉方向中文']}
持倉數量: {state['持倉數量']} 口
入場價: {state['入場價']:.0f}
停損價: {state['停損價']:.0f}
目標價: {state['目標價']:.0f}
日損點數: {state['日損點數']:.0f}
週損點數: {state['週損點數']:.0f}

=== 未實現損益 ===
"""
        if self.engine.open_position and price is not None:
            md += f"""當前價格: {price:.0f}
未實現點數: {upnl_pts:.0f}
未實現金額 (扣手續費估算): NT$ {upnl_money:,.0f}

=== 風險檢查 ===
風險門檻 (1.5x 保證金): NT$ {state['風險門檻']:,.0f}
目前權益 / 門檻: {(state['權益']/state['風險門檻'])*100:.1f}%
"""
        else:
            md += """當前無持倉，或未提供即時價格 → 無未實現損益計算
"""
        md += f"""
=== 備註 ===
此報表由 trade_report_engine 連接 live_75k_final.py 的 LiveTMFEngine 產生。
可定時執行或接入真實行情源 (如 Shioaji / 期貨商 API) 進行即時監控。
"""
        with open(output, 'w', encoding='utf-8') as f:
            f.write(md)
        print(f"實時監控報表已生成: {output}")
        print(f"權益: NT$ {state['權益']:,.0f} | 持倉: {state['持倉狀態']} | 報酬: {state['報酬 %']:.1f}%")
        if self.engine.open_position and price is not None:
            print(f"未實現: {upnl_pts:.0f} 點 = NT$ {upnl_money:,.0f}")
        return md

if __name__ == "__main__":
    # 範例: 從 parquet 讀取最新價格作為模擬即時價格
    try:
        df_price = pd.read_parquet(r"C:\Users\yp86n\OneDrive\文件\Default Project\mtxf-ai-trader\data_store\TXF_15m_2006_20260901.parquet")
        latest_price = df_price.sort_index()['close'].iloc[-1]
        print(f"模擬即時價格 (最新 15m 收盤): {latest_price}")
    except Exception as e:
        latest_price = 23000.0
        print(f"無法讀取即時價格，使用預設: {latest_price} (錯誤: {e})")

    monitor = LiveMonitorReport(engine=engine_mod.LiveTMFEngine(), initial_equity=75000)
    # 模擬一個持倉狀態（僅示範）
    # 實際執行時，引擎會從實盤狀態反映真實持倉
    monitor.generate_report(price=latest_price, output="live_monitor_report.md")
    print("\n=== 即時監控已啟動 ===")
    print("可定時執行: python live_monitor_report.py")
    print("或整合至實盤引擎的事件循環中")
