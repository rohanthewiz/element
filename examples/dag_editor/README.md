# DAG Editor

A visual DAG (Directed Acyclic Graph) editor for creating Airflow-style workflows. Built with the Element library and rweb server.

## Features

- **Visual Node Editor**: Drag-and-drop interface for creating and positioning task nodes
- **Node Connections**: Click and drag to create dependencies between nodes
- **Multiple Node Types**:
  - K8s Pod Operator: Kubernetes-based task execution
  - External Sensor: Dependencies on other DAGs
  - Python Operator: Custom Python code execution
- **Properties Panel**: Edit node and DAG properties
- **Monaco Editor Integration**: Code editing with syntax highlighting
- **Python Export**: Generate Airflow-compatible Python DAG files
- **Live Dependency View**: See `node1 >> node2` style dependencies in real-time

## Running the Example

```bash
cd examples/dag_editor
go run .
```

Then visit: http://localhost:8080

## Usage

### Creating a DAG

1. Click "+ New DAG" on the home page
2. Enter a name for your DAG
3. You'll be taken to the visual editor

### Adding Nodes

- Click "+ K8s Task" to add a Kubernetes Pod Operator task
- Click "+ External Sensor" to add a dependency on another DAG

### Connecting Nodes

1. Click the "C" (Connect) button in the toolbar
2. Click and drag from a node's output port (right side) to another node's input port (left side)
3. Release to create the connection

### Editing Properties

- Click on any node to edit its properties in the right panel
- Switch to the "DAG" tab to edit overall DAG settings
- Switch to the "Code" tab to see generated Python code

### Deleting Nodes/Connections

1. Click the "D" (Delete) button in the toolbar
2. Click on a node or connection to delete it

### Exporting

Click "Export" to download the generated Airflow Python DAG file.

## Generated Code Format

The editor generates Python code compatible with Airflow, including:

- DAG definition with configurable schedule, concurrency, and default arguments
- K8s Pod Operator tasks with environment variables and resource configs
- External Task Sensors for cross-DAG dependencies
- Dependency chains in the format:

```python
external_sensor >> start_status
start_status >> process_data
process_data >> finish_status
```

## Architecture

```
dag_editor/
├── main.go      # Server, routes, and Element-based UI rendering
├── models.go    # Data models (DAG, Node, Edge) and Python code generation
├── styles.go    # CSS styles
├── scripts.go   # JavaScript for the interactive editor
└── README.md    # This file
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/dags` | List all DAGs |
| GET | `/api/dags/{id}` | Get a specific DAG |
| POST | `/api/dags` | Create a new DAG |
| PUT | `/api/dags/{id}` | Update a DAG |
| DELETE | `/api/dags/{id}` | Delete a DAG |
| GET | `/api/dags/{id}/export` | Export DAG as Python file |

## Keyboard Shortcuts (Planned)

- `S`: Switch to Select mode
- `C`: Switch to Connect mode
- `D`: Switch to Delete mode
- `Ctrl+S`: Save DAG

## Extending

### Adding New Node Types

1. Add the type to `NodeType` constants in `models.go`
2. Update `GeneratePythonCode()` to handle the new type
3. Add UI button in `main.go` toolbar
4. Update `showNodeProperties()` in `scripts.go` for the properties form

### Adding New Properties

1. Add field to the appropriate struct in `models.go`
2. Update the JSON serialization (automatic with struct tags)
3. Add form field in `showNodeProperties()` in `scripts.go`
4. Update `GeneratePythonCode()` if the property affects code generation
