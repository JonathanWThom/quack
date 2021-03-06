# Quack

CLI for tweet-sized private journal entries.

![Build](https://github.com/JonathanWThom/quack/workflows/Build/badge.svg)

## What is it?

Let me start by saying that this is mostly a coding exercise for me, and the
"usefulness" of this program might be neglible. But if I were to market Quack, I
would say something like "Quack is a secure, lightweight, private journaling
application that aims to be cloud-agnostic." Others might call it a CRUD app. 

It works like this. You install Quack and run it with some variables present in
the environment. One of those is your QUACKWORD, which is your key to the
castle. When you enter a message, it is encrypted, and can only be read with the
right QUACKWORD. If you include no other variables, each message is stored as a
file in `$HOME/.quack`. If however, you include the credentials for an S3 or
GCS Bucket, your messages will be stored there, and you'll be able to read or write
to them (with Quack) from anywhere.

Oh, and your messages can't be longer than 280 characters.

## Installation

_By far the easiest way to install Quack is with Docker._

1. Pull the image: `docker pull
   docker.pkg.github.com/jonathanwthom/quack/quack:latest`

2. Create a file to include your environment variables, e.g. `.env`.

    For AWS storage:
    ```
    QUACK_AWS_ACCESS_KEY_ID=<access-key-id-goes-here>
    QUACK_AWS_SECRET_ACCESS_KEY=<secret-access-key-goes-here>
    QUACK_S3_BUCKET_NAME=<s3-bucket-name-goes-here>
    QUACK_S3_BUCKET_REGION=<s3-bucket-region-goes-here>
    QUACKWORD=<quackword-goes-here>
    ```
    
    For Google Cloud storage:
    ```
    QUACK_GOOGLE_BUCKET_NAME=<gcs-bucket-name-goes-here>
    QUACK_GOOGLE_APPLICATION_CREDENTIALS=<path-to-application-credentials-json-goes-here>
    QUACKWORD=<quackword-goes-here>
    ```

    For local filesystem storage:
    ```
    QUACKWORD=<quackword-goes-here>
    ```

3. Run an interactive shell in a container:
    For AWS:
    ```
    docker run -it --env-file .env docker.pkg.github.com/jonathanwthom/quack/quack:latest /bin/sh
    ```

    For Google Cloud, include credentials file as a volume:
    ```
    docker run -it --env-file .env -v $(pwd)/.google-application-credentials.json:/.google-application-credentials.json docker.pkg.github.com/jonathanwthom/quack/quack:latest /bin/sh
    ```

    If you don't want to store your credentials in just a plain file, you can
    pass them in on the fly to the container:
    ```
    docker run -it -e QUACKWORD=my-quackword ...etc 
    ```
 
4. If you want to run Docker, and not use the cloud, you'll need to persist your
   messages to a volume.

   ```
   docker run -it -e QUACKWORD=my-quackword -v $HOME/.quack:/root/quack docker.pkg.github.com/jonathanwthom/quack/quack:latest /bin/sh
   ``` 

5. The `quack` executable will be loaded. Run `quack -h` to see all options for
   usage. Current options are:
   ```
   delete      Delete an entry
   help        Help about any command
   new         Create a new entry
   quackword   Reset your QUACKWORD
   read        Read last 10 entries 
        -s, --search string   Search entries by text
        -v, --verbose         Display entries in verbose mode
        -d, --date string     Search entries by date in format:  "March 9, 2020"
        -n, --number int      Return last n entries
   ```
   You can add `-h` to any command to read more, e.g. `quack read -h`

_If you have Go installed, you can also build from source._

1. Clone the repo:
    ```
    git clone https://github.com/JonathanWThom/quack.git`
    ```

    OR

    ```
    git clone git@github.com:JonathanWThom/quack.git
    ```

2. Move into the directory and install the binary:
    ```
    cd quack
    go install
    ```

3. Add the the appropriate environment variables to your shell configuration
   file. E.g. in .zshrc, add:
    ```
    export QUACK_AWS_ACCESS_KEY_ID=<access-key-id-goes-here>
    export QUACK_AWS_SECRET_ACCESS_KEY=<secret-access-key-goes-here>
    export QUACK_S3_BUCKET_NAME=<s3-bucket-name-goes-here>
    export QUACK_S3_BUCKET_REGION=<s3-bucket-region-goes-here>
    export QUACKWORD=<quackword-goes-here>
    ```

4. Invoke `quack` as described above.

## Development

After cloning the repo, run `go build ./...` and then `./quack <some-command>`.
Tests can be run with `go test ./...`. You can either set your environment
variables within your shell session/environment, or within a `.env` file at the
root of the project. This project uses the awesome [Cobra framework](https://github.com/spf13/cobra).

The TODO list for the project can be found in [features.md](https://github.com/JonathanWThom/quack/blob/master/features.md).

## License

Apache License 2.0

