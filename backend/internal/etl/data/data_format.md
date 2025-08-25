# Data Format Documentation

This document captures the actual response structures and data formats from our ETL extractors based on testing results.

## 📱 Instagram Extractor

### API Configuration
- **API Name**: Instagram Premium API 2023
- **Host**: `instagram-premium-api-2023.p.rapidapi.com`
- **Endpoint**: `/v1/hashtag/medias/top/recent/chunk`
- **Status**: ✅ **FULLY FUNCTIONAL**

### Response Structure
The Instagram API returns an **array** with 2 elements:
1. **First element**: Array of Instagram posts
2. **Second element**: Cursor/pagination token (string)

### Sample Response Format
```json
[
  [
    {
      "accessibility_caption": null,
      "caption_text": "🧑‍⚖️ Un juez federal de EU autorizó la detención...",
      "code": "DNr3fXv3g30",
      "comment_count": 23,
      "comments_disabled": false,
      "has_liked": false,
      "id": "3705299166803398132_1448725107",
      "image_versions": [
        {
          "height": 1350,
          "url": "https://scontent-lax3-1.cdninstagram.com/v/t51.2885-15/...",
          "width": 1080
        }
      ],
      "is_paid_partnership": false,
      "like_and_view_counts_disabled": false,
      "like_count": 555,
      "location": null,
      "media_type": 1,
      "pk": "3705299166803398132",
      "play_count": 0,
      "product_type": "feed",
      "resources": [],
      "sponsor_tags": [],
      "taken_at": "2025-08-23T05:15:06Z",
      "taken_at_ts": 1755926106,
      "thumbnail_url": "https://scontent-lax3-1.cdninstagram.com/v/t51.2885-15/...",
      "title": "",
      "user": {
        "full_name": "Animal Político",
        "id": "1448725107",
        "is_private": false,
        "is_verified": true,
        "pk": "1448725107",
        "profile_pic_url": "https://scontent-lax3-1.cdninstagram.com/v/t51.2885-19/...",
        "username": "pajaropolitico"
      },
      "usertags": [],
      "video_dash_manifest": "",
      "video_duration": 0,
      "video_url": null,
      "video_versions": [],
      "view_count": 0
    }
  ],
  "WyJmZDEyNzlmMTg5MWE0YWFjYTE4YjZkYmU0ZTZmYWQ0NyIsW10sMV0="
]
```

### Key Data Fields
- **`code`**: Unique post identifier (e.g., "DNr3fXv3g30")
- **`caption_text`**: Post content/caption
- **`like_count`**: Number of likes
- **`comment_count`**: Number of comments
- **`user.username`**: Instagram username
- **`user.full_name`**: Display name
- **`taken_at`**: Post timestamp
- **`media_type`**: Type of media (1=image, 2=video, 8=carousel)
- **`image_versions`**: Array of image URLs in different sizes
- **`pk`**: Internal post ID

### Go Struct Mapping
```go
type InstagramResponse struct {
    Status  string      `json:"status"`
    Data    interface{} `json:"data,omitempty"` // For backward compatibility
    Error   string      `json:"error,omitempty"`
    Hashtag string      `json:"hashtag,omitempty"`
    MaxID   string      `json:"max_id,omitempty"`
    MediaID string      `json:"media_id,omitempty"`
    Amount  int         `json:"amount,omitempty"`
    
    // Direct API response fields for array structure
    Posts  []interface{} `json:"posts,omitempty"`  // Posts data from first array element
    Cursor string        `json:"cursor,omitempty"` // Cursor token from second array element
}
```

---

## 📺 YouTube Extractor

### API Configuration
- **API Name**: YouTube Data API v3 (via RapidAPI)
- **Host**: `youtube138.p.rapidapi.com`
- **Status**: ✅ **FULLY FUNCTIONAL**

### Search Videos Response Structure
**Endpoint**: `/search/?q={query}&hl={language}&gl={geo}`

```json
{
  "status": "success",
  "contents": [
    {
      "type": "video",
      "video": {
        "author": {
          "avatar": [
            {
              "height": 68,
              "url": "https://yt3.ggpht.com/8Lwf4LCR2VmxD2JKiozRu7Lo2jGdnhRs42NawHmMN_xJ8TdW-30e3J9DhumEksivp1Esog4A=s88-c-k-c0x00ffffff-no-rj",
              "width": 68
            }
          ],
          "badges": [
            {
              "text": "Official Artist Channel",
              "type": "OFFICIAL_ARTIST_CHANNEL"
            }
          ],
          "canonicalBaseUrl": null,
          "channelId": "UCxoq-PAQeAdk_zyg8YS0JqA",
          "title": "Luis Fonsi",
          "badges": ["CC"],
          "descriptionSnippet": "#LuisFonsi #Despacito #Imposible #Calypso #EchamelaCulpa #NadaEsImposible #NothingisImpossible #LF Music video by Luis ...",
          "isLiveNow": false,
          "lengthSeconds": 282,
          "movingThumbnails": [
            {
              "height": 180,
              "url": "https://i.ytimg.com/an_webp/kJQP7kiw5Fk/mqdefault_6s.webp?du=3000&sqp=CPPnpZUG&rs=AOn4CLCdydAShRnBZW65sOXSvW8fJPi2MA",
              "width": 320
            }
          ],
          "publishedTimeText": "5 years ago",
          "stats": {
            "views": 7870471715
          },
          "thumbnails": [
            {
              "height": 202,
              "url": "https://i.ytimg.com/vi/kJQP7kiw5Fk/hq720.jpg?sqp=-oaymwEcCOgCEMoBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&rs=AOn4CLBCg9eudi8DoM9M-qjPgJBGGkuIgA",
              "width": 360
            }
          ],
          "title": "Luis Fonsi - Despacito ft. Daddy Yankee",
          "videoId": "kJQP7kiw5Fk"
        }
      }
    }
  ],
  "cursorNext": "eyJ2ZXJ0aWNhbCI6ZmFsc2UsImN1cnNvck5leHQiOiJ4b0FJQUFBQUFBQUFCRlRtbG5lMlZ1YVhOemRXVnpkQ0J6YVdkcGJpQnBkQ0JwYmlCbGNuSnZjbWxrWlNCaGJtUWdkMmxrWlNCamJHRnpjM2R2Y21RPSJ9",
  "estimatedResults": 1000000,
  "filterGroups": [...],
  "refinements": [...],
  "tag": "covid 19",
  "geo": "US"
}
```

### Video Comments Response Structure
**Endpoint**: `/video/comments/?id={videoId}`

```json
{
  "status": "success",
  "comments": [
    {
      "author": {
        "avatar": [
          {
            "height": 48,
            "url": "https://yt3.ggpht.com/p4FCAFrFhKYog3VTmKU-YHXF0YuBhV74l4dV6o1OHk2lzIHdOHBFH6rLuMcKkgqT5hgUv6gg1A=s48-c-k-c0x00ffffff-no-rj",
            "width": 48
          }
        ],
        "badges": [],
        "channelId": "UC9Hq2JZT1aKlG8nJqd7CB-A",
        "isChannelOwner": false,
        "title": "Armaan singh"
      },
      "commentId": "UgwW1S95-KHLW2bUI-14AaABAg",
      "content": "2017: People came to listen song. 2021: People come to check views.",
      "creatorHeart": false,
      "cursorReplies": "Eg0SC2tKUVA3a2l3NUZrGAYyfhpLEhpVZ3dXMVM5NS1LSExXMmJVSS0xNEFhQUJBZyICCAAqGFVDTHA4UkJoUUh1OXdTc3E2MmpfTWQ2QTILa0pRUDdraXc1RmtAAUgKQi9jb21tZW50LXJlcGxpZXMtaXRlbS1VZ3dXMVM5NS1LSExXMmJVSS0xNEFhQUJBZw%3D%3D",
      "pinned": {
        "status": false,
        "text": null
      },
      "publishedTimeText": "1 year ago",
      "stats": {
        "replies": 496,
        "votes": 414947
      }
    }
  ],
  "totalCommentsCount": 1234567,
  "filters": [...],
  "video_id": "kJQP7kiw5Fk"
}
```

### Key Data Fields

#### Search Results
- **`contents`**: Array of video results
- **`estimatedResults`**: Total number of search results
- **`cursorNext`**: Pagination token for next page
- **`video.videoId`**: Unique video identifier
- **`video.title`**: Video title
- **`video.author.title`**: Channel name
- **`video.stats.views`**: View count
- **`video.publishedTimeText`**: Upload time

#### Comments
- **`comments`**: Array of comment objects
- **`totalCommentsCount`**: Total number of comments
- **`commentId`**: Unique comment identifier
- **`content`**: Comment text
- **`author.title`**: Commenter name
- **`stats.votes`**: Like/dislike count
- **`publishedTimeText`**: Comment timestamp

### Go Struct Mapping
```go
type YouTubeResponse struct {
    Status  string      `json:"status"`
    Data    interface{} `json:"data,omitempty"` // Kept for compatibility if other APIs use it
    Error   string      `json:"error,omitempty"`
    Tag     string      `json:"tag,omitempty"` // Used for query now
    Geo     string      `json:"geo,omitempty"`
    VideoID string      `json:"video_id,omitempty"`

    // Direct API response fields for search
    Contents         []interface{} `json:"contents,omitempty"`
    CursorNext      string        `json:"cursorNext,omitempty"`
    EstimatedResults int64         `json:"estimatedResults,omitempty"`
    FilterGroups    interface{}   `json:"filterGroups,omitempty"`
    Refinements     interface{}   `json:"refinements,omitempty"`

    // Comments API response fields
    Comments            []interface{} `json:"comments,omitempty"`
    TotalCommentsCount  int64         `json:"totalCommentsCount,omitempty"`
    Filters             interface{}   `json:"filters,omitempty"`
}
```

---

## 📰 Indonesian News Extractor

### API Configuration
- **API Name**: Indonesia News API (via RapidAPI)
- **Host**: `indonesia-news.p.rapidapi.com`
- **Status**: ✅ **FULLY FUNCTIONAL** (3/3 sources working)
- **Note**: TEMPO removed due to server-side issues (502 Bad Gateway)

### Response Structure by Source

#### CNN Indonesia
**Endpoint**: `/search/cnn?query={query}&page={page}&limit={limit}`

```json
{
  "items": [
    {
      "date": {
        "created": "Senin, 25/08/2025 16:19:39",
        "created_timestamp": "1756113579",
        "publish": "Senin, 25/08/2025 16:43:19",
        "publish_timestamp": "1756114999"
      },
      "editor": "Indra Hendriana",
      "idberita": 1266202,
      "idkanal": 625,
      "images": {
        "caption": "",
        "cover": "https://akcdn.detik.net.id/visual/2025/08/25/bri-1756112640600_169.jpeg?w=480",
        "story": "https://akcdn.detik.net.id/visual/2025/08/25/bri-1756112640600_169.jpeg?w=480"
      },
      "namakanal": "Corporate Action",
      "namaparent": "Ekonomi",
      "namasubkanal": "Berita Corporate Action",
      "parentkanal": 624,
      "penulis": "",
      "reporter": "",
      "streditorialreview": "",
      "subtitle": "",
      "summary": "Haluan Bali, brand fashion lokal asal Bali, sukses menembus pasar internasional dengan produk berkelanjutan yang menggabungkan seni tradisional dan modern.",
      "tag": "bri|bbri|pemberdayaan umkm|pengusaha muda brilian|haluan bali|pengusaha muda|keberlanjutan",
      "title": "Pengusaha Muda Bali Pasarkan Fashion Sentuhan Digital di Pameran BRI",
      "type_content": 1,
      "typechannel": 13,
      "url": "https://www.cnnindonesia.com/ekonomi/20250825161939-625-1266202/pengusaha-muda-bali-pasarkan-fashion-sentuhan-digital-di-pameran-bri"
    }
  ],
  "metadata": {
    "limit": 3,
    "page": 1,
    "total": 10000,
    "total_page": 3334
  }
}
```

#### DETIK
**Endpoint**: `/search/detik?keyword={keyword}&limit={limit}&page={page}`

```json
{
  "item": [
    {
      "city": "Solo",
      "date": {
        "created": "Senin, 25 Agustus 2025 17:30 WIB",
        "created_timestamp": "1756117802",
        "publish": "Senin, 25 Agustus 2025 18:01 WIB",
        "publish_timestamp": "1756119708"
      },
      "editor": "ams",
      "id_content_type": "1.0.0.0",
      "idkanal": "1891",
      "idnews": "8078683",
      "images": {
        "caption": "",
        "cover": "https://akcdn.detik.net.id/community/media/visual/2025/08/25/hotel-agas-solo-1756118454881_169.jpeg?w=240&q=90"
      },
      "kanal_parent_name": "detikJateng",
      "link": "https://www.detik.com/jateng/bisnis/d-8078683/berdiri-sejak-1996-hotel-agas-solo-kini-dijual-rp-120-miliar",
      "locname": "Solo, Kota Surakarta, Jawa Tengah, Indonesia",
      "penulis": "afn",
      "reporter": "",
      "subtitle": "",
      "summary": "Hotel legendaris yang berada di tengah Kota Solo, Agas dijual seharga Rp 120 miliar. Hotel yang berada di bawah flyover Manahan itu kukut usai 29 tahun berdiri.",
      "title": "Berdiri sejak 1996. Hotel Agas Solo Kini Dijual Rp 120 Miliar"
    }
  ],
  "metadata": {
    "limit": 3,
    "page": 1,
    "total": 10000,
    "total_page": 3334
  }
}
```

#### KOMPAS
**Endpoint**: `/search/kompas?command={command}&page={page}&limit={limit}`

```json
{
  "xml": {
    "pencarian": {
      "item": [
        {
          "channel": "TREN",
          "description": "Varian baru Covid-19 bernama XFG atau \"Stratus\" kini menjadi perhatian Organisasi Kesehatan Dunia (WHO). Apa gejala khas varian XFG?",
          "guid": ".xml.2025.07.27.220220065",
          "kanal": "Tren",
          "link": "http://www.kompas.com/tren/read/2025/07/27/220220065/gejala-khas-varian-baru-covid-19-xfg-atau-stratus",
          "photo": "https://asset.kompas.com/crops/NVk1T8X7qwwD-Jc-T-87zCEwkxw=/159x699:1428x1546/195x98/data/photo/2025/07/27/688624154bcb8.png",
          "pubDate": "2025-07-27 22:02:20",
          "section": "Tren",
          "thumb": "https://asset.kompas.com/crops/NVk1T8X7qwwD-Jc-T-87zCEwkxw=/159x699:1428x1546/195x98/data/photo/2025/07/27/688624154bcb8.png",
          "title": "Gejala Khas Varian Baru Covid-19 XFG atau Stratus"
        }
      ]
    }
  }
}
```

### Key Data Fields by Source

#### CNN Indonesia
- **`title`**: News headline
- **`summary`**: Article summary/description
- **`url`**: Full article URL
- **`date.created`**: Creation timestamp
- **`date.publish`**: Publication timestamp
- **`editor`**: Editor name
- **`namakanal`**: Channel name
- **`namaparent`**: Parent category
- **`images.cover`**: Cover image URL
- **`tag`**: Comma-separated tags

#### DETIK
- **`title`**: News headline
- **`summary`**: Article summary
- **`link`**: Full article URL
- **`date.created`**: Creation timestamp
- **`date.publish`**: Publication timestamp
- **`editor`**: Editor name
- **`penulis`**: Author name
- **`city`**: Location
- **`kanal_parent_name`**: Channel name
- **`images.cover`**: Cover image URL

#### KOMPAS
- **`title`**: News headline
- **`description`**: Article description
- **`link`**: Full article URL
- **`pubDate`**: Publication date
- **`channel`**: Channel category
- **`kanal`**: Channel name
- **`section`**: Section name
- **`photo`**: Photo URL
- **`thumb`**: Thumbnail URL
- **`guid`**: Unique identifier

### Go Struct Mapping
```go
type IndonesiaNewsResponse struct {
    Items    []interface{}          `json:"items,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
    Status   string                 `json:"status,omitempty"`
    Error    string                 `json:"error,omitempty"`
    Source   string                 `json:"source,omitempty"`
    Query    string                 `json:"query,omitempty"`
    Params   interface{}            `json:"params,omitempty"`
}
```

---

## 📰 Real-Time News Data API

### API Configuration
- **API Name**: Real-Time News Data API (via RapidAPI)
- **Host**: `real-time-news-data.p.rapidapi.com`
- **Status**: ✅ **FULLY FUNCTIONAL**

### Response Structure
**Endpoint**: `/search?query={query}&limit={limit}&time_published={time}&country={country}&lang={lang}`

```json
{
  "status": "OK",
  "request_id": "8513c4a0-7f82-4dbe-98e9-515d74d5b46f",
  "data": [
    {
      "article_id": "CBMitAFBVV95cUxQenBNQWZ3eVJaa2s3RklpV1Z1LTg5bVRTbHNWZ3Z2Q3ZhcXY5cVI0U1k0VWVBa0xSYkNrYmxURDF3THZWczRJd1JxNy1MeGpNTl9OZlhKQkZSOEJlZUJUYVVhZkVuSS1FQ0RzR1ExcmRydzFkeWlBRnFCTmVlSk1uMTZGTFhXbzhwS1h4WDVpeVNFSHJOQnpYNTFvRmEzNnB0Z2RWVzVwTV9QT2dwRFNCdGNHNzE",
      "title": "Masalah Kendali Optimal Model Matematika Epidemi COVID-19 dengan Parameter Fuzzy di Indonesia",
      "link": "https://unair.ac.id/masalah-kendali-optimal-model-matematika-epidemi-covid-19-dengan-parameter-fuzzy-di-indonesia/",
      "snippet": "kami menganalisis model pengendalian optimal epidemi COVID-19 di Indonesia dengan Mempertimbangkan Kebijakan Pemerintah.",
      "photo_url": "https://unair.ac.id/wp-content/uploads/2025/08/5ed517a5b7a42-1024x683.webp",
      "thumbnail_url": "https://news.google.com/api/attachments/CC8iL0NnNURaV1YwVVV0V2FFeG9UWFV4VFJERUF4aW1CU2dLTWdrWlFwSkpubWs5YUFJ=-w200-h200-p-df-rw",
      "published_datetime_utc": "2025-08-25T03:56:00.000Z",
      "authors": [],
      "source_url": "https://unair.ac.id",
      "source_name": "unair.ac.id",
      "source_logo_url": "https://encrypted-tbn1.gstatic.com/faviconV2?url=https://unair.ac.id&client=NEWS_360&size=256&type=FAVICON&fallback_opts=TYPE,SIZE,URL",
      "source_favicon_url": "https://encrypted-tbn1.gstatic.com/faviconV2?url=https://unair.ac.id&client=NEWS_360&size=96&type=FAVICON&fallback_opts=TYPE,SIZE,URL",
      "source_publication_id": "CAAqJQgKIh9DQklTRVFnTWFnMEtDM1Z1WVdseUxtRmpMbWxrS0FBUAE",
      "related_topics": []
    }
  ]
}
```

### Key Data Fields
- **`article_id`**: Unique article identifier
- **`title`**: News headline
- **`link`**: Full article URL
- **`snippet`**: Article summary/description
- **`photo_url`**: High-resolution image URL
- **`thumbnail_url`**: Thumbnail image URL
- **`published_datetime_utc`**: Publication timestamp in UTC
- **`authors`**: Array of author names
- **`source_name`**: News source name (e.g., "unair.ac.id", "TIMES Indonesia")
- **`source_url`**: Source website URL
- **`source_logo_url`**: Source logo image URL
- **`source_favicon_url`**: Source favicon URL
- **`source_publication_id`**: Google News publication ID
- **`related_topics`**: Array of related topic tags

### Go Struct Mapping
```go
type RealTimeNewsResponse struct {
    Status    string      `json:"status"`
    RequestID string      `json:"request_id"`
    Data      interface{} `json:"data,omitempty"`
    Error     interface{} `json:"error,omitempty"` // Can be string or object
    Query     string      `json:"query,omitempty"`
    Country   string      `json:"country,omitempty"`
    Lang      string      `json:"lang,omitempty"`
    Limit     int         `json:"limit,omitempty"`
}
```

### Supported Parameters
- **`query`**: Search term (e.g., "covid 19")
- **`limit`**: Number of results (1-100)
- **`time_published`**: Time filter ("anytime", "past_24h", "past_week", "past_month")
- **`country`**: Country code (e.g., "ID", "US", "GB")
- **`lang`**: Language code (e.g., "id", "en")

---

## 🔄 ETL Flow Summary

### Instagram Flow
1. **Search**: `GetHashtagMedia("covid19", "")` → Returns array of posts
2. **Extract**: Parse posts array for metadata (likes, comments, captions, user info)
3. **Transform**: Clean and structure post data
4. **Load**: Store in database

### YouTube Flow
1. **Search**: `SearchVideos("covid 19", "en", "US")` → Returns video results
2. **Extract Video IDs**: Parse `contents[].video.videoId` from search results
3. **Get Comments**: `GetVideoComments(videoId)` for each video
4. **Transform**: Combine video metadata with comments
5. **Load**: Store in database

### Indonesian News Flow
1. **Search**: `SearchNews(source, "covid", params)` → Returns news items
2. **Extract**: Parse source-specific response structure for news metadata
3. **Transform**: Clean and structure news data
4. **Load**: Store in database

### Real-Time News Flow
1. **Search**: `SearchNews(query, country, lang, limit, timePublished)` → Returns news articles
2. **Extract**: Parse articles array for metadata (title, snippet, link, source, etc.)
3. **Transform**: Clean and structure news data
4. **Load**: Store in database

---

## 📊 Data Extraction Status

| Extractor | Status | Data Retrieved | Notes |
|-----------|--------|----------------|-------|
| **Instagram** | ✅ Working | 30 posts with full metadata | Array response structure handled |
| **YouTube** | ✅ Working | Videos + Comments | Two-step process: search then comments |
| **Indonesian News** | ✅ Working | 3/3 sources working | CNN, DETIK, KOMPAS all functional |
| **Real-Time News** | ✅ Working | News articles | Real-time news data retrieval |

---

## 🚀 Next Steps

1. **✅ All Extractors Tested and Working**
2. **Implement data transformation logic**
3. **Set up database storage**
4. **Create unified ETL pipeline**

---

*Last Updated: Based on testing results from Instagram, YouTube, Indonesian News, and Real-Time News extractors*
