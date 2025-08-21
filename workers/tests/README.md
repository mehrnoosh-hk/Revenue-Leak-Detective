# Test Suite Documentation

This directory contains comprehensive tests for the Revenue Leak Detective Agent workflow.

## Files

- `conftest.py` - Test configuration and shared fixtures
- `test_run.py` - Complete test suite for the agent workflow

## Test Structure

The test suite is organized into the following test classes:

### TestSampleNode
- **Purpose**: Unit tests for the `sample_node` function
- **Coverage**: Input validation, output format, immutability, edge cases
- **Special Cases**: Unicode, whitespace, multiline, numeric strings

### TestCreateWorkflow
- **Purpose**: Unit tests for the `create_workflow` function
- **Coverage**: Workflow creation, compilation, multiple instances

### TestRunAgent
- **Purpose**: Unit and integration tests for the `run_agent` function
- **Coverage**: Default/custom parameters, dry-run mode, various message types
- **Special Features**: Logging verification in dry-run mode

### TestMainFunction
- **Purpose**: Integration tests for the CLI interface
- **Coverage**: Default args, custom flags, custom messages
- **Testing Approach**: Uses monkeypatch to simulate command-line arguments

### TestErrorHandling
- **Purpose**: Error handling and edge case tests
- **Coverage**: Missing keys, invalid data types, exception chaining
- **Note**: All errors are wrapped in `RuntimeError` by the implementation

### TestIntegration
- **Purpose**: End-to-end integration tests
- **Coverage**: Complete workflow testing, consistency verification
- **Scope**: Tests the entire pipeline from input to output

### TestPerformance
- **Purpose**: Performance and timing tests
- **Coverage**: Execution time limits, large input handling
- **Thresholds**: < 1s for normal operations, < 2s for large inputs

## Test Fixtures

The test suite includes comprehensive fixtures in `conftest.py`:

### State Fixtures
- `sample_state` - Basic test message
- `empty_state` - Empty message
- `default_state` - Default "Mehrnoosh" message
- `long_message_state` - 1000 character message
- `special_chars_state` - Special characters
- `unicode_state` - Unicode and emoji characters
- `whitespace_state` - Various whitespace characters
- `numeric_string_state` - Numeric string
- `multiline_state` - Multiline text

### Data Fixtures
- `test_messages` - List of various test messages
- `edge_case_messages` - Edge case scenarios
- `log_capture` - Logging verification fixture

## Test Markers

The test suite uses custom pytest markers for organization:

- `@pytest.mark.unit` - Unit tests (37 tests)
- `@pytest.mark.integration` - Integration tests (7 tests)
- `@pytest.mark.performance` - Performance tests (3 tests)
- `@pytest.mark.error_handling` - Error handling tests (4 tests)
- `@pytest.mark.slow` - Slow-running tests (1 test)

## Running Tests

### Run All Tests
```bash
python -m pytest tests/test_run.py -v
```

### Run by Marker
```bash
# Unit tests only
python -m pytest tests/test_run.py -m "unit" -v

# Integration tests only
python -m pytest tests/test_run.py -m "integration" -v

# Performance tests only
python -m pytest tests/test_run.py -m "performance" -v

# Error handling tests only
python -m pytest tests/test_run.py -m "error_handling" -v
```

### Run by Test Class
```bash
# Sample node tests only
python -m pytest tests/test_run.py::TestSampleNode -v

# Error handling tests only  
python -m pytest tests/test_run.py::TestErrorHandling -v
```

## Test Coverage

The test suite provides comprehensive coverage of:

✅ **Function-level testing** - All public functions tested  
✅ **Input validation** - Various input types and edge cases  
✅ **Error handling** - Exception scenarios and error chaining  
✅ **Integration testing** - End-to-end workflow testing  
✅ **Performance testing** - Timing and large input handling  
✅ **CLI testing** - Command-line interface functionality  
✅ **Logging verification** - Log message capture and validation  
✅ **Unicode support** - International characters and emojis  
✅ **State immutability** - Ensuring inputs are not modified  
✅ **Consistency testing** - Reproducible results  

## Best Practices Implemented

1. **Descriptive test names** - Clear test purpose and expectations
2. **Parametrized tests** - Efficient testing of multiple scenarios
3. **Fixture usage** - Reusable test data and setup
4. **Proper test organization** - Logical grouping using classes
5. **Marker usage** - Easy test filtering and categorization
6. **Error message validation** - Specific exception matching
7. **Performance boundaries** - Reasonable execution time limits
8. **Integration testing** - Real workflow verification
9. **Documentation** - Comprehensive docstrings and comments
10. **No mocking** - Real function testing for authentic behavior

## Test Statistics

- **Total Tests**: 53
- **Unit Tests**: 37 (70%)
- **Integration Tests**: 7 (13%)
- **Performance Tests**: 3 (6%)
- **Error Handling Tests**: 4 (8%)
- **Other Tests**: 2 (3%)

All tests pass consistently and provide reliable validation of the agent workflow functionality.
