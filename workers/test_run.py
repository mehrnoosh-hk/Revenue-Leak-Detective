import pytest
from unittest.mock import patch

# Import the components from run.py
from run import State, sample_node, workflow, graph


class TestState:
    """Test cases for the State TypedDict."""

    def test_state_creation(self):
        """Test that State can be created with required fields."""
        state: State = {"message": "test message"}
        assert state["message"] == "test message"

    def test_state_type_annotation(self):
        """Test that State has correct type annotations."""
        # Check if State is a TypedDict
        assert hasattr(State, "__annotations__")
        assert State.__annotations__ == {"message": str}


class TestSampleNode:
    """Test cases for the sample_node function."""

    def test_sample_node_basic_functionality(self):
        """Test that sample_node processes state correctly."""
        input_state: State = {"message": "John"}
        result = sample_node(input_state)

        assert result["message"] == "Welcome to my agent! John"
        assert isinstance(result, dict)
        # Ensure the original state is not modified
        assert input_state["message"] == "John"

    def test_sample_node_empty_message(self):
        """Test sample_node with empty message."""
        input_state: State = {"message": ""}
        result = sample_node(input_state)

        assert result["message"] == "Welcome to my agent! "

    def test_sample_node_special_characters(self):
        """Test sample_node with special characters in message."""
        input_state: State = {"message": "Test@123!"}
        result = sample_node(input_state)

        assert result["message"] == "Welcome to my agent! Test@123!"

    def test_sample_node_long_message(self):
        """Test sample_node with a long message."""
        long_message = "A" * 1000
        input_state: State = {"message": long_message}
        result = sample_node(input_state)

        assert result["message"] == f"Welcome to my agent! {long_message}"
        assert len(result["message"]) == len("Welcome to my agent! ") + 1000


class TestWorkflow:
    """Test cases for the workflow configuration."""

    def test_workflow_has_nodes(self):
        """Test that workflow has the expected nodes."""
        # Access the workflow's internal structure
        nodes = workflow.nodes
        assert "sample" in nodes

    def test_workflow_compilation(self):
        """Test that workflow compiles successfully."""
        compiled_graph = workflow.compile()
        assert compiled_graph is not None


class TestGraph:
    """Test cases for the compiled graph."""

    def test_graph_invoke_basic(self):
        """Test basic graph invocation."""
        initial_state: State = {"message": "TestUser"}
        result = graph.invoke(initial_state)

        assert "message" in result
        assert result["message"] == "Welcome to my agent! TestUser"

    def test_graph_invoke_different_inputs(self):
        """Test graph with different input messages."""
        test_cases = [
            {"input": "Alice", "expected": "Welcome to my agent! Alice"},
            {"input": "Bob", "expected": "Welcome to my agent! Bob"},
            {"input": "123", "expected": "Welcome to my agent! 123"},
        ]

        for case in test_cases:
            initial_state: State = {"message": case["input"]}
            result = graph.invoke(initial_state)
            assert result["message"] == case["expected"]

    def test_graph_invoke_preserves_state_structure(self):
        """Test that graph invocation preserves the state structure."""
        initial_state: State = {"message": "StructureTest"}
        result = graph.invoke(initial_state)

        # Ensure result is still a valid State
        assert isinstance(result, dict)
        assert "message" in result
        assert isinstance(result["message"], str)


class TestMainExecution:
    """Test cases for the main execution block."""

    @patch("builtins.print")
    @patch("run.graph")
    def test_main_execution(self, mock_graph, mock_print):
        """Test the main execution block."""
        # Mock the graph.invoke method
        mock_result = {"message": "Welcome to my agent! Mehrnoosh"}
        mock_graph.invoke.return_value = mock_result

        # Verify graph.invoke was called with correct initial state
        expected_initial_state = {"message": "Mehrnoosh"}
        mock_graph.invoke.assert_called_with(expected_initial_state)

        # Verify print was called with the result
        mock_print.assert_called_with(mock_result)


class TestIntegration:
    """Integration tests for the complete workflow."""

    def test_end_to_end_workflow(self):
        """Test the complete workflow from start to finish."""
        # Test the actual workflow as it would run
        initial_state: State = {"message": "IntegrationTest"}
        result = graph.invoke(initial_state)

        # Verify the complete transformation
        assert result["message"] == "Welcome to my agent! IntegrationTest"

    def test_workflow_with_original_example(self):
        """Test with the same input as in the original main block."""
        initial_state: State = {"message": "Mehrnoosh"}
        result = graph.invoke(initial_state)

        assert result["message"] == "Welcome to my agent! Mehrnoosh"


if __name__ == "__main__":
    # Run the tests
    pytest.main([__file__])
