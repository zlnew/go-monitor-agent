package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"horizonx-server/internal/core/metrics/collector/os"
	"horizonx-server/internal/domain"
)

func (a *Agent) handleInitCommand(ctx context.Context, cmd domain.WsAgentCommand) {
	jobID := cmd.Payload.JobID

	go func() {
		a.log.Info("starting job", "job_id", jobID)

		startURL := fmt.Sprintf("%s/agent/jobs/%d/start", a.cfg.AgentTargetAPIURL, jobID)
		if err := a.httpPost(startURL, nil); err != nil {
			a.log.Error("failed to start job", "job_id", jobID, "error", err)
			return
		}
		osInfoCollector := os.NewCollector(a.log)
		osInfo, err := osInfoCollector.Collect(ctx)
		if err != nil {
			a.log.Error("failed to collect os info", "job_id", jobID, "error", err.Error())
			return
		}

		finishURL := fmt.Sprintf("%s/agent/jobs/%d/finish", a.cfg.AgentTargetAPIURL, jobID)
		outputLog, err := json.Marshal(osInfo)
		if err != nil {
			a.log.Error("failed to marshal os info", "job_id", jobID, "error", err)
			return
		}
		payload := &domain.JobFinishRequest{
			Status:    domain.JobSuccess,
			OutputLog: string(outputLog),
		}
		if err := a.httpPost(finishURL, payload); err != nil {
			a.log.Error("failed to finish job", "job_id", jobID, "error", err)
			return
		}

		a.log.Info("job finished successfully", "job_id", jobID)
	}()
}

func (a *Agent) httpPost(url string, body any) error {
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.cfg.AgentServerID.String()+"."+a.cfg.AgentServerAPIToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("http error: %s", resp.Status)
	}
	return nil
}
