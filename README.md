# LingerAPI
LingerAPI is a tool designed to retrieve information from TikTok profiles using the /api/tiktok?username= endpoint. It returns results in JSON format, including potential bio links and social media profiles found within a user's biography.

# Usage
To use LingerAPI, send a GET request to /api/tiktok?username=, where username is the TikTok username for which you want to obtain information.

Example request:

```
GET /api/tiktok?username=johndoe
```

Response:
```
{
    "username": "username",
    "followers": 831800,
    "bio": [
        "https://linktr.ee/username"
    ],
    "soclinks": [
        "https://www.facebook.com/username,
        "https://www.instagram.com/username,
        "https://instagram.com/username,
        "https://youtube.com/user/username,
        "https://instagram.com/username,
    ]
}
```

# Filter Configuration
You can configure filters to process links found within user biographies. Filters are stored in text files in the filter directory. You can edit and customize these files to define rules for analyzing bio links.
