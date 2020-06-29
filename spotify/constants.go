package spotify

const (
	authURL string = "https://accounts.spotify.com"
	baseURL string = "https://api.spotify.com"

	// ImageUpload 	             Write access to user-provided images.
	ImageUpload string = "ugc-image-upload"

	// PlaylistReadCollaborative Include collaborative playlists when requesting a user's playlists.
	PlaylistReadCollaborative string = "playlist-read-collaborative"
	// PlaylistModifyPrivate 	 Write access to a user's private playlists.
	PlaylistModifyPrivate string = "playlist-modify-private"
	// PlaylistModifyPublic 	 Write access to a user's public playlists.
	PlaylistModifyPublic string = "playlist-modify-public"
	// PlaylistReadPrivate 		 Read access to user's private playlists.
	PlaylistReadPrivate string = "playlist-read-private"

	// UserLibraryModify         Write/delete access to a user's "Your Music" library.
	UserLibraryModify string = "user-library-modify"
	// UserLibraryRead           Read access to a user's "Your Music" library.
	UserLibraryRead string = "user-library-read"

	// UserReadPrivate           Read access to user’s subscription details (type of user account).
	UserReadPrivate string = "user-read-private"
	// UserReadEmail             Read access to user’s email address.
	UserReadEmail string = "user-read-email"

	// UserFollowModify          Write/delete access to the list of artists and other users that the user follows.
	UserFollowModify string = "user-follow-modify"
	// UserFollowRead            Read access to the list of artists and other users that the user follows.
	UserFollowRead string = "user-follow-read"
)
