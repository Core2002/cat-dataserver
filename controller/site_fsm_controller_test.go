package controller

import (
	"testing"

	"fifu.fun/cat-dataserver/repository"
)

func TestNewSiteFSMController(t *testing.T) {
	repo := repository.NewSiteFSMRepository()
	siteRepo := repository.NewSiteRepository()
	ctrl := NewSiteFSMController(repo, siteRepo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
	if ctrl.siteRepo != siteRepo {
		t.Error("Controller siteRepo does not match input siteRepo")
	}
}
