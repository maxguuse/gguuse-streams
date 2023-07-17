# gguuse-streams

Just a simple twitch chatbot that I'm writing for learning purposes.

# Functionality

It has following basic commands:
- `help`: shows all available commands.
- `setmessage`:
  - Usage: !setmessage <command> <message> - updates `command` with `message`, if `message` is empty, deletes command
- `newannouncement`:
  - Usage: !newannouncement <id> <repetition_interval> <message> - updates `id`-announcement with `message`. Announcement repeats every `repetition_interval` minutes. `id` can be any combination of any symbols.
- `stopannouncement`:
  - Usage: !stopannouncement <id> - stops announcement `id`.

Developer can predefine some commands right in code. <br/>
Commands created with `setmessage` stored in directory `json_commands` in files named `<channel>_commands.json`. <br/>
Announcements stored in directory `json_announcements` in files named `<channel>_announcements.json`. <br/>
