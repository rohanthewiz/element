// Package main provides data models for the DAG editor application.
package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// NodeType represents the type of a DAG node/task
type NodeType string

const (
	NodeTypeK8sPodOperator   NodeType = "k8s_pod_operator"
	NodeTypeExternalSensor   NodeType = "external_sensor"
	NodeTypePythonOperator   NodeType = "python_operator"
	NodeTypeBashOperator     NodeType = "bash_operator"
)

// Node represents a task in the DAG
type Node struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        NodeType          `json:"type"`
	X           float64           `json:"x"`           // Canvas X position
	Y           float64           `json:"y"`           // Canvas Y position
	JobConfig   string            `json:"job_config"`  // For k8s_pod_operator
	EnvVars     map[string]string `json:"env_vars"`
	ImageVersion string           `json:"image_version"`

	// For external sensors
	ExternalDagID  string `json:"external_dag_id"`
	ExternalTaskID string `json:"external_task_id"`

	// Resource configuration
	CPURequest    string `json:"cpu_request"`
	MemoryRequest string `json:"memory_request"`
	CPULimit      string `json:"cpu_limit"`
	MemoryLimit   string `json:"memory_limit"`

	// Timeout
	StartupTimeoutSeconds int `json:"startup_timeout_seconds"`

	// Custom Python code (for python_operator or custom logic)
	PythonCode string `json:"python_code"`
}

// Edge represents a dependency between two nodes
type Edge struct {
	ID     string `json:"id"`
	From   string `json:"from"`   // Source node ID
	To     string `json:"to"`     // Target node ID
}

// DAGDefaultArgs represents default arguments for a DAG
type DAGDefaultArgs struct {
	Owner            string `json:"owner"`
	StartDate        string `json:"start_date"`
	DependOnPast     bool   `json:"depend_on_past"`
	EmailOnFailure   bool   `json:"email_on_failure"`
	EmailOnRetry     bool   `json:"email_on_retry"`
	Email            string `json:"email"`
	Retries          int    `json:"retries"`
	RetryDelayMins   int    `json:"retry_delay_mins"`
	ExecutionTimeoutHours int `json:"execution_timeout_hours"`
	IsPausedOnCreation bool `json:"is_paused_on_creation"`
}

// DAG represents an Airflow-style DAG
type DAG struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Schedule     string         `json:"schedule"`     // cron expression or None
	Concurrency  int            `json:"concurrency"`
	Tags         []string       `json:"tags"`
	DefaultArgs  DAGDefaultArgs `json:"default_args"`
	Nodes        []Node         `json:"nodes"`
	Edges        []Edge         `json:"edges"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Store provides thread-safe storage for DAGs
type Store struct {
	mu   sync.RWMutex
	dags map[string]*DAG
}

// NewStore creates a new storage instance
func NewStore() *Store {
	return &Store{
		dags: make(map[string]*DAG),
	}
}

// GetDAG retrieves a DAG by ID
func (s *Store) GetDAG(id string) (*DAG, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dag, ok := s.dags[id]
	return dag, ok
}

// ListDAGs returns all DAGs
func (s *Store) ListDAGs() []*DAG {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dags := make([]*DAG, 0, len(s.dags))
	for _, dag := range s.dags {
		dags = append(dags, dag)
	}
	return dags
}

// SaveDAG saves or updates a DAG
func (s *Store) SaveDAG(dag *DAG) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dag.UpdatedAt = time.Now()
	if dag.CreatedAt.IsZero() {
		dag.CreatedAt = time.Now()
	}
	s.dags[dag.ID] = dag
}

// DeleteDAG removes a DAG by ID
func (s *Store) DeleteDAG(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dags[id]; ok {
		delete(s.dags, id)
		return true
	}
	return false
}

// GeneratePythonCode generates Airflow Python code for a DAG
func (dag *DAG) GeneratePythonCode(imageRepoHost, imageRepoName string) string {
	var sb strings.Builder

	// Imports
	sb.WriteString(`from datetime import timedelta
from airflow import DAG
from airflow.providers.cncf.kubernetes.operators.pod import KubernetesPodOperator
from airflow.sensors.external_task import ExternalTaskSensor
from airflow.exceptions import AirflowSkipException
from airflow.models import DagRun
from airflow.utils.session import provide_session
from kubernetes.client import V1LocalObjectReference, V1ResourceRequirements

`)

	// Configuration variables
	sb.WriteString(fmt.Sprintf(`# Configuration
image_repo_host = "%s"
image_repo_name = "%s"
namespace = "airflow"
platform_env = "prod"
cpu_req_limit = "4"

`, imageRepoHost, imageRepoName))

	// Helper functions
	sb.WriteString(`# ===== HELPERS =====
def build_image_name(image_version):
    return f"{image_repo_host}/{image_repo_name}:{image_version}"


@provide_session
def most_recent_successful_dag(dt, external_dag_id, session=None, **_):
    search_start_date = dt - timedelta(hours=23)
    search_end_date = dt + timedelta(hours=8)

    successful_runs = DagRun.find(
        dag_id=external_dag_id,
        state="success",
        execution_start_date=search_start_date,
        execution_end_date=search_end_date,
        session=session
    )

    if successful_runs:
        successful_runs.sort(key=lambda x: x.execution_date, reverse=True)
        return successful_runs[0].execution_date
    else:
        raise AirflowSkipException(f"No successful run found for DAG dependency: {external_dag_id}"
                                   f" between {search_start_date} and {search_end_date}")


def create_k8s_pod_operator(
    dag: DAG,
    task_id,
    job_config,
    image_version,
    env_vars=None,
    startup_timeout_seconds=None,
    resource_requests=None,
    resource_limits=None,
):
    if resource_requests is None:
        resource_requests = {"cpu": "1", "memory": "8Gi"}
    if resource_limits is None:
        resource_limits = {"cpu": cpu_req_limit, "memory": "28Gi"}

    return KubernetesPodOperator(
        name=f"{task_id}-{platform_env}-pod",
        namespace=namespace,
        env_vars=env_vars,
        image=build_image_name(image_version),
        image_pull_secrets=[V1LocalObjectReference(name="image-secret")],
        image_pull_policy="IfNotPresent",
        arguments=[job_config],
        container_resources=V1ResourceRequirements(
            requests=resource_requests,
            limits=resource_limits,
        ),
        task_id=task_id,
        startup_timeout_seconds=startup_timeout_seconds,
        hostnetwork=False,
        in_cluster=False,
        is_delete_operator_pod=True,
        termination_grace_period=60,
        dag=dag,
    )


`)

	// DAG definition
	schedule := dag.Schedule
	if schedule == "" || schedule == "None" {
		schedule = "None"
	} else {
		schedule = fmt.Sprintf("'%s'", schedule)
	}

	tagsStr := "[]"
	if len(dag.Tags) > 0 {
		tags := make([]string, len(dag.Tags))
		for i, t := range dag.Tags {
			tags[i] = fmt.Sprintf("'%s'", t)
		}
		tagsStr = fmt.Sprintf("[%s]", strings.Join(tags, ", "))
	}

	sb.WriteString(fmt.Sprintf(`#----- DAG -----
dag_id = '%s'
dag_description = "%s"
dag_schedule = %s
dag_concurrency = %d
dag_tags = %s

dag_default_args = {
    'owner': '%s',
    'start_date': '%s',
    'depend_on_past': %s,
    'email_on_failure': %s,
    'email_on_retry': %s,
    'email': ['%s'],
    'retries': %d,
    'retry_delay': timedelta(minutes=%d),
    'execution_timeout': timedelta(hours=%d),
    'is_paused_upon_creation': %s
}

the_dag = DAG(
    dag_id=dag_id,
    default_args=dag_default_args,
    schedule=dag_schedule,
    description=dag_description,
    catchup=False,
    concurrency=dag_concurrency,
    tags=dag_tags
)
#----- END DAG -----

`,
		dag.ID,
		dag.Description,
		schedule,
		dag.Concurrency,
		tagsStr,
		dag.DefaultArgs.Owner,
		dag.DefaultArgs.StartDate,
		pythonBool(dag.DefaultArgs.DependOnPast),
		pythonBool(dag.DefaultArgs.EmailOnFailure),
		pythonBool(dag.DefaultArgs.EmailOnRetry),
		dag.DefaultArgs.Email,
		dag.DefaultArgs.Retries,
		dag.DefaultArgs.RetryDelayMins,
		dag.DefaultArgs.ExecutionTimeoutHours,
		pythonBool(dag.DefaultArgs.IsPausedOnCreation),
	))

	// Tasks
	sb.WriteString("#----- TASKS -----\n\n")

	// Build a map of node IDs to variable names
	nodeVarNames := make(map[string]string)
	for _, node := range dag.Nodes {
		varName := toSnakeCase(node.Name)
		nodeVarNames[node.ID] = varName
	}

	for _, node := range dag.Nodes {
		varName := nodeVarNames[node.ID]

		switch node.Type {
		case NodeTypeExternalSensor:
			sb.WriteString(fmt.Sprintf(`%s = ExternalTaskSensor(
    task_id='%s',
    external_dag_id='%s',
    external_task_id='%s',
    execution_date_fn=lambda dt, **kwargs: most_recent_successful_dag(dt, '%s'),
    mode='reschedule',
    timeout=14400,
    poke_interval=60,
    dag=the_dag,
)

`, varName, node.Name, node.ExternalDagID, node.ExternalTaskID, node.ExternalDagID))

		case NodeTypeK8sPodOperator:
			envVarsStr := "None"
			if len(node.EnvVars) > 0 {
				envPairs := make([]string, 0, len(node.EnvVars))
				for k, v := range node.EnvVars {
					envPairs = append(envPairs, fmt.Sprintf("'%s': '%s'", k, v))
				}
				envVarsStr = fmt.Sprintf("{%s}", strings.Join(envPairs, ", "))
			}

			imageVersion := node.ImageVersion
			if imageVersion == "" {
				imageVersion = "{RECENT_PUSHED_VERSION}"
			}

			startupTimeout := node.StartupTimeoutSeconds
			if startupTimeout == 0 {
				startupTimeout = 3600
			}

			sb.WriteString(fmt.Sprintf(`%s = create_k8s_pod_operator(
    dag=the_dag,
    task_id='%s',
    job_config='%s',
    image_version='%s',
    env_vars=%s,
    startup_timeout_seconds=%d,
)

`, varName, node.Name, node.JobConfig, imageVersion, envVarsStr, startupTimeout))

		case NodeTypePythonOperator:
			if node.PythonCode != "" {
				sb.WriteString(fmt.Sprintf(`# Custom Python code for %s
%s

`, varName, node.PythonCode))
			}
		}
	}

	// Dependencies
	sb.WriteString("\n# Dependencies\n")
	for _, edge := range dag.Edges {
		fromVar := nodeVarNames[edge.From]
		toVar := nodeVarNames[edge.To]
		if fromVar != "" && toVar != "" {
			sb.WriteString(fmt.Sprintf("%s >> %s\n", fromVar, toVar))
		}
	}

	sb.WriteString("#----- END TASKS -----\n")

	return sb.String()
}

func pythonBool(b bool) string {
	if b {
		return "True"
	}
	return "False"
}

func toSnakeCase(s string) string {
	// Replace hyphens and spaces with underscores
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	// Remove any characters that aren't alphanumeric or underscore
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	return strings.ToLower(result.String())
}

// CreateSampleDAG creates a sample DAG for demonstration
func CreateSampleDAG() *DAG {
	return &DAG{
		ID:          "sample-dag",
		Name:        "Sample DAG",
		Description: "A sample DAG demonstrating the editor",
		Schedule:    "None",
		Concurrency: 1,
		Tags:        []string{"sample", "demo"},
		DefaultArgs: DAGDefaultArgs{
			Owner:                 "data-team",
			StartDate:             time.Now().Format("2006-01-02 15:04:05"),
			DependOnPast:          false,
			EmailOnFailure:        true,
			EmailOnRetry:          false,
			Email:                 "team@example.com",
			Retries:               0,
			RetryDelayMins:        10,
			ExecutionTimeoutHours: 60,
			IsPausedOnCreation:    false,
		},
		Nodes: []Node{
			{
				ID:             "node-1",
				Name:           "external-sensor",
				Type:           NodeTypeExternalSensor,
				X:              100,
				Y:              100,
				ExternalDagID:  "upstream-dag",
				ExternalTaskID: "upstream-task",
			},
			{
				ID:           "node-2",
				Name:         "start-status",
				Type:         NodeTypeK8sPodOperator,
				X:            300,
				Y:            100,
				JobConfig:    "common-configs/job_flow_status.json",
				ImageVersion: "{RECENT_PUSHED_VERSION}",
				EnvVars: map[string]string{
					"JOB_FLOW_STATUS": "Started",
					"ENV":             "prod",
				},
				StartupTimeoutSeconds: 3600,
			},
			{
				ID:           "node-3",
				Name:         "process-data",
				Type:         NodeTypeK8sPodOperator,
				X:            500,
				Y:            100,
				JobConfig:    "datamart-configs/process_data.json",
				ImageVersion: "{RECENT_PUSHED_VERSION}",
				EnvVars: map[string]string{
					"PARALLEL_JOBS": "true",
				},
				StartupTimeoutSeconds: 3600,
			},
			{
				ID:           "node-4",
				Name:         "finish-status",
				Type:         NodeTypeK8sPodOperator,
				X:            700,
				Y:            100,
				JobConfig:    "common-configs/job_flow_status.json",
				ImageVersion: "{RECENT_PUSHED_VERSION}",
				EnvVars: map[string]string{
					"JOB_FLOW_STATUS": "Finished",
					"ENV":             "prod",
				},
				StartupTimeoutSeconds: 3600,
			},
		},
		Edges: []Edge{
			{ID: "edge-1", From: "node-1", To: "node-2"},
			{ID: "edge-2", From: "node-2", To: "node-3"},
			{ID: "edge-3", From: "node-3", To: "node-4"},
		},
	}
}
