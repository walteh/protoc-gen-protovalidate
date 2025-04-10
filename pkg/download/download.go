package download

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(ctx context.Context, language string, releaseTag string) ([]byte, []byte, error) {

	formatz, err := getJson(ctx, language, releaseTag)
	if err != nil {
		return nil, nil, err
	}

	formatted, err := json.MarshalIndent(formatz, "", "\t")
	if err != nil {
		return nil, nil, err
	}

	body, err := downloadTarGz(ctx, formatz["tarball_url"].(string))
	if err != nil {
		return nil, nil, err
	}

	return formatted, body, nil
}

func getJson(ctx context.Context, language string, releaseTag string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://api.github.com/repos/bufbuild/protovalidate-%s/releases/%s", language, releaseTag)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download protovalidate-%s-latest.json: %s", language, resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var formatz map[string]interface{}
	err = json.Unmarshal(body, &formatz)
	if err != nil {
		return nil, err
	}

	return formatz, nil
}

func downloadTarGz(ctx context.Context, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
