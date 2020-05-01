# Quack

CLI for tweet-sized private journal entries.

### What is it?

Let me start by saying that this is mostly a coding exercise for me, and the
"usefulness" of this program might be neglible. But if I were to market Quack, I
would say something like "Quack is a secure, lightweight, private journaling
application that aims be cloud-agnostic." Others might call it a CRUD app. 

It works like this. You install Quack and run it with some variables present in
the environment. One of those is your QUACKWORD, which is your keys to the
castle. When you enter a message, it is encrypted, and can only be read with the
right QUACKWORD. If you include no other variables, each message is stored as a
file in `$HOME/.quack`. If however, you include the credentials for an S3
bucket, your messages will be stored there, and you'll be able to read or write
to them (with Quack) from anywhere.

Oh, and your messages can't be longer than 280 characters.    

### Installation

By far the easiest way to install Quack is with Docker.

1. Pull the image: `docker pull
   docker.pkg.github.com/jonathanwthom/quack/quack:latest`

2. Create a file to include your environment variables, e.g. `.env`. It could
   look like this: 
    ```
    AWS_ACCESS_KEY_ID=<access-key-id-goes-here>
    AWS_SECRET_ACCESS_KEY=<secret-access-key-goes-here>
    S3_BUCKET_NAME=<s3-bucket-name-goes-here>
    S3_BUCKET_REGION=<s3-bucket-region-goes-here>
    QUACKWORD=<quackword-goes-here>
    ```

3. Run an interactive shell in a container:
    ```
    docker run -it --env-file .env docker.pkg.github.com/jonathanwthom/quack/quack:latest /bin/sh`
    ```

    If you don't want to store your credentials in just a plain file, you can
    pass them in on the fly to the container:
    ```
    docker run -it -e QUACKWORD=my-quackword ...etc 
    ```
 
4. If you want to run Docker, and not use the cloud, you'll need to persist your
   messages to a volume.

   ```
   docker run -it -e QUACKWORD=quacky -v $HOME/.quack:/root/quack docker.pkg.github.com/jonathanwthom/quack/quack:latest /bin/sh
   ``` 

5. The `quack` executable will be loaded. Run `quack -h` to see all options for
   usage. Current options are:
   ```
   delete      Delete an entry
   help        Help about any command
   new         Create a new entry
   read        Read all entries
        -s, --search string   Search entries by text
        -v, --verbose         Display entries in verbose mode
   ```
   You can add `-h` to any command to read more, e.g. `quack read -h`

### Development

After cloning the repo, run `go build ./...` and then `./quack <some-command>`.
Tests can be run with `go test ./...`. You can either set your environment
variables within your shell session/environment, or within a `.env` file at the
root of the project. 

### License

Apache License 2.0

