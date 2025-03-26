package acrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestToStringArray tests the ToStringArray method of MiniRandomCodes.
func TestToStringArray(t *testing.T) {
	codes := MiniRandomCodes{"code1", "code2"}
	strArray := codes.ToStringArray()

	if len(strArray) != len(codes) {
		t.Errorf("Expected length %d, got %d", len(codes), len(strArray))
	}

	for i, code := range codes {
		if strArray[i] != code {
			t.Errorf("Expected string %s, got %s", code, strArray[i])
		}
	}
}

// TestMatchAndRemove tests the MatchAndRemove method of MiniRandomCodes.
func TestMatchAndRemove(t *testing.T) {
	codes := &MiniRandomCodes{"code1", "code2", "code3"}
	target := "code2"
	found := codes.MatchAndRemove(target)

	if !found {
		t.Errorf("Expected to find %s", target)
	}

	if len(*codes) != 2 {
		t.Errorf("Expected length 2, got %d", len(*codes))
	}

	for _, code := range *codes {
		if code == target {
			t.Errorf("Did not expect to find %s", target)
		}
	}
}

// TestGenerate tests the Generate method of MiniRandomCodes.
func TestGenerate(t *testing.T) {
	codes := &MiniRandomCodes{}
	err := codes.Generate(5, 10, "-")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(*codes) != 5 {
		t.Errorf("Expected 5 codes, got %d", len(*codes))
	}

	for _, code := range *codes {
		if len(code) != 11 { // 10 characters + 1 divider
			t.Errorf("Expected code length 11, got %d", len(code))
		}
	}
}

// TestGenerateMiniRandomCodes tests the standalone GenerateMiniRandomCodes function.
func TestGenerateMiniRandomCodes(t *testing.T) {
	codes, err := GenerateMiniRandomCodes(5, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(codes) != 5 {
		t.Errorf("Expected 5 codes, got %d", len(codes))
	}

	for _, code := range codes {
		if len(code) != 10 {
			t.Errorf("Expected code length 10, got %d", len(code))
		}
	}
}

// TestGenerateWithInvalidParams tests the error handling of the Generate method.
func TestGenerateWithInvalidParams(t *testing.T) {
	codes := &MiniRandomCodes{}
	err := codes.Generate(-1, 10, "-")
	if err == nil {
		t.Errorf("Expected an error for invalid count, got nil")
	}

	err = codes.Generate(5, -1, "-")
	if err == nil {
		t.Errorf("Expected an error for invalid length, got nil")
	}
}

// TestGenerateMiniRandomCodesInvalidParams checks the error handling of GenerateMiniRandomCodes.
func TestGenerateMiniRandomCodesInvalidParams(t *testing.T) {
	_, err := GenerateMiniRandomCodes(-1, 6)
	if err == nil {
		t.Error("GenerateMiniRandomCodes should return an error for negative count")
	}

	_, err = GenerateMiniRandomCodes(5, -1)
	if err == nil {
		t.Error("GenerateMiniRandomCodes should return an error for negative length")
	}
}

func TestMiniRandomCodes_GenerateOne(t *testing.T) {
	codes := MiniRandomCodes{}
	err := codes.Generate(1, 6, "")
	assert.NoError(t, err)
	assert.Len(t, codes, 1)
	assert.Len(t, codes[0], 6)

	codes = MiniRandomCodes{}
	err = codes.Generate(1, 6, "-")
	assert.NoError(t, err)
	assert.Len(t, codes, 1)
	assert.Len(t, codes[0], 7) // 6 characters + 1 divider
}

func TestMiniRandomCodes_Generate(t *testing.T) {
	codes := MiniRandomCodes{}
	err := codes.Generate(16, 10, "-")
	assert.NoError(t, err)
	assert.Len(t, codes, 16)

	for _, code := range codes {
		assert.Len(t, code, 11) // 10 characters + 1 divider
	}

	set1 := codes[:8]
	assert.Len(t, set1, 8)
	set2 := codes[8:]
	assert.Len(t, set2, 8)

	codes = MiniRandomCodes{}
	err = codes.Generate(16, 10, "")
	assert.NoError(t, err)
	assert.Len(t, codes, 16)

	for _, code := range codes {
		assert.Len(t, code, 10)
	}

	codes = MiniRandomCodes{}
	err = codes.Generate(16, 9, "")
	assert.NoError(t, err)
	assert.Len(t, codes, 16)

	for _, code := range codes {
		assert.Len(t, code, 9)
	}

	codes = MiniRandomCodes{}
	err = codes.Generate(16, 9, "-")
	assert.NoError(t, err)
	assert.Len(t, codes, 16)

	for _, code := range codes {
		assert.Len(t, code, 10) // 9 characters + 1 divider
	}
}
