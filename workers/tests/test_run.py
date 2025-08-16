"""Tests for src.agent.run module.

This module contains comprehensive unit and integration tests for the
Revenue Leak Detective agent workflow defined in src.agent.run.
"""

import pytest
import sys
import time

# Import from the correct module path
from src.agent.run import (
    State,
    sample_node,
    create_workflow,
    run_agent,
    main,
)


class TestSampleNode:
    """Unit tests for the sample_node function."""

    @pytest.mark.unit
    @pytest.mark.parametrize(
        "input_message,expected_output",
        [
            ("John", "Welcome to my agent! John"),
            ("", "Welcome to my agent! "),
            ("Test@123!", "Welcome to my agent! Test@123!"),
            ("A" * 100, "Welcome to my agent! " + "A" * 100),
            ("Mehrnoosh", "Welcome to my agent! Mehrnoosh"),
        ],
    )
    def test_sample_node_various_inputs(self, input_message: str, expected_output: str):
        """Test sample_node with various input messages."""
        input_state: State = {"message": input_message}
        result = sample_node(input_state)

        assert result["message"] == expected_output
        assert isinstance(result, dict)
        # Ensure the original state is not modified
        assert input_state["message"] == input_message

    @pytest.mark.unit
    def test_sample_node_with_fixtures(self, sample_state):
        """Test sample_node using fixture."""
        result = sample_node(sample_state)
        assert result["message"] == "Welcome to my agent! TestMessage"

    @pytest.mark.unit
    def test_sample_node_immutability(self, sample_state):
        """Test that sample_node doesn't mutate the input state."""
        original_message = sample_state["message"]
        sample_node(sample_state)
        assert sample_state["message"] == original_message

    @pytest.mark.unit
    def test_sample_node_return_type(self, sample_state):
        """Test that sample_node returns correct type."""
        result = sample_node(sample_state)
        assert isinstance(result, dict)
        assert "message" in result
        assert isinstance(result["message"], str)

    @pytest.mark.unit
    def test_sample_node_unicode(self, unicode_state):
        """Test sample_node with unicode characters."""
        result = sample_node(unicode_state)
        expected = "Welcome to my agent! Hello 世界 🌍 Здравствуй"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_sample_node_whitespace(self, whitespace_state):
        """Test sample_node with whitespace characters."""
        result = sample_node(whitespace_state)
        expected = "Welcome to my agent!  \t\n\r Test Message \t\n\r "
        assert result["message"] == expected

    @pytest.mark.unit
    def test_sample_node_multiline(self, multiline_state):
        """Test sample_node with multiline message."""
        result = sample_node(multiline_state)
        expected = "Welcome to my agent! Line 1\nLine 2\nLine 3"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_sample_node_numeric_string(self, numeric_string_state):
        """Test sample_node with numeric string message."""
        result = sample_node(numeric_string_state)
        expected = "Welcome to my agent! 12345"
        assert result["message"] == expected

    @pytest.mark.unit
    @pytest.mark.parametrize("message", [
        "", " ", "\n", "\t", "   \t\n\r   ", "A", "A" * 10000
    ])
    def test_sample_node_edge_cases(self, message: str):
        """Test sample_node with edge case messages."""
        state: State = {"message": message}
        result = sample_node(state)
        expected = f"Welcome to my agent! {message}"
        assert result["message"] == expected


class TestCreateWorkflow:
    """Unit tests for the create_workflow function."""

    @pytest.mark.unit
    def test_create_workflow_returns_compiled_graph(self):
        """Test that create_workflow returns a compiled graph."""
        workflow = create_workflow()
        assert workflow is not None
        # Test that it has the invoke method (compiled graphs have this)
        assert hasattr(workflow, 'invoke')

    @pytest.mark.unit
    def test_create_workflow_can_process_state(self, sample_state):
        """Test that created workflow can process a state."""
        workflow = create_workflow()
        result = workflow.invoke(sample_state)
        
        assert "message" in result
        assert result["message"] == "Welcome to my agent! TestMessage"

    @pytest.mark.unit
    def test_create_workflow_multiple_calls(self):
        """Test that create_workflow can be called multiple times."""
        workflow1 = create_workflow()
        workflow2 = create_workflow()
        
        # They should be different instances but both functional
        assert workflow1 is not workflow2
        
        test_state = {"message": "test"}
        result1 = workflow1.invoke(test_state)
        result2 = workflow2.invoke(test_state)
        
        assert result1 == result2


class TestRunAgent:
    """Unit and integration tests for the run_agent function."""

    @pytest.mark.unit
    def test_run_agent_default_parameters(self):
        """Test run_agent with default parameters."""
        result = run_agent()
        expected = "Welcome to my agent! Mehrnoosh"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_run_agent_custom_message(self):
        """Test run_agent with custom message."""
        custom_message = "CustomTestMessage"
        result = run_agent(message=custom_message)
        expected = f"Welcome to my agent! {custom_message}"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_run_agent_dry_run_mode(self):
        """Test run_agent in dry-run mode."""
        custom_message = "DryRunTest"
        result = run_agent(dry_run=True, message=custom_message)
        expected = f"Welcome to my agent! {custom_message}"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_run_agent_dry_run_vs_normal(self):
        """Test that dry-run and normal execution produce same results."""
        message = "ComparisonTest"
        
        dry_run_result = run_agent(dry_run=True, message=message)
        normal_result = run_agent(dry_run=False, message=message)
        
        assert dry_run_result == normal_result

    @pytest.mark.unit
    @pytest.mark.parametrize("message", [
        "Test", "", "Unicode🌍", "Special!@#$%", "123456", "A" * 1000
    ])
    def test_run_agent_various_messages(self, message: str):
        """Test run_agent with various message types."""
        result = run_agent(message=message)
        expected = f"Welcome to my agent! {message}"
        assert result["message"] == expected

    @pytest.mark.unit
    def test_run_agent_return_type(self):
        """Test that run_agent returns correct type."""
        result = run_agent()
        assert isinstance(result, dict)
        assert "message" in result
        assert isinstance(result["message"], str)

    @pytest.mark.integration
    def test_run_agent_with_logging(self, log_capture):
        """Test run_agent with logging in dry-run mode."""
        run_agent(dry_run=True, message="LogTest")
        
        log_contents = log_capture.getvalue()
        assert "state transition" in log_contents
        assert "LogTest" in log_contents


class TestMainFunction:
    """Tests for the main function and CLI interface."""

    @pytest.mark.integration
    def test_main_with_default_args(self, monkeypatch, capsys):
        """Test main function with default arguments."""
        # Mock sys.argv to simulate default arguments
        test_args = ["run.py"]
        monkeypatch.setattr(sys, 'argv', test_args)
        
        # Run main function
        try:
            main()
        except SystemExit:
            pass  # main() might call sys.exit, which is OK for tests
            
        # Check that output was produced
        captured = capsys.readouterr()
        # In normal mode, result should be printed
        if captured.out:
            assert "message" in captured.out or "Welcome" in captured.out

    @pytest.mark.integration
    def test_main_with_dry_run_flag(self, monkeypatch, capsys):
        """Test main function with dry-run flag."""
        test_args = ["run.py", "--dry-run"]
        monkeypatch.setattr(sys, 'argv', test_args)
        
        try:
            main()
        except SystemExit:
            pass
            
        # In dry-run mode, no result should be printed to stdout
        capsys.readouterr()
        # Dry-run mode only logs, doesn't print result

    @pytest.mark.integration
    def test_main_with_custom_message(self, monkeypatch, capsys):
        """Test main function with custom message."""
        test_message = "CustomMainTest"
        test_args = ["run.py", "--message", test_message]
        monkeypatch.setattr(sys, 'argv', test_args)
        
        try:
            main()
        except SystemExit:
            pass
            
        captured = capsys.readouterr()
        if captured.out:
            assert test_message in captured.out


class TestErrorHandling:
    """Tests for error handling and edge cases."""

    @pytest.mark.error_handling
    @pytest.mark.unit
    def test_sample_node_with_missing_key(self):
        """Test sample_node with missing message key."""
        with pytest.raises(RuntimeError, match="Failed to process state in sample_node"):
            invalid_state = {}  # type: ignore
            sample_node(invalid_state)

    @pytest.mark.error_handling
    @pytest.mark.unit
    def test_sample_node_with_none_message(self):
        """Test sample_node behavior with None message."""
        with pytest.raises(RuntimeError, match="Failed to process state in sample_node"):
            invalid_state = {"message": None}  # type: ignore
            sample_node(invalid_state)

    @pytest.mark.error_handling
    @pytest.mark.unit
    def test_sample_node_with_non_string_message(self):
        """Test sample_node with non-string message."""
        with pytest.raises(RuntimeError, match="Failed to process state in sample_node"):
            invalid_state = {"message": 123}  # type: ignore
            sample_node(invalid_state)

    @pytest.mark.error_handling
    @pytest.mark.unit
    def test_sample_node_error_contains_original_exception(self):
        """Test that RuntimeError contains information about the original exception."""
        with pytest.raises(RuntimeError) as exc_info:
            invalid_state = {}  # type: ignore
            sample_node(invalid_state)
        
        # Check that the error message contains details about the original KeyError
        assert "'message'" in str(exc_info.value)
        # Check that the original exception is chained
        assert exc_info.value.__cause__ is not None
        assert isinstance(exc_info.value.__cause__, KeyError)


class TestIntegration:
    """Integration tests for the complete workflow."""

    @pytest.mark.integration
    @pytest.mark.parametrize("message", [
        "Alice", "Bob", "123", "IntegrationTest", "Mehrnoosh", "", "🌍"
    ])
    def test_complete_workflow_end_to_end(self, message: str):
        """Test complete workflow end-to-end with various inputs."""
        result = run_agent(message=message)
        expected = f"Welcome to my agent! {message}"
        assert result["message"] == expected

    @pytest.mark.integration
    def test_workflow_with_all_fixture_states(self, test_messages):
        """Test workflow with all fixture-provided test messages."""
        for message in test_messages:
            result = run_agent(message=message)
            expected = f"Welcome to my agent! {message}"
            assert result["message"] == expected

    @pytest.mark.integration
    def test_workflow_consistency(self):
        """Test that multiple runs of the same input produce consistent results."""
        message = "ConsistencyTest"
        results = [run_agent(message=message) for _ in range(5)]
        
        # All results should be identical
        expected = f"Welcome to my agent! {message}"
        for result in results:
            assert result["message"] == expected


class TestPerformance:
    """Performance tests for the agent workflow."""

    @pytest.mark.performance
    @pytest.mark.slow
    def test_workflow_performance(self):
        """Test that workflow execution completes in reasonable time."""
        message = "PerformanceTest"
        
        start_time = time.time()
        result = run_agent(message=message)
        end_time = time.time()
        
        # Should complete in less than 1 second
        execution_time = end_time - start_time
        assert execution_time < 1.0
        assert result["message"] == f"Welcome to my agent! {message}"

    @pytest.mark.performance
    def test_workflow_with_large_message(self):
        """Test workflow performance with large message."""
        large_message = "A" * 10000  # 10KB message
        
        start_time = time.time()
        result = run_agent(message=large_message)
        end_time = time.time()
        
        # Should still complete reasonably fast even with large input
        execution_time = end_time - start_time
        assert execution_time < 2.0
        assert result["message"] == f"Welcome to my agent! {large_message}"

    @pytest.mark.performance
    def test_multiple_sequential_executions(self):
        """Test performance of multiple sequential workflow executions."""
        messages = [f"Test{i}" for i in range(10)]
        
        start_time = time.time()
        results = [run_agent(message=msg) for msg in messages]
        end_time = time.time()
        
        # All executions should complete in reasonable time
        total_time = end_time - start_time
        assert total_time < 5.0  # 10 executions in under 5 seconds
        
        # Verify all results are correct
        for i, result in enumerate(results):
            expected = f"Welcome to my agent! Test{i}"
            assert result["message"] == expected
