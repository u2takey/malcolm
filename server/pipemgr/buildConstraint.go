package pipemgr

import (
	"github.com/Sirupsen/logrus"
	"github.com/arlert/malcolm/model"
	// "github.com/arlert/malcolm/utils"
)

func matchConstraint(build *model.Build) bool {
	if build.CurrentStep == 0 {
		return true
	}
	if build.CurrentStep >= len(build.Steps) {
		return false
	}
	curstep := build.Steps[build.CurrentStep]
	if curstep.Config.Constraint == nil {
		return true
	}
	constraint := curstep.Config.Constraint
	laststep := build.Steps[build.CurrentStep-1]
	logrus.Debugln("matchConstraint", constraint.MatchState, laststep.Status.StateDetail)
	if (constraint.MatchState == model.MatchStateAlways) ||
		(laststep.Status.StateDetail == model.StateCompleteDetailSuccess &&
			constraint.MatchState == model.MatchStateSuccess) ||
		(laststep.Status.StateDetail == model.StateCompleteDetailFailed &&
			constraint.MatchState == model.MatchStateFail) {
		return true
	}
	return false
}
