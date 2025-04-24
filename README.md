# neopoll-cli
A simple polling CLI app

## Commands

| Command           | Description                                                          |
|-------------------|----------------------------------------------------------------------|
| `completion`      | Generate the autocompletion script for the specified shell           |
| `generate-sample` | Create a sample poll JSON file                                       |
| `help`            | Help about any command                                               |
| `login`           | Request an OTP & verify (prompts for email and, if needed, token)    |
| `logout`          | Delete saved session                                                 |
| `start`           | Start a new poll session (prompts for password)                      |
| `next`            | Proceed to the next poll item                                        |
| `reveal`          | Reveal the answers to the current question                           |
| `stop`            | Stop the poll session and save the results locally                   |
| `refresh-token`   | Refresh saved session (uses stored refresh token)                    |

---

## Quickstart

1. **Authenticate**  
   ```bash
   # Primary flow: prompt for email → send OTP → prompt for token → save session
   $ neopoll login
   
   # If you already have an OTP:
   $ neopoll login --email you@example.com --token 123456
   ```
   _This saves your Supabase session to `~/.neopoll/session.json`._

2. **Generate a sample poll file**  
   ```bash
   # writes to ./sample-poll.json by default
   $ neopoll generate-sample
   
   # or choose a path
   $ neopoll generate-sample --output=my-poll.json
   ```
   Edit the resulting JSON to define your own questions and options.

3. **Start the poll**  
   ```bash
   $ neopoll start my-poll.json
   ```
   You’ll be prompted to choose a “room” password that participants must enter to join.  
   The command will print your poll URL (e.g. `https://poll.app/abc-123`).

4. **Control the live poll**  
   - Advance to the next question:  
     ```bash
     $ neopoll next
     ```
   - Reveal the current answer:  
     ```bash
     $ neopoll reveal
     ```
   - Stop the poll and pull down results:  
     ```bash
     # defaults to ./results.json
     $ neopoll stop

     # or specify your own output path
     $ neopoll stop --output=my-results.json
     ```
   Participants’ votes will be saved locally for analysis.

---

### Other Handy Commands

- **`neopoll logout`** — delete your saved session (you’ll have to `login` again).  
- **`neopoll refresh-token`** — renew your session if your access token expires.  
- **`neopoll completion [bash|zsh|fish|powershell]`** — install shell auto-completion.  

For full details on any command:
```bash
$ neopoll help <command>
```