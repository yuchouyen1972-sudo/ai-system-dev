# AI 工廠持續計畫（進度唯一真相來源）

> 更新：2026-11-13
> 狀態：Spec Kit 與現有開發流程相互整合優化已完成

## 已完成

- Spec Kit (specify-cli 1.4.0) 已安裝，整合代理為 opencode，腳本類型 ps
- `specs/` 子目錄已初始化，內含 `.specify/`（extensions/integrations/memory/scripts/templates/workflows）
- 已安裝擴充：bug (Bug Triage Workflow v1.0.2)、assess (Idea Assessment Pipeline v1.0.2)
- 規範文件已修訂：FACTORY.md、SOP.md、STARTUP.md、AGENTS.md（新增第十條 Spec Kit 規格層整合）
- `.gitignore` 已建立，排除 `.opencode/`、*.pkl、*.csv、*.log 等
- 整合範本：TMF 參數配置驗證（spec/plan/tasks 三文件、validate_config 加入 tmf_engine.py、tests/test_config_validation.py 7 passed）
- git commit 已完成：cb623a6 feat(spec-kit): integrate SDD specification layer with existing dev workflow
- git push：未執行（本地無 remote 儲存庫，待使用者提供 remote URL）

## 下一個待辦

- 使用者在 OpenCode 對話中執行 `/speckit.converge` 完成 SDD 收斂驗證（Converged）
- 如需 push 至遠端，請提供 remote URL 並執行 `git remote add origin <url>` 後 `git push -u origin master`

## 相依關係

- .omo/plans/ 與 specs/.specify/memory/ 互不覆蓋：前者管整體里程碑，後者管個別規格 SDD 狀態
- 四件套（pytest/ruff/format/mypy）＋ CI 門禁維持不變，所有整合後代碼須通過
- worktimer 紅線機制貫穿整個 SDD＋實現流程