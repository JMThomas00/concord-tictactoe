package game

import "github.com/JMThomas00/Concord/sdk/arcade"

// The unlockables: piece sets (an X and an O that belong together) and
// boards, earned with Gold Stars. Each player sees the game in their own
// set on their own board.

// pieceSet is an X and an O, each 9 x 8 pixels (one rune a pixel, "." see-
// through), the colour role of each rune, and a one-character stand-in
// for boards too small for pictures.
type pieceSet struct {
	arcade.Unlockable
	x, o       []string
	xPal, oPal map[rune]string
	xMini      string
	oMini      string
	xRole      string // the mini piece's colour
	oRole      string
}

const (
	kindPieces = "pieces"
	kindBoard  = "board"
)

func set(id, name string, tier arcade.Tier, blurb string) arcade.Unlockable {
	return arcade.Unlockable{ID: id, Name: name, Tier: tier, Kind: kindPieces, Blurb: blurb}
}

var pieceSets = []pieceSet{
	{Unlockable: set("classic", "CLASSIC", arcade.Starter, "Pencil on paper."),
		x:    []string{"XX.....XX", ".XX...XX.", "..XX.XX..", "...XXX...", "...XXX...", "..XX.XX..", ".XX...XX.", "XX.....XX"},
		o:    []string{"..OOOOO..", ".OO...OO.", "OO.....OO", "OO.....OO", "OO.....OO", "OO.....OO", ".OO...OO.", "..OOOOO.."},
		xPal: map[rune]string{'X': "pink"}, oPal: map[rune]string{'O': "cyan"},
		xMini: "X", oMini: "O", xRole: "pink", oRole: "cyan"},
	{Unlockable: set("chunky", "CHUNKY", arcade.Starter, "Pressed down hard."),
		x:    []string{"XXX...XXX", ".XXX.XXX.", "..XXXXX..", "...XXX...", "..XXXXX..", ".XXX.XXX.", "XXX...XXX", "........."},
		o:    []string{"..OOOOO..", ".OOOOOOO.", "OOO...OOO", "OO.....OO", "OO.....OO", "OOO...OOO", ".OOOOOOO.", "..OOOOO.."},
		xPal: map[rune]string{'X': "pink"}, oPal: map[rune]string{'O': "cyan"},
		xMini: "X", oMini: "O", xRole: "pink", oRole: "cyan"},
	{Unlockable: set("grapes", "GRAPES & LEAVES", arcade.Common, "The house special."),
		x:    []string{"....gg...", "...PpP...", "..PpPPP..", ".PPpPPPP.", "..PPPpP..", "...PPP...", "....P....", "........."},
		o:    []string{"......gG.", "....gggG.", "..ggggGg.", ".gggGggg.", ".ggGgggg.", "gGgggg...", "G.gg.....", "........."},
		xPal: map[rune]string{'P': "purple", 'p': "hi", 'g': "green"}, oPal: map[rune]string{'g': "green", 'G': "greenD"},
		xMini: "●", oMini: "♣", xRole: "purple", oRole: "green"},
	{Unlockable: set("hearts", "HEARTS & STARS", arcade.Common, "Gold stars all round."),
		x:    []string{".HH...HH.", "HHHH.HHHH", "HHHHHHHHH", "HHHHHHHHH", ".HHHHHHH.", "..HHHHH..", "...HHH...", "....H...."},
		o:    []string{"....S....", "...SSS...", "SSSSSSSSS", ".SSSSSSS.", "..SSSSS..", "..SS.SS..", ".SS...SS.", "........."},
		xPal: map[rune]string{'H': "pink"}, oPal: map[rune]string{'S': "yellow"},
		xMini: "♥", oMini: "★", xRole: "pink", oRole: "yellow"},
	{Unlockable: set("sunmoon", "SUN & MOON", arcade.Common, "Day against night."),
		x:    []string{"S...S...S", ".S.SSS.S.", "..SSSSS..", "SSSSSSSSS", "..SSSSS..", ".S.SSS.S.", "S...S...S", "........."},
		o:    []string{"...MMM...", "..MMM....", ".MMM.....", ".MMM.....", ".MMM.....", "..MMM....", "...MMM...", "........."},
		xPal: map[rune]string{'S': "yellow"}, oPal: map[rune]string{'M': "cyan"},
		xMini: "☼", oMini: "☾", xRole: "yellow", oRole: "cyan"},
	{Unlockable: set("fruit", "CHERRY & LEMON", arcade.Common, "Sweet against sour."),
		x:    []string{"....g....", "...gg....", "..g..g...", ".g....g..", "RR....RR.", "RRR..RRR.", "RRR..RRR.", ".R....R.."},
		o:    []string{"...YYY...", "..YYYYY..", ".YYYYYYY.", "YYYYYYYYY", ".YYYYYYY.", "..YYYYY..", "...YYY...", "........."},
		xPal: map[rune]string{'g': "green", 'R': "red"}, oPal: map[rune]string{'Y': "yellow"},
		xMini: "♦", oMini: "●", xRole: "red", oRole: "yellow"},
	{Unlockable: set("swords", "SWORDS & SHIELDS", arcade.Rare, "Attack against defence."),
		x:    []string{"....W....", "....W....", "....W....", "....W....", "....W....", "..ooooo..", "....o....", "....o...."},
		o:    []string{".SSSSSSS.", "SSSSCSSSS", "SSSCCCSSS", "SSSSCSSSS", ".SSSSSSS.", "..SSSSS..", "...SSS...", "....S...."},
		xPal: map[rune]string{'W': "fg", 'o': "orange"}, oPal: map[rune]string{'S': "cyan", 'C': "fg"},
		xMini: "†", oMini: "◘", xRole: "fg", oRole: "cyan"},
	{Unlockable: set("space", "ROCKET & UFO", arcade.Rare, "Take me to your leader."),
		x:    []string{"....W....", "...WWW...", "...WRW...", "...WWW...", "..WWWWW..", ".W.WWW.W.", "...OOO...", "....O...."},
		o:    []string{"...CCC...", "..CCCCC..", "GGGGGGGGG", ".GYGYGYG.", "..GGGGG..", "..Y...Y..", ".Y.....Y.", "........."},
		xPal: map[rune]string{'W': "fg", 'R': "red", 'O': "orange"}, oPal: map[rune]string{'C': "cyan", 'G': "comment", 'Y': "yellow"},
		xMini: "↑", oMini: "◊", xRole: "fg", oRole: "cyan"},
	{Unlockable: set("snacks", "PIZZA & DONUT", arcade.Rare, "Lunch is a contest."),
		x:    []string{"OOOOOOOOO", ".YRYYYRY.", ".YYYRYY..", "..YRYYY..", "..YYYR...", "...YY....", "...Y.....", "........."},
		o:    []string{"..PPPPP..", ".PPPPPPP.", "PPWPPPWPP", "PPP...PPP", "PPP...PPP", "PPWPPPPPP", ".OOOOOOO.", "..OOOOO.."},
		xPal: map[rune]string{'O': "orange", 'Y': "yellow", 'R': "red"}, oPal: map[rune]string{'P': "pink", 'W': "fg", 'O': "orange"},
		xMini: "▼", oMini: "◎", xRole: "yellow", oRole: "pink"},
	{Unlockable: set("cars", "TINY CARS", arcade.Legendary, "Straight from the garage."),
		x:     []string{".........", ".........", "...pP....", "..PPPP...", "CCCCCCCCy", "CCCCCCCCC", ".kk...kk.", "........."},
		o:     []string{".........", ".........", "....Pp...", "...PPPP..", "yCCCCCCCC", "CCCCCCCCC", ".kk...kk.", "........."},
		xPal:  map[rune]string{'P': "purple", 'p': "hi", 'C': "pink", 'y': "yellow", 'k': "tire"},
		oPal:  map[rune]string{'P': "purple", 'p': "hi", 'C': "cyan", 'y': "yellow", 'k': "tire"},
		xMini: "▶", oMini: "◀", xRole: "pink", oRole: "cyan"},
	{Unlockable: set("kitchen", "DUCK & TOASTER", arcade.Legendary, "Nobody knows why."),
		x:    []string{"...YY....", "..YYYk...", "..YYYYoo.", "...YY....", ".YYYYYY..", "YYYYYYYY.", ".YYYYYY..", "........."},
		o:    []string{".W.....W.", "..W...W..", "GGGGGGGGG", "GOOOGOOOG", "GGGGGGGGG", "GGGGGGGRG", "GGGGGGGGG", "k.......k"},
		xPal: map[rune]string{'Y': "yellow", 'k': "tire", 'o': "orange"}, oPal: map[rune]string{'W': "fg", 'G': "comment", 'O': "orange", 'R': "red", 'k': "tire"},
		xMini: "◆", oMini: "■", xRole: "yellow", oRole: "comment"},
	{Unlockable: set("feet", "SOCK & SANDAL", arcade.Legendary, "A fashion crime."),
		x:    []string{"...WWW...", "...RRR...", "...WWW...", "...WWW...", "...WWWW..", "..WWWWWW.", "..WWWWWW.", "........."},
		o:    []string{".........", "..B...B..", "..BB.BB..", "...BBB...", "OOOOOOOOO", "OOOOOOOOO", "BBBBBBBBB", "........."},
		xPal: map[rune]string{'W': "fg", 'R': "red"}, oPal: map[rune]string{'B': "orangeD", 'O': "orange"},
		xMini: "♠", oMini: "▬", xRole: "fg", oRole: "orange"},
}

func boardStyle(id, name string, tier arcade.Tier, blurb string) arcade.Unlockable {
	return arcade.Unlockable{ID: id, Name: name, Tier: tier, Kind: kindBoard, Blurb: blurb}
}

var boardStyles = []arcade.Unlockable{
	boardStyle("classic", "CLASSIC", arcade.Starter, "Lines, nothing else."),
	boardStyle("notebook", "NOTEBOOK", arcade.Starter, "The back page in maths."),
	boardStyle("chalk", "CHALKBOARD", arcade.Common, "Squeak."),
	boardStyle("wood", "WOODEN", arcade.Common, "Carved by grandad."),
	boardStyle("neon", "NEON", arcade.Rare, "Open all night."),
	boardStyle("picnic", "PICNIC BLANKET", arcade.Rare, "Mind the ants."),
	boardStyle("vineyard", "VINEYARD TRELLIS", arcade.Legendary, "Grown, not drawn."),
}

// Board ids are prefixed in the unlockables list so they can't clash with
// a piece set's ("classic" is both).
const boardPrefix = "board:"

// unlockables is every piece set and board, for the arcade.
func unlockables() []arcade.Unlockable {
	var out []arcade.Unlockable
	for _, s := range pieceSets {
		out = append(out, s.Unlockable)
	}
	for _, b := range boardStyles {
		b.ID = boardPrefix + b.ID
		out = append(out, b)
	}
	return out
}

// pieces looks up a piece set (the first starter if id is unknown).
func pieces(id string) *pieceSet {
	for i := range pieceSets {
		if pieceSets[i].ID == id {
			return &pieceSets[i]
		}
	}
	return &pieceSets[0]
}

// ghostOf maps every coloured rune of a picture to "ghost": a locked
// silhouette.
func ghostOf(rows []string) map[rune]string {
	pal := map[rune]string{}
	for _, r := range rows {
		for _, ch := range r {
			if ch != '.' {
				pal[ch] = "ghost"
			}
		}
	}
	return pal
}

// sprite is seat's piece ("X" or "O") in set s: its rows and palette.
func (s *pieceSet) sprite(piece string, ghost bool) ([]string, map[rune]string) {
	rows, pal := s.x, s.xPal
	if piece == "O" {
		rows, pal = s.o, s.oPal
	}
	if ghost {
		pal = ghostOf(rows)
	}
	return rows, pal
}

func (s *pieceSet) mini(piece string, ghost bool) (string, string) {
	ch, role := s.xMini, s.xRole
	if piece == "O" {
		ch, role = s.oMini, s.oRole
	}
	if ghost {
		role = "ghost"
	}
	return ch, role
}
