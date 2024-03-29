/*
  Annoyance Discord Bot by IMcPwn.
  For the latest code and contact information visit: http://imcpwn.com

  Copyright (c) 2016, IMcPwn 
  All rights reserved.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions are met:

  * Redistributions of source code must retain the above copyright notice, this
    list of conditions and the following disclaimer.

  * Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

  * Neither the name of Annoyance nor the names of its
    contributors may be used to endorse or promote products derived from
    this software without specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
  SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
  CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
  OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

package main

import (
    "fmt"
    "io/ioutil"
    "flag"
    "math/rand"

    "github.com/bwmarrin/discordgo"
    "github.com/bwmarrin/dgvoice"
)

// The folder that contains the MP3s used
var FOLDER *string

func main() {
    TOKEN := flag.String("t", "", "Discord authentication token")
    FOLDER = flag.String("f", "", "Folder that contains the MP3s to play")
    flag.Parse()
	
    if *FOLDER == "" {
    	flag.Usage()
        fmt.Println("-f option is required")
        return
    }
    if *TOKEN == "" {
    	flag.Usage()
        fmt.Println("-t option is required")
        return
    }

    dg, err := discordgo.New(*TOKEN)
    if err != nil {
        flag.Usage()
        fmt.Println(err)
        return
    }

    dg.AddHandler(VoiceStateUpdate)
    dg.AddHandler(messageCreate)

    // Open the websocket and begin listening.
    err = dg.Open()
    if err != nil {
        flag.Usage()
        fmt.Println(err)
        return
    }

    // Make sure we're logged in successfully
    prefix, err := dg.User("@me")
    if err != nil {
        flag.Usage()
        fmt.Println(err)
        return
    }
    fmt.Println("Logged in as " + prefix.Username)

    // Set Discord status to "away"
    dg.UpdateStatus(1, "")

    fmt.Println("Welcome to Annoyance! Press enter to quit.")
    var input string
    fmt.Scanln(&input)
    return
}

// This function is called whenever there is a voice state update.
// This function is responsible for playing the MP3s.
// TODO: Only play MP3s when a user joins a channel.
// This may be possible by caching all the previous voice states.
// var VoiceStateCache map[string]*discordgo.VoiceState
func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
    fmt.Println("[*] VoiceStateUpdate Called")
    if v.ChannelID == "" {
        fmt.Println("[X] Invalid channel")
        return
    }
    if len(s.VoiceConnections) != 0 {
        fmt.Println("[X] Already speaking")
        return
    }
    // 1/20 chance of not ignoring the call 
    // This is so the bot is not triggered every time there is a voice update.
    if rand.Intn(20) != 1 {
        fmt.Println("[X] Ignoring call")
        return
    }
    fmt.Println("[*] Responding to call")

    fmt.Println("[*] Joining Channel ID #" + v.ChannelID)
    // Join the server unmuted and deafened
    dgv, err := s.ChannelVoiceJoin(v.GuildID, v.ChannelID, false, true)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Start loop and attempt to play all files in the given folder
    fmt.Println("[*] Reading Folder: ", *FOLDER)
    files, err := ioutil.ReadDir(*FOLDER)
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, f := range files {
        fmt.Println("[*] PlayAudioFile:", f.Name())
        // Say we're "playing" the name of the audio file
        s.UpdateStatus(0, f.Name())
        dgvoice.PlayAudioFile(dgv, fmt.Sprintf("%s/%s", *FOLDER, f.Name()))
    }
    // Set Discord status to away
    s.UpdateStatus(1, "")
    dgv.Disconnect()
    dgv.Close()
}

// This function will be called every time a new message is created 
// on any channel that the autenticated user has access to.
// This function is responsible for responding to @mentions.
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
            fmt.Println("[*] Mentioned. Handling commands.")
            help := "Hello. I'm a bot. I may or may not follow you into channels. "
            _, err = s.ChannelMessageSend(m.ChannelID, "@" + m.Author.Username + m.Author.Discriminator + " " + help)
            if err != nil {
                fmt.Println(err)
                return
            }
        }
}
