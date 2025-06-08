package customers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Customer Registration", func() {
	It("registers a new customer successfully", func() {
		url := "http://localhost:80/v1.0/customers/register"
		payload := map[string]any{
			"email":    "e2e_test_user_" + time.Now().Format("150405") + "@example.com",
			"password": "strongpassword123",
			"name":     "E2E Test User",
		}
		body, err := json.Marshal(payload)
		Expect(err).NotTo(HaveOccurred())

		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		Expect(err).NotTo(HaveOccurred())
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		var result struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Email).To(Equal(payload["email"]))
		Expect(result.Name).To(Equal(payload["name"]))
		Expect(result.ID).To(MatchRegexp(`^[a-fA-F0-9]{24}$`))
		Expect(result.CreatedAt).NotTo(BeEmpty())
	})
})
