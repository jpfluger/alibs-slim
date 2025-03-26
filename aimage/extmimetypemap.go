package aimage

import (
	"github.com/jpfluger/alibs-slim/autils"
	"mime"
	"strings"
)

type ExtMimeType struct {
	Ext        string
	Mime       string
	CleanedExt string
	Tags       []string
}

func (emt ExtMimeType) GetCleanedExt() string {
	if emt.CleanedExt == "" {
		return emt.Ext
	}
	return emt.CleanedExt
}

func (emt ExtMimeType) GetTags() []string {
	var arr []string
	// The Mime has a "/". Add a qualifying tag for the prefix
	parts := strings.SplitN(emt.Mime, "/", 2)
	if len(parts) > 0 {
		arr = append(arr, parts[0]) // Add the prefix (e.g., "image" for "image/jpeg")
	}
	// Append existing tags if any
	if emt.Tags != nil && len(emt.Tags) > 0 {
		arr = append(arr, emt.Tags...)
	}
	return arr
}

type ExtMimeTypeMap map[ImageType]ExtMimeType

func (emap ExtMimeTypeMap) FindExtMime(ext string) *ExtMimeType {
	emt, exists := emap[ImageType(autils.StripExtensionPrefix(ext))]
	if !exists {
		return nil
	}
	return &emt
}

func (emap ExtMimeTypeMap) FindMime(ext string) string {
	emt, exists := emap[ImageType(autils.StripExtensionPrefix(ext))]
	if !exists {
		return ""
	}
	return emt.Mime
}

func (emap ExtMimeTypeMap) GetCleanedExt(ext string) string {
	emt, exists := emap[ImageType(autils.StripExtensionPrefix(ext))]
	if !exists {
		return ""
	}
	return emt.GetCleanedExt()
}

// customExtMimeTypes holds the map of file extensions to their corresponding MIME types.
var customExtMimeTypes ExtMimeTypeMap

func init() {
	customExtMimeTypes = ExtMimeTypeMap{
		// Images
		"jpg":  {Ext: "jpg", Mime: "image/jpeg"},
		"jpeg": {Ext: "jpeg", Mime: "image/jpeg", CleanedExt: "jpg"},
		"png":  {Ext: "png", Mime: "image/png"},
		"gif":  {Ext: "gif", Mime: "image/gif"},
		"bmp":  {Ext: "bmp", Mime: "image/bmp"},
		"webp": {Ext: "webp", Mime: "image/webp"},
		"tiff": {Ext: "tiff", Mime: "image/tiff"},
		"tif":  {Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"},
		"svg":  {Ext: "svg", Mime: "image/svg+xml"},
		"ico":  {Ext: "ico", Mime: "image/vnd.microsoft.icon"},
		"heic": {Ext: "heic", Mime: "image/heic"},
		"heif": {Ext: "heif", Mime: "image/heif"},

		// Video
		"mp4":  {Ext: "mp4", Mime: "video/mp4"},
		"m4v":  {Ext: "m4v", Mime: "video/x-m4v"},
		"amp4": {Ext: "amp4", Mime: "video/mp4"}, // not widely used. mp4 standard recognizes video.
		"mov":  {Ext: "mov", Mime: "video/quicktime"},
		"avi":  {Ext: "avi", Mime: "video/x-msvideo"},
		"wmv":  {Ext: "wmv", Mime: "video/x-ms-wmv"},
		"flv":  {Ext: "flv", Mime: "video/x-flv"},
		"mkv":  {Ext: "mkv", Mime: "video/x-matroska"},
		"webm": {Ext: "webm", Mime: "video/webm"},
		"3gp":  {Ext: "3gp", Mime: "video/3gpp"},
		"3g2":  {Ext: "3g2", Mime: "video/3gpp2"},

		// Audio
		"mp3":  {Ext: "mp3", Mime: "audio/mpeg"},
		"m4a":  {Ext: "m4a", Mime: "audio/mp4"},
		"aac":  {Ext: "aac", Mime: "audio/aac"},
		"wav":  {Ext: "wav", Mime: "audio/wav"},
		"flac": {Ext: "flac", Mime: "audio/flac"},
		"ogg":  {Ext: "ogg", Mime: "audio/ogg"},
		"opus": {Ext: "opus", Mime: "audio/opus"},
		"amr":  {Ext: "amr", Mime: "audio/amr"},
		"aiff": {Ext: "aiff", Mime: "audio/aiff"},
		"wma":  {Ext: "wma", Mime: "audio/x-ms-wma"},

		// Documents
		"pdf":  {Ext: "pdf", Mime: "application/pdf", Tags: []string{"document"}},
		"doc":  {Ext: "doc", Mime: "application/msword", Tags: []string{"document", "ms-office"}},
		"docx": {Ext: "docx", Mime: "application/vnd.openxmlformats-officedocument.wordprocessingml.document", Tags: []string{"document", "ms-office"}},
		"xls":  {Ext: "xls", Mime: "application/vnd.ms-excel", Tags: []string{"document", "ms-office"}},
		"xlsx": {Ext: "xlsx", Mime: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Tags: []string{"document", "ms-office"}},
		"ppt":  {Ext: "ppt", Mime: "application/vnd.ms-powerpoint", Tags: []string{"document", "ms-office"}},
		"pptx": {Ext: "pptx", Mime: "application/vnd.openxmlformats-officedocument.presentationml.presentation", Tags: []string{"document", "ms-office"}},
		"txt":  {Ext: "txt", Mime: "text/plain", Tags: []string{"text", "document"}},
		"rtf":  {Ext: "rtf", Mime: "application/rtf", Tags: []string{"document"}},
		"csv":  {Ext: "csv", Mime: "text/csv", Tags: []string{"document"}},

		// OpenOffice / LibreOffice
		"odt": {Ext: "odt", Mime: "application/vnd.oasis.opendocument.text", Tags: []string{"document", "open-office"}},
		"ods": {Ext: "ods", Mime: "application/vnd.oasis.opendocument.spreadsheet", Tags: []string{"document", "open-office"}},
		"odp": {Ext: "odp", Mime: "application/vnd.oasis.opendocument.presentation", Tags: []string{"document", "open-office"}},
		"odg": {Ext: "odg", Mime: "application/vnd.oasis.opendocument.graphics", Tags: []string{"document", "open-office"}},
		"ott": {Ext: "ott", Mime: "application/vnd.oasis.opendocument.text-template", Tags: []string{"document", "open-office"}},
		"ots": {Ext: "ots", Mime: "application/vnd.oasis.opendocument.spreadsheet-template", Tags: []string{"document", "open-office"}},
		"otp": {Ext: "otp", Mime: "application/vnd.oasis.opendocument.presentation-template", Tags: []string{"document", "open-office"}},

		// Other Document Types
		"epub":     {Ext: "epub", Mime: "application/epub+zip", Tags: []string{"book-pub"}},
		"mobi":     {Ext: "mobi", Mime: "application/x-mobipocket-ebook", Tags: []string{"book-pub"}},
		"ps":       {Ext: "ps", Mime: "application/postscript", Tags: []string{"graphics"}},
		"ai":       {Ext: "ai", Mime: "application/postscript", Tags: []string{"graphics"}},
		"afdesign": {Ext: "afdesign", Mime: "application/x-affinity-designer", Tags: []string{"graphics"}},
		"afphoto":  {Ext: "afphoto", Mime: "application/x-affinity-photo", Tags: []string{"graphics"}},
		"afpub":    {Ext: "afpub", Mime: "application/x-affinity-publisher", Tags: []string{"graphics"}},

		// Web/Programming Languages
		"html": {Ext: "html", Mime: "text/html", Tags: []string{"web-template"}},
		"htm":  {Ext: "htm", Mime: "text/html"},
		"css":  {Ext: "css", Mime: "text/css"},
		"js":   {Ext: "js", Mime: "application/javascript"},
		"json": {Ext: "json", Mime: "application/json"},
		"xml":  {Ext: "xml", Mime: "application/xml"},
		"yaml": {Ext: "yaml", Mime: "application/x-yaml", CleanedExt: "yaml"},
		"yml":  {Ext: "yml", Mime: "application/x-yaml", CleanedExt: "yaml"},
		"php":  {Ext: "php", Mime: "application/x-httpd-php"},
		"java": {Ext: "java", Mime: "text/x-java-source"},
		"py":   {Ext: "py", Mime: "text/x-script.python"},
		"rb":   {Ext: "rb", Mime: "application/x-ruby"},
		"go":   {Ext: "go", Mime: "text/x-go"},
		"c":    {Ext: "c", Mime: "text/x-c"},
		"cpp":  {Ext: "cpp", Mime: "text/x-c++"},
		"h":    {Ext: "h", Mime: "text/x-c"},
		"hpp":  {Ext: "hpp", Mime: "text/x-c++"},
		"ts":   {Ext: "ts", Mime: "application/typescript"},
		"tsx":  {Ext: "tsx", Mime: "application/typescript"},
		"md":   {Ext: "md", Mime: "text/markdown"},

		// Go HTML Templates
		"gohtml": {Ext: "gohtml", Mime: "text/html", Tags: []string{"web-template"}},

		// SQL Scripts
		"sql": {Ext: "sql", Mime: "application/sql", Tags: []string{"database"}},
	}
}

func GetCleanedExt(ext string) string {
	return customExtMimeTypes.GetCleanedExt(ext)
}

func GetMimeType(ext string) string {
	return customExtMimeTypes.FindMime(ext)
}

// SetExtMimeType allows higher-order apps to add or update a MIME type in the global MimeTypes map.
func SetExtMimeType(ext, mimeType string, tags []string, cleanedExt ...string) {
	ext = autils.StripExtensionPrefix(ext)
	if ext == "" {
		return
	}

	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	cleaned := ext
	if len(cleanedExt) > 0 {
		cleaned = autils.StripExtensionPrefix(cleanedExt[0])
	}

	customExtMimeTypes[ImageType(ext)] = ExtMimeType{
		Ext:        ext,
		Mime:       mimeType,
		CleanedExt: cleaned,
		Tags:       tags,
	}
}

// CleanMimeType normalizes the MIME type based on the file extension or returns the detected MIME type.
func CleanMimeType(mimeType, fileNameOrExt string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	ext := autils.StripExtensionPrefix(fileNameOrExt)

	if ext != "" {
		emt := customExtMimeTypes.FindExtMime(ext)
		if emt != nil {
			mimeType = emt.Mime
			ext = emt.GetCleanedExt()
		}
	}

	// Use Go's built-in MIME type detection for file extensions
	if mimeType == "" && ext != "" {
		mimeType = mime.TypeByExtension("." + ext)
	}

	return mimeType
}
