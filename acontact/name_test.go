package acontact

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName_Validate(t *testing.T) {
	t.Run("Valid person name", func(t *testing.T) {
		n := &Name{First: "John", Last: "Doe"}
		assert.NoError(t, n.Validate())
	})

	t.Run("Valid entity name", func(t *testing.T) {
		n := &Name{Company: "TechCorp"}
		assert.NoError(t, n.Validate())
	})

	t.Run("Invalid name", func(t *testing.T) {
		n := &Name{}
		assert.Error(t, n.Validate())
	})
}

func TestName_IsPerson(t *testing.T) {
	assert.True(t, (&Name{First: "John"}).IsPerson())
	assert.True(t, (&Name{Last: "Doe"}).IsPerson())
	assert.False(t, (&Name{Company: "TechCorp"}).IsPerson())
}

func TestName_IsEntity(t *testing.T) {
	assert.True(t, (&Name{Company: "TechCorp"}).IsEntity())
	assert.False(t, (&Name{First: "John"}).IsEntity())
}

func TestName_IsBothPersonAndEntity(t *testing.T) {
	assert.True(t, (&Name{First: "John", Company: "TechCorp"}).IsBothPersonAndEntity())
	assert.False(t, (&Name{First: "John"}).IsBothPersonAndEntity())
}

func TestName_GetFirstLastName(t *testing.T) {
	assert.Equal(t, "John Doe", (&Name{First: "John", Last: "Doe"}).GetFirstLastName())
	assert.Equal(t, "John", (&Name{First: "John"}).GetFirstLastName())
	assert.Equal(t, "Doe", (&Name{Last: "Doe"}).GetFirstLastName())
	assert.Equal(t, "", (&Name{}).GetFirstLastName())
}

func TestName_GetName(t *testing.T) {
	assert.Equal(t, "John Doe", (&Name{First: "John", Last: "Doe"}).GetName())
	assert.Equal(t, "TechCorp", (&Name{Company: "TechCorp"}).GetName())
	assert.Equal(t, "", (&Name{}).GetName())
}

func TestName_GetNamePlusCompany(t *testing.T) {
	assert.Equal(t, "John Doe (TechCorp)", (&Name{First: "John", Last: "Doe", Company: "TechCorp"}).GetNamePlusCompany())
	assert.Equal(t, "John Doe", (&Name{First: "John", Last: "Doe"}).GetNamePlusCompany())
	assert.Equal(t, "TechCorp", (&Name{Company: "TechCorp"}).GetNamePlusCompany())
}

func TestName_GetMail(t *testing.T) {
	assert.Equal(t, []string{"TechCorp", "John Doe"}, (&Name{First: "John", Last: "Doe", Company: "TechCorp"}).GetMail())
	assert.Equal(t, []string{"John Doe"}, (&Name{First: "John", Last: "Doe"}).GetMail())
	assert.Equal(t, []string{"TechCorp"}, (&Name{Company: "TechCorp"}).GetMail())
}

func TestName_GetMailWithTitleDepartment(t *testing.T) {
	t.Run("Both person and entity with title", func(t *testing.T) {
		name := Name{
			First:   "John",
			Last:    "Doe",
			Title:   "Dr.",
			Company: "TechCorp",
		}
		result := name.GetMailWithTitleDepartment()
		expected := []string{"Dr. John Doe", "TechCorp"}
		assert.Equal(t, expected, result)
	})

	t.Run("Person with title and department", func(t *testing.T) {
		name := Name{
			First:      "Jane",
			Last:       "Smith",
			Title:      "Ms.",
			Department: "Engineering",
		}
		result := name.GetMailWithTitleDepartment()
		expected := []string{"Ms. Jane Smith", "Engineering"}
		assert.Equal(t, expected, result)
	})

	t.Run("Entity only", func(t *testing.T) {
		name := Name{
			Company: "TechCorp",
		}
		result := name.GetMailWithTitleDepartment()
		expected := []string{"TechCorp"}
		assert.Equal(t, expected, result)
	})

	t.Run("Person without title or department", func(t *testing.T) {
		name := Name{
			First: "John",
			Last:  "Doe",
		}
		result := name.GetMailWithTitleDepartment()
		expected := []string{"John Doe"}
		assert.Equal(t, expected, result)
	})
}

func TestName_CanSign(t *testing.T) {
	t.Run("Valid person with company", func(t *testing.T) {
		name := Name{
			First:   "John",
			Last:    "Doe",
			Title:   "Dr.",
			Company: "TechCorp",
		}
		err := name.CanSign(true)
		assert.NoError(t, err)
	})

	t.Run("Valid person without company", func(t *testing.T) {
		name := Name{
			First: "Jane",
			Last:  "Smith",
			Title: "Ms.",
		}
		err := name.CanSign(false)
		assert.NoError(t, err)
	})

	t.Run("Missing title", func(t *testing.T) {
		name := Name{
			First:   "Jane",
			Last:    "Smith",
			Company: "TechCorp",
		}
		err := name.CanSign(true)
		assert.Error(t, err)
		assert.Equal(t, "title is empty", err.Error())
	})

	t.Run("Missing name", func(t *testing.T) {
		name := Name{
			Title:   "Dr.",
			Company: "TechCorp",
		}
		err := name.CanSign(true)
		assert.Error(t, err)
		assert.Equal(t, "name is empty", err.Error())
	})

	t.Run("Missing company", func(t *testing.T) {
		name := Name{
			First: "John",
			Last:  "Doe",
			Title: "Mr.",
		}
		err := name.CanSign(true)
		assert.Error(t, err)
		assert.Equal(t, "company is empty", err.Error())
	})
}

func TestName_GetSignMap(t *testing.T) {
	t.Run("Valid person with company", func(t *testing.T) {
		name := Name{
			First:   "John",
			Last:    "Doe",
			Title:   "Dr.",
			Company: "TechCorp",
		}
		signMap, err := name.GetSignMap(true)
		assert.NoError(t, err)
		expected := map[string]string{
			"title":   "Dr.",
			"name":    "John Doe",
			"company": "TechCorp",
		}
		assert.Equal(t, expected, signMap)
	})

	t.Run("Valid person without company", func(t *testing.T) {
		name := Name{
			First: "Jane",
			Last:  "Smith",
			Title: "Ms.",
		}
		signMap, err := name.GetSignMap(false)
		assert.NoError(t, err)
		expected := map[string]string{
			"title": "Ms.",
			"name":  "Jane Smith",
		}
		assert.Equal(t, expected, signMap)
	})

	t.Run("Invalid person (missing title)", func(t *testing.T) {
		name := Name{
			First: "John",
			Last:  "Doe",
		}
		signMap, err := name.GetSignMap(false)
		assert.Error(t, err)
		assert.Nil(t, signMap)
		assert.Equal(t, "title is empty", err.Error())
	})

	t.Run("Invalid person (missing name)", func(t *testing.T) {
		name := Name{
			Title: "Dr.",
		}
		signMap, err := name.GetSignMap(false)
		assert.Error(t, err)
		assert.Nil(t, signMap)
		assert.Equal(t, "name is empty", err.Error())
	})

	t.Run("Invalid person (missing company)", func(t *testing.T) {
		name := Name{
			First: "Jane",
			Last:  "Smith",
			Title: "Ms.",
		}
		signMap, err := name.GetSignMap(true)
		assert.Error(t, err)
		assert.Nil(t, signMap)
		assert.Equal(t, "company is empty", err.Error())
	})
}
