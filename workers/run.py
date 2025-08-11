from typing import TypedDict
from langgraph.graph import StateGraph, START, END


class State(TypedDict):
    message: str


def sample_node(state: State) -> State:
    return {"message": "Welcome to my agent! " + state["message"]}


workflow = StateGraph(State)

workflow.add_node("sample", sample_node)

workflow.add_edge(START, "sample")

workflow.add_edge("sample", END)

graph = workflow.compile()

# Run the graph
if __name__ == "__main__":
    initial_state: State = {"message": "Mehrnoosh"}
    result = graph.invoke(
        initial_state
    )
    print(result)
