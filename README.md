# go-disposable-email

This small tool is to generate a disposable temporary email using different domains quickly in 1 second, receive your mail and attachments, and save them to a folder!

## what is disposable emails?
A disposable email is a unique email address that is temporary. It expires after a set amount of time or a set number of uses!

## How to Use

1. Clone the repository to your local machine:
    ```bash
    git clone https://github.com/go-nerds/go-disposable-email.git
    ```

2. Build the project:
   ```bash
   go build
   ```
3. Run the executable file:
   ```bash
    go run .
    Use the arrow keys to navigate: ↓ ↑ → ←
    ? Select Domain:
    > 1secmail.com
        1secmail.org
        1secmail.net
        kzccv.com
    ↓   qiott.com

   ```
4. The tool will display for you the temporary email address and will watch every 5 seconds for any inbox message.
5. After checking the messages and attachments, it will save them in a file in the same selected domain.
## LICENSE

This project is licensed under the MIT License. See the [LICENSE](https://github.com/go-nerds/go-disposable-email/blob/main/LICENSE) file for details.
