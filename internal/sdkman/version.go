package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/alex-held/devctl/pkg/aarch"
)

type VersionService service

var logger = log.New()

func (s *VersionService) All(ctx context.Context, sdk string, arch aarch.Arch) (versions []string, err error) {
	req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("candidates/%s/%s/versions/all", sdk, arch), http.NoBody)

	if err != nil {
		logger.
			WithError(err).
			WithField("sdk", sdk).
			Errorln("could not build http.Request for VersionService.Default")
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		logger.
			WithError(err).
			WithField("sdk", sdk).
			WithField("url", req.URL.String()).
			Errorln("http.Request failed for VersionService.All")
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, errors.Wrap(err, "read from response body failed")
	}
	defer resp.Body.Close()
	body := buf.String()
	versions = strings.Split(body, ",")
	return versions, nil
}

func (s *VersionService) Default(ctx context.Context, sdk string) (version string, err error) {
	req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("candidates/default/%s", sdk), http.NoBody)

	if err != nil {
		logger.
			WithError(err).
			WithField("sdk", sdk).
			Errorln("could not build http.Request for VersionService.Default")
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		logger.
			WithError(err).
			WithField("sdk", sdk).
			WithField("url", req.URL.String()).
			Errorln("http.Request failed for VersionService.Default")
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return "", errors.Wrap(err, "read from response body failed")
	}
	defer resp.Body.Close()
	body := buf.String()
	return body, nil
}
