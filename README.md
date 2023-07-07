# gguuse-streams

Just a simple twitch chatbot that I'm writing for learning purposes.

# Functionality

It has two basic commands:
- `help`: shows all available commands.
- `setmessage`: creates new command. `setmessage` is not available for anyone except for user of bot whose nick stored in enviroment variable "NICK".

Bot stores all commands except `help` and `setmessage` in commands.json file.
