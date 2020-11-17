package agent

import (
	"testing"
)

func TestUnwaste(t *testing.T) {
	che := Chemical{1000, 1000, 1000, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	che.Unwaste()
	if che.WASTE != 500 {
		t.Error("Expected 500, got ", che.WASTE)
	}
}

func TestAddGlucose(t *testing.T) {
	che := Chemical{1024, 1000, 1000, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	che.AddGluckose()
	if che.GLUC != 1090 {
		t.Error("Expected 1090, got ", che.GLUC)
	}
}
