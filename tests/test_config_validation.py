"""TMF 參數配置驗證單元測試 (四件套驗證)

涵蓋 FR-003 至 FR-007：保證金足夠、停損範圍、tick_value 非零、風險比例警告、結構化回傳。
"""

from tmf_engine import Account, TMFDecisionEngine, TMFSpec, validate_config


def _normal_spec_account():
    return TMFSpec(), Account(equity=75000)


def test_valid_normal_params():
    """FR-001/FR-007: 正常參數應回傳 valid=True, errors 為空, 結構完整."""
    spec, account = _normal_spec_account()
    result = validate_config(spec, account)
    assert isinstance(result, dict)
    assert result["valid"] is True
    assert result["errors"] == []
    assert "warnings" in result
    assert "checks" in result
    for check in result["checks"]:
        assert "name" in check
        assert "passed" in check
        assert "detail" in check


def test_margin_insufficient():
    """FR-003: 權益小於 1.5 倍原始保證金應 valid=False."""
    spec, account = _normal_spec_account()
    account.equity = 10000  # 小於 16100 * 1.5 = 24150
    result = validate_config(spec, account)
    assert result["valid"] is False
    assert any("margin" in e.lower() for e in result["errors"])


def test_tick_value_zero():
    """FR-006: tick_value = 0 應 valid=False."""
    spec, account = _normal_spec_account()
    spec.tick_value = 0
    result = validate_config(spec, account)
    assert result["valid"] is False
    assert any("tick_value" in e.lower() for e in result["errors"])


def test_stop_range_invalid():
    """FR-005: min_stop_points > max_stop_points 應 valid=False."""
    spec, account = _normal_spec_account()
    spec.min_stop_points = 80
    spec.max_stop_points = 40
    result = validate_config(spec, account)
    assert result["valid"] is False
    assert any("stop_range" in e.lower() for e in result["errors"])


def test_risk_ratio_warning():
    """FR-004: max_risk_pct > 5% 應為警告而非錯誤，valid 仍為 True."""
    spec, account = _normal_spec_account()
    account.max_risk_pct = 0.10  # 10% > 5%
    result = validate_config(spec, account)
    assert result["valid"] is True
    assert len(result["warnings"]) > 0
    assert any("risk" in w.lower() for w in result["warnings"])


def test_engine_config_validation_attached():
    """FR-002: TMFDecisionEngine 初始化後 config_validation 應為 dict 且非 None."""
    engine = TMFDecisionEngine(Account(equity=75000))
    assert engine.config_validation is not None
    assert isinstance(engine.config_validation, dict)
    assert "valid" in engine.config_validation


def test_engine_main_program_still_works():
    """SC-004: 現有 tmf_engine.py 主程式仍可正常執行."""
    engine = TMFDecisionEngine(Account(equity=75000))
    assert engine.account.equity == 75000
    assert engine.config_validation["valid"] is True
