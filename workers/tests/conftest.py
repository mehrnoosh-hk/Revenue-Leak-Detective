"""Test configuration and shared fixtures for the Revenue Leak Detective Agent.

This module provides pytest fixtures and configuration for testing the agent workflow,
including state fixtures and test utilities.
"""

import pytest
from typing import List
import logging
from io import StringIO

# Import from the correct module path
from src.agent.run import State


# State fixtures for different test scenarios
@pytest.fixture
def sample_state() -> State:
    """Fixture providing a sample state for testing."""
    return {"message": "TestMessage"}


@pytest.fixture
def empty_state() -> State:
    """Fixture providing an empty state for testing."""
    return {"message": ""}


@pytest.fixture
def default_state() -> State:
    """Fixture providing the default state used in run.py."""
    return {"message": "Mehrnoosh"}


@pytest.fixture
def long_message_state() -> State:
    """Fixture providing a state with a long message."""
    return {"message": "A" * 1000}


@pytest.fixture
def special_chars_state() -> State:
    """Fixture providing a state with special characters."""
    return {"message": "Test@123!#$%^&*()"}


@pytest.fixture
def unicode_state() -> State:
    """Fixture providing a state with unicode characters."""
    return {"message": "Hello 世界 🌍 Здравствуй"}


@pytest.fixture
def whitespace_state() -> State:
    """Fixture providing a state with various whitespace characters."""
    return {"message": " \t\n\r Test Message \t\n\r "}


@pytest.fixture
def numeric_string_state() -> State:
    """Fixture providing a state with numeric string message."""
    return {"message": "12345"}


@pytest.fixture
def multiline_state() -> State:
    """Fixture providing a state with multiline message."""
    return {"message": "Line 1\nLine 2\nLine 3"}


# Logging fixtures
@pytest.fixture
def log_capture():
    """Fixture to capture log messages during tests."""
    log_stream = StringIO()
    handler = logging.StreamHandler(log_stream)
    handler.setLevel(logging.INFO)

    # Get the logger from the run module
    logger = logging.getLogger("src.agent.run")
    original_level = logger.level
    original_handlers = logger.handlers.copy()

    logger.setLevel(logging.INFO)
    logger.addHandler(handler)

    yield log_stream

    # Cleanup
    logger.removeHandler(handler)
    logger.handlers = original_handlers
    logger.setLevel(original_level)


# Test data fixtures
@pytest.fixture
def test_messages() -> List[str]:
    """Fixture providing a list of test messages for parametrized tests."""
    return [
        "Simple message",
        "",  # Empty message
        "A" * 100,  # Long message
        "Special!@#$%^&*()Characters",
        "Numbers123456789",
        "Mixed Case Message",
        "\t\nWhitespace\r\n",
        "🎉 Emojis 🚀 Test 🎯",
        "Mehrnoosh",  # Default from run.py
        "Hello 世界",  # Unicode
        "Multi\nLine\nMessage",
        "123.45",  # Numeric string
    ]


@pytest.fixture
def edge_case_messages() -> List[str]:
    """Fixture providing edge case messages for testing."""
    return [
        "",  # Empty string
        " ",  # Single space
        "\n",  # Newline only
        "\t",  # Tab only
        "   \t\n\r   ",  # Mixed whitespace
        "A",  # Single character
        "A" * 10000,  # Very long string
    ]


# Configuration fixtures
@pytest.fixture(autouse=True)
def setup_test_environment():
    """Automatically set up test environment for each test."""
    # Setup code before test
    yield
    # Cleanup code after test (if needed)
    pass


# Pytest configuration
def pytest_configure(config):
    """Configure pytest with custom markers."""
    config.addinivalue_line("markers", "integration: mark test as an integration test")
    config.addinivalue_line("markers", "unit: mark test as a unit test")
    config.addinivalue_line("markers", "performance: mark test as a performance test")
    config.addinivalue_line(
        "markers", "error_handling: mark test as an error handling test"
    )
    config.addinivalue_line("markers", "slow: mark test as slow running")
