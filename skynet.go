/* Skynet Discord Chat Bot by IMcPwn.

 * Copyright 2016 IMcPwn 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.

 * For the latest code and contact information visit: http://imcpwn.com
 */

package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "math/rand"

    "github.com/bwmarrin/discordgo"
    "github.com/bwmarrin/dgvoice"
)

func main() {

    TOKEN := os.Getenv("TOKEN")
    /* Alternatively you can use username/password instead of a token
    USER := os.Getenv("USER")
    PASS := os.Getenv("PASS")
    dg, err := discordogo.New(USER, PASS)
    */
    dg, err := discordgo.New(TOKEN)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Register messageCreate as a callback for the messageCreate events.
    dg.AddHandler(messageCreate)

    // Open the websocket and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Make sure we're logged in successfully
    prefix, err := dg.User("@me")
    if err != nil {
        fmt.Println(err)
        fmt.Println("Make sure TOKEN is defined and valid\n" +
        "Windows example: set TOKEN=change_to_token\nLinux example: export TOKEN=change_to_token")
        return
    }
    fmt.Println("Logged in as " + prefix.Username)

    fmt.Println("Welcome to Skynet! Press enter to quit.")
    var input string
    fmt.Scanln(&input)
    return
}

func handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
    GUILD := os.Getenv("GUILD")
    fileName := "test.mp3"
    // Connect to voice channel.
    // NOTE: Setting mute to false, deaf to true.
    dgv, err := s.ChannelVoiceJoin(GUILD, m.ChannelID, false, true)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Say we're "playing" the file name
    s.UpdateStatus(0, fileName)
    dgv.PlayAudioFile(dgv, fileName)
    // TODO: defer?
    dgv.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated user has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
        if len(m.Mentions) < 1 {
            return
        }
        prefix, err := s.User("@me")
        if err != nil {
            fmt.Println(err)
            return
        }
        if m.Mentions[0].ID == prefix.ID  {
            fmt.Println("Mentioned. Handling commands.")
            handleCommands(s, m)
        }
}
