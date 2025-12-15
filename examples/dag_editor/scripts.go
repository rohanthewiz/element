package main

func homePageScript() string {
	return `
async function createNewDAG() {
    const name = prompt('Enter DAG name:');
    if (!name) return;

    const dag = {
        id: 'dag-' + Date.now(),
        name: name,
        description: 'New DAG',
        schedule: 'None',
        concurrency: 1,
        tags: [],
        default_args: {
            owner: 'data-team',
            start_date: new Date().toISOString().split('T')[0] + ' 00:00:00',
            depend_on_past: false,
            email_on_failure: true,
            email_on_retry: false,
            email: 'team@example.com',
            retries: 0,
            retry_delay_mins: 10,
            execution_timeout_hours: 60,
            is_paused_on_creation: false
        },
        nodes: [],
        edges: []
    };

    try {
        const response = await fetch('/api/dags', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(dag)
        });

        if (response.ok) {
            const savedDag = await response.json();
            window.location.href = '/editor/' + savedDag.id;
        }
    } catch (err) {
        alert('Error creating DAG: ' + err.message);
    }
}

async function deleteDAG(id) {
    if (!confirm('Are you sure you want to delete this DAG?')) return;

    try {
        const response = await fetch('/api/dags/' + id, { method: 'DELETE' });
        if (response.ok) {
            window.location.reload();
        }
    } catch (err) {
        alert('Error deleting DAG: ' + err.message);
    }
}

async function exportDAG(id) {
    window.open('/api/dags/' + id + '/export', '_blank');
}
`
}

func editorScript() string {
	return `
// Global state
let dag = window.dagData;
let selectedNode = null;
let currentTool = 'select';
let isDragging = false;
let dragOffset = { x: 0, y: 0 };
let isConnecting = false;
let connectionStart = null;
let tempLine = null;
let monacoEditor = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initCanvas();
    initToolbar();
    initPanelTabs();
    initMonaco();
    renderNodes();
    renderConnections();
    updateDependenciesOutput();
});

// Initialize canvas
function initCanvas() {
    const canvas = document.getElementById('canvas');

    canvas.addEventListener('click', (e) => {
        if (e.target === canvas && currentTool === 'select') {
            selectNode(null);
        }
    });
}

// Initialize toolbar
function initToolbar() {
    document.querySelectorAll('.tool-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.tool-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            currentTool = btn.dataset.tool;
            document.getElementById('canvas').style.cursor =
                currentTool === 'connect' ? 'crosshair' :
                currentTool === 'delete' ? 'not-allowed' : 'default';
        });
    });
}

// Initialize panel tabs
function initPanelTabs() {
    document.querySelectorAll('.panel-tab').forEach(tab => {
        tab.addEventListener('click', () => {
            document.querySelectorAll('.panel-tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.panel-section').forEach(s => s.classList.remove('active'));
            tab.classList.add('active');
            document.getElementById(tab.dataset.panel + '-panel').classList.add('active');

            if (tab.dataset.panel === 'code') {
                updateMonacoContent();
            }
        });
    });
}

// Initialize Monaco Editor
function initMonaco() {
    require.config({ paths: { vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs' } });
    require(['vs/editor/editor.main'], function () {
        monacoEditor = monaco.editor.create(document.getElementById('monaco-editor'), {
            value: '',
            language: 'python',
            theme: 'vs-dark',
            minimap: { enabled: false },
            fontSize: 12,
            automaticLayout: true,
            readOnly: false,
            scrollBeyondLastLine: false
        });

        // Listen for changes when editing node code
        monacoEditor.onDidChangeModelContent(() => {
            if (selectedNode && selectedNode.type === 'python_operator') {
                selectedNode.python_code = monacoEditor.getValue();
            }
        });
    });
}

function updateMonacoContent() {
    if (!monacoEditor) return;

    if (selectedNode && selectedNode.python_code) {
        monacoEditor.setValue(selectedNode.python_code);
    } else {
        // Generate preview code
        const code = generatePythonPreview();
        monacoEditor.setValue(code);
    }
}

function generatePythonPreview() {
    let code = '# Generated Airflow DAG Preview\n\n';
    code += '# Dependencies:\n';
    dag.edges.forEach(edge => {
        const fromNode = dag.nodes.find(n => n.id === edge.from);
        const toNode = dag.nodes.find(n => n.id === edge.to);
        if (fromNode && toNode) {
            code += toSnakeCase(fromNode.name) + ' >> ' + toSnakeCase(toNode.name) + '\n';
        }
    });
    return code;
}

function toSnakeCase(str) {
    return str.replace(/-/g, '_').replace(/\s/g, '_').toLowerCase();
}

// Add node
function addNode(type) {
    const id = 'node-' + Date.now();
    const canvas = document.getElementById('canvas');
    const rect = canvas.getBoundingClientRect();

    const node = {
        id: id,
        name: type === 'external_sensor' ? 'external-sensor-' + dag.nodes.length : 'task-' + dag.nodes.length,
        type: type,
        x: 100 + (dag.nodes.length % 4) * 200,
        y: 100 + Math.floor(dag.nodes.length / 4) * 150,
        job_config: '',
        env_vars: {},
        image_version: '{RECENT_PUSHED_VERSION}',
        external_dag_id: type === 'external_sensor' ? 'upstream-dag' : '',
        external_task_id: type === 'external_sensor' ? 'upstream-task' : '',
        startup_timeout_seconds: 3600
    };

    dag.nodes.push(node);
    renderNode(node);
    selectNode(node);
    updateDependenciesOutput();
}

// Render all nodes
function renderNodes() {
    dag.nodes.forEach(node => renderNode(node));
}

// Render single node
function renderNode(node) {
    const canvas = document.getElementById('canvas');
    const existingEl = document.getElementById(node.id);
    if (existingEl) existingEl.remove();

    const el = document.createElement('div');
    el.id = node.id;
    el.className = 'node' + (node.type === 'external_sensor' ? ' sensor' : '');
    el.style.left = node.x + 'px';
    el.style.top = node.y + 'px';

    const typeLabel = node.type === 'external_sensor' ? 'External Sensor' :
                      node.type === 'k8s_pod_operator' ? 'K8s Pod' : node.type;

    el.innerHTML = ` + "`" + `
        <div class="node-header">
            <span>${node.name}</span>
            <span class="node-type">${typeLabel}</span>
        </div>
        <div class="node-body">
            ${node.type === 'external_sensor'
                ? ` + "`" + `DAG: ${node.external_dag_id || 'Not set'}` + "`" + `
                : ` + "`" + `Config: ${node.job_config || 'Not set'}` + "`" + `}
        </div>
        <div class="node-ports">
            <div class="port port-in" data-port="in"></div>
            <div class="port port-out" data-port="out"></div>
        </div>
    ` + "`" + `;

    // Node events
    el.addEventListener('mousedown', (e) => {
        if (e.target.classList.contains('port')) return;

        if (currentTool === 'delete') {
            deleteNode(node.id);
            return;
        }

        selectNode(node);
        isDragging = true;
        dragOffset = {
            x: e.clientX - node.x,
            y: e.clientY - node.y
        };
    });

    // Port events for connections
    el.querySelectorAll('.port').forEach(port => {
        port.addEventListener('mousedown', (e) => {
            e.stopPropagation();
            if (currentTool !== 'connect') return;

            isConnecting = true;
            connectionStart = {
                nodeId: node.id,
                port: port.dataset.port,
                x: node.x + (port.dataset.port === 'out' ? el.offsetWidth : 0),
                y: node.y + el.offsetHeight / 2
            };

            // Create temp line
            const svg = document.getElementById('connections-svg');
            tempLine = document.createElementNS('http://www.w3.org/2000/svg', 'path');
            tempLine.classList.add('temp-connection');
            svg.appendChild(tempLine);
        });

        port.addEventListener('mouseup', (e) => {
            if (!isConnecting || !connectionStart) return;
            if (connectionStart.nodeId === node.id) return;

            // Create edge
            if (connectionStart.port === 'out' && port.dataset.port === 'in') {
                addEdge(connectionStart.nodeId, node.id);
            } else if (connectionStart.port === 'in' && port.dataset.port === 'out') {
                addEdge(node.id, connectionStart.nodeId);
            }
        });
    });

    canvas.appendChild(el);
}

// Mouse move handler
document.addEventListener('mousemove', (e) => {
    if (isDragging && selectedNode) {
        const canvas = document.getElementById('canvas');
        const rect = canvas.getBoundingClientRect();

        selectedNode.x = Math.max(0, e.clientX - rect.left - dragOffset.x);
        selectedNode.y = Math.max(0, e.clientY - rect.top - dragOffset.y);

        const el = document.getElementById(selectedNode.id);
        el.style.left = selectedNode.x + 'px';
        el.style.top = selectedNode.y + 'px';

        renderConnections();
    }

    if (isConnecting && tempLine) {
        const canvas = document.getElementById('canvas');
        const rect = canvas.getBoundingClientRect();
        const endX = e.clientX - rect.left;
        const endY = e.clientY - rect.top;

        const d = createCurvePath(connectionStart.x, connectionStart.y, endX, endY);
        tempLine.setAttribute('d', d);
    }
});

// Mouse up handler
document.addEventListener('mouseup', () => {
    isDragging = false;

    if (isConnecting) {
        isConnecting = false;
        connectionStart = null;
        if (tempLine) {
            tempLine.remove();
            tempLine = null;
        }
    }
});

// Select node
function selectNode(node) {
    document.querySelectorAll('.node').forEach(el => el.classList.remove('selected'));
    selectedNode = node;

    if (node) {
        document.getElementById(node.id).classList.add('selected');
        showNodeProperties(node);
    } else {
        document.getElementById('node-form').classList.add('hidden');
        document.querySelector('.empty-selection').classList.remove('hidden');
    }
}

// Show node properties
function showNodeProperties(node) {
    const form = document.getElementById('node-form');
    form.classList.remove('hidden');
    document.querySelector('.empty-selection').classList.add('hidden');

    let html = ` + "`" + `
        <div class="form-group">
            <label>Name</label>
            <input type="text" class="form-input" value="${node.name}" onchange="updateNodeProperty('name', this.value)">
        </div>
        <div class="form-group">
            <label>Type</label>
            <select class="form-input" onchange="updateNodeProperty('type', this.value)">
                <option value="k8s_pod_operator" ${node.type === 'k8s_pod_operator' ? 'selected' : ''}>K8s Pod Operator</option>
                <option value="external_sensor" ${node.type === 'external_sensor' ? 'selected' : ''}>External Sensor</option>
                <option value="python_operator" ${node.type === 'python_operator' ? 'selected' : ''}>Python Operator</option>
            </select>
        </div>
    ` + "`" + `;

    if (node.type === 'external_sensor') {
        html += ` + "`" + `
            <div class="form-group">
                <label>External DAG ID</label>
                <input type="text" class="form-input" value="${node.external_dag_id || ''}"
                       onchange="updateNodeProperty('external_dag_id', this.value)">
            </div>
            <div class="form-group">
                <label>External Task ID</label>
                <input type="text" class="form-input" value="${node.external_task_id || ''}"
                       onchange="updateNodeProperty('external_task_id', this.value)">
            </div>
        ` + "`" + `;
    } else if (node.type === 'k8s_pod_operator') {
        html += ` + "`" + `
            <div class="form-group">
                <label>Job Config</label>
                <input type="text" class="form-input" value="${node.job_config || ''}"
                       onchange="updateNodeProperty('job_config', this.value)"
                       placeholder="path/to/config.json">
            </div>
            <div class="form-group">
                <label>Image Version</label>
                <input type="text" class="form-input" value="${node.image_version || ''}"
                       onchange="updateNodeProperty('image_version', this.value)">
            </div>
            <div class="form-group">
                <label>Startup Timeout (seconds)</label>
                <input type="number" class="form-input" value="${node.startup_timeout_seconds || 3600}"
                       onchange="updateNodeProperty('startup_timeout_seconds', parseInt(this.value))">
            </div>
            <div class="form-group">
                <label>Environment Variables</label>
                <div class="env-vars-list" id="env-vars-list">
                    ${renderEnvVars(node.env_vars || {})}
                </div>
                <button type="button" class="add-env-var" onclick="addEnvVar()">+ Add Variable</button>
            </div>
        ` + "`" + `;
    }

    form.innerHTML = html;
}

function renderEnvVars(envVars) {
    let html = '';
    for (const [key, value] of Object.entries(envVars)) {
        html += ` + "`" + `
            <div class="env-var-row">
                <input type="text" value="${key}" placeholder="KEY" onchange="updateEnvVarKey(this, '${key}')">
                <input type="text" value="${value}" placeholder="value" onchange="updateEnvVarValue('${key}', this.value)">
                <button type="button" onclick="removeEnvVar('${key}')">x</button>
            </div>
        ` + "`" + `;
    }
    return html;
}

function addEnvVar() {
    if (!selectedNode) return;
    if (!selectedNode.env_vars) selectedNode.env_vars = {};
    const key = 'NEW_VAR_' + Object.keys(selectedNode.env_vars).length;
    selectedNode.env_vars[key] = '';
    showNodeProperties(selectedNode);
}

function updateEnvVarKey(input, oldKey) {
    if (!selectedNode || !selectedNode.env_vars) return;
    const value = selectedNode.env_vars[oldKey];
    delete selectedNode.env_vars[oldKey];
    selectedNode.env_vars[input.value] = value;
}

function updateEnvVarValue(key, value) {
    if (!selectedNode || !selectedNode.env_vars) return;
    selectedNode.env_vars[key] = value;
}

function removeEnvVar(key) {
    if (!selectedNode || !selectedNode.env_vars) return;
    delete selectedNode.env_vars[key];
    showNodeProperties(selectedNode);
}

// Update node property
function updateNodeProperty(prop, value) {
    if (!selectedNode) return;
    selectedNode[prop] = value;
    renderNode(selectedNode);
    document.getElementById(selectedNode.id).classList.add('selected');
    updateDependenciesOutput();
}

// Delete node
function deleteNode(nodeId) {
    dag.nodes = dag.nodes.filter(n => n.id !== nodeId);
    dag.edges = dag.edges.filter(e => e.from !== nodeId && e.to !== nodeId);

    document.getElementById(nodeId)?.remove();
    renderConnections();
    selectNode(null);
    updateDependenciesOutput();
}

// Add edge
function addEdge(fromId, toId) {
    // Check if edge already exists
    const exists = dag.edges.some(e => e.from === fromId && e.to === toId);
    if (exists) return;

    dag.edges.push({
        id: 'edge-' + Date.now(),
        from: fromId,
        to: toId
    });

    renderConnections();
    updateDependenciesOutput();
}

// Render connections
function renderConnections() {
    const svg = document.getElementById('connections-svg');
    // Remove existing paths (except temp)
    svg.querySelectorAll('path:not(.temp-connection)').forEach(p => p.remove());

    dag.edges.forEach(edge => {
        const fromEl = document.getElementById(edge.from);
        const toEl = document.getElementById(edge.to);
        if (!fromEl || !toEl) return;

        const fromNode = dag.nodes.find(n => n.id === edge.from);
        const toNode = dag.nodes.find(n => n.id === edge.to);

        const fromX = fromNode.x + fromEl.offsetWidth;
        const fromY = fromNode.y + fromEl.offsetHeight / 2;
        const toX = toNode.x;
        const toY = toNode.y + toEl.offsetHeight / 2;

        const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        path.setAttribute('d', createCurvePath(fromX, fromY, toX, toY));
        path.dataset.edgeId = edge.id;

        path.addEventListener('click', () => {
            if (currentTool === 'delete') {
                dag.edges = dag.edges.filter(e => e.id !== edge.id);
                renderConnections();
                updateDependenciesOutput();
            }
        });
        path.style.pointerEvents = 'stroke';
        path.style.cursor = currentTool === 'delete' ? 'not-allowed' : 'default';

        svg.appendChild(path);
    });
}

function createCurvePath(x1, y1, x2, y2) {
    const midX = (x1 + x2) / 2;
    const cp1x = midX;
    const cp1y = y1;
    const cp2x = midX;
    const cp2y = y2;
    return ` + "`" + `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${x2} ${y2}` + "`" + `;
}

// Update dependencies output
function updateDependenciesOutput() {
    const output = document.getElementById('dependencies-output');
    let deps = '';

    dag.edges.forEach(edge => {
        const fromNode = dag.nodes.find(n => n.id === edge.from);
        const toNode = dag.nodes.find(n => n.id === edge.to);
        if (fromNode && toNode) {
            deps += toSnakeCase(fromNode.name) + ' >> ' + toSnakeCase(toNode.name) + '\n';
        }
    });

    output.textContent = deps || 'No dependencies defined';
}

// Save DAG
async function saveDAG() {
    // Update DAG properties from form
    const dagForm = document.getElementById('dag-form');
    if (dagForm) {
        dag.name = dagForm.querySelector('#dag-name')?.value || dag.name;
        dag.description = dagForm.querySelector('#dag-description')?.value || dag.description;
        dag.schedule = dagForm.querySelector('#dag-schedule')?.value || dag.schedule;
        dag.concurrency = parseInt(dagForm.querySelector('#dag-concurrency')?.value) || dag.concurrency;

        const tagsInput = dagForm.querySelector('#dag-tags')?.value || '';
        dag.tags = tagsInput.split(',').map(t => t.trim()).filter(t => t);

        dag.default_args.owner = dagForm.querySelector('#dag-owner')?.value || dag.default_args.owner;
        dag.default_args.email = dagForm.querySelector('#dag-email')?.value || dag.default_args.email;
    }

    try {
        const response = await fetch('/api/dags/' + dag.id, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(dag)
        });

        if (response.ok) {
            alert('DAG saved successfully!');
        } else {
            throw new Error('Failed to save');
        }
    } catch (err) {
        alert('Error saving DAG: ' + err.message);
    }
}

// Export DAG
function exportDAG() {
    window.open('/api/dags/' + dag.id + '/export', '_blank');
}
`
}
