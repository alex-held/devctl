package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type VersionService service
type SDKVersion struct{}

var logger = log.New()

func (s *VersionService) Default(ctx context.Context, sdk string) (version *semver.Version, err error) {
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
		return nil, errors.Wrap(err, "read from response body failed")
	}
	defer resp.Body.Close()
	body := buf.String()
	sdkSemVer, err := semver.ParseTolerant(body)
	if err != nil {
		err = errors.Wrap(err, "parsing semver failed")
		logger.
			WithError(err).
			WithField("body", body).
			WithField("url", req.URL.String()).
			WithField("sdk", sdk).
			Errorln("response is no valid semver")
		return nil, err
	}

	return &sdkSemVer, err
}
