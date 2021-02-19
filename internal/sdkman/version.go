package sdkman

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/alex-held/devctl/pkg/aarch"
	"github.com/blang/semver"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type VersionService service
type SDKVersion struct{}

var logger = log.New()

func (s *VersionService) All(ctx context.Context, sdk string, arch aarch.Arch) (versions []semver.Version, err error)  {

	req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("candidates/%s/%s/versions/all", sdk, arch.String()), http.NoBody)
	
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

	for _, v := range  strings.Split(body, ",") {
		ver, err := semver.ParseTolerant(v)
		if err != nil {
			return nil, err
		}
		versions = append(versions, ver)
	}
	
	return versions, nil
}


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
