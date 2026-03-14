# Simply File Share

A simple, lightweight file sharing web application built with Go. Upload, download, and manage files through a clean web interface.

## Features

- **File Upload**: Upload files with automatic timestamping to prevent naming conflicts
- **File Download**: Download files with download count tracking
- **File Management**: View file metadata including size, upload time, and download count
- **File Deletion**: Delete files with confirmation prompts
- **Web Interface**: Clean, responsive UI built with Bootstrap
- **REST API**: JSON API for programmatic access
- **Persistent Storage**: File metadata stored in JSON format
- **Embedded Templates**: No external template files needed

## Prerequisites

- Go 1.17 or higher
- Web browser for accessing the interface

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/simply-file-share.git
cd simply-file-share
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o file-share main.go
```

## Usage

### Running the Server

Start the server with default port (8022):
```bash
./file-share
```

Or specify a custom port:
```bash
./file-share -p 8080
```

The server will start and display the port it's running on.

### Accessing the Web Interface

Open your web browser and navigate to `http://localhost:8022` (or your specified port).

### Web Interface Features

- **Upload Files**: Click "Choose file" and select a file, then click "Upload"
- **View Files**: See all uploaded files in a table with metadata
- **Download Files**: Click the "Download" button next to any file
- **Delete Files**: Click the "Delete" button (with confirmation prompt)

## API Endpoints

The application provides a REST API for programmatic access:

### GET /
Serves the main web interface.

### POST /upload
Upload a file.

**Parameters:**
- `file`: Multipart form file

**Response:** Redirects to home page on success.

### GET /download/{filename}
Download a specific file.

**Parameters:**
- `filename`: The name of the file to download

### GET /files
Get JSON list of all uploaded files.

**Response:**
```json
[
  {
    "name": "060102150405_filename.txt",
    "size": 1024,
    "upload_time": "2023-01-02T15:04:05Z",
    "download_count": 5
  }
]
```

### POST /delete/{filename}
Delete a specific file.

**Parameters:**
- `filename`: The name of the file to delete

**Response:** Redirects to home page on success.

## File Storage

- **Files**: Stored in the `myfile/` directory
- **Metadata**: Stored in `data.json` in the root directory
- **Naming**: Files are prefixed with timestamp (YYMMDDHHMMSS_) to prevent conflicts

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router and URL matcher
- [Bootstrap](https://getbootstrap.com/) - CSS framework for the web interface

## Development

To modify the web interface, edit `templates/index.html`. The template is embedded in the binary using Go's `embed` package.

## License

See LICENSE.txt for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Security Notes

- This is a basic file sharing application intended for local/private use
- No authentication or authorization is implemented
- Files are served directly from the filesystem
- Consider implementing security measures for production use</content>
<parameter name="filePath">/home/aisuma/Documents/project/simply-file-share/README.md