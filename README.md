# LingerAPI
🔗 LingerAPI is a tool designed to retrieve information from TikTok profiles using the /api/tiktok/ endpoint and Linktree profiles using /api/bio/. It returns results in JSON format, including potential bio links and social media profiles found within a user's biography.

![image](https://github.com/trashplusplus/LingerAPI/assets/19663951/2ba39fdc-f0ff-457f-a514-2bdc60a18415)


# Usage
🔗 To use LingerAPI, send a GET request to /api/tiktok?username=, where username is the TikTok username for which you want to obtain information.

# Install

Clone the LingerAPI repository from GitHub:
```
git clone https://github.com/trashplusplus/LingerAPI.git
```
Change your working directory to the LingerAPI project folder:
```
cd LingerAPI
```
Initialize the Go Module:
```
git mod init LingerAPI
```
Download and Install Dependencies:
```
git mod tidy
```
To run the LingerAPI:
```
go run .
```

# Example

🔗 Possible endpoints
```
GET /api/tiktok/
GET /api/bio/
```


➡️ Request /api/tiktok/:

```
GET /api/tiktok?username=username
```

📃 Response:
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
➡️ Request /api/bio/:
```
GET /api/bio?username=https://linktr.ee/deftones
```
📃 Response:
```
{
    "bio": [
        "https://linktr.ee/deftones"
    ],
    "soclinks": [
        "https://youtube.com/playlist?list=PLNRsYvRgbfmpyK7YFeyn5_OPekQUiSLVa",
        "https://www.instagram.com/ar/701930743692711/",
        "https://www.facebook.com/fbcameraeffects/tryit/372370910427346/",
        "https://youtu.be/KUDbj0oeAj0",
        "https://youtube.com/c/deftones",
        "https://www.facebook.com/deftones/",
        "https://instagram.com/deftones",
        "https://www.youtube.com/c/deftones",
        "https://www.twitch.tv/deftonesofficial"
    ]
}
```

# Filter Configuration
You can configure filters to process links found within user biographies. Filters are stored in text files in the filter directory. You can edit and customize these files to define rules for analyzing bio links.
