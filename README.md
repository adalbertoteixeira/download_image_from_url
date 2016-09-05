# Download image from URL

Utility that does the following:

- Reads the instagram object stored in the database;
- Downloads the image;
- Updates the database with the filename of the downloaded / created image;
- Creates the different sizes of the images (large, medium and thumbnail).

# install

Along with the install and build, you'll need to have a file with the following vars:

```
  FILE_PATH=/path/to/the/directory/where/the/subfolders/will/be/created/
  DATABASE_NAME=dbname
  DATABASE_USERNAME=dbuser
  DATABASE_PASSWORD=dbpassword
  DATABASE_HOST=localhost
  DATABASE_PORT=5432
  IMAGEMAGICK_PATH=/path/to/imagemagick/convert
```

## Usage

```go install download_image_from_url```

and

```ENV_VARS_FILE="/.vars_file" $GOPATH/bin/download_image_from_url``` (or whatever you use)
