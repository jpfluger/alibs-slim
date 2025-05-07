package azb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViewPort_BasicDimensions(t *testing.T) {
	vp := ViewPort{
		Window:  ViewPortSize{W: 1280, H: 720},
		Content: ViewPortSize{W: 1024, H: 600},
	}

	assert.Equal(t, 1280, vp.Width())
	assert.Equal(t, 720, vp.Height())
	assert.Equal(t, 1024, vp.ContentWidth())
	assert.Equal(t, 600, vp.ContentHeight())
}

func TestViewPort_DeviceType(t *testing.T) {
	tests := []struct {
		width     int
		expected  string
		isPhone   bool
		isTablet  bool
		isDesktop bool
	}{
		{0, "unknown", false, false, false},
		{320, "phone", true, false, false},
		{767, "phone", true, false, false},
		{768, "tablet", false, true, false},
		{1023, "tablet", false, true, false},
		{1024, "desktop", false, false, true},
		{1920, "desktop", false, false, true},
	}

	for _, tt := range tests {
		vp := ViewPort{
			Window: ViewPortSize{W: tt.width, H: 800},
		}
		assert.Equal(t, tt.expected, vp.DeviceType(), "DeviceType() failed for width=%d", tt.width)
		assert.Equal(t, tt.isPhone, vp.IsPhone(), "IsPhone() failed for width=%d", tt.width)
		assert.Equal(t, tt.isTablet, vp.IsTablet(), "IsTablet() failed for width=%d", tt.width)
		assert.Equal(t, tt.isDesktop, vp.IsDesktop(), "IsDesktop() failed for width=%d", tt.width)
	}
}

func TestViewPort_BootstrapType(t *testing.T) {
	tests := []struct {
		width    int
		expected string
	}{
		{0, "xs"},
		{575, "xs"},
		{576, "sm"},
		{767, "sm"},
		{768, "md"},
		{991, "md"},
		{992, "lg"},
		{1199, "lg"},
		{1200, "xl"},
		{1399, "xl"},
		{1400, "xxl"},
		{2000, "xxl"},
	}

	for _, tt := range tests {
		vp := ViewPort{Window: ViewPortSize{W: tt.width}}
		assert.Equal(t, tt.expected, vp.BootstrapType(), "BootstrapType failed for width %d", tt.width)
	}
}

func TestViewPort_BootstrapComparisons(t *testing.T) {
	vp := ViewPort{Window: ViewPortSize{W: 800}}

	assert.True(t, vp.IsLtBootstrap("lg"))
	assert.False(t, vp.IsLteBootstrap("md"))
	assert.False(t, vp.IsLtBootstrap("md"))
	assert.True(t, vp.IsGteBootstrap("sm"))
	assert.False(t, vp.IsGtBootstrap("xl"))
	assert.False(t, vp.IsGtBootstrap("unknown")) // graceful fallback
}

func TestViewPort_BootstrapNamedChecks(t *testing.T) {
	tests := []struct {
		width             int
		isXs, isSm, isMd  bool
		isLg, isXl, isXxl bool
	}{
		{500, true, false, false, false, false, false},
		{600, false, true, false, false, false, false},
		{800, false, false, true, false, false, false},
		{1000, false, false, false, true, false, false},
		{1250, false, false, false, false, true, false},
		{1450, false, false, false, false, false, true},
	}

	for _, tt := range tests {
		vp := ViewPort{Window: ViewPortSize{W: tt.width}}

		assert.Equal(t, tt.isXs, vp.IsXs(), "IsXs failed for width %d", tt.width)
		assert.Equal(t, tt.isSm, vp.IsSm(), "IsSm failed for width %d", tt.width)
		assert.Equal(t, tt.isMd, vp.IsMd(), "IsMd failed for width %d", tt.width)
		assert.Equal(t, tt.isLg, vp.IsLg(), "IsLg failed for width %d", tt.width)
		assert.Equal(t, tt.isXl, vp.IsXl(), "IsXl failed for width %d", tt.width)
		assert.Equal(t, tt.isXxl, vp.IsXxl(), "IsXxl failed for width %d", tt.width)
	}
}
