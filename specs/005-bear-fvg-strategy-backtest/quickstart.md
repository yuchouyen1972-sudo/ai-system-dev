# Quickstart: 005 驗證指南 驗證指南

1. `uv run python scripts/test_strategy_bear_fvg.py` → 合成案例結構斷言 ALL PASS（SC-003）＋既有 27 項回歸（SC-006）
2. `uv run python scripts/strategy_bear_fvg.py` → 全史回測 <5 分（SC-001）→ `docs/validation_005/` 產出（SC-002/004/005）
3. 判定閱讀：verdict 以滑價 2＋成本上限格為準（FR-008）；對照歸因表核對 D5

