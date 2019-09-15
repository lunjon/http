package rest

import (
	"fmt"
	"net/http"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/spf13/cobra"
)

type SandboxHandler struct{}

func (h *SandboxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`{"url": "%s", "method": "%s"}`, r.URL, r.Method)
	_, _ = w.Write([]byte(body))
}

func StartSandbox(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt(constants.SandboxPortFlagName)

	fmt.Printf("Starting server at localhost:%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &SandboxHandler{})
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}
}
