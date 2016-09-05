# Download image from URL

Utility that does the following:

- Reads the instagram object stored in the database;
- Downloads the image;
- Updates the database with the filename of the downloaded / created image;
- Creates the different sizes of the images (large, medium and thumbnail).

## Usage

`go install download_image_from_url`
and
`ENV_VARS_FILE="/.vars_file" $GOPATH/bin/download_image_from_url` (or whatever you use)
