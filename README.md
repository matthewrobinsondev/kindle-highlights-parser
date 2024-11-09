# What
A lightweight kindle clippings parser, which will correlate to markdown files.

# Why
When reading my kindle I highlight a lot of information for a variety of reason, but the main one is usually because the part will resonate with me. I use obsidian for my notes pretty much every day, what I wanted to be able to do was take my kindle highlights & put them straight into my obsidian vault, with some extra parsing. Having to manually get the clippings file from the kindle sucks tho, sort it out Jeff.

# How
1. Plug your kindle in & copy the file from `documents/My Clippings.txt`
2. Place it in this repository
3. Create a `config.toml`
4. Add a `notes_directory="Documents/where ever your notes are"`
5. Run the application:
    - `go build && ./kindle-notes-parser`
    - `go run .`

# TODOS
I have a bunch in the code, I'll probably move to issues but there reminders for refactors once its working.
