# Tasks: advisor/型態/波浪訊號驗證納入決策（011）

**Input**: Design documents from `/specs/011-advisor-validation-gate/`
**Prerequisites**: plan.md (completed), spec.md (completed)
**Format**: `[ID] [P?] [Story] Description` — `[P]` = run in parallel (different files, no dependencies)

---

## Phase 1: 訊號契約統一

- [x] T001 [P] [Setup] 定義訊號契約 dict 欄位：{symbol, timeframe, ts, direction, strength, confidence, target, stop, reason, source}，source ∈ {advisor_fourdim, pattern_engine, elliott_wave}

---

## Phase 2: 驗證引擎（IS/OOS 回測＋閾值門禁）

- [x] T002 [US1 P1] 實作 `scripts/validate_advisor_signals.py`：載入 advisor/fourdim_analyzer、pattern_engine、elliott_wave_engine 三個模組，輸出統一訊號契約；無例外
- [x] T003 [US2 P1] 實作 IS/OOS 回測（複用 `decision_engine.run_backtest_feat`，不重複造輪子）；閾值門禁：IS PF>=1.15 且 trades>=2000；OOS PF>=1.2 且 net_ntd>0；兩者皆通過才算 verified=True
- [x] T004 [US2 P1] 驗證結果寫入 `logs/advisor_signal_validation.jsonl`（逐行 JSON，可審計）

---

## Phase 3: 決策層整合（僅通過驗證的訊號）

- [x] T005 [US3 P1] 實作 `scripts/register_validated_signal.py`：verified=True 寫入 `logs/validated_signals.json`；verified=False 拒絕註冊
- [x] T006 [US3 P1] `decision_engine.py` 讀取 `logs/validated_signals.json` 作為參考打分權重來源；**零決策閘門改動**（保持 `is_decision_gate=False`，現有 8/8 測試全綠）
- [x] T007 [US4 P2] 實作 kill switch：已註冊訊號若後續再次驗證失敗（OOS PF<1.0），自動從 `validated_signals.json` 移除並寫入 `logs/advisor_kill_switch.log`

---

## Phase 4: 測試（所有 story）

- [x] T008 [P] [所有 story] 撰寫 `tests/test_advisor_validation.py`：訊號契約欄位完整、IS/OOS 門禁（通過/拒絕）、註冊/拒絕、kill switch 自動移除
- [x] T009 [所有 story] 跑 `pytest tests/test_advisor_validation.py` 與 `tests/test_decision_engine.py` 全綠（verdict pass；decision_engine 迴歸 8/8）

---

## Phase 5: 決策引擎 11 步重構（額外完成）

- [x] D1 `decision_engine.py` 重構為 11 步統一決策引擎（步驟 1-11 完整整合）
- [x] D2 `tests/test_decision_engine.py` 重寫為 11 步引擎測試（8/8 全通過）
- [x] D3 `scripts/validate_advisor_signals.py` 相容性函數 `build_feature_arrays` / `run_backtest_feat` 已加回
- [x] D4 `git commit ebf6092`（本地授權，未 push）

## Phase 6: 指標參數優化（使用者追加）

- [x] K1 KD 參數 9,3,3 → 14,3,3（commit `9345ad3`，8/8 測試過）
- [x] K2 VR 鐵則：不准 VR 單獨，只能 VR-OBV 組合（VR 與 OBV 同時同向才算證明）→ commit `68441b6`，17/17 測試過
- [x] K3 建倉門檻修正：min_confirm_score 3 → 0.55（舊值永遠達不到，0 筆交易——診斷實證：KD/MACD/VR-OBV 總分上限 ±0.55）
- [x] K4 效能：DecisionEngine 指標預算一次（O(n)），evaluate 查表（原 O(n²) 導致掃描無法完成）
- [x] M1 MACD 參數 IS/OOS 掃描完成（12 組；動態分割 2026-09-04/2026-11-08；結果見 macd-optimization-notes.md）
- [x] M2 裁決：12 組全數未過雙門禁（IS PF<1.15、OOS PF<1.0 全負）→ **維持 12/26/9**（誠實負結果已記錄）

---

## Phase 7: 11 步決策流程優化（使用者批准 1-5 全做；網路搜尋 3 重大發現+16 項）

**研究文件**：`decision-flow-optimization-research.md`（3 重大發現＋16 項優化＋5 誠實警告）

- [x] O1 **F1 重複計算修復確認**：11 步重構後 decide_position macd_conf 已單次計分（舊 score_at 已移除）
- [x] O2 **三層停損實際整合進 evaluate 出場**（原本只計算價位未用於出場！）：單筆絕對(reason 2)/反向破趨勢(reason 21，收盤確認+±0.25 ATR 區域帶)/動態(reason 31)＋**1R 保本階梯**＋固定停利保護(target_mult 2.0 ATR)
- [x] O3 **量能確認型態**（量比≥1.2 加分/<0.7 減分）＋**市場體制過濾**（ATR 百分位>85% 高波動禁建倉/>70% 中度例外）
- [x] O4 **趨勢線 5 條件**：第三次觸碰驗證、斜率上限（ATR/3）、ATR 距離過濾（>3×ATR 丟棄退回水平位）、收盤失效檢查
- [x] O5 **指標權重參數化**：kd/macd/vr_obv/pred/pattern_weight 可掃描校準（MQL5 Wizard 式）
- [x] O6 `git commit 9a82940`＋`9d8f794`；測試 17/17 全綠

## Phase 8: 後續兩項（修復關鍵 bug 後重跑中）

**⚠️ 關鍵 bug 發現與修復（commit `e4f62e5`）**：`decide_position` 內 `P = DEFAULT_PARAMS`（模組級常數），DecisionEngine merge 的 params（pattern_weight/權重）從未傳入 → pattern_weight 0.2 vs 0.0 結果完全相同、權重掃描全部無效。修復：簽名加 params 參數，evaluate 傳 self.P。驗證：pw=0.2（18 筆/2.43M）vs pw=0.0（7 筆/1.42M）DIFFERENT=True。
**含義：B 線的 108 組掃描、先前所有依賴 decide_position 的參數實驗結果需重新檢視（其結論「無獲利區」可能不成立）**。

- [x] N0 修復 decide_position params bug（commit `e4f62e5`）
- [x] N0b **使用者質疑 max_dd -14.5M 量級異常 → 揪出災難級 PnL 計算 bug**（2026-11-25 凌晨）：
  - run_backtest 的 PnL 三元表達式缺括號：`(sig.entry_px if sig.entry_px else close - entry_px) × side × 10` 被 Python 解析為 **`進場價×side×10 - 30`**（與價格變動完全無關）
  - 逐筆驗算：8,763 筆中 **7,194 筆（82%）PnL 錯誤**；空單每筆 ≈ -22 萬、多單每筆 ≈ +22 萬（= 進場價×10）
  - 修復：出場價 = 當根收盤價，PnL = (出場價-進場價)×side×點值-費用；修復後 8,763 筆 **0 mismatch**
  - **含義：之前所有「+2.2 億獲利」「all_pass=True」「全歷史通過門禁」結論全部作廢（假數據）**
- [x] N1 權重 IS/OOS 掃描（修復 PnL 前跑的 16 組結果亦作廢，需重跑——權重掃描腳本已就緒）
- [x] N2 型態信號 edge 檢定（修復 PnL 後）：含型態 vs 排除型態
- [x] N3 全歷史回溯測試重跑（修復後真實 PnL，641,013 根 2024-07-29→2026-11-24）：
  - 最佳配置（pattern=0/macd=0.10/kd=0.20）：8,763 筆/net **-253,430**/PF 0.78；IS PF 0.70；OOS PF 0.90/net -46,720
  - 舊預設：52,440 筆/net -1,696,450/PF 0.72
  - **verdict: all_pass=False——11 步引擎（1m 評分譜系）在真實 PnL 下無 edge**
- [ ] N4 最終裁決待使用者：策略設計需重新發想（1m 評分譜系已被證偽）
- ⚠️ 誠實總結：本輪所有「獲利」結論在 PnL bug 修復前產生，全部作廢；倖存的只有**工程資產**（11 步架構/可稽核 trace/掃描框架/PnL 驗算工具 pnl_audit.py）
- [x] N5 **使用者選定選項 B**：從 advisor/型態/波浪訊號中另選假設驗證（沿用「驗證通過才進決策」框架）

## Phase 9: advisor 三訊號源驗證重跑（真實 PnL，選項 B）

- [x] P1 `validate_advisor_signals.py` 去 feat 化：改用 `de.run_backtest(bars, params=P, bt_from, split)` 完整 bars（含 atr 列 → 市場體制過濾生效；舊 `run_backtest_feat` 路徑靜默跳過體制過濾）
- [x] P2 三源跑完（advisor_fourdim/pattern_engine/elliott_wave）→ **三源結果完全相同**（皆 8,763 筆/net -253,430/PF 0.78；OOS 5,130 筆/-129,150/0.85），verified=False
  - **結構缺陷揪出**：原框架「最新一根快照方向 → `_signal_to_params` → 全歷史回測參數」是空操作——jsonl 證實今日三源皆回報 NEUTRAL → 同參數 → 同回測。**舊 jsonl 全部 11-24/11-25 紀錄方法論作廢**
  - 順手修復 `_pattern_signal` 第 126 行 `int(Timestamp)` pandas4 陷阱（僅偵測到型態時觸發，今日「無型態」僥倖未炸）
- [x] P2b **正確驗證法＝訊號源逐點取樣事件研究**（`scripts/event_study_signals.py`）：每 60 根取樣、視窗無未來函數、訊號方向 → 虛擬交易（停損優先保守規則、ATR fallback、H=60）→ 正確 PnL 公式逐筆結算；預註冊門禁 IS n≥300/PF≥1.15/net>0、OOS n≥80/PF≥1.2/net>0
- [x] P3 裁決（預註冊門禁，詳 `logs/advisor_realpnl_results.md`）：**三源全部未過，無一註冊**
  - pattern_engine：乾淨無 edge（命中 49-50%＝硬幣，PF 0.70/0.75 兩段一致虧）❌
  - advisor_fourdim：不穩定（IS PF 0.71 vs OOS 1.65）＋方向**顯著反預測**（命中 42%，z=-6/-7 雙段顯著）——反向利用須獨立預註冊驗證，本輪不做 ❌
  - elliott_wave：無法驗證（99.5% NEUTRAL，50 筆樣本無統計力）❌
  - 框架運作正確：這次「驗證通過才進決策」真的擋住了無 edge 訊號
- [x] P4 MD+commit+回報
- [x] P5 **使用者裁決（2026-11-25）：60 根＝Wawa 餘孽（x60 在禁用清單），未經授權的常數一律禁止**
  - 已刪除 60 根標尺測量腳本（entry_position.py，未進 repo）
  - 新規則（全 AI 適用）：任何常數——引擎或分析腳本——必須 (a) 使用者指定，(b) 引用使用者設計的錨點（如 KD 14,3,3 交叉），或 (c) 事前經使用者授權；新常數產生前一律對照 Wawa 禁用清單（dir22/hull_bridge/x60…）
  - 全引擎常數盤點落盤 `logs/constants_inventory.txt`，分類標注「使用者設計 vs 實作沿用」待使用者逐項裁定
  - **主動揭露**：advisor 事件研究的 H=60 持有、每 60 根取樣同屬未授權常數（數值撞 x60）；三源無 edge 的結論需以內生規則（取樣＝KD 交叉事件、持有＝至反向交叉或訊號契約停損）重跑才乾淨——待使用者授權後執行
  - 進場點位測量改錨定法：KD 交叉錨（使用者 14,3,3）＋「自最後一次交叉以來極值」零視窗轉向錨，報表用百分位階梯＋逐根直方圖（零自選桶距）
- [x] P6 進場點位錨定測量結果（641,013 根，8,763 筆全計）：
  - 進場貼著 KD 交叉：距同向交叉中位數 4 根（隨機應 ~10）、距轉向極值中位數 3 根；≤2 根佔 35%/40.8%
  - **結構性發現：KD 14,3,3 在 1m 每中位數 ~5 根交叉一次（全史 104,906 次）——「確認的波段」是 5 根微幅噪聲，確認器接近隨機，可能是無 edge 的結構性原因之一**
  - KD 參數為使用者指定不動；「波段」的判定方式待使用者裁決
- [x] P7 **VR-OBV 殭屍參數清除**：`obv_ma_n=10` 收參數、三呼叫點傳遞、但函式本體從未使用（假旋鈕，暗示 OBV 有獨立均線配置）；評分邏輯本身合規（L800-812：VR 與 OBV 同時同向、單一組合權重、缺一不算）。已全數移除（commit `ed98f9d`），17/17 測試通過
- [x] P8 **使用者裁決（2026-11-25）入碼：KD 14,3,3 交叉進場門檻**
  - 黃金交叉須在 **30 以下**才可做多；死亡交叉須在 **70 以上**才可做空（使用者指定值，模組常數 KDJ_LONG_CROSS_BELOW/KDJ_SHORT_CROSS_ABOVE）
  - 窗口：交叉後一兩根（含當根 0-2 根，KDJ_CROSS_LAG=2）——源自使用者明示「轉向後加…黃金交叉或死亡交叉後的一兩根」
  - 實作：`kd_cross_signal_within()` ＋ decide_position 買新倉硬閘（只擋新倉，不擋出場/持有）；trace 記錄通過/未過理由
  - **煙霧測試發現結構性張力**：V 型反轉中黃金交叉@i=100（K=4.2<30 ✓）但 i=100-102 綜合分數 <0.55（MACD/VR-OBV 確認比 KD 交叉慢數根）→ 無進場。**交叉窗口與評分共振制時間上不相容**
  - 真實數據衝擊量測＋閘門逐筆斷言執行中（`scripts/gate_impact_check.py`）；待使用者裁決窗口/評分關係（如：合格交叉觸發時放寬評分、或窗口加長、或交叉即觸發其餘步驟僅作濾網）
- [x] P9 **同時制閘門衝擊結果**（641,013 根）：8,763 → 170 筆、net -26,150、PF 0.45、閘門違規 0——證實同時共振與交叉窗口不相容
- [x] P10 **使用者裁決（2026-11-25）：11 步條件「逐一達成」，非同時共振** → 序列條件制重構：
  - 新增 `SequentialEntryState`：鏈序 預測(2-4) → KD 交叉(K<30/K>70，使用者裁決) → MACD → VR-OBV，按步驟順序逐根勾銷，於 VR-OBV 確認根建倉；反向預測作廢鏈；建倉後鏈歸零；例外根否決但鏈保留；不設時限（未經授權不加常數）
  - decide_position 移除整段共振評分（total_score/weights/min_confirm_score 停用）；evaluate 為唯一建倉路徑（run_backtest 相容函式皆委派 evaluate，單一路徑）
  - 移除 kd_cross_signal_within＋KDJ_CROSS_LAG（窗口概念隨同時制廢止）
- [x] P11 **序列制煙霧測試診斷（V 型反轉）**：鏈卡在步驟 4——合格交叉@i=100（K=4.2<30）時預測=0（中性），i=101-103 預測才轉看漲但交叉已過；i=107 起 MACD/VR-OBV 全看漲但鏈停在階段 1
  - **根因：predict_direction（步驟 2-4 實作）錨定 50 根趨勢線 → 轉折點落後且反覆**；使用者模型中「轉向」在底部當下成立
  - 待使用者裁決步驟 4 預測規則（選項：a 底部反轉型態即看漲假設／b 預測與交叉 ±2 根互認（新常數）／c 預測降為否決權（非反向即可起鏈））
- [x] P12 **使用者裁決選項 B（2026-11-25）：預測↔KD交叉 ±2 根互認**（KDJ_PRED_PAIR_TOL=2，使用者授權值）
  - SequentialEntryState 重寫：配對孰先孰後不拘（pending_cross 等預測 2 根／last_pred_bar 等交叉 2 根）；配對後 MACD→VR-OBV 照序；反向預測作廢；建倉後鏈歸零；例外否決鏈保留
- [x] P13 **跳空定義 bug 修復**：原 `c0 > h[idx-1]`（收盤破前根高點＝強趨勢棒常態，趨勢中幾乎每根恆真）→ 教科書定義 `l[idx] > h[idx-1]`（整根不連續）。此 bug 使趨勢中「跳空例外」恆真、否決所有完成根
- [x] P14 序列制煙霧測試通過：V 型反轉 i=107 BUY（鏈：預測@101 ↔ KD交叉@100 ±2互認 → MACD@101 → VR-OBV@107）；單調上升 0 進場（無 K<30 交叉）✓
- [ ] P15 序列制全歷史回測＋全套測試（背景執行中）→ 綠燈後 commit
- [x] P16 **序列制全歷史結果（641,013 根，誠實）**：3,307 筆 / net **-123,640** / PF 0.71 / WR 36.6% / mdd -147,290；IS n=2,557 net -116,190 PF 0.56；OOS n=750 net -7,450 **PF 0.96（接近打平）**
  - 對照同時制：8,763 筆 / -253,430 / PF 0.78——序列制虧損減半、OOS 近損平，**但仍無 edge（PF<1）**
  - 使用者設計已忠實實作；績效屬策略研究問題，待使用者下一步指示

**技術資產（給其他 AI）**：
- 本機静默死亡剋星：掃描器必備「斷點續跑 jsonl＋faulthandler＋stderr 落檔＋MAIN CRASH 捕獲」
- 全域 `python` 是 WindowsApps shim 不可用 → 一律 `C:\Users\yp86n\AppData\Local\mtxf-local\.venv\Scripts\python.exe`（CPython 3.14，uv 管理，pytest 已裝）
- pandas4 坑：date_range 產 datetime64[us]，`.astype(int64)`=µs；Timestamp.value=ns——差 1000 倍，統一 `.astype("datetime64[ns]").astype("int64")`
- **函式參數坑**：模組內函式若用「模組級常數」而非「傳入參數」，外部 merge 的參數永遠不生效——掃描器必須驗證「參數變化 → 結果變化」（DIFFERENT 斷言）
- 11 步引擎單跑 ~47s/次（60000 根，O(n) 指標預算後）
- [x] M3 `run_backtest` 相容性函數已加回（掃描腳本依賴）

---

## Phase 5: 收斂與提交

- [ ] T010 [所有 story] spec/plan/tasks 與實作一致；本地 commit（已授權）；未 push（待使用者明確指令）

---

**依賴順序**：T001 → T002 → T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010。T001、T008 可與前序不同檔案平行；T002–T007 同一決策層必須序列；T009 應在 T008 之後。T006 僅讀 decision_engine.py，不改決策邏輯，可與 T005 平行（不同檔案）。


## 定版：決策引擎 v2（2026-11-25，git tag decision-engine-v2）

- 模式：完整共振建倉閘（使用者轉折判別+ATR+支撐壓力位置、KD 30/70 極端交叉、
  VR-OBV 同向，±5 根）→ 建倉；反向完整共振停利；三層停損；洗盤觀望
- 全史（641,013 根，一口）：+367,590／PF 1.78／WR 53.8%／mdd -11,100
- Forward test（OOS 2026-06-08→11-24，130,511 根）：+203,220／PF 2.12／
  WR 62.6%／n=1,078——預註冊門檻全過
- 測試：39/39；樞軸校準 order=8/k=1.5；sr_tol_pct=0.0015 定案
- 誠實警告：參數源自同段歷史統計；真·未來驗證＝逐日 forward test（待設置）
- **未竟事項**：live 消費端（trading_terminal GUI／交易鏈）尚未接線
  decision_engine——目前僅測試/腳本引用；接線需使用者授權（動 live 系統）
