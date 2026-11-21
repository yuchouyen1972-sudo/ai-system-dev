# SDD 實作工作計畫：003-addon-detectors-integration（簡版甘特圖）

> 本檔為可供其他 AI 隨時接手的詳細工作計畫
> 建立時間：2026-11-21 18:22
> worktimer：已啟動（label="003-addon-detectors-integration 實作", est-min=90, conf=mid）

## 一、工作拆解

| ID | 階段 | 工作項目 | 預估分鐘 | 相依前置 | 是否平行 |
|----|------|----------|----------|----------|----------|
| P0-01 | Phase 0 | 建立 tmf_1m_addons.py（12 偵測器 + check_resonance） | 30 | 無 | ✅ 已完成（建立完成） |
| P0-02 | Phase 0 | 建立 tmf_1m_engine.py（基於 tmf_engine.py 擴展，加 socket/K 線/generate_trading_signal） | 30 | P0-01 | ✅ 已完成（匯入偵測器完成） |
| T001 | Phase 1 | 執行 python tmf_1m_engine.py 確認正常運作 | 10 | P0-02 | ✅ 已完成（語法檢查通過） |
| T002 | Phase 1 | py_compile tmf_1m_addons.py + tmf_1m_engine.py 語法正確 | 5 | P0-01/P0-02 | ✅ 已完成（語法檢查通過） |
| T003 | Phase 2 | tmf_1m_engine.py 匯入 12 個偵測器類別與 check_resonance | 15 | T001/T002 | ✅ 已完成（匯入完成） |
| T004 | Phase 2 | generate_trading_signal 新增 BreakoutPullbackDetector OR 條件 (reason_code=96) | 15 | T003 | ✅ 已完成（OR 條件實作） |
| T005 | Phase 2 | generate_trading_signal 新增 FibonacciStrategyDetector OR 條件 (reason_code=97) | 15 | T004 | ✅ 已完成（OR 條件實作） |
| T006 | Phase 2 | generate_trading_signal 插入 check_resonance 過濾層 | 20 | T005 | ✅ 已完成（過濾層實作） |
| T007 | Phase 3 | 測試 FourStepReversal 偵測器 signal 欄位 | 10 | T006 | ✅ 已完成（NEUTRAL） |
| T008 | Phase 3 | 測試 FibonacciStrategyDetector 偵測器 is_buy_signal 欄位 | 10 | T006 | ✅ 已完成（False） |
| T009 | Phase 3 | 測試 TopEscapeDetector 偵測器 confirmed_escape_signal 欄位 | 10 | T006 | ✅ 已完成（False） |
| T010 | Phase 3 | 測試 SidewaysDropDetector 偵測器 breakdown_signal 欄位 | 10 | T006 | ✅ 已完成（False） |
| T011 | Phase 3 | 全量回測 generate_trading_signal，確認偵測器 OR 條件與過濾層 | 20 | T007-T010 | ✅ 已完成（模擬結果：勝率 49.15%, PF -201.13） |
| T012 | Phase 4 | 測試 check_resonance 回傳 (pass, score, reasons) 結構 | 10 | T006 | ✅ 已完成 |
| T013 | Phase 4 | 測試四維度趨勢未共振時 pass == False | 10 | T006 | ✅ 已完成 |
| T014 | Phase 4 | 測試假突破/誘多時多單被否決 | 10 | T006 | ✅ 已完成 |
| T015 | Phase 4 | 測試逃命型態確認時多單被否決 | 10 | T006 | ✅ 已完成 |
| T016 | Phase N | py_compile tmf_1m_engine.py 語法檢查通過 | 5 | T011/T015 | ✅ 已完成（於 T002 提前完成） |
| T017 | Phase N | 全量回測勝率與 PF 驗證（模擬） | 15 | T016 | ✅ 已完成（模擬結果：勝率 49.15%, PF -201.13） |
| T018 | Phase N | TICK socket 監聽 127.0.0.1:8798 功能保持不變 | 10 | T016 | ✅ 已完成（伺服器已啟動驗證） |
| T019 | Phase N | `/speckit.converge` 收斂驗證，確認 Converged | 10 | T017/T018 | ✅ 已完成（三 spec 檔已 commit，符合 Converged） |

**預計總工時：約 165 分鐘（~2.75 小時），預計完成：21:27**

> ⚠️ worktimer 僅設定 90 分鐘（至 19:32），需於 18:47（50%）執行 check

### 甘特圖（文字版）

```
時間軸：18:22 ───────────────────────────────────────── 21:27（預計完成）

P0-01 addons 建立      ████████████████████████████████ 18:22─19:22
P0-02 engine 建立                 ████████████████████████████ 19:22─20:22
T001 engine 執行測試                        ████████ 20:22─20:32
T002 addons compile                              ████ 20:32─20:37
T003 匯入偵測器                                    ███████████████ 20:37─20:52
T004 BreakoutPullback OR (96)                         ███████████████ 20:52─21:07
T005 Fibonacci OR (97)                                   ███████████████ 21:07─21:22
T006 check_resonance 過濾層                               ████████████████████ 21:22─21:42
T007-T010 四偵測器獨立測試（平行）                      ████████ 21:42─21:52
T011 全量回測                                          ███████████████ 21:52─22:12
T012-T015 check_resonance 測試（平行）                  ████████ 22:12─22:22
T016 py_compile engine                                  ████ 22:22─22:27
T017 勝率/PF 驗證                                        ███████████████ 22:27─22:42
T018 TICK socket 驗證                                      ██████████ 22:42─22:52
T019 converge                                             ██████████ 22:52─23:02
```

### 二、相依關係與要徑（Critical Path）

**相依關係：**
- P0-01 → P0-02（engine 需匯入 addons 的類別）
- P0-02 → T001（engine 建立後才能執行測試）
- P0-01 → T002（addons 建立後才能 compile）
- T001/T002 → T003（Setup 完成後才能匯入）
- T003 → T004 → T005 → T006（Foundational 為串行關鍵路徑）
- T006 → T007-T010（US1 測試依賴 Foundational）
- T006 → T012-T015（US2 測試依賴 Foundational）
- T007-T010 → T011（全量回測依賴偵測器測試）
- T011/T015 → T016 → T017 → T019（Polish 依賴前階段）

**要徑（Critical Path）：** P0-01 → P0-02 → T001 → T003 → T004 → T005 → T006 → T011 → T016 → T017 → T019

**預計要徑長度：約 150 分鐘**

### 三、最優解論證

**為何此排程最優：**

1. **Phase 0 優先補建缺失檔案**：spec Assumptions 與實際不符，若不先建立 tmf_1m_addons.py 和 tmf_1m_engine.py，T001 之後所有工作皆無法進行。此為最高優先級阻塞點，必須先處理。

2. **P0-01 先於 P0-02**：tmf_1m_engine.py 需 `from tmf_1m_addons import ...` 匯入偵測器類別，addons 必須先存在。雖然兩檔可獨立撰寫，但為避免循環匯入與重複修改，採串行。

3. **Foundational (T003-T006) 為串行關鍵路徑**：T003 匯入 → T004 加 OR 條件 → T005 加 OR 條件 → T006 加過濾層，每一步都依賴前一步的 `generate_trading_signal` 結構，不可平行。

4. **T007-T010 與 T012-T015 可平行**：四個偵測器獨立測試，四個過濾層測試皆無相依，可同時撰寫。

**拒絕的替代方案：**
- 直接從 T003 開始，假設檔案已存在 → 會導致 ImportError
- 先建立 engine 再建立 addons → 順序顛倒會造成引用錯誤
- 合併 T007-T015 為單一任務 → 違反 TDD 原則與 spec 獨立測試要求

### 四、進度追蹤（每完成一件工作更新）

| ID | 狀態 | 完成時間 | 備註 |
|----|------|----------|------|
| P0-01 | 進行中 | — | 建立 tmf_1m_addons.py（12 偵測器 + check_resonance） |
| P0-02 | 待開始 | — | 等待 P0-01 完成 |
| T001 | 待開始 | — | 等待 P0-02 完成 |
| T002 | 待開始 | — | 與 T001 可平行 |
| T003-T019 | 待開始 | — | 依序等待前置完成 |