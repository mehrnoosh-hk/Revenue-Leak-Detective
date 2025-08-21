"""
Agent runner module for Revenue Leak Detective.

This module provides the main entry point for running the agent workflow.
It supports dry-run mode for testing and validation purposes.
"""

import argparse
import logging
import sys
from typing import TypedDict

from langgraph.graph import StateGraph, START, END


# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)

logger = logging.getLogger(__name__)


class State(TypedDict):
    """State type definition for the agent workflow."""

    message: str


def sample_node(state: State) -> State:
    """
    Sample node that processes the state message.

    Args:
        state: Current state with message

    Returns:
        Updated state with processed message
    Raises:
        RuntimeError: If processing fails
    """
    try:
        return {"message": "Welcome to my agent! " + state["message"]}
    except Exception as e:
        logger.error(f"Error processing state in sample_node: {e}")
        raise RuntimeError(f"Failed to process state in sample_node: {e}") from e


def create_workflow():
    """
    Create and configure the agent workflow.

    Returns:
        Compiled StateGraph workflow
    Raises:
        RuntimeError: If workflow creation fails
    """
    try:
        workflow = StateGraph(State)
        workflow.add_node("sample", sample_node)
        workflow.add_edge(START, "sample")
        workflow.add_edge("sample", END)

        return workflow.compile()

    except Exception as e:
        logger.error(f"Failed to create workflow: {e}")
        raise RuntimeError(f"Workflow creation failed: {e}") from e


def run_agent(dry_run: bool = False, message: str = "Mehrnoosh") -> dict:
    """
    Run the agent workflow.

    Args:
        dry_run: If True, only log state transitions without executing
        message: Initial message for the workflow

    Returns:
        Result dictionary from the workflow execution
    """
    initial_state: State = {"message": message}

    if dry_run:
        logger.info(
            "state transition: %s -> %s",
            initial_state,
            {"message": f"Welcome to my agent! {message}"},
        )
        return {"message": f"Welcome to my agent! {message}"}

    graph = create_workflow()
    result = graph.invoke(initial_state)
    return result


def main() -> None:
    """Main entry point for the agent runner."""
    parser = argparse.ArgumentParser(
        description="Revenue Leak Detective Agent Runner",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )

    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Log state transitions without executing the workflow",
    )

    parser.add_argument(
        "--message",
        type=str,
        default="Mehrnoosh",
        help="Initial message for the workflow (default: %(default)s)",
    )

    args = parser.parse_args()

    try:
        result = run_agent(dry_run=args.dry_run, message=args.message)

        if not args.dry_run:
            print(result)

    except Exception as e:
        logger.error("Failed to run agent: %s", e)
        sys.exit(1)


if __name__ == "__main__":
    main()
