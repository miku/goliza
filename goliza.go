// ----------------------------------------------------------------------
//  eliza.go
//
//  A Go port of the classic ELIZA chatbot
//  Originally by Joe Strout (Python version)
//  Ported to Go for goliza project
// ----------------------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const banner = `
░█▀▀░█░░░▀█▀░▀▀█░█▀█
░█▀▀░█░░░░█░░▄▀░░█▀█
░▀▀▀░▀▀▀░▀▀▀░▀▀▀░▀░▀
`

type Eliza struct {
	keys   []*regexp.Regexp
	values [][]string
}

func NewEliza() *Eliza {
	keys := make([]*regexp.Regexp, len(gPats))
	values := make([][]string, len(gPats))

	for i, pat := range gPats {
		// Compile with case-insensitive flag
		keys[i] = regexp.MustCompile("(?i)" + pat.pattern)
		values[i] = pat.responses
	}

	return &Eliza{
		keys:   keys,
		values: values,
	}
}

// translate replaces words in str according to the provided dictionary
func (e *Eliza) translate(str string, dict map[string]string) string {
	words := strings.Fields(strings.ToLower(str))
	for i, word := range words {
		if val, ok := dict[word]; ok {
			words[i] = val
		}
	}
	return strings.Join(words, " ")
}

// respond takes a string and returns a response based on pattern matching
func (e *Eliza) respond(str string) string {
	// Find a match among keys
	for i, key := range e.keys {
		match := key.FindStringSubmatch(str)
		if match != nil {
			// Found a match - choose randomly from available responses
			resp := e.values[i][rand.Intn(len(e.values[i]))]

			// Stuff in reflected text where indicated by %1, %2, etc.
			for {
				pos := strings.Index(resp, "%")
				if pos == -1 {
					break
				}

				if pos+1 < len(resp) {
					// Extract the group number
					numStr := string(resp[pos+1])
					if num, err := strconv.Atoi(numStr); err == nil && num < len(match) {
						// Replace %N with the reflected captured group
						reflected := e.translate(match[num], gReflections)
						resp = resp[:pos] + reflected + resp[pos+2:]
					} else {
						break
					}
				} else {
					break
				}
			}

			// Fix munged punctuation at the end
			if strings.HasSuffix(resp, "?.") {
				resp = resp[:len(resp)-2] + "."
			}
			if strings.HasSuffix(resp, "??") {
				resp = resp[:len(resp)-2] + "?"
			}

			return resp
		}
	}

	return ""
}

// ----------------------------------------------------------------------
// gReflections: translation table to convert things you say into
// things the computer says back, e.g. "I am" --> "you are"
// ----------------------------------------------------------------------
var gReflections = map[string]string{
	"am":     "are",
	"was":    "were",
	"i":      "you",
	"i'd":    "you would",
	"i've":   "you have",
	"i'll":   "you will",
	"my":     "your",
	"are":    "am",
	"you've": "I have",
	"you'll": "I will",
	"your":   "my",
	"yours":  "mine",
	"you":    "me",
	"me":     "you",
}

// ----------------------------------------------------------------------
// gPats: the main response table. Each element contains a regexp
// pattern and a list of possible responses, with group-macros
// labelled as %1, %2, etc.
// ----------------------------------------------------------------------
type pattern struct {
	pattern   string
	responses []string
}

var gPats = []pattern{
	{
		pattern: `I need (.*)`,
		responses: []string{
			"Why do you need %1?",
			"Would it really help you to get %1?",
			"Are you sure you need %1?",
		},
	},
	{
		pattern: `Why don\'?t you ([^\?]*)\??`,
		responses: []string{
			"Do you really think I don't %1?",
			"Perhaps eventually I will %1.",
			"Do you really want me to %1?",
		},
	},
	{
		pattern: `Why can\'?t I ([^\?]*)\??`,
		responses: []string{
			"Do you think you should be able to %1?",
			"If you could %1, what would you do?",
			"I don't know -- why can't you %1?",
			"Have you really tried?",
		},
	},
	{
		pattern: `I can\'?t (.*)`,
		responses: []string{
			"How do you know you can't %1?",
			"Perhaps you could %1 if you tried.",
			"What would it take for you to %1?",
		},
	},
	{
		pattern: `I am (.*)`,
		responses: []string{
			"Did you come to me because you are %1?",
			"How long have you been %1?",
			"How do you feel about being %1?",
		},
	},
	{
		pattern: `I\'?m (.*)`,
		responses: []string{
			"How does being %1 make you feel?",
			"Do you enjoy being %1?",
			"Why do you tell me you're %1?",
			"Why do you think you're %1?",
		},
	},
	{
		pattern: `Are you ([^\?]*)\??`,
		responses: []string{
			"Why does it matter whether I am %1?",
			"Would you prefer it if I were not %1?",
			"Perhaps you believe I am %1.",
			"I may be %1 -- what do you think?",
		},
	},
	{
		pattern: `What (.*)`,
		responses: []string{
			"Why do you ask?",
			"How would an answer to that help you?",
			"What do you think?",
		},
	},
	{
		pattern: `How (.*)`,
		responses: []string{
			"How do you suppose?",
			"Perhaps you can answer your own question.",
			"What is it you're really asking?",
		},
	},
	{
		pattern: `Because (.*)`,
		responses: []string{
			"Is that the real reason?",
			"What other reasons come to mind?",
			"Does that reason apply to anything else?",
			"If %1, what else must be true?",
		},
	},
	{
		pattern: `(.*) sorry (.*)`,
		responses: []string{
			"There are many times when no apology is needed.",
			"What feelings do you have when you apologize?",
		},
	},
	{
		pattern: `Hello(.*)`,
		responses: []string{
			"Hello... I'm glad you could drop by today.",
			"Hi there... how are you today?",
			"Hello, how are you feeling today?",
		},
	},
	{
		pattern: `I think (.*)`,
		responses: []string{
			"Do you doubt %1?",
			"Do you really think so?",
			"But you're not sure %1?",
		},
	},
	{
		pattern: `(.*) friend (.*)`,
		responses: []string{
			"Tell me more about your friends.",
			"When you think of a friend, what comes to mind?",
			"Why don't you tell me about a childhood friend?",
		},
	},
	{
		pattern: `Yes`,
		responses: []string{
			"You seem quite sure.",
			"OK, but can you elaborate a bit?",
		},
	},
	{
		pattern: `(.*) computer(.*)`,
		responses: []string{
			"Are you really talking about me?",
			"Does it seem strange to talk to a computer?",
			"How do computers make you feel?",
			"Do you feel threatened by computers?",
		},
	},
	{
		pattern: `Is it (.*)`,
		responses: []string{
			"Do you think it is %1?",
			"Perhaps it's %1 -- what do you think?",
			"If it were %1, what would you do?",
			"It could well be that %1.",
		},
	},
	{
		pattern: `It is (.*)`,
		responses: []string{
			"You seem very certain.",
			"If I told you that it probably isn't %1, what would you feel?",
		},
	},
	{
		pattern: `Can you ([^\?]*)\??`,
		responses: []string{
			"What makes you think I can't %1?",
			"If I could %1, then what?",
			"Why do you ask if I can %1?",
		},
	},
	{
		pattern: `Can I ([^\?]*)\??`,
		responses: []string{
			"Perhaps you don't want to %1.",
			"Do you want to be able to %1?",
			"If you could %1, would you?",
		},
	},
	{
		pattern: `You are (.*)`,
		responses: []string{
			"Why do you think I am %1?",
			"Does it please you to think that I'm %1?",
			"Perhaps you would like me to be %1.",
			"Perhaps you're really talking about yourself?",
		},
	},
	{
		pattern: `You\'?re (.*)`,
		responses: []string{
			"Why do you say I am %1?",
			"Why do you think I am %1?",
			"Are we talking about you, or me?",
		},
	},
	{
		pattern: `I don\'?t (.*)`,
		responses: []string{
			"Don't you really %1?",
			"Why don't you %1?",
			"Do you want to %1?",
		},
	},
	{
		pattern: `I feel (.*)`,
		responses: []string{
			"Good, tell me more about these feelings.",
			"Do you often feel %1?",
			"When do you usually feel %1?",
			"When you feel %1, what do you do?",
		},
	},
	{
		pattern: `I have (.*)`,
		responses: []string{
			"Why do you tell me that you've %1?",
			"Have you really %1?",
			"Now that you have %1, what will you do next?",
		},
	},
	{
		pattern: `I would (.*)`,
		responses: []string{
			"Could you explain why you would %1?",
			"Why would you %1?",
			"Who else knows that you would %1?",
		},
	},
	{
		pattern: `Is there (.*)`,
		responses: []string{
			"Do you think there is %1?",
			"It's likely that there is %1.",
			"Would you like there to be %1?",
		},
	},
	{
		pattern: `My (.*)`,
		responses: []string{
			"I see, your %1.",
			"Why do you say that your %1?",
			"When your %1, how do you feel?",
		},
	},
	{
		pattern: `You (.*)`,
		responses: []string{
			"We should be discussing you, not me.",
			"Why do you say that about me?",
			"Why do you care whether I %1?",
		},
	},
	{
		pattern: `Why (.*)`,
		responses: []string{
			"Why don't you tell me the reason why %1?",
			"Why do you think %1?",
		},
	},
	{
		pattern: `I want (.*)`,
		responses: []string{
			"What would it mean to you if you got %1?",
			"Why do you want %1?",
			"What would you do if you got %1?",
			"If you got %1, then what would you do?",
		},
	},
	{
		pattern: `(.*) mother(.*)`,
		responses: []string{
			"Tell me more about your mother.",
			"What was your relationship with your mother like?",
			"How do you feel about your mother?",
			"How does this relate to your feelings today?",
			"Good family relations are important.",
		},
	},
	{
		pattern: `(.*) father(.*)`,
		responses: []string{
			"Tell me more about your father.",
			"How did your father make you feel?",
			"How do you feel about your father?",
			"Does your relationship with your father relate to your feelings today?",
			"Do you have trouble showing affection with your family?",
		},
	},
	{
		pattern: `(.*) child(.*)`,
		responses: []string{
			"Did you have close friends as a child?",
			"What is your favorite childhood memory?",
			"Do you remember any dreams or nightmares from childhood?",
			"Did the other children sometimes tease you?",
			"How do you think your childhood experiences relate to your feelings today?",
		},
	},
	{
		pattern: `(.*)\?`,
		responses: []string{
			"Why do you ask that?",
			"Please consider whether you can answer your own question.",
			"Perhaps the answer lies within yourself?",
			"Why don't you tell me?",
		},
	},
	{
		pattern: `quit`,
		responses: []string{
			"Thank you for talking with me.",
			"Good-bye.",
			"Thank you, that will be $150. Have a good day!",
		},
	},
	{
		pattern: `(.*)`,
		responses: []string{
			"Please tell me more.",
			"Let's change focus a bit... Tell me about your family.",
			"Can you elaborate on that?",
			"Why do you say that %1?",
			"I see.",
			"Very interesting.",
			"%1.",
			"I see. And what does that tell you?",
			"How does that make you feel?",
			"How do you feel when you say that?",
		},
	},
}

// ----------------------------------------------------------------------
//
//	commandInterface: main REPL loop
//
// ----------------------------------------------------------------------
func commandInterface() {
	fmt.Println(banner)
	fmt.Println()
	fmt.Println("paper (1966): https://web.stanford.edu/class/cs124/p36-weizenabaum.pdf")
	fmt.Println("go implementation (2026): https://git.sr.ht/~miku/goliza")
	fmt.Println()
	fmt.Println("Therapist")
	fmt.Println("---------")
	fmt.Println("Talk to the program by typing in plain English, using normal upper-")
	fmt.Println("and lower-case letters and punctuation. Enter \"quit\" when done.")
	fmt.Println(strings.Repeat("=", 72))
	fmt.Println("Hello. How are you feeling today?")

	therapist := NewEliza()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		s := scanner.Text()
		if s == "quit" {
			fmt.Println(therapist.respond(s))
			break
		}

		// Remove trailing punctuation
		s = strings.TrimRight(s, "!.")

		response := therapist.respond(s)
		if response != "" {
			fmt.Println(response)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	commandInterface()
}
