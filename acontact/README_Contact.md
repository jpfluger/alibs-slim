# Contact

The `acontact` package provides a comprehensive framework for managing contacts in SaaS applications or other systems requiring sophisticated contact management capabilities. At its core, it defines flexible and extensible data structures for storing, managing, and validating contact information for individuals, businesses, or hybrid entities.

---

## Core Structures

### 1. **ContactCore**
The `ContactCore` struct represents the foundational structure for a contact, including:

- **Name**: Contains personal or organizational identifying fields like `First`, `Last`, `Company`, `Title`, and `Department`.
- **Communication Channels**: Collections of email addresses (`Emails`), phone numbers (`Phones`), mailing addresses (`Mails`), and URLs (`Urls`).
- **Metadata**:
    - Dates for significant events (`Dates`).
    - Images (`Images`).
    - Tags for categorization (`Tags`).
- **Details**: A generic `json.RawMessage` field for custom-structured data.

### 2. **Contact**
A `Contact` extends `ContactCore` by adding a globally unique identifier (`CID`), making it suitable for use in systems that require contact deduplication or tracking.

### 3. **ContactUID**
A `ContactUID` further extends `Contact` with a `UID` to link the contact to a specific user, along with metadata for:
- Default status (`IsDefault`).
- Entity types (`EntityTypes`).
- Person types (`PersonTypes`).

---

## Use Cases

The `acontact` package is versatile and can be implemented in various SaaS applications or systems. Below are **10 use cases**:

### 1. **CRM (Customer Relationship Management)**
- Manage customer contact information, including names, emails, phones, and company details.
- Categorize contacts using `Tags` for segmentation and targeted marketing campaigns.

### 2. **Document Signing Platforms**
- Ensure a person or entity has the necessary information (`Title`, `Name`, `Company`) to sign legal documents.
- Use `CanSign` validation for automated workflows in digital signature applications.

### 3. **Helpdesk or Support Systems**
- Store contact details for customers and their associated companies to streamline ticket handling.
- Link contacts to specific users via `ContactUID` for personalized support.

### 4. **Event Management**
- Manage registrants' details, including their `Title` and `Department` for personalized communication.
- Use `Dates` to track important milestones like registration deadlines or event attendance.

### 5. **E-Commerce Platforms**
- Store billing and shipping addresses (`Mails`) and phone numbers for customer orders.
- Use `Tags` to categorize customers based on purchase history or loyalty programs.

### 6. **Human Resource Management (HRM)**
- Store employee contact information, including `Title`, `Department`, and `Company` affiliation.
- Use `Dates` to track employment anniversaries or other HR-relevant dates.

### 7. **B2B Relationship Management**
- Manage business-to-business relationships with detailed `EntityTypes` and `PersonTypes`.
- Track multiple points of contact within a single organization.

### 8. **Email Marketing Campaigns**
- Use `Emails` and `Tags` to segment and target specific customer groups for email campaigns.
- Maintain clean, validated email lists using the `Validate` functions.

### 9. **Appointment Scheduling**
- Manage customer contact information, including phone numbers and emails, for reminders and notifications.
- Track relevant dates (`Dates`) for scheduled appointments.

### 10. **Real Estate Applications**
- Store client and company details for buyers, sellers, and agents.
- Use `Tags` to categorize contacts as "Prospect," "Lead," or "Closed Deal."

---

## Key Features

1. **Flexible Design**:
    - Accommodates both individuals and organizations with fields like `Title`, `Department`, and `Company`.

2. **Extensible Metadata**:
    - Use `Details` for custom data specific to your application.

3. **Robust Validation**:
    - Built-in validation functions like `CanSign` ensure data integrity for critical operations.

4. **Comprehensive Support for Communication Channels**:
    - Includes `Emails`, `Phones`, `Mails`, and `Urls`.

5. **Hierarchical Relationships**:
    - `ContactUID` links contacts to specific users, supporting multi-tenant systems.

6. **Categorization and Tagging**:
    - Use `Tags` to organize and group contacts efficiently.

---

## Getting Started

### Installation

```bash
go get github.com/yourorganization/acontact
```

### Example Usage

#### Create a Contact
```go
contact := ContactCore{
	Name: Name{
		First:      "John",
		Last:       "Doe",
		Company:    "TechCorp",
		Title:      "Dr.",
		Department: "Engineering",
	},
	Emails: Emails{{Type: "work", Address: "john.doe@techcorp.com"}},
	Phones: Phones{{Type: "mobile", Number: "+1234567890"}},
	Mails:  Mails{{Type: "home", Address: "123 Main St, Anytown, USA"}},
	Tags:   atags.TagArrStrings{"VIP", "NewsletterSubscriber"},
}
```

#### Validate for Signing
```go
err := contact.Name.CanSign(true)
if err != nil {
	fmt.Printf("Contact cannot sign: %v\n", err)
} else {
	fmt.Println("Contact is qualified to sign.")
}
```

#### Get Mail Information
```go
mail := contact.Name.GetMailWithTitleDepartment()
fmt.Println("Mail Information:", mail)
```

#### Save and Load Custom Details
```go
type CustomDetails struct {
	Role        string `json:"role"`
	Preferences string `json:"preferences"`
}

details := CustomDetails{Role: "Manager", Preferences: "EmailOnly"}
contact.SaveDetails(details)

var loadedDetails CustomDetails
contact.LoadDetails(&loadedDetails)
fmt.Printf("Loaded Details: %+v\n", loadedDetails)
```
