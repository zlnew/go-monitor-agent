CREATE TABLE IF NOT EXISTS logs (
  id BIGSERIAL PRIMARY KEY,

  timestamp TIMESTAMPTZ NOT NULL,

  level TEXT NOT NULL,
  source TEXT NOT NULL,
  action TEXT NOT NULL,

  trace_id UUID NOT NULL,

  job_id BIGINT NULL,
  server_id UUID NULL,
  application_id BIGINT NULL,
  deployment_id BIGINT NULL,

  message TEXT NOT NULL,
  context JSONB NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_logs_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE SET NULL,
  CONSTRAINT fk_logs_server FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE SET NULL,
  CONSTRAINT fk_logs_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE SET NULL,
  CONSTRAINT fk_logs_deployment FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_trace_id ON logs(trace_id);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);
