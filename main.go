package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var zalgo_up []string
var zalgo_down []string
var zalgo_mid []string

func init() {
	// data set of leet unicode chars
	//---------------------------------------------------

	//those go UP
	zalgo_up = []string{
		"\u030d" /*     ̍     */, "\u030e" /*     ̎     */, "\u0304" /*     ̄     */, "\u0305", /*     ̅     */
		"\u033f" /*     ̿     */, "\u0311" /*     ̑     */, "\u0306" /*     ̆     */, "\u0310", /*     ̐     */
		"\u0352" /*     ͒     */, "\u0357" /*     ͗     */, "\u0351" /*     ͑     */, "\u0307", /*     ̇     */
		"\u0308" /*     ̈     */, "\u030a" /*     ̊     */, "\u0342" /*     ͂     */, "\u0343", /*     ̓     */
		"\u0344" /*     ̈́     */, "\u034a" /*     ͊     */, "\u034b" /*     ͋     */, "\u034c", /*     ͌     */
		"\u0303" /*     ̃     */, "\u0302" /*     ̂     */, "\u030c" /*     ̌     */, "\u0350", /*     ͐     */
		"\u0300" /*     ̀     */, "\u0301" /*     ́     */, "\u030b" /*     ̋     */, "\u030f", /*     ̏     */
		"\u0312" /*     ̒     */, "\u0313" /*     ̓     */, "\u0314" /*     ̔     */, "\u033d", /*     ̽     */
		"\u0309" /*     ̉     */, "\u0363" /*     ͣ     */, "\u0364" /*     ͤ     */, "\u0365", /*     ͥ     */
		"\u0366" /*     ͦ     */, "\u0367" /*     ͧ     */, "\u0368" /*     ͨ     */, "\u0369", /*     ͩ     */
		"\u036a" /*     ͪ     */, "\u036b" /*     ͫ     */, "\u036c" /*     ͬ     */, "\u036d", /*     ͭ     */
		"\u036e" /*     ͮ     */, "\u036f" /*     ͯ     */, "\u033e" /*     ̾     */, "\u035b", /*     ͛     */
		"\u0346" /*     ͆     */, "\u031a", /*     ̚     */
	}

	//those go DOWN
	zalgo_down = []string{
		"\u0316" /*     ̖     */, "\u0317" /*     ̗     */, "\u0318" /*     ̘     */, "\u0319", /*     ̙     */
		"\u031c" /*     ̜     */, "\u031d" /*     ̝     */, "\u031e" /*     ̞     */, "\u031f", /*     ̟     */
		"\u0320" /*     ̠     */, "\u0324" /*     ̤     */, "\u0325" /*     ̥     */, "\u0326", /*     ̦     */
		"\u0329" /*     ̩     */, "\u032a" /*     ̪     */, "\u032b" /*     ̫     */, "\u032c", /*     ̬     */
		"\u032d" /*     ̭     */, "\u032e" /*     ̮     */, "\u032f" /*     ̯     */, "\u0330", /*     ̰     */
		"\u0331" /*     ̱     */, "\u0332" /*     ̲     */, "\u0333" /*     ̳     */, "\u0339", /*     ̹     */
		"\u033a" /*     ̺     */, "\u033b" /*     ̻     */, "\u033c" /*     ̼     */, "\u0345", /*     ͅ     */
		"\u0347" /*     ͇     */, "\u0348" /*     ͈     */, "\u0349" /*     ͉     */, "\u034d", /*     ͍     */
		"\u034e" /*     ͎     */, "\u0353" /*     ͓     */, "\u0354" /*     ͔     */, "\u0355", /*     ͕     */
		"\u0356" /*     ͖     */, "\u0359" /*     ͙     */, "\u035a" /*     ͚     */, "\u0323", /*     ̣     */
	}

	//those always stay in the middle
	zalgo_mid = []string{
		"\u0315" /*     ̕     */, "\u031b" /*     ̛     */, "\u0340" /*     ̀     */, "\u0341", /*     ́     */
		"\u0358" /*     ͘     */, "\u0321" /*     ̡     */, "\u0322" /*     ̢     */, "\u0327", /*     ̧     */
		"\u0328" /*     ̨     */, "\u0334" /*     ̴     */, "\u0335" /*     ̵     */, "\u0336", /*     ̶     */
		"\u034f" /*     ͏     */, "\u035c" /*     ͜     */, "\u035d" /*     ͝     */, "\u035e", /*     ͞     */
		"\u035f" /*     ͟     */, "\u0360" /*     ͠     */, "\u0362" /*     ͢     */, "\u0338", /*     ̸     */
		"\u0337" /*     ̷     */, "\u0361" /*     ͡     */, "\u0489", /*     ҉_     */
	}
}

// rand funcs
//---------------------------------------------------

//gets a random char from a zalgo char table
func rand_zalgo(array []string) string {
	index := rand.Intn(len(array))
	return array[index]
}

//lookup char to know if its a zalgo char or not
func is_zalgo_char(c string) bool {
	for _, cc := range zalgo_up {
		if c == cc {
			return true
		}
	}

	for _, cc := range zalgo_mid {
		if c == cc {
			return true
		}
	}

	for _, cc := range zalgo_down {
		if c == cc {
			return true
		}
	}

	return false
}

func main() {
	ctx := context.Background()

	var size string
	var up, middle, down bool
	middle, down = true, true

	cmd := &cobra.Command{
		Use:  "eeemo [text]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if size != "mini" && size != "normal" && size != "maxi" {
				return errors.New("invalid size")
			}

			var input string

			file := os.Stdin
			fi, err := file.Stat()
			if err != nil {
				return err
			}
			fsize := fi.Size()
			if fsize == 0 {
				if len(args) > 0 {
					input = args[0]
				}
			}

			if fsize > 0 {
				// scanner := bufio.NewScanner(os.Stdin)
				// for scanner.Scan() {
				// 	fmt.Println(scanner.Text()) // Println will add back the final '\n'
				// }
				// if err := scanner.Err(); err != nil {
				// 	fmt.Fprintln(os.Stderr, "reading standard input:", err)
				// }
				b, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}

				input = string(b)
			}

			if input == "" {
				return errors.New("couldn't read input")
			}

			fmt.Print(
				HECOMES(input, size, up, middle, down),
			)

			return nil
		},
	}

	cmd.Flags().StringVarP(&size, "size", "", "mini", "'mini', 'normal' or 'maxi'")
	cmd.Flags().BoolVarP(&up, "up", "", true, "")
	cmd.Flags().BoolVarP(&middle, "middle", "", true, "")
	cmd.Flags().BoolVarP(&down, "down", "", false, "")

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// main shit
//---------------------------------------------------
func HECOMES(txt string, size string, up bool, middle bool, down bool) string {
	newtxt := ""

	for _, c := range strings.Split(txt, "") {
		if is_zalgo_char(c) {
			continue
		}

		var numUp, numMid, numDown int

		newtxt += string(c)

		if size == "mini" {
			numUp = rand.Intn(8)
			numMid = rand.Intn(2)
			numDown = rand.Intn(8)
		} else if size == "normal" {
			numUp = rand.Intn(16)/2 + 1
			numMid = rand.Intn(6) / 2
			numDown = rand.Intn(16)/2 + 1
		} else { // maxi
			numUp = rand.Intn(64)/4 + 3
			numMid = rand.Intn(16)/4 + 1
			numDown = rand.Intn(64)/4 + 3
		}

		if up {
			for j := 0; j < numUp; j++ {
				newtxt += rand_zalgo(zalgo_up)
			}
		}

		if middle {
			for j := 0; j < numMid; j++ {
				newtxt += rand_zalgo(zalgo_mid)
			}
		}

		if down {
			for j := 0; j < numDown; j++ {
				newtxt += rand_zalgo(zalgo_down)
			}
		}
	}

	// for (var i = 0; i < txt.length; i++) {
	//   if (is_zalgo_char(txt.substr(i, 1))) continue;

	//   var num_up;
	//   var num_mid;
	//   var num_down;

	//add the normal character
	//   newtxt += txt.substr(i, 1);

	//   //options
	//   if (document.getElementById("zalgo_opt_mini").checked) {
	// 	num_up = rand(8);
	// 	num_mid = rand(2);
	// 	num_down = rand(8);
	//   } else if (document.getElementById("zalgo_opt_normal").checked) {
	// 	num_up = rand(16) / 2 + 1;
	// 	num_mid = rand(6) / 2;
	// 	num_down = rand(16) / 2 + 1;
	//   } //maxi
	//   else {
	// 	num_up = rand(64) / 4 + 3;
	// 	num_mid = rand(16) / 4 + 1;
	// 	num_down = rand(64) / 4 + 3;
	//   }

	//   if (document.getElementById("zalgo_opt_up").checked)
	// 	for (var j = 0; j < num_up; j++) newtxt += rand_zalgo(zalgo_up);
	//   if (document.getElementById("zalgo_opt_mid").checked)
	// 	for (var j = 0; j < num_mid; j++) newtxt += rand_zalgo(zalgo_mid);
	//   if (document.getElementById("zalgo_opt_down").checked)
	// 	for (var j = 0; j < num_down; j++) newtxt += rand_zalgo(zalgo_down);
	// }

	//result is in nextxt, display that

	//remove all children of lulz_container
	// var container = document.getElementById("lulz_container");
	// while (container.childNodes.length)
	//   container.removeChild(container.childNodes[0]);

	// //build blocks for each line & create a <br />
	// var lines = newtxt.split("\n");
	// for (var i = 0; i < lines.length; i++) {
	//   var n = document.createElement("text");
	//   n.innerHTML = lines[i];
	//   container.appendChild(n);
	//   var nl = document.createElement("br");
	//   container.appendChild(nl);
	// }

	return newtxt
}
