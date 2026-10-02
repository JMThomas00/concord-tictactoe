# Tic-tac-toe

Tic-tac-toe for the terminal and for
[Concord](https://github.com/JMThomas00/Concord) channels, from the same
program. In Concord the X and O are pictures drawn in your theme's colors,
with a sound for each move, a win and a draw.

## Play it on your own computer

**Download:** from the [Releases](https://github.com/JMThomas00/concord-tictactoe/releases)
page, get the zip for your system (`concord-tictactoe_windows_amd64.zip`,
`concord-tictactoe_darwin_arm64.zip` for Apple silicon, `concord-tictactoe_linux_amd64.zip`, ...),
unzip it, and run the program inside from a terminal:

```sh
./concord-tictactoe          # Windows: .\concord-tictactoe.exe
```

On macOS, if it's blocked as being from an unidentified developer, run
`xattr -d com.apple.quarantine concord-tictactoe` once. **Or with Go installed:**
`go install github.com/JMThomas00/concord-tictactoe@latest`, then run `concord-tictactoe`.

It starts with a menu: two players on one keyboard, against the computer at
three levels (the hardest never loses), or over the network. For a network
game, one player hosts and is shown their address and a 6-character code; the
other chooses join and types both.

## Play it on a Concord server

You need to be the server owner, or have the **Manage Plugins** permission.

1. In Concord, open **Server Settings → Plugins** and press **I** (install).
2. Type `JMThomas00/concord-tictactoe` and press Enter. Concord downloads the latest
   release for the server's own system, verifies it, and starts it: no
   restart, no files to edit.
3. Open **Server Settings → Channels**, create a channel, and choose
   **Tic-tac-toe** as its type. Its options:
   - **Seating**: *seats* (one board; sit down with M, everyone else
     watches), *challenge* (a lobby where members challenge each other), or
     *private* (your own games with opponents you pick).
   - **Allow spectators**, **Computer opponent**, and **Computer strength**
     (easy, normal or hard).
4. Select the channel and press **Tab** (or click the board) so your keys go
   to the game. **M** opens the table menu: sit down, play the computer,
   resign, rematch. **Esc** gives the keyboard back to Concord, and **Tab**
   moves on to the member list.

The pieces are pictures in terminals that can show them (Windows Terminal,
iTerm2, WezTerm, Kitty, Ghostty, foot, and others). Elsewhere, and when
playing standalone, they're drawn with block characters in the same colors.
Either way they follow the viewer's Concord theme, including custom ones.

To update later: select it in **Server Settings → Plugins**, press **U**, then
Enter. A failed update rolls back by itself.

## Playing

- **Arrow keys** (or h j k l) move the cursor; **Enter** or **Space** places your mark.
- **1–9** place it directly, numbered like a phone keypad read left to right:
  `1 2 3` on top, `7 8 9` at the bottom.
- **?** shows the rules and keys.
- **M** opens the table menu: sit, stand, resign, rematch, play the computer
  (playing standalone, Tab does too).
- **Esc** hands the keyboard back to Concord (or closes the help screen).
  **Tab** moves on to the member list.

The last move is marked with a dot, and a winning line is highlighted.

## Layout

- `engine/`: the rules and the computer player (minimax at the hardest level).
- `game/`: connects the engine to the Concord SDK's table kit, the board you
  play on, and which piece picture matches the viewer's theme.
- `client/`: the pictures and sounds Concord sends to members' clients.
- `tools/gen.go`: `go run tools/gen.go` regenerates `client/` and
  `game/pieces_gen.go` (a picture of each piece in every built-in theme's red and cyan).
- `release.go`: `go run release.go` builds the release zips Concord installs.

Tag a version (`git tag v0.1.0 && git push --tags`) and the workflow publishes them.
