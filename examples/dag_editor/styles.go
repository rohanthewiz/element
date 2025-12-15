package main

func getStyles() string {
	return `
* { box-sizing: border-box; margin: 0; padding: 0; }

:root {
    --primary: #3498db;
    --primary-dark: #2980b9;
    --danger: #e74c3c;
    --success: #27ae60;
    --warning: #f39c12;
    --bg: #f5f7fa;
    --bg-dark: #2c3e50;
    --text: #333;
    --text-light: #666;
    --border: #ddd;
    --node-bg: #fff;
    --node-border: #3498db;
    --sensor-border: #9b59b6;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: var(--bg);
    color: var(--text);
    line-height: 1.5;
}

.container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 1rem;
}

/* Buttons */
.btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    text-decoration: none;
    transition: all 0.2s;
}

.btn:hover { opacity: 0.9; }
.btn-primary { background: var(--primary); color: white; }
.btn-secondary { background: #95a5a6; color: white; }
.btn-danger { background: var(--danger); color: white; }
.btn-sm { padding: 0.25rem 0.75rem; font-size: 0.8rem; }

/* Header */
header {
    background: var(--bg-dark);
    color: white;
    padding: 1rem 0;
}

.header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

header h1 {
    font-size: 1.5rem;
    font-weight: 600;
}

/* Main content */
main {
    padding: 2rem 0;
}

main h2 {
    margin-bottom: 1.5rem;
    font-size: 1.25rem;
    color: var(--text-light);
}

/* DAG Grid */
.dag-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
}

.dag-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    overflow: hidden;
    transition: transform 0.2s, box-shadow 0.2s;
}

.dag-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.dag-card-header {
    padding: 1rem;
    background: var(--bg);
    border-bottom: 1px solid var(--border);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.dag-card-header h3 {
    font-size: 1rem;
    font-weight: 600;
}

.dag-schedule {
    font-size: 0.75rem;
    color: var(--text-light);
    background: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
}

.dag-card-body {
    padding: 1rem;
}

.dag-card-body p {
    color: var(--text-light);
    font-size: 0.875rem;
    margin-bottom: 0.75rem;
}

.dag-stats {
    font-size: 0.8rem;
    color: var(--text-light);
    margin-bottom: 0.75rem;
}

.dag-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.tag {
    font-size: 0.7rem;
    padding: 0.2rem 0.5rem;
    background: var(--bg);
    border-radius: 3px;
    color: var(--text-light);
}

.dag-card-actions {
    padding: 0.75rem 1rem;
    background: var(--bg);
    display: flex;
    gap: 0.5rem;
}

.empty-state {
    text-align: center;
    padding: 3rem;
    color: var(--text-light);
}

/* Editor Layout */
.editor-app {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 1rem;
    background: var(--bg-dark);
    color: white;
}

.toolbar-left, .toolbar-center, .toolbar-right {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.dag-name {
    font-weight: 600;
    margin-left: 1rem;
}

.tool-btn {
    width: 32px;
    height: 32px;
    border: none;
    border-radius: 4px;
    background: rgba(255,255,255,0.1);
    color: white;
    cursor: pointer;
    font-weight: bold;
}

.tool-btn:hover {
    background: rgba(255,255,255,0.2);
}

.tool-btn.active {
    background: var(--primary);
}

.editor-main {
    display: flex;
    flex: 1;
    overflow: hidden;
}

/* Canvas */
.canvas-container {
    flex: 1;
    position: relative;
    background:
        linear-gradient(90deg, var(--border) 1px, transparent 1px),
        linear-gradient(var(--border) 1px, transparent 1px);
    background-size: 20px 20px;
    overflow: auto;
}

.canvas {
    position: absolute;
    top: 0;
    left: 0;
    min-width: 100%;
    min-height: 100%;
}

.connections-svg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
}

.connections-svg path {
    fill: none;
    stroke: #666;
    stroke-width: 2;
    marker-end: url(#arrowhead);
}

.connections-svg path.temp-connection {
    stroke: var(--primary);
    stroke-dasharray: 5,5;
}

/* Nodes */
.node {
    position: absolute;
    min-width: 180px;
    background: var(--node-bg);
    border: 2px solid var(--node-border);
    border-radius: 8px;
    cursor: move;
    user-select: none;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.node.sensor {
    border-color: var(--sensor-border);
}

.node.selected {
    box-shadow: 0 0 0 3px rgba(52,152,219,0.3);
}

.node-header {
    padding: 0.5rem 0.75rem;
    background: var(--node-border);
    color: white;
    font-size: 0.8rem;
    font-weight: 600;
    border-radius: 5px 5px 0 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.node.sensor .node-header {
    background: var(--sensor-border);
}

.node-type {
    font-size: 0.65rem;
    opacity: 0.8;
    text-transform: uppercase;
}

.node-body {
    padding: 0.5rem 0.75rem;
    font-size: 0.75rem;
    color: var(--text-light);
}

.node-ports {
    display: flex;
    justify-content: space-between;
    padding: 0 0.5rem 0.5rem;
}

.port {
    width: 12px;
    height: 12px;
    background: white;
    border: 2px solid var(--node-border);
    border-radius: 50%;
    cursor: crosshair;
}

.node.sensor .port {
    border-color: var(--sensor-border);
}

.port:hover {
    background: var(--primary);
    border-color: var(--primary);
}

.port-in { margin-left: -6px; }
.port-out { margin-right: -6px; }

/* Properties Panel */
.properties-panel {
    width: 350px;
    background: white;
    border-left: 1px solid var(--border);
    display: flex;
    flex-direction: column;
}

.panel-header {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border);
}

.panel-header h3 {
    font-size: 0.875rem;
    margin-bottom: 0.5rem;
}

.panel-tabs {
    display: flex;
    gap: 0.25rem;
}

.panel-tab {
    padding: 0.25rem 0.75rem;
    border: none;
    background: var(--bg);
    border-radius: 4px;
    font-size: 0.75rem;
    cursor: pointer;
}

.panel-tab.active {
    background: var(--primary);
    color: white;
}

.panel-content {
    flex: 1;
    overflow-y: auto;
}

.panel-section {
    display: none;
    padding: 1rem;
}

.panel-section.active {
    display: block;
}

.empty-selection {
    text-align: center;
    padding: 2rem;
    color: var(--text-light);
}

/* Forms */
.form-group {
    margin-bottom: 1rem;
}

.form-group label {
    display: block;
    font-size: 0.75rem;
    font-weight: 600;
    margin-bottom: 0.25rem;
    color: var(--text-light);
}

.form-input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border);
    border-radius: 4px;
    font-size: 0.875rem;
}

.form-input:focus {
    outline: none;
    border-color: var(--primary);
}

textarea.form-input {
    resize: vertical;
    min-height: 60px;
}

/* Monaco Editor */
.monaco-container {
    height: 400px;
    border: 1px solid var(--border);
    border-radius: 4px;
}

/* Dependencies Panel */
.dependencies-panel {
    padding: 0.75rem 1rem;
    background: var(--bg-dark);
    color: white;
    max-height: 150px;
    overflow-y: auto;
}

.dependencies-panel h4 {
    font-size: 0.75rem;
    margin-bottom: 0.5rem;
    opacity: 0.8;
}

.dependencies-panel pre {
    font-family: 'SF Mono', Consolas, monospace;
    font-size: 0.8rem;
    line-height: 1.4;
}

.hidden { display: none !important; }

/* Env vars editor */
.env-vars-list {
    margin-top: 0.5rem;
}

.env-var-row {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
    align-items: center;
}

.env-var-row input {
    flex: 1;
    padding: 0.25rem 0.5rem;
    border: 1px solid var(--border);
    border-radius: 3px;
    font-size: 0.75rem;
}

.env-var-row button {
    padding: 0.25rem 0.5rem;
    border: none;
    background: var(--danger);
    color: white;
    border-radius: 3px;
    cursor: pointer;
    font-size: 0.75rem;
}

.add-env-var {
    padding: 0.25rem 0.5rem;
    border: 1px dashed var(--border);
    background: none;
    border-radius: 3px;
    cursor: pointer;
    font-size: 0.75rem;
    color: var(--text-light);
    width: 100%;
}

.add-env-var:hover {
    border-color: var(--primary);
    color: var(--primary);
}
`
}
