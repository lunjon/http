package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
)

type SandboxHandler struct {

}

func (h *SandboxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`{"url": "%s", "method": "%s"}`, r.URL, r.Method)
	_, _ = w.Write([]byte(body))
}

func startSandbox(cmd *cobra.Command, args []string) {
	portStr := cmd.Flag(SandboxPortFlagName).Value.String()

	port, _ := strconv.Atoi(portStr)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &SandboxHandler{})
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}

	fmt.Printf("Started server at localhost:%d...", port)
}