package user

// Profile represents a user's profile information
type Profile struct {
	name    string
	picture string
}

// NewProfile creates a new Profile
func NewProfile(name, picture string) Profile {
	return Profile{
		name:    name,
		picture: picture,
	}
}

// Name returns the user's display name
func (p Profile) Name() string {
	return p.name
}

// Picture returns the URL of the user's profile picture
func (p Profile) Picture() string {
	return p.picture
}

// WithName creates a new Profile with updated name
func (p Profile) WithName(name string) Profile {
	return Profile{
		name:    name,
		picture: p.picture,
	}
}

// WithPicture creates a new Profile with updated picture
func (p Profile) WithPicture(picture string) Profile {
	return Profile{
		name:    p.name,
		picture: picture,
	}
}

// IsEmpty returns true if the profile has no name and no picture
func (p Profile) IsEmpty() bool {
	return p.name == "" && p.picture == ""
}
