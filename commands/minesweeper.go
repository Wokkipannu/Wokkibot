package commands

import (
	"fmt"
	"math/rand"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var minesweeperCommand = discord.SlashCommandCreate{
	Name:        "minesweeper",
	Description: "Start a mine sweeper game",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "width",
			Description: "Width of the board",
			Required:    false,
		},
		discord.ApplicationCommandOptionInt{
			Name:        "height",
			Description: "Height of the board",
			Required:    false,
		},
		discord.ApplicationCommandOptionInt{
			Name:        "mines",
			Description: "Number of mines",
			Required:    false,
		},
	},
}

const (
	EMOJI_MINE    = "💣"
	EMOJI_FLAG    = "🚩"
	EMOJI_COVERED = "⬜"
	EMOJI_EMPTY   = "⬛"
	EMOJI_ONE     = "1️⃣"
	EMOJI_TWO     = "2️⃣"
	EMOJI_THREE   = "3️⃣"
	EMOJI_FOUR    = "4️⃣"
	EMOJI_FIVE    = "5️⃣"
	EMOJI_SIX     = "6️⃣"
	EMOJI_SEVEN   = "7️⃣"
	EMOJI_EIGHT   = "8️⃣"
	EMOJI_UP      = "⬆️"
	EMOJI_DOWN    = "⬇️"
	EMOJI_LEFT    = "⬅️"
	EMOJI_RIGHT   = "➡️"
	EMOJI_CURSOR  = "🧑"
)

const (
	DEFAULT_WIDTH  = 8
	DEFAULT_HEIGHT = 8
	DEFAULT_MINES  = 10
)

type Cell struct {
	isMine     bool
	isRevealed bool
	isFlagged  bool
	mineCount  int
}

type Board struct {
	width    int
	height   int
	mines    int
	cells    [][]Cell
	cursorX  int
	cursorY  int
	gameOver bool
	authorID string
}

func newBoard(width, height, mines int, authorID string) *Board {
	if width <= 0 {
		width = DEFAULT_WIDTH
	}
	if height <= 0 {
		height = DEFAULT_HEIGHT
	}
	if mines <= 0 {
		mines = DEFAULT_MINES
	}

	maxMines := (width * height) - 1
	if mines > maxMines {
		mines = maxMines
	}

	board := &Board{
		width:    width,
		height:   height,
		mines:    mines,
		cells:    make([][]Cell, height),
		cursorX:  0,
		cursorY:  0,
		gameOver: false,
		authorID: authorID,
	}

	for i := range board.cells {
		board.cells[i] = make([]Cell, width)
	}

	minesPlaced := 0
	for minesPlaced < mines {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if !board.cells[y][x].isMine {
			board.cells[y][x].isMine = true
			minesPlaced++
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !board.cells[y][x].isMine {
				count := 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						newY, newX := y+dy, x+dx
						if newY >= 0 && newY < height && newX >= 0 && newX < width {
							if board.cells[newY][newX].isMine {
								count++
							}
						}
					}
				}
				board.cells[y][x].mineCount = count
			}
		}
	}

	return board
}

func (b *Board) String() string {
	var result string
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			cell := b.cells[y][x]
			cellEmoji := ""

			if x == b.cursorX && y == b.cursorY {
				cellEmoji = EMOJI_CURSOR
			} else if cell.isFlagged {
				cellEmoji = EMOJI_FLAG
			} else if !cell.isRevealed {
				cellEmoji = EMOJI_COVERED
			} else if cell.isMine {
				cellEmoji = EMOJI_MINE
			} else if cell.mineCount == 0 {
				cellEmoji = EMOJI_EMPTY
			} else {
				switch cell.mineCount {
				case 1:
					cellEmoji = EMOJI_ONE
				case 2:
					cellEmoji = EMOJI_TWO
				case 3:
					cellEmoji = EMOJI_THREE
				case 4:
					cellEmoji = EMOJI_FOUR
				case 5:
					cellEmoji = EMOJI_FIVE
				case 6:
					cellEmoji = EMOJI_SIX
				case 7:
					cellEmoji = EMOJI_SEVEN
				case 8:
					cellEmoji = EMOJI_EIGHT
				}
			}
			result += cellEmoji
		}
		result += "\n"
	}
	return result
}

func UpdateSweeperMessage(e discord.User, board *Board, gameOver bool) discord.MessageUpdate {
	embedBuilder := discord.NewEmbedBuilder().
		SetTitlef("%s's Minesweeper Game", e.Username).
		AddField("Width", fmt.Sprintf("%d", board.width), true).
		AddField("Height", fmt.Sprintf("%d", board.height), true).
		AddField("Mines", fmt.Sprintf("%d", board.mines), true).
		AddField("", board.String(), false).
		SetColor(utils.RGBToInteger(255, 215, 0))

	if gameOver {
		embedBuilder.SetTitle(fmt.Sprintf("%s's Minesweeper Game - Game Over!", e.Username))
	}

	builder := discord.NewMessageUpdateBuilder().
		SetEmbeds(embedBuilder.Build())

	if gameOver {
		builder.ClearContainerComponents()
	}

	return builder.Build()
}

func HandleMinesweeper(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		width := DEFAULT_WIDTH
		height := DEFAULT_HEIGHT
		mines := DEFAULT_MINES

		data := e.SlashCommandInteractionData()

		if w, ok := data.OptInt("width"); ok {
			width = w
		}

		if h, ok := data.OptInt("height"); ok {
			height = h
		}

		if m, ok := data.OptInt("mines"); ok {
			mines = m
		}

		if width > 20 || height > 20 || mines > (width*height)-1 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Width and height must be less than 20 and mines must be less than (width*height)-1").Build())
		}

		board := newBoard(width, height, mines, e.User().ID.String())
		b.Games[e.User().ID] = board

		embedBuilder := discord.NewEmbedBuilder().
			SetTitlef("%s's Minesweeper Game", e.User().Username).
			AddField("Width", fmt.Sprintf("%d", board.width), true).
			AddField("Height", fmt.Sprintf("%d", board.height), true).
			AddField("Mines", fmt.Sprintf("%d", board.mines), true).
			AddField("", board.String(), false).
			SetColor(utils.RGBToInteger(255, 215, 0))

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(embedBuilder.Build()).
			AddActionRow(
				discord.NewPrimaryButton("Up", "minesweeper/up").WithEmoji(discord.ComponentEmoji{Name: EMOJI_UP}),
				discord.NewPrimaryButton("Down", "minesweeper/down").WithEmoji(discord.ComponentEmoji{Name: EMOJI_DOWN}),
				discord.NewPrimaryButton("Left", "minesweeper/left").WithEmoji(discord.ComponentEmoji{Name: EMOJI_LEFT}),
				discord.NewPrimaryButton("Right", "minesweeper/right").WithEmoji(discord.ComponentEmoji{Name: EMOJI_RIGHT}),
			).
			AddActionRow(
				discord.NewPrimaryButton("Flag", "minesweeper/flag").WithEmoji(discord.ComponentEmoji{Name: EMOJI_FLAG}),
				discord.NewPrimaryButton("Reveal", "minesweeper/reveal").WithEmoji(discord.ComponentEmoji{Name: EMOJI_COVERED}),
			).
			Build())
	}
}

func HandleMinesweeperFlagAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	cell := &board.cells[board.cursorY][board.cursorX]
	if !cell.isRevealed {
		cell.isFlagged = !cell.isFlagged
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}

func HandleMinesweeperRevealAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	cell := &board.cells[board.cursorY][board.cursorX]
	if !cell.isFlagged && !cell.isRevealed {
		cell.isRevealed = true
		if cell.isMine {
			board.gameOver = true
			for y := 0; y < board.height; y++ {
				for x := 0; x < board.width; x++ {
					if board.cells[y][x].isMine {
						board.cells[y][x].isRevealed = true
					}
				}
			}
			return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, true))
		} else if cell.mineCount == 0 {
			board.revealEmptyAdjacent(board.cursorX, board.cursorY)
		}
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}

func (b *Board) revealEmptyAdjacent(x, y int) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			newX, newY := x+dx, y+dy
			if newX >= 0 && newX < b.width && newY >= 0 && newY < b.height {
				cell := &b.cells[newY][newX]
				if !cell.isRevealed && !cell.isFlagged {
					cell.isRevealed = true
					if cell.mineCount == 0 {
						b.revealEmptyAdjacent(newX, newY)
					}
				}
			}
		}
	}
}

func HandleMinesweeperUpAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	if board.cursorY > 0 {
		board.cursorY--
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}

func HandleMinesweeperDownAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	if board.cursorY < board.height-1 {
		board.cursorY++
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}

func HandleMinesweeperLeftAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	if board.cursorX > 0 {
		board.cursorX--
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}

func HandleMinesweeperRightAction(b *wokkibot.Wokkibot, e *handler.ComponentEvent) error {
	board, ok := b.Games[e.User().ID].(*Board)
	if !ok || board.gameOver || board.authorID != e.User().ID.String() {
		return nil
	}

	if board.cursorX < board.width-1 {
		board.cursorX++
	}

	return e.UpdateMessage(UpdateSweeperMessage(e.User(), board, false))
}
