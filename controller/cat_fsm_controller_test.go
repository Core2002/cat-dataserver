package controller

import (
	"testing"

	"fifu.fun/cat-dataserver/repository"
)

func TestNewCatFSMController(t *testing.T) {
	repo := repository.NewCatFSMRepository()
	ctrl := NewCatFSMController(repo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}
