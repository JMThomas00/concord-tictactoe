# Tic-tac-toe

Tic-tac-toe for the terminal and for
[Concord](https://github.com/JMThomas00/Concord) channels, from the same
program. In Concord it's a little arcade cabinet: a title screen where the
computer plays itself, a menu, pixel-art pieces in your theme's colours,
chiptune sounds, a Hall of Fame, and **Gold Stars** to spend on new piece
sets and boards. A draw is a cat's game, with the cat.

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
4. Select the channel and press **Tab** (or click it) so your keys go to the
   game. Press **Enter** on the title screen, then pick from the menu.

## The arcade

Everyone who opens the channel starts on the title screen. The menu:

- **1 PLAYER VS CPU** `◂ NORMAL ▸`: a game of your own against the computer
  (←/→ picks easy, normal or hard; the hard one never loses). Leave and come
  back, and it picks up where you were.
- **TAKE A SEAT** (seats channels) sits you at the channel's table; in a
  challenge channel it's **2 PLAYERS** (challenge someone) and **WATCH**, in a
  private one **NEW GAME** and **YOUR GAMES**.
- **SETS**: your piece set and board. 12 piece sets (classic, chunky, grapes
  and leaves, hearts and stars ... up to tiny cars, a duck and a toaster, and
  a sock and sandal) and 7 boards (notebook, chalkboard, wooden, neon, picnic
  blanket, vineyard trellis). Everyone sees the game in their own.
- **HALL OF FAME**, **HOW TO PLAY** and **OPTIONS** (your sound and effects).

You earn a **Gold Star** for each new achievement, every three wins in a row
and every ten games. Spend one in SETS on a locked set or board: you're
offered three and pick one.

After a game, **Enter** shows the results; Enter again asks for a rematch,
which starts once both players have. **Esc** goes back a screen (on the title
screen it gives the keyboard back to Concord). A pane smaller than 64 x 24
gets the plain board instead.

To update later: select it in **Server Settings → Plugins**, press **U**, then
Enter. A failed update rolls back by itself.

## Playing

- **Arrow keys** (or h j k l) move the cursor; **Enter** or **Space** places your mark.
- **1–9** place it directly, numbered like a phone keypad read left to right:
  `1 2 3` on top, `7 8 9` at the bottom.
- **?** shows the rules and keys.
- **M** opens the table menu: sit, stand, resign, rematch, play the computer
  (playing standalone, Tab does too).
- **Esc** goes back to the menu (or closes the help screen).

The last move is marked with a dot, and a winning line is highlighted.

## Layout

- `engine/`: the rules and the computer player (minimax at the hardest level).
- `game/`: connects the engine to the Concord SDK's table kit, the board you
  play on (`board.go`), the piece sets and boards (`sets.go`, `draw.go`) and the
  arcade's personality (`arcade.go`: attract mode, the cat).
- `client/`: the pictures and sounds Concord sends to members' clients.
- `tools/gen.go`: `go run tools/gen.go` regenerates `client/` and
  `game/pieces_gen.go` (a picture of each piece in every built-in theme's red and cyan, for the
  plain board) and the arcade sound kit.
- `release.go`: `go run release.go` builds the release zips Concord installs.

Tag a version (`git tag v0.1.0 && git push --tags`) and the workflow publishes them.

## License

MIT License — see LICENSE file for details.
