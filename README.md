```sh
export GPG_TTY=`tty`

echo -n 'password' > /Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceprivtest-pass.txt


gpg --output file.txt --decrypt file.txt.gpg 
```