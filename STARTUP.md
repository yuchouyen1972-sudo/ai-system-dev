# STARTUP — 新 session 開機協議（opencode 啟動即執行）

> 載入時機：每個新 session 第一次回應前。目的：讓使用者一開就看到「上次停在哪、下一步是什麼」。
> 規則：只偵察、不開工。任何實作／pipeline／push 都要使用者明確一句確認才動手。

## 開機五步（依序，不可跳步）

1. **讀計畫真相**：`.omo/plans/ai-factory-continuation-plan.md`（進度唯一真相來源）。
2. **查實機狀態**：`git status --porcelain`＋`git log --oneline -5`（有無未 commit／上次停在哪）。
3. **驗 key 狀態**：9 個環境變數只報 PRESENT/MISSING（長度可，不印值）：
   `OPENROUTER_API_KEY GROQ_API_KEY GEMINI_API_KEY CLOUDFLARE_API_TOKEN CLOUDFLARE_ACCOUNT_ID FREEMODEL_API_KEY MISTRAL_API_KEY SILICONFLOW_API_KEY ZHIPU_API_KEY`
4. **跑體檢**：`ai-dev-team doctor --no-speed`（Ollama＋遠端啟用狀態）＋ `specify check`（Spec Kit CLI 可用性與 `opencode` 整合狀態）。
5. **回報＋等確認**：用繁體中文回報「上次進度／現況／建議下一步」，然後停住等使用者指令。禁止擅自開工、禁止擅自 push。

## 輸出格式（固定）

```
開機狀態（<日期>）：
- 上次進度：<計畫檔最後狀態，一句話>
- 工作區：<clean／N 檔未 commit>
- key：<PRESENT 幾個／MISSING 哪幾個>
- doctor：<reachable＋啟用中 provider 數>
- specify：<specify-cli 版本＋opencode 整合狀態>
- 今日要徑：<排隊中 critical path 一句話＋下一個待辦工作單號>
- 建議下一步：<一件事>
等你一句確認再動手。
```
